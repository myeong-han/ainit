package provider

import (
	"context"
	"testing"
)

func TestGetDefaultProvidersContainsLatestModels(t *testing.T) {
	providers := GetAvailableProviders()
	if len(providers) == 0 {
		t.Fatal("expected available providers, got none")
	}

	foundAnthropic := false
	for _, p := range providers {
		if p.ID == "anthropic" {
			foundAnthropic = true
			foundSonnet37 := false
			for _, m := range p.Models {
				if m.ID == "claude-3-7-sonnet" || m.ID == "claude-3-5-sonnet-20241022" {
					foundSonnet37 = true
					break
				}
			}
			if !foundSonnet37 {
				t.Errorf("expected anthropic models to include latest sonnet, got %v", p.Models)
			}
		}
	}

	if !foundAnthropic {
		t.Error("expected anthropic provider in list")
	}
}

func TestFetchLiveModelsFallback(t *testing.T) {
	ctx := context.Background()
	// Fetching with empty key should return updated fallback model catalog gracefully
	models, err := FetchLiveModelsForProvider(ctx, "openai", "")
	if err != nil {
		t.Fatalf("unexpected error fetching models: %v", err)
	}

	if len(models) == 0 {
		t.Error("expected models list, got empty")
	}
}
