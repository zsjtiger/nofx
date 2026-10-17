package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crux/intel"

	"github.com/gin-gonic/gin"
)

// intelRunTimeout bounds one analysis run, which fetches and summarises several pages.
const intelRunTimeout = 5 * time.Minute

// WithIntelFetcher enables intelligence analysis with the given page fetcher.
// Without it, the intel run endpoint returns 503.
func (s *Server) WithIntelFetcher(f intel.Fetcher) *Server {
	s.intelFetcher = f
	return s
}

// handleIntelList returns the newest analysed pages for a symbol.
func (s *Server) handleIntelList(c *gin.Context) {
	userID := c.GetString("user_id")
	symbol := strings.TrimSpace(c.Query("symbol"))
	if symbol == "" {
		SafeBadRequest(c, "symbol is required")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, err := s.store.Intel().ListBySymbol(userID, strings.ToUpper(symbol), limit)
	if err != nil {
		SafeInternalError(c, "List intel items", err)
		return
	}
	c.JSON(http.StatusOK, items)
}

type intelRunRequest struct {
	Symbol  string   `json:"symbol" binding:"required"`
	URLs    []string `json:"urls" binding:"required"`
	ModelID string   `json:"model_id" binding:"required"`
}

// handleIntelRun fetches each URL through the crawler, asks the chosen AI model
// to analyse it, and stores the results.
func (s *Server) handleIntelRun(c *gin.Context) {
	if s.intelFetcher == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "intelligence crawler is not configured (set CRAWL4AI_URL)"})
		return
	}
	var req intelRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SafeBadRequest(c, "Invalid request parameters")
		return
	}
	userID := c.GetString("user_id")

	complete := func(systemPrompt, userPrompt string) (string, error) {
		return s.runRealAITest(userID, req.ModelID, systemPrompt, userPrompt)
	}
	service := intel.NewService(s.intelFetcher, complete, s.store.Intel())

	ctx, cancel := context.WithTimeout(c.Request.Context(), intelRunTimeout)
	defer cancel()
	results, err := service.Run(ctx, intel.Job{UserID: userID, Symbol: req.Symbol, URLs: req.URLs})
	if err != nil {
		SafeBadRequest(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}
