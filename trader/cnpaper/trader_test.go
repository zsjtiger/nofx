package cnpaper

import (
	"context"
	"errors"
	"testing"
	"time"

	"nofx/trader/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeQuoter returns a fixed quote per symbol.
type fakeQuoter map[string]Quote

func (f fakeQuoter) Quote(_ context.Context, symbol string) (Quote, error) {
	q, ok := f[symbol]
	if !ok {
		return Quote{}, errors.New("no quote for " + symbol)
	}
	return q, nil
}

var cst = time.FixedZone("CST", 8*3600)

// newTestTrader returns a trader whose clock is set to the given Beijing time.
func newTestTrader(t *testing.T, now time.Time, quotes fakeQuoter) *Trader {
	t.Helper()
	tr, err := New(1_000_000, quotes)
	require.NoError(t, err)
	tr.now = func() time.Time { return now }
	return tr
}

var (
	// 2026-10-14 10:00 Beijing, a trading day inside the morning session.
	morning = time.Date(2026, 10, 14, 10, 0, 0, 0, cst)
	// 2026-10-15 10:00 Beijing, the next trading day.
	nextDay = time.Date(2026, 10, 15, 10, 0, 0, 0, cst)

	normalQuote = Quote{Price: 100, PrevClose: 100, UpLimit: 110, DownLimit: 90, Volume: 5000}
)

func TestNewValidatesArguments(t *testing.T) {
	_, err := New(0, fakeQuoter{})
	assert.Error(t, err)
	_, err = New(100, nil)
	assert.Error(t, err)
}

func TestOpenLongFillsWholeLotsAndChargesFees(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": normalQuote})

	res, err := tr.OpenLong("600519.SH", 300, 1)
	require.NoError(t, err)
	assert.Equal(t, "FILLED", res["status"])
	assert.Equal(t, 300.0, res["executedQty"])

	// 300 * 100 = 30,000 plus commission at 0.025% = 7.5, above the 5 CNY minimum.
	bal, err := tr.GetBalance()
	require.NoError(t, err)
	assert.InDelta(t, 1_000_000-30_000-7.5, bal["totalWalletBalance"], 0.001)

	pos, err := tr.GetPositions()
	require.NoError(t, err)
	require.Len(t, pos, 1)
	assert.Equal(t, 300.0, pos[0]["positionAmt"])
	assert.Equal(t, "long", pos[0]["side"])
	assert.Equal(t, int64(morning.UnixMilli()), pos[0]["createdTime"])
}

func TestOpenLongRejectsOddLotAndLeverage(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": normalQuote})

	_, err := tr.OpenLong("600519.SH", 150, 1)
	assert.ErrorContains(t, err, "multiple of 100")

	_, err = tr.OpenLong("600519.SH", 100, 3)
	assert.ErrorContains(t, err, "not allowed")
}

func TestShortingIsRejected(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": normalQuote})

	_, err := tr.OpenShort("600519.SH", 100, 1)
	assert.ErrorContains(t, err, "short selling is not allowed")

	_, err = tr.CloseShort("600519.SH", 0)
	assert.Error(t, err)
}

func TestOpenLongBlockedAtLimitUpAndWhenSuspended(t *testing.T) {
	limitUp := Quote{Price: 110, PrevClose: 100, UpLimit: 110, DownLimit: 90, Volume: 5000}
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": limitUp})
	_, err := tr.OpenLong("600519.SH", 100, 1)
	assert.ErrorContains(t, err, "limit-up")

	suspended := Quote{Price: 100, PrevClose: 100, UpLimit: 110, DownLimit: 90, Volume: 0}
	tr2 := newTestTrader(t, morning, fakeQuoter{"600519.SH": suspended})
	_, err = tr2.OpenLong("600519.SH", 100, 1)
	assert.ErrorContains(t, err, "suspended")
}

func TestOpenLongRejectsInsufficientCash(t *testing.T) {
	tr, err := New(1000, fakeQuoter{"600519.SH": normalQuote})
	require.NoError(t, err)
	tr.now = func() time.Time { return morning }
	_, err = tr.OpenLong("600519.SH", 100, 1)
	assert.ErrorContains(t, err, "insufficient cash")
}

func TestSameDaySellBlockedByT1(t *testing.T) {
	quotes := fakeQuoter{"600519.SH": normalQuote}
	tr := newTestTrader(t, morning, quotes)
	_, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)

	_, err = tr.CloseLong("600519.SH", 0)
	assert.ErrorContains(t, err, "T+1")

	tr.now = func() time.Time { return nextDay }
	_, err = tr.CloseLong("600519.SH", 0)
	assert.NoError(t, err, "selling on the next trading day is allowed")
}

func TestCloseLongRecordsRealizedPnLAndStampDuty(t *testing.T) {
	quotes := fakeQuoter{"600519.SH": normalQuote}
	tr := newTestTrader(t, morning, quotes)
	_, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)

	tr.now = func() time.Time { return nextDay }
	quotes["600519.SH"] = Quote{Price: 105, PrevClose: 100, UpLimit: 110, DownLimit: 90, Volume: 5000}
	_, err = tr.CloseLong("600519.SH", 0)
	require.NoError(t, err)

	closed, err := tr.GetClosedPnL(time.Time{}, 10)
	require.NoError(t, err)
	require.Len(t, closed, 1)
	assert.Equal(t, "manual", closed[0].CloseType)
	assert.InDelta(t, 105, closed[0].ExitPrice, 0.001)
	// Sale 10,500; commission 5 (minimum) plus stamp duty 5.25 => proceeds 10,489.75.
	// Cost 10,000 + buy commission 5 => realised P&L 484.75.
	assert.InDelta(t, 484.75, closed[0].RealizedPnL, 0.01)

	pos, err := tr.GetPositions()
	require.NoError(t, err)
	assert.Empty(t, pos)
}

