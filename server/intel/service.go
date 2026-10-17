package intel

import (
	"context"
	"fmt"
	"strings"

	"crux/store"
)

// MaxURLsPerRun bounds how many pages one analysis run may fetch.
const MaxURLsPerRun = 5

// Service fetches pages, analyses them, and stores the results.
type Service struct {
	fetcher  Fetcher
	complete Completer
	store    *store.IntelStore
}

// NewService wires a page fetcher, an AI completer, and the intel store.
func NewService(fetcher Fetcher, complete Completer, st *store.IntelStore) *Service {
	return &Service{fetcher: fetcher, complete: complete, store: st}
}

// Job is one analysis request for a single stock.
type Job struct {
	UserID string
	Symbol string
	URLs   []string
}

// Result reports the outcome for one URL. Item is nil when the URL failed.
type Result struct {
	URL   string           `json:"url"`
	Item  *store.IntelItem `json:"item,omitempty"`
	Error string           `json:"error,omitempty"`
}

// Run analyses each URL in order. A failing URL is reported in its Result and
// does not stop the others. Run returns early only when ctx is cancelled.
func (s *Service) Run(ctx context.Context, job Job) ([]Result, error) {
	symbol := strings.ToUpper(strings.TrimSpace(job.Symbol))
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	if len(job.URLs) == 0 {
		return nil, fmt.Errorf("at least one url is required")
	}
	if len(job.URLs) > MaxURLsPerRun {
		return nil, fmt.Errorf("at most %d urls per run", MaxURLsPerRun)
	}

	results := make([]Result, 0, len(job.URLs))
	for _, url := range job.URLs {
		if err := ctx.Err(); err != nil {
			return results, err
		}
		results = append(results, s.analyzeOne(ctx, job.UserID, symbol, url))
	}
	return results, nil
}

func (s *Service) analyzeOne(ctx context.Context, userID, symbol, url string) Result {
	res := Result{URL: url}

	text, err := s.fetcher.Fetch(ctx, url)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	finding, err := Analyze(s.complete, symbol, text)
	if err != nil {
		res.Error = err.Error()
		return res
	}

	item := &store.IntelItem{
		UserID:    userID,
		Symbol:    symbol,
		SourceURL: url,
		Summary:   finding.Summary,
		Sentiment: finding.Sentiment,
		Risk:      finding.Risk,
		Symbols:   strings.Join(verifiedSymbols(finding.Symbols, text), ","),
	}
	if err := s.store.Create(item); err != nil {
		res.Error = fmt.Sprintf("save result: %v", err)
		return res
	}
	res.Item = item
	return res
}

// verifiedSymbols keeps the reported symbols whose code appears in the page text,
// which filters out codes the model made up. The symbol itself is not included.
func verifiedSymbols(reported []string, pageText string) []string {
	var out []string
	for _, sym := range reported {
		if sym == "" || !strings.Contains(pageText, symbolCode(sym)) {
			continue
		}
		out = append(out, sym)
	}
	return out
}

// symbolCode returns the digits of an A-share code (600519.SH -> 600519) or the
// ticker itself (AAPL -> AAPL), which is what usually appears in article text.
func symbolCode(sym string) string {
	if i := strings.IndexByte(sym, '.'); i > 0 {
		return sym[:i]
	}
	return sym
}
