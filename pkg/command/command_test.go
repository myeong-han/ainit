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
		if c.Name == "/set-name" {
			t.Error("/set-name should be removed from available slash commands")
		}
	}

	if !foundGitInit {
		t.Error("expected /git-init in available slash commands")
	}
}

func TestGitInitWithNameUpdatesProjectName(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	// Execute /git-init my-service-app
	res, err := cmdEngine.Execute("/git-init my-service-app")
	if err != nil {
		t.Fatalf("expected no error executing /git-init my-service-app, got: %v", err)
	}

	if cfg.Step1.ProjectName != "my-service-app" {
		t.Errorf("expected project name to be updated to 'my-service-app', got '%s'", cfg.Step1.ProjectName)
	}

	if res.Action != ActionGitInit {
		t.Errorf("expected action ActionGitInit, got %v", res.Action)
	}
}
