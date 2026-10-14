package store

import "strings"

// MissingRequiredExchangeCredentialFields lists the credential fields an exchange type still needs.
func MissingRequiredExchangeCredentialFields(exchangeType, apiKey, secretKey string) []string {
	switch strings.ToLower(strings.TrimSpace(exchangeType)) {
	case "alpaca_paper":
		return missingNamedFields(
			namedField{"api_key", apiKey},
			namedField{"secret_key", secretKey},
		)
	case "cn_paper":
		return missingNamedFields(
			namedField{"api_key (tushare token)", apiKey},
		)
	default:
		return []string{"exchange_type"}
	}
}

type namedField struct {
	name  string
	value string
}

func missingNamedFields(fields ...namedField) []string {
	var missing []string
	for _, f := range fields {
		if strings.TrimSpace(f.value) == "" {
			missing = append(missing, f.name)
		}
	}
	return missing
}

// IsVisibleAIModel reports whether a saved model should appear in the settings list.
func IsVisibleAIModel(model *AIModel) bool {
	if model == nil {
		return false
	}
	return model.Enabled ||
		strings.TrimSpace(string(model.APIKey)) != "" ||
		strings.TrimSpace(model.CustomAPIURL) != "" ||
		strings.TrimSpace(model.CustomModelName) != ""
}

// IsVisibleExchange reports whether a saved exchange account should appear in the list.
func IsVisibleExchange(exchange *Exchange) bool {
	if exchange == nil {
		return false
	}
	return exchange.Enabled ||
		strings.TrimSpace(string(exchange.APIKey)) != "" ||
		strings.TrimSpace(string(exchange.SecretKey)) != ""
}

func IsVisibleTrader(trader *Trader) bool {
	if trader == nil {
		return false
	}
	return strings.TrimSpace(trader.Name) != "" &&
		strings.TrimSpace(trader.AIModelID) != "" &&
		strings.TrimSpace(trader.ExchangeID) != ""
}

func IsVisibleStrategy(strategy *Strategy) bool {
	if strategy == nil {
		return false
	}
	return strings.TrimSpace(strategy.Name) != ""
}

// MarketForExchangeType returns the market an exchange type trades. Exchanges
// added before per-market support are perpetual crypto markets.

// MarketForExchangeType returns the market an exchange type trades, or an empty
// string for an unknown exchange type.
func MarketForExchangeType(exchangeType string) string {
	switch strings.ToLower(strings.TrimSpace(exchangeType)) {
	case "alpaca_paper":
		return "us_equity"
	case "cn_paper":
		return "cn_equity"
	default:
		return ""
	}
}
