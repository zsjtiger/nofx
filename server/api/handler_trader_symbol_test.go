package api

import (
	"testing"

	"nofx/trader/types"
)

func TestIsSupportedTraderSymbol(t *testing.T) {
	tests := []struct {
		name   string
		market types.Market
		symbol string
		want   bool
	}{
		{name: "US ticker", market: types.MarketUSEquity, symbol: "AAPL", want: true},
		{name: "US ticker lowercase", market: types.MarketUSEquity, symbol: "msft", want: true},
		{name: "US class share", market: types.MarketUSEquity, symbol: "BRK.B", want: true},
		{name: "US rejects A-share code", market: types.MarketUSEquity, symbol: "600519.SH", want: false},
		{name: "A-share with suffix", market: types.MarketCNEquity, symbol: "600519.SH", want: true},
		{name: "A-share bare code", market: types.MarketCNEquity, symbol: "000001", want: true},
		{name: "A-share rejects US ticker", market: types.MarketCNEquity, symbol: "AAPL", want: false},
		{name: "A-share rejects crypto pair", market: types.MarketCNEquity, symbol: "BTCUSDT", want: false},
		{name: "empty slot ignored", market: types.MarketUSEquity, symbol: "  ", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportedTraderSymbol(tt.market, tt.symbol); got != tt.want {
				t.Fatalf("isSupportedTraderSymbol(%s, %q) = %v, want %v", tt.market, tt.symbol, got, tt.want)
			}
		})
	}
}
