package providers

import (
	"net/http"

	"crux/llm"
)

const (
	DefaultGeminiBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
	DefaultGeminiModel   = "gemini-3.1-pro"
)

func init() {
	llm.RegisterProvider(llm.ProviderGemini, func(opts ...llm.ClientOption) llm.AIClient {
		return NewGeminiClientWithOptions(opts...)
	})
}

type GeminiClient struct {
	*llm.Client
}

func (c *GeminiClient) BaseClient() *llm.Client { return c.Client }

// NewGeminiClient creates Gemini client (backward compatible)
func NewGeminiClient() llm.AIClient {
	return NewGeminiClientWithOptions()
}

// NewGeminiClientWithOptions creates Gemini client (supports options pattern)
func NewGeminiClientWithOptions(opts ...llm.ClientOption) llm.AIClient {
	geminiOpts := []llm.ClientOption{
		llm.WithProvider(llm.ProviderGemini),
		llm.WithModel(DefaultGeminiModel),
		llm.WithBaseURL(DefaultGeminiBaseURL),
	}

	allOpts := append(geminiOpts, opts...)
	baseClient := llm.NewClient(allOpts...).(*llm.Client)

	geminiClient := &GeminiClient{
		Client: baseClient,
	}

	baseClient.Hooks = geminiClient
	return geminiClient
}

func (c *GeminiClient) SetAPIKey(apiKey string, customURL string, customModel string) {
	c.APIKey = apiKey

	if len(apiKey) > 8 {
		c.Log.Infof("🔧 [MCP] Gemini API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])
	}
	if customURL != "" {
		c.BaseURL = customURL
		c.Log.Infof("🔧 [MCP] Gemini using custom BaseURL: %s", customURL)
	} else {
		c.Log.Infof("🔧 [MCP] Gemini using default BaseURL: %s", c.BaseURL)
	}
	if customModel != "" {
		c.Model = customModel
		c.Log.Infof("🔧 [MCP] Gemini using custom Model: %s", customModel)
	} else {
		c.Log.Infof("🔧 [MCP] Gemini using default Model: %s", c.Model)
	}
}

// Gemini OpenAI-compatible API uses standard Bearer auth
func (c *GeminiClient) SetAuthHeader(reqHeaders http.Header) {
	c.Client.SetAuthHeader(reqHeaders)
}
