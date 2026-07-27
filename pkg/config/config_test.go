package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg.Step1.ProjectName != "unknown" {
		t.Errorf("expected default project name to be 'unknown', got '%s'", cfg.Step1.ProjectName)
	}

	expectedKubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if cfg.Step2.KubeconfigPath != expectedKubeconfig {
		t.Errorf("expected default kubeconfig path to be '%s', got '%s'", expectedKubeconfig, cfg.Step2.KubeconfigPath)
	}

	if cfg.Step2.CI != "jenkins" {
		t.Errorf("expected default CI to be 'jenkins', got '%s'", cfg.Step2.CI)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := NewDefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error validating default config: %v", err)
	}
}
