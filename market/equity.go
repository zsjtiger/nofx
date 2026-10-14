package market

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nofx/provider/alpaca"
	"nofx/provider/tushare"
	"nofx/trader/types"
)

// Exchange names for equity accounts. They match the exchange_type values in the store.
const (
	ExchangeCNPaper     = "cn_paper"
	ExchangeAlpacaPaper = "alpaca_paper"
)

// minEquityBars is the history needed for the indicators in BuildDataFromKlines.
const minEquityBars = 30

// equityHistoryDays is the calendar window requested for daily bars.
const equityHistoryDays = 200

// EquityKlineSource returns daily klines for a canonical symbol, oldest first.
type EquityKlineSource func(ctx context.Context, symbol string, limit int) ([]Kline, error)

var (
	equityMu      sync.RWMutex
	equitySources = map[string]EquityKlineSource{}
)

// IsEquityExchange reports whether an exchange type trades stocks rather than perpetuals.
func IsEquityExchange(exchange string) bool {
	switch exchange {
	case ExchangeCNPaper, ExchangeAlpacaPaper:
		return true
	default:
		return false
	}
}

// RegisterEquitySource sets the daily-bar source for an equity exchange type.
// Each exchange type has one source per process; a trader registers its own on creation.
func RegisterEquitySource(exchange string, src EquityKlineSource) {
	equityMu.Lock()
	defer equityMu.Unlock()
	equitySources[exchange] = src
}

func equitySource(exchange string) EquityKlineSource {
	equityMu.RLock()
	defer equityMu.RUnlock()
	return equitySources[exchange]
}

// equityMarket maps an equity exchange type to its market.
func equityMarket(exchange string) (types.Market, error) {
	switch exchange {
	case ExchangeCNPaper:
		return types.MarketCNEquity, nil
	case ExchangeAlpacaPaper:
		return types.MarketUSEquity, nil
	default:
		return "", fmt.Errorf("%s is not an equity exchange", exchange)
	}
}

// getEquityData builds market data for a stock from daily bars.
// Intraday bars are not available from these sources, so the "primary" series is daily.
func getEquityData(symbol, exchange string) (*Data, error) {
	m, err := equityMarket(exchange)
	if err != nil {
		return nil, err
	}
	canonical, err := types.NormalizeSymbol(m, symbol)
	if err != nil {
		return nil, err
	}
	src := equitySource(exchange)
	if src == nil {
		return nil, fmt.Errorf("no daily bar source registered for %s", exchange)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	klines, err := src(ctx, canonical, equityHistoryDays)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily bars for %s: %w", canonical, err)
	}
	if len(klines) < minEquityBars {
		return nil, fmt.Errorf("only %d daily bars for %s, need %d", len(klines), canonical, minEquityBars)
	}
	data, err := BuildDataFromKlines(canonical, klines, klines)
	if err != nil {
		return nil, err
	}
	// BuildDataFromKlines applies crypto symbol rules, so restore the equity symbol.
	data.Symbol = canonical
	return data, nil
}

// CNDailyKlineSource reads forward-adjusted daily bars for A-shares from Tushare.
func CNDailyKlineSource(client *tushare.Client) EquityKlineSource {
	return func(ctx context.Context, symbol string, limit int) ([]Kline, error) {
		end := time.Now()
		start := end.AddDate(0, 0, -limit)
		startStr, endStr := start.Format("20060102"), end.Format("20060102")

		bars, err := client.DailyBars(ctx, symbol, startStr, endStr)
		if err != nil {
			return nil, err
		}
		factors, err := client.AdjFactors(ctx, symbol, startStr, endStr)
		if err != nil {
			return nil, err
		}
		adjusted, err := tushare.ForwardAdjust(bars, factors)
		if err != nil {
			return nil, err
		}
		return cnBarsToKlines(adjusted)
	}
}

func cnBarsToKlines(bars []tushare.Bar) ([]Kline, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	out := make([]Kline, 0, len(bars))
	for _, b := range bars {
		// A-share daily bars close at 15:00 Beijing time.
		closeAt, err := time.ParseInLocation("20060102 15:04", b.TradeDate+" 15:00", loc)
		if err != nil {
			return nil, fmt.Errorf("bad trade date %q: %w", b.TradeDate, err)
		}
		out = append(out, Kline{
			OpenTime:  closeAt.Add(-6 * time.Hour).UnixMilli(),
			CloseTime: closeAt.UnixMilli(),
			Open:      b.Open,
			High:      b.High,
			Low:       b.Low,
			Close:     b.Close,
			Volume:    b.Vol * 100, // Tushare reports lots of 100 shares
		})
	}
	return out, nil
}

// USDailyKlineSource reads daily bars for US stocks from Alpaca market data.
func USDailyKlineSource(client *alpaca.Client) EquityKlineSource {
	return func(ctx context.Context, symbol string, limit int) ([]Kline, error) {
		bars, err := client.GetBars(ctx, symbol, "1Day", limit)
		if err != nil {
			return nil, err
		}
		out := make([]Kline, 0, len(bars))
		for _, b := range bars {
			out = append(out, Kline{
				OpenTime:  b.Timestamp.UnixMilli(),
				CloseTime: b.Timestamp.Add(24 * time.Hour).UnixMilli(),
				Open:      b.Open,
				High:      b.High,
				Low:       b.Low,
				Close:     b.Close,
				Volume:    float64(b.Volume),
			})
		}
		return out, nil
	}
}

// NormalizeFor canonicalises a symbol for the given exchange type. Equity
// exchanges use their own symbol rules; every other exchange keeps the
// crypto perpetual rules in Normalize.
func NormalizeFor(exchange, symbol string) string {
	m, err := equityMarket(exchange)
	if err != nil {
		return Normalize(symbol)
	}
	canonical, err := types.NormalizeSymbol(m, symbol)
	if err != nil {
		return symbol
	}
	return canonical
}
