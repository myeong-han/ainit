package tui

import (
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

func TestKeyInputModeAndRealConnectionTesting(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Switch to wizard mode Step 0 (AI Provider)
	m.mode = ModeWizard
	m.currentStep = Step0
	m.cursor = 2 // Auth Method / API Key input

	m.toggleCurrentField()

	if m.mode != ModeKeyInput {
		t.Errorf("expected ModeKeyInput after selecting API Key field, got %v", m.mode)
	}
}
