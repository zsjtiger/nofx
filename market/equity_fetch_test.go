package market

import (
	"context"
	"errors"
	"testing"

	"nofx/trader/types"
)

func TestFetchEquityKlinesUsesCanonicalSymbolAndClamps(t *testing.T) {
	var gotSymbol string
	var gotLimit int
	RegisterEquitySource(types.MarketUSEquity, func(_ context.Context, symbol string, limit int) ([]Kline, error) {
		gotSymbol, gotLimit = symbol, limit
		return make([]Kline, 500), nil
	})
	t.Cleanup(func() { RegisterEquitySource(types.MarketUSEquity, nil) })

	klines, err := FetchEquityKlines("aapl", types.MarketUSEquity, 9999)
	if err != nil {
		t.Fatal(err)
	}
	if gotSymbol != "AAPL" {
		t.Fatalf("symbol = %q, want AAPL", gotSymbol)
	}
	if gotLimit != equityHistoryDays {
		t.Fatalf("limit = %d, want %d", gotLimit, equityHistoryDays)
	}
	if len(klines) != equityHistoryDays {
		t.Fatalf("returned %d bars, want %d", len(klines), equityHistoryDays)
	}
}

func TestFetchEquityKlinesWithoutSource(t *testing.T) {
	RegisterEquitySource(types.MarketCNEquity, nil)
	_, err := FetchEquityKlines("600519", types.MarketCNEquity, 10)
	if !errors.Is(err, ErrNoEquitySource) {
		t.Fatalf("err = %v, want ErrNoEquitySource", err)
	}
}

func TestFetchEquityKlinesRejectsBadSymbol(t *testing.T) {
	if _, err := FetchEquityKlines("BTCUSDT", types.MarketUSEquity, 10); err == nil {
		t.Fatal("expected error for crypto pair on US market")
	}
}
