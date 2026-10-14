package api

import (
	"crux/trader/types"
	"fmt"
	"net/http"
	"strings"
	"time"

	"crux/logger"
	"crux/store"
	"crux/trader"

	"github.com/gin-gonic/gin"
)

// handleSyncBalance Sync exchange balance to initial_balance (Option B: Manual Sync + Option C: Smart Detection)
func (s *Server) handleSyncBalance(c *gin.Context) {
	userID := c.GetString("user_id")
	traderID := c.Param("id")

	logger.Infof("🔄 User %s requested balance sync for trader %s", userID, traderID)

	// Get trader configuration from database (including exchange info)
	fullConfig, err := s.store.Trader().GetFullConfig(userID, traderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trader does not exist"})
		return
	}

	traderConfig := fullConfig.Trader
	exchangeCfg := fullConfig.Exchange

	if exchangeCfg == nil || !exchangeCfg.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exchange not configured or not enabled"})
		return
	}

	tempTrader, createErr := buildExchangeProbeTrader(exchangeCfg, userID)
	if createErr != nil {
		logger.Infof("⚠️ Failed to create temporary trader: %v", createErr)
		SafeInternalError(c, "Failed to connect to exchange", createErr)
		return
	}

	// Query actual balance
	balanceInfo, balanceErr := tempTrader.GetBalance()
	if balanceErr != nil {
		logger.Infof("⚠️ Failed to query exchange balance: %v", balanceErr)
		SafeInternalError(c, "Failed to query balance", balanceErr)
		return
	}

	// Extract total equity (for P&L calculation, we need total account value, not available balance)
	actualBalance, found := extractExchangeTotalEquity(balanceInfo)
	if !found {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get total equity"})
		return
	}

	s.exchangeAccountStateCache.Invalidate(userID)

	oldBalance := traderConfig.InitialBalance

	// Smart balance change detection
	changePercent := ((actualBalance - oldBalance) / oldBalance) * 100
	changeType := "increase"
	if changePercent < 0 {
		changeType = "decrease"
	}

	logger.Infof("✓ Queried actual exchange balance: %.2f %s (current config: %.2f %s, change: %.2f%%)",
		actualBalance, brokerCurrency(exchangeCfg.ExchangeType), oldBalance, brokerCurrency(exchangeCfg.ExchangeType), changePercent)

	// Update initial_balance in database
	err = s.store.Trader().UpdateInitialBalance(userID, traderID, actualBalance)
	if err != nil {
		logger.Infof("❌ Failed to update initial_balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	// Reload traders into memory
	err = s.traderManager.LoadUserTradersFromStore(s.store, userID)
	if err != nil {
		logger.Infof("⚠️ Failed to reload user traders into memory: %v", err)
	}

	logger.Infof("✅ Synced balance: %.2f → %.2f %s (%s %.2f%%)", oldBalance, actualBalance, brokerCurrency(exchangeCfg.ExchangeType), changeType, changePercent)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Balance synced successfully",
		"old_balance":    oldBalance,
		"new_balance":    actualBalance,
		"change_percent": changePercent,
		"change_type":    changeType,
	})
}

