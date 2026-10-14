package llm_test

import (
	"fmt"
	"net/http"
	"time"

	"crux/llm"
	"crux/llm/providers"
)

// ============================================================
// Example 1: Basic Usage (Backward Compatible)
// ============================================================

func Example_backward_compatible() {
	// Old code continues to work without modification
	client := llm.New()
	client.SetAPIKey("sk-xxx", "https://api.custom.com", "gpt-4")

	// Usage
	result, _ := client.CallWithMessages("system prompt", "user prompt")
	fmt.Println(result)
}

func Example_deepseek_backward_compatible() {
	// DeepSeek old code continues to work
	client := providers.NewDeepSeekClient()
	client.SetAPIKey("sk-xxx", "", "")

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

// ============================================================
// Example 2: New Recommended Usage (Options Pattern)
// ============================================================

func Example_new_client_basic() {
	// Use default configuration
	client := llm.NewClient()

	// Use DeepSeek
	client = llm.NewClient(
		llm.WithDeepSeekConfig("sk-xxx"),
	)

	// Use Qwen
	client = llm.NewClient(
		llm.WithQwenConfig("sk-xxx"),
	)

	_ = client
}

func Example_new_client_with_options() {
	// Combine multiple options
	client := llm.NewClient(
		llm.WithDeepSeekConfig("sk-xxx"),
		llm.WithTimeout(60*time.Second),
		llm.WithMaxRetries(5),
		llm.WithMaxTokens(4000),
		llm.WithTemperature(0.7),
	)

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

// ============================================================
// Example 3: Custom Logger
// ============================================================

// CustomLogger custom logger example
type CustomLogger struct{}

func (l *CustomLogger) Debugf(format string, args ...any) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}

func (l *CustomLogger) Infof(format string, args ...any) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *CustomLogger) Warnf(format string, args ...any) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

func (l *CustomLogger) Errorf(format string, args ...any) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

func Example_custom_logger() {
	// Use custom logger
	customLogger := &CustomLogger{}

	client := llm.NewClient(
		llm.WithDeepSeekConfig("sk-xxx"),
		llm.WithLogger(customLogger),
	)

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

func Example_no_logger_for_testing() {
	// Disable logging during testing
	client := llm.NewClient(
		llm.WithLogger(llm.NewNoopLogger()),
	)

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

// ============================================================
// Example 4: Custom HTTP Client
// ============================================================

func Example_custom_http_client() {
	// Custom HTTP client (add proxy, TLS, etc.)
	customHTTP := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			// Custom TLS, connection pool, etc.
		},
	}

	client := llm.NewClient(
		llm.WithDeepSeekConfig("sk-xxx"),
		llm.WithHTTPClient(customHTTP),
	)

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

// ============================================================
// Example 5: DeepSeek Client (New API)
// ============================================================

func Example_deepseek_new_api() {
	// Basic usage
	client := providers.NewDeepSeekClientWithOptions(
		llm.WithAPIKey("sk-xxx"),
	)

	// Advanced usage
	client = providers.NewDeepSeekClientWithOptions(
		llm.WithAPIKey("sk-xxx"),
		llm.WithLogger(&CustomLogger{}),
		llm.WithTimeout(90*time.Second),
		llm.WithMaxTokens(8000),
	)

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

// ============================================================
// Example 6: Qwen Client (New API)
// ============================================================

func Example_qwen_new_api() {
	// Basic usage
	client := providers.NewQwenClientWithOptions(
		llm.WithAPIKey("sk-xxx"),
	)

	// Advanced usage
	client = providers.NewQwenClientWithOptions(
		llm.WithAPIKey("sk-xxx"),
		llm.WithLogger(&CustomLogger{}),
		llm.WithTimeout(90*time.Second),
	)

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

// ============================================================
// Example 7: Migration Example in trader/auto_trader.go
// ============================================================

func Example_trader_migration() {
	// Old code (continues to work)
	oldStyleClient := func(apiKey, customURL, customModel string) llm.AIClient {
		client := providers.NewDeepSeekClient()
		client.SetAPIKey(apiKey, customURL, customModel)
		return client
	}

	// New code (recommended)
	newStyleClient := func(apiKey, customURL, customModel string) llm.AIClient {
		opts := []llm.ClientOption{
			llm.WithAPIKey(apiKey),
		}

		if customURL != "" {
			opts = append(opts, llm.WithBaseURL(customURL))
		}

		if customModel != "" {
			opts = append(opts, llm.WithModel(customModel))
		}

		return providers.NewDeepSeekClientWithOptions(opts...)
	}

	// Both approaches work
	_ = oldStyleClient("sk-xxx", "", "")
	_ = newStyleClient("sk-xxx", "", "")
}

// ============================================================
// Example 8: Testing Scenarios
// ============================================================

// MockHTTPClient Mock HTTP client
type MockHTTPClient struct {
	Response string
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Return preset response
	return &http.Response{
		StatusCode: 200,
		Body:       nil, // Need to implement in actual tests
	}, nil
}

func Example_testing_with_mock() {
	client := llm.NewClient(
		llm.WithLogger(llm.NewNoopLogger()), // Disable logging
	)

	result, _ := client.CallWithMessages("system", "user")
	fmt.Println(result)
}

// ============================================================
// Example 9: Environment-Specific Configuration
// ============================================================

func Example_environment_specific() {
	// Development environment: detailed logging
	devClient := llm.NewClient(
		llm.WithDeepSeekConfig("sk-xxx"),
		llm.WithLogger(&CustomLogger{}), // Detailed logging
	)

	// Production environment: structured logging + timeout protection
	prodClient := llm.NewClient(
		llm.WithDeepSeekConfig("sk-xxx"),
		llm.WithTimeout(30*time.Second),
		llm.WithMaxRetries(3),
	)

	_, _ = devClient.CallWithMessages("system", "user")
	_, _ = prodClient.CallWithMessages("system", "user")
}

// ============================================================
// Example 10: Complete Real-World Example
// ============================================================

func Example_real_world_usage() {
	// Create client with complete configuration
	client := providers.NewDeepSeekClientWithOptions(
		llm.WithAPIKey("sk-xxxxxxxxxx"),
		llm.WithTimeout(60*time.Second),
		llm.WithMaxRetries(5),
		llm.WithMaxTokens(4000),
		llm.WithTemperature(0.5),
		llm.WithLogger(&CustomLogger{}),
	)

	// Use client
	systemPrompt := "You are a professional quantitative trading advisor"
	userPrompt := "Analyze current BTC trend"

	result, err := client.CallWithMessages(systemPrompt, userPrompt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("AI response: %s\n", result)
}
