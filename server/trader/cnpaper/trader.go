// Package cnpaper is a paper-trading broker for A-shares. It keeps cash and
// positions in memory and fills orders at the latest quote, applying the
// A-share rules from trader/types: whole lots, T+1, price limits, sessions
// and no shorting. It is for simulation only and never sends orders to a broker.
package cnpaper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"nofx/logger"
	"nofx/trader/types"
)

// Quote is the latest market state for one A-share.
type Quote struct {
	Price     float64 // last price
	PrevClose float64
	UpLimit   float64 // limit-up price for the day
	DownLimit float64 // limit-down price for the day
	Volume    float64 // zero means suspended or no trading
}

// Quoter supplies quotes. The Tushare-backed implementation lives in quote.go.
type Quoter interface {
	Quote(ctx context.Context, symbol string) (Quote, error)
}

// Fee schedule for A-share trading (commission on both sides, stamp duty on sells).
const (
	commissionRate = 0.00025
	minCommission  = 5.0
	stampDutyRate  = 0.0005
)

type position struct {
	symbol     string
	quantity   float64
	avgCost    float64 // includes buy commissions
	entryTime  time.Time
	stopLoss   float64
	takeProfit float64
}

type fill struct {
	orderID  string
	symbol   string
	side     string
	quantity float64
	price    float64
	fee      float64
	time     time.Time
}

var _ types.Trader = (*Trader)(nil)

// Trader implements types.Trader for paper A-share trading.
type Trader struct {
	mu        sync.Mutex
	caps      types.Capabilities
	quoter    Quoter
	now       func() time.Time
	cash      float64
	positions map[string]*position
	fills     map[string]fill
	closed    []types.ClosedPnLRecord
	nextID    int
	statePath string // when set, the ledger is saved here after every change
}

// state is the persisted form of the ledger.
type state struct {
	Cash      float64                 `json:"cash"`
	NextID    int                     `json:"next_id"`
	Positions []positionState         `json:"positions"`
	Fills     map[string]fill         `json:"fills"`
	Closed    []types.ClosedPnLRecord `json:"closed"`
}

type positionState struct {
	Symbol     string    `json:"symbol"`
	Quantity   float64   `json:"quantity"`
	AvgCost    float64   `json:"avg_cost"`
	EntryTime  time.Time `json:"entry_time"`
	StopLoss   float64   `json:"stop_loss"`
	TakeProfit float64   `json:"take_profit"`
}

// New returns a paper trader with the given starting cash in CNY. The ledger
// lives in memory only; use NewPersistent to keep it across restarts.
func New(initialCash float64, quoter Quoter) (*Trader, error) {
	return NewPersistent(initialCash, quoter, "")
}

// NewPersistent returns a paper trader whose ledger is saved to statePath after
// every change. If the file exists, the saved ledger is loaded and initialCash
// is ignored. An empty statePath keeps the ledger in memory only.
func NewPersistent(initialCash float64, quoter Quoter, statePath string) (*Trader, error) {
	if initialCash <= 0 {
		return nil, fmt.Errorf("initial cash must be positive")
	}
	if quoter == nil {
		return nil, fmt.Errorf("quoter is required")
	}
	caps, err := types.CapabilitiesOf(types.MarketCNEquity)
	if err != nil {
		return nil, err
	}
	t := &Trader{
		caps:      caps,
		quoter:    quoter,
		now:       time.Now,
		cash:      initialCash,
		positions: map[string]*position{},
		fills:     map[string]fill{},
		statePath: statePath,
	}
	if statePath != "" {
		if err := t.load(); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// load restores the ledger from statePath when the file exists.
func (t *Trader) load() error {
	raw, err := os.ReadFile(t.statePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read paper ledger: %w", err)
	}
	var st state
	if err := json.Unmarshal(raw, &st); err != nil {
		return fmt.Errorf("parse paper ledger %s: %w", t.statePath, err)
	}
	t.cash = st.Cash
	t.nextID = st.NextID
	t.closed = st.Closed
	if st.Fills != nil {
		t.fills = st.Fills
	}
	t.positions = make(map[string]*position, len(st.Positions))
	for _, p := range st.Positions {
		t.positions[p.Symbol] = &position{
			symbol:     p.Symbol,
			quantity:   p.Quantity,
			avgCost:    p.AvgCost,
			entryTime:  p.EntryTime,
			stopLoss:   p.StopLoss,
			takeProfit: p.TakeProfit,
		}
	}
	return nil
}

// save writes the ledger atomically: a temp file is written, then renamed over the old one.
func (t *Trader) save() error {
	if t.statePath == "" {
		return nil
	}
	st := state{Cash: t.cash, NextID: t.nextID, Fills: t.fills, Closed: t.closed}
	for _, p := range t.positions {
		st.Positions = append(st.Positions, positionState{
			Symbol:     p.symbol,
			Quantity:   p.quantity,
			AvgCost:    p.avgCost,
			EntryTime:  p.entryTime,
			StopLoss:   p.stopLoss,
			TakeProfit: p.takeProfit,
		})
	}
	raw, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(t.statePath), 0o755); err != nil {
		return err
	}
	tmp := t.statePath + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, t.statePath)
}

