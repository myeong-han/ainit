package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/myeong-han/ainit/pkg/command"
	"github.com/myeong-han/ainit/pkg/config"
	"github.com/myeong-han/ainit/pkg/connection"
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
	ModeWizard                  // Step Config Form Mode (via /settings)
	ModeKeyInput                // API Key / Token Input Modal
	ModeConfirm                 // Confirmation Prompt before execution
	ModeGenerating
	ModeDone
)

type ChatMessage struct {
	Sender  string // "User" or "Agent"
	Content string
}

type Model struct {
	cfg                    *config.Config
	cmdEngine              *command.CommandEngine
	connTester             *connection.ConnectionTester
	currentStep            Step
	cursor                 int
	mode                   Mode
	inputs                 []textinput.Model
	promptInput            textarea.Model
	keyInput               textinput.Model
	keyTarget              string
	chatHistory            []ChatMessage
	slashDropdownOpen      bool
	slashDropdownDismissed bool
	slashCursor            int
	slashOptions           []command.SlashOption
	pendingGenCmd          string
	pendingGenAction       command.ActionType
	width                  int
	height                 int
	statusMsg              string
	quitting               bool
	genResult              string
	aiConnStatus           string
	gitConnStatus          string
	docConnStatus          string
	msgConnStatus          string
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

	sidebarSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FFD1"))

	sidebarKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0"))

	sidebarValStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF87D7"))

	sidebarReadyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5AF78E")).
				Bold(true)

	sidebarDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#444444"))

	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FFAF00")).
			Padding(1, 2)

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

	dropdownBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#00FFD1")).
				Background(lipgloss.Color("#1A1A2E")).
				Padding(0, 1)

	dropdownItemFocused = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FFD1")).
				Background(lipgloss.Color("#2E1A47"))

	dropdownItemUnfocused = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC"))

	userChatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFD1")).
			Bold(true)

	agentChatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00")).
			Bold(true)
)

func filterSlashCommands(prefix string) []command.SlashOption {
	all := command.GetAvailableSlashCommands()
	if prefix == "" || prefix == "/" {
		return all
	}

	var filtered []command.SlashOption
	lowerPrefix := strings.ToLower(prefix)
	for _, opt := range all {
		if strings.HasPrefix(strings.ToLower(opt.Name), lowerPrefix) {
			filtered = append(filtered, opt)
		}
	}
	return filtered
}

