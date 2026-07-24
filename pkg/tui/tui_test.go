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

func TestSlashCommandSelectAutocompletesWithoutInstantExecution(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Simulate typing '/' in promptInput
	m.promptInput.SetValue("/")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2 := updatedModel.(Model)

	if !m2.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be true after typing '/'")
	}

	// Press Enter on dropdown item to SELECT
	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m3 := updatedModel.(Model)

	// Dropdown must be closed, but mode MUST STILL BE ModePromptInput (Not executed or mode-switched!)
	if m3.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be false after selecting item")
	}

	if m3.mode != ModePromptInput {
		t.Errorf("expected mode to remain ModePromptInput for editing arguments, got %v", m3.mode)
	}

	// Prompt input should contain the autocompleted command with a trailing space
	if !strings.HasPrefix(m3.promptInput.Value(), "/set-name") && !strings.HasPrefix(m3.promptInput.Value(), "/git-init") {
		t.Errorf("expected promptInput value to be autocompleted, got '%s'", m3.promptInput.Value())
	}
}

func TestGenerationConfirmationFlow(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	m.promptInput.SetValue("/gen-all")
	m.slashDropdownOpen = false // Ensure dropdown is closed for direct command execution
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	if m2.mode != ModeConfirm {
		t.Errorf("expected /gen-all to enter ModeConfirm, got mode %v", m2.mode)
	}
}