func (t *Trader) quote(symbol string) (Quote, error) {
	q, err := t.quoter.Quote(context.Background(), symbol)
	if err != nil {
		return Quote{}, err
	}
	if q.Price <= 0 {
		return Quote{}, fmt.Errorf("no price for %s", symbol)
	}
	return q, nil
}

func commission(amount float64) float64 {
	return math.Max(minCommission, amount*commissionRate)
}

// GetBalance returns cash, market value and unrealised P&L. It first applies any triggered stops.
func (t *Trader) GetBalance() (map[string]interface{}, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.checkProtectionsLocked()

	var marketValue, unrealized float64
	for _, p := range t.positions {
		q, err := t.quote(p.symbol)
		if err != nil {
			return nil, err
		}
		marketValue += p.quantity * q.Price
		unrealized += p.quantity * (q.Price - p.avgCost)
	}
	return map[string]interface{}{
		"totalWalletBalance":    t.cash,
		"availableBalance":      t.cash,
		"totalUnrealizedProfit": unrealized,
		"totalEquity":           t.cash + marketValue,
	}, nil
}

// GetPositions returns open long positions in the shape the AutoTrader expects.
func (t *Trader) GetPositions() ([]map[string]interface{}, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.checkProtectionsLocked()

	out := make([]map[string]interface{}, 0, len(t.positions))
	for _, p := range t.positions {
		q, err := t.quote(p.symbol)
		if err != nil {
			return nil, err
		}
		out = append(out, map[string]interface{}{
			"symbol":           p.symbol,
			"side":             "long",
			"positionAmt":      p.quantity,
			"entryPrice":       p.avgCost,
			"markPrice":        q.Price,
			"unRealizedProfit": p.quantity * (q.Price - p.avgCost),
			"leverage":         1,
			"liquidationPrice": 0.0,
			"createdTime":      p.entryTime.UnixMilli(),
		})
	}
	return out, nil
}

// OpenLong buys whole lots at the latest price, subject to limit-up and suspension.
func (t *Trader) OpenLong(symbol string, quantity float64, leverage int) (map[string]interface{}, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := types.ValidateSideAllowed(t.caps, leverage, false); err != nil {
		return nil, err
	}
	if err := types.ValidateOrderQuantity(t.caps, quantity, true, false); err != nil {
		return nil, err
	}
	q, err := t.quote(symbol)
	if err != nil {
		return nil, err
	}
	if err := checkFillable(q, true); err != nil {
		return nil, err
	}
	amount := quantity * q.Price
	fee := commission(amount)
	if amount+fee > t.cash {
		return nil, fmt.Errorf("insufficient cash: need %.2f, have %.2f", amount+fee, t.cash)
	}

	id := t.newOrderID()
	t.cash -= amount + fee
	p, ok := t.positions[symbol]
	if !ok {
		p = &position{symbol: symbol, entryTime: t.now()}
		t.positions[symbol] = p
	}
	totalCost := p.avgCost*p.quantity + amount + fee
	p.quantity += quantity
	p.avgCost = totalCost / p.quantity
	t.recordFill(id, symbol, "BUY", quantity, q.Price, fee)
	t.persist()
	return orderResult(id, q.Price, quantity), nil
}

// OpenShort is always rejected: A-share cash accounts cannot sell short.
func (t *Trader) OpenShort(symbol string, quantity float64, leverage int) (map[string]interface{}, error) {
	return nil, types.ValidateSideAllowed(t.caps, leverage, true)
}

// CloseLong sells shares. quantity 0 closes the whole position.
func (t *Trader) CloseLong(symbol string, quantity float64) (map[string]interface{}, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.closeLocked(symbol, quantity, "manual")
}

