package connection

import (
	"context"
	"testing"
)

func TestVerifyAIProviderConnectionWithApiKey(t *testing.T) {
	tester := NewConnectionTester()

	// Testing with dummy key should return auth error cleanly without panicking
	res := tester.TestAIProvider("openai", "invalid-api-key-test")
	if res.Connected {
		t.Error("expected invalid API key to fail connection test")
	}

	if res.ErrorMessage == "" {
		t.Error("expected error message for invalid API key")
	}
}

func TestVerifyGitProviderConnectionWithToken(t *testing.T) {
	tester := NewConnectionTester()

	ctx := context.Background()
	res := tester.TestGitProvider(ctx, "github", "ghp_invalid_token_test")
	if res.Connected {
		t.Error("expected invalid GitHub token to fail connection test")
	}
}

func TestVerifyDocSyncConnection(t *testing.T) {
	tester := NewConnectionTester()

	res := tester.TestDocSync("notion", "secret_invalid_notion_token")
	if res.Connected {
		t.Error("expected invalid Notion token to fail connection test")
	}
}
