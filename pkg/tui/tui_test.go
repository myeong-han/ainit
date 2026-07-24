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

	sidebar := m.renderRightSidebarNav()
	if !strings.Contains(sidebar, "App Name: unknown") {
		t.Errorf("expected right sidebar to display 'App Name: unknown', got:\n%s", sidebar)
	}
}

func TestSlashCommandSelectAutocompletesWithoutInstantExecution(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	m.promptInput.SetValue("/")
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2 := updatedModel.(Model)

	if !m2.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be true after typing '/'")
	}

	// Press Enter on dropdown item (/git-init)
	updatedModel, _ = m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m3 := updatedModel.(Model)

	if m3.slashDropdownOpen {
		t.Error("expected slashDropdownOpen to be false after selecting item")
	}

	if !strings.HasPrefix(m3.promptInput.Value(), "/git-init") {
		t.Errorf("expected promptInput value to be autocompleted with '/git-init', got '%s'", m3.promptInput.Value())
	}
}

func TestGenerationConfirmationFlow(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	m.promptInput.SetValue("/gen-all")
	m.slashDropdownOpen = false
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updatedModel.(Model)

	if m2.mode != ModeConfirm {
		t.Errorf("expected /gen-all to enter ModeConfirm, got mode %v", m2.mode)
	}
}
