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

func TestSlashCommandDropdownActivationAndSelection(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Simulate typing '/' in promptInput
	m.promptInput.SetValue("/")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2 := updatedModel.(Model)

	if !m2.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be true after typing '/'")
	}

	// Press Tab to autocomplete current selected dropdown item (/set-confs)
	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyTab})
	m3 := updatedModel.(Model)

	if !strings.HasPrefix(m3.promptInput.Value(), "/set-confs") {
		t.Errorf("expected promptInput value to be autocompleted with '/set-confs', got '%s'", m3.promptInput.Value())
	}
}

func TestDynamicTerminalResponsiveLayout(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m2 := updatedModel.(Model)

	if m2.width != 120 || m2.height != 30 {
		t.Errorf("expected terminal dimensions 120x30, got %dx%d", m2.width, m2.height)
	}
}
