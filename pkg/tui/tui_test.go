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

func TestRenderStepBodyAlignment(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	body := m.renderStepBody()
	lines := strings.Split(body, "\n")

	// Ensure no unexpected horizontal displacement (e.g. Licensing Mode should be aligned with Primary Model)
	foundLicensing := false
	foundPrimary := false

	for _, line := range lines {
		if strings.Contains(line, "Licensing Mode") {
			foundLicensing = true
			if !strings.HasPrefix(line, "  Licensing Mode:") {
				t.Errorf("expected 'Licensing Mode:' to start with '  Licensing Mode:', got '%s'", line)
			}
		}
		if strings.Contains(line, "Primary Model") {
			foundPrimary = true
			if !strings.HasPrefix(line, "  Primary Model:") {
				t.Errorf("expected 'Primary Model:' to start with '  Primary Model:', got '%s'", line)
			}
		}
	}

	if !foundLicensing || !foundPrimary {
		t.Errorf("expected step 0 fields in body output, licensing: %v, primary: %v", foundLicensing, foundPrimary)
	}
}
