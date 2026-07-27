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

func TestRightSidebarRendersAllConfigs(t *testing.T) {
	cfg := config.NewDefaultConfig()
	m := NewModel(cfg)

	sidebarView := m.renderRightSidebarNav()

	expectedSectionsAndKeys := []string{
		"CONFIG STATUS NAV",
		"Session :",
		"App Name:",
		"Step 0: AI Licensing",
		"• Prov :",
		"• Model:",
		"• Auth :",
		"• Fallb:",
		"Step 1: Arch Spec",
		"• Style:",
		"• Repo :",
		"• Diag :",
		"Step 2: MCP Connections",
		"• Git  :",
		"• K8s  :",
		"• CI/CD:",
		"• Doc  :",
		"• Msg  :",
		"Step 3: Harness & TDD",
		"• Commit:",
		"• PRTpl :",
		"• TDD   :",
		"• Sbox  :",
		"Step 4: Release Pipeline",
		"• SemVer:",
		"• Change:",
		"• Sync  :",
		"• Alert :",
	}

	for _, key := range expectedSectionsAndKeys {
		if !strings.Contains(sidebarView, key) {
			t.Errorf("expected sidebar status nav to contain '%s', got:\n%s", key, sidebarView)
		}
	}
}
