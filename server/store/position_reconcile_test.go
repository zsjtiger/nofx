package store

import (
	"testing"
	"time"

	"crux/store/position"
)

func newReconcileTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := New(t.TempDir() + "/nofx.db")
	if err != nil {
		t.Fatalf("store.New failed: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func openRow(t *testing.T, st *Store, traderID, exchangeID, symbol, side string, qty, pnl float64, entryMs int64) int64 {
	t.Helper()
	pos := &position.TraderPosition{
		TraderID:      traderID,
		ExchangeID:    exchangeID,
		ExchangeType:  "alpaca_paper",
		Symbol:        symbol,
		Side:          side,
		Quantity:      qty,
		EntryQuantity: qty,
		EntryPrice:    100,
		EntryTime:     entryMs,
		RealizedPnL:   pnl,
		Status:        "OPEN",
		Source:        "sync",
		CreatedAt:     entryMs,
		UpdatedAt:     entryMs,
	}
	if err := st.Position().CreateOpenPosition(pos); err != nil {
		t.Fatalf("create open position: %v", err)
	}
	return pos.ID
}

func TestReconcileClosesZombiesKeepsLiveAndTrims(t *testing.T) {
	st := newReconcileTestStore(t)
	const exch = "ex-hl"
	base := time.Now().Add(-48 * time.Hour).UnixMilli()

	// Zombie under the CURRENT trader id: exchange holds nothing for IBM.
	zombieID := openRow(t, st, "trader-now", exch, "IBM", "LONG", 6.8, -20.34, base)

	// Zombie left by a PRIOR autopilot incarnation on the SAME exchange —
	// this is the case a per-trader-id reconcile would miss.
	legacyID := openRow(t, st, "trader-old", exch, "NVDA", "SHORT", 6.94, -3.5, base+500)

	// Duplicates for SPY (different incarnations): newest survives trimmed,
	// older closes.
	oldSP := openRow(t, st, "trader-old", exch, "SPY", "LONG", 0.07, -1.5, base+1000)
	newSP := openRow(t, st, "trader-now", exch, "SPY", "LONG", 0.124, 2.5, base+2000)

	// Healthy row exactly matching live — untouched.
	healthy := openRow(t, st, "trader-now", exch, "AAPL", "LONG", 0.01, 0, base+3000)

	// Row on a DIFFERENT exchange account — must be out of scope.
	otherExch := openRow(t, st, "trader-now", "ex-other", "MSFT", "LONG", 2.0, 0, base+4000)

	live := map[string]float64{
		position.LivePositionKey("SPY", "long"):  0.057,
		position.LivePositionKey("AAPL", "long"): 0.01,
	}

	closed, err := st.Position().ReconcileOpenPositionsWithLive(exch, live)
	if err != nil {
		t.Fatalf("reconcile failed: %v", err)
	}
	if closed != 3 {
		t.Fatalf("expected 3 zombies closed (IBM + legacy NVDA + old SPY), got %d", closed)
	}

	get := func(id int64) *position.TraderPosition {
		t.Helper()
		var pos position.TraderPosition
		if err := st.GormDB().First(&pos, id).Error; err != nil {
			t.Fatalf("load position %d: %v", id, err)
		}
		return &pos
	}

	if dram := get(zombieID); dram.Status != "CLOSED" || dram.RealizedPnL != -20.34 || dram.CloseReason != "reconcile" {
		t.Fatalf("IBM zombie should close via reconcile keeping PnL, got %+v", dram)
	}
	if sol := get(legacyID); sol.Status != "CLOSED" {
		t.Fatalf("legacy-incarnation NVDA zombie on same exchange should close, got %+v", sol)
	}
	if sp := get(oldSP); sp.Status != "CLOSED" {
		t.Fatalf("older duplicate SPY row should close, got %+v", sp)
	}
	if sp := get(newSP); sp.Status != "OPEN" || sp.Quantity > 0.0571 || sp.Quantity < 0.0569 {
		t.Fatalf("newest SPY row should stay open trimmed to live 0.057, got status=%s qty=%v", sp.Status, sp.Quantity)
	}
	if btc := get(healthy); btc.Status != "OPEN" || btc.Quantity != 0.01 {
		t.Fatalf("healthy row must be untouched, got %+v", btc)
	}
	if eth := get(otherExch); eth.Status != "OPEN" {
		t.Fatalf("row on a different exchange must be out of scope, got %+v", eth)
	}
}

func TestReconcileNoOpenRowsIsNoop(t *testing.T) {
	st := newReconcileTestStore(t)
	closed, err := st.Position().ReconcileOpenPositionsWithLive("ex-empty", map[string]float64{})
	if err != nil || closed != 0 {
		t.Fatalf("expected clean noop, got closed=%d err=%v", closed, err)
	}
}
