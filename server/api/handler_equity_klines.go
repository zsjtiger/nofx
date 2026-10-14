package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"crux/market"
	"crux/trader/types"

	"github.com/gin-gonic/gin"
)

// handleGetEquityKlines returns daily bars for a stock: GET /equity/klines?market=&symbol=&limit=.
func (s *Server) handleGetEquityKlines(c *gin.Context) {
	m, err := types.ResolveMarket(c.Query("market"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit := 120
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
			return
		}
		limit = n
	}
	klines, err := market.FetchEquityKlines(c.Query("symbol"), m, limit)
	if err != nil {
		if errors.Is(err, market.ErrNoEquitySource) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no market data source for this market yet; create a trader for it first"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	symbol := market.NormalizeFor(m, c.Query("symbol"))
	c.JSON(http.StatusOK, gin.H{"market": m, "symbol": symbol, "klines": klines})
}
