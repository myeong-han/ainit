package config

import (
	"errors"
	"strings"
)

// Step0Config defines AI Licensing and LLM setup
type Step0Config struct {
	LicensingMode string // "subscription", "apikey", "local"
	PrimaryModel  string // "claude-3-5-sonnet", "gpt-4o", "gemini-1.5-pro"
	FallbackModel string // "gpt-4o-mini", "claude-3-5-haiku"
	APIKey        string
}

// Step1Config defines Architecture & Specification setup
type Step1Config struct {
	ProjectName       string // e.g., "my-awesome-msa"
	ArchitectureStyle string // "msa", "monolith", "eda"
	RepoStructure     string // "monorepo", "multirepo"
	GenerateSequence  bool
	GenerateTraffic   bool
	GenerateGitOps    bool
	GenerateERD       bool
}

// Step2Config defines MCP Integrations setup
type Step2Config struct {
	GitProvider     string   // "github", "bitbucket"
	K8sTarget       string   // "local", "remote", "none"
	CICDTools       []string // "jenkins", "argocd", "harbor"
	DocTools        []string // "notion", "confluence"
	MessengerAlerts []string // "slack", "discord"
}

// Step3Config defines Git Init & Harness Coding setup
type Step3Config struct {
	CommitConvention string // "conventional", "gitmoji", "issue-prefix", "custom"
	IssueKeyFormat   string // e.g., "PROJ-\\d+"
	PRTemplateStyle  string // "standard", "minimal", "jira"
	AutoPRLabeling   bool
	TDDMode          bool
	LocalSandboxTest bool
}

// Step4Config defines Release & Operations setup
type Step4Config struct {
	VersioningStrategy string // "semver"
	AutoChangelog      bool
	ReleaseNotesSync   bool
	DeployAlert        bool
}

// Config is the root configuration holding all 5-step parameters
type Config struct {
	Step0 Step0Config
	Step1 Step1Config
	Step2 Step2Config
	Step3 Step3Config
	Step4 Step4Config
}

// NewDefaultConfig returns a Config initialized with sensible defaults
func NewDefaultConfig() *Config {
	return &Config{
		Step0: Step0Config{
			LicensingMode: "subscription",
			PrimaryModel:  "claude-3-5-sonnet",
			FallbackModel: "gpt-4o-mini",
		},
		Step1: Step1Config{
			ProjectName:       "ainit-project",
			ArchitectureStyle: "msa",
			RepoStructure:     "monorepo",
			GenerateSequence:  true,
			GenerateTraffic:   true,
			GenerateGitOps:    true,
			GenerateERD:       true,
		},
		Step2: Step2Config{
			GitProvider:     "github",
			K8sTarget:       "local",
			CICDTools:       []string{"argocd"},
			DocTools:        []string{"notion"},
			MessengerAlerts: []string{"slack"},
		},
		Step3: Step3Config{
			CommitConvention: "conventional",
			IssueKeyFormat:   `[A-Z]+-\d+`,
			PRTemplateStyle:  "standard",
			AutoPRLabeling:   true,
			TDDMode:          true,
			LocalSandboxTest: true,
		},
		Step4: Step4Config{
			VersioningStrategy: "semver",
			AutoChangelog:      true,
			ReleaseNotesSync:   true,
			DeployAlert:        true,
		},
	}
}

// Validate checks if the required fields in Config are properly populated
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Step1.ProjectName) == "" {
		return errors.New("project name cannot be empty")
	}
	if strings.TrimSpace(c.Step2.GitProvider) == "" {
		return errors.New("git provider is required")
	}
	return nil
}
