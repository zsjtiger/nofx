// Package rules holds pure trading rules and order helpers used by the
// trader engine. Nothing here depends on AutoTrader state.
package rules

import (
	"fmt"
	"strings"
	"time"

	"nofx/kernel"
)

// ValidateProtectionPrices checks the stop-loss and take-profit of an open order
// against the market price. The stop must sit on the loss side and the take-profit
// on the profit side.
func ValidateProtectionPrices(action string, marketPrice, stopLoss, takeProfit float64) error {
	if marketPrice <= 0 || stopLoss <= 0 {
		return fmt.Errorf("market price and stop loss must be positive")
	}
	switch action {
	case "open_long":
		if stopLoss >= marketPrice {
			return fmt.Errorf("long stop loss %.8f must be below market price %.8f", stopLoss, marketPrice)
		}
		if takeProfit <= marketPrice {
			return fmt.Errorf("long take profit %.8f must be above market price %.8f", takeProfit, marketPrice)
		}
	case "open_short":
		if stopLoss <= marketPrice {
			return fmt.Errorf("short stop loss %.8f must be above market price %.8f", stopLoss, marketPrice)
		}
		if takeProfit <= 0 || takeProfit >= marketPrice {
			return fmt.Errorf("short take profit %.8f must be positive and below market price %.8f", takeProfit, marketPrice)
		}
	default:
		return fmt.Errorf("unsupported open action %q", action)
	}
	return nil
}

// PositionPricePnLPct converts the margin-based UnrealizedPnLPct reported for
// a position into the underlying price-move percentage.
func PositionPricePnLPct(pos *kernel.PositionInfo) float64 {
	if pos == nil {
		return 0
	}
	if pos.Leverage > 1 {
		return pos.UnrealizedPnLPct / float64(pos.Leverage)
	}
	return pos.UnrealizedPnLPct
}

func IsOpenAction(action string) bool {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "open_long", "open_short":
		return true
	default:
		return false
	}
}

func IsCloseAction(action string) bool {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "close_long", "close_short":
		return true
	default:
		return false
	}
}

func CloseActionSide(action string) string {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "close_long":
		return "long"
	case "close_short":
		return "short"
	default:
		return ""
	}
}

func OpenActionSide(action string) string {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case "open_long":
		return "long"
	case "open_short":
		return "short"
	default:
		return ""
	}
}

// CanonicalPositionSymbol normalizes a symbol reported by the exchange. Exchange
// symbols are already canonical, so only case and whitespace are normalised.
func CanonicalPositionSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))
}

func FindContextPosition(ctx *kernel.Context, symbol string, side string) *kernel.PositionInfo {
	if ctx == nil {
		return nil
	}
	for i := range ctx.Positions {
		pos := &ctx.Positions[i]
		if CanonicalPositionSymbol(pos.Symbol) == symbol && strings.EqualFold(pos.Side, side) {
			return pos
		}
	}
	return nil
}

func FindAnyContextPosition(ctx *kernel.Context, symbol string) *kernel.PositionInfo {
	if ctx == nil {
		return nil
	}
	for i := range ctx.Positions {
		pos := &ctx.Positions[i]
		if CanonicalPositionSymbol(pos.Symbol) == symbol {
			return pos
		}
	}
	return nil
}

func RoundDuration(d time.Duration) string {
	if d < time.Minute {
		return "0m"
	}
	return d.Round(time.Minute).String()
}
