package command

import (
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestGetAvailableSlashCommands(t *testing.T) {
	cmds := GetAvailableSlashCommands()
	if len(cmds) == 0 {
		t.Fatal("expected non-empty list of slash commands")
	}

	foundGitInit := false
	for _, c := range cmds {
		if c.Name == "/git-init" {
			foundGitInit = true
		}
	}

	if !foundGitInit {
		t.Error("expected /git-init in available slash commands")
	}
}

func TestParseAndExecuteGitInit(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	result, err := cmdEngine.Execute("/git-init myeong-han/ainit")
	if err != nil {
		t.Fatalf("expected no error executing /git-init, got: %v", err)
	}

	if result.Action != ActionGitInit {
		t.Errorf("expected action ActionGitInit, got %v", result.Action)
	}
}

func TestParseAndExecuteSetConfs(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	_, err := cmdEngine.Execute("/set-confs --provider openai --arch monolith --git bitbucket")
	if err != nil {
		t.Fatalf("expected no error executing /set-confs, got: %v", err)
	}

	if cfg.Step0.ProviderID != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", cfg.Step0.ProviderID)
	}
}
