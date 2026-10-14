package rules

import (
	"fmt"
	"math"
)

// ValidateRiskReward checks that the distance to the take-profit is at least
// minRatio times the distance to the stop-loss, measured from the market price.
// A minRatio of zero or less disables the check.
func ValidateRiskReward(price, stopLoss, takeProfit, minRatio float64) error {
	if minRatio <= 0 {
		return nil
	}
	risk := math.Abs(price - stopLoss)
	reward := math.Abs(takeProfit - price)
	if risk <= 0 {
		return fmt.Errorf("stop loss %.4f must differ from market price %.4f", stopLoss, price)
	}
	if ratio := reward / risk; ratio+1e-9 < minRatio {
		return fmt.Errorf("risk/reward %.2f is below the minimum %.2f", ratio, minRatio)
	}
	return nil
}

// ExposureRoom returns how much new position value fits under the gross
// exposure limit equity × maxUsage, given the value already held. It never
// returns a negative number. A maxUsage of zero or less means 1 (no leverage).
func ExposureRoom(currentExposure, equity, maxUsage float64) float64 {
	if maxUsage <= 0 {
		maxUsage = 1
	}
	room := equity*maxUsage - currentExposure
	if room < 0 {
		return 0
	}
	return room
}
