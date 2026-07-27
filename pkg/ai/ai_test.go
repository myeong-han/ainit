package ai

import (
	"context"
	"strings"
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestMockAIGeneratorBuildsArchitectureSpecAndMermaid(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cfg.Step1.ProjectName = "payment-gateway"
	cfg.Step1.ArchitectureStyle = "msa"

	engine := NewEngine(cfg)
	userPrompt := "Build a high-performance payment gateway service with Redis cache and Kafka event streaming."

	ctx := context.Background()
	doc, err := engine.GenerateArchitectureDoc(ctx, userPrompt)
	if err != nil {
		t.Fatalf("unexpected error generating architecture doc: %v", err)
	}

	if !strings.Contains(doc, "# Architecture Specification: payment-gateway") {
		t.Errorf("expected doc to contain project title, got:\n%s", doc)
	}

	if !strings.Contains(doc, "```mermaid") {
		t.Errorf("expected doc to contain Mermaid diagrams, got:\n%s", doc)
	}

	if !strings.Contains(doc, "**Architecture Style**: `msa`") {
		t.Errorf("expected doc to mention MSA architecture style, got:\n%s", doc)
	}
}
