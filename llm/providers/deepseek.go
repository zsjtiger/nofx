package providers

import (
	"net/http"

	"nofx/llm"
)

func init() {
	llm.RegisterProvider(llm.ProviderDeepSeek, func(opts ...llm.ClientOption) llm.AIClient {
		return NewDeepSeekClientWithOptions(opts...)
	})
}

type DeepSeekClient struct {
	*llm.Client
}

func (c *DeepSeekClient) BaseClient() *llm.Client { return c.Client }

// NewDeepSeekClient creates DeepSeek client (backward compatible)
//
// Deprecated: Recommend using NewDeepSeekClientWithOptions for better flexibility
func NewDeepSeekClient() llm.AIClient {
	return NewDeepSeekClientWithOptions()
}

// NewDeepSeekClientWithOptions creates DeepSeek client (supports options pattern)
func NewDeepSeekClientWithOptions(opts ...llm.ClientOption) llm.AIClient {
	deepseekOpts := []llm.ClientOption{
		llm.WithProvider(llm.ProviderDeepSeek),
		llm.WithModel(llm.DefaultDeepSeekModel),
		llm.WithBaseURL(llm.DefaultDeepSeekBaseURL),
	}

	allOpts := append(deepseekOpts, opts...)
	baseClient := llm.NewClient(allOpts...).(*llm.Client)

	dsClient := &DeepSeekClient{
		Client: baseClient,
	}

	baseClient.Hooks = dsClient
	return dsClient
}

func (dsClient *DeepSeekClient) SetAPIKey(apiKey string, customURL string, customModel string) {
	dsClient.APIKey = apiKey

	if len(apiKey) > 8 {
		dsClient.Log.Infof("🔧 [MCP] DeepSeek API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])
	}
	if customURL != "" {
		dsClient.BaseURL = customURL
		dsClient.Log.Infof("🔧 [MCP] DeepSeek using custom BaseURL: %s", customURL)
	} else {
		dsClient.Log.Infof("🔧 [MCP] DeepSeek using default BaseURL: %s", dsClient.BaseURL)
	}
	if customModel != "" {
		dsClient.Model = customModel
		dsClient.Log.Infof("🔧 [MCP] DeepSeek using custom Model: %s", customModel)
	} else {
		dsClient.Log.Infof("🔧 [MCP] DeepSeek using default Model: %s", dsClient.Model)
	}
}

func (dsClient *DeepSeekClient) SetAuthHeader(reqHeaders http.Header) {
	dsClient.Client.SetAuthHeader(reqHeaders)
}
