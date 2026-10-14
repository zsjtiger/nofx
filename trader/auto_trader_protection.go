package trader

import (
	"fmt"
	"nofx/logger"
	"nofx/trader/types"
	"strings"
)

func (at *AutoTrader) closeUnprotectedPosition(symbol, side string, quantity float64, protectionErr error) error {
	logger.Infof("  🚨 %v; emergency-closing %s %s", protectionErr, symbol, side)
	if closeErr := at.emergencyClosePositionAndVerify(symbol, side, quantity); closeErr != nil {
		return fmt.Errorf("%w; emergency close failed: %v", protectionErr, closeErr)
	}
	return fmt.Errorf("%w; opened position was emergency-closed", protectionErr)
}

func (at *AutoTrader) emergencyClosePositionAndVerify(symbol, side string, quantity float64) error {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		var err error
		if side == "long" {
			_, err = at.trader.CloseLong(symbol, quantity)
		} else if side == "short" {
			_, err = at.trader.CloseShort(symbol, quantity)
		} else {
			return fmt.Errorf("unknown position direction: %s", side)
		}
		if err != nil {
			lastErr = err
		} else {
			positions, err := getFreshPositions(at.trader)
			if err != nil {
				lastErr = fmt.Errorf("verify flat position: %w", err)
			} else {
				residual := false
				for _, pos := range positions {
					if universeBaseKey(fmt.Sprint(pos["symbol"])) == universeBaseKey(symbol) && strings.EqualFold(fmt.Sprint(pos["side"]), side) {
						residual = true
						break
					}
				}
				if !residual {
					if err := at.trader.CancelAllOrders(symbol); err != nil {
						return fmt.Errorf("position is flat but protection-order cleanup failed: %w", err)
					}
					openOrders, err := at.trader.GetOpenOrders(symbol)
					if err != nil {
						return fmt.Errorf("position is flat but protection-order verification failed: %w", err)
					}
					if len(openOrders) != 0 {
						return fmt.Errorf("position is flat but %d open orders remain for %s", len(openOrders), symbol)
					}
					return nil
				}
				lastErr = fmt.Errorf("residual %s position remains after close attempt %d", side, attempt)
			}
		}
	}
	return lastErr
}

// positionCacheInvalidator is implemented by traders that cache positions.
type positionCacheInvalidator interface {
	InvalidatePositionCache()
}

func getFreshPositions(tr types.Trader) ([]map[string]interface{}, error) {
	if invalidator, ok := tr.(positionCacheInvalidator); ok {
		invalidator.InvalidatePositionCache()
	}
	return tr.GetPositions()
}

// enforceVergexSignalPolicy turns the current direction board into a strict
// position state machine. Detail data can explain a signal, but cannot reverse
// or prematurely exit it.
