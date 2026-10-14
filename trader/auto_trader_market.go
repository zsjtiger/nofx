package trader

import (
	"nofx/kernel"
	"nofx/trader/types"
)

// marketCapabilities returns the rules for this trader's market. A trader built
// without a market (for example in tests) gets the legacy crypto rules.
func (at *AutoTrader) marketCapabilities() types.Capabilities {
	if at.capabilities.Market == "" {
		caps, _ := types.CapabilitiesOf(types.MarketCrypto)
		return caps
	}
	return at.capabilities
}

// applyMarketRules enforces the market's rules before an open order is sent.
// It returns the quantity to send: buys are floored to the market's lot size.
func (at *AutoTrader) applyMarketRules(decision *kernel.Decision, quantity float64, isShort bool) (float64, error) {
	caps := at.marketCapabilities()
	if err := types.ValidateSideAllowed(caps, decision.Leverage, isShort); err != nil {
		return 0, err
	}
	qty := types.FloorToLot(caps, quantity)
	if err := types.ValidateOrderQuantity(caps, qty, true, false); err != nil {
		return 0, err
	}
	return qty, nil
}
