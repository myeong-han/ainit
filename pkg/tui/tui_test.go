package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/myeong-han/ainit/pkg/config"
)

func TestNewModel(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	if m.currentStep != Step0 {
		t.Errorf("expected initial step to be Step0, got %v", m.currentStep)
	}

	if m.mode != ModeWizard {
		t.Errorf("expected initial mode to be ModeWizard, got %v", m.mode)
	}
}

func TestSlashCommandExecutionInTUI(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Transition to Prompt Input Mode
	m.mode = ModePromptInput
	m.promptInput.SetValue("/set-confs --provider openai --arch monolith")

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	if m2.cfg.Step0.ProviderID != "openai" {
		t.Errorf("expected slash command to update provider to 'openai', got '%s'", m2.cfg.Step0.ProviderID)
	}

	if m2.cfg.Step1.ArchitectureStyle != "monolith" {
		t.Errorf("expected slash command to update arch style to 'monolith', got '%s'", m2.cfg.Step1.ArchitectureStyle)
	}
}

func TestTwoColumnLayoutWithRightSidebar(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	view := m.View()
	if !strings.Contains(view, "CONFIG STATUS") && !strings.Contains(view, "STATUS NAV") {
		t.Errorf("expected view to render right sidebar status nav, got:\n%s", view)
	}
}
