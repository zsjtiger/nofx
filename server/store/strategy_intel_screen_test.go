package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrategyConfigRoundTripsIntelScreen(t *testing.T) {
	raw := `{"ai_config":{"coin_source":{"source_type":"static","static_coins":["600519.SH"],
		"intel_screen":{"enabled":true,"exclude_high_risk":true,"min_sentiment":-0.2,"lookback_hours":48}}}}`
	var cfg StrategyConfig
	require.NoError(t, json.Unmarshal([]byte(raw), &cfg))
	assert.True(t, cfg.CoinSource.IntelScreen.Enabled)
	assert.True(t, cfg.CoinSource.IntelScreen.ExcludeHighRisk)
	assert.InDelta(t, -0.2, cfg.CoinSource.IntelScreen.MinSentiment, 1e-9)
	assert.Equal(t, 48, cfg.CoinSource.IntelScreen.LookbackHours)

	out, err := json.Marshal(cfg)
	require.NoError(t, err)
	var again StrategyConfig
	require.NoError(t, json.Unmarshal(out, &again))
	assert.Equal(t, cfg.CoinSource.IntelScreen, again.CoinSource.IntelScreen)
}
