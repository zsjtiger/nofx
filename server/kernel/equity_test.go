package kernel

import (
	"testing"

	"crux/store"
	"crux/trader/types"

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
	return store.RiskControlConfig{MaxPositions: 3, MaxPositionValueRatio: 0.5, MinConfidence: 70}
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
	decisions := []Decision{{Symbol: "600519.SH", Action: "open_short", PositionSizeUSD: 1000, StopLoss: 110, Confidence: 80}}
	err := ValidateEquityDecisions(decisions, caps, 100000, testRisk())
	assert.ErrorContains(t, err, "open_short is not allowed in cn_equity")

	decisions = []Decision{{Symbol: "600519.SH", Action: "close_short"}}
	assert.Error(t, ValidateEquityDecisions(decisions, caps, 100000, testRisk()))
}

func TestValidateEquityRequiresStopAndSize(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	noStop := []Decision{{Symbol: "AAPL", Action: "open_long", PositionSizeUSD: 1000, Confidence: 80}}
	assert.ErrorContains(t, ValidateEquityDecisions(noStop, caps, 10000, testRisk()), "stop_loss")

	noSize := []Decision{{Symbol: "AAPL", Action: "open_long", StopLoss: 90, Confidence: 80}}
	assert.ErrorContains(t, ValidateEquityDecisions(noSize, caps, 10000, testRisk()), "position_size_usd")
}

func TestValidateEquityCapsPositionValue(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	decisions := []Decision{{Symbol: "AAPL", Action: "open_long", PositionSizeUSD: 9000, StopLoss: 90, Confidence: 80}}
	require.NoError(t, ValidateEquityDecisions(decisions, caps, 10000, testRisk()))
	assert.InDelta(t, 5000, decisions[0].PositionSizeUSD, 0.001, "capped at equity x 0.5")
}

func TestValidateEquityAllowsUSShort(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	decisions := []Decision{{Symbol: "TSLA", Action: "open_short", PositionSizeUSD: 1000, StopLoss: 260, Confidence: 80}}
	assert.NoError(t, ValidateEquityDecisions(decisions, caps, 10000, testRisk()))
}

func TestValidateEquityHoldPassesThrough(t *testing.T) {
	caps := equityCaps(t, types.MarketCNEquity)
	decisions := []Decision{{Symbol: "600519.SH", Action: "hold"}, {Symbol: "000001.SZ", Action: "wait"}}
	assert.NoError(t, ValidateEquityDecisions(decisions, caps, 100000, testRisk()))
}

func TestValidateEquityLowConfidenceOpenBecomesWait(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	decisions := []Decision{{Symbol: "AAPL", Action: "open_long", PositionSizeUSD: 1000, StopLoss: 90, Confidence: 40}}
	require.NoError(t, ValidateEquityDecisions(decisions, caps, 10000, testRisk()))
	assert.Equal(t, "wait", decisions[0].Action, "confidence below the minimum must not open a position")
}

func TestValidateEquityMinConfidenceZeroAcceptsAny(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	risk := testRisk()
	risk.MinConfidence = 0
	decisions := []Decision{{Symbol: "AAPL", Action: "open_long", PositionSizeUSD: 1000, StopLoss: 90}}
	require.NoError(t, ValidateEquityDecisions(decisions, caps, 10000, risk))
	assert.Equal(t, "open_long", decisions[0].Action)
}

func TestEquityPromptStatesEnforcedRiskRules(t *testing.T) {
	caps := equityCaps(t, types.MarketUSEquity)
	risk := testRisk()
	risk.MinRiskRewardRatio = 2
	risk.MaxMarginUsage = 1
	risk.ProfitProtection = true
	prompt := BuildEquitySystemPrompt(caps, risk, 10000, "")
	assert.Contains(t, prompt, "at least 2.0 times the stop-loss distance")
	assert.Contains(t, prompt, "must not exceed 10000.00 USD")
	assert.Contains(t, prompt, "Profit protection is on")

	risk.ProfitProtection = false
	assert.NotContains(t, BuildEquitySystemPrompt(caps, risk, 10000, ""), "Profit protection is on")
}
