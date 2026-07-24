package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg.Step0.ProviderID != "anthropic" {
		t.Errorf("expected default provider anthropic, got %s", cfg.Step0.ProviderID)
	}

	// Default project name must be 'unknown'
	if cfg.Step1.ProjectName != "unknown" {
		t.Errorf("expected default project name 'unknown', got '%s'", cfg.Step1.ProjectName)
	}

	if cfg.Step3.CommitConvention != "conventional" {
		t.Errorf("expected default commit convention conventional, got %s", cfg.Step3.CommitConvention)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := NewDefaultConfig()
	err := cfg.Validate()
	if err != nil {
		t.Errorf("expected config to be valid, got %v", err)
	}
}
