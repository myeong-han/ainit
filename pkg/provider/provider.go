package provider

// ModelSpec defines detailed capability specs for an LLM model (OpenCode compatible)
type ModelSpec struct {
	ID            string
	Name          string
	ContextWindow int    // e.g. 128000, 200000, 2000000
	MaxOutput     int    // e.g. 4096, 8192, 16384
	Description   string // e.g. "Most intelligent model", "Low latency"
	SupportsImages bool
}

// ProviderSpec defines LLM provider credentials and available model catalog
type ProviderSpec struct {
	ID             string
	Name           string
	DefaultAuth    string // "oauth", "apikey", "local"
	RequiresURL    bool   // true for Ollama / Custom Endpoints
	DefaultBaseURL string
	Models         []ModelSpec
}

// Catalog holds all supported LLM Providers and Models inspired by OpenCode
var catalog = []ProviderSpec{
	{
		ID:          "anthropic",
		Name:        "Anthropic",
		DefaultAuth: "oauth",
		Models: []ModelSpec{
			{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet", ContextWindow: 200000, MaxOutput: 8192, Description: "State-of-the-art coding & reasoning", SupportsImages: true},
			{ID: "claude-3-5-haiku", Name: "Claude 3.5 Haiku", ContextWindow: 200000, MaxOutput: 8192, Description: "Lightning fast responses", SupportsImages: false},
			{ID: "claude-3-opus-20240229", Name: "Claude 3 Opus", ContextWindow: 200000, MaxOutput: 4096, Description: "Complex analytical tasks", SupportsImages: true},
		},
	},
	{
		ID:          "openai",
		Name:        "OpenAI",
		DefaultAuth: "apikey",
		Models: []ModelSpec{
			{ID: "gpt-4o", Name: "GPT-4o", ContextWindow: 128000, MaxOutput: 16384, Description: "Flagship high-intelligence model", SupportsImages: true},
			{ID: "gpt-4o-mini", Name: "GPT-4o Mini", ContextWindow: 128000, MaxOutput: 16384, Description: "Affordable fast model", SupportsImages: true},
			{ID: "o3-mini", Name: "o3-mini", ContextWindow: 200000, MaxOutput: 100000, Description: "STEM & coding reasoning model", SupportsImages: false},
			{ID: "o1", Name: "o1", ContextWindow: 200000, MaxOutput: 100000, Description: "Deep reasoning model", SupportsImages: true},
		},
	},
	{
		ID:          "google",
		Name:        "Google Gemini",
		DefaultAuth: "apikey",
		Models: []ModelSpec{
			{ID: "gemini-1.5-pro", Name: "Gemini 1.5 Pro", ContextWindow: 2000000, MaxOutput: 8192, Description: "2M massive context window", SupportsImages: true},
			{ID: "gemini-1.5-flash", Name: "Gemini 1.5 Flash", ContextWindow: 1000000, MaxOutput: 8192, Description: "High speed multimodal model", SupportsImages: true},
			{ID: "gemini-2.0-flash-exp", Name: "Gemini 2.0 Flash Exp", ContextWindow: 1000000, MaxOutput: 8192, Description: "Next-gen experimental model", SupportsImages: true},
		},
	},
	{
		ID:          "deepseek",
		Name:        "DeepSeek",
		DefaultAuth: "apikey",
		Models: []ModelSpec{
			{ID: "deepseek-chat", Name: "DeepSeek V3 (Chat)", ContextWindow: 64000, MaxOutput: 8192, Description: "High performance open weights", SupportsImages: false},
			{ID: "deepseek-reasoner", Name: "DeepSeek R1 (Reasoner)", ContextWindow: 64000, MaxOutput: 8192, Description: "Open reasoning model", SupportsImages: false},
		},
	},
	{
		ID:          "openrouter",
		Name:        "OpenRouter",
		DefaultAuth: "apikey",
		Models: []ModelSpec{
			{ID: "openrouter/auto", Name: "OpenRouter Auto Router", ContextWindow: 128000, MaxOutput: 8192, Description: "Best value auto-routing", SupportsImages: true},
			{ID: "anthropic/claude-3.5-sonnet", Name: "Claude 3.5 Sonnet (OpenRouter)", ContextWindow: 200000, MaxOutput: 8192, Description: "Routed via OpenRouter API", SupportsImages: true},
		},
	},
	{
		ID:             "ollama",
		Name:           "Ollama (Local)",
		DefaultAuth:    "local",
		RequiresURL:    true,
		DefaultBaseURL: "http://localhost:11434",
		Models: []ModelSpec{
			{ID: "qwen2.5-coder:32b", Name: "Qwen 2.5 Coder 32B", ContextWindow: 32000, MaxOutput: 8192, Description: "Local coding powerhouse", SupportsImages: false},
			{ID: "llama3.3:70b", Name: "Llama 3.3 70B", ContextWindow: 128000, MaxOutput: 8192, Description: "Local open foundation model", SupportsImages: false},
			{ID: "deepseek-r1:14b", Name: "DeepSeek R1 14B", ContextWindow: 64000, MaxOutput: 8192, Description: "Local distil reasoning model", SupportsImages: false},
		},
	},
}

// GetAvailableProviders returns the full catalog of supported AI providers
func GetAvailableProviders() []ProviderSpec {
	return catalog
}

// GetModelsForProvider returns all supported models under a specific provider ID
func GetModelsForProvider(providerID string) []ModelSpec {
	for _, p := range catalog {
		if p.ID == providerID {
			return p.Models
		}
	}
	return nil
}

// GetProviderByID retrieves provider details by ID
func GetProviderByID(providerID string) *ProviderSpec {
	for _, p := range catalog {
		if p.ID == providerID {
			return &p
		}
	}
	return nil
}
