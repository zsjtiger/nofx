package trader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAlpacaPaperTraderRequiresKeys(t *testing.T) {
	_, err := newAlpacaPaperTrader(AutoTraderConfig{Name: "t"})
	assert.ErrorContains(t, err, "alpaca API key and secret")

	tr, err := newAlpacaPaperTrader(AutoTraderConfig{Name: "t", AlpacaAPIKey: "k", AlpacaSecretKey: "s"})
	require.NoError(t, err)
	assert.NotNil(t, tr)
}

func TestNewCNPaperTraderRequiresTokenAndUsesDefaultCash(t *testing.T) {
	_, err := newCNPaperTrader(AutoTraderConfig{Name: "t"})
	assert.ErrorContains(t, err, "tushare token")

	tr, err := newCNPaperTrader(AutoTraderConfig{Name: "t", TushareToken: "tok"})
	require.NoError(t, err)
	bal, err := tr.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, float64(DefaultCNPaperCash), bal["totalWalletBalance"], "no positions means no quote calls")

	tr, err = newCNPaperTrader(AutoTraderConfig{Name: "t", TushareToken: "tok", PaperInitialCash: 250000})
	require.NoError(t, err)
	bal, err = tr.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, 250000.0, bal["totalWalletBalance"])
}

func TestNewAutoTraderRejectsUnknownMarket(t *testing.T) {
	_, err := NewAutoTrader(AutoTraderConfig{Name: "t", Market: "forex", Exchange: "cn_paper"}, nil, "default")
	assert.ErrorContains(t, err, "unknown market")
}
