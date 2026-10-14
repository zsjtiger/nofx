package providers

import (
	"net/http"

	"crux/llm"
)

const (
	DefaultOpenAIBaseURL = "https://api.openai.com/v1"
	DefaultOpenAIModel   = "gpt-5.4"
)

func init() {
	llm.RegisterProvider(llm.ProviderOpenAI, func(opts ...llm.ClientOption) llm.AIClient {
		return NewOpenAIClientWithOptions(opts...)
	})
}

type OpenAIClient struct {
	*llm.Client
}

func (c *OpenAIClient) BaseClient() *llm.Client { return c.Client }

// NewOpenAIClient creates OpenAI client (backward compatible)
func NewOpenAIClient() llm.AIClient {
	return NewOpenAIClientWithOptions()
}

// NewOpenAIClientWithOptions creates OpenAI client (supports options pattern)
func NewOpenAIClientWithOptions(opts ...llm.ClientOption) llm.AIClient {
	openaiOpts := []llm.ClientOption{
		llm.WithProvider(llm.ProviderOpenAI),
		llm.WithModel(DefaultOpenAIModel),
		llm.WithBaseURL(DefaultOpenAIBaseURL),
	}

	allOpts := append(openaiOpts, opts...)
	baseClient := llm.NewClient(allOpts...).(*llm.Client)

	openaiClient := &OpenAIClient{
		Client: baseClient,
	}

	baseClient.Hooks = openaiClient
	return openaiClient
}

func (c *OpenAIClient) SetAPIKey(apiKey string, customURL string, customModel string) {
	c.APIKey = apiKey

	if len(apiKey) > 8 {
		c.Log.Infof("🔧 [MCP] OpenAI API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])
	}
	if customURL != "" {
		c.BaseURL = customURL
		c.Log.Infof("🔧 [MCP] OpenAI using custom BaseURL: %s", customURL)
	} else {
		c.Log.Infof("🔧 [MCP] OpenAI using default BaseURL: %s", c.BaseURL)
	}
	if customModel != "" {
		c.Model = customModel
		c.Log.Infof("🔧 [MCP] OpenAI using custom Model: %s", customModel)
	} else {
		c.Log.Infof("🔧 [MCP] OpenAI using default Model: %s", c.Model)
	}
}

// OpenAI uses standard Bearer auth
func (c *OpenAIClient) SetAuthHeader(reqHeaders http.Header) {
	c.Client.SetAuthHeader(reqHeaders)
}
