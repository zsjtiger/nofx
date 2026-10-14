package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMarket(t *testing.T) {
	m, err := ParseMarket(" CN_Equity ")
	require.NoError(t, err)
	assert.Equal(t, MarketCNEquity, m)

	_, err = ParseMarket("forex")
	assert.Error(t, err)
}

func TestCapabilitiesOf(t *testing.T) {
	cn, err := CapabilitiesOf(MarketCNEquity)
	require.NoError(t, err)
	assert.Equal(t, "CNY", cn.Currency)
	assert.Equal(t, "Asia/Shanghai", cn.Timezone)
	assert.False(t, cn.AllowShort)
	assert.False(t, cn.AllowLeverage)
	assert.Equal(t, 100, cn.LotSize)
	assert.True(t, cn.TPlusOne)

	us, err := CapabilitiesOf(MarketUSEquity)
	require.NoError(t, err)
	assert.Equal(t, "USD", us.Currency)
	assert.Equal(t, "America/New_York", us.Timezone)
	assert.True(t, us.AllowShort)
	assert.Equal(t, 0, us.LotSize)
	assert.False(t, us.TPlusOne)

	_, err = CapabilitiesOf("bogus")
	assert.Error(t, err)
}

func TestNormalizeSymbol(t *testing.T) {
	tests := []struct {
		name    string
		market  Market
		input   string
		want    string
		wantErr bool
	}{
		{name: "us lower case", market: MarketUSEquity, input: "aapl", want: "AAPL"},
		{name: "us class share", market: MarketUSEquity, input: "brk.b", want: "BRK.B"},
		{name: "us rejects digits", market: MarketUSEquity, input: "AAP1", wantErr: true},
		{name: "us rejects too long", market: MarketUSEquity, input: "ABCDEF", wantErr: true},
		{name: "cn shanghai inferred", market: MarketCNEquity, input: "600519", want: "600519.SH"},
		{name: "cn star market inferred", market: MarketCNEquity, input: "688981", want: "688981.SH"},
		{name: "cn shenzhen inferred", market: MarketCNEquity, input: "000001", want: "000001.SZ"},
		{name: "cn chinext inferred", market: MarketCNEquity, input: "300750", want: "300750.SZ"},
		{name: "cn beijing inferred", market: MarketCNEquity, input: "830799", want: "830799.BJ"},
		{name: "cn explicit suffix kept", market: MarketCNEquity, input: "000001.sz", want: "000001.SZ"},
		{name: "cn bad suffix", market: MarketCNEquity, input: "600519.HK", wantErr: true},
		{name: "cn bad length", market: MarketCNEquity, input: "60051", wantErr: true},
		{name: "cn unknown prefix", market: MarketCNEquity, input: "100001", wantErr: true},
		{name: "crypto passthrough", market: MarketCrypto, input: "btcusdt", want: "BTCUSDT"},
		{name: "empty", market: MarketUSEquity, input: "  ", wantErr: true},
		{name: "unknown market", market: "bogus", input: "AAPL", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeSymbol(tt.market, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateOrderQuantity(t *testing.T) {
	cn, _ := CapabilitiesOf(MarketCNEquity)
	us, _ := CapabilitiesOf(MarketUSEquity)

	tests := []struct {
		name           string
		caps           Capabilities
		qty            float64
		isBuy          bool
		closesPosition bool
		wantErr        bool
	}{
		{name: "cn buy whole lot", caps: cn, qty: 300, isBuy: true},
		{name: "cn buy odd lot rejected", caps: cn, qty: 150, isBuy: true, wantErr: true},
		{name: "cn sell odd lot partial rejected", caps: cn, qty: 150, isBuy: false, wantErr: true},
		{name: "cn sell odd lot closing position", caps: cn, qty: 150, isBuy: false, closesPosition: true},
		{name: "cn fractional rejected", caps: cn, qty: 100.5, isBuy: true, wantErr: true},
		{name: "cn zero rejected", caps: cn, qty: 0, isBuy: true, wantErr: true},
		{name: "us fractional shares allowed", caps: us, qty: 1.5, isBuy: true},
		{name: "us whole share", caps: us, qty: 7, isBuy: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderQuantity(tt.caps, tt.qty, tt.isBuy, tt.closesPosition)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSideAllowed(t *testing.T) {
	cn, _ := CapabilitiesOf(MarketCNEquity)
	us, _ := CapabilitiesOf(MarketUSEquity)

	assert.Error(t, ValidateSideAllowed(cn, 1, true), "A-shares cannot short")
	assert.Error(t, ValidateSideAllowed(cn, 2, false), "A-shares cannot use leverage")
	assert.NoError(t, ValidateSideAllowed(cn, 1, false))
	assert.NoError(t, ValidateSideAllowed(us, 1, true), "US equities can short")
	assert.Error(t, ValidateSideAllowed(us, 2, false), "US equities here are unlevered")
}

func TestResolveMarket(t *testing.T) {
	m, err := ResolveMarket("")
	require.NoError(t, err)
	assert.Equal(t, MarketCrypto, m, "empty market keeps legacy crypto behaviour")

	m, err = ResolveMarket("cn_equity")
	require.NoError(t, err)
	assert.Equal(t, MarketCNEquity, m)

	_, err = ResolveMarket("forex")
	assert.Error(t, err)
}

func TestFloorToLot(t *testing.T) {
	cn, _ := CapabilitiesOf(MarketCNEquity)
	us, _ := CapabilitiesOf(MarketUSEquity)

	assert.Equal(t, 300.0, FloorToLot(cn, 350))
	assert.Equal(t, 0.0, FloorToLot(cn, 99))
	assert.Equal(t, 1.5, FloorToLot(us, 1.5), "markets without lots are unchanged")
	assert.Equal(t, -5.0, FloorToLot(cn, -5), "non-positive input is left for ValidateOrderQuantity to reject")
}
