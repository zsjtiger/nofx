// Package alpacapaper trades US equities through the Alpaca paper-trading API.
// Orders go to the paper account only; the base URL is configurable so the
// adapter can be pointed at a test server.
package alpacapaper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"crux/trader/types"
)

const (
	// DefaultTradingURL is the Alpaca paper-trading endpoint.
	DefaultTradingURL = "https://paper-api.alpaca.markets/v2"
	// DefaultDataURL is the Alpaca market-data endpoint.
	DefaultDataURL = "https://data.alpaca.markets/v2"
)

// Trader implements types.Trader on an Alpaca paper account.
type Trader struct {
	mu         sync.Mutex
	keyID      string
	secret     string
	tradingURL string
	dataURL    string
	http       *http.Client
	caps       types.Capabilities
	// Take-profit levels are kept here because Alpaca cannot hold a stop and a
	// limit exit on the same shares without an OCO order. They are checked on each read.
	takeProfit map[string]float64
}

// New returns a paper trader. Empty URLs use the public Alpaca endpoints.
func New(keyID, secret, tradingURL, dataURL string) (*Trader, error) {
	if keyID == "" || secret == "" {
		return nil, errors.New("alpaca API key and secret are required")
	}
	if tradingURL == "" {
		tradingURL = DefaultTradingURL
	}
	if dataURL == "" {
		dataURL = DefaultDataURL
	}
	caps, err := types.CapabilitiesOf(types.MarketUSEquity)
	if err != nil {
		return nil, err
	}
	return &Trader{
		keyID:      keyID,
		secret:     secret,
		tradingURL: strings.TrimRight(tradingURL, "/"),
		dataURL:    strings.TrimRight(dataURL, "/"),
		http:       &http.Client{Timeout: 15 * time.Second},
		caps:       caps,
		takeProfit: map[string]float64{},
	}, nil
}

type apiError struct {
	Status int
	Body   string
}

func (e *apiError) Error() string { return fmt.Sprintf("alpaca HTTP %d: %s", e.Status, e.Body) }

func (t *Trader) do(method, base, path string, query url.Values, body any, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	u := base + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(context.Background(), method, u, reader)
	if err != nil {
		return err
	}
	req.Header.Set("APCA-API-KEY-ID", t.keyID)
	req.Header.Set("APCA-API-SECRET-KEY", t.secret)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := t.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return &apiError{Status: resp.StatusCode, Body: strings.TrimSpace(string(raw))}
	}
	if out != nil && len(raw) > 0 {
		return json.Unmarshal(raw, out)
	}
	return nil
}

type account struct {
	Cash        string `json:"cash"`
	Equity      string `json:"equity"`
	BuyingPower string `json:"buying_power"`
	Status      string `json:"status"`
}

type apiPosition struct {
	Symbol        string `json:"symbol"`
	Qty           string `json:"qty"`
	Side          string `json:"side"`
	AvgEntryPrice string `json:"avg_entry_price"`
	CurrentPrice  string `json:"current_price"`
	UnrealizedPL  string `json:"unrealized_pl"`
	CreatedAt     string `json:"created_at"`
}

type apiOrder struct {
	ID             string `json:"id"`
	Symbol         string `json:"symbol"`
	Side           string `json:"side"`
	Type           string `json:"type"`
	Qty            string `json:"qty"`
	FilledQty      string `json:"filled_qty"`
	FilledAvgPrice string `json:"filled_avg_price"`
	StopPrice      string `json:"stop_price"`
	LimitPrice     string `json:"limit_price"`
	Status         string `json:"status"`
}

