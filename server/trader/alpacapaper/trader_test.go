package alpacapaper

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeAlpaca is an in-memory stand-in for the Alpaca trading and data APIs.
type fakeAlpaca struct {
	t         *testing.T
	mu        sync.Mutex
	positions []map[string]any
	orders    []map[string]any
	price     float64
	requests  []recorded
	orderSeq  int
	forceErr  int
}

type recorded struct {
	Method string
	Path   string
	Body   map[string]any
}

func (f *fakeAlpaca) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		if r.Header.Get("APCA-API-KEY-ID") != "key" || r.Header.Get("APCA-API-SECRET-KEY") != "secret" {
			http.Error(w, `{"message":"forbidden"}`, http.StatusForbidden)
			return
		}
		var body map[string]any
		if raw, _ := io.ReadAll(r.Body); len(raw) > 0 {
			require.NoError(f.t, json.Unmarshal(raw, &body))
		}
		f.requests = append(f.requests, recorded{Method: r.Method, Path: r.URL.Path, Body: body})
		if f.forceErr > 0 {
			f.forceErr--
			http.Error(w, `{"message":"insufficient buying power"}`, http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/account":
			_, _ = io.WriteString(w, `{"cash":"50000","equity":"52000","buying_power":"100000","status":"ACTIVE"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v2/positions":
			_ = json.NewEncoder(w).Encode(f.positions)
		case r.Method == http.MethodPost && r.URL.Path == "/v2/orders":
			f.orderSeq++
			id := "order-" + string(rune('0'+f.orderSeq))
			status := "new"
			if body["type"] == "market" {
				status = "filled" // market orders fill immediately in the fake
				if body["side"] == "sell" {
					f.positions = removeLong(f.positions, body["symbol"])
				}
			}
			f.orders = append(f.orders, map[string]any{"id": id, "type": body["type"], "side": body["side"], "symbol": body["symbol"], "status": status, "qty": body["qty"]})
			_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "status": "accepted"})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/orders":
			open := []map[string]any{}
			for _, o := range f.orders {
				if o["status"] != "filled" {
					open = append(open, o)
				}
			}
			_ = json.NewEncoder(w).Encode(open)
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/v2/orders/"):
			id := strings.TrimPrefix(r.URL.Path, "/v2/orders/")
			kept := f.orders[:0]
			for _, o := range f.orders {
				if o["id"] != id {
					kept = append(kept, o)
				}
			}
			f.orders = kept
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v2/orders/"):
			id := strings.TrimPrefix(r.URL.Path, "/v2/orders/")
			_, _ = io.WriteString(w, `{"id":"`+id+`","status":"filled","filled_qty":"10","filled_avg_price":"101.5"}`)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v2/stocks/") && strings.HasSuffix(r.URL.Path, "/trades/latest"):
			_ = json.NewEncoder(w).Encode(map[string]any{"trade": map[string]any{"p": f.price}})
		default:
			http.NotFound(w, r)
		}
	})
}

// removeLong drops the long position for a symbol, as a full market sell would.
func removeLong(positions []map[string]any, symbol any) []map[string]any {
	out := positions[:0]
	for _, p := range positions {
		if p["symbol"] == symbol && p["side"] == "long" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func newFake(t *testing.T) (*Trader, *fakeAlpaca) {
	t.Helper()
	f := &fakeAlpaca{t: t, price: 150}
	srv := httptest.NewServer(f.handler())
	t.Cleanup(srv.Close)
	tr, err := New("key", "secret", srv.URL+"/v2", srv.URL+"/v2")
	require.NoError(t, err)
	return tr, f
}

func (f *fakeAlpaca) last() recorded {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.requests[len(f.requests)-1]
}

// lastPost returns the most recent POST, ignoring reads that happened after it.
func (f *fakeAlpaca) lastPost() recorded {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := len(f.requests) - 1; i >= 0; i-- {
		if f.requests[i].Method == http.MethodPost {
			return f.requests[i]
		}
	}
	return recorded{}
}

func TestNewRequiresCredentials(t *testing.T) {
	_, err := New("", "secret", "", "")
	assert.Error(t, err)
	tr, err := New("key", "secret", "", "")
	require.NoError(t, err)
	assert.Equal(t, DefaultTradingURL, tr.tradingURL)
	assert.Equal(t, DefaultDataURL, tr.dataURL)
}

func TestBadCredentialsSurfaceAsAPIError(t *testing.T) {
	tr, _ := newFake(t)
	tr.keyID = "wrong"
	_, err := tr.GetBalance()
	require.Error(t, err)
	var apiErr *apiError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status)
}

func TestGetBalance(t *testing.T) {
	tr, f := newFake(t)
	f.positions = []map[string]any{{"symbol": "AAPL", "qty": "10", "side": "long", "unrealized_pl": "120.5"}}

	bal, err := tr.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, 50000.0, bal["totalWalletBalance"])
	assert.Equal(t, 100000.0, bal["availableBalance"])
	assert.Equal(t, 52000.0, bal["totalEquity"])
	assert.Equal(t, 120.5, bal["totalUnrealizedProfit"])
}

func TestGetPositionsLongAndShort(t *testing.T) {
	tr, f := newFake(t)
	f.positions = []map[string]any{
		{"symbol": "AAPL", "qty": "10", "side": "long", "avg_entry_price": "100", "current_price": "110", "unrealized_pl": "100", "created_at": "2026-10-14T13:30:00Z"},
		{"symbol": "TSLA", "qty": "-5", "side": "short", "avg_entry_price": "200", "current_price": "190", "unrealized_pl": "50", "created_at": "2026-10-14T13:31:00Z"},
	}
	pos, err := tr.GetPositions()
	require.NoError(t, err)
	require.Len(t, pos, 2)
	assert.Equal(t, "long", pos[0]["side"])
	assert.Equal(t, 10.0, pos[0]["positionAmt"])
	assert.Equal(t, "short", pos[1]["side"])
	assert.Equal(t, 5.0, pos[1]["positionAmt"], "short sizes are reported as positive")
	assert.Equal(t, time.Date(2026, 10, 14, 13, 30, 0, 0, time.UTC).UnixMilli(), pos[0]["createdTime"])
}

func TestOpenLongSendsMarketBuy(t *testing.T) {
	tr, f := newFake(t)
	res, err := tr.OpenLong("AAPL", 2.5)
	require.NoError(t, err)
	assert.Equal(t, "order-1", res["orderId"])

	req := f.last()
	assert.Equal(t, http.MethodPost, req.Method)
	assert.Equal(t, "AAPL", req.Body["symbol"])
	assert.Equal(t, "buy", req.Body["side"])
	assert.Equal(t, "market", req.Body["type"])
	assert.Equal(t, "2.5", req.Body["qty"], "fractional shares are sent as-is")
	assert.Equal(t, "day", req.Body["time_in_force"])
}

func TestOpenShortSendsSell(t *testing.T) {
	tr, f := newFake(t)
	_, err := tr.OpenShort("TSLA", 3)
	require.NoError(t, err)
	assert.Equal(t, "sell", f.last().Body["side"])

}

func TestCloseLongAllAndRejectsMissing(t *testing.T) {
	tr, f := newFake(t)
	f.positions = []map[string]any{{"symbol": "AAPL", "qty": "10", "side": "long"}}

	_, err := tr.CloseLong("AAPL", 0)
	require.NoError(t, err)
	assert.Equal(t, "10", f.last().Body["qty"], "quantity 0 closes the full position")
	assert.Equal(t, "sell", f.last().Body["side"])

	_, err = tr.CloseLong("MSFT", 0)
	assert.ErrorContains(t, err, "no open long position")
}

func TestCloseShortBuysToCover(t *testing.T) {
	tr, f := newFake(t)
	f.positions = []map[string]any{{"symbol": "TSLA", "qty": "-4", "side": "short"}}
	_, err := tr.CloseShort("TSLA", 0)
	require.NoError(t, err)
	assert.Equal(t, "buy", f.last().Body["side"])
	assert.Equal(t, "4", f.last().Body["qty"])
}

func TestSetStopLossPlacesStopOrder(t *testing.T) {
	tr, f := newFake(t)
	require.NoError(t, tr.SetStopLoss("AAPL", "LONG", 10, 95))
	body := f.last().Body
	assert.Equal(t, "stop", body["type"])
	assert.Equal(t, "sell", body["side"])
	assert.Equal(t, "95.00", body["stop_price"])
	assert.Equal(t, "gtc", body["time_in_force"])

	require.NoError(t, tr.SetStopLoss("TSLA", "SHORT", 3, 210))
	assert.Equal(t, "buy", f.last().Body["side"], "a short is protected by a buy stop")

	assert.Error(t, tr.SetStopLoss("AAPL", "LONG", 10, 0))
}

func TestCancelStopLossOnlyCancelsStops(t *testing.T) {
	tr, f := newFake(t)
	f.orders = []map[string]any{
		{"id": "stop-1", "symbol": "AAPL", "type": "stop", "side": "sell", "qty": "10", "status": "new"},
		{"id": "limit-1", "symbol": "AAPL", "type": "limit", "side": "sell", "qty": "10", "status": "new"},
	}
	require.NoError(t, tr.CancelStopLossOrders("AAPL"))

	require.Len(t, f.orders, 1)
	assert.Equal(t, "limit-1", f.orders[0]["id"], "limit orders are left alone")
}

func TestTakeProfitClosesWhenPriceReached(t *testing.T) {
	tr, f := newFake(t)
	f.positions = []map[string]any{{"symbol": "AAPL", "qty": "10", "side": "long", "current_price": "149"}}
	require.NoError(t, tr.SetTakeProfit("AAPL", "LONG", 10, 150))

	pos, err := tr.GetPositions()
	require.NoError(t, err)
	assert.Len(t, pos, 1, "price below target: nothing closed yet")
	assert.Empty(t, f.lastPost().Method, "no POST yet means no order was sent")

	f.mu.Lock()
	f.positions[0]["current_price"] = "151"
	f.mu.Unlock()
	pos, err = tr.GetPositions()
	require.NoError(t, err)
	assert.Empty(t, pos, "position is closed once the target is reached")
	last := f.lastPost()
	assert.Equal(t, "sell", last.Body["side"])
	assert.Equal(t, "10", last.Body["qty"])

	orders, err := tr.GetOpenOrders("AAPL")
	require.NoError(t, err)
	assert.Empty(t, orders, "take-profit is cleared once it fires")
}

func TestGetOpenOrdersIncludesLocalTakeProfit(t *testing.T) {
	tr, f := newFake(t)
	f.orders = []map[string]any{{"id": "stop-1", "symbol": "AAPL", "type": "stop", "side": "sell", "qty": "10", "stop_price": "95", "status": "new"}}
	require.NoError(t, tr.SetTakeProfit("AAPL", "LONG", 10, 150))

	orders, err := tr.GetOpenOrders("AAPL")
	require.NoError(t, err)
	require.Len(t, orders, 2)
	assert.Equal(t, "STOP", strings.ToUpper(orders[0].Type)[:4])
	assert.Equal(t, "TAKE_PROFIT_MARKET", orders[1].Type)
	assert.Equal(t, 150.0, orders[1].StopPrice)
}

func TestGetMarketPriceUsesDataAPI(t *testing.T) {
	tr, _ := newFake(t)
	price, err := tr.GetMarketPrice("AAPL")
	require.NoError(t, err)
	assert.Equal(t, 150.0, price)
}

func TestGetOrderStatusMapsAlpacaStatus(t *testing.T) {
	tr, _ := newFake(t)
	st, err := tr.GetOrderStatus("AAPL", "order-1")
	require.NoError(t, err)
	assert.Equal(t, "FILLED", st["status"])
	assert.Equal(t, 101.5, st["avgPrice"])
	assert.Equal(t, 10.0, st["executedQty"])

	assert.Equal(t, "CANCELED", mapOrderStatus("expired"))
	assert.Equal(t, "NEW", mapOrderStatus("accepted"))
}

func TestFormatQuantity(t *testing.T) {
	assert.Equal(t, "2.5", formatQty(2.5))
	assert.Equal(t, "3", formatQty(3))
	assert.Equal(t, "0.000001", formatQty(0.0000012345))
}

func TestOrderRejectionSurfaces(t *testing.T) {
	tr, f := newFake(t)
	f.forceErr = 1
	_, err := tr.OpenLong("AAPL", 1)
	assert.ErrorContains(t, err, "insufficient buying power")
}
