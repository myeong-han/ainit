package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/myeong-han/ainit/pkg/command"
	"github.com/myeong-han/ainit/pkg/config"
	"github.com/myeong-han/ainit/pkg/generator"
	"github.com/myeong-han/ainit/pkg/provider"
)

type Step int

const (
	Step0 Step = iota // AI Licensing & Provider
	Step1             // Architecture Spec
	Step2             // MCP Tooling
	Step3             // Harness TDD & Conventions
	Step4             // Release Pipeline
)

func (s Step) String() string {
	switch s {
	case Step0:
		return "0. AI Licensing"
	case Step1:
		return "1. Arch Spec"
	case Step2:
		return "2. MCP Tooling"
	case Step3:
		return "3. Harness TDD"
	case Step4:
		return "4. Release Pipeline"
	default:
		return ""
	}
}

type Mode int

const (
	ModePromptInput Mode = iota // Default Main Chatting Mode
	ModeWizard                  // Step Config Form Mode (via /set-confs)
	ModeGenerating
	ModeDone
)

type ChatMessage struct {
	Sender  string // "User" or "Agent"
	Content string
}

type Model struct {
	cfg         *config.Config
	cmdEngine   *command.CommandEngine
	currentStep Step
	cursor      int
	mode        Mode
	inputs      []textinput.Model
	promptInput textarea.Model
	chatHistory []ChatMessage
	width       int
	height      int
	statusMsg   string
	quitting    bool
	genResult   string
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFD1")).
			Background(lipgloss.Color("#2E1A47")).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8A8A8A")).
				Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFD1"))

	sidebarHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#0087D7")).
				Padding(0, 1)

	focusedLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FFD1"))

	unfocusedLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D1D1D1"))

	focusedValueStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF87D7")).
				Underline(true)

	unfocusedValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#BBBBBB"))

	userChatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFD1")).
			Bold(true)

	agentChatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00")).
			Bold(true)
)

