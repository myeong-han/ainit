package session

import (
	"strings"
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestSessionManagerCreationAndSaveLoadWithApiKey(t *testing.T) {
	tmpBaseDir := t.TempDir()
	mgr := NewSessionManagerWithDir(tmpBaseDir)

	sess, err := mgr.CreateSession("test-payment-app")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	cfg := config.NewDefaultConfig()
	cfg.Step1.ProjectName = "test-payment-app"
	cfg.Step0.ProviderID = "gemini"
	cfg.Step0.ApiKey = "sk-test-secret-api-key-12345"

	history := []ChatMessageItem{
		{Sender: "User", Content: "/settings"},
		{Sender: "Agent", Content: "API Key Registered"},
	}

	err = mgr.SaveSession(sess.ID, cfg, history)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Load Session
	loadedSess, loadedCfg, loadedHist, err := mgr.LoadSession(sess.ID)
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}

	if loadedSess.ID != sess.ID {
		t.Errorf("loaded session ID mismatch: expected %s, got %s", sess.ID, loadedSess.ID)
	}

	if loadedCfg.Step0.ApiKey != "sk-test-secret-api-key-12345" {
		t.Errorf("expected ApiKey to be restored as 'sk-test-secret-api-key-12345', got '%s'", loadedCfg.Step0.ApiKey)
	}

	if len(loadedHist) != 2 || !strings.Contains(loadedHist[1].Content, "Registered") {
		t.Errorf("loaded chat history mismatch, got %v", loadedHist)
	}
}
