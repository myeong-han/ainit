package command

import (
	"path/filepath"
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestGetAvailableSlashCommands(t *testing.T) {
	cmds := GetAvailableSlashCommands()
	foundGitInit := false
	foundSettings := false
	foundResume := false

	for _, cmd := range cmds {
		if cmd.Name == "/git-init" {
			foundGitInit = true
		}
		if cmd.Name == "/settings" {
			foundSettings = true
		}
		if cmd.Name == "/resume" {
			foundResume = true
		}
	}

	if !foundGitInit || !foundSettings || !foundResume {
		t.Errorf("expected /git-init, /settings and /resume in slash commands list, got %v", cmds)
	}
}

func TestGitInitWithNameUpdatesProjectName(t *testing.T) {
	cfg := config.NewDefaultConfig()
	engine := NewCommandEngine(cfg)

	targetDir := filepath.Join(".", "my-service-app")

	res, err := engine.Execute("/git-init my-service-app")
	if err != nil {
		t.Fatalf("unexpected error executing /git-init: %v", err)
	}

	if cfg.Step1.ProjectName != "my-service-app" {
		t.Errorf("expected ProjectName to be updated to 'my-service-app', got '%s'", cfg.Step1.ProjectName)
	}

	if res.Action != ActionGitInit {
		t.Errorf("expected ActionGitInit, got %v", res.Action)
	}

	_ = targetDir
}
