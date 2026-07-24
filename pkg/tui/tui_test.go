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

func TestDynamicFilteringAndAutoFocusDropdown(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	// Simulate typing '/gen-g' by setting promptInput to '/gen-' and sending 'g' key
	m.promptInput.SetValue("/gen-")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m2 := updatedModel.(Model)

	if !m2.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be true for /gen-g prefix match")
	}

	if len(m2.slashOptions) != 1 || m2.slashOptions[0].Name != "/gen-gitops" {
		t.Errorf("expected filtered options to contain only '/gen-gitops', got %v", m2.slashOptions)
	}

	// Press Tab to autocomplete '/gen-gitops'
	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyTab})
	m3 := updatedModel.(Model)

	if !strings.HasPrefix(m3.promptInput.Value(), "/gen-gitops") {
		t.Errorf("expected promptInput value to be autocompleted with '/gen-gitops', got '%s'", m3.promptInput.Value())
	}
}
