package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	ModeWizard Mode = iota
	ModePromptInput
	ModeGenerating
	ModeDone
)

type Model struct {
	cfg         *config.Config
	currentStep Step
	cursor      int
	mode        Mode
	inputs      []textinput.Model
	promptInput textarea.Model
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

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(84)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFD1"))

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

	badgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFD1")).
			Italic(true)

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
	ta.Placeholder = "Enter your plain text architecture requirements here...\ne.g. 'Build a React frontend with Go microservices using Kafka event broker and Postgres DB'"
	ta.SetWidth(78)
	ta.SetHeight(8)

	return Model{
		cfg:         cfg,
		currentStep: Step0,
		cursor:      0,
		mode:        ModeWizard,
		inputs:      []textinput.Model{ti},
		promptInput: ta,
		statusMsg:   "Press 'Tab' / '←/→' to switch steps | '↑/↓' to navigate | 'Enter/Space' to cycle options",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
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
		return 3 // 0: Changelog, 1: ReleaseNotes, 2: Alerts, 3: [PROCEED TO ARCHITECTURE PROMPT INPUT]
	default:
		return 0
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.mode == ModePromptInput {
		var cmd tea.Cmd
		switch msg := msg.(tea.Msg).(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "esc":
				m.mode = ModeWizard
				m.statusMsg = "Returned to Wizard settings"
				return m, nil
			case "ctrl+s", "ctrl+d": // Submit Architecture Prompt
				m.mode = ModeGenerating
				promptText := m.promptInput.Value()
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

	// ModeWizard update logic
	switch msg := msg.(tea.Msg).(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

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
				// Transition to Plain Text Architecture Prompt Input Mode
				m.mode = ModePromptInput
				m.promptInput.Focus()
				m.statusMsg = "Type your architecture requirements in plain text. Press 'Ctrl+S' or 'Ctrl+D' when finished."
				return m, textarea.Blink
			}
			m.toggleCurrentField()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
			m.statusMsg = "Type your architecture requirements in plain text. Press 'Ctrl+S' or 'Ctrl+D' when finished."
		}
	}
}

func (m Model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD1")).Render("\n Thanks for using Agentic-Init (ainit)! See you soon.\n\n")
	}

	var sb strings.Builder

	sb.WriteString(titleStyle.Render("⚡ Agentic-Init (ainit) Harness TUI Engineering Tool"))
	sb.WriteString("\n\n")

	if m.mode == ModePromptInput {
		sb.WriteString(headerStyle.Render("📝 Step 5: Input Architecture Requirements (Plain Text)"))
		sb.WriteString("\n\n")
		sb.WriteString(m.promptInput.View())
		sb.WriteString("\n\n")
		sb.WriteString(boxStyle.Render(fmt.Sprintf(
			" Selected Config Summary:\n  • Provider: %s (%s)\n  • Arch: %s (%s)\n  • Git: %s\n  • Commit: %s",
			m.cfg.Step0.ProviderID, m.cfg.Step0.PrimaryModel,
			m.cfg.Step1.ArchitectureStyle, m.cfg.Step1.RepoStructure,
			m.cfg.Step2.GitProvider, m.cfg.Step3.CommitConvention,
		)))
		sb.WriteString("\n\n")
		sb.WriteString(statusStyle.Render(" Status: "))
		sb.WriteString(m.statusMsg)
		sb.WriteString("\n")
		return sb.String()
	}

	if m.mode == ModeDone {
		sb.WriteString(headerStyle.Render("🎉 Generation Pipeline Completed!"))
		sb.WriteString("\n\n")
		sb.WriteString(boxStyle.Render(m.genResult + "\n\nOutput files:\n • docs/ARCHITECTURE_SPEC.md\n • AGENTS.md\n • CLAUDE.md\n • .cursorrules\n • .github/copilot-instructions.md"))
		sb.WriteString("\n\n")
		sb.WriteString(statusStyle.Render(" Press Enter or 'q' to exit."))
		sb.WriteString("\n")
		return sb.String()
	}

	// Step Tabs Header
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

	// Main Box Content
	body := m.renderStepBody()
	sb.WriteString(boxStyle.Render(body))
	sb.WriteString("\n\n")

	// Footer Status Bar
	sb.WriteString(statusStyle.Render(" Status: "))
	sb.WriteString(m.statusMsg)
	sb.WriteString("\n")

	return sb.String()
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

	return fmt.Sprintf("%s%-20s %s\n", prefix, lblStyle.Render(label), valStyle.Render(value))
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

		sb.WriteString(headerStyle.Render("🤖 Step 0: OpenCode Multi-Provider & Model Catalog"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "AI Provider:", "["+provName+"]"))
		sb.WriteString(m.renderRow(1, "Primary Model:", "["+m.cfg.Step0.PrimaryModel+"]"))
		sb.WriteString(m.renderRow(2, "Auth Method:", "["+m.cfg.Step0.LicensingMode+"]"))
		sb.WriteString(m.renderRow(3, "Fallback Model:", "["+m.cfg.Step0.FallbackModel+"]"))
		sb.WriteString("\n")
		sb.WriteString(badgeStyle.Render("  [OpenCode Inspired Catalog: Anthropic, OpenAI, Gemini, DeepSeek, OpenRouter, Ollama]\n"))
		sb.WriteString(hintStyle.Render("  [Press Enter on AI Provider to switch provider and dynamically load models]"))

	case Step1:
		sb.WriteString(headerStyle.Render("🏗️ Step 1: Architecture Spec & Mermaid Generation"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Arch Style:", "["+m.cfg.Step1.ArchitectureStyle+"]"))
		sb.WriteString(m.renderRow(1, "Repo Structure:", "["+m.cfg.Step1.RepoStructure+"]"))
		sb.WriteString(m.renderRow(2, "Sequence Diagram:", fmt.Sprintf("[%v]", m.cfg.Step1.GenerateSequence)))
		sb.WriteString(m.renderRow(3, "GitOps Diagram:", fmt.Sprintf("[%v]", m.cfg.Step1.GenerateGitOps)))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Press Enter/Space to toggle MSA/Monolith or diagram specs]"))

	case Step2:
		sb.WriteString(headerStyle.Render("🔌 Step 2: MCP Tooling & Infrastructure Connections"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Git Provider:", "["+m.cfg.Step2.GitProvider+"] (Required)"))
		sb.WriteString(m.renderRow(1, "K8s Target:", "["+m.cfg.Step2.K8sTarget+"]"))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", unfocusedLabelStyle.Render("CI/CD Tools:"), unfocusedValueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.CICDTools))))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", unfocusedLabelStyle.Render("Doc Sync:"), unfocusedValueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.DocTools))))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", unfocusedLabelStyle.Render("Messenger Webhook:"), unfocusedValueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.MessengerAlerts))))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Press Enter/Space on Git/K8s to cycle provider targets]"))

	case Step3:
		sb.WriteString(headerStyle.Render("🛠️ Step 3: Harness TDD, Commit & PR Conventions"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Commit Convention:", "["+m.cfg.Step3.CommitConvention+"]"))
		sb.WriteString(m.renderRow(1, "PR Template Style:", "["+m.cfg.Step3.PRTemplateStyle+"]"))
		sb.WriteString(m.renderRow(2, "TDD First Mode:", fmt.Sprintf("[%v]", m.cfg.Step3.TDDMode)))
		sb.WriteString(m.renderRow(3, "Local Sandbox Test:", fmt.Sprintf("[%v]", m.cfg.Step3.LocalSandboxTest)))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Configures Conventional Commits, PR Template & TDD loop]"))

	case Step4:
		sb.WriteString(headerStyle.Render("🚀 Step 4: Release Pipeline & Submit Setup"))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderRow(0, "Auto Changelog:", fmt.Sprintf("[%v]", m.cfg.Step4.AutoChangelog)))
		sb.WriteString(m.renderRow(1, "Release Notes Sync:", fmt.Sprintf("[%v]", m.cfg.Step4.ReleaseNotesSync)))
		sb.WriteString(m.renderRow(2, "Deploy Alerts:", fmt.Sprintf("[%v]", m.cfg.Step4.DeployAlert)))
		sb.WriteString("\n")
		sb.WriteString(m.renderRow(3, "👉 PROCEED TO PLAIN TEXT ARCHITECTURE INPUT 👈", ""))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Select item 3 or press Enter to input plain text architecture requirements]"))
	}

	return sb.String()
}
