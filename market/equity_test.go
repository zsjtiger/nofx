package market

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nofx/provider/tushare"
	"nofx/trader/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dailyKlines builds n ascending daily bars ending on the given date.
func dailyKlines(n int, last time.Time, startPrice float64) []Kline {
	out := make([]Kline, 0, n)
	for i := 0; i < n; i++ {
		day := last.AddDate(0, 0, i-(n-1))
		p := startPrice + float64(i)*0.5
		out = append(out, Kline{OpenTime: day.UnixMilli(), CloseTime: day.Add(24 * time.Hour).UnixMilli(), Open: p, High: p + 1, Low: p - 1, Close: p, Volume: 1000})
	}
	return out
}

func TestGetForMarketUsesRegisteredSource(t *testing.T) {
	RegisterEquitySource(types.MarketCNEquity, func(_ context.Context, symbol string, _ int) ([]Kline, error) {
		assert.Equal(t, "600519.SH", symbol, "symbol is canonicalised before the source is called")
		return dailyKlines(60, time.Date(2026, 10, 14, 15, 0, 0, 0, time.UTC), 100), nil
	})
	t.Cleanup(func() { RegisterEquitySource(types.MarketCNEquity, nil) })

	data, err := GetForMarket("600519", types.MarketCNEquity)
	require.NoError(t, err)
	assert.Equal(t, "600519.SH", data.Symbol)
	assert.InDelta(t, 100+59*0.5, data.CurrentPrice, 0.001)
	assert.NotNil(t, data.LongerTermContext)
}

func TestGetForMarketErrors(t *testing.T) {
	_, err := GetForMarket("600519", types.MarketCNEquity)
	require.Error(t, err, "no source registered")

	RegisterEquitySource(types.MarketUSEquity, func(context.Context, string, int) ([]Kline, error) {
		return dailyKlines(5, time.Now(), 10), nil
	})
	t.Cleanup(func() { RegisterEquitySource(types.MarketUSEquity, nil) })
	_, err = GetForMarket("AAPL", types.MarketUSEquity)
	assert.ErrorContains(t, err, "need 30")

	_, err = GetForMarket("6005", types.MarketCNEquity)
	assert.Error(t, err, "invalid A-share code is rejected before any source call")
}

func TestCNBarsToKlinesUsesBeijingCloseTime(t *testing.T) {
	klines, err := cnBarsToKlines([]tushare.Bar{{TradeDate: "20261014", Open: 10, High: 11, Low: 9, Close: 10.5, Vol: 25}})
	require.NoError(t, err)
	require.Len(t, klines, 1)
	// 15:00 Beijing on 14 Oct is 07:00 UTC.
	assert.Equal(t, time.Date(2026, 10, 14, 7, 0, 0, 0, time.UTC).UnixMilli(), klines[0].CloseTime)
	assert.Equal(t, 2500.0, klines[0].Volume, "lots are converted to shares")

	_, err = cnBarsToKlines([]tushare.Bar{{TradeDate: "bad"}})
	assert.Error(t, err)
}

func TestCNDailyKlineSourceAdjustsPrices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		require.NoError(t, json.Unmarshal(body, &req))
		var resp string
		switch req["api_name"] {
		case "daily":
			resp = `{"code":0,"data":{"fields":["trade_date","open","high","low","close","pre_close","vol","amount"],
				"items":[["20261013",10,11,9,10,10,100,1],["20261014",20,21,19,20,10,100,1]]}}`
		case "adj_factor":
			resp = `{"code":0,"data":{"fields":["trade_date","adj_factor"],"items":[["20261013",1],["20261014",2]]}}`
		default:
			http.Error(w, "unexpected", http.StatusTeapot)
			return
		}
		_, _ = io.WriteString(w, resp)
	}))
	t.Cleanup(srv.Close)

	src := CNDailyKlineSource(tushare.NewClientWithBaseURL("tok", srv.URL))
	klines, err := src(context.Background(), "600519.SH", 10)
	require.NoError(t, err)
	require.Len(t, klines, 2)
	assert.Equal(t, 5.0, klines[0].Close, "earlier bar is forward-adjusted by 1/2")
	assert.Equal(t, 20.0, klines[1].Close, "latest bar is unchanged")
}

func TestNormalizeFor(t *testing.T) {
	assert.Equal(t, "600519.SH", NormalizeFor(types.MarketCNEquity, "600519"))
	assert.Equal(t, "AAPL", NormalizeFor(types.MarketUSEquity, "aapl"))
	assert.Equal(t, "NOT-A-CODE", NormalizeFor(types.MarketCNEquity, "not-a-code"), "invalid input is upper-cased and trimmed")
}
