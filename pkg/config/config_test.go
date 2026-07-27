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

	if !cfg.Step3.PlaywrightCI {
		t.Error("expected default PlaywrightCI to be true")
	}

	if !cfg.Step3.LocalTestLoop {
		t.Error("expected default LocalTestLoop to be true")
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := NewDefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error validating default config: %v", err)
	}
}
