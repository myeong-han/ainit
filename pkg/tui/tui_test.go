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

	// Default mode must be ModePromptInput (Chatting view)
	if m.mode != ModePromptInput {
		t.Errorf("expected default mode to be ModePromptInput (Chat), got %v", m.mode)
	}
}

func TestDynamicTerminalResponsiveLayout(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Simulate Terminal Window Resize event to 120x30
	updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m2 := updatedModel.(Model)

	if m2.width != 120 || m2.height != 30 {
		t.Errorf("expected terminal dimensions 120x30, got %dx%d", m2.width, m2.height)
	}

	view := m2.View()
	if !strings.Contains(view, "CONFIG STATUS NAV") {
		t.Errorf("expected right sidebar status nav to be rendered at right edge, got:\n%s", view)
	}
}

func TestSlashCommandSetConfsNavigation(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Input /set-confs command in chat to switch to wizard form mode
	m.promptInput.SetValue("/set-confs")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	if m2.mode != ModeWizard {
		t.Errorf("expected /set-confs to navigate to ModeWizard, got %v", m2.mode)
	}
}
