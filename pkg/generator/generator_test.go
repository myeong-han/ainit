package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateAgentContextFilesContainsPlaywrightAndTestLoopConcepts(t *testing.T) {
	tmpDir := t.TempDir()

	err := GenerateAgentContextFiles(tmpDir)
	if err != nil {
		t.Fatalf("failed to generate agent context files: %v", err)
	}

	agentsFile := filepath.Join(tmpDir, "AGENTS.md")
	content, err := os.ReadFile(agentsFile)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}

	strContent := string(content)

	expectedKeywords := []string{
		"Headless Playwright",
		"Goal-Driven Iterative Loop",
		"Local Test-Centric Verification",
		"make test",
	}

	for _, kw := range expectedKeywords {
		if !strings.Contains(strContent, kw) {
			t.Errorf("expected AGENTS.md to contain '%s', got:\n%s", kw, strContent)
		}
	}
}
