package trader

import (
	"fmt"
	"nofx/kernel"
	"nofx/trader/types"
	"time"
)

// marketCapabilities returns the rules for this trader's market. A trader with
// no market has zero capabilities, which allow no shorts and no leverage.
func (at *AutoTrader) marketCapabilities() types.Capabilities {
	return at.capabilities
}

// applyMarketRules enforces the market's rules before an open order is sent.
// It returns the quantity to send: buys are floored to the market's lot size.
func (at *AutoTrader) applyMarketRules(decision *kernel.Decision, quantity float64, isShort bool, now time.Time) (float64, error) {
	caps := at.marketCapabilities()
	if err := types.ValidateSideAllowed(caps, decision.Leverage, isShort); err != nil {
		return 0, err
	}
	open, err := caps.InTradingSession(now)
	if err != nil {
		return 0, err
	}
	if !open {
		return 0, fmt.Errorf("%s is outside the regular trading session", caps.Market)
	}
	qty := types.FloorToLot(caps, quantity)
	if err := types.ValidateOrderQuantity(caps, qty, true, false); err != nil {
		return 0, err
	}
	return qty, nil
}

// enforceSettlement blocks a close that the market's settlement rules forbid.
func (at *AutoTrader) enforceSettlement(entryTime, now time.Time) error {
	return types.ValidateSettlement(at.marketCapabilities(), entryTime, now)
}

// validateProtectionPrices checks the stop-loss and take-profit of an open order
// against the market price. The stop must sit on the loss side and the take-profit
// on the profit side.
func validateProtectionPrices(action string, marketPrice, stopLoss, takeProfit float64) error {
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
