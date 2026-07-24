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

func TestEnhancedRightSidebarVisibility(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	sidebar := m.renderRightSidebarNav()

	// Should contain clear section dividers and status icons
	if !strings.Contains(sidebar, "─────────────") {
		t.Errorf("expected sidebar to contain horizontal section dividers, got:\n%s", sidebar)
	}

	if !strings.Contains(sidebar, "READY") && !strings.Contains(sidebar, "ACTIVE") {
		t.Errorf("expected sidebar to display status badges (READY/ACTIVE), got:\n%s", sidebar)
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
