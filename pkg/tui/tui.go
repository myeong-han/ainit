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
	Step2             // MCP Integrations
	Step3             // Git & Harness TDD & Conventions
	Step4             // Release Pipeline
)

func (s Step) String() string {
	switch s {
	case Step0:
		return "0. AI Licensing"
	case Step1:
		return "1. Architecture Spec"
	case Step2:
		return "2. MCP Tooling"
	case Step3:
		return "3. Harness TDD & Conventions"
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

// Styles using Lipgloss for rich UI rendering
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFD1")).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("#00FFD1")).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262")).
				Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFD1")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00")).
			Italic(true)
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
		statusMsg:   "Press 'Tab' or '←/→' to navigate steps. 'Ctrl+C' or 'q' to quit.",
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

	// Step Tabs
	var tabs []string
	steps := []Step{Step0, Step1, Step2, Step3, Step4}
	for _, s := range steps {
		if s == m.currentStep {
			tabs = append(tabs, activeTabStyle.Render(s.String()))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(s.String()))
		}
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	sb.WriteString("\n\n")

	// Step Body View
	body := m.renderStepBody()
	sb.WriteString(boxStyle.Render(body))
	sb.WriteString("\n\n")

	// Footer / Navigation Guide
	sb.WriteString(statusStyle.Render(m.statusMsg))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) renderStepBody() string {
	var sb strings.Builder

	switch m.currentStep {
	case Step0:
		sb.WriteString(focusedStyle.Render("🤖 Step 0: AI Licensing & Model Setup\n\n"))
		sb.WriteString(fmt.Sprintf("  Licensing Mode: [%s]\n", m.cfg.Step0.LicensingMode))
		sb.WriteString(fmt.Sprintf("  Primary Model:  [%s]\n", m.cfg.Step0.PrimaryModel))
		sb.WriteString(fmt.Sprintf("  Fallback Model: [%s]\n", m.cfg.Step0.FallbackModel))
		sb.WriteString("\n  [Use Up/Down to select, Enter to toggle]")

	case Step1:
		sb.WriteString(focusedStyle.Render("🏗️  Step 1: Architecture Spec & Mermaid Generation\n\n"))
		sb.WriteString(fmt.Sprintf("  Project Name:      %s\n", m.cfg.Step1.ProjectName))
		sb.WriteString(fmt.Sprintf("  Arch Style:        [%s]\n", m.cfg.Step1.ArchitectureStyle))
		sb.WriteString(fmt.Sprintf("  Repo Structure:    [%s]\n", m.cfg.Step1.RepoStructure))
		sb.WriteString(fmt.Sprintf("  Mermaid Diagrams:  Sequence(%v), Traffic(%v), GitOps(%v)\n",
			m.cfg.Step1.GenerateSequence, m.cfg.Step1.GenerateTraffic, m.cfg.Step1.GenerateGitOps))

	case Step2:
		sb.WriteString(focusedStyle.Render("🔌 Step 2: MCP Tooling & Infrastructure Connections\n\n"))
		sb.WriteString(fmt.Sprintf("  Git Provider:      [%s] (Required)\n", m.cfg.Step2.GitProvider))
		sb.WriteString(fmt.Sprintf("  K8s Target:        [%s]\n", m.cfg.Step2.K8sTarget))
		sb.WriteString(fmt.Sprintf("  CI/CD Tools:       %v\n", m.cfg.Step2.CICDTools))
		sb.WriteString(fmt.Sprintf("  Doc Sync Tools:    %v\n", m.cfg.Step2.DocTools))
		sb.WriteString(fmt.Sprintf("  Messenger Webhook: %v\n", m.cfg.Step2.MessengerAlerts))

	case Step3:
		sb.WriteString(focusedStyle.Render("🛠️  Step 3: Harness TDD, Commit/PR Conventions\n\n"))
		sb.WriteString(fmt.Sprintf("  Commit Convention: [%s]\n", m.cfg.Step3.CommitConvention))
		sb.WriteString(fmt.Sprintf("  Issue Format:      [%s]\n", m.cfg.Step3.IssueKeyFormat))
		sb.WriteString(fmt.Sprintf("  PR Template:       [%s]\n", m.cfg.Step3.PRTemplateStyle))
		sb.WriteString(fmt.Sprintf("  TDD First Mode:    [%v]\n", m.cfg.Step3.TDDMode))
		sb.WriteString(fmt.Sprintf("  Local Sandbox:     [%v]\n", m.cfg.Step3.LocalSandboxTest))

	case Step4:
		sb.WriteString(focusedStyle.Render("🚀 Step 4: Release & Deployment Pipeline\n\n"))
		sb.WriteString(fmt.Sprintf("  Versioning:        [%s]\n", m.cfg.Step4.VersioningStrategy))
		sb.WriteString(fmt.Sprintf("  Auto Changelog:    [%v]\n", m.cfg.Step4.AutoChangelog))
		sb.WriteString(fmt.Sprintf("  Release Notes Sync:[%v]\n", m.cfg.Step4.ReleaseNotesSync))
		sb.WriteString(fmt.Sprintf("  Deploy Alerts:     [%v]\n", m.cfg.Step4.DeployAlert))
	}

	return sb.String()
}
