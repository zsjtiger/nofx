package kernel

import (
	"testing"

	"nofx/store"
	"nofx/trader/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func equityCaps(t *testing.T, m types.Market) types.Capabilities {
	t.Helper()
	caps, err := types.CapabilitiesOf(m)
	require.NoError(t, err)
	return caps
}

func testRisk() store.RiskControlConfig {
	return store.RiskControlConfig{MaxPositions: 3, AltcoinMaxPositionValueRatio: 0.5, MinConfidence: 70}
}

func TestEngineRequiresStockMarket(t *testing.T) {
	_, err := (&StrategyEngine{}).capabilities()
	assert.Error(t, err, "an engine without a market refuses to decide")

	_, err = (&StrategyEngine{market: "crypto"}).capabilities()
	assert.Error(t, err, "crypto is no longer a market")

	caps, err := (&StrategyEngine{}).WithMarket(types.MarketCNEquity).capabilities()
	require.NoError(t, err)
	assert.Equal(t, types.MarketCNEquity, caps.Market)
}

func TestCNEquityPromptStatesRules(t *testing.T) {
	prompt := BuildEquitySystemPrompt(equityCaps(t, types.MarketCNEquity), testRisk(), 100000, "Prefer liquid large caps.")

	assert.Contains(t, prompt, "A-share equity trader")
	assert.Contains(t, prompt, "Short selling is NOT allowed")
	assert.NotContains(t, prompt, "open_short |", "short actions are not offered as options")
	assert.Contains(t, prompt, "Leverage is NOT allowed")
	assert.Contains(t, prompt, "whole lots of 100 shares")
	assert.Contains(t, prompt, "T+1")
	assert.Contains(t, prompt, "CNY")
	assert.Contains(t, prompt, "At most 3 open positions")
	assert.Contains(t, prompt, "50000.00 CNY", "position cap is equity x ratio")
	assert.Contains(t, prompt, "confidence of at least 70")
	assert.Contains(t, prompt, "Prefer liquid large caps.")
}

func TestUSEquityPromptAllowsShorts(t *testing.T) {
	prompt := BuildEquitySystemPrompt(equityCaps(t, types.MarketUSEquity), testRisk(), 10000, "")
	assert.Contains(t, prompt, "Short selling is allowed")
	assert.Contains(t, prompt, "open_short")
	assert.NotContains(t, prompt, "T+1")
	assert.NotContains(t, prompt, "whole lots")
	assert.NotContains(t, prompt, "strategy owner")
}

func TestEngineBuildsEquityPromptForStockMarket(t *testing.T) {
	cfg := store.GetDefaultStrategyConfig("en")
	engine := NewStrategyEngine(&cfg).WithMarket(types.MarketUSEquity)
	prompt := engine.BuildSystemPrompt(10000, "")
	assert.Contains(t, prompt, "US equity trader")
}

func TestValidateEquityRejectsShortOnCN(t *testing.T) {
	caps := equityCaps(t, types.MarketCNEquity)
	decisions := []Decision{{Symbol: "600519.SH", Action: "open_short", Leverage: 1, PositionSizeUSD: 1000, StopLoss: 110}}
	err := ValidateEquityDecisions(decisions, caps, 100000, testRisk())
	assert.ErrorContains(t, err, "open_short is not allowed in cn_equity")

	decisions = []Decision{{Symbol: "600519.SH", Action: "close_short", Leverage: 1}}
	assert.Error(t, ValidateEquityDecisions(decisions, caps, 100000, testRisk()))
}

func TestValidateEquityClampsLeverageToOne(t *testing.T) {
	caps := equityCaps(t, types.MarketCNEquity)
	decisions := []Decision{{Symbol: "600519.SH", Action: "open_long", Leverage: 5, PositionSizeUSD: 10000, StopLoss: 95}}
	require.NoError(t, ValidateEquityDecisions(decisions, caps, 100000, testRisk()))
	assert.Equal(t, 1, decisions[0].Leverage)
}

func TestValidateEquityRequiresStopAndSize(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	noStop := []Decision{{Symbol: "AAPL", Action: "open_long", Leverage: 1, PositionSizeUSD: 1000}}
	assert.ErrorContains(t, ValidateEquityDecisions(noStop, caps, 10000, testRisk()), "stop_loss")

	noSize := []Decision{{Symbol: "AAPL", Action: "open_long", Leverage: 1, StopLoss: 90}}
	assert.ErrorContains(t, ValidateEquityDecisions(noSize, caps, 10000, testRisk()), "position_size_usd")
}

func TestValidateEquityCapsPositionValue(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	decisions := []Decision{{Symbol: "AAPL", Action: "open_long", Leverage: 1, PositionSizeUSD: 9000, StopLoss: 90}}
	require.NoError(t, ValidateEquityDecisions(decisions, caps, 10000, testRisk()))
	assert.InDelta(t, 5000, decisions[0].PositionSizeUSD, 0.001, "capped at equity x 0.5")
}

func TestValidateEquityAllowsUSShort(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	decisions := []Decision{{Symbol: "TSLA", Action: "open_short", Leverage: 1, PositionSizeUSD: 1000, StopLoss: 260}}
	assert.NoError(t, ValidateEquityDecisions(decisions, caps, 10000, testRisk()))
}

func TestValidateEquityHoldPassesThrough(t *testing.T) {
	caps := equityCaps(t, types.MarketCNEquity)
	decisions := []Decision{{Symbol: "600519.SH", Action: "hold", Leverage: 1}, {Symbol: "000001.SZ", Action: "wait"}}
	assert.NoError(t, ValidateEquityDecisions(decisions, caps, 100000, testRisk()))
}
