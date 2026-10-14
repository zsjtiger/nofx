package trader

import (
	"testing"
	"time"

	"crux/kernel"
	"crux/trader/types"

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

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH"}, 350, false, cnOpenTime)
	require.NoError(t, err)
	assert.Equal(t, 300.0, qty, "buys are floored to whole lots")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH"}, 99, false, cnOpenTime)
	assert.Error(t, err, "less than one lot cannot be bought")

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH"}, 300, true, cnOpenTime)
	assert.Error(t, err, "A-shares cannot be shorted")

}

func TestApplyMarketRulesUSEquityKeepsFractionalShares(t *testing.T) {
	at := newMarketTestTrader(t, types.MarketUSEquity)

	qty, err := at.applyMarketRules(&kernel.Decision{Symbol: "AAPL"}, 1.5, false, usOpenTime)
	require.NoError(t, err)
	assert.Equal(t, 1.5, qty)

	_, err = at.applyMarketRules(&kernel.Decision{Symbol: "AAPL"}, 2, true, usOpenTime)
	assert.NoError(t, err, "US equities may be shorted")
}

func TestApplyMarketRulesWithoutMarketFailsClosed(t *testing.T) {
	at := &AutoTrader{}

	_, err := at.applyMarketRules(&kernel.Decision{Symbol: "AAPL"}, 1, true, usOpenTime)
	assert.Error(t, err, "a trader with no market allows no shorts")

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

	_, err := at.applyMarketRules(&kernel.Decision{Symbol: "600519.SH"}, 100, false, closed)
	assert.Error(t, err, "orders during the lunch break are rejected")
}
