package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg.Step1.ProjectName != "unknown" {
		t.Errorf("expected default project name to be 'unknown', got '%s'", cfg.Step1.ProjectName)
	}

	if cfg.Step2.CI != "jenkins" {
		t.Errorf("expected default CI to be 'jenkins', got '%s'", cfg.Step2.CI)
	}

	if cfg.Step2.CD != "argocd" {
		t.Errorf("expected default CD to be 'argocd', got '%s'", cfg.Step2.CD)
	}

	if cfg.Step2.Doc != "notion" {
		t.Errorf("expected default Doc to be 'notion', got '%s'", cfg.Step2.Doc)
	}

	if cfg.Step2.Messenger != "slack" {
		t.Errorf("expected default Messenger to be 'slack', got '%s'", cfg.Step2.Messenger)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := NewDefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error validating default config: %v", err)
	}
}
