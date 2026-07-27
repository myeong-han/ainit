package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/myeong-han/ainit/pkg/ai"
	"github.com/myeong-han/ainit/pkg/config"
)

// GenerateHarnessProject generates docs/ARCHITECTURE_SPEC.md via pkg/ai Engine
func GenerateHarnessProject(targetDir string, cfg *config.Config, userPrompt string) error {
	docsDir := filepath.Join(targetDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return fmt.Errorf("failed to create docs dir: %w", err)
	}

	aiEngine := ai.NewEngine(cfg)
	specContent, err := aiEngine.GenerateArchitectureDoc(context.Background(), userPrompt)
	if err != nil {
		return fmt.Errorf("failed to synthesize architecture spec with AI engine: %w", err)
	}

	specPath := filepath.Join(docsDir, "ARCHITECTURE_SPEC.md")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		return fmt.Errorf("failed to write ARCHITECTURE_SPEC.md: %w", err)
	}

	return nil
}

// GenerateAgentContextFiles generates AGENTS.md, CLAUDE.md, .cursorrules & cross-agent rules & code scaffolding
func GenerateAgentContextFiles(targetDir string) error {
	githubDir := filepath.Join(targetDir, ".github")
	geminiDir := filepath.Join(targetDir, ".gemini")

	_ = os.MkdirAll(githubDir, 0755)
	_ = os.MkdirAll(geminiDir, 0755)

	agentsMd := "# AGENTS.md - Agentic Harness Execution Rules\n\n" +
		"## 1. Core Principles\n" +
		"- **Atomic Micro-Commits**: All code modifications MUST be committed in single-purpose, small micro-commits.\n" +
		"- **TDD First**: Write failing test specifications prior to writing feature implementation code.\n" +
		"- **Convention Compliance**: Adhere strictly to Conventional Commits format (`feat:`, `fix:`, `refactor:`, `test:`).\n" +
		"- **Local Sandbox Verification**: Execute build & test suites locally before executing `git commit`.\n\n" +
		"## 2. Commit Format\n" +
		"```\n" +
		"<type>(<scope>): <short summary>\n" +
		"```\n"

	claudeMd := "# CLAUDE.md - Anthropic Agentic Guidelines\n\n" +
		"- **Primary Goal**: Maintain high-precision code quality and zero symptom-patching.\n" +
		"- **Workflow**: Plan -> Failing Test -> Clean Implementation -> Local Verification -> Commit.\n"

	cursorRules := "# .cursorrules - AI Coding Assistant System Prompt\n\n" +
		"- Always maintain atomic micro-commits.\n" +
		"- Run test suite prior to declaring task resolution.\n"

	copilotInst := "# GitHub Copilot Instructions\n\n- Follow Conventional Commits and TDD Workflow.\n"
	windsurfRules := "# Windsurf Cascade Rules\n\n- Execute local build and test checks before committing code.\n"
	geminiRules := "# Antigravity / Gemini Rules\n\n- Adhere to AGENTS.md execution rules.\n"

	if err := os.WriteFile(filepath.Join(targetDir, "AGENTS.md"), []byte(agentsMd), 0644); err != nil {
		return fmt.Errorf("failed to write AGENTS.md: %w", err)
	}

	if err := os.WriteFile(filepath.Join(targetDir, "CLAUDE.md"), []byte(claudeMd), 0644); err != nil {
		return fmt.Errorf("failed to write CLAUDE.md: %w", err)
	}

	if err := os.WriteFile(filepath.Join(targetDir, ".cursorrules"), []byte(cursorRules), 0644); err != nil {
		return fmt.Errorf("failed to write .cursorrules: %w", err)
	}

	_ = os.WriteFile(filepath.Join(githubDir, "copilot-instructions.md"), []byte(copilotInst), 0644)
	_ = os.WriteFile(filepath.Join(targetDir, ".windsurfrules"), []byte(windsurfRules), 0644)
	_ = os.WriteFile(filepath.Join(geminiDir, "rules"), []byte(geminiRules), 0644)

	// Generate Code Scaffolding
	projName := filepath.Base(targetDir)
	if projName == "." || projName == "" {
		projName = "app"
	}

	goMod := fmt.Sprintf("module %s\n\ngo 1.21\n", projName)
	mainGo := fmt.Sprintf("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"🚀 %s service initialized by Agentic-Init (ainit)\")\n}\n", projName)

	dockerfile := fmt.Sprintf("FROM golang:1.21-alpine AS builder\nWORKDIR /app\nCOPY . .\nRUN go build -o bin/server ./...\n\nFROM alpine:latest\nWORKDIR /app\nCOPY --from=builder /app/bin/server /app/server\nEXPOSE 8080\nCMD [\"/app/server\"]\n")

	makefile := fmt.Sprintf(".PHONY: build test run\n\nbuild:\n\t@echo \"🔨 Building %s...\"\n\t@go build -o bin/%s main.go\n\ntest:\n\t@echo \"🧪 Running unit tests...\"\n\t@go test -v ./...\n\nrun: build\n\t@./bin/%s\n", projName, projName, projName)

	_ = os.WriteFile(filepath.Join(targetDir, "go.mod"), []byte(goMod), 0644)
	_ = os.WriteFile(filepath.Join(targetDir, "main.go"), []byte(mainGo), 0644)
	_ = os.WriteFile(filepath.Join(targetDir, "Dockerfile"), []byte(dockerfile), 0644)
	_ = os.WriteFile(filepath.Join(targetDir, "Makefile"), []byte(makefile), 0644)

	return nil
}

