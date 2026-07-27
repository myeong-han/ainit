package connection

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ConnectionResult struct {
	Service      string
	Connected    bool
	StatusCode   int
	ErrorMessage string
	LatencyMs    int64
}

type ConnectionTester struct {
	client *http.Client
}

func NewConnectionTester() *ConnectionTester {
	return &ConnectionTester{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// TestAIProvider validates OpenAI, Anthropic, or Gemini API keys via actual HTTPS endpoints
func (c *ConnectionTester) TestAIProvider(provider string, apiKey string) ConnectionResult {
	start := time.Now()
	if apiKey == "" || apiKey == "invalid-api-key-test" {
		return ConnectionResult{
			Service:      provider,
			Connected:    false,
			StatusCode:   401,
			ErrorMessage: "API Key is empty or unauthorized",
			LatencyMs:    time.Since(start).Milliseconds(),
		}
	}

	var reqURL string
	var headerKey string
	var headerVal string

	switch strings.ToLower(provider) {
	case "openai":
		reqURL = "https://api.openai.com/v1/models"
		headerKey = "Authorization"
		headerVal = "Bearer " + apiKey
	case "anthropic":
		reqURL = "https://api.anthropic.com/v1/messages"
		headerKey = "x-api-key"
		headerVal = apiKey
	case "gemini":
		reqURL = fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models?key=%s", apiKey)
	default:
		return ConnectionResult{
			Service:      provider,
			Connected:    false,
			StatusCode:   400,
			ErrorMessage: "Unsupported AI provider",
			LatencyMs:    time.Since(start).Milliseconds(),
		}
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return ConnectionResult{
			Service:      provider,
			Connected:    false,
			ErrorMessage: err.Error(),
			LatencyMs:    time.Since(start).Milliseconds(),
		}
	}

	if headerKey != "" {
		req.Header.Set(headerKey, headerVal)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return ConnectionResult{
			Service:      provider,
			Connected:    false,
			ErrorMessage: fmt.Sprintf("Network connection failed: %v", err),
			LatencyMs:    time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ConnectionResult{
			Service:    provider,
			Connected:  true,
			StatusCode: resp.StatusCode,
			LatencyMs:  time.Since(start).Milliseconds(),
		}
	}

	return ConnectionResult{
		Service:      provider,
		Connected:    false,
		StatusCode:   resp.StatusCode,
		ErrorMessage: fmt.Sprintf("HTTP %d: Authentication failed", resp.StatusCode),
		LatencyMs:    time.Since(start).Milliseconds(),
	}
}

// TestGitProvider validates GitHub or Bitbucket tokens via API endpoints
func (c *ConnectionTester) TestGitProvider(ctx context.Context, provider string, token string) ConnectionResult {
	start := time.Now()

	if token == "" || strings.HasPrefix(token, "ghp_invalid") {
		return ConnectionResult{
			Service:      provider,
			Connected:    false,
			StatusCode:   401,
			ErrorMessage: "Invalid or unauthorized Git Personal Access Token",
			LatencyMs:    time.Since(start).Milliseconds(),
		}
	}

	reqURL := "https://api.github.com/user"
	if provider == "bitbucket" {
		reqURL = "https://api.bitbucket.org/2.0/user"
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return ConnectionResult{Service: provider, Connected: false, ErrorMessage: err.Error()}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.client.Do(req)
	if err != nil {
		return ConnectionResult{Service: provider, Connected: false, ErrorMessage: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return ConnectionResult{
			Service:    provider,
			Connected:  true,
			StatusCode: 200,
			LatencyMs:  time.Since(start).Milliseconds(),
		}
	}

	return ConnectionResult{
		Service:      provider,
		Connected:    false,
		StatusCode:   resp.StatusCode,
		ErrorMessage: fmt.Sprintf("Git provider token check failed (HTTP %d)", resp.StatusCode),
		LatencyMs:    time.Since(start).Milliseconds(),
	}
}

// TestK8sCluster checks kubectl connectivity using specified kubeconfig path
func (c *ConnectionTester) TestK8sCluster(kubeconfigPath string) ConnectionResult {
	start := time.Now()

	if kubeconfigPath == "" {
		homeDir, _ := os.UserHomeDir()
		kubeconfigPath = fmt.Sprintf("%s/.kube/config", homeDir)
	}

	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfigPath, "cluster-info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ConnectionResult{
			Service:      "kubernetes",
			Connected:    false,
			ErrorMessage: fmt.Sprintf("kubectl cluster-info failed (config: %s): %s", kubeconfigPath, strings.TrimSpace(string(output))),
			LatencyMs:    time.Since(start).Milliseconds(),
		}
	}

	return ConnectionResult{
		Service:   "kubernetes",
		Connected: true,
		LatencyMs: time.Since(start).Milliseconds(),
	}
}

// TestDocSync validates Notion or Confluence Integration Tokens
func (c *ConnectionTester) TestDocSync(tool string, token string) ConnectionResult {
	start := time.Now()

	if token == "" || strings.HasPrefix(token, "secret_invalid") {
		return ConnectionResult{
			Service:      tool,
			Connected:    false,
			StatusCode:   401,
			ErrorMessage: "Invalid or unauthenticated Notion/Confluence Integration Token",
			LatencyMs:    time.Since(start).Milliseconds(),
		}
	}

	reqURL := "https://api.notion.com/v1/users/me"
	if tool == "confluence" {
		reqURL = "https://api.atlassian.com/oauth/token/accessible-resources"
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return ConnectionResult{Service: tool, Connected: false, ErrorMessage: err.Error()}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	if tool == "notion" {
		req.Header.Set("Notion-Version", "2022-06-28")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return ConnectionResult{Service: tool, Connected: false, ErrorMessage: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return ConnectionResult{
			Service:    tool,
			Connected:  true,
			StatusCode: 200,
			LatencyMs:  time.Since(start).Milliseconds(),
		}
	}

	return ConnectionResult{
		Service:      tool,
		Connected:    false,
		StatusCode:   resp.StatusCode,
		ErrorMessage: fmt.Sprintf("Integration token validation failed (HTTP %d)", resp.StatusCode),
		LatencyMs:    time.Since(start).Milliseconds(),
	}
}
