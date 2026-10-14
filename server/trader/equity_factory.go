package trader

import (
	"fmt"
	"path/filepath"

	"crux/logger"
	"crux/market"
	"crux/provider/alpaca"
	"crux/provider/tushare"
	"crux/trader/alpacapaper"
	"crux/trader/cnpaper"
	"crux/trader/types"
)

// DefaultCNPaperCash is the starting cash for A-share paper accounts, in CNY.
const DefaultCNPaperCash = 1_000_000

// paperStatePath is where a cn_paper trader keeps its ledger, one file per trader.
func paperStatePath(config AutoTraderConfig) string {
	id := config.ID
	if id == "" {
		id = "default_trader"
	}
	return filepath.Join(PaperStateDir, id+".json")
}

// PaperStateDir holds paper-trading ledgers. It is relative to the working directory.
var PaperStateDir = "data/paper"

// newAlpacaPaperTrader builds a US-equity paper trader and registers its daily-bar source.
func newAlpacaPaperTrader(config AutoTraderConfig) (Trader, error) {
	if config.AlpacaAPIKey == "" || config.AlpacaSecretKey == "" {
		return nil, fmt.Errorf("alpaca API key and secret are required for alpaca_paper")
	}
	t, err := alpacapaper.New(config.AlpacaAPIKey, config.AlpacaSecretKey, "", "")
	if err != nil {
		return nil, err
	}
	market.RegisterEquitySource(types.MarketUSEquity, market.USDailyKlineSource(alpaca.NewClientWithKeys(config.AlpacaAPIKey, config.AlpacaSecretKey)))
	logger.Infof("🏦 [%s] Using Alpaca paper trading (US equities)", config.Name)
	return t, nil
}

// newCNPaperTrader builds an A-share paper broker priced from Tushare and registers its daily-bar source.
func newCNPaperTrader(config AutoTraderConfig) (Trader, error) {
	if config.TushareToken == "" {
		return nil, fmt.Errorf("tushare token is required for cn_paper")
	}
	cash := config.PaperInitialCash
	if cash <= 0 {
		cash = DefaultCNPaperCash
	}
	client := tushare.NewClient(config.TushareToken)
	t, err := cnpaper.NewPersistent(cash, cnpaper.NewTushareQuoter(client), paperStatePath(config))
	if err != nil {
		return nil, err
	}
	market.RegisterEquitySource(types.MarketCNEquity, market.CNDailyKlineSource(client))
	logger.Infof("🏦 [%s] Using A-share paper trading (CNY %.0f starting cash)", config.Name, cash)
	return t, nil
}
