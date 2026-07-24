package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/myeong-han/ainit/pkg/config"
)

// Step represents the current active step in the TUI wizard
type Step int

const (
	Step0 Step = iota // AI Licensing
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

// Model defines the Bubble Tea TUI application state
type Model struct {
	cfg         *config.Config
	currentStep Step
	cursor      int
	inputs      []textinput.Model
	width       int
	height      int
	statusMsg   string
	quitting    bool
}

// Lipgloss Style Definitions for Clean & Modern Aesthetic
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
			Width(78)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFD1"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D1D1"))

	valueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF87D7"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00")).
			Bold(true)
)

// NewModel initializes the TUI model with the given configuration
func NewModel(cfg *config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-project"
	ti.SetValue(cfg.Step1.ProjectName)
	ti.Focus()

	return Model{
		cfg:         cfg,
		currentStep: Step0,
		cursor:      0,
		inputs:      []textinput.Model{ti},
		statusMsg:   "Press 'Tab' or '←/→' to navigate steps | 'Ctrl+C' or 'q' to quit",
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.cursor++

		case "enter":
			m.statusMsg = fmt.Sprintf("Confirmed settings for %s", m.currentStep.String())
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFD1")).Render("\n Thanks for using Agentic-Init (ainit)! See you soon.\n\n")
	}

	var sb strings.Builder

	// Header Title
	sb.WriteString(titleStyle.Render("⚡ Agentic-Init (ainit) Harness TUI Engineering Tool"))
	sb.WriteString("\n\n")

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

func (m Model) renderStepBody() string {
	var sb strings.Builder

	switch m.currentStep {
	case Step0:
		sb.WriteString(headerStyle.Render("🤖 Step 0: AI Licensing & Model Setup"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Licensing Mode:"), valueStyle.Render("["+m.cfg.Step0.LicensingMode+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Primary Model:"), valueStyle.Render("["+m.cfg.Step0.PrimaryModel+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Fallback Model:"), valueStyle.Render("["+m.cfg.Step0.FallbackModel+"]")))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Use Up/Down to navigate fields, Enter to toggle options]"))

	case Step1:
		sb.WriteString(headerStyle.Render("🏗️ Step 1: Architecture Spec & Mermaid Generation"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Project Name:"), valueStyle.Render(m.cfg.Step1.ProjectName)))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Arch Style:"), valueStyle.Render("["+m.cfg.Step1.ArchitectureStyle+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Repo Structure:"), valueStyle.Render("["+m.cfg.Step1.RepoStructure+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Mermaid Diagrams:"),
			valueStyle.Render(fmt.Sprintf("Sequence(%v), Traffic(%v), GitOps(%v)",
				m.cfg.Step1.GenerateSequence, m.cfg.Step1.GenerateTraffic, m.cfg.Step1.GenerateGitOps))))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Generates Sequence, Ingress Traffic, GitOps & ERD specs]"))

	case Step2:
		sb.WriteString(headerStyle.Render("🔌 Step 2: MCP Tooling & Infrastructure Connections"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Git Provider:"), valueStyle.Render("["+m.cfg.Step2.GitProvider+"] (Required)")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("K8s Target:"), valueStyle.Render("["+m.cfg.Step2.K8sTarget+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("CI/CD Tools:"), valueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.CICDTools))))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Doc Sync:"), valueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.DocTools))))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Messenger Webhook:"), valueStyle.Render(fmt.Sprintf("%v", m.cfg.Step2.MessengerAlerts))))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Press Enter on K8s or Git to run interactive Health Check]"))

	case Step3:
		sb.WriteString(headerStyle.Render("🛠️ Step 3: Harness TDD, Commit & PR Conventions"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Commit Convention:"), valueStyle.Render("["+m.cfg.Step3.CommitConvention+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Issue Key Format:"), valueStyle.Render("["+m.cfg.Step3.IssueKeyFormat+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("PR Template Style:"), valueStyle.Render("["+m.cfg.Step3.PRTemplateStyle+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("TDD First Mode:"), valueStyle.Render(fmt.Sprintf("[%v]", m.cfg.Step3.TDDMode))))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Local Sandbox Test:"), valueStyle.Render(fmt.Sprintf("[%v]", m.cfg.Step3.LocalSandboxTest))))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Auto-generates AGENTS.md, CLAUDE.md & Commit Msg Git Hooks]"))

	case Step4:
		sb.WriteString(headerStyle.Render("🚀 Step 4: Release & Deployment Pipeline"))
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Versioning:"), valueStyle.Render("["+m.cfg.Step4.VersioningStrategy+"]")))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Auto Changelog:"), valueStyle.Render(fmt.Sprintf("[%v]", m.cfg.Step4.AutoChangelog))))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Release Notes Sync:"), valueStyle.Render(fmt.Sprintf("[%v]", m.cfg.Step4.ReleaseNotesSync))))
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", labelStyle.Render("Deploy Alerts:"), valueStyle.Render(fmt.Sprintf("[%v]", m.cfg.Step4.DeployAlert))))
		sb.WriteString("\n")
		sb.WriteString(hintStyle.Render("  [Triggers SemVer tagging and notifies Notion & Slack]"))
	}

	return sb.String()
}
