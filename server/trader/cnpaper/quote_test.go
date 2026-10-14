package cnpaper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"crux/provider/tushare"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTushareQuoterUsesLatestDailyBarAndLimits(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		require.NoError(t, json.Unmarshal(body, &req))
		var resp string
		switch req["api_name"] {
		case "daily":
			resp = `{"code":0,"data":{"fields":["trade_date","open","high","low","close","pre_close","vol","amount"],
				"items":[["20261014",101,103,99,102,100,8000,1],["20261013",99,100,98,100,99,7000,1]]}}`
		case "stk_limit":
			resp = `{"code":0,"data":{"fields":["trade_date","up_limit","down_limit"],"items":[["20261014",112.2,91.8]]}}`
		default:
			http.Error(w, "unexpected", http.StatusTeapot)
			return
		}
		_, _ = io.WriteString(w, resp)
	}))
	t.Cleanup(srv.Close)

	client := tushare.NewClientWithBaseURL("tok", srv.URL)
	q := NewTushareQuoter(client)
	q.now = func() time.Time { return time.Date(2026, 10, 14, 15, 30, 0, 0, cst) }

	quote, err := q.Quote(context.Background(), "600519.SH")
	require.NoError(t, err)
	assert.Equal(t, 102.0, quote.Price, "latest daily close")
	assert.Equal(t, 100.0, quote.PrevClose)
	assert.Equal(t, 112.2, quote.UpLimit)
	assert.Equal(t, 91.8, quote.DownLimit)
	assert.Equal(t, 8000.0, quote.Volume)
}

func TestTushareQuoterNoBarsIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"code":0,"data":{"fields":["trade_date"],"items":[]}}`)
	}))
	t.Cleanup(srv.Close)

	q := NewTushareQuoter(tushare.NewClientWithBaseURL("tok", srv.URL))
	_, err := q.Quote(context.Background(), "600519.SH")
	assert.ErrorContains(t, err, "no daily bars")
}
