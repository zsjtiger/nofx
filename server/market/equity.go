package market

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"crux/provider/alpaca"
	"crux/provider/tushare"
	"crux/trader/types"
)

// minEquityBars is the history needed for the indicators in BuildDataFromKlines.
const minEquityBars = 30

// equityHistoryDays is the calendar window requested for daily bars.
const equityHistoryDays = 200

// EquityKlineSource returns daily klines for a canonical symbol, oldest first.
type EquityKlineSource func(ctx context.Context, symbol string, limit int) ([]Kline, error)

var (
	equityMu      sync.RWMutex
	equitySources = map[types.Market]EquityKlineSource{}
)

// RegisterEquitySource sets the daily-bar source for a stock market. A trader
// registers the source for its market when it is created.
func RegisterEquitySource(m types.Market, src EquityKlineSource) {
	equityMu.Lock()
	defer equityMu.Unlock()
	equitySources[m] = src
}

func equitySource(m types.Market) EquityKlineSource {
	equityMu.RLock()
	defer equityMu.RUnlock()
	return equitySources[m]
}

// GetForMarket returns market data for a stock in the given market. Data is built
// from daily bars, because the equity sources do not provide intraday bars, so
// the "primary" series is daily.
func GetForMarket(symbol string, m types.Market) (*Data, error) {
	canonical, err := types.NormalizeSymbol(m, symbol)
	if err != nil {
		return nil, err
	}
	src := equitySource(m)
	if src == nil {
		return nil, fmt.Errorf("no daily bar source registered for %s", m)
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
	// BuildDataFromKlines may rewrite the symbol, so restore the canonical equity symbol.
	data.Symbol = canonical
	return data, nil
}

// ErrNoEquitySource is returned when no daily-bar source is registered for a market,
// which happens until a trader for that market has been created.
var ErrNoEquitySource = errors.New("no daily bar source registered")

// FetchEquityKlines returns up to limit daily bars for a stock in market m, oldest
// first. The limit is clamped to the equity history window.
func FetchEquityKlines(symbol string, m types.Market, limit int) ([]Kline, error) {
	canonical, err := types.NormalizeSymbol(m, symbol)
	if err != nil {
		return nil, err
	}
	src := equitySource(m)
	if src == nil {
		return nil, fmt.Errorf("%w for %s", ErrNoEquitySource, m)
	}
	if limit <= 0 || limit > equityHistoryDays {
		limit = equityHistoryDays
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	klines, err := src(ctx, canonical, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily bars for %s: %w", canonical, err)
	}
	if len(klines) > limit {
		klines = klines[len(klines)-limit:]
	}
	return klines, nil
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

// NormalizeFor canonicalises a symbol for the given market. Input that is not
// a valid symbol is upper-cased and trimmed, so it still matches exchange data.
func NormalizeFor(m types.Market, symbol string) string {
	canonical, err := types.NormalizeSymbol(m, symbol)
	if err != nil {
		return strings.ToUpper(strings.TrimSpace(symbol))
	}
	return canonical
}
