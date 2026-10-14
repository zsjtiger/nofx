package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstimateTokensGrowsWithSymbolsAndBars(t *testing.T) {
	small := GetDefaultStrategyConfig("en")
	small.CoinSource.StaticCoins = []string{"AAPL"}
	large := GetDefaultStrategyConfig("en")
	large.CoinSource.StaticCoins = []string{"AAPL", "MSFT", "NVDA", "AMZN", "GOOGL"}
	large.Indicators.Klines.PrimaryCount = 120

	assert.Greater(t, large.EstimateTokens().Total, small.EstimateTokens().Total)
}

func TestEstimateTokensReportsEveryModel(t *testing.T) {
	cfg := GetDefaultStrategyConfig("en")
	estimate := cfg.EstimateTokens()
	assert.NotEmpty(t, estimate.ModelLimits)
	for _, limit := range estimate.ModelLimits {
		assert.Greater(t, limit.ContextLimit, 0, limit.Name)
		assert.Contains(t, []string{"ok", "warning", "danger"}, limit.Level)
	}
}

func TestGetContextLimit(t *testing.T) {
	assert.Greater(t, GetContextLimit("deepseek"), 0)
	assert.Greater(t, GetContextLimit("unknown-provider"), 0, "unknown providers get a safe default")
}
