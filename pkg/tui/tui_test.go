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

func TestTwoColumnLayoutWithRightSidebar(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	view := m.View()
	// Should render the right sidebar containing step status overview
	if !strings.Contains(view, "CONFIG STATUS") && !strings.Contains(view, "STATUS NAV") {
		t.Errorf("expected view to render right sidebar status nav, got:\n%s", view)
	}

	if !strings.Contains(view, "Step 0:") || !strings.Contains(view, "Step 1:") {
		t.Errorf("expected sidebar to display step 0 and step 1 configurations, got:\n%s", view)
	}
}

func TestTransitionToArchitecturePromptInputMode(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)
	m.currentStep = Step4
	m.cursor = 3 // Submit button on Step 4

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	if m2.mode != ModePromptInput {
		t.Errorf("expected transition to ModePromptInput after submit on Step4, got %v", m2.mode)
	}
}