func TestCloseLongRejectsOverSelling(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": normalQuote})
	_, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)

	tr.now = func() time.Time { return nextDay }
	_, err = tr.CloseLong("600519.SH", 200)
	assert.ErrorContains(t, err, "cannot sell")
}

func TestPartialSellLeavesOddLotRemainder(t *testing.T) {
	quotes := fakeQuoter{"600519.SH": normalQuote}
	tr := newTestTrader(t, morning, quotes)
	_, err := tr.OpenLong("600519.SH", 300, 1)
	require.NoError(t, err)

	tr.now = func() time.Time { return nextDay }
	_, err = tr.CloseLong("600519.SH", 150)
	assert.ErrorContains(t, err, "multiple of 100", "partial sells must be whole lots")

	_, err = tr.CloseLong("600519.SH", 200)
	require.NoError(t, err)
	pos, _ := tr.GetPositions()
	require.Len(t, pos, 1)
	assert.Equal(t, 100.0, pos[0]["positionAmt"])
}

func TestCloseSellAtLimitDownBlocked(t *testing.T) {
	quotes := fakeQuoter{"600519.SH": normalQuote}
	tr := newTestTrader(t, morning, quotes)
	_, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)

	tr.now = func() time.Time { return nextDay }
	quotes["600519.SH"] = Quote{Price: 90, PrevClose: 100, UpLimit: 110, DownLimit: 90, Volume: 5000}
	_, err = tr.CloseLong("600519.SH", 0)
	assert.ErrorContains(t, err, "limit-down")
}

func TestStopLossTriggersOnRead(t *testing.T) {
	quotes := fakeQuoter{"600519.SH": normalQuote}
	tr := newTestTrader(t, morning, quotes)
	_, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)
	require.NoError(t, tr.SetStopLoss("600519.SH", "LONG", 100, 95))

	orders, err := tr.GetOpenOrders("600519.SH")
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, "STOP_MARKET", orders[0].Type)

	// Price falls through the stop on the next trading day.
	tr.now = func() time.Time { return nextDay }
	quotes["600519.SH"] = Quote{Price: 94, PrevClose: 100, UpLimit: 110, DownLimit: 90, Volume: 5000}
	pos, err := tr.GetPositions()
	require.NoError(t, err)
	assert.Empty(t, pos, "stop-loss closes the position")

	closed, _ := tr.GetClosedPnL(time.Time{}, 10)
	require.Len(t, closed, 1)
	assert.Equal(t, "stop_loss", closed[0].CloseType)
}

func TestStopBlockedByT1StaysPending(t *testing.T) {
	quotes := fakeQuoter{"600519.SH": normalQuote}
	tr := newTestTrader(t, morning, quotes)
	_, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)
	require.NoError(t, tr.SetStopLoss("600519.SH", "LONG", 100, 95))

	// Still the same day: the stop is hit but T+1 blocks the sale.
	quotes["600519.SH"] = Quote{Price: 94, PrevClose: 100, UpLimit: 110, DownLimit: 90, Volume: 5000}
	pos, err := tr.GetPositions()
	require.NoError(t, err)
	assert.Len(t, pos, 1, "position stays open until it can be sold")
}

func TestTakeProfitAndCancel(t *testing.T) {
	quotes := fakeQuoter{"600519.SH": normalQuote}
	tr := newTestTrader(t, morning, quotes)
	_, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)
	require.NoError(t, tr.SetTakeProfit("600519.SH", "LONG", 100, 108))

	require.NoError(t, tr.CancelTakeProfitOrders("600519.SH"))
	orders, _ := tr.GetOpenOrders("600519.SH")
	assert.Empty(t, orders)
}

func TestProtectionRequiresOpenPosition(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": normalQuote})
	assert.Error(t, tr.SetStopLoss("600519.SH", "LONG", 100, 95))
	assert.Error(t, tr.SetTakeProfit("600519.SH", "LONG", 100, 110))
}

func TestGetOrderStatusAfterFill(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": normalQuote})
	res, err := tr.OpenLong("600519.SH", 100, 1)
	require.NoError(t, err)

	st, err := tr.GetOrderStatus("600519.SH", res["orderId"].(string))
	require.NoError(t, err)
	assert.Equal(t, "FILLED", st["status"])
	assert.Equal(t, 100.0, st["executedQty"])

	_, err = tr.GetOrderStatus("600519.SH", "missing")
	assert.Error(t, err)
}

func TestSessionAndFormatHelpers(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{"600519.SH": normalQuote})

	qty, err := tr.FormatQuantity("600519.SH", 350)
	require.NoError(t, err)
	assert.Equal(t, "300", qty)

	price, err := tr.GetMarketPrice("600519.SH")
	require.NoError(t, err)
	assert.Equal(t, 100.0, price)

	assert.NoError(t, tr.SetLeverage("600519.SH", 1))
	assert.Error(t, tr.SetLeverage("600519.SH", 2))
	_, err = tr.PlaceLimitOrder(&types.LimitOrderRequest{})
	assert.Error(t, err)
}

func TestNoPriceMeansNoFill(t *testing.T) {
	tr := newTestTrader(t, morning, fakeQuoter{})
	_, err := tr.OpenLong("600519.SH", 100, 1)
	assert.ErrorContains(t, err, "no quote")
}
