package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrategyConfigMarshalUsesNestedAIConfig(t *testing.T) {
	cfg := GetDefaultStrategyConfig("en")
	raw, err := json.Marshal(cfg)
	require.NoError(t, err)

	var shape map[string]any
	require.NoError(t, json.Unmarshal(raw, &shape))
	assert.Contains(t, shape, "ai_config")
	assert.NotContains(t, shape, "grid_config", "grid configuration is no longer part of the schema")
}

func TestStrategyConfigUnmarshalLegacyFlatLayout(t *testing.T) {
	raw := `{"coin_source":{"source_type":"static","static_coins":["aapl","msft"]},"risk_control":{"max_positions":2},"prompt_sections":{"role_definition":"x"}}`
	var cfg StrategyConfig
	require.NoError(t, json.Unmarshal([]byte(raw), &cfg))
	assert.Equal(t, []string{"aapl", "msft"}, cfg.CoinSource.StaticCoins)
	assert.Equal(t, 2, cfg.RiskControl.MaxPositions)
	assert.Equal(t, "x", cfg.PromptSections.RoleDefinition)
}

func TestStrategyConfigIgnoresRetiredFields(t *testing.T) {
	raw := `{"ai_config":{"coin_source":{"source_type":"vergex_signal","use_ai500":true},"risk_control":{"btc_eth_max_leverage":10}}}`
	var cfg StrategyConfig
	require.NoError(t, json.Unmarshal([]byte(raw), &cfg))
	cfg.NormalizeProductSchema()
	assert.Equal(t, "static", cfg.CoinSource.SourceType, "any retired source is coerced to a static list")
}

func TestNormalizeProductSchemaCleansSymbols(t *testing.T) {
	cfg := StrategyConfig{
		CoinSource: CoinSourceConfig{StaticCoins: []string{" aapl ", "AAPL", "", "600519.sh"}},
	}
	cfg.NormalizeProductSchema()
	assert.Equal(t, []string{"AAPL", "600519.SH"}, cfg.CoinSource.StaticCoins)
	assert.Equal(t, "1d", cfg.Indicators.Klines.PrimaryTimeframe)
	assert.Equal(t, []string{"1d"}, cfg.Indicators.Klines.SelectedTimeframes)
	assert.Equal(t, "en", cfg.Language)
}

func TestClampLimitsBoundsRiskSettings(t *testing.T) {
	cfg := GetDefaultStrategyConfig("en")
	cfg.RiskControl.MaxPositions = 99
	cfg.RiskControl.AltcoinMaxPositionValueRatio = 50
	cfg.RiskControl.MinConfidence = 5
	cfg.ClampLimits()

	assert.Equal(t, MaxPositions, cfg.RiskControl.MaxPositions)
	assert.Equal(t, MaxPositionValue, cfg.RiskControl.AltcoinMaxPositionValueRatio)
	assert.Equal(t, MinConfidence, cfg.RiskControl.MinConfidence)
}

func TestMergeStrategyConfigAppliesNestedPatch(t *testing.T) {
	base := GetDefaultStrategyConfig("en")
	merged, err := MergeStrategyConfig(base, map[string]any{
		"ai_config": map[string]any{
			"risk_control": map[string]any{"max_positions": 2},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 2, merged.RiskControl.MaxPositions)
	assert.Equal(t, base.CoinSource.StaticCoins, merged.CoinSource.StaticCoins, "untouched fields keep their values")
}

func TestStrategyClampWarningsReportChanges(t *testing.T) {
	before := GetDefaultStrategyConfig("en")
	after := before
	after.RiskControl.MaxPositions = 1
	warnings := StrategyClampWarnings(before, after, "en")
	require.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], "max_positions")
}
