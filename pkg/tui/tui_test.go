package tui

import (
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

	if m.cfg.Step1.ProjectName != "ainit-project" {
		t.Errorf("expected config project name 'ainit-project', got '%s'", m.cfg.Step1.ProjectName)
	}
}

func TestModelStepNavigation(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Simulate pressing 'tab' or 'down' / 'enter' to advance step
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m2 := updatedModel.(Model)

	if m2.currentStep != Step1 {
		t.Errorf("expected step after right arrow to be Step1, got %v", m2.currentStep)
	}

	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m3 := updatedModel.(Model)

	if m3.currentStep != Step0 {
		t.Errorf("expected step after left arrow to return to Step0, got %v", m3.currentStep)
	}
}
