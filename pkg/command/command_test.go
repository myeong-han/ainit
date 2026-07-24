package command

import (
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestParseAndExecuteSetConfs(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	// Test /set-confs command
	result, err := cmdEngine.Execute("/set-confs --provider openai --arch monolith --git bitbucket")
	if err != nil {
		t.Fatalf("expected no error executing /set-confs, got: %v", err)
	}

	if cfg.Step0.ProviderID != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", cfg.Step0.ProviderID)
	}

	if cfg.Step1.ArchitectureStyle != "monolith" {
		t.Errorf("expected arch 'monolith', got '%s'", cfg.Step1.ArchitectureStyle)
	}

	if cfg.Step2.GitProvider != "bitbucket" {
		t.Errorf("expected git 'bitbucket', got '%s'", cfg.Step2.GitProvider)
	}

	if result.Action != ActionSetConfs {
		t.Errorf("expected action ActionSetConfs, got %v", result.Action)
	}
}

func TestParseGenDocs(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	result, err := cmdEngine.Execute("/gen-docs")
	if err != nil {
		t.Fatalf("expected no error executing /gen-docs, got: %v", err)
	}

	if result.Action != ActionGenDocs {
		t.Errorf("expected action ActionGenDocs, got %v", result.Action)
	}
}

func TestParseGenCodes(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cmdEngine := NewCommandEngine(cfg)

	result, err := cmdEngine.Execute("/gen-codes")
	if err != nil {
		t.Fatalf("expected no error executing /gen-codes, got: %v", err)
	}

	if result.Action != ActionGenCodes {
		t.Errorf("expected action ActionGenCodes, got %v", result.Action)
	}
}
