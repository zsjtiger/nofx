package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"nofx/market"
	"nofx/trader/types"

	"github.com/gin-gonic/gin"
)

func runEquityKlines(t *testing.T, query string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/equity/klines?"+query, nil)
	(&Server{}).handleGetEquityKlines(c)
	return w
}

func TestEquityKlinesRejectsMissingMarket(t *testing.T) {
	if w := runEquityKlines(t, "symbol=AAPL"); w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestEquityKlinesRejectsBadLimit(t *testing.T) {
	if w := runEquityKlines(t, "market=us_equity&symbol=AAPL&limit=abc"); w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestEquityKlinesWithoutSourceIsUnavailable(t *testing.T) {
	market.RegisterEquitySource(types.MarketCNEquity, nil)
	if w := runEquityKlines(t, "market=cn_equity&symbol=600519"); w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", w.Code)
	}
}
