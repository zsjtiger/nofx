package intel

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCrawlerServer(t *testing.T, handler http.HandlerFunc) *Crawl4AI {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewCrawl4AI(srv.URL)
}

func TestCrawl4AIFetchSendsURLAndReadsStringMarkdown(t *testing.T) {
	var gotPath string
	var gotBody crawlRequest
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(raw, &gotBody))
		_, _ = w.Write([]byte(`{"success":true,"results":[{"url":"https://example.com/a","success":true,"markdown":"贵州茅台发布年报"}]}`))
	})

	text, err := c.Fetch(context.Background(), "https://example.com/a")
	require.NoError(t, err)
	assert.Equal(t, "贵州茅台发布年报", text)
	assert.Equal(t, "/crawl", gotPath)
	assert.Equal(t, []string{"https://example.com/a"}, gotBody.URLs)
}

func TestCrawl4AIFetchReadsObjectMarkdown(t *testing.T) {
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"results":[{"success":true,"markdown":{"raw_markdown":"raw text","fit_markdown":"fit"}}]}`))
	})
	text, err := c.Fetch(context.Background(), "https://example.com/a")
	require.NoError(t, err)
	assert.Equal(t, "raw text", text)
}

func TestCrawl4AIFetchReportsCrawlFailure(t *testing.T) {
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"results":[{"success":false,"error_message":"timeout"}]}`))
	})
	_, err := c.Fetch(context.Background(), "https://example.com/a")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestCrawl4AIFetchRejectsEmptyPage(t *testing.T) {
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":true,"results":[{"success":true,"markdown":"   "}]}`))
	})
	_, err := c.Fetch(context.Background(), "https://example.com/a")
	require.Error(t, err)
}

func TestCrawl4AIFetchHandlesHTTPError(t *testing.T) {
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	_, err := c.Fetch(context.Background(), "https://example.com/a")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestCrawl4AIFetchBlocksInternalTargets(t *testing.T) {
	called := false
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	_, err := c.Fetch(context.Background(), "http://127.0.0.1/admin")
	require.Error(t, err)
	assert.False(t, called, "internal target must be rejected before reaching the crawler")
}

// TestCrawl4AIFetchParsesRecordedServerResponse runs the client against a
// /crawl response recorded from Crawl4AI 0.9.5 (testdata/crawl4ai_success.json).
func TestCrawl4AIFetchParsesRecordedServerResponse(t *testing.T) {
	recorded, err := os.ReadFile("testdata/crawl4ai_success.json")
	require.NoError(t, err)
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(recorded)
	})
	text, err := c.Fetch(context.Background(), "https://example.com/茅台半年报")
	require.NoError(t, err)
	assert.Contains(t, text, "贵州茅台 2026 年半年报摘要")
}

func TestCrawl4AIFetchReportsTopLevelFailure(t *testing.T) {
	c := newCrawlerServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":false,"msg":"browser unavailable","results":[]}`))
	})
	_, err := c.Fetch(context.Background(), "https://example.com/a")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "browser unavailable")
}
