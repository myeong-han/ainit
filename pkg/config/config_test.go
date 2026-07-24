package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg.Step0.LicensingMode != "subscription" {
		t.Errorf("expected default licensing mode 'subscription', got '%s'", cfg.Step0.LicensingMode)
	}

	if cfg.Step1.ArchitectureStyle != "msa" {
		t.Errorf("expected default architecture 'msa', got '%s'", cfg.Step1.ArchitectureStyle)
	}

	if cfg.Step2.GitProvider != "github" {
		t.Errorf("expected default git provider 'github', got '%s'", cfg.Step2.GitProvider)
	}

	if cfg.Step3.CommitConvention != "conventional" {
		t.Errorf("expected default commit convention 'conventional', got '%s'", cfg.Step3.CommitConvention)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.Step1.ProjectName = ""

	err := cfg.Validate()
	if err == nil {
		t.Error("expected error when ProjectName is empty, got nil")
	}

	cfg.Step1.ProjectName = "my-awesome-app"
	err = cfg.Validate()
	if err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}