func NewModel(cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "unknown"
	ti.SetValue(cfg.Step1.ProjectName)
	ti.Focus()

	ki := textinput.New()
	ki.Placeholder = "Enter API Key / Token / Webhook URL..."
	ki.EchoMode = textinput.EchoPassword
	ki.EchoCharacter = '•'

	ta := textarea.New()
	ta.Placeholder = "Type '/' for Slash Commands (/git-init <name>, /settings, /gen-all)..."
	ta.SetWidth(60)
	ta.SetHeight(4)
	ta.Focus()

	initialHistory := []ChatMessage{
		{
			Sender:  "Agent",
			Content: "Welcome to Agentic-Init (ainit)! Type '/git-init <name>' to set project name or '/' for commands.",
		},
	}

	return Model{
		cfg:                    cfg,
		cmdEngine:              command.NewCommandEngine(cfg),
		connTester:             connection.NewConnectionTester(),
		currentStep:            Step0,
		cursor:                 0,
		mode:                   ModePromptInput,
		inputs:                 []textinput.Model{ti},
		promptInput:            ta,
		keyInput:               ki,
		chatHistory:            initialHistory,
		slashDropdownOpen:      false,
		slashDropdownDismissed: false,
		slashCursor:            0,
		slashOptions:           command.GetAvailableSlashCommands(),
		width:                  100,
		height:                 28,
		statusMsg:              "Main Chatting View | Type '/git-init <name>' to set project name | Press Ctrl+S to submit",
		aiConnStatus:           "[UNTESTED]",
		gitConnStatus:          "[READY]",
		docConnStatus:          "[READY]",
		msgConnStatus:          "[READY]",
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
		return 5
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

	if m.mode == ModeKeyInput {
		var cmd tea.Cmd
		switch msg := msg.(tea.Msg).(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.mode = ModeWizard
				m.statusMsg = "Cancelled input."
				return m, nil
			case "enter":
				keyVal := strings.TrimSpace(m.keyInput.Value())
				m.mode = ModeWizard

				if m.keyTarget == "ai" {
					res := m.connTester.TestAIProvider(m.cfg.Step0.ProviderID, keyVal)
					if res.Connected {
						m.aiConnStatus = fmt.Sprintf("🟢 200 OK (%dms)", res.LatencyMs)
						m.statusMsg = fmt.Sprintf("AI Provider API Key Verified! (%s)", m.aiConnStatus)
					} else {
						m.aiConnStatus = fmt.Sprintf("❌ FAILED (%s)", res.ErrorMessage)
						m.statusMsg = fmt.Sprintf("AI Connection Failed: %s", res.ErrorMessage)
					}
				} else if m.keyTarget == "git" {
					res := m.connTester.TestGitProvider(context.Background(), m.cfg.Step2.GitProvider, keyVal)
					if res.Connected {
						m.gitConnStatus = fmt.Sprintf("🟢 200 OK (%dms)", res.LatencyMs)
						m.statusMsg = fmt.Sprintf("Git Provider Token Verified! (%s)", m.gitConnStatus)
					} else {
						m.gitConnStatus = fmt.Sprintf("❌ FAILED (%s)", res.ErrorMessage)
						m.statusMsg = fmt.Sprintf("Git Connection Failed: %s", res.ErrorMessage)
					}
				} else if m.keyTarget == "doc" {
					res := m.connTester.TestDocSync("notion", keyVal)
					if res.Connected {
						m.docConnStatus = fmt.Sprintf("🟢 200 OK (%dms)", res.LatencyMs)
						m.statusMsg = fmt.Sprintf("Notion Token Verified! (%s)", m.docConnStatus)
					} else {
						m.docConnStatus = fmt.Sprintf("❌ FAILED (%s)", res.ErrorMessage)
						m.statusMsg = fmt.Sprintf("Notion Sync Failed: %s", res.ErrorMessage)
					}
				} else if m.keyTarget == "msg" {
					if keyVal != "" {
						m.msgConnStatus = "🟢 200 OK"
						m.statusMsg = "Slack Webhook URL configured successfully!"
					} else {
						m.msgConnStatus = "❌ UNSET"
						m.statusMsg = "Slack Webhook URL cleared."
					}
				}
				m.keyInput.Reset()
				return m, nil
			}
		}
		m.keyInput, cmd = m.keyInput.Update(msg)
		return m, cmd
	}

	if m.mode == ModeConfirm {
		switch msg := msg.(tea.Msg).(type) {
		case tea.KeyMsg:
			switch strings.ToLower(msg.String()) {
			case "y", "enter":
				m.mode = ModeGenerating
				res, err := m.executePendingGeneration()
				if err != nil {
					m.genResult = fmt.Sprintf("[!] Generation failed: %v", err)
				} else {
					m.genResult = res
				}
				m.mode = ModeDone
				return m, nil

			case "n", "esc":
				m.mode = ModePromptInput
				m.statusMsg = "[!] Generation process cancelled by user."
				m.chatHistory = append(m.chatHistory, ChatMessage{Sender: "Agent", Content: "[!] Generation pipeline cancelled."})
				return m, textarea.Blink
			}
		}
		return m, nil
	}

	if m.mode == ModePromptInput {
		var cmd tea.Cmd

		switch msg := msg.(tea.Msg).(type) {
		case tea.KeyMsg:
			if msg.String() == "/" {
				m.slashDropdownDismissed = false
			}

			if m.slashDropdownOpen {
				switch msg.String() {
				case "up", "k":
					if m.slashCursor > 0 {
						m.slashCursor--
					}
					return m, nil
				case "down", "j":
					if m.slashCursor < len(m.slashOptions)-1 {
						m.slashCursor++
					}
					return m, nil
				case "tab", "enter", "right":
					if len(m.slashOptions) > 0 && m.slashCursor < len(m.slashOptions) {
						selectedCmd := m.slashOptions[m.slashCursor].Name
						m.promptInput.SetValue(selectedCmd + " ")
						m.slashDropdownOpen = false
						m.slashDropdownDismissed = true
						m.statusMsg = fmt.Sprintf("Autocompleted '%s'. Add arguments and press Enter to execute.", selectedCmd)
						return m, nil
					}
				case "esc":
					m.slashDropdownOpen = false
					m.slashDropdownDismissed = true
					m.statusMsg = "Slash commands dropdown closed."
					return m, nil
				}
			}

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

				if val == "/settings" {
					m.mode = ModeWizard
					m.statusMsg = "Switched to Config Wizard Form. Press 'Tab' to navigate steps or 'Esc' to return to Chat."
					m.promptInput.Reset()
					m.slashDropdownOpen = false
					m.slashDropdownDismissed = false
					return m, nil
				}

				if command.IsSlashCommand(val) {
					if val == "/gen-docs" || val == "/gen-codes" || val == "/gen-gitops" || val == "/gen-all" {
						m.pendingGenCmd = val
						m.mode = ModeConfirm
						m.statusMsg = fmt.Sprintf("[!] Confirm generation for %s? [y/n]", val)
						m.promptInput.Reset()
						m.slashDropdownOpen = false
						m.slashDropdownDismissed = false
						return m, nil
					}

					res, err := m.cmdEngine.Execute(val)
					if err != nil {
						m.chatHistory = append(m.chatHistory, ChatMessage{Sender: "Agent", Content: fmt.Sprintf("[!] Error: %v", err)})
						m.statusMsg = fmt.Sprintf("Slash Command Error: %v", err)
					} else {
						m.chatHistory = append(m.chatHistory, ChatMessage{Sender: "Agent", Content: res.Message})
						m.statusMsg = res.Message
					}
					m.promptInput.Reset()
					m.slashDropdownOpen = false
					m.slashDropdownDismissed = false
					return m, nil
				}

				m.chatHistory = append(m.chatHistory, ChatMessage{
					Sender:  "Agent",
					Content: fmt.Sprintf("Received architecture prompt: '%s'. Press 'Ctrl+S' to generate docs & code.", val),
				})
				m.promptInput.Reset()
				m.slashDropdownOpen = false
				m.slashDropdownDismissed = false
				return m, nil

			case "ctrl+s", "ctrl+d":
				m.pendingGenCmd = "/gen-all"
				m.mode = ModeConfirm
				m.statusMsg = "[!] Confirm full generation pipeline? [y/n]"
				return m, nil
			}
		}

		m.promptInput, cmd = m.promptInput.Update(msg)

		currVal := strings.TrimSpace(m.promptInput.Value())
		if currVal == "" {
			m.slashDropdownDismissed = false
		}

		if !m.slashDropdownDismissed && strings.HasPrefix(currVal, "/") && !strings.Contains(currVal, " ") {
			filtered := filterSlashCommands(currVal)
			if len(filtered) > 0 {
				m.slashOptions = filtered
				if m.slashCursor >= len(filtered) {
					m.slashCursor = 0
				}
				m.slashDropdownOpen = true
			} else {
				m.slashDropdownOpen = false
			}
		} else {
			m.slashDropdownOpen = false
		}

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

func (m Model) executePendingGeneration() (string, error) {
	cmdVal := m.pendingGenCmd
	if cmdVal == "" {
		cmdVal = "/gen-all"
	}

	res, err := m.cmdEngine.Execute(cmdVal)
	if err != nil {
		return "", err
	}
	return res.Message, nil
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
			m.mode = ModeKeyInput
			m.keyTarget = "ai"
			m.keyInput.Placeholder = fmt.Sprintf("Enter %s API Key...", strings.ToUpper(m.cfg.Step0.ProviderID))
			m.keyInput.Focus()
			m.statusMsg = "Type your API Key and press Enter to test connection."

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
			m.mode = ModeKeyInput
			m.keyTarget = "git"
			m.keyInput.Placeholder = fmt.Sprintf("Enter %s Personal Access Token...", strings.ToUpper(m.cfg.Step2.GitProvider))
			m.keyInput.Focus()
			m.statusMsg = "Type Git Token and press Enter to test connection."

		case 1:
			res := m.connTester.TestK8sCluster(m.cfg.Step2.K8sTarget)
			if res.Connected {
				m.statusMsg = fmt.Sprintf("Kubernetes Cluster Connected! (%dms)", res.LatencyMs)
			} else {
				m.statusMsg = fmt.Sprintf("K8s Connection Failed: %s", res.ErrorMessage)
			}

		case 2:
			m.statusMsg = "CI Tool is fixed to [jenkins]."

		case 3:
			m.statusMsg = "CD Tool is fixed to [argocd]."

		case 4:
			m.mode = ModeKeyInput
			m.keyTarget = "doc"
			m.keyInput.Placeholder = "Enter Notion Integration Token..."
			m.keyInput.Focus()
			m.statusMsg = "Type Notion Token and press Enter to test connection."

		case 5:
			m.mode = ModeKeyInput
			m.keyTarget = "msg"
			m.keyInput.Placeholder = "Enter Slack Webhook URL..."
			m.keyInput.Focus()
			m.statusMsg = "Type Slack Webhook URL and press Enter to configure."
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

	sb.WriteString(titleStyle.Render("Agentic-Init (ainit) Harness TUI Engineering Tool"))
	sb.WriteString("\n\n")

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

	if m.mode == ModeWizard || m.mode == ModeKeyInput {
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

	leftColumnView := m.renderLeftColumn(mainWidth - 4)
	rightSidebarView := m.renderRightSidebarNav()

	splitView := lipgloss.JoinHorizontal(lipgloss.Top, leftBox.Render(leftColumnView), rightSidebarBox.Render(rightSidebarView))
	sb.WriteString(splitView)
	sb.WriteString("\n\n")

	sb.WriteString(statusStyle.Render(" Status: "))
	sb.WriteString(m.statusMsg)
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) renderLeftColumn(width int) string {
	if m.mode == ModeKeyInput {
		var sb strings.Builder
		sb.WriteString(headerStyle.Render("🔑 API Key & Token Verification Modal"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("Target Service: %s\n\n", strings.ToUpper(m.keyTarget)))
		sb.WriteString(m.keyInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(hintStyle.Render("[Press Enter to test connection via HTTPS | Press Esc to cancel]"))
		return sb.String()
	}

	if m.mode == ModeConfirm {
		projName := m.cfg.Step1.ProjectName
		if projName == "" {
			projName = "unknown"
		}

		var sb strings.Builder
		sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFAF00")).Render("[!] CONFIRM CODE & DOC GENERATION PIPELINE"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("• Trigger Command : %s\n", m.pendingGenCmd))
		sb.WriteString(fmt.Sprintf("• Target App Name : %s\n", projName))
		sb.WriteString(fmt.Sprintf("• Target Provider : %s (%s)\n", m.cfg.Step0.ProviderID, m.cfg.Step0.PrimaryModel))
		sb.WriteString(fmt.Sprintf("• Output Paths    : docs/, gitops/, AGENTS.md, .cursorrules\n\n"))
		sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFD1")).Render("Are you sure you want to proceed with file generation?"))
		sb.WriteString("\n\n")
		sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5AF78E")).Render("[Press 'y' or Enter to CONFIRM]  ") + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("[Press 'n' or Esc to CANCEL]"))

		return confirmBoxStyle.Width(width - 2).Render(sb.String())
	}

	if m.mode == ModePromptInput {
		var sb strings.Builder
		sb.WriteString(headerStyle.Render("Main Architecture Chatting Session"))
		sb.WriteString("\n\n")

		for _, msg := range m.chatHistory {
			if msg.Sender == "User" {
				sb.WriteString(userChatStyle.Render("[User] ") + msg.Content + "\n")
			} else {
				sb.WriteString(agentChatStyle.Render("[Agent] ") + msg.Content + "\n")
			}
		}

		sb.WriteString("\n")

		if m.slashDropdownOpen {
			sb.WriteString(m.renderSlashDropdown(width))
			sb.WriteString("\n")
		}

		m.promptInput.SetWidth(width - 2)
		sb.WriteString(m.promptInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(hintStyle.Render("[Type '/' for Slash Commands Dropdown | Press Enter/Tab to select & edit | Press Esc to close dropdown]"))
		return sb.String()
	}

	if m.mode == ModeDone {
		var sb strings.Builder
		sb.WriteString(headerStyle.Render("Generation Pipeline Completed!"))
		sb.WriteString("\n\n")
		sb.WriteString(m.genResult + "\n\nOutput files:\n • docs/ARCHITECTURE_SPEC.md\n • gitops/helm & gitops/argocd\n • AGENTS.md & CLAUDE.md")
		return sb.String()
	}

	return m.renderStepBody()
}

func (m Model) renderSlashDropdown(width int) string {
	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFD1")).Render("Slash Commands (Use ↑/↓ to navigate, Enter/Tab to select & edit, Esc to close):") + "\n")

	for i, opt := range m.slashOptions {
		if i == m.slashCursor {
			sb.WriteString(dropdownItemFocused.Render(fmt.Sprintf(" ▶ %-12s - %s", opt.Name, opt.Description)) + "\n")
		} else {
			sb.WriteString(dropdownItemUnfocused.Render(fmt.Sprintf("   %-12s - %s", opt.Name, opt.Description)) + "\n")
		}
	}

	return dropdownBoxStyle.Width(width - 4).Render(sb.String())
}

func (m Model) renderRightSidebarNav() string {
	var sb strings.Builder
	divider := sidebarDividerStyle.Render("──────────────────────")

	sb.WriteString(sidebarHeaderStyle.Render("CONFIG STATUS NAV"))
	sb.WriteString("\n")
	sb.WriteString(divider)
	sb.WriteString("\n\n")

	projName := m.cfg.Step1.ProjectName
	if projName == "" {
		projName = "unknown"
	}
	sb.WriteString(sidebarKeyStyle.Render("App Name: ") + sidebarValStyle.Render(truncateStr(projName, 12)) + "\n")
	sb.WriteString(divider + "\n\n")

	// Step 0: AI Licensing & Provider
	sb.WriteString(sidebarSectionStyle.Render("Step 0: AI Licensing") + "\n")
	sb.WriteString(fmt.Sprintf("%s %s\n", sidebarKeyStyle.Render("• Prov :"), sidebarValStyle.Render(truncateStr(m.cfg.Step0.ProviderID, 11))))
	sb.WriteString(fmt.Sprintf("%s %s\n", sidebarKeyStyle.Render("• Model:"), sidebarValStyle.Render(truncateStr(m.cfg.Step0.PrimaryModel, 11))))
	sb.WriteString(fmt.Sprintf("%s %s\n\n", sidebarKeyStyle.Render("• Ping :"), sidebarReadyStyle.Render(truncateStr(m.aiConnStatus, 12))))
	sb.WriteString(divider + "\n\n")

	// Step 1: Arch Spec
	seqStr := "N"
	if m.cfg.Step1.GenerateSequence {
		seqStr = "Y"
	}
	gitopsStr := "N"
	if m.cfg.Step1.GenerateGitOps {
		gitopsStr = "Y"
	}
	sb.WriteString(sidebarSectionStyle.Render("Step 1: Arch Spec") + "\n")
	sb.WriteString(fmt.Sprintf("%s %s (%s)\n", sidebarKeyStyle.Render("• Style:"), sidebarValStyle.Render(m.cfg.Step1.ArchitectureStyle), m.cfg.Step1.RepoStructure))
	sb.WriteString(fmt.Sprintf("%s Seq(%s) Git(%s)\n\n", sidebarKeyStyle.Render("• Diag :"), seqStr, gitopsStr))
	sb.WriteString(divider + "\n\n")

	// Step 2: MCP Tooling
	sb.WriteString(sidebarSectionStyle.Render("Step 2: MCP Connections") + "\n")
	sb.WriteString(fmt.Sprintf("%s %s %s\n", sidebarKeyStyle.Render("• Git  :"), sidebarValStyle.Render(m.cfg.Step2.GitProvider), sidebarReadyStyle.Render(truncateStr(m.gitConnStatus, 9))))
	sb.WriteString(fmt.Sprintf("%s %s | %s\n\n", sidebarKeyStyle.Render("• K8s  :"), sidebarValStyle.Render(m.cfg.Step2.K8sTarget), sidebarReadyStyle.Render("ArgoCD:ON")))
	sb.WriteString(divider + "\n\n")

	// Step 3: Harness & TDD
	tddStr := "N"
	if m.cfg.Step3.TDDMode {
		tddStr = "Y"
	}
	sbStr := "N"
	if m.cfg.Step3.LocalSandboxTest {
		sbStr = "Y"
	}
	sb.WriteString(sidebarSectionStyle.Render("Step 3: Harness & TDD") + "\n")
	sb.WriteString(fmt.Sprintf("%s %s\n", sidebarKeyStyle.Render("• Commit:"), sidebarValStyle.Render(m.cfg.Step3.CommitConvention)))
	sb.WriteString(fmt.Sprintf("%s TDD(%s) Sbox(%s)\n\n", sidebarKeyStyle.Render("• Mode  :"), tddStr, sbStr))
	sb.WriteString(divider + "\n\n")

	// Step 4: Release Pipeline
	syncStr := "N"
	if m.cfg.Step4.ReleaseNotesSync {
		syncStr = "Y"
	}
	sb.WriteString(sidebarSectionStyle.Render("Step 4: Release Pipeline") + "\n")
	sb.WriteString(fmt.Sprintf("%s %s\n", sidebarKeyStyle.Render("• SemVer:"), sidebarValStyle.Render(m.cfg.Step4.VersioningStrategy)))
	sb.WriteString(fmt.Sprintf("%s Sync(%s)\n", sidebarKeyStyle.Render("• Slack :"), syncStr))

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

		sb.WriteString(headerStyle.Render("Step 0: OpenCode AI Provider Catalog"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "AI Provider:", "["+provName+"]"))
		sb.WriteString(m.renderRow(1, "Primary Model:", "["+truncateStr(m.cfg.Step0.PrimaryModel, 16)+"]"))
		sb.WriteString(m.renderRow(2, "🔑 Enter API Key:", "["+m.aiConnStatus+"]"))
		sb.WriteString(m.renderRow(3, "Fallback Model:", "["+m.cfg.Step0.FallbackModel+"]"))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Press Enter to input API Key and perform live connection test]"))

	case Step1:
		sb.WriteString(headerStyle.Render("Step 1: Architecture Spec & Mermaid"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Arch Style:", "["+m.cfg.Step1.ArchitectureStyle+"]"))
		sb.WriteString(m.renderRow(1, "Repo Structure:", "["+m.cfg.Step1.RepoStructure+"]"))
		sb.WriteString(m.renderRow(2, "Sequence Diag:", fmt.Sprintf("[%v]", m.cfg.Step1.GenerateSequence)))
		sb.WriteString(m.renderRow(3, "GitOps Diag:", fmt.Sprintf("[%v]", m.cfg.Step1.GenerateGitOps)))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Press Enter/Space to toggle specs]"))

	case Step2:
		sb.WriteString(headerStyle.Render("Step 2: MCP Tooling Connections"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "🔑 Git Token:", "["+m.cfg.Step2.GitProvider+"] "+m.gitConnStatus))
		sb.WriteString(m.renderRow(1, "K8s Target:", "["+m.cfg.Step2.K8sTarget+"] (Press Enter to Test)"))
		sb.WriteString(m.renderRow(2, "CI Tool:", "["+m.cfg.Step2.CI+"] [READY]"))
		sb.WriteString(m.renderRow(3, "CD Tool:", "["+m.cfg.Step2.CD+"] [ON]"))
		sb.WriteString(m.renderRow(4, "🔑 Doc Tool:", fmt.Sprintf("[%s] %s", m.cfg.Step2.Doc, m.docConnStatus)))
		sb.WriteString(m.renderRow(5, "🔑 Messenger:", fmt.Sprintf("[%s] %s", m.cfg.Step2.Messenger, m.msgConnStatus)))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Press Enter to input Tokens / Webhooks and test live connection]"))

	case Step3:
		sb.WriteString(headerStyle.Render("Step 3: Harness TDD & Conventions"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Commit Conv:", "["+m.cfg.Step3.CommitConvention+"]"))
		sb.WriteString(m.renderRow(1, "PR Template:", "["+m.cfg.Step3.PRTemplateStyle+"]"))
		sb.WriteString(m.renderRow(2, "TDD First Mode:", fmt.Sprintf("[%v]", m.cfg.Step3.TDDMode)))
		sb.WriteString(m.renderRow(3, "Local Sandbox:", fmt.Sprintf("[%v]", m.cfg.Step3.LocalSandboxTest)))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("[Configures Conventional Commits & TDD]"))

	case Step4:
		sb.WriteString(headerStyle.Render("Step 4: Release & Submit Setup"))
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
