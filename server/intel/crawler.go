// Package intel collects public web pages about a stock, asks an AI model to
// summarise each one, and stores the results for the user's strategies and
// dashboards. Page retrieval is delegated to a Crawl4AI server.
package intel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"crux/security"
)

// Fetcher returns the readable text of a page.
type Fetcher interface {
	Fetch(ctx context.Context, rawURL string) (string, error)
}

// maxCrawlResponseBytes caps what is read from the crawler, which can return
// large markdown bodies.
const maxCrawlResponseBytes = 8 << 20

// Crawl4AI calls the /crawl endpoint of a Crawl4AI server.
type Crawl4AI struct {
	baseURL string
	client  *http.Client
}

// NewCrawl4AI returns a client for the Crawl4AI server at baseURL, for example
// http://crawl4ai:11235.
func NewCrawl4AI(baseURL string) *Crawl4AI {
	return &Crawl4AI{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 2 * time.Minute},
	}
}

type crawlRequest struct {
	URLs []string `json:"urls"`
}

type crawlResponse struct {
	Success bool          `json:"success"`
	Msg     string        `json:"msg"`
	Results []crawlResult `json:"results"`
}

type crawlResult struct {
	URL          string          `json:"url"`
	Success      bool            `json:"success"`
	ErrorMessage string          `json:"error_message"`
	Markdown     json.RawMessage `json:"markdown"`
}

// Fetch asks the crawler for one page and returns its markdown. The target URL
// is checked against the SSRF rules before it is sent.
func (c *Crawl4AI) Fetch(ctx context.Context, rawURL string) (string, error) {
	if err := security.ValidateURL(rawURL); err != nil {
		return "", fmt.Errorf("url rejected: %w", err)
	}
	body, err := json.Marshal(crawlRequest{URLs: []string{rawURL}})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/crawl", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("crawl4ai request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxCrawlResponseBytes))
	if err != nil {
		return "", fmt.Errorf("read crawl4ai response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("crawl4ai returned HTTP %d", resp.StatusCode)
	}

	var parsed crawlResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("parse crawl4ai response: %w", err)
	}
	if !parsed.Success {
		return "", fmt.Errorf("crawl4ai request failed: %s", parsed.Msg)
	}
	if len(parsed.Results) == 0 {
		return "", fmt.Errorf("crawl4ai returned no results")
	}
	result := parsed.Results[0]
	if !result.Success {
		return "", fmt.Errorf("crawl failed: %s", result.ErrorMessage)
	}
	text := markdownText(result.Markdown)
	if strings.TrimSpace(text) == "" {
		return "", fmt.Errorf("page has no readable text")
	}
	return text, nil
}

// markdownText accepts both shapes Crawl4AI uses for markdown: a plain string,
// or an object with raw_markdown and fit_markdown fields.
func markdownText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var obj struct {
		RawMarkdown string `json:"raw_markdown"`
		FitMarkdown string `json:"fit_markdown"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil {
		if obj.RawMarkdown != "" {
			return obj.RawMarkdown
		}
		return obj.FitMarkdown
	}
	return ""
}
