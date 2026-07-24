package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitCheckResult struct {
	Installed bool
	Version   string
}

type RepoVerifyResult struct {
	Valid   bool
	Message string
}

type GitInitResult struct {
	Success bool
	Action  string // "cloned" or "initialized"
	WorkDir string
	Message string
}

func CheckGitInstallation() (*GitCheckResult, error) {
	cmd := exec.Command("git", "--version")
	out, err := cmd.Output()
	if err != nil {
		return &GitCheckResult{Installed: false, Version: ""}, err
	}

	verStr := strings.TrimSpace(string(out))
	return &GitCheckResult{
		Installed: true,
		Version:   verStr,
	}, nil
}

func VerifyRepoPath(provider, repoPath string) RepoVerifyResult {
	parts := strings.Split(strings.TrimSpace(repoPath), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return RepoVerifyResult{
			Valid:   false,
			Message: fmt.Sprintf("Invalid repo format '%s'. Must be in 'owner/repository' format.", repoPath),
		}
	}

	return RepoVerifyResult{
		Valid:   true,
		Message: fmt.Sprintf("Valid %s repository format: %s", strings.Title(provider), repoPath),
	}
}

// CheckRemoteExists executes git ls-remote to verify if remote repo exists
func CheckRemoteExists(remoteURL string) bool {
	cmd := exec.Command("git", "ls-remote", remoteURL)
	err := cmd.Run()
	return err == nil
}

// InitOrCloneRepository checks remote repository existence.
// If it exists, clones it and updates work-dir. Otherwise, initializes a new repository.
func InitOrCloneRepository(provider, repoPath, targetDir string) (*GitInitResult, error) {
	parts := strings.Split(strings.TrimSpace(repoPath), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("repoPath must be in 'owner/repo' format, got '%s'", repoPath)
	}

	remoteURL := fmt.Sprintf("https://github.com/%s.git", repoPath)
	if strings.ToLower(provider) == "bitbucket" {
		remoteURL = fmt.Sprintf("https://bitbucket.org/%s.git", repoPath)
	}

	exists := CheckRemoteExists(remoteURL)

	if exists {
		// Clone existing repository
		if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
			// Directory exists, pull latest
			cmd := exec.Command("git", "pull", "origin", "main")
			cmd.Dir = targetDir
			_ = cmd.Run()
			return &GitInitResult{
				Success: true,
				Action:  "pulled",
				WorkDir: targetDir,
				Message: fmt.Sprintf("Remote repository '%s' exists. Updated working directory at '%s'.", repoPath, targetDir),
			}, nil
		}

		cmd := exec.Command("git", "clone", remoteURL, targetDir)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to clone repository: %s, output: %s", err, string(out))
		}

		return &GitInitResult{
			Success: true,
			Action:  "cloned",
			WorkDir: targetDir,
			Message: fmt.Sprintf("Remote repository '%s' cloned successfully. Updated work-dir: %s", repoPath, targetDir),
		}, nil
	}

	// Remote repo does not exist: Initialize new repository locally & configure remote
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	initCmd := exec.Command("git", "init")
	initCmd.Dir = targetDir
	if err := initCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to git init: %v", err)
	}

	remoteAddCmd := exec.Command("git", "remote", "add", "origin", remoteURL)
	remoteAddCmd.Dir = targetDir
	_ = remoteAddCmd.Run()

	readmePath := filepath.Join(targetDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		_ = os.WriteFile(readmePath, []byte(fmt.Sprintf("# %s\nInitialized by Agentic-Init (`ainit`)\n", parts[1])), 0644)
	}

	return &GitInitResult{
		Success: true,
		Action:  "initialized",
		WorkDir: targetDir,
		Message: fmt.Sprintf("Remote '%s' not found. Initialized new local repository at '%s' with remote '%s'.", repoPath, targetDir, remoteURL),
	}, nil
}
