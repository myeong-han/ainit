package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseConventionalCommitsAndBuildChangelog(t *testing.T) {
	rawGitLog := `45ae883 refactor(command): remove /set-name, set default App Name to unknown
db683ad fix(tui): clamp sidebar text lines and section dividers
01b2b98 feat(tui): dynamically filter slash commands and auto-focus
41c7054 docs: update README.md with /git-init <name>`

	changelog := BuildChangelogFromLog("1.0.0", rawGitLog)

	if !strings.Contains(changelog, "## [1.0.0]") {
		t.Errorf("expected changelog to contain version header, got:\n%s", changelog)
	}

	if !strings.Contains(changelog, "### 🚀 Features") || !strings.Contains(changelog, "dynamically filter slash commands") {
		t.Errorf("expected changelog to parse feature commits, got:\n%s", changelog)
	}

	if !strings.Contains(changelog, "### 🐛 Bug Fixes") || !strings.Contains(changelog, "clamp sidebar text lines") {
		t.Errorf("expected changelog to parse bug fix commits, got:\n%s", changelog)
	}

	// Test writing to CHANGELOG.md file
	tmpDir := t.TempDir()
	changelogFile := filepath.Join(tmpDir, "CHANGELOG.md")

	err := WriteChangelogFile(changelogFile, "1.0.0", rawGitLog)
	if err != nil {
		t.Fatalf("unexpected error writing changelog: %v", err)
	}

	data, err := os.ReadFile(changelogFile)
	if err != nil {
		t.Fatalf("failed to read written changelog file: %v", err)
	}

	if !strings.Contains(string(data), "## [1.0.0]") {
		t.Errorf("written changelog content mismatch:\n%s", string(data))
	}
}
