package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustCaps(t *testing.T, m Market) Capabilities {
	t.Helper()
	c, err := CapabilitiesOf(m)
	require.NoError(t, err)
	return c
}

func TestInTradingSessionCN(t *testing.T) {
	cn := mustCaps(t, MarketCNEquity)
	// 2026-10-14 is a Wednesday. Beijing time is UTC+8.
	at := func(hhmm string) time.Time {
		ts, err := time.Parse("2006-01-02 15:04", "2026-10-14 "+hhmm)
		require.NoError(t, err)
		return time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), 0, 0, time.FixedZone("CST", 8*3600))
	}
	tests := []struct {
		name string
		at   time.Time
		want bool
	}{
		{name: "before open", at: at("09:29"), want: false},
		{name: "open", at: at("09:30"), want: true},
		{name: "morning close excluded", at: at("11:30"), want: false},
		{name: "lunch break", at: at("12:00"), want: false},
		{name: "afternoon", at: at("14:59"), want: true},
		{name: "after close", at: at("15:00"), want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cn.InTradingSession(tt.at)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	sat := time.Date(2026, 10, 17, 10, 0, 0, 0, time.FixedZone("CST", 8*3600))
	got, err := cn.InTradingSession(sat)
	require.NoError(t, err)
	assert.False(t, got, "weekends are closed")
}

func TestInTradingSessionUSAndCrypto(t *testing.T) {
	us := mustCaps(t, MarketUSEquity)
	// 2026-10-14 10:00 New York time (EDT, UTC-4) is inside the regular session.
	open := time.Date(2026, 10, 14, 14, 0, 0, 0, time.UTC)
	got, err := us.InTradingSession(open)
	require.NoError(t, err)
	assert.True(t, got)

	closed := time.Date(2026, 10, 14, 21, 0, 0, 0, time.UTC)
	got, err = us.InTradingSession(closed)
	require.NoError(t, err)
	assert.False(t, got)

	crypto := mustCaps(t, MarketCrypto)
	got, err = crypto.InTradingSession(closed)
	require.NoError(t, err)
	assert.True(t, got, "crypto has no session table and is always open")
}

func TestSameTradingDayUsesMarketTimezone(t *testing.T) {
	cn := mustCaps(t, MarketCNEquity)
	// 15:00 UTC on 14 Oct is 23:00 on 14 Oct in Beijing; 17:00 UTC is 01:00 on 15 Oct.
	a := time.Date(2026, 10, 14, 15, 0, 0, 0, time.UTC)
	b := time.Date(2026, 10, 14, 17, 0, 0, 0, time.UTC)
	same, err := cn.SameTradingDay(a, b)
	require.NoError(t, err)
	assert.False(t, same, "Beijing dates differ even though UTC dates are the same")
}

func TestValidateSettlementT1(t *testing.T) {
	cn := mustCaps(t, MarketCNEquity)
	us := mustCaps(t, MarketUSEquity)
	bought := time.Date(2026, 10, 14, 2, 0, 0, 0, time.UTC) // 10:00 Beijing
	sameDay := time.Date(2026, 10, 14, 6, 30, 0, 0, time.UTC)
	nextDay := time.Date(2026, 10, 15, 2, 0, 0, 0, time.UTC)

	assert.Error(t, ValidateSettlement(cn, bought, sameDay), "same-day sell blocked")
	assert.NoError(t, ValidateSettlement(cn, bought, nextDay), "next-day sell allowed")
	assert.NoError(t, ValidateSettlement(us, bought, sameDay), "US equities are not T+1 locked")
	assert.NoError(t, ValidateSettlement(cn, time.Time{}, sameDay), "unknown entry time is not blocked here")
}

func TestPriceLimits(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		prevClose float64
		isST      bool
		wantLow   float64
		wantHigh  float64
	}{
		{name: "main board", symbol: "600519.SH", prevClose: 100, wantLow: 90, wantHigh: 110},
		{name: "chinext", symbol: "300750.SZ", prevClose: 100, wantLow: 80, wantHigh: 120},
		{name: "star", symbol: "688981.SH", prevClose: 100, wantLow: 80, wantHigh: 120},
		{name: "beijing", symbol: "830799.BJ", prevClose: 100, wantLow: 70, wantHigh: 130},
		{name: "st", symbol: "600000.SH", prevClose: 100, isST: true, wantLow: 95, wantHigh: 105},
		{name: "rounds to cents", symbol: "000001.SZ", prevClose: 10.33, wantLow: 9.30, wantHigh: 11.36},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			low, high, err := PriceLimits(tt.symbol, tt.prevClose, tt.isST)
			require.NoError(t, err)
			assert.InDelta(t, tt.wantLow, low, 0.001)
			assert.InDelta(t, tt.wantHigh, high, 0.001)
		})
	}

	_, _, err := PriceLimits("600519.SH", 0, false)
	assert.Error(t, err)
}

func TestValidatePriceWithinLimits(t *testing.T) {
	assert.NoError(t, ValidatePriceWithinLimits("600519.SH", 100, 110, false), "upper limit is allowed")
	assert.Error(t, ValidatePriceWithinLimits("600519.SH", 100, 110.01, false), "above limit-up rejected")
	assert.Error(t, ValidatePriceWithinLimits("600519.SH", 100, 89.99, false), "below limit-down rejected")
}
