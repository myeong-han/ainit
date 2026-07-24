package provider

import (
	"testing"
)

func TestGetProviders(t *testing.T) {
	providers := GetAvailableProviders()
	if len(providers) == 0 {
		t.Fatal("expected non-empty list of available providers")
	}

	foundAnthropic := false
	for _, p := range providers {
		if p.ID == "anthropic" {
			foundAnthropic = true
			if len(p.Models) == 0 {
				t.Errorf("expected anthropic models to be non-empty")
			}
		}
	}

	if !foundAnthropic {
		t.Errorf("expected anthropic provider to exist in catalog")
	}
}

func TestGetModelsByProvider(t *testing.T) {
	models := GetModelsForProvider("openai")
	if len(models) == 0 {
		t.Fatal("expected openai models, got empty list")
	}

	foundGPT4o := false
	for _, m := range models {
		if m.ID == "gpt-4o" {
			foundGPT4o = true
			if m.ContextWindow == 0 {
				t.Errorf("expected context window > 0 for gpt-4o")
			}
		}
	}

	if !foundGPT4o {
		t.Errorf("expected gpt-4o model under openai provider")
	}
}