// handleClosePosition One-click close position
func (s *Server) handleClosePosition(c *gin.Context) {
	userID := c.GetString("user_id")
	traderID := c.Param("id")

	var req struct {
		Symbol string `json:"symbol" binding:"required"`
		Side   string `json:"side" binding:"required"` // "LONG" or "SHORT"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter error: symbol and side are required"})
		return
	}

	logger.Infof("🔻 User %s requested position close: trader=%s, symbol=%s, side=%s", userID, traderID, req.Symbol, req.Side)

	// Get trader configuration from database (including exchange info)
	fullConfig, err := s.store.Trader().GetFullConfig(userID, traderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trader does not exist"})
		return
	}

	exchangeCfg := fullConfig.Exchange

	if exchangeCfg == nil || !exchangeCfg.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exchange not configured or not enabled"})
		return
	}

	// Create temporary trader to execute close position
	var tempTrader trader.Trader
	var createErr error

	// Use ExchangeType (e.g. "alpaca_paper") instead of ExchangeID (which is now UUID)
	// Convert EncryptedString fields to string
	tempTrader, createErr = buildExchangeProbeTrader(exchangeCfg, userID)

	if createErr != nil {
		logger.Infof("⚠️ Failed to create temporary trader: %v", createErr)
		SafeInternalError(c, "Failed to connect to exchange", createErr)
		return
	}

	// Get current position info BEFORE closing (to get quantity and price)
	positions, err := tempTrader.GetPositions()
	if err != nil {
		logger.Infof("⚠️ Failed to get positions: %v", err)
	}

	var posQty float64
	var entryPrice float64
	for _, pos := range positions {
		if pos["symbol"] == req.Symbol && pos["side"] == strings.ToLower(req.Side) {
			if amt, ok := pos["positionAmt"].(float64); ok {
				posQty = amt
				if posQty < 0 {
					posQty = -posQty // Make positive
				}
			}
			if price, ok := pos["entryPrice"].(float64); ok {
				entryPrice = price
			}
			break
		}
	}

	// Execute close position operation
	var result map[string]interface{}
	var closeErr error

	if req.Side == "LONG" {
		result, closeErr = tempTrader.CloseLong(req.Symbol, 0) // 0 means close all
	} else if req.Side == "SHORT" {
		result, closeErr = tempTrader.CloseShort(req.Symbol, 0) // 0 means close all
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "side must be LONG or SHORT"})
		return
	}

	if closeErr != nil {
		logger.Infof("❌ Close position failed: symbol=%s, side=%s, error=%v", req.Symbol, req.Side, closeErr)
		SafeInternalError(c, "Close position", closeErr)
		return
	}

	logger.Infof("✅ Position closed successfully: symbol=%s, side=%s, qty=%.6f, result=%v", req.Symbol, req.Side, posQty, result)

	// Backfill the just-closed fill immediately. Manual closes may happen while
	// the bot runtime is stopped, so the background OrderSync loop is not enough.
	if syncErr := s.syncOrdersAfterManualClose(tempTrader, traderID, exchangeCfg.ID, exchangeCfg.ExchangeType); syncErr != nil {
		logger.Infof("  ⚠️ Manual close sync failed: %v", syncErr)
		s.recordClosePositionOrder(traderID, exchangeCfg.ID, exchangeCfg.ExchangeType, req.Symbol, req.Side, posQty, entryPrice, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Position closed successfully",
		"symbol":  req.Symbol,
		"side":    req.Side,
		"result":  result,
	})
}

func (s *Server) syncOrdersFromExchange(exchangeTrader trader.Trader, traderID, exchangeID, exchangeType string) error {
	return fmt.Errorf("order sync is not available for exchange type %s", exchangeType)
}

func (s *Server) syncOrdersAfterManualClose(exchangeTrader trader.Trader, traderID, exchangeID, exchangeType string) error {
	var lastErr error
	for attempt := 1; attempt <= 4; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt-1) * 500 * time.Millisecond)
		}
		if err := s.syncOrdersFromExchange(exchangeTrader, traderID, exchangeID, exchangeType); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("manual close sync did not run")
}

// recordClosePositionOrder records a manual close. Paper brokers fill synchronously,
// so the order is written here with its fill.
func (s *Server) recordClosePositionOrder(traderID, exchangeID, exchangeType, symbol, side string, quantity, exitPrice float64, result map[string]interface{}) {
	// Check if order was placed (skip if NO_POSITION)
	status, _ := result["status"].(string)
	if status == "NO_POSITION" {
		logger.Infof("  ⚠️ No position to close, skipping order record")
		return
	}

	// Get order ID from result
	var orderID string
	switch v := result["orderId"].(type) {
	case int64:
		orderID = fmt.Sprintf("%d", v)
	case float64:
		orderID = fmt.Sprintf("%.0f", v)
	case string:
		orderID = v
	default:
		orderID = fmt.Sprintf("%v", v)
	}

	if orderID == "" || orderID == "0" {
		logger.Infof("  ⚠️ Order ID is empty, skipping record")
		return
	}

	// Determine order action based on side
	var orderAction string
	if side == "LONG" {
		orderAction = "close_long"
	} else {
		orderAction = "close_short"
	}

	// Use entry price if exit price not available
	if exitPrice == 0 {
		exitPrice = quantity * 100 // Rough estimate if we don't have price
	}

	// Estimate fee
	fee := exitPrice * quantity * 0.0004

	// Create order record - DIRECTLY as FILLED (market orders fill immediately)
	orderRecord := &store.TraderOrder{
		TraderID:        traderID,
		ExchangeID:      exchangeID,
		ExchangeType:    exchangeType,
		ExchangeOrderID: orderID,
		Symbol:          symbol,
		PositionSide:    side,
		OrderAction:     orderAction,
		Type:            "MARKET",
		Side:            getSideFromAction(orderAction),
		Quantity:        quantity,
		Price:           0, // Market order
		Status:          "FILLED",
		FilledQuantity:  quantity,
		AvgFillPrice:    exitPrice,
		Commission:      fee,
		FilledAt:        time.Now().UTC().UnixMilli(),
		CreatedAt:       time.Now().UTC().UnixMilli(),
		UpdatedAt:       time.Now().UTC().UnixMilli(),
	}

	if err := s.store.Order().CreateOrder(orderRecord); err != nil {
		logger.Infof("  ⚠️ Failed to record order: %v", err)
		return
	}

	logger.Infof("  ✅ Order recorded as FILLED: %s [%s] %s qty=%.6f price=%.6f", orderID, orderAction, symbol, quantity, exitPrice)

	// Create fill record immediately
	tradeID := fmt.Sprintf("%s-%d", orderID, time.Now().UnixNano())
	fillRecord := &store.TraderFill{
		TraderID:        traderID,
		ExchangeID:      exchangeID,
		ExchangeType:    exchangeType,
		OrderID:         orderRecord.ID,
		ExchangeOrderID: orderID,
		ExchangeTradeID: tradeID,
		Symbol:          symbol,
		Side:            getSideFromAction(orderAction),
		Price:           exitPrice,
		Quantity:        quantity,
		QuoteQuantity:   exitPrice * quantity,
		Commission:      fee,
		CommissionAsset: brokerCurrency(exchangeType),
		RealizedPnL:     0,
		IsMaker:         false,
		CreatedAt:       time.Now().UTC().UnixMilli(),
	}

	if err := s.store.Order().CreateFill(fillRecord); err != nil {
		logger.Infof("  ⚠️ Failed to record fill: %v", err)
	} else {
		logger.Infof("  ✅ Fill record created: price=%.6f qty=%.6f", exitPrice, quantity)
	}
}

