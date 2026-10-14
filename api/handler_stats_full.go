package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleStatisticsFull returns the full set of computed performance metrics for
// a single trader: win rate, profit factor, Sharpe ratio, max drawdown, and the
// average win/loss amounts. These are derived from the trader's CLOSED positions
// via store.Position().GetFullStatsByTraderFilters — the same computation the
// strategy engine feeds to the AI, so the dashboard and the model see identical
// numbers.
//
// The existing GET /statistics endpoint only returns cycle/position counts; this
// endpoint exposes the richer trade-quality metrics the terminal dashboard needs.
func (s *Server) handleStatisticsFull(c *gin.Context) {
	_, traderID, err := s.getTraderFromQuery(c)
	if err != nil {
		SafeBadRequest(c, "Invalid trader ID")
		return
	}

	trader, err := s.traderManager.GetTrader(traderID)
	if err != nil {
		SafeNotFound(c, "Trader")
		return
	}

	store := trader.GetStore()
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Store not available"})
		return
	}

	traderIDs := []string{trader.GetID()}
	var traderIDPatterns []string

	stats, err := store.Position().GetFullStatsByTraderFilters(traderIDs, traderIDPatterns, trader.GetInitialBalance())
	if err != nil {
		SafeInternalError(c, "Get full statistics", err)
		return
	}

	c.JSON(http.StatusOK, stats)
}
