package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntelRunReturns503WhenCrawlerNotConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	s := &Server{}
	router.POST("/api/intel/run", s.handleIntelRun)

	body := `{"symbol":"600519.SH","urls":["https://example.com"],"model_id":"m1"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/intel/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Contains(t, rec.Body.String(), "CRAWL4AI_URL")
}
