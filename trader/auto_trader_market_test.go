package trader

import (
	"testing"
	"time"

	"nofx/kernel"
	"nofx/trader/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Sessions used by the order-gate tests: 10:00 Beijing on a Wednesday, and
// 10:00 New York on the same date.
var (
	cnOpenTime = time.Date(2026, 10, 14, 2, 0, 0, 0, time.UTC)
	usOpenTime = time.Date(2026, 10, 14, 14, 0, 0, 0, time.UTC)
)

func newMarketTestTrader(t *testing.T, m types.Market) *AutoTrader {
	t.Helper()
	caps, err := types.CapabilitiesOf(m)
	require.NoError(t, err)
	return &AutoTrader{capabilities: caps}
}

func TestApplyMarketRulesCNEquity(t *testing.T) {
	at := newMarketTestTrader(t, types.MarketCNEquity)

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 1}, 350, false, cnOpenTime)
	require.NoError(t, err)
	assert.Equal(t, 300.0, qty, "buys are floored to whole lots")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 1}, 99, false, cnOpenTime)
	assert.Error(t, err, "less than one lot cannot be bought")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 1}, 300, true, cnOpenTime)
	assert.Error(t, err, "A-shares cannot be shorted")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 3}, 300, false, cnOpenTime)
	assert.Error(t, err, "A-shares cannot use leverage")
}

func TestApplyMarketRulesUSEquityKeepsFractionalShares(t *testing.T) {
	at := newMarketTestTrader(t, types.MarketUSEquity)

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "AAPL", Leverage: 1}, 1.5, false, usOpenTime)
	require.NoError(t, err)
	assert.Equal(t, 1.5, qty)

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "AAPL", Leverage: 1}, 2, true, usOpenTime)
	assert.NoError(t, err, "US equities may be shorted")
}

func TestApplyMarketRulesZeroCapabilitiesUseCryptoRules(t *testing.T) {
	at := &AutoTrader{}

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "BTCUSDT", Leverage: 10}, 0.123, true, cnOpenTime)
	require.NoError(t, err, "a trader without a market keeps legacy crypto behaviour")
	assert.Equal(t, 0.123, qty)
}

func TestEnforceSettlementCNBlocksSameDaySell(t *testing.T) {
	at := newMarketTestTrader(t, types.MarketCNEquity)
	bought := time.Date(2026, 10, 14, 2, 0, 0, 0, time.UTC) // 10:00 Beijing
	sameDay := time.Date(2026, 10, 14, 6, 30, 0, 0, time.UTC)
	nextDay := time.Date(2026, 10, 15, 2, 0, 0, 0, time.UTC)

	assert.Error(t, at.enforceSettlement(bought, sameDay))
	assert.NoError(t, at.enforceSettlement(bought, nextDay))
}

func TestEnforceSettlementUSAndCryptoAllowSameDay(t *testing.T) {
	bought := time.Date(2026, 10, 14, 2, 0, 0, 0, time.UTC)
	sameDay := time.Date(2026, 10, 14, 6, 30, 0, 0, time.UTC)

	assert.NoError(t, newMarketTestTrader(t, types.MarketUSEquity).enforceSettlement(bought, sameDay))
	assert.NoError(t, (&AutoTrader{}).enforceSettlement(bought, sameDay))
}

func TestApplyMarketRulesRejectsClosedSession(t *testing.T) {
	at := newMarketTestTrader(t, types.MarketCNEquity)
	closed := time.Date(2026, 10, 14, 4, 0, 0, 0, time.UTC) // 12:00 Beijing, lunch break

	_, err := at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH", Leverage: 1}, 100, false, closed)
	assert.Error(t, err, "orders during the lunch break are rejected")
}
