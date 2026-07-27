package release

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

type CommitItem struct {
	Type    string
	Scope   string
	Message string
}

// BuildChangelogFromLog parses raw git log output into a structured Markdown CHANGELOG
func BuildChangelogFromLog(version string, rawGitLog string) string {
	lines := strings.Split(rawGitLog, "\n")

	var feats []string
	var fixes []string
	var refactors []string
	var docs []string
	var others []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, " ", 2)
		msg := trimmed
		if len(parts) == 2 {
			msg = parts[1]
		}

		lowerMsg := strings.ToLower(msg)

		switch {
		case strings.HasPrefix(lowerMsg, "feat"):
			feats = append(feats, msg)
		case strings.HasPrefix(lowerMsg, "fix"):
			fixes = append(fixes, msg)
		case strings.HasPrefix(lowerMsg, "refactor"):
			refactors = append(refactors, msg)
		case strings.HasPrefix(lowerMsg, "docs"):
			docs = append(docs, msg)
		default:
			others = append(others, msg)
		}
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("## [%s] - %s\n\n", version, time.Now().Format("2006-01-02")))

	if len(feats) > 0 {
		buf.WriteString("### 🚀 Features\n")
		for _, f := range feats {
			buf.WriteString(fmt.Sprintf("- %s\n", f))
		}
		buf.WriteString("\n")
	}

	if len(fixes) > 0 {
		buf.WriteString("### 🐛 Bug Fixes\n")
		for _, f := range fixes {
			buf.WriteString(fmt.Sprintf("- %s\n", f))
		}
		buf.WriteString("\n")
	}

	if len(refactors) > 0 {
		buf.WriteString("### 🛠️ Refactoring\n")
		for _, r := range refactors {
			buf.WriteString(fmt.Sprintf("- %s\n", r))
		}
		buf.WriteString("\n")
	}

	if len(docs) > 0 {
		buf.WriteString("### 📚 Documentation\n")
		for _, d := range docs {
			buf.WriteString(fmt.Sprintf("- %s\n", d))
		}
		buf.WriteString("\n")
	}

	if len(others) > 0 {
		buf.WriteString("### 📦 Maintenance & Other\n")
		for _, o := range others {
			buf.WriteString(fmt.Sprintf("- %s\n", o))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// WriteChangelogFile writes or prepends changelog content to CHANGELOG.md
func WriteChangelogFile(filePath string, version string, rawGitLog string) error {
	newContent := BuildChangelogFromLog(version, rawGitLog)

	var finalContent string
	if existing, err := os.ReadFile(filePath); err == nil {
		finalContent = newContent + "\n" + string(existing)
	} else {
		finalContent = "# CHANGELOG\n\n" + newContent
	}

	return os.WriteFile(filePath, []byte(finalContent), 0644)
}

// BuildSlackReleasePayload creates JSON payload for Slack release webhook notification
func BuildSlackReleasePayload(appName string, version string, changelog string) string {
	return fmt.Sprintf(`{
	"text": "🚀 *New Release Deployed: %s v%s*",
	"attachments": [
		{
			"color": "#00FFD1",
			"fields": [
				{
					"title": "Release Notes",
					"value": "%s",
					"short": false
				}
			]
		}
	]
}`, appName, version, strings.ReplaceAll(changelog, "\n", "\\n"))
}
