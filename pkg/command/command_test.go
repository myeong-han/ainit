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

	foundGenGitOps := false
	foundGenAll := false
	for _, c := range cmds {
		if c.Name == "/gen-gitops" {
			foundGenGitOps = true
		}
		if c.Name == "/gen-all" {
			foundGenAll = true
		}
	}

	if !foundGenGitOps || !foundGenAll {
		t.Errorf("expected /gen-gitops and /gen-all in available slash commands, got gitops:%v all:%v", foundGenGitOps, foundGenAll)
	}
}

func TestSetConfsProjectNameOption(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	_, err := cmdEngine.Execute("/set-confs --name my-custom-service --provider openai")
	if err != nil {
		t.Fatalf("expected no error executing /set-confs with --name, got: %v", err)
	}

	if cfg.Step1.ProjectName != "my-custom-service" {
		t.Errorf("expected project name 'my-custom-service', got '%s'", cfg.Step1.ProjectName)
	}
}

func TestGenGitOpsAndGenAll(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	resGitOps, err := cmdEngine.Execute("/gen-gitops")
	if err != nil {
		t.Fatalf("expected no error executing /gen-gitops, got: %v", err)
	}
	if resGitOps.Action != ActionGenGitOps {
		t.Errorf("expected action ActionGenGitOps, got %v", resGitOps.Action)
	}

	resAll, err := cmdEngine.Execute("/gen-all")
	if err != nil {
		t.Fatalf("expected no error executing /gen-all, got: %v", err)
	}
	if resAll.Action != ActionGenAll {
		t.Errorf("expected action ActionGenAll, got %v", resAll.Action)
	}
}