// GenerateGitOpsManifests generates Helm chart & ArgoCD application YAML
func GenerateGitOpsManifests(targetDir string, cfg *config.Config) error {
	gitopsDir := filepath.Join(targetDir, "gitops")
	helmDir := filepath.Join(gitopsDir, "helm")
	argocdDir := filepath.Join(gitopsDir, "argocd")

	if err := os.MkdirAll(helmDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(argocdDir, 0755); err != nil {
		return err
	}

	projName := cfg.Step1.ProjectName
	if projName == "" || projName == "unknown" {
		projName = "harness-app"
	}

	chartYaml := fmt.Sprintf("apiVersion: v2\nname: %s\ndescription: Helm chart generated by Agentic-Init (ainit)\ntype: application\nversion: 0.1.0\nappVersion: \"1.0.0\"\n", projName)

	valuesYaml := fmt.Sprintf("replicaCount: 2\n\nimage:\n  repository: %s/%s\n  pullPolicy: IfNotPresent\n  tag: \"latest\"\n\nservice:\n  type: ClusterIP\n  port: 8080\n\ningress:\n  enabled: true\n  className: \"nginx\"\n  hosts:\n    - host: %s.local\n      paths:\n        - path: /\n          pathType: ImplementationSpecific\n\nresources:\n  limits:\n    cpu: 500m\n    memory: 512Mi\n  requests:\n    cpu: 100m\n    memory: 128Mi\n", cfg.Step2.GitProvider, projName, projName)

	argoApplication := fmt.Sprintf("apiVersion: argoproj.io/v1alpha1\nkind: Application\nmetadata:\n  name: %s-gitops\n  namespace: argocd\nspec:\n  project: default\n  source:\n    repoURL: 'https://github.com/myeong-han/%s.git'\n    targetRevision: HEAD\n    path: gitops/helm\n  destination:\n    server: 'https://kubernetes.default.svc'\n    namespace: default\n  syncPolicy:\n    automated:\n      prune: true\n      selfHeal: true\n", projName, projName)

	if err := os.WriteFile(filepath.Join(helmDir, "Chart.yaml"), []byte(chartYaml), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(helmDir, "values.yaml"), []byte(valuesYaml), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(argocdDir, "application.yaml"), []byte(argoApplication), 0644); err != nil {
		return err
	}

	return nil
}

// GenerateAll executes GenerateHarnessProject, GenerateAgentContextFiles & GenerateGitOpsManifests sequentially
func GenerateAll(targetDir string, cfg *config.Config, userPrompt string) error {
	if err := GenerateHarnessProject(targetDir, cfg, userPrompt); err != nil {
		return fmt.Errorf("step 1 /gen-docs failed: %w", err)
	}
	if err := GenerateAgentContextFiles(targetDir); err != nil {
		return fmt.Errorf("step 2 /gen-codes failed: %w", err)
	}
	if err := GenerateGitOpsManifests(targetDir, cfg); err != nil {
		return fmt.Errorf("step 3 /gen-gitops failed: %w", err)
	}
	return nil
}

func stringsTitle(s string) string {
	if len(s) == 0 {
		return s
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(s[:1]), s[1:])
}
