package providers

import (
	"net/http"

	"nofx/llm"
)

const (
	DefaultGrokBaseURL = "https://api.x.ai/v1"
	DefaultGrokModel   = "grok-3-latest"
)

func init() {
	llm.RegisterProvider(llm.ProviderGrok, func(opts ...llm.ClientOption) llm.AIClient {
		return NewGrokClientWithOptions(opts...)
	})
}

type GrokClient struct {
	*llm.Client
}

func (c *GrokClient) BaseClient() *llm.Client { return c.Client }

// NewGrokClient creates Grok client (backward compatible)
func NewGrokClient() llm.AIClient {
	return NewGrokClientWithOptions()
}

// NewGrokClientWithOptions creates Grok client (supports options pattern)
func NewGrokClientWithOptions(opts ...llm.ClientOption) llm.AIClient {
	grokOpts := []llm.ClientOption{
		llm.WithProvider(llm.ProviderGrok),
		llm.WithModel(DefaultGrokModel),
		llm.WithBaseURL(DefaultGrokBaseURL),
	}

	allOpts := append(grokOpts, opts...)
	baseClient := llm.NewClient(allOpts...).(*llm.Client)

	grokClient := &GrokClient{
		Client: baseClient,
	}

	baseClient.Hooks = grokClient
	return grokClient
}

func (c *GrokClient) SetAPIKey(apiKey string, customURL string, customModel string) {
	c.APIKey = apiKey

	if len(apiKey) > 8 {
		c.Log.Infof("🔧 [MCP] Grok API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])
	}
	if customURL != "" {
		c.BaseURL = customURL
		c.Log.Infof("🔧 [MCP] Grok using custom BaseURL: %s", customURL)
	} else {
		c.Log.Infof("🔧 [MCP] Grok using default BaseURL: %s", c.BaseURL)
	}
	if customModel != "" {
		c.Model = customModel
		c.Log.Infof("🔧 [MCP] Grok using custom Model: %s", customModel)
	} else {
		c.Log.Infof("🔧 [MCP] Grok using default Model: %s", c.Model)
	}
}

// Grok uses standard OpenAI-compatible API with Bearer auth
func (c *GrokClient) SetAuthHeader(reqHeaders http.Header) {
	c.Client.SetAuthHeader(reqHeaders)
}
