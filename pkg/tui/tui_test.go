package tui

import (
	"strings"
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestNewModelDefaultIsChatModeAndUnknownAppName(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	if m.mode != ModePromptInput {
		t.Errorf("expected default mode to be ModePromptInput (Chat), got %v", m.mode)
	}

	if m.cfg.Step1.ProjectName != "unknown" {
		t.Errorf("expected initial project name to be 'unknown', got '%s'", m.cfg.Step1.ProjectName)
	}
}

func TestStep2SeparatedFixedTools(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)
	m.currentStep = Step2

	body := m.renderStepBody()

	expectedItems := []string{
		"CI Tool:",
		"CD Tool:",
		"Doc Tool:",
		"Messenger:",
		"jenkins",
		"argocd",
		"notion",
		"slack",
	}

	for _, item := range expectedItems {
		if !strings.Contains(body, item) {
			t.Errorf("expected Step 2 body to contain '%s', got:\n%s", item, body)
		}
	}
}
