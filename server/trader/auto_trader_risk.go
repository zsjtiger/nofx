package trader

import (
	"fmt"
	"math"
	"nofx/logger"
	"nofx/store"
	"nofx/trader/rules"
	"time"
)

const (
	// The monitor arms only once the underlying PRICE has moved +5% in the
	// position's favor (leverage-independent — at 10x the old margin-basis
	// check armed at a +0.5% price wiggle and strangled every winner), then
	// closes if the position gives back 40% of its peak profit.
	drawdownClosePriceGainPct = 5.0
	drawdownCloseGivebackPct  = 40.0
	// Keep total planned margin below 100% so every configured slot retains
	// room for exchange overhead and fees. This is an allocation target, not
	// the per-position hard cap in RiskControlConfig.
	autopilotBookMarginBudget = 0.96
)

// shouldDrawdownClose reports whether the profit-protection close should fire.
// pricePnLPct is the price-basis move in the position's favor; drawdownPct is
// the relative giveback from the position's peak profit.
func shouldDrawdownClose(pricePnLPct, drawdownPct float64) bool {
	return pricePnLPct > drawdownClosePriceGainPct && drawdownPct >= drawdownCloseGivebackPct
}

// startDrawdownMonitor starts drawdown monitoring
func (at *AutoTrader) startDrawdownMonitor() {
	at.monitorWg.Add(1)
	go func() {
		defer at.monitorWg.Done()

		ticker := time.NewTicker(1 * time.Minute) // Check every minute
		defer ticker.Stop()

		logger.Info("📊 Started position drawdown monitoring (check every minute)")

		for {
			select {
			case <-ticker.C:
				at.checkPositionDrawdown()
			case <-at.stopMonitorCh:
				logger.Info("⏹ Stopped position drawdown monitoring")
				return
			}
		}
	}()
}

// checkPositionDrawdown checks position drawdown situation
func (at *AutoTrader) checkPositionDrawdown() {
	if !at.profitProtectionEnabled() {
		return
	}
	// Get current positions
	positions, err := at.trader.GetPositions()
	if err != nil {
		logger.Infof("❌ Drawdown monitoring: failed to get positions: %v", err)
		return
	}

	for _, pos := range positions {
		symbol := pos["symbol"].(string)
		side := pos["side"].(string)
		entryPrice := pos["entryPrice"].(float64)
		markPrice := pos["markPrice"].(float64)
		quantity := pos["positionAmt"].(float64)
		if quantity < 0 {
			quantity = -quantity // Short position quantity is negative, convert to positive
		}

		// Guard: skip if entry price is zero (prevents division by zero panic)
		if entryPrice <= 0 {
			logger.Warnf("⚠️ Drawdown monitoring: %s %s has zero entry price, skipping", symbol, side)
			continue
		}

		// Calculate current P&L percentage
		leverage := 10 // Default value
		if lev, ok := pos["leverage"].(float64); ok {
			leverage = int(lev)
		}

		// Price-basis move drives the close decision so the trigger point does
		// not tighten as leverage grows; the margin-basis (leveraged) value is
		// only kept for the peak cache shown alongside margin-based PnL% in
		// prompts.
		var pricePnLPct float64
		if side == "long" {
			pricePnLPct = ((markPrice - entryPrice) / entryPrice) * 100
		} else {
			pricePnLPct = ((entryPrice - markPrice) / entryPrice) * 100
		}
		currentPnLPct := pricePnLPct * float64(leverage)

		// Construct unique position identifier (distinguish long/short)
		posKey := symbol + "_" + side

		// Get historical peak profit for this position
		at.peakPnLCacheMutex.RLock()
		peakPnLPct, exists := at.peakPnLCache[posKey]
		at.peakPnLCacheMutex.RUnlock()

		if !exists {
			// If no historical peak record, use current P&L as initial value
			peakPnLPct = currentPnLPct
			at.UpdatePeakPnL(symbol, side, currentPnLPct)
		} else {
			// Update peak cache
			at.UpdatePeakPnL(symbol, side, currentPnLPct)
		}

		// Calculate drawdown (magnitude of decline from peak)
		var drawdownPct float64
		if peakPnLPct > 0 && currentPnLPct < peakPnLPct {
			drawdownPct = ((peakPnLPct - currentPnLPct) / peakPnLPct) * 100
		}

		// Check close position condition: price move > +5% and drawdown >= 40%
		if shouldDrawdownClose(pricePnLPct, drawdownPct) {
			logger.Infof("🚨 Drawdown close position condition triggered: %s %s | Price move: %.2f%% | Current profit: %.2f%% | Peak profit: %.2f%% | Drawdown: %.2f%%",
				symbol, side, pricePnLPct, currentPnLPct, peakPnLPct, drawdownPct)

			// Execute close position
			if err := at.emergencyClosePosition(symbol, side); err != nil {
				logger.Infof("❌ Drawdown close position failed (%s %s): %v", symbol, side, err)
			} else {
				logger.Infof("✅ Drawdown close position succeeded: %s %s", symbol, side)
				// Clear cache for this position after closing
				at.ClearPeakPnLCache(symbol, side)
			}
		} else if pricePnLPct > drawdownClosePriceGainPct {
			// Record situations close to close position condition (for debugging)
			logger.Infof("📊 Drawdown monitoring: %s %s | Price move: %.2f%% | Profit: %.2f%% | Peak: %.2f%% | Drawdown: %.2f%%",
				symbol, side, pricePnLPct, currentPnLPct, peakPnLPct, drawdownPct)
		}
	}
}

