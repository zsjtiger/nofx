package cnpaper

import (
	"context"
	"fmt"
	"time"

	"crux/provider/tushare"
)

// TushareQuoter prices A-shares from Tushare daily data. Tushare has no free
// intraday feed, so the price is the latest daily close. Paper fills therefore
// reflect end-of-day prices, not intraday execution.
type TushareQuoter struct {
	client *tushare.Client
	now    func() time.Time
}

// NewTushareQuoter returns a quoter that reads through the given Tushare client.
func NewTushareQuoter(client *tushare.Client) *TushareQuoter {
	return &TushareQuoter{client: client, now: time.Now}
}

// Quote returns the latest daily bar and that day's limit prices.
func (q *TushareQuoter) Quote(ctx context.Context, symbol string) (Quote, error) {
	end := q.now().Format("20060102")
	start := q.now().AddDate(0, 0, -15).Format("20060102")
	bars, err := q.client.DailyBars(ctx, symbol, start, end)
	if err != nil {
		return Quote{}, err
	}
	if len(bars) == 0 {
		return Quote{}, fmt.Errorf("no daily bars for %s between %s and %s", symbol, start, end)
	}
	last := bars[len(bars)-1]
	up, down, err := q.client.StockLimit(ctx, symbol, last.TradeDate)
	if err != nil {
		return Quote{}, err
	}
	return Quote{
		Price:     last.Close,
		PrevClose: last.PreClose,
		UpLimit:   up,
		DownLimit: down,
		Volume:    last.Vol,
	}, nil
}
