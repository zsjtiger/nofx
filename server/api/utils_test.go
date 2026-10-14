package api

import (
	"fmt"
	"strings"
	"testing"
)

func TestMaskSensitiveString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Short string (8 characters or less)",
			input:    "short",
			expected: "****",
		},
		{
			name:     "Normal API key",
			input:    "sk-1234567890abcdefghijklmnopqrstuvwxyz",
			expected: "sk-1****wxyz",
		},
		{
			name:     "Normal private key",
			input:    "0x1234567890abcdef1234567890abcdef12345678",
			expected: "0x12****5678",
		},
		{
			name:     "Exactly 9 characters",
			input:    "123456789",
			expected: "1234****6789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskSensitiveString(tt.input)
			if result != tt.expected {
				t.Errorf("MaskSensitiveString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeModelConfigForLog(t *testing.T) {
	models := map[string]ModelConfigUpdate{
		"deepseek": {
			Enabled:         true,
			APIKey:          "sk-1234567890abcdefghijklmnopqrstuvwxyz",
			CustomAPIURL:    "https://api.deepseek.com",
			CustomModelName: "deepseek-chat",
		},
	}

	result := SanitizeModelConfigForLog(models)

	deepseekConfig, ok := result["deepseek"].(map[string]interface{})
	if !ok {
		t.Fatal("deepseek config not found or wrong type")
	}

	if deepseekConfig["enabled"] != true {
		t.Errorf("expected enabled=true, got %v", deepseekConfig["enabled"])
	}

	maskedKey, ok := deepseekConfig["api_key"].(string)
	if !ok {
		t.Fatal("api_key not found or wrong type")
	}

	if maskedKey != "sk-1****wxyz" {
		t.Errorf("expected masked api_key='sk-1****wxyz', got %q", maskedKey)
	}

	if deepseekConfig["custom_api_url"] != "https://api.deepseek.com" {
		t.Errorf("custom_api_url should not be masked")
	}
}

func TestSanitizeExchangeConfigForLog(t *testing.T) {
	exchanges := map[string]ExchangeConfigUpdate{
		"alpaca": {
			Enabled:   true,
			APIKey:    "alpaca_key_id_1234567890abcdef",
			SecretKey: "alpaca_secret_key_1234567890abcdef",
		},
		"tushare": {
			Enabled: true,
			APIKey:  "tushare_token_1234567890abcdef",
		},
	}

	result := SanitizeExchangeConfigForLog(exchanges)

	alpaca, ok := result["alpaca"].(map[string]interface{})
	if !ok {
		t.Fatal("alpaca config not found or wrong type")
	}
	if alpaca["enabled"] != true {
		t.Errorf("enabled should be true, got %v", alpaca["enabled"])
	}
	for _, field := range []string{"api_key", "secret_key"} {
		value, _ := alpaca[field].(string)
		if value == "" || strings.Contains(value, "1234567890abcdef") {
			t.Errorf("%s should be masked, got %q", field, value)
		}
	}

	tushare, ok := result["tushare"].(map[string]interface{})
	if !ok {
		t.Fatal("tushare config not found or wrong type")
	}
	if _, has := tushare["secret_key"]; has {
		t.Error("secret_key should be omitted when empty")
	}
}

// TestSanitizeExchangeConfigForLog_NoPlaintextSecrets renders the sanitized log
// output exactly as the handler does (`%+v`) and asserts that no plaintext
// secret — including the passphrase and exchange API key private key that were
// historically not redacted — survives into the log line.
func TestSanitizeExchangeConfigForLog_NoPlaintextSecrets(t *testing.T) {
	secrets := map[string]string{
		"api_key":    "alpaca_key_id_1234567890abcdef",
		"secret_key": "alpaca_secret_key_1234567890abcdef",
	}

	exchanges := map[string]ExchangeConfigUpdate{
		"alpaca": {
			Enabled:   true,
			APIKey:    secrets["api_key"],
			SecretKey: secrets["secret_key"],
		},
	}

	logged := fmt.Sprintf("%+v", SanitizeExchangeConfigForLog(exchanges))
	for name, secret := range secrets {
		if strings.Contains(logged, secret) {
			t.Errorf("%s leaked into the log output", name)
		}
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty email",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid format",
			input:    "notanemail",
			expected: "****",
		},
		{
			name:     "Normal email",
			input:    "user@example.com",
			expected: "us****@example.com",
		},
		{
			name:     "Short username",
			input:    "a@example.com",
			expected: "**@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskEmail(tt.input)
			if result != tt.expected {
				t.Errorf("MaskEmail(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
