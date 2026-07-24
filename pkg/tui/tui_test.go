package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/myeong-han/ainit/pkg/config"
)

func TestNewModelDefaultIsChatMode(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	if m.mode != ModePromptInput {
		t.Errorf("expected default mode to be ModePromptInput (Chat), got %v", m.mode)
	}
}

func TestGenerationConfirmationFlow(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Simulate executing /gen-all
	m.promptInput.SetValue("/gen-all")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	// Must enter ModeConfirm first to prevent unintended generation
	if m2.mode != ModeConfirm {
		t.Errorf("expected /gen-all to enter ModeConfirm, got mode %v", m2.mode)
	}

	// Press 'n' to cancel generation
	cancelledModel, _ := m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m3 := cancelledModel.(Model)

	if m3.mode != ModePromptInput {
		t.Errorf("expected pressing 'n' to return to ModePromptInput, got %v", m3.mode)
	}
	if !strings.Contains(m3.statusMsg, "cancelled") {
		t.Errorf("expected statusMsg to indicate cancellation, got '%s'", m3.statusMsg)
	}
}

func TestEnhancedRightSidebarVisibility(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	sidebar := m.renderRightSidebarNav()

	if !strings.Contains(sidebar, "─────────────") {
		t.Errorf("expected sidebar to contain horizontal section dividers, got:\n%s", sidebar)
	}
}
