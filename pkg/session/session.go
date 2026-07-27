package session

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/myeong-han/ainit/pkg/config"
)

type ChatMessageItem struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type Session struct {
	ID          string    `json:"id"`
	ProjectName string    `json:"project_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SessionSummary struct {
	ID          string    `json:"id"`
	ProjectName string    `json:"project_name"`
	UpdatedAt   time.Time `json:"updated_at"`
	MsgCount    int       `json:"msg_count"`
}

type SessionManager struct {
	baseDir string
}

func NewSessionManager() *SessionManager {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = os.Getenv("HOME")
	}
	baseDir := filepath.Join(homeDir, ".ainit", "sessions")
	return NewSessionManagerWithDir(baseDir)
}

func NewSessionManagerWithDir(baseDir string) *SessionManager {
	_ = os.MkdirAll(baseDir, 0755)
	return &SessionManager{baseDir: baseDir}
}

func generateHashID(projName string) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%d", projName, time.Now().UnixNano())))
	fullHash := hex.EncodeToString(hasher.Sum(nil))
	return fullHash[:12]
}

// CreateSession initializes a new hash-based session directory
func (m *SessionManager) CreateSession(projName string) (*Session, error) {
	if projName == "" {
		projName = "unknown"
	}

	sessionID := generateHashID(projName)
	sessDir := filepath.Join(m.baseDir, sessionID)
	if err := os.MkdirAll(sessDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session dir: %w", err)
	}

	sess := &Session{
		ID:          sessionID,
		ProjectName: projName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	sessMetaBytes, _ := json.MarshalIndent(sess, "", "  ")
	_ = os.WriteFile(filepath.Join(sessDir, "session.json"), sessMetaBytes, 0644)

	// Save default config bound to this session
	defCfg := config.NewDefaultConfig()
	defCfg.Step1.ProjectName = projName
	cfgBytes, _ := json.MarshalIndent(defCfg, "", "  ")
	_ = os.WriteFile(filepath.Join(sessDir, "config.json"), cfgBytes, 0644)

	return sess, nil
}

// SaveSession persists session metadata, session-bound /settings config, and chat history
func (m *SessionManager) SaveSession(sessionID string, cfg *config.Config, history []ChatMessageItem) error {
	sessDir := filepath.Join(m.baseDir, sessionID)
	if err := os.MkdirAll(sessDir, 0755); err != nil {
		return fmt.Errorf("failed to ensure session dir: %w", err)
	}

	sess := &Session{
		ID:          sessionID,
		ProjectName: cfg.Step1.ProjectName,
		UpdatedAt:   time.Now(),
	}

	sessMetaBytes, _ := json.MarshalIndent(sess, "", "  ")
	if err := os.WriteFile(filepath.Join(sessDir, "session.json"), sessMetaBytes, 0644); err != nil {
		return fmt.Errorf("failed to write session.json: %w", err)
	}

	cfgBytes, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(filepath.Join(sessDir, "config.json"), cfgBytes, 0644); err != nil {
		return fmt.Errorf("failed to write config.json: %w", err)
	}

	histBytes, _ := json.MarshalIndent(history, "", "  ")
	if err := os.WriteFile(filepath.Join(sessDir, "history.json"), histBytes, 0644); err != nil {
		return fmt.Errorf("failed to write history.json: %w", err)
	}

	return nil
}

// LoadSession restores a session, its bound /settings config, and chat history
func (m *SessionManager) LoadSession(sessionID string) (*Session, *config.Config, []ChatMessageItem, error) {
	sessDir := filepath.Join(m.baseDir, sessionID)

	sessBytes, err := os.ReadFile(filepath.Join(sessDir, "session.json"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("session '%s' not found: %w", sessionID, err)
	}

	var sess Session
	if err := json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid session metadata: %w", err)
	}

	var cfg config.Config
	cfgBytes, err := os.ReadFile(filepath.Join(sessDir, "config.json"))
	if err == nil {
		_ = json.Unmarshal(cfgBytes, &cfg)
	} else {
		cfg = *config.NewDefaultConfig()
		cfg.Step1.ProjectName = sess.ProjectName
	}

	var history []ChatMessageItem
	histBytes, err := os.ReadFile(filepath.Join(sessDir, "history.json"))
	if err == nil {
		_ = json.Unmarshal(histBytes, &history)
	}

	return &sess, &cfg, history, nil
}

// ListSessions lists all saved sessions sorted by last updated timestamp
func (m *SessionManager) ListSessions() ([]SessionSummary, error) {
	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var summaries []SessionSummary
	for _, entry := range entries {
		if entry.IsDir() {
			sessID := entry.Name()
			sess, cfg, hist, err := m.LoadSession(sessID)
			if err == nil && sess != nil {
				proj := sess.ProjectName
				if cfg != nil && cfg.Step1.ProjectName != "" {
					proj = cfg.Step1.ProjectName
				}
				summaries = append(summaries, SessionSummary{
					ID:          sess.ID,
					ProjectName: proj,
					UpdatedAt:   sess.UpdatedAt,
					MsgCount:    len(hist),
				})
			}
		}
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].UpdatedAt.After(summaries[j].UpdatedAt)
	})

	return summaries, nil
}
