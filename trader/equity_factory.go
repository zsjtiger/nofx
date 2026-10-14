package trader

import (
	"fmt"

	"nofx/logger"
	"nofx/market"
	"nofx/provider/alpaca"
	"nofx/provider/tushare"
	"nofx/trader/alpacapaper"
	"nofx/trader/cnpaper"
)

// defaultCNPaperCash is the starting cash for A-share paper accounts, in CNY.
const defaultCNPaperCash = 1_000_000

// newAlpacaPaperTrader builds a US-equity paper trader and registers its daily-bar source.
func newAlpacaPaperTrader(config AutoTraderConfig) (Trader, error) {
	if config.AlpacaAPIKey == "" || config.AlpacaSecretKey == "" {
		return nil, fmt.Errorf("alpaca API key and secret are required for alpaca_paper")
	}
	t, err := alpacapaper.New(config.AlpacaAPIKey, config.AlpacaSecretKey, "", "")
	if err != nil {
		return nil, err
	}
	market.RegisterEquitySource(market.ExchangeAlpacaPaper, market.USDailyKlineSource(alpaca.NewClientWithKeys(config.AlpacaAPIKey, config.AlpacaSecretKey)))
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
		cash = defaultCNPaperCash
	}
	client := tushare.NewClient(config.TushareToken)
	t, err := cnpaper.New(cash, cnpaper.NewTushareQuoter(client))
	if err != nil {
		return nil, err
	}
	market.RegisterEquitySource(market.ExchangeCNPaper, market.CNDailyKlineSource(client))
	logger.Infof("🏦 [%s] Using A-share paper trading (CNY %.0f starting cash)", config.Name, cash)
	return t, nil
}
