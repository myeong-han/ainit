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

func TestFieldToggleAndSelection(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Step 0, Cursor 0 (LicensingMode) -> Press Enter/Space to cycle from "subscription" to "apikey"
	if m.cfg.Step0.LicensingMode != "subscription" {
		t.Fatalf("expected initial licensing mode to be 'subscription', got '%s'", m.cfg.Step0.LicensingMode)
	}

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m2 := updatedModel.(Model)

	if m2.cfg.Step0.LicensingMode != "apikey" {
		t.Errorf("expected licensing mode after space press to be 'apikey', got '%s'", m2.cfg.Step0.LicensingMode)
	}

	// Move cursor down to Field 1 (Primary Model) and press Enter to cycle to "gpt-4o"
	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyDown})
	m3 := updatedModel.(Model)
	if m3.cursor != 1 {
		t.Fatalf("expected cursor to be 1, got %d", m3.cursor)
	}

	updatedModel, _ = m3.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m4 := updatedModel.(Model)
	if m4.cfg.Step0.PrimaryModel != "gpt-4o" {
		t.Errorf("expected primary model after enter press to be 'gpt-4o', got '%s'", m4.cfg.Step0.PrimaryModel)
	}
}

func TestRenderStepBodyAlignment(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	body := m.renderStepBody()
	lines := strings.Split(body, "\n")

	foundLicensing := false
	foundPrimary := false

	for _, line := range lines {
		if strings.Contains(line, "Licensing Mode") {
			foundLicensing = true
		}
		if strings.Contains(line, "Primary Model") {
			foundPrimary = true
		}
	}

	if !foundLicensing || !foundPrimary {
		t.Errorf("expected step 0 fields in body output, licensing: %v, primary: %v", foundLicensing, foundPrimary)
	}
}
