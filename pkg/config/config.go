package config

import (
	"errors"
	"os"
	"path/filepath"
)

type Step0Config struct {
	ProviderID    string `json:"provider_id"`
	PrimaryModel  string `json:"primary_model"`
	LicensingMode string `json:"licensing_mode"`
	FallbackModel string `json:"fallback_model"`
}

type Step1Config struct {
	ProjectName       string `json:"project_name"`
	ArchitectureStyle string `json:"architecture_style"`
	RepoStructure     string `json:"repo_structure"`
	GenerateSequence  bool   `json:"generate_sequence"`
	GenerateGitOps    bool   `json:"generate_gitops"`
}

type Step2Config struct {
	GitProvider    string `json:"git_provider"`
	K8sTarget      string `json:"k8s_target"`
	KubeconfigPath string `json:"kubeconfig_path"`
	CI             string `json:"ci"`
	CD             string `json:"cd"`
	Doc            string `json:"doc"`
	Messenger      string `json:"messenger"`
}

type Step3Config struct {
	CommitConvention string `json:"commit_convention"`
	PRTemplateStyle  string `json:"pr_template_style"`
	TDDMode          bool   `json:"tdd_mode"`
	LocalSandboxTest bool   `json:"local_sandbox_test"`
}

type Step4Config struct {
	VersioningStrategy string `json:"versioning_strategy"`
	AutoChangelog      bool   `json:"auto_changelog"`
	ReleaseNotesSync   bool   `json:"release_notes_sync"`
	DeployAlert        bool   `json:"deploy_alert"`
}

type Config struct {
	Step0 Step0Config `json:"step0"`
	Step1 Step1Config `json:"step1"`
	Step2 Step2Config `json:"step2"`
	Step3 Step3Config `json:"step3"`
	Step4 Step4Config `json:"step4"`
}

func NewDefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	return &Config{
		Step0: Step0Config{
			ProviderID:    "anthropic",
			PrimaryModel:  "claude-3-5-sonnet",
			LicensingMode: "subscription",
			FallbackModel: "claude-3-5-haiku",
		},
		Step1: Step1Config{
			ProjectName:       "unknown",
			ArchitectureStyle: "msa",
			RepoStructure:     "monorepo",
			GenerateSequence:  true,
			GenerateGitOps:    true,
		},
		Step2: Step2Config{
			GitProvider:    "github",
			K8sTarget:      "local",
			KubeconfigPath: filepath.Join(homeDir, ".kube", "config"),
			CI:             "jenkins",
			CD:             "argocd",
			Doc:            "notion",
			Messenger:      "slack",
		},
		Step3: Step3Config{
			CommitConvention: "conventional",
			PRTemplateStyle:  "standard",
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

func (c *Config) Validate() error {
	if c.Step0.ProviderID == "" {
		return errors.New("step0.provider_id cannot be empty")
	}
	if c.Step1.ProjectName == "" {
		return errors.New("step1.project_name cannot be empty")
	}
	return nil
}