func NewModel(cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-project"
	ti.SetValue(cfg.Step1.ProjectName)
	ti.Focus()

	ta := textarea.New()
	ta.Placeholder = "Type architecture requirements or Slash Commands (/set-confs, /gen-docs, /gen-codes, /help)..."
	ta.SetWidth(60)
	ta.SetHeight(4)
	ta.Focus()

	initialHistory := []ChatMessage{
		{
			Sender:  "Agent",
			Content: "Welcome to Agentic-Init (ainit)! Type your architecture requirements in plain text or use /set-confs to modify settings.",
		},
	}

	return Model{
		cfg:         cfg,
		cmdEngine:   command.NewCommandEngine(cfg),
		currentStep: Step0,
		cursor:      0,
		mode:        ModePromptInput, // Default Main View is Chat
		inputs:      []textinput.Model{ti},
		promptInput: ta,
		chatHistory: initialHistory,
		width:       100, // Default width fallback
		height:      28,  // Default height fallback
		statusMsg:   "Main Chatting View | Type /set-confs to open config form | Press Ctrl+S to submit",
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) getMaxCursorForStep() int {
	switch m.currentStep {
	case Step0:
		return 3
	case Step1:
		return 3
	case Step2:
		return 1
	case Step3:
		return 3
	case Step4:
		return 3
	default:
		return 0
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(tea.Msg).(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.mode == ModePromptInput {
		var cmd tea.Cmd
		switch msg := msg.(tea.Msg).(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			case "enter":
				val := strings.TrimSpace(m.promptInput.Value())
				if val == "" {
					return m, nil
				}

				m.chatHistory = append(m.chatHistory, ChatMessage{Sender: "User", Content: val})

				if val == "/set-confs" {
					m.mode = ModeWizard
					m.statusMsg = "Switched to Config Wizard Form. Press 'Tab' to navigate steps or 'Esc' to return to Chat."
					m.promptInput.Reset()
					return m, nil
				}

				if command.IsSlashCommand(val) {
					res, err := m.cmdEngine.Execute(val)
					if err != nil {
						m.chatHistory = append(m.chatHistory, ChatMessage{Sender: "Agent", Content: fmt.Sprintf("❌ Error: %v", err)})
						m.statusMsg = fmt.Sprintf("Slash Command Error: %v", err)
					} else {
						m.chatHistory = append(m.chatHistory, ChatMessage{Sender: "Agent", Content: res.Message})
						m.statusMsg = res.Message
						if res.Action == command.ActionGenDocs || res.Action == command.ActionGenCodes {
							m.mode = ModeDone
							m.genResult = res.Message
						}
					}
					m.promptInput.Reset()
					return m, nil
				}

				// Plain text message
				m.chatHistory = append(m.chatHistory, ChatMessage{
					Sender:  "Agent",
					Content: fmt.Sprintf("Received architecture prompt: '%s'. Press 'Ctrl+S' to generate docs & code.", val),
				})
				m.promptInput.Reset()
				return m, nil

			case "ctrl+s", "ctrl+d":
				m.mode = ModeGenerating
				promptText := m.promptInput.Value()
				if strings.TrimSpace(promptText) == "" && len(m.chatHistory) > 0 {
					promptText = m.chatHistory[len(m.chatHistory)-1].Content
				}
				if strings.TrimSpace(promptText) == "" {
					promptText = "Default MSA architecture with Go backend & React frontend"
				}

				err := generator.GenerateHarnessProject(".", m.cfg, promptText)
				if err != nil {
					m.genResult = fmt.Sprintf("❌ Error generating project: %v", err)
				} else {
					m.genResult = "🎉 Architecture Spec, Mermaid Diagrams, Git & Agent Rules Generated Successfully!"
				}
				m.mode = ModeDone
				return m, nil
			}
		}
		m.promptInput, cmd = m.promptInput.Update(msg)
		return m, cmd
	}

	if m.mode == ModeWizard {
		switch msg := msg.(tea.Msg).(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit

			case "esc":
				m.mode = ModePromptInput
				m.statusMsg = "Returned to Main Chatting View"
				return m, textarea.Blink

			case "right", "tab":
				if m.currentStep < Step4 {
					m.currentStep++
					m.cursor = 0
				}

			case "left", "shift+tab":
				if m.currentStep > Step0 {
					m.currentStep--
					m.cursor = 0
				}

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				maxCursor := m.getMaxCursorForStep()
				if m.cursor < maxCursor {
					m.cursor++
				}

			case "enter", " ":
				if m.currentStep == Step4 && m.cursor == 3 {
					m.mode = ModePromptInput
					m.promptInput.Focus()
					m.statusMsg = "Returned to Main Chatting View"
					return m, textarea.Blink
				}
				m.toggleCurrentField()
			}
		}
		return m, nil
	}

	if m.mode == ModeDone {
		switch msg := msg.(tea.Msg).(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q", "enter", "esc":
				m.quitting = true
				return m, tea.Quit
			}
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) toggleCurrentField() {
	switch m.currentStep {
	case Step0:
		providers := provider.GetAvailableProviders()
		switch m.cursor {
		case 0:
			currIdx := -1
			for i, p := range providers {
				if p.ID == m.cfg.Step0.ProviderID {
					currIdx = i
					break
				}
			}
			nextIdx := (currIdx + 1) % len(providers)
			nextProvider := providers[nextIdx]
			m.cfg.Step0.ProviderID = nextProvider.ID
			m.cfg.Step0.LicensingMode = nextProvider.DefaultAuth
			if len(nextProvider.Models) > 0 {
				m.cfg.Step0.PrimaryModel = nextProvider.Models[0].ID
			}
			m.statusMsg = fmt.Sprintf("Switched AI Provider to [%s] (%s)", nextProvider.Name, nextProvider.DefaultAuth)

		case 1:
			models := provider.GetModelsForProvider(m.cfg.Step0.ProviderID)
			if len(models) > 0 {
				currIdx := -1
				for i, mod := range models {
					if mod.ID == m.cfg.Step0.PrimaryModel {
						currIdx = i
						break
					}
				}
				nextIdx := (currIdx + 1) % len(models)
				m.cfg.Step0.PrimaryModel = models[nextIdx].ID
				m.statusMsg = fmt.Sprintf("Selected Primary Model: [%s]", models[nextIdx].Name)
			}

		case 2:
			switch m.cfg.Step0.LicensingMode {
			case "subscription":
				m.cfg.Step0.LicensingMode = "apikey"
			case "apikey":
				m.cfg.Step0.LicensingMode = "local"
			default:
				m.cfg.Step0.LicensingMode = "subscription"
			}
			m.statusMsg = fmt.Sprintf("Auth Mode changed to [%s]", m.cfg.Step0.LicensingMode)

		case 3:
			switch m.cfg.Step0.FallbackModel {
			case "gpt-4o-mini":
				m.cfg.Step0.FallbackModel = "claude-3-5-haiku"
			case "claude-3-5-haiku":
				m.cfg.Step0.FallbackModel = "gemini-1.5-flash"
			case "gemini-1.5-flash":
				m.cfg.Step0.FallbackModel = "none"
			default:
				m.cfg.Step0.FallbackModel = "gpt-4o-mini"
			}
			m.statusMsg = fmt.Sprintf("Fallback Model set to [%s]", m.cfg.Step0.FallbackModel)
		}

	case Step1:
		switch m.cursor {
		case 0:
			switch m.cfg.Step1.ArchitectureStyle {
			case "msa":
				m.cfg.Step1.ArchitectureStyle = "monolith"
			case "monolith":
				m.cfg.Step1.ArchitectureStyle = "eda"
			default:
				m.cfg.Step1.ArchitectureStyle = "msa"
			}
			m.statusMsg = fmt.Sprintf("Changed Arch Style to [%s]", m.cfg.Step1.ArchitectureStyle)

		case 1:
			if m.cfg.Step1.RepoStructure == "monorepo" {
				m.cfg.Step1.RepoStructure = "multirepo"
			} else {
				m.cfg.Step1.RepoStructure = "monorepo"
			}
			m.statusMsg = fmt.Sprintf("Changed Repo Layout to [%s]", m.cfg.Step1.RepoStructure)

		case 2:
			m.cfg.Step1.GenerateSequence = !m.cfg.Step1.GenerateSequence
			m.statusMsg = fmt.Sprintf("Toggled Sequence Diagram Generation: [%v]", m.cfg.Step1.GenerateSequence)

		case 3:
			m.cfg.Step1.GenerateGitOps = !m.cfg.Step1.GenerateGitOps
			m.statusMsg = fmt.Sprintf("Toggled GitOps Diagram Generation: [%v]", m.cfg.Step1.GenerateGitOps)
		}

	case Step2:
		switch m.cursor {
		case 0:
			if m.cfg.Step2.GitProvider == "github" {
				m.cfg.Step2.GitProvider = "bitbucket"
			} else {
				m.cfg.Step2.GitProvider = "github"
			}
			m.statusMsg = fmt.Sprintf("Selected Git Provider: [%s]", m.cfg.Step2.GitProvider)

		case 1:
			switch m.cfg.Step2.K8sTarget {
			case "local":
				m.cfg.Step2.K8sTarget = "remote"
			case "remote":
				m.cfg.Step2.K8sTarget = "none"
			default:
				m.cfg.Step2.K8sTarget = "local"
			}
			m.statusMsg = fmt.Sprintf("Selected K8s Target: [%s]", m.cfg.Step2.K8sTarget)
		}

	case Step3:
		switch m.cursor {
		case 0:
			switch m.cfg.Step3.CommitConvention {
			case "conventional":
				m.cfg.Step3.CommitConvention = "gitmoji"
			case "gitmoji":
				m.cfg.Step3.CommitConvention = "issue-prefix"
			case "issue-prefix":
				m.cfg.Step3.CommitConvention = "custom"
			default:
				m.cfg.Step3.CommitConvention = "conventional"
			}
			m.statusMsg = fmt.Sprintf("Changed Commit Convention to [%s]", m.cfg.Step3.CommitConvention)

		case 1:
			switch m.cfg.Step3.PRTemplateStyle {
			case "standard":
				m.cfg.Step3.PRTemplateStyle = "minimal"
			case "minimal":
				m.cfg.Step3.PRTemplateStyle = "jira"
			default:
				m.cfg.Step3.PRTemplateStyle = "standard"
			}
			m.statusMsg = fmt.Sprintf("Selected PR Template Style: [%s]", m.cfg.Step3.PRTemplateStyle)

		case 2:
			m.cfg.Step3.TDDMode = !m.cfg.Step3.TDDMode
			m.statusMsg = fmt.Sprintf("Toggled TDD First Mode: [%v]", m.cfg.Step3.TDDMode)

		case 3:
			m.cfg.Step3.LocalSandboxTest = !m.cfg.Step3.LocalSandboxTest
			m.statusMsg = fmt.Sprintf("Toggled Local Sandbox Build Check: [%v]", m.cfg.Step3.LocalSandboxTest)
		}

	case Step4:
		switch m.cursor {
		case 0:
			m.cfg.Step4.AutoChangelog = !m.cfg.Step4.AutoChangelog
			m.statusMsg = fmt.Sprintf("Toggled Auto Changelog: [%v]", m.cfg.Step4.AutoChangelog)

		case 1:
			m.cfg.Step4.ReleaseNotesSync = !m.cfg.Step4.ReleaseNotesSync
			m.statusMsg = fmt.Sprintf("Toggled Release Notes Sync: [%v]", m.cfg.Step4.ReleaseNotesSync)

		case 2:
			m.cfg.Step4.DeployAlert = !m.cfg.Step4.DeployAlert
			m.statusMsg = fmt.Sprintf("Toggled Deploy Webhook Alerts: [%v]", m.cfg.Step4.DeployAlert)

		case 3:
			m.mode = ModePromptInput
			m.promptInput.Focus()
			m.statusMsg = "Returned to Main Chatting View"
		}
	}
}

func (m Model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD1")).Render("\n Thanks for using Agentic-Init (ainit)! See you soon.\n\n")
	}

	var sb strings.Builder

	// Title
	sb.WriteString(titleStyle.Render("⚡ Agentic-Init (ainit) Harness TUI Engineering Tool"))
	sb.WriteString("\n\n")

	// Dynamic Layout Width & Height Calculation
	sidebarWidth := 30
	mainWidth := m.width - sidebarWidth - 6
	if mainWidth < 40 {
		mainWidth = 40
	}
	boxHeight := m.height - 8
	if boxHeight < 14 {
		boxHeight = 14
	}

	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(mainWidth).
		Height(boxHeight)

	rightSidebarBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFD1")).
		Padding(1, 1).
		Width(sidebarWidth).
		Height(boxHeight)

	// Step Tabs Header (Only shown when in Wizard Mode)
	if m.mode == ModeWizard {
		var tabs []string
		steps := []Step{Step0, Step1, Step2, Step3, Step4}
		for _, s := range steps {
			if s == m.currentStep {
				tabs = append(tabs, activeTabStyle.Render(" ▶ "+s.String()+" "))
			} else {
				tabs = append(tabs, inactiveTabStyle.Render("   "+s.String()+" "))
			}
		}
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, tabs...))
		sb.WriteString("\n\n")
	}

	// 2-Column Responsive Layout: Left Main Flex View + Right Fixed Sidebar Nav
	leftColumnView := m.renderLeftColumn(mainWidth - 4)
	rightSidebarView := m.renderRightSidebarNav()

	splitView := lipgloss.JoinHorizontal(lipgloss.Top, leftBox.Render(leftColumnView), rightSidebarBox.Render(rightSidebarView))
	sb.WriteString(splitView)
	sb.WriteString("\n\n")

	// Footer Status Bar
	sb.WriteString(statusStyle.Render(" Status: "))
	sb.WriteString(m.statusMsg)
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) renderLeftColumn(width int) string {
	if m.mode == ModePromptInput {
		var sb strings.Builder
		sb.WriteString(headerStyle.Render("💬 Main Architecture Chatting Session"))
		sb.WriteString("\n\n")

		// Render Chat History
		for _, msg := range m.chatHistory {
			if msg.Sender == "User" {
				sb.WriteString(userChatStyle.Render("👤 User: ") + msg.Content + "\n")
			} else {
				sb.WriteString(agentChatStyle.Render("🤖 Agent: ") + msg.Content + "\n")
			}
		}

		sb.WriteString("\n")
		m.promptInput.SetWidth(width - 2)
		sb.WriteString(m.promptInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(hintStyle.Render("[Type /set-confs for config form | Type /help for commands | Press Ctrl+S to Submit]"))
		return sb.String()
	}

	if m.mode == ModeDone {
		var sb strings.Builder
		sb.WriteString(headerStyle.Render("🎉 Generation Pipeline Completed!"))
		sb.WriteString("\n\n")
		sb.WriteString(m.genResult + "\n\nOutput files:\n • docs/ARCHITECTURE_SPEC.md\n • AGENTS.md\n • CLAUDE.md\n • .cursorrules\n • .github/copilot-instructions.md")
		return sb.String()
	}

	return m.renderStepBody()
}

func (m Model) renderRightSidebarNav() string {
	var sb strings.Builder
	sb.WriteString(sidebarHeaderStyle.Render("📊 CONFIG STATUS NAV"))
	sb.WriteString("\n\n")

	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFD1")).Render("Step 0: AI Licensing") + "\n")
	sb.WriteString(fmt.Sprintf("• Prov: %s\n", truncateStr(m.cfg.Step0.ProviderID, 14)))
	sb.WriteString(fmt.Sprintf("• Model: %s\n", truncateStr(m.cfg.Step0.PrimaryModel, 13)))
	sb.WriteString(fmt.Sprintf("• Auth: [%s]\n\n", m.cfg.Step0.LicensingMode))

	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFD1")).Render("Step 1: Arch Spec") + "\n")
	sb.WriteString(fmt.Sprintf("• Style: %s (%s)\n", m.cfg.Step1.ArchitectureStyle, m.cfg.Step1.RepoStructure))
	sb.WriteString(fmt.Sprintf("• Diag: Seq(%v) GitOps(%v)\n\n", m.cfg.Step1.GenerateSequence, m.cfg.Step1.GenerateGitOps))

	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFD1")).Render("Step 2: MCP Connections") + "\n")
	sb.WriteString(fmt.Sprintf("• Git: %s (READY)\n", m.cfg.Step2.GitProvider))
	sb.WriteString(fmt.Sprintf("• K8s: %s | ArgoCD: ON\n\n", m.cfg.Step2.K8sTarget))

	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFD1")).Render("Step 3: Harness & TDD") + "\n")
	sb.WriteString(fmt.Sprintf("• Commit: %s\n", m.cfg.Step3.CommitConvention))
	sb.WriteString(fmt.Sprintf("• TDD: %v | Sandbox: %v\n\n", m.cfg.Step3.TDDMode, m.cfg.Step3.LocalSandboxTest))

	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFD1")).Render("Step 4: Release Pipeline") + "\n")
	sb.WriteString(fmt.Sprintf("• SemVer: %s\n", m.cfg.Step4.VersioningStrategy))
	sb.WriteString(fmt.Sprintf("• Sync: Notion/Slack (%v)\n", m.cfg.Step4.ReleaseNotesSync))

	return sb.String()
}

func truncateStr(s string, max int) string {
	if len(s) > max {
		return s[:max-2] + ".."
	}
	return s
}

func (m Model) renderRow(index int, label string, value string) string {
	prefix := "  "
	lblStyle := unfocusedLabelStyle
	valStyle := unfocusedValueStyle

	if index == m.cursor {
		prefix = "▶ "
		lblStyle = focusedLabelStyle
		valStyle = focusedValueStyle
	}

	return fmt.Sprintf("%s%-18s %s\n", prefix, lblStyle.Render(label), valStyle.Render(value))
}

func (m Model) renderStepBody() string {
	var sb strings.Builder

	switch m.currentStep {
	case Step0:
		providerObj := provider.GetProviderByID(m.cfg.Step0.ProviderID)
		provName := m.cfg.Step0.ProviderID
		if providerObj != nil {
			provName = providerObj.Name
		}

		sb.WriteString(headerStyle.Render("🤖 Step 0: OpenCode AI Provider Catalog"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "AI Provider:", "["+provName+"]"))
		sb.WriteString(m.renderRow(1, "Primary Model:", "["+truncateStr(m.cfg.Step0.PrimaryModel, 16)+"]"))
		sb.WriteString(m.renderRow(2, "Auth Method:", "["+m.cfg.Step0.LicensingMode+"]"))
		sb.WriteString(m.renderRow(3, "Fallback Model:", "["+m.cfg.Step0.FallbackModel+"]"))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Press Enter to cycle OpenCode Providers]"))

	case Step1:
		sb.WriteString(headerStyle.Render("🏗️ Step 1: Architecture Spec & Mermaid"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Arch Style:", "["+m.cfg.Step1.ArchitectureStyle+"]"))
		sb.WriteString(m.renderRow(1, "Repo Structure:", "["+m.cfg.Step1.RepoStructure+"]"))
		sb.WriteString(m.renderRow(2, "Sequence Diag:", fmt.Sprintf("[%v]", m.cfg.Step1.GenerateSequence)))
		sb.WriteString(m.renderRow(3, "GitOps Diag:", fmt.Sprintf("[%v]", m.cfg.Step1.GenerateGitOps)))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Press Enter/Space to toggle specs]"))

	case Step2:
		sb.WriteString(headerStyle.Render("🔌 Step 2: MCP Tooling Connections"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Git Provider:", "["+m.cfg.Step2.GitProvider+"] (READY)"))
		sb.WriteString(m.renderRow(1, "K8s Target:", "["+m.cfg.Step2.K8sTarget+"]"))
		sb.WriteString(fmt.Sprintf("  %-18s %s\n", unfocusedLabelStyle.Render("CI/CD Tools:"), unfocusedValueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.CICDTools))))
		sb.WriteString(fmt.Sprintf("  %-18s %s\n", unfocusedLabelStyle.Render("Doc Sync:"), unfocusedValueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.DocTools))))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Press Enter on Git/K8s to toggle]"))

	case Step3:
		sb.WriteString(headerStyle.Render("🛠️ Step 3: Harness TDD & Conventions"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Commit Conv:", "["+m.cfg.Step3.CommitConvention+"]"))
		sb.WriteString(m.renderRow(1, "PR Template:", "["+m.cfg.Step3.PRTemplateStyle+"]"))
		sb.WriteString(m.renderRow(2, "TDD First Mode:", fmt.Sprintf("[%v]", m.cfg.Step3.TDDMode)))
		sb.WriteString(m.renderRow(3, "Local Sandbox:", fmt.Sprintf("[%v]", m.cfg.Step3.LocalSandboxTest)))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Configures Conventional Commits & TDD]"))

	case Step4:
		sb.WriteString(headerStyle.Render("🚀 Step 4: Release & Submit Setup"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Auto Changelog:", fmt.Sprintf("[%v]", m.cfg.Step4.AutoChangelog)))
		sb.WriteString(m.renderRow(1, "Release Sync:", fmt.Sprintf("[%v]", m.cfg.Step4.ReleaseNotesSync)))
		sb.WriteString(m.renderRow(2, "Deploy Alerts:", fmt.Sprintf("[%v]", m.cfg.Step4.DeployAlert)))
		sb.WriteString("\n")
		sb.WriteString(m.renderRow(3, "👉 RETURN TO MAIN CHATTING 👈", ""))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Press Enter to return to main chat]"))
	}

	return sb.String()
}