// CloseShort is always rejected: there are no short positions to close.
func (t *Trader) CloseShort(symbol string, quantity float64) (map[string]interface{}, error) {
	return nil, types.ValidateSideAllowed(t.caps, 1, true)
}

func (t *Trader) closeLocked(symbol string, quantity float64, reason string) (map[string]interface{}, error) {
	p, ok := t.positions[symbol]
	if !ok {
		return nil, fmt.Errorf("no open position for %s", symbol)
	}
	if quantity == 0 {
		quantity = p.quantity
	}
	if quantity > p.quantity {
		return nil, fmt.Errorf("cannot sell %v shares of %s: holding %v", quantity, symbol, p.quantity)
	}
	closesPosition := quantity == p.quantity
	if err := types.ValidateOrderQuantity(t.caps, quantity, false, closesPosition); err != nil {
		return nil, err
	}
	if err := types.ValidateSettlement(t.caps, p.entryTime, t.now()); err != nil {
		return nil, err
	}
	q, err := t.quote(symbol)
	if err != nil {
		return nil, err
	}
	if err := checkFillable(q, false); err != nil {
		return nil, err
	}

	amount := quantity * q.Price
	fee := commission(amount) + amount*stampDutyRate
	proceeds := amount - fee
	id := t.newOrderID()
	t.cash += proceeds
	t.recordFill(id, symbol, "SELL", quantity, q.Price, fee)

	t.closed = append(t.closed, types.ClosedPnLRecord{
		Symbol:      symbol,
		Side:        "long",
		EntryPrice:  p.avgCost,
		ExitPrice:   q.Price,
		Quantity:    quantity,
		RealizedPnL: proceeds - quantity*p.avgCost,
		Fee:         fee,
		Leverage:    1,
		EntryTime:   p.entryTime,
		ExitTime:    t.now(),
		OrderID:     id,
		CloseType:   reason,
	})

	p.quantity -= quantity
	if p.quantity <= 0 {
		delete(t.positions, symbol)
	}
	t.persist()
	return orderResult(id, q.Price, quantity), nil
}

// persist saves the ledger after a change. A failed write is logged, not returned,
// because the trade has already happened in memory and the caller cannot undo it.
func (t *Trader) persist() {
	if err := t.save(); err != nil {
		logger.Errorf("[cnpaper] failed to save paper ledger %s: %v", t.statePath, err)
	}
}

// checkFillable blocks fills at the daily limits or while suspended.
func checkFillable(q Quote, buying bool) error {
	if q.Volume <= 0 {
		return errors.New("security is suspended or has no trades")
	}
	if buying && q.UpLimit > 0 && q.Price >= q.UpLimit {
		return fmt.Errorf("price %.2f is at limit-up %.2f; buy order cannot fill", q.Price, q.UpLimit)
	}
	if !buying && q.DownLimit > 0 && q.Price <= q.DownLimit {
		return fmt.Errorf("price %.2f is at limit-down %.2f; sell order cannot fill", q.Price, q.DownLimit)
	}
	return nil
}

// checkProtectionsLocked closes positions whose stop-loss or take-profit has been reached.
// A triggered sell that T+1 or the limit rules block stays pending and is retried later.
func (t *Trader) checkProtectionsLocked() {
	for symbol, p := range t.positions {
		q, err := t.quote(symbol)
		if err != nil {
			continue
		}
		hitStop := p.stopLoss > 0 && q.Price <= p.stopLoss
		hitTarget := p.takeProfit > 0 && q.Price >= p.takeProfit
		if !hitStop && !hitTarget {
			continue
		}
		reason := "take_profit"
		if hitStop {
			reason = "stop_loss"
		}
		_, _ = t.closeLocked(symbol, 0, reason)
	}
}

// SetLeverage accepts only 1x, since A-shares are unlevered.
func (t *Trader) SetLeverage(symbol string, leverage int) error {
	return types.ValidateSideAllowed(t.caps, leverage, false)
}

// SetMarginMode is a no-op: cash accounts have no margin mode.
func (t *Trader) SetMarginMode(symbol string, isCrossMargin bool) error { return nil }

// GetMarketPrice returns the latest quote price.
func (t *Trader) GetMarketPrice(symbol string) (float64, error) {
	q, err := t.quote(symbol)
	if err != nil {
		return 0, err
	}
	return q.Price, nil
}

