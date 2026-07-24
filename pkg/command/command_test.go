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

	foundSetName := false
	for _, c := range cmds {
		if c.Name == "/set-name" {
			foundSetName = true
		}
	}

	if !foundSetName {
		t.Error("expected /set-name in available slash commands")
	}
}

func TestParseAndExecuteSetName(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	res, err := cmdEngine.Execute("/set-name awesome-microservice")
	if err != nil {
		t.Fatalf("expected no error executing /set-name, got: %v", err)
	}

	if cfg.Step1.ProjectName != "awesome-microservice" {
		t.Errorf("expected project name to be 'awesome-microservice', got '%s'", cfg.Step1.ProjectName)
	}

	if res.Action != ActionSetName {
		t.Errorf("expected action ActionSetName, got %v", res.Action)
	}
}

func TestGitInitUsesSetNameProjectName(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	// First set project name via /set-name
	_, err := cmdEngine.Execute("/set-name test-repo-service")
	if err != nil {
		t.Fatalf("expected no error executing /set-name, got: %v", err)
	}

	// Then execute /git-init without positional arguments
	res, err := cmdEngine.Execute("/git-init")
	if err != nil {
		t.Fatalf("expected no error executing /git-init, got: %v", err)
	}

	if res.Action != ActionGitInit {
		t.Errorf("expected action ActionGitInit, got %v", res.Action)
	}
}
