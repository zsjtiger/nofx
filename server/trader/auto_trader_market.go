package trader

import (
	"fmt"
	"nofx/kernel"
	"nofx/trader/types"
	"time"
)

// marketCapabilities returns the rules for this trader's market. A trader with
// no market has zero capabilities, which allow no shorts.
func (at *AutoTrader) marketCapabilities() types.Capabilities {
	return at.capabilities
}

// marketCurrency returns the account currency of a market, or "" when unknown.
func marketCurrency(m types.Market) string {
	caps, err := types.CapabilitiesOf(m)
	if err != nil {
		return ""
	}
	return caps.Currency
}

// applyMarketRules enforces the market's rules before an open order is sent.
// It returns the quantity to send: buys are floored to the market's lot size.
func (at *AutoTrader) applyMarketRules(decision *kernel.Decision, quantity float64, isShort bool, now time.Time) (float64, error) {
	caps := at.marketCapabilities()
	if err := types.ValidateSideAllowed(caps, isShort); err != nil {
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
