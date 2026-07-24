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

	if m.cfg.Step1.ProjectName != "ainit-project" {
		t.Errorf("expected config project name 'ainit-project', got '%s'", m.cfg.Step1.ProjectName)
	}
}

func TestModelStepNavigation(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

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

func TestOpenCodeProviderAndModelCascadingSelection(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Step 0, Cursor 0 (Provider Selection): Initial is "anthropic"
	if m.cfg.Step0.ProviderID != "anthropic" {
		t.Fatalf("expected initial provider to be 'anthropic', got '%s'", m.cfg.Step0.ProviderID)
	}

	// Press Enter to cycle Provider from "anthropic" to "openai"
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	if m2.cfg.Step0.ProviderID != "openai" {
		t.Errorf("expected provider after enter to be 'openai', got '%s'", m2.cfg.Step0.ProviderID)
	}

	// Primary Model should automatically update to an OpenAI model (e.g., "gpt-4o")
	if m2.cfg.Step0.PrimaryModel != "gpt-4o" {
		t.Errorf("expected primary model to cascade update to 'gpt-4o', got '%s'", m2.cfg.Step0.PrimaryModel)
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
