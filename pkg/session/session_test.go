package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/myeong-han/ainit/pkg/config"
)

func TestSessionManagerCreationAndSaveLoad(t *testing.T) {
	tmpBaseDir := t.TempDir()
	mgr := NewSessionManagerWithDir(tmpBaseDir)

	// Create new session
	sess, err := mgr.CreateSession("test-payment-app")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	if sess.ID == "" || len(sess.ID) < 8 {
		t.Errorf("expected valid hash-based session ID, got '%s'", sess.ID)
	}

	sessionDir := filepath.Join(tmpBaseDir, sess.ID)
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		t.Errorf("expected session directory to exist at '%s'", sessionDir)
	}

	// Prepare config and chat history bound to this session
	cfg := config.NewDefaultConfig()
	cfg.Step1.ProjectName = "test-payment-app"
	cfg.Step0.ProviderID = "openai"

	history := []ChatMessageItem{
		{Sender: "User", Content: "/git-init test-payment-app"},
		{Sender: "Agent", Content: "Project Name set to test-payment-app"},
	}

	// Save Session
	err = mgr.SaveSession(sess.ID, cfg, history)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// List Sessions
	sessions, err := mgr.ListSessions()
	if err != nil {
		t.Fatalf("failed to list sessions: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session summary, got %d", len(sessions))
	}

	if sessions[0].ProjectName != "test-payment-app" {
		t.Errorf("expected project name 'test-payment-app', got '%s'", sessions[0].ProjectName)
	}

	// Load Session
	loadedSess, loadedCfg, loadedHist, err := mgr.LoadSession(sess.ID)
	if err != nil {
		t.Fatalf("failed to load session: %v", err)
	}

	if loadedSess.ID != sess.ID {
		t.Errorf("loaded session ID mismatch: expected %s, got %s", sess.ID, loadedSess.ID)
	}

	if loadedCfg.Step0.ProviderID != "openai" {
		t.Errorf("loaded config provider mismatch: expected openai, got %s", loadedCfg.Step0.ProviderID)
	}

	if len(loadedHist) != 2 || !strings.Contains(loadedHist[1].Content, "test-payment-app") {
		t.Errorf("loaded chat history mismatch, got %v", loadedHist)
	}
}
