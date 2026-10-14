package tushare

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeTushare serves canned responses keyed by api_name and records the last request.
func fakeTushare(t *testing.T, responses map[string]string) (*Client, *map[string]any) {
	t.Helper()
	var last map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &last))
		name, _ := last["api_name"].(string)
		resp, ok := responses[name]
		if !ok {
			http.Error(w, "unexpected api "+name, http.StatusTeapot)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, resp)
	}))
	t.Cleanup(srv.Close)
	c := NewClient("test-token")
	c.baseURL = srv.URL
	return c, &last
}

func TestCallSendsTokenAndParsesTable(t *testing.T) {
	c, last := fakeTushare(t, map[string]string{
		"daily": `{"code":0,"msg":"","data":{"fields":["trade_date","close"],"items":[["20261014",10.5]]}}`,
	})
	table, err := c.Call(context.Background(), "daily", map[string]string{"ts_code": "600519.SH"}, "trade_date,close")
	require.NoError(t, err)
	assert.Equal(t, "test-token", (*last)["token"])
	assert.Equal(t, []string{"trade_date", "close"}, table.Fields)
	assert.Equal(t, "20261014", table.Rows()[0]["trade_date"])
	assert.Equal(t, 10.5, table.Rows()[0]["close"])
}

func TestCallErrors(t *testing.T) {
	c, _ := fakeTushare(t, map[string]string{
		"bad_code": `{"code":40203,"msg":"no permission","data":null}`,
		"bad_json": `not json`,
	})
	_, err := c.Call(context.Background(), "bad_code", nil, "")
	assert.ErrorContains(t, err, "no permission")

	_, err = c.Call(context.Background(), "bad_json", nil, "")
	assert.ErrorContains(t, err, "invalid response")

	_, err = c.Call(context.Background(), "unknown", nil, "")
	assert.ErrorContains(t, err, "HTTP 418")

	empty := NewClient("")
	_, err = empty.Call(context.Background(), "daily", nil, "")
	assert.ErrorContains(t, err, "token is not configured")
}

func TestDailyBarsSortedOldestFirst(t *testing.T) {
	c, last := fakeTushare(t, map[string]string{
		"daily": `{"code":0,"data":{"fields":["trade_date","open","high","low","close","pre_close","vol","amount"],
			"items":[["20261014",11,12,10.5,11.5,10.8,1000,5000],["20261013",10,11,9.8,10.8,10,800,4000]]}}`,
	})
	bars, err := c.DailyBars(context.Background(), "600519.SH", "20261001", "20261014")
	require.NoError(t, err)
	require.Len(t, bars, 2)
	assert.Equal(t, "20261013", bars[0].TradeDate)
	assert.Equal(t, 11.5, bars[1].Close)
	assert.Equal(t, "600519.SH", (*last)["params"].(map[string]any)["ts_code"])
}

func TestAdjFactorsAndForwardAdjust(t *testing.T) {
	c, _ := fakeTushare(t, map[string]string{
		"adj_factor": `{"code":0,"data":{"fields":["trade_date","adj_factor"],"items":[["20261014",2],["20261013",1]]}}`,
	})
	factors, err := c.AdjFactors(context.Background(), "600519.SH", "20261001", "20261014")
	require.NoError(t, err)
	assert.Equal(t, 1.0, factors["20261013"])

	bars := []Bar{
		{TradeDate: "20261013", Open: 10, High: 11, Low: 9.5, Close: 10, PreClose: 9.9},
		{TradeDate: "20261014", Open: 20, High: 22, Low: 19, Close: 21, PreClose: 20},
	}
	adj, err := ForwardAdjust(bars, factors)
	require.NoError(t, err)
	// Earlier bar scaled by 1/2; latest bar is unchanged.
	assert.Equal(t, 5.0, adj[0].Open)
	assert.Equal(t, 5.0, adj[0].Close)
	assert.Equal(t, 21.0, adj[1].Close)
}

func TestForwardAdjustMissingFactors(t *testing.T) {
	bars := []Bar{{TradeDate: "20261013", Close: 10}, {TradeDate: "20261014", Close: 11}}

	_, err := ForwardAdjust(bars, map[string]float64{"20261013": 1})
	assert.ErrorContains(t, err, "missing adjustment factor for 20261014")

	adj, err := ForwardAdjust(bars, map[string]float64{"20261014": 1})
	assert.ErrorContains(t, err, "20261013")
	assert.Equal(t, 10.0, adj[0].Close, "bars without a factor are returned unadjusted")

	none, err := ForwardAdjust(nil, nil)
	assert.NoError(t, err)
	assert.Empty(t, none)
}

func TestStockLimit(t *testing.T) {
	c, _ := fakeTushare(t, map[string]string{
		"stk_limit": `{"code":0,"data":{"fields":["trade_date","up_limit","down_limit"],"items":[["20261014",110,90]]}}`,
	})
	up, down, err := c.StockLimit(context.Background(), "600519.SH", "20261014")
	require.NoError(t, err)
	assert.Equal(t, 110.0, up)
	assert.Equal(t, 90.0, down)

	c2, _ := fakeTushare(t, map[string]string{
		"stk_limit": `{"code":0,"data":{"fields":["trade_date","up_limit","down_limit"],"items":[]}}`,
	})
	_, _, err = c2.StockLimit(context.Background(), "600519.SH", "20261014")
	assert.ErrorContains(t, err, "no limit prices")
}

func TestIsST(t *testing.T) {
	c, _ := fakeTushare(t, map[string]string{
		"stock_basic": `{"code":0,"data":{"fields":["ts_code","name"],"items":[["600000.SH","*ST浦发"]]}}`,
	})
	st, err := c.IsST(context.Background(), "600000.SH")
	require.NoError(t, err)
	assert.True(t, st)
}

func TestIsSTName(t *testing.T) {
	assert.True(t, IsSTName("ST华仪"))
	assert.True(t, IsSTName("*ST浦发"))
	assert.False(t, IsSTName("贵州茅台"))
	assert.False(t, IsSTName("ABC"))
}
