package trader

import (
	"testing"

	"nofx/kernel"
	"nofx/trader/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMarketTestTrader(t *testing.T, m types.Market) *AutoTrader {
	t.Helper()
	caps, err := types.CapabilitiesOf(m)
	require.NoError(t, err)
	return &AutoTrader{capabilities: caps}
}

func TestApplyMarketRulesCNEquity(t *testing.T) {
	at := newMarketTestTrader(t, types.MarketCNEquity)

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 1}, 350, false)
	require.NoError(t, err)
	assert.Equal(t, 300.0, qty, "buys are floored to whole lots")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 1}, 99, false)
	assert.Error(t, err, "less than one lot cannot be bought")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 1}, 300, true)
	assert.Error(t, err, "A-shares cannot be shorted")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 3}, 300, false)
	assert.Error(t, err, "A-shares cannot use leverage")
}

func TestApplyMarketRulesUSEquityKeepsFractionalShares(t *testing.T) {
	at := newMarketTestTrader(t, types.MarketUSEquity)

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "AAPL", Leverage: 1}, 1.5, false)
	require.NoError(t, err)
	assert.Equal(t, 1.5, qty)

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "AAPL", Leverage: 1}, 2, true)
	assert.NoError(t, err, "US equities may be shorted")
}

func TestApplyMarketRulesZeroCapabilitiesUseCryptoRules(t *testing.T) {
	at := &AutoTrader{}

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "BTCUSDT", Leverage: 10}, 0.123, true)
	require.NoError(t, err, "a trader without a market keeps legacy crypto behaviour")
	assert.Equal(t, 0.123, qty)
}