// emergencyClosePosition emergency close position function
func (at *AutoTrader) emergencyClosePosition(symbol, side string) error {
	switch side {
	case "long":
		order, err := at.trader.CloseLong(symbol, 0) // 0 = close all
		if err != nil {
			return err
		}
		logger.Infof("✅ Emergency close long position succeeded, order ID: %v", order["orderId"])
	case "short":
		order, err := at.trader.CloseShort(symbol, 0) // 0 = close all
		if err != nil {
			return err
		}
		logger.Infof("✅ Emergency close short position succeeded, order ID: %v", order["orderId"])
	default:
		return fmt.Errorf("unknown position direction: %s", side)
	}

	return nil
}

// GetPeakPnLCache gets peak profit cache
func (at *AutoTrader) GetPeakPnLCache() map[string]float64 {
	at.peakPnLCacheMutex.RLock()
	defer at.peakPnLCacheMutex.RUnlock()

	// Return a copy of the cache
	cache := make(map[string]float64)
	for k, v := range at.peakPnLCache {
		cache[k] = v
	}
	return cache
}

// UpdatePeakPnL updates peak profit cache
func (at *AutoTrader) UpdatePeakPnL(symbol, side string, currentPnLPct float64) {
	at.peakPnLCacheMutex.Lock()
	defer at.peakPnLCacheMutex.Unlock()

	posKey := symbol + "_" + side
	if peak, exists := at.peakPnLCache[posKey]; exists {
		// Update peak (if long, take larger value; if short, currentPnLPct is negative, also compare)
		if currentPnLPct > peak {
			at.peakPnLCache[posKey] = currentPnLPct
		}
	} else {
		// First time recording
		at.peakPnLCache[posKey] = currentPnLPct
	}
}

// ClearPeakPnLCache clears peak cache for specified position
func (at *AutoTrader) ClearPeakPnLCache(symbol, side string) {
	at.peakPnLCacheMutex.Lock()
	defer at.peakPnLCacheMutex.Unlock()

	posKey := symbol + "_" + side
	delete(at.peakPnLCache, posKey)
}

// ============================================================================
// Risk Control Helpers
// ============================================================================

// enforcePositionValueRatio checks and enforces position value ratio limits (CODE ENFORCED)
// Returns the adjusted position size (capped if necessary) and whether the position was capped
// positionSizeUSD: the original position size in USD
// equity: the account equity
// symbol: the trading symbol
func (at *AutoTrader) enforcePositionValueRatio(positionSizeUSD float64, equity float64, symbol string) (float64, bool) {
	if at.config.StrategyConfig == nil {
		return positionSizeUSD, false
	}

	riskControl := at.config.StrategyConfig.RiskControl

	// Stocks use one position-value tier: equity times the altcoin ratio.
	maxPositionValueRatio := riskControl.AltcoinMaxPositionValueRatio
	if maxPositionValueRatio <= 0 {
		maxPositionValueRatio = 1.0
	}

	// Calculate max allowed position value = equity × ratio
	maxPositionValue := equity * maxPositionValueRatio

	// Check if position size exceeds limit
	if positionSizeUSD > maxPositionValue {
		logger.Infof("  ⚠️ [RISK CONTROL] Position %.2f %s exceeds limit (equity %.2f × %.1fx = %.2f %s max for %s), capping",
			positionSizeUSD, at.marketCapabilities().Currency, equity, maxPositionValueRatio, maxPositionValue, at.marketCapabilities().Currency, symbol)
		return maxPositionValue, true
	}

	return positionSizeUSD, false
}

