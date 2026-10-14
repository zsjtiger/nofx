package trader

import (
	"crux/kernel"
	"crux/market"
	"crux/store"
	"crux/trader/rules"
	"fmt"
	"strings"
	"time"
)

const (
	// Exit gates, validated by decision replay (2026-07-26, 4154 cycles,
	// 3-fold robustness): gates beat no-gates by 34 pts and the old rigid
	// 4h/8h by 16 pts of worst-fold score; the searched optimum sits at these
	// values. Thresholds are price-move percentages.
	autopilotMinHoldDuration        = 90 * time.Minute
	autopilotNoiseCloseHoldDuration = 3 * time.Hour
	// Re-entering a just-closed symbol was a consistent loss source: the
	// replay's top-20 configs cluster tightly at ~4h.
	autopilotReentryCooldown      = 4 * time.Hour
	earlyCloseStopLossBypassPct   = -3.0
	earlyCloseTakeProfitBypassPct = 8.0
	noiseCloseLossFloorPct        = -2.0
	noiseCloseProfitCeilingPct    = 3.0
)

func (at *AutoTrader) normalizedDecisionSymbol(symbol string) string {
	return market.NormalizeFor(at.marketCapabilities().Market, strings.TrimSpace(symbol))
}

func (at *AutoTrader) tradeThrottleReason(decision kernel.Decision, ctx *kernel.Context) string {
	if ctx == nil {
		return ""
	}

	switch {
	case rules.IsOpenAction(decision.Action):
		return at.openThrottleReason(decision, ctx)
	case rules.IsCloseAction(decision.Action):
		return at.closeThrottleReason(decision, ctx)
	default:
		return ""
	}
}

func (at *AutoTrader) openThrottleReason(decision kernel.Decision, ctx *kernel.Context) string {
	symbol := at.normalizedDecisionSymbol(decision.Symbol)
	if symbol == "" {
		return ""
	}

	if pos := rules.FindAnyContextPosition(ctx, symbol); pos != nil {
		return fmt.Sprintf("trade throttle: %s already has an open %s position; manage or close it before opening another side", symbol, pos.Side)
	}

	if order := at.findRecentCloseOrder(symbol, time.Now().Add(-autopilotReentryCooldown)); order != nil {
		age := time.Since(time.UnixMilli(order.CreatedAt))
		remaining := autopilotReentryCooldown - age
		if remaining < 0 {
			remaining = 0
		}
		return fmt.Sprintf("trade throttle: %s was closed %s ago; wait %s before re-entry", symbol, rules.RoundDuration(age), rules.RoundDuration(remaining))
	}

	return ""
}

func (at *AutoTrader) closeThrottleReason(decision kernel.Decision, ctx *kernel.Context) string {
	symbol := at.normalizedDecisionSymbol(decision.Symbol)
	side := rules.CloseActionSide(decision.Action)
	if symbol == "" || side == "" {
		return ""
	}

	pos := rules.FindContextPosition(ctx, symbol, side)
	pnlPct := 0.0
	entryTime := int64(0)
	if pos != nil {
		pnlPct = pos.UnrealizedPnLPct
		entryTime = pos.UpdateTime
	}

	if order := at.findRecentOpenOrder(symbol, side, time.Now().Add(-autopilotNoiseCloseHoldDuration)); order != nil && order.CreatedAt > entryTime {
		entryTime = order.CreatedAt
	}
	if entryTime <= 0 {
		return ""
	}

	heldFor := time.Since(time.UnixMilli(entryTime))
	if heldFor < 0 {
		heldFor = 0
	}
	if heldFor >= autopilotMinHoldDuration {
		if heldFor >= autopilotNoiseCloseHoldDuration ||
			pnlPct <= noiseCloseLossFloorPct ||
			pnlPct >= noiseCloseProfitCeilingPct {
			return ""
		}

		remaining := autopilotNoiseCloseHoldDuration - heldFor
		return fmt.Sprintf(
			"trade throttle: %s %s has been held for %s with price PnL %.2f%%; it is still inside the noise band %.1f%% to %.1f%%, so wait about %s before a flat/small close",
			symbol,
			side,
			rules.RoundDuration(heldFor),
			pnlPct,
			noiseCloseLossFloorPct,
			noiseCloseProfitCeilingPct,
			rules.RoundDuration(remaining),
		)
	}

	// Do not block true risk exits or unusually strong take-profit exits.
	if pnlPct <= earlyCloseStopLossBypassPct || pnlPct >= earlyCloseTakeProfitBypassPct {
		return ""
	}

	remaining := autopilotMinHoldDuration - heldFor
	return fmt.Sprintf(
		"trade throttle: %s %s has only been held for %s with price PnL %.2f%%; min AI-managed hold is %s unless price loss <= %.1f%% or price profit >= %.1f%%",
		symbol,
		side,
		rules.RoundDuration(heldFor),
		pnlPct,
		rules.RoundDuration(autopilotMinHoldDuration),
		earlyCloseStopLossBypassPct,
		earlyCloseTakeProfitBypassPct,
	) + fmt.Sprintf("; wait about %s", rules.RoundDuration(remaining))
}

func (at *AutoTrader) recentOrders(limit int) ([]*store.TraderOrder, error) {
	if at == nil || at.store == nil {
		return nil, nil
	}
	return at.store.Order().GetTraderOrders(at.id, limit)
}

func (at *AutoTrader) findRecentCloseOrder(symbol string, since time.Time) *store.TraderOrder {
	orders, err := at.recentOrders(100)
	if err != nil {
		at.logWarnf("⚠️ Trade throttle could not read recent close orders: %v", err)
		return nil
	}
	sinceMs := since.UTC().UnixMilli()
	for _, order := range orders {
		if order == nil || order.CreatedAt < sinceMs || isCanceledOrder(order) {
			continue
		}
		if at.normalizedDecisionSymbol(order.Symbol) == symbol && rules.IsCloseAction(order.OrderAction) {
			return order
		}
	}
	return nil
}

func (at *AutoTrader) findRecentOpenOrder(symbol string, side string, since time.Time) *store.TraderOrder {
	orders, err := at.recentOrders(100)
	if err != nil {
		at.logWarnf("⚠️ Trade throttle could not read recent open orders: %v", err)
		return nil
	}
	sinceMs := since.UTC().UnixMilli()
	for _, order := range orders {
		if order == nil || order.CreatedAt < sinceMs || isCanceledOrder(order) {
			continue
		}
		if at.normalizedDecisionSymbol(order.Symbol) == symbol &&
			strings.EqualFold(rules.OpenActionSide(order.OrderAction), side) {
			return order
		}
	}
	return nil
}

func isCanceledOrder(order *store.TraderOrder) bool {
	status := strings.ToUpper(strings.TrimSpace(order.Status))
	return status == "CANCELED" || status == "CANCELLED" || status == "REJECTED" || status == "EXPIRED"
}
