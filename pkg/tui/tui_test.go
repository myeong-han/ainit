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

func TestStep2CICDAndDocSyncSelectionToggle(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)
	m.currentStep = Step2

	if m.getMaxCursorForStep() != 3 {
		t.Errorf("expected Step2 maxCursor to be 3, got %d", m.getMaxCursorForStep())
	}

	// Move cursor to index 2 (CI/CD Tools) and toggle
	m.cursor = 2
	m.toggleCurrentField()
	if len(m.cfg.Step2.CICDTools) == 0 {
		t.Error("expected CICDTools to be updated after toggle")
	}

	// Move cursor to index 3 (Doc Sync) and toggle
	m.cursor = 3
	m.toggleCurrentField()
	if len(m.cfg.Step2.DocTools) == 0 {
		t.Error("expected DocTools to be updated after toggle")
	}

	body := m.renderStepBody()
	if !strings.Contains(body, "CI/CD Tools:") || !strings.Contains(body, "Doc Sync:") {
		t.Errorf("expected step body to contain interactive rows for CI/CD and Doc Sync, got:\n%s", body)
	}
}