// minPositionSize returns the smallest position the strategy allows.
func (at *AutoTrader) minPositionSize() float64 {
	if at.config.StrategyConfig == nil || at.config.StrategyConfig.RiskControl.MinPositionSize <= 0 {
		return store.MinPositionValue
	}
	return at.config.StrategyConfig.RiskControl.MinPositionSize
}

// enforceMinPositionSize checks minimum position size (CODE ENFORCED)
func (at *AutoTrader) enforceMinPositionSize(positionSizeUSD float64) error {
	minSize := at.minPositionSize()
	if positionSizeUSD < minSize {
		return fmt.Errorf("❌ [RISK CONTROL] Position %.2f %s below minimum (%.2f %s)", positionSizeUSD, at.marketCapabilities().Currency, minSize, at.marketCapabilities().Currency)
	}
	return nil
}

// minRiskReward returns the strategy's minimum take-profit to stop-loss ratio, or 0 when unset.
func (at *AutoTrader) minRiskReward() float64 {
	if at.config.StrategyConfig == nil {
		return 0
	}
	return at.config.StrategyConfig.RiskControl.MinRiskRewardRatio
}

// currentExposure sums the market value of the open positions.
func currentExposure(positions []map[string]interface{}) float64 {
	total := 0.0
	for _, pos := range positions {
		qty, _ := pos["positionAmt"].(float64)
		price, _ := pos["markPrice"].(float64)
		total += math.Abs(qty) * price
	}
	return total
}

// enforceExposureLimit caps a new position so the total open value stays within
// equity × MaxMarginUsage. It rejects the order when no room is left for the minimum size.
func (at *AutoTrader) enforceExposureLimit(positions []map[string]interface{}, positionSizeUSD, equity float64) (float64, error) {
	if at.config.StrategyConfig == nil {
		return positionSizeUSD, nil
	}
	exposure := currentExposure(positions)
	room := rules.ExposureRoom(exposure, equity, at.config.StrategyConfig.RiskControl.MaxMarginUsage)
	if positionSizeUSD <= room {
		return positionSizeUSD, nil
	}
	logger.Infof("  ⚠️ [RISK CONTROL] Position %.2f %s exceeds exposure room %.2f %s (open %.2f, equity %.2f), capping",
		positionSizeUSD, at.marketCapabilities().Currency, room, at.marketCapabilities().Currency, exposure, equity)
	if room < at.minPositionSize() {
		return 0, fmt.Errorf("❌ [RISK CONTROL] No exposure room left (open %.2f %s, limit %.2f %s)",
			exposure, at.marketCapabilities().Currency, equity*at.config.StrategyConfig.RiskControl.MaxMarginUsage, at.marketCapabilities().Currency)
	}
	return room, nil
}

// enforceMaxPositions checks maximum positions count (CODE ENFORCED)
func (at *AutoTrader) enforceMaxPositions(currentPositionCount int) error {
	if at.config.StrategyConfig == nil {
		return nil
	}

	maxPositions := at.config.StrategyConfig.RiskControl.MaxPositions
	if maxPositions <= 0 {
		maxPositions = 3 // Default: 3 positions
	}

	if currentPositionCount >= maxPositions {
		return fmt.Errorf("❌ [RISK CONTROL] Already at max positions (%d/%d)", currentPositionCount, maxPositions)
	}
	return nil
}

// getSideFromAction converts order action to side (BUY/SELL)
func getSideFromAction(action string) string {
	switch action {
	case "open_long", "close_short":
		return "BUY"
	case "open_short", "close_long":
		return "SELL"
	default:
		return "BUY"
	}
}

// profitProtectionEnabled reports whether the strategy turns on profit protection.
func (at *AutoTrader) profitProtectionEnabled() bool {
	return at.config.StrategyConfig != nil && at.config.StrategyConfig.RiskControl.ProfitProtection
}