// SetStopLoss records a stop-loss on an open position. It is checked on each read.
func (t *Trader) SetStopLoss(symbol string, positionSide string, quantity, stopPrice float64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.positions[symbol]
	if !ok {
		return fmt.Errorf("no open position for %s", symbol)
	}
	if stopPrice <= 0 {
		return fmt.Errorf("stop price must be positive")
	}
	p.stopLoss = stopPrice
	t.persist()
	return nil
}

// SetTakeProfit records a take-profit on an open position.
func (t *Trader) SetTakeProfit(symbol string, positionSide string, quantity, takeProfitPrice float64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.positions[symbol]
	if !ok {
		return fmt.Errorf("no open position for %s", symbol)
	}
	if takeProfitPrice <= 0 {
		return fmt.Errorf("take-profit price must be positive")
	}
	p.takeProfit = takeProfitPrice
	t.persist()
	return nil
}

// CancelStopLossOrders clears the stop-loss on every position for the symbol.
func (t *Trader) CancelStopLossOrders(symbol string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if p, ok := t.positions[symbol]; ok {
		p.stopLoss = 0
	}
	t.persist()
	return nil
}

// CancelTakeProfitOrders clears the take-profit for the symbol.
func (t *Trader) CancelTakeProfitOrders(symbol string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if p, ok := t.positions[symbol]; ok {
		p.takeProfit = 0
	}
	t.persist()
	return nil
}

// CancelAllOrders clears protections for the symbol. Paper orders always fill at once.
func (t *Trader) CancelAllOrders(symbol string) error {
	if err := t.CancelStopLossOrders(symbol); err != nil {
		return err
	}
	return t.CancelTakeProfitOrders(symbol)
}

// CancelStopOrders clears both protections for the symbol.
func (t *Trader) CancelStopOrders(symbol string) error { return t.CancelAllOrders(symbol) }

// FormatQuantity returns a whole-lot quantity as text, since A-shares trade in whole shares.
func (t *Trader) FormatQuantity(symbol string, quantity float64) (string, error) {
	return fmt.Sprintf("%d", int64(types.FloorToLot(t.caps, quantity))), nil
}

// GetOrderStatus returns the recorded fill for a paper order.
func (t *Trader) GetOrderStatus(symbol string, orderID string) (map[string]interface{}, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	f, ok := t.fills[orderID]
	if !ok {
		return nil, fmt.Errorf("unknown order %s", orderID)
	}
	return map[string]interface{}{
		"status":      "FILLED",
		"avgPrice":    f.price,
		"executedQty": f.quantity,
		"commission":  f.fee,
	}, nil
}

// GetClosedPnL returns closed-trade records on or after startTime, newest last.
func (t *Trader) GetClosedPnL(startTime time.Time, limit int) ([]types.ClosedPnLRecord, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []types.ClosedPnLRecord
	for _, r := range t.closed {
		if r.ExitTime.Before(startTime) {
			continue
		}
		out = append(out, r)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

// GetOpenOrders returns the stop-loss and take-profit protections that are set.
func (t *Trader) GetOpenOrders(symbol string) ([]types.OpenOrder, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []types.OpenOrder
	for _, p := range t.positions {
		if symbol != "" && p.symbol != symbol {
			continue
		}
		if p.stopLoss > 0 {
			out = append(out, types.OpenOrder{Symbol: p.symbol, Side: "SELL", PositionSide: "LONG", Type: "STOP_MARKET", StopPrice: p.stopLoss, Quantity: p.quantity, Status: "NEW"})
		}
		if p.takeProfit > 0 {
			out = append(out, types.OpenOrder{Symbol: p.symbol, Side: "SELL", PositionSide: "LONG", Type: "TAKE_PROFIT_MARKET", StopPrice: p.takeProfit, Quantity: p.quantity, Status: "NEW"})
		}
	}
	return out, nil
}

func (t *Trader) newOrderID() string {
	t.nextID++
	return fmt.Sprintf("cnpaper-%d", t.nextID)
}

func (t *Trader) recordFill(id, symbol, side string, qty, price, fee float64) {
	t.fills[id] = fill{orderID: id, symbol: symbol, side: side, quantity: qty, price: price, fee: fee, time: t.now()}
}

func orderResult(id string, price, qty float64) map[string]interface{} {
	return map[string]interface{}{
		"orderId":     id,
		"status":      "FILLED",
		"avgPrice":    price,
		"executedQty": qty,
	}
}
