package ai

import (
	"context"
	"strings"
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestGenerateChatResponseWithMockOrRealAPI(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cfg.Step0.ProviderID = "gemini"
	cfg.Step0.PrimaryModel = "gemini-2.0-flash"

	engine := NewEngine(cfg)

	resp, err := engine.GenerateChatResponse(context.Background(), "너는 동작하니?", "")
	if err != nil {
		t.Fatalf("unexpected error generating chat response: %v", err)
	}

	if resp == "" {
		t.Error("expected non-empty AI chat response")
	}

	if !strings.Contains(resp, "gemini") && !strings.Contains(resp, "Agentic-Init") {
		t.Errorf("expected response to mention Provider or Agentic-Init context, got: %s", resp)
	}
}
