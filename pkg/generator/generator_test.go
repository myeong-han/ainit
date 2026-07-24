package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateAgentContextFiles(t *testing.T) {
	tempDir := t.TempDir()

	err := GenerateAgentContextFiles(tempDir)
	if err != nil {
		t.Fatalf("expected no error generating agent context files, got: %v", err)
	}

	expectedFiles := []string{
		"AGENTS.md",
		"CLAUDE.md",
		".cursorrules",
		filepath.Join(".github", "copilot-instructions.md"),
		".windsurfrules",
		filepath.Join(".gemini", "rules"),
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(tempDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("expected file %s to be generated, but it does not exist", file)
		}
	}
}