func parseF(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// GetBalance returns cash, equity and unrealised P&L from the paper account.
func (t *Trader) GetBalance() (map[string]interface{}, error) {
	var acc account
	if err := t.do(http.MethodGet, t.tradingURL, "/account", nil, nil, &acc); err != nil {
		return nil, err
	}
	positions, err := t.fetchPositions()
	if err != nil {
		return nil, err
	}
	var unrealized float64
	for _, p := range positions {
		unrealized += parseF(p.UnrealizedPL)
	}
	return map[string]interface{}{
		"totalWalletBalance":    parseF(acc.Cash),
		"availableBalance":      parseF(acc.BuyingPower),
		"totalUnrealizedProfit": unrealized,
		"totalEquity":           parseF(acc.Equity),
	}, nil
}

func (t *Trader) fetchPositions() ([]apiPosition, error) {
	var out []apiPosition
	err := t.do(http.MethodGet, t.tradingURL, "/positions", nil, nil, &out)
	return out, err
}

// GetPositions returns open positions, long or short, with positive sizes.
// It first closes any long position whose take-profit level has been reached.
func (t *Trader) GetPositions() ([]map[string]interface{}, error) {
	if err := t.enforceTakeProfit(); err != nil {
		return nil, err
	}
	raw, err := t.fetchPositions()
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, p := range raw {
		qty := math.Abs(parseF(p.Qty))
		side := "long"
		if strings.EqualFold(p.Side, "short") {
			side = "short"
		}
		createdMs := int64(0)
		if ts, err := time.Parse(time.RFC3339Nano, p.CreatedAt); err == nil {
			createdMs = ts.UnixMilli()
		}
		out = append(out, map[string]interface{}{
			"symbol":           p.Symbol,
			"side":             side,
			"positionAmt":      qty,
			"entryPrice":       parseF(p.AvgEntryPrice),
			"markPrice":        parseF(p.CurrentPrice),
			"unRealizedProfit": parseF(p.UnrealizedPL),
			"createdTime":      createdMs,
		})
	}
	return out, nil
}

// enforceTakeProfit sells long positions whose current price has reached the take-profit level.
func (t *Trader) enforceTakeProfit() error {
	t.mu.Lock()
	levels := make(map[string]float64, len(t.takeProfit))
	for sym, tp := range t.takeProfit {
		levels[sym] = tp
	}
	t.mu.Unlock()
	if len(levels) == 0 {
		return nil
	}
	raw, err := t.fetchPositions()
	if err != nil {
		return err
	}
	for _, p := range raw {
		tp, ok := levels[p.Symbol]
		if !ok || strings.EqualFold(p.Side, "short") {
			continue
		}
		if parseF(p.CurrentPrice) < tp {
			continue
		}
		if _, err := t.submitMarket(p.Symbol, "sell", math.Abs(parseF(p.Qty))); err != nil {
			return err
		}
		if err := t.CancelTakeProfitOrders(p.Symbol); err != nil {
			return err
		}
	}
	return nil
}

// OpenLong buys with a market order. Fractional share quantities are allowed.
func (t *Trader) OpenLong(symbol string, quantity float64) (map[string]interface{}, error) {
	if err := types.ValidateSideAllowed(t.caps, false); err != nil {
		return nil, err
	}
	return t.submitMarket(symbol, "buy", quantity)
}

// OpenShort sells short with a market order. The margin account must allow shorting.
func (t *Trader) OpenShort(symbol string, quantity float64) (map[string]interface{}, error) {
	if err := types.ValidateSideAllowed(t.caps, true); err != nil {
		return nil, err
	}
	return t.submitMarket(symbol, "sell", quantity)
}

// CloseLong sells the long position. quantity 0 closes all shares.
func (t *Trader) CloseLong(symbol string, quantity float64) (map[string]interface{}, error) {
	qty, err := t.resolveCloseQty(symbol, "long", quantity)
	if err != nil {
		return nil, err
	}
	return t.submitMarket(symbol, "sell", qty)
}

// CloseShort buys to cover the short position. quantity 0 closes all.
func (t *Trader) CloseShort(symbol string, quantity float64) (map[string]interface{}, error) {
	qty, err := t.resolveCloseQty(symbol, "short", quantity)
	if err != nil {
		return nil, err
	}
	return t.submitMarket(symbol, "buy", qty)
}

func (t *Trader) resolveCloseQty(symbol, side string, quantity float64) (float64, error) {
	positions, err := t.GetPositions()
	if err != nil {
		return 0, err
	}
	for _, p := range positions {
		if p["symbol"] != symbol || p["side"] != side {
			continue
		}
		held := p["positionAmt"].(float64)
		if quantity == 0 || quantity > held {
			return held, nil
		}
		return quantity, nil
	}
	return 0, fmt.Errorf("no open %s position for %s", side, symbol)
}

func (t *Trader) submitMarket(symbol, side string, qty float64) (map[string]interface{}, error) {
	if err := types.ValidateOrderQuantity(t.caps, qty, side == "buy", false); err != nil {
		return nil, err
	}
	body := map[string]any{
		"symbol":        symbol,
		"qty":           formatQty(qty),
		"side":          side,
		"type":          "market",
		"time_in_force": "day",
	}
	var order apiOrder
	if err := t.do(http.MethodPost, t.tradingURL, "/orders", nil, body, &order); err != nil {
		return nil, err
	}
	return map[string]interface{}{"orderId": order.ID, "status": order.Status}, nil
}

// GetMarketPrice returns the latest trade price from the data API.
func (t *Trader) GetMarketPrice(symbol string) (float64, error) {
	var out struct {
		Trade struct {
			P float64 `json:"p"`
		} `json:"trade"`
	}
	if err := t.do(http.MethodGet, t.dataURL, "/stocks/"+url.PathEscape(symbol)+"/trades/latest", nil, nil, &out); err != nil {
		return 0, err
	}
	if out.Trade.P <= 0 {
		return 0, fmt.Errorf("no trade price for %s", symbol)
	}
	return out.Trade.P, nil
}

// SetStopLoss places a good-till-cancel stop order that sells the position.
func (t *Trader) SetStopLoss(symbol string, positionSide string, quantity, stopPrice float64) error {
	if stopPrice <= 0 {
		return fmt.Errorf("stop price must be positive")
	}
	side := "sell"
	if strings.EqualFold(positionSide, "SHORT") {
		side = "buy"
	}
	body := map[string]any{
		"symbol":        symbol,
		"qty":           formatQty(quantity),
		"side":          side,
		"type":          "stop",
		"stop_price":    strconv.FormatFloat(stopPrice, 'f', 2, 64),
		"time_in_force": "gtc",
	}
	return t.do(http.MethodPost, t.tradingURL, "/orders", nil, body, nil)
}

// SetTakeProfit records a take-profit level, enforced on each read by closing the position.
func (t *Trader) SetTakeProfit(symbol string, positionSide string, quantity, takeProfitPrice float64) error {
	if takeProfitPrice <= 0 {
		return fmt.Errorf("take-profit price must be positive")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.takeProfit[symbol] = takeProfitPrice
	return nil
}

// CancelStopLossOrders cancels open stop orders for the symbol.
func (t *Trader) CancelStopLossOrders(symbol string) error {
	return t.cancelOpenOrders(symbol, func(o apiOrder) bool { return o.Type == "stop" || o.Type == "stop_limit" })
}

// CancelTakeProfitOrders clears the take-profit level for the symbol.
func (t *Trader) CancelTakeProfitOrders(symbol string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.takeProfit, symbol)
	return nil
}

// CancelAllOrders cancels every open order for the symbol and clears the take-profit.
func (t *Trader) CancelAllOrders(symbol string) error {
	if err := t.cancelOpenOrders(symbol, func(apiOrder) bool { return true }); err != nil {
		return err
	}
	return t.CancelTakeProfitOrders(symbol)
}

// CancelStopOrders cancels stop orders and clears the take-profit for the symbol.
func (t *Trader) CancelStopOrders(symbol string) error { return t.CancelAllOrders(symbol) }

func (t *Trader) cancelOpenOrders(symbol string, match func(apiOrder) bool) error {
	orders, err := t.fetchOpenOrders(symbol)
	if err != nil {
		return err
	}
	for _, o := range orders {
		if !match(o) {
			continue
		}
		if err := t.do(http.MethodDelete, t.tradingURL, "/orders/"+o.ID, nil, nil, nil); err != nil {
			return err
		}
	}
	return nil
}

func (t *Trader) fetchOpenOrders(symbol string) ([]apiOrder, error) {
	q := url.Values{"status": {"open"}, "limit": {"100"}}
	if symbol != "" {
		q.Set("symbols", symbol)
	}
	var out []apiOrder
	err := t.do(http.MethodGet, t.tradingURL, "/orders", q, nil, &out)
	return out, err
}

// FormatQuantity returns the quantity as text with up to six decimals, trailing zeros removed.
func (t *Trader) FormatQuantity(symbol string, quantity float64) (string, error) {
	return formatQty(quantity), nil
}

func formatQty(q float64) string {
	s := strconv.FormatFloat(q, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

// GetOrderStatus maps an Alpaca order to the status fields the AutoTrader reads.
func (t *Trader) GetOrderStatus(symbol string, orderID string) (map[string]interface{}, error) {
	var o apiOrder
	if err := t.do(http.MethodGet, t.tradingURL, "/orders/"+url.PathEscape(orderID), nil, nil, &o); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"status":      mapOrderStatus(o.Status),
		"avgPrice":    parseF(o.FilledAvgPrice),
		"executedQty": parseF(o.FilledQty),
		"commission":  0.0,
	}, nil
}

func mapOrderStatus(s string) string {
	switch strings.ToLower(s) {
	case "filled":
		return "FILLED"
	case "canceled", "expired", "rejected", "done_for_day":
		return "CANCELED"
	default:
		return "NEW"
	}
}

// GetClosedPnL is not reconstructed from Alpaca activities yet. The AutoTrader
// records closed trades from its own position data, so this returns no rows.
func (t *Trader) GetClosedPnL(startTime time.Time, limit int) ([]types.ClosedPnLRecord, error) {
	return nil, nil
}

// GetOpenOrders returns open stop orders and the local take-profit level.
func (t *Trader) GetOpenOrders(symbol string) ([]types.OpenOrder, error) {
	orders, err := t.fetchOpenOrders(symbol)
	if err != nil {
		return nil, err
	}
	var out []types.OpenOrder
	for _, o := range orders {
		out = append(out, types.OpenOrder{
			OrderID:   o.ID,
			Symbol:    o.Symbol,
			Side:      strings.ToUpper(o.Side),
			Type:      strings.ToUpper(o.Type),
			StopPrice: parseF(o.StopPrice),
			Price:     parseF(o.LimitPrice),
			Quantity:  parseF(o.Qty),
			Status:    strings.ToUpper(o.Status),
		})
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for sym, tp := range t.takeProfit {
		if symbol != "" && sym != symbol {
			continue
		}
		out = append(out, types.OpenOrder{Symbol: sym, Type: "TAKE_PROFIT_MARKET", StopPrice: tp, Status: "NEW"})
	}
	return out, nil
}

var _ types.Trader = (*Trader)(nil)
