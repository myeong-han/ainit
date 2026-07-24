package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitCheckResult contains git CLI health information
type GitCheckResult struct {
	Installed bool
	Version   string
}

// RepoVerifyResult contains repository validation information
type RepoVerifyResult struct {
	Valid   bool
	Message string
}

// CheckGitInstallation executes git --version to verify local git CLI installation
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

// VerifyRepoPath checks whether the repository path matches 'owner/repo' format
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
