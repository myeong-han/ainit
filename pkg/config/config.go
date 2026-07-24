package config

import (
	"errors"
	"strings"
)

type Step0Config struct {
	LicensingMode string // "subscription" | "apikey" | "local"
	PrimaryModel  string // "claude-3-5-sonnet", "gpt-4o", etc.
	FallbackModel string // "gpt-4o-mini", "claude-3-5-haiku", etc.
	ProviderID    string // "anthropic", "openai", "gemini", "deepseek", "openrouter", "ollama"
}

type Step1Config struct {
	ProjectName       string // Defaults to "unknown"
	ArchitectureStyle string // "msa" | "monolith" | "eda"
	RepoStructure     string // "monorepo" | "multirepo"
	GenerateSequence  bool
	GenerateGitOps    bool
}

type Step2Config struct {
	GitProvider string // "github" | "bitbucket"
	K8sTarget   string // "local" | "remote" | "none"
	CICDTools   []string
	DocTools    []string
}

type Step3Config struct {
	CommitConvention string // "conventional" | "gitmoji" | "issue-prefix" | "custom"
	PRTemplateStyle  string // "standard" | "minimal" | "jira"
	TDDMode          bool
	LocalSandboxTest bool
}

type Step4Config struct {
	AutoChangelog      bool
	ReleaseNotesSync   bool
	DeployAlert        bool
	VersioningStrategy string // "semver" | "calver"
}

type Config struct {
	Step0 Step0Config
	Step1 Step1Config
	Step2 Step2Config
	Step3 Step3Config
	Step4 Step4Config
}

func NewDefaultConfig() *Config {
	return &Config{
		Step0: Step0Config{
			LicensingMode: "subscription",
			PrimaryModel:  "claude-3-5-sonnet",
			FallbackModel: "gpt-4o-mini",
			ProviderID:    "anthropic",
		},
		Step1: Step1Config{
			ProjectName:       "unknown", // Initial loading default is "unknown"
			ArchitectureStyle: "msa",
			RepoStructure:     "monorepo",
			GenerateSequence:  true,
			GenerateGitOps:    true,
		},
		Step2: Step2Config{
			GitProvider: "github",
			K8sTarget:   "local",
			CICDTools:   []string{"argo-cd", "jenkins"},
			DocTools:    []string{"notion", "confluence"},
		},
		Step3: Step3Config{
			CommitConvention: "conventional",
			PRTemplateStyle:  "standard",
			TDDMode:          true,
			LocalSandboxTest: true,
		},
		Step4: Step4Config{
			AutoChangelog:      true,
			ReleaseNotesSync:   true,
			DeployAlert:        true,
			VersioningStrategy: "semver",
		},
	}
}

func (c *Config) Validate() error {
	if c.Step0.ProviderID == "" {
		return errors.New("Step 0: ProviderID cannot be empty")
	}
	if c.Step1.ProjectName == "" {
		c.Step1.ProjectName = "unknown"
	}
	if !strings.EqualFold(c.Step1.ArchitectureStyle, "msa") &&
		!strings.EqualFold(c.Step1.ArchitectureStyle, "monolith") &&
		!strings.EqualFold(c.Step1.ArchitectureStyle, "eda") {
		return errors.New("Step 1: Invalid architecture style")
	}
	return nil
}
