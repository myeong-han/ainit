package git

import (
	"testing"
)

func TestCheckGitConnection(t *testing.T) {
	// Test git CLI availability or local repository status
	res, err := CheckGitInstallation()
	if err != nil {
		t.Fatalf("expected git CLI to be available, got error: %v", err)
	}

	if !res.Installed {
		t.Error("expected git to be installed locally")
	}
}

func TestVerifyRemoteRepository(t *testing.T) {
	// Test validating public repository format
	status := VerifyRepoPath("github", "myeong-han/ainit")
	if !status.Valid {
		t.Errorf("expected valid repo path, got invalid: %s", status.Message)
	}

	invalidStatus := VerifyRepoPath("github", "invalid-repo-format")
	if invalidStatus.Valid {
		t.Error("expected invalid repo path format for single string without slash")
	}
}
