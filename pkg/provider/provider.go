package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ModelInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProviderInfo struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	DefaultAuth string      `json:"default_auth"`
	Models      []ModelInfo `json:"models"`
}

// GetAvailableProviders returns the list of supported OpenCode AI Providers with updated model catalogs
func GetAvailableProviders() []ProviderInfo {
	return []ProviderInfo{
		{
			ID:          "anthropic",
			Name:        "Anthropic Claude",
			Description: "Advanced reasoning & agentic capabilities",
			DefaultAuth: "subscription",
			Models: []ModelInfo{
				{ID: "claude-3-7-sonnet", Name: "Claude 3.7 Sonnet (Latest)", Description: "Flagship hybrid reasoning model"},
				{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet v2", Description: "High-precision code & architecture"},
				{ID: "claude-3-5-haiku-20241022", Name: "Claude 3.5 Haiku", Description: "Ultra-fast lightweight model"},
			},
		},
		{
			ID:          "openai",
			Name:        "OpenAI GPT",
			Description: "State-of-the-art multi-modal models",
			DefaultAuth: "apikey",
			Models: []ModelInfo{
				{ID: "gpt-4o", Name: "GPT-4o (Omni)", Description: "Versatile flagship model"},
				{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Description: "Fast & cost-effective lightweight model"},
				{ID: "o3-mini", Name: "o3-mini (Reasoning)", Description: "Next-gen reasoning engine"},
				{ID: "o1", Name: "o1 (Deep Thought)", Description: "Complex math & algorithm synthesis"},
			},
		},
		{
			ID:          "gemini",
			Name:        "Google Gemini",
			Description: "Long-context & high-speed multi-modal models",
			DefaultAuth: "apikey",
			Models: []ModelInfo{
				{ID: "gemini-2.0-flash", Name: "Gemini 2.0 Flash (Latest)", Description: "Next-gen real-time speed model"},
				{ID: "gemini-2.0-pro-exp", Name: "Gemini 2.0 Pro Experimental", Description: "Deep reasoning & architecture analysis"},
				{ID: "gemini-1.5-pro", Name: "Gemini 1.5 Pro", Description: "2M Token massive context window"},
			},
		},
		{
			ID:          "ollama",
			Name:        "Ollama (Local / Open-Source)",
			Description: "On-premise zero-data leakage local LLMs",
			DefaultAuth: "local",
			Models: []ModelInfo{
				{ID: "llama3.3:70b", Name: "Llama 3.3 70B", Description: "Local open-weights flagship"},
				{ID: "qwen2.5-coder:32b", Name: "Qwen 2.5 Coder 32B", Description: "Specialized local code assistant"},
				{ID: "deepseek-r1:32b", Name: "DeepSeek R1 32B", Description: "Local open-reasoning model"},
			},
		},
	}
}

func GetProviderByID(id string) *ProviderInfo {
	providers := GetAvailableProviders()
	for _, p := range providers {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

func GetModelsForProvider(providerID string) []ModelInfo {
	p := GetProviderByID(providerID)
	if p != nil {
		return p.Models
	}
	return nil
}

type openAIModelList struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type geminiModelList struct {
	Models []struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"models"`
}

type ollamaTagsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

// FetchLiveModelsForProvider dynamically fetches real-time available models from provider REST APIs
func FetchLiveModelsForProvider(ctx context.Context, providerID string, apiKey string) ([]ModelInfo, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	switch strings.ToLower(providerID) {
	case "openai":
		if apiKey == "" {
			return GetModelsForProvider("openai"), nil
		}
		req, err := http.NewRequestWithContext(ctx, "GET", "https://api.openai.com/v1/models", nil)
		if err != nil {
			return GetModelsForProvider("openai"), nil
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return GetModelsForProvider("openai"), nil
		}
		defer resp.Body.Close()

		var list openAIModelList
		if err := json.NewDecoder(resp.Body).Decode(&list); err == nil && len(list.Data) > 0 {
			var fetched []ModelInfo
			for _, m := range list.Data {
				if strings.HasPrefix(m.ID, "gpt-4") || strings.HasPrefix(m.ID, "o1") || strings.HasPrefix(m.ID, "o3") {
					fetched = append(fetched, ModelInfo{ID: m.ID, Name: m.ID + " (Live)", Description: "Fetched via OpenAI Live API"})
				}
			}
			if len(fetched) > 0 {
				return fetched, nil
			}
		}
		return GetModelsForProvider("openai"), nil

	case "gemini":
		if apiKey == "" {
			return GetModelsForProvider("gemini"), nil
		}
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models?key=%s", apiKey)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return GetModelsForProvider("gemini"), nil
		}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return GetModelsForProvider("gemini"), nil
		}
		defer resp.Body.Close()

		var list geminiModelList
		if err := json.NewDecoder(resp.Body).Decode(&list); err == nil && len(list.Models) > 0 {
			var fetched []ModelInfo
			for _, m := range list.Models {
				cleanID := strings.TrimPrefix(m.Name, "models/")
				fetched = append(fetched, ModelInfo{ID: cleanID, Name: m.DisplayName + " (Live)", Description: "Fetched via Gemini Live API"})
			}
			if len(fetched) > 0 {
				return fetched, nil
			}
		}
		return GetModelsForProvider("gemini"), nil

	case "ollama":
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:11434/api/tags", nil)
		if err != nil {
			return GetModelsForProvider("ollama"), nil
		}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return GetModelsForProvider("ollama"), nil
		}
		defer resp.Body.Close()

		var tags ollamaTagsResponse
		if err := json.NewDecoder(resp.Body).Decode(&tags); err == nil && len(tags.Models) > 0 {
			var fetched []ModelInfo
			for _, m := range tags.Models {
				fetched = append(fetched, ModelInfo{ID: m.Name, Name: m.Name + " (Local)", Description: "Fetched via Ollama Daemon API"})
			}
			if len(fetched) > 0 {
				return fetched, nil
			}
		}
		return GetModelsForProvider("ollama"), nil

	default:
		return GetModelsForProvider(providerID), nil
	}
}
