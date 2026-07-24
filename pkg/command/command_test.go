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

	foundSetConfs := false
	for _, c := range cmds {
		if c.Name == "/set-confs" {
			foundSetConfs = true
		}
	}

	if !foundSetConfs {
		t.Error("expected /set-confs in available slash commands")
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