// pollAndUpdateOrderStatus Poll order status and update with fill data
func (s *Server) pollAndUpdateOrderStatus(orderRecordID int64, traderID, exchangeID, exchangeType, orderID, symbol, orderAction string, tempTrader trader.Trader) {
	var actualPrice float64
	var actualQty float64
	var fee float64

	// Wait a bit for order to be filled
	time.Sleep(500 * time.Millisecond)

	// Poll GetOrderStatus until the order is filled
	for i := 0; i < 5; i++ {
		status, err := tempTrader.GetOrderStatus(symbol, orderID)
		if err != nil {
			logger.Infof("  ⚠️ GetOrderStatus failed (attempt %d/5): %v", i+1, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if err == nil {
			statusStr, _ := status["status"].(string)
			if statusStr == "FILLED" {
				// Get actual fill price
				if avgPrice, ok := status["avgPrice"].(float64); ok && avgPrice > 0 {
					actualPrice = avgPrice
				}
				// Get actual executed quantity
				if execQty, ok := status["executedQty"].(float64); ok && execQty > 0 {
					actualQty = execQty
				}
				// Get commission/fee
				if commission, ok := status["commission"].(float64); ok {
					fee = commission
				}

				logger.Infof("  ✅ Order filled: avgPrice=%.6f, qty=%.6f, fee=%.6f", actualPrice, actualQty, fee)

				// Update order status to FILLED
				if err := s.store.Order().UpdateOrderStatus(orderRecordID, "FILLED", actualQty, actualPrice, fee); err != nil {
					logger.Infof("  ⚠️ Failed to update order status: %v", err)
					return
				}

				// Record fill details
				tradeID := fmt.Sprintf("%s-%d", orderID, time.Now().UnixNano())
				fillRecord := &store.TraderFill{
					TraderID:        traderID,
					ExchangeID:      exchangeID,
					ExchangeType:    exchangeType,
					OrderID:         orderRecordID,
					ExchangeOrderID: orderID,
					ExchangeTradeID: tradeID,
					Symbol:          symbol,
					Side:            getSideFromAction(orderAction),
					Price:           actualPrice,
					Quantity:        actualQty,
					QuoteQuantity:   actualPrice * actualQty,
					Commission:      fee,
					CommissionAsset: brokerCurrency(exchangeType),
					RealizedPnL:     0,
					IsMaker:         false,
					CreatedAt:       time.Now().UTC().UnixMilli(),
				}

				if err := s.store.Order().CreateFill(fillRecord); err != nil {
					logger.Infof("  ⚠️ Failed to record fill: %v", err)
				} else {
					logger.Infof("  📝 Fill recorded: price=%.6f, qty=%.6f", actualPrice, actualQty)
				}

				return
			} else if statusStr == "CANCELED" || statusStr == "EXPIRED" || statusStr == "REJECTED" {
				logger.Infof("  ⚠️ Order %s, updating status", statusStr)
				s.store.Order().UpdateOrderStatus(orderRecordID, statusStr, 0, 0, 0)
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	logger.Infof("  ⚠️ Failed to confirm order fill after polling, order may still be pending")
}

// getSideFromAction Get order side (BUY/SELL) from order action
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

// brokerCurrency returns the account currency of a broker type, or "" when unknown.
func brokerCurrency(exchangeType string) string {
	caps, err := types.CapabilitiesOf(types.Market(store.MarketForExchangeType(exchangeType)))
	if err != nil {
		return ""
	}
	return caps.Currency
}
