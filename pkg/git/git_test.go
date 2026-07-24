package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckGitConnection(t *testing.T) {
	res, err := CheckGitInstallation()
	if err != nil {
		t.Fatalf("expected git CLI to be available, got error: %v", err)
	}

	if !res.Installed {
		t.Error("expected git to be installed locally")
	}
}

func TestVerifyRemoteRepository(t *testing.T) {
	status := VerifyRepoPath("github", "myeong-han/ainit")
	if !status.Valid {
		t.Errorf("expected valid repo path, got invalid: %s", status.Message)
	}
}

func TestInitOrCloneRepositoryExistingRepo(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "ainit-test")

	// Test initializing or cloning with a known public repository
	res, err := InitOrCloneRepository("github", "myeong-han/ainit", targetPath)
	if err != nil {
		t.Fatalf("expected no error executing InitOrCloneRepository, got: %v", err)
	}

	if !res.Success {
		t.Errorf("expected successful git-init/clone, got message: %s", res.Message)
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("expected target directory %s to be created", targetPath)
	}
}
