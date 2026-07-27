package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/myeong-han/ainit/pkg/ai"
	"github.com/myeong-han/ainit/pkg/config"
)

// GenerateHarnessProject produces architecture docs, agent rules, and code scaffolding
func GenerateHarnessProject(targetDir string, cfg *config.Config, prompt string) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	engine := ai.NewEngine(cfg)
	markdownDoc, err := engine.GenerateArchitectureDoc(context.Background(), prompt)
	if err != nil {
		return fmt.Errorf("failed to generate architecture spec: %w", err)
	}

	docsDir := filepath.Join(targetDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return fmt.Errorf("failed to create docs dir: %w", err)
	}

	docFile := filepath.Join(docsDir, "ARCHITECTURE_SPEC.md")
	if err := os.WriteFile(docFile, []byte(markdownDoc), 0644); err != nil {
		return fmt.Errorf("failed to write ARCHITECTURE_SPEC.md: %w", err)
	}

	if err := GenerateAgentContextFiles(targetDir); err != nil {
		return fmt.Errorf("failed to generate agent context files: %w", err)
	}

	if err := GenerateGitOpsManifests(targetDir, cfg); err != nil {
		return fmt.Errorf("failed to generate GitOps manifests: %w", err)
	}

	return nil
}

// GenerateAgentContextFiles builds AGENTS.md, CLAUDE.md, .cursorrules, and code scaffolding
func GenerateAgentContextFiles(targetDir string) error {
	agentsContent := `# AGENTS.md - Agentic Harness Execution Rules

## 1. Core Principles & CI/Dev Execution Rules
- **Atomic Micro-Commits**: All code modifications MUST be committed in single-purpose, small micro-commits.
- **TDD First**: Write failing test specifications prior to writing feature implementation code.
- **CI Pipeline Rule (Headless Playwright)**: In CI scripts and automated test pipelines, ALWAYS execute E2E/UI verification via Headless Playwright ('npx playwright test --headed=false').
- **Development Goal Loop**: Development phase MUST continuously run a Goal-Driven Iterative Loop until the explicit objective is satisfied.
- **Local Test-Centric Verification**: Loop completion is strictly judged by Local Test-Centric Verification ('make test' & unit/integration test suite PASS). Never declare done without green local test execution.

## 2. Commit Format
` + "```" + `
<type>(<scope>): <short summary>
` + "```" + `
- Types: feat, fix, refactor, test, docs, chore
`

	claudeContent := `# CLAUDE.md - Context & Coding Guidelines

## Build & Test Commands
- Local Test Suite Execution (Local Test-Centric Verification): make test
- Build Binary: make build
- CI E2E Automation: npx playwright test --headed=false

## Architecture & Development Guidelines
- Always execute Goal-Driven Iterative Loop in development until all local test suites pass cleanly.
- Enforce Headless Playwright for E2E validation in CI pipelines.
`

	cursorRulesContent := `# .cursorrules - Agent Execution Context

rule "Headless Playwright in CI":
  description: "Always run Playwright in headless mode within CI/CD pipelines"
  command: "npx playwright test --headed=false"

rule "Goal-Driven Iterative Loop":
  description: "Continuously iterate in development phase until objective is achieved"

rule "Local Test-Centric Verification":
  description: "Judge loop completion strictly based on local test suite results ('make test')"
`

	files := map[string]string{
		"AGENTS.md":                            agentsContent,
		"CLAUDE.md":                            claudeContent,
		".cursorrules":                         cursorRulesContent,
		".github/copilot-instructions.md":      agentsContent,
		".windsurfrules":                       cursorRulesContent,
		filepath.Join(".gemini", "rules"):      agentsContent,
	}

	for path, content := range files {
		fullPath := filepath.Join(targetDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create dir for %s: %w", path, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
	}

	// Scaffolding: Makefile with Playwright & Dev-Loop targets
	makefilePath := filepath.Join(targetDir, "Makefile")
	makefileContent := `.PHONY: test build test-e2e dev-loop

test:
	@echo "🧪 Running Local Test-Centric Verification..."
	go test -v ./...

test-e2e:
	@echo "🎭 Running CI Headless Playwright Tests..."
	npx playwright test --headed=false

dev-loop:
	@echo "🔁 Running Goal-Driven Iterative Development Loop..."
	@make test

build:
	@echo "🔨 Building binary..."
	go build -o bin/app ./cmd/...
`
	_ = os.WriteFile(makefilePath, []byte(makefileContent), 0644)

	return nil
}

// GenerateGitOpsManifests produces Helm Charts & ArgoCD Application YAML
func GenerateGitOpsManifests(targetDir string, cfg *config.Config) error {
	gitopsDir := filepath.Join(targetDir, "gitops")
	helmDir := filepath.Join(gitopsDir, "helm", cfg.Step1.ProjectName)
	argocdDir := filepath.Join(gitopsDir, "argocd")

	if err := os.MkdirAll(helmDir, 0755); err != nil {
		return fmt.Errorf("failed to create helm dir: %w", err)
	}
	if err := os.MkdirAll(argocdDir, 0755); err != nil {
		return fmt.Errorf("failed to create argocd dir: %w", err)
	}

	// Chart.yaml
	chartContent := fmt.Sprintf(`apiVersion: v2
name: %s
description: Helm chart for %s
type: application
version: 0.1.0
appVersion: "1.0.0"
`, cfg.Step1.ProjectName, cfg.Step1.ProjectName)
	_ = os.WriteFile(filepath.Join(helmDir, "Chart.yaml"), []byte(chartContent), 0644)

	// values.yaml
	valuesContent := fmt.Sprintf(`replicaCount: 2
image:
  repository: %s/%s
  pullPolicy: IfNotPresent
  tag: "latest"
service:
  type: ClusterIP
  port: 8080
`, cfg.Step2.GitProvider, cfg.Step1.ProjectName)
	_ = os.WriteFile(filepath.Join(helmDir, "values.yaml"), []byte(valuesContent), 0644)

	// ArgoCD Application YAML
	argoContent := fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: %s-gitops
  namespace: argocd
spec:
  project: default
  source:
    repoURL: 'https://github.com/myeong-han/%s.git'
    targetRevision: HEAD
    path: gitops/helm/%s
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: default
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
`, cfg.Step1.ProjectName, cfg.Step1.ProjectName, cfg.Step1.ProjectName)
	_ = os.WriteFile(filepath.Join(argocdDir, "application.yaml"), []byte(argoContent), 0644)

	return nil
}

// GenerateAll executes the complete scaffolding pipeline
func GenerateAll(targetDir string, cfg *config.Config, prompt string) error {
	return GenerateHarnessProject(targetDir, cfg, prompt)
}
