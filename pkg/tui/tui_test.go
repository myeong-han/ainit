package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

func TestDropdownDismissPersistence(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Simulate typing '/'
	m.promptInput.SetValue("/")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2 := updatedModel.(Model)

	if !m2.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be true after typing '/'")
	}

	// Press Esc to dismiss dropdown
	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m3 := updatedModel.(Model)

	if m3.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be false after pressing Esc")
	}

	// Continue typing 'g' when dismissed - dropdown MUST REMAIN CLOSED
	m3.promptInput.SetValue("/g")
	updatedModel, _ = m3.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m4 := updatedModel.(Model)

	if m4.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to REMAIN FALSE after typing additional characters when dismissed")
	}
}

func TestSlashCommandSelectAutocompletesWithoutInstantExecution(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	m.promptInput.SetValue("/")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2 := updatedModel.(Model)

	// Select dropdown item
	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m3 := updatedModel.(Model)

	if m3.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be false after selecting item")
	}

	if !strings.HasPrefix(m3.promptInput.Value(), "/git-init") {
		t.Errorf("expected promptInput value to be autocompleted with '/git-init', got '%s'", m3.promptInput.Value())
	}
}
