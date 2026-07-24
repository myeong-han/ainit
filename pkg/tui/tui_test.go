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

func TestTransitionToArchitecturePromptInputMode(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)
	m.currentStep = Step4
	m.cursor = 3 // Submit / Next button on Step 4

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	if m2.mode != ModePromptInput {
		t.Errorf("expected transition to ModePromptInput after submit on Step4, got %v", m2.mode)
	}
}

func TestRenderStepBodyAlignment(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	body := m.renderStepBody()
	lines := strings.Split(body, "\n")

	foundProvider := false
	foundPrimary := false

	for _, line := range lines {
		if strings.Contains(line, "AI Provider") {
			foundProvider = true
		}
		if strings.Contains(line, "Primary Model") {
			foundPrimary = true
		}
	}

	if !foundProvider || !foundPrimary {
		t.Errorf("expected step 0 provider & primary model fields, provider: %v, primary: %v", foundProvider, foundPrimary)
	}
}
