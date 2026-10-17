package trader

import (
	"testing"
	"time"

	"crux/store"
	"crux/trader/types"

	"github.com/stretchr/testify/assert"
)

func intelItem(symbol, mentions string, sentiment float64, risk string, hoursAgo float64) store.IntelItem {
	return store.IntelItem{
		Symbol:    symbol,
		Symbols:   mentions,
		Sentiment: sentiment,
		Risk:      risk,
		CreatedAt: time.Now().Add(-time.Duration(hoursAgo * float64(time.Hour))),
	}
}

func TestRankDiscoveredSymbolsRanksMentionsFromPositivePages(t *testing.T) {
	since := time.Now().Add(-72 * time.Hour)
	items := []store.IntelItem{
		intelItem("600519.SH", "000858.SZ,601318.SH", 0.8, store.IntelRiskNone, 1),
		intelItem("000001.SZ", "000858.SZ", 0.4, store.IntelRiskLow, 2),
		intelItem("300750.SZ", "601318.SH", -0.5, store.IntelRiskNone, 3), // negative page: its mentions drop
	}
	got := rankDiscoveredSymbols(items, types.MarketCNEquity, 0, map[string]bool{}, since)
	// 600519 avg 0.8; 000858 avg (0.8+0.4)/2=0.6; 000001 avg 0.4; 601318 avg (0.8-0.5)/2=0.15
	assert.Equal(t, []string{"600519.SH", "000858.SZ", "000001.SZ", "601318.SH"}, got)
}

func TestRankDiscoveredSymbolsSkipsHighRiskExcludedAndOldPages(t *testing.T) {
	since := time.Now().Add(-72 * time.Hour)
	items := []store.IntelItem{
		intelItem("600519.SH", "000858.SZ", 0.9, store.IntelRiskHigh, 1),   // high risk page
		intelItem("601318.SH", "300750.SZ", 0.9, store.IntelRiskNone, 100), // outside window
		intelItem("000001.SZ", "600036.SH", 0.5, store.IntelRiskNone, 1),
	}
	exclude := map[string]bool{"600036.SH": true}
	got := rankDiscoveredSymbols(items, types.MarketCNEquity, 0, exclude, since)
	assert.Equal(t, []string{"000001.SZ"}, got)
}

func TestRankDiscoveredSymbolsHonoursMinSentimentAndNormalizesCase(t *testing.T) {
	since := time.Now().Add(-72 * time.Hour)
	items := []store.IntelItem{
		intelItem("aapl", "msft", 0.3, store.IntelRiskNone, 1),
		intelItem("nvda", "", 0.1, store.IntelRiskNone, 1),
	}
	got := rankDiscoveredSymbols(items, types.MarketUSEquity, 0.2, map[string]bool{}, since)
	assert.Equal(t, []string{"AAPL", "MSFT"}, got)
}

func TestRankDiscoveredSymbolsDropsSymbolsForOtherMarket(t *testing.T) {
	since := time.Now().Add(-72 * time.Hour)
	items := []store.IntelItem{intelItem("600519.SH", "AAPL", 0.6, store.IntelRiskNone, 1)}
	got := rankDiscoveredSymbols(items, types.MarketCNEquity, 0, map[string]bool{}, since)
	assert.Equal(t, []string{"600519.SH"}, got, "a US ticker is not valid on the A-share market")
}
