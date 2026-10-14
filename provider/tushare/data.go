package tushare

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// Bar is one daily bar. Vol is in lots of 100 shares, as Tushare reports it.
type Bar struct {
	TradeDate string  // YYYYMMDD
	Open      float64 // raw (unadjusted) prices
	High      float64
	Low       float64
	Close     float64
	PreClose  float64
	Vol       float64
	Amount    float64 // thousand CNY
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case nil:
		return 0, false
	default:
		return 0, false
	}
}

func getFloat(row map[string]any, key string) float64 {
	f, _ := toFloat(row[key])
	return f
}

func getString(row map[string]any, key string) string {
	s, _ := row[key].(string)
	return s
}

// DailyBars returns unadjusted daily bars between start and end (YYYYMMDD), oldest first.
func (c *Client) DailyBars(ctx context.Context, tsCode, start, end string) ([]Bar, error) {
	table, err := c.Call(ctx, "daily", map[string]string{
		"ts_code":    tsCode,
		"start_date": start,
		"end_date":   end,
	}, "trade_date,open,high,low,close,pre_close,vol,amount")
	if err != nil {
		return nil, err
	}
	bars := make([]Bar, 0, len(table.Items))
	for _, row := range table.Rows() {
		bars = append(bars, Bar{
			TradeDate: getString(row, "trade_date"),
			Open:      getFloat(row, "open"),
			High:      getFloat(row, "high"),
			Low:       getFloat(row, "low"),
			Close:     getFloat(row, "close"),
			PreClose:  getFloat(row, "pre_close"),
			Vol:       getFloat(row, "vol"),
			Amount:    getFloat(row, "amount"),
		})
	}
	sort.Slice(bars, func(i, j int) bool { return bars[i].TradeDate < bars[j].TradeDate })
	return bars, nil
}

// AdjFactors returns the adjustment factor for each trade date.
func (c *Client) AdjFactors(ctx context.Context, tsCode, start, end string) (map[string]float64, error) {
	table, err := c.Call(ctx, "adj_factor", map[string]string{
		"ts_code":    tsCode,
		"start_date": start,
		"end_date":   end,
	}, "trade_date,adj_factor")
	if err != nil {
		return nil, err
	}
	out := make(map[string]float64, len(table.Items))
	for _, row := range table.Rows() {
		out[getString(row, "trade_date")] = getFloat(row, "adj_factor")
	}
	return out, nil
}

// ForwardAdjust applies forward (前复权) adjustment: each price is scaled by
// factor(date) / factor(latest date), so the most recent prices stay unchanged.
// Bars without a factor are returned unadjusted and reported in the error.
func ForwardAdjust(bars []Bar, factors map[string]float64) ([]Bar, error) {
	if len(bars) == 0 {
		return nil, nil
	}
	latest := bars[len(bars)-1].TradeDate
	latestFactor, ok := factors[latest]
	if !ok || latestFactor == 0 {
		return nil, fmt.Errorf("missing adjustment factor for %s", latest)
	}
	out := make([]Bar, len(bars))
	var missing []string
	for i, b := range bars {
		f, ok := factors[b.TradeDate]
		if !ok || f == 0 {
			missing = append(missing, b.TradeDate)
			out[i] = b
			continue
		}
		ratio := f / latestFactor
		out[i] = b
		out[i].Open = round2(b.Open * ratio)
		out[i].High = round2(b.High * ratio)
		out[i].Low = round2(b.Low * ratio)
		out[i].Close = round2(b.Close * ratio)
		out[i].PreClose = round2(b.PreClose * ratio)
	}
	if len(missing) > 0 {
		return out, fmt.Errorf("missing adjustment factor for %s", strings.Join(missing, ","))
	}
	return out, nil
}

func round2(v float64) float64 {
	return float64(int64(v*100+signHalf(v))) / 100
}

func signHalf(v float64) float64 {
	if v < 0 {
		return -0.5
	}
	return 0.5
}

// StockLimit returns the limit-up and limit-down prices for a trade date.
// Tushare reports these per day, which covers ST names and board differences.
func (c *Client) StockLimit(ctx context.Context, tsCode, tradeDate string) (up, down float64, err error) {
	table, err := c.Call(ctx, "stk_limit", map[string]string{
		"ts_code":    tsCode,
		"trade_date": tradeDate,
	}, "trade_date,up_limit,down_limit")
	if err != nil {
		return 0, 0, err
	}
	rows := table.Rows()
	if len(rows) == 0 {
		return 0, 0, fmt.Errorf("no limit prices for %s on %s", tsCode, tradeDate)
	}
	return getFloat(rows[0], "up_limit"), getFloat(rows[0], "down_limit"), nil
}

// IsST reports whether the stock is marked ST (special treatment) by its short name.
func (c *Client) IsST(ctx context.Context, tsCode string) (bool, error) {
	table, err := c.Call(ctx, "stock_basic", map[string]string{"ts_code": tsCode}, "ts_code,name")
	if err != nil {
		return false, err
	}
	rows := table.Rows()
	if len(rows) == 0 {
		return false, fmt.Errorf("unknown ts_code %s", tsCode)
	}
	return IsSTName(getString(rows[0], "name")), nil
}

// IsSTName checks a short name for the ST or *ST prefix.
func IsSTName(name string) bool {
	n := strings.ToUpper(strings.TrimSpace(name))
	return strings.HasPrefix(n, "ST") || strings.HasPrefix(n, "*ST")
}
