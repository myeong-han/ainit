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

func TestRightSidebarLinesDoNotOverflowOrWrap(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	sidebarRaw := m.renderRightSidebarNav()
	lines := strings.Split(sidebarRaw, "\n")

	for i, line := range lines {
		cleanLine := stripAnsi(line)
		if len([]rune(cleanLine)) > 26 {
			t.Errorf("line %d in sidebar exceeds max width (26): length %d, content: '%s'", i, len([]rune(cleanLine)), cleanLine)
		}
	}
}

func stripAnsi(str string) string {
	var sb strings.Builder
	inAnsi := false
	for _, r := range str {
		if r == '\x1b' {
			inAnsi = true
			continue
		}
		if inAnsi {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inAnsi = false
			}
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}
