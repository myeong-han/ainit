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

func TestNoVariableWidthEmojisInSidebar(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	sidebarRaw := m.renderRightSidebarNav()

	// Variable-width emojis cause terminal font engine rendering misalignment
	unstableEmojis := []string{"🤖", "👤", "💬", "📊", "🏗️", "🔌", "🛠️", "🚀", "🟢", "⚡", "❌", "⚠️"}
	for _, emoji := range unstableEmojis {
		if strings.Contains(sidebarRaw, emoji) {
			t.Errorf("sidebar rendering contains variable-width emoji '%s' which causes terminal font alignment bugs", emoji)
		}
	}
}
