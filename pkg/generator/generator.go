package generator

import (
	"os"
	"path/filepath"
)

// GenerateAgentContextFiles generates agent context rules for all major AI agent platforms,
// establishing AGENTS.md as the Single Source of Truth.
func GenerateAgentContextFiles(targetDir string) error {
	files := map[string]string{
		"AGENTS.md": `# Agentic-Init (ainit) Execution Rules

## General Principles
1. **Atomic Micro-Commits**: All generated code must be committed in small, single-purpose commits.
2. **Commit Convention Compliance**: Every commit message MUST strictly adhere to the configured Commit Convention.
3. **PR Convention Compliance**: Every Pull Request MUST follow the generated PR template and satisfy all checklist criteria.
4. **TDD First**: Write failing test specifications prior to writing feature code.
5. **Local Sandbox Verification**: Execute build & test suites locally before executing git commit.
6. **Documentation Sync**: Keep Notion and Confluence documentation up-to-date via MCP integrations.
`,
		"CLAUDE.md": `# CLAUDE.md - Agent Context

> ⚠️ **CRITICAL AGENT INSTRUCTION**: This repository uses [AGENTS.md](./AGENTS.md) as the Single Source of Truth.
> Claude Code MUST strictly read, respect, and obey all instructions in AGENTS.md.
`,
		".cursorrules": `# Cursor Rules Configuration

# CRITICAL: Always refer to AGENTS.md for all project execution rules, TDD loops, and commit conventions.
# AGENTS.md is the Single Source of Truth for this codebase.

- Always read AGENTS.md before making any code modifications.
- Execute unit tests (make test) before finalizing any changes.
`,
		filepath.Join(".github", "copilot-instructions.md"): `# GitHub Copilot Custom Instructions

All AI Copilot and Agent behaviors MUST be strictly governed by [AGENTS.md](../AGENTS.md).
Refer to AGENTS.md as the Single Source of Truth.
`,
		".windsurfrules": `# Windsurf Rules Configuration

# Refer to AGENTS.md for all project execution rules, TDD loops, and commit conventions.
# AGENTS.md is the Single Source of Truth for this codebase.
`,
		filepath.Join(".gemini", "rules"): `# Gemini & Antigravity Agent Rules

# Refer to AGENTS.md for all project execution rules, TDD loops, and commit conventions.
# AGENTS.md is the Single Source of Truth for this codebase.
`,
	}

	for relPath, content := range files {
		fullPath := filepath.Join(targetDir, relPath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}
