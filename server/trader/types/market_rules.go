package types

import (
	"fmt"
	"math"
	"strings"
	"time"
	_ "time/tzdata" // embed zone data so session checks work on minimal hosts
)

// Session is a daily trading window in the market's local time.
type Session struct {
	Start string // "HH:MM"
	End   string // "HH:MM", exclusive
}

// marketSessions lists the regular trading windows per market.
var marketSessions = map[Market][]Session{
	MarketCNEquity: {{Start: "09:30", End: "11:30"}, {Start: "13:00", End: "15:00"}},
	MarketUSEquity: {{Start: "09:30", End: "16:00"}},
}

func (c Capabilities) location() (*time.Location, error) {
	if c.Timezone == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(c.Timezone)
}

// InTradingSession reports whether t falls inside a regular session on a weekday.
// Exchange holidays are not modelled; markets with no session table are always open.
func (c Capabilities) InTradingSession(t time.Time) (bool, error) {
	sessions, ok := marketSessions[c.Market]
	if !ok {
		return true, nil
	}
	loc, err := c.location()
	if err != nil {
		return false, err
	}
	local := t.In(loc)
	if local.Weekday() == time.Saturday || local.Weekday() == time.Sunday {
		return false, nil
	}
	minute := local.Hour()*60 + local.Minute()
	for _, s := range sessions {
		start, err := parseClock(s.Start)
		if err != nil {
			return false, err
		}
		end, err := parseClock(s.End)
		if err != nil {
			return false, err
		}
		if minute >= start && minute < end {
			return true, nil
		}
	}
	return false, nil
}

func parseClock(hhmm string) (int, error) {
	var h, m int
	if _, err := fmt.Sscanf(hhmm, "%d:%d", &h, &m); err != nil {
		return 0, fmt.Errorf("invalid clock %q: %w", hhmm, err)
	}
	return h*60 + m, nil
}

// SameTradingDay reports whether a and b fall on the same calendar date in the market's timezone.
func (c Capabilities) SameTradingDay(a, b time.Time) (bool, error) {
	loc, err := c.location()
	if err != nil {
		return false, err
	}
	ay, am, ad := a.In(loc).Date()
	by, bm, bd := b.In(loc).Date()
	return ay == by && am == bm && ad == bd, nil
}

// ValidateSettlement blocks closing a position that T+1 still locks.
// A-shares bought on the current trading day cannot be sold until the next day.
func ValidateSettlement(c Capabilities, entryTime, now time.Time) error {
	if !c.TPlusOne || entryTime.IsZero() {
		return nil
	}
	same, err := c.SameTradingDay(entryTime, now)
	if err != nil {
		return err
	}
	if same {
		return fmt.Errorf("shares bought on %s cannot be sold until the next trading day (T+1)", entryTime.In(mustLoc(c)).Format("2006-01-02"))
	}
	return nil
}

func mustLoc(c Capabilities) *time.Location {
	loc, err := c.location()
	if err != nil {
		return time.UTC
	}
	return loc
}

// PriceLimits returns the daily limit-up and limit-down prices for an A-share.
// Main board ±10%, ChiNext and STAR ±20%, Beijing ±30%, ST names ±5%.
// Prices are rounded to the nearest cent.
func PriceLimits(symbol string, prevClose float64, isST bool) (low, high float64, err error) {
	if prevClose <= 0 {
		return 0, 0, fmt.Errorf("previous close must be positive, got %v", prevClose)
	}
	code, _, _ := strings.Cut(symbol, ".")
	var pct float64
	switch {
	case isST:
		pct = 0.05
	case strings.HasPrefix(code, "300"), strings.HasPrefix(code, "301"), strings.HasPrefix(code, "688"):
		pct = 0.20
	case strings.HasPrefix(code, "4"), strings.HasPrefix(code, "8"):
		pct = 0.30
	default:
		pct = 0.10
	}
	round := func(v float64) float64 { return math.Round(v*100) / 100 }
	return round(prevClose * (1 - pct)), round(prevClose * (1 + pct)), nil
}

// ValidatePriceWithinLimits rejects an A-share order priced outside the daily limits.
func ValidatePriceWithinLimits(symbol string, prevClose, price float64, isST bool) error {
	low, high, err := PriceLimits(symbol, prevClose, isST)
	if err != nil {
		return err
	}
	if price > high+1e-9 || price < low-1e-9 {
		return fmt.Errorf("price %.2f is outside the daily limit [%.2f, %.2f] for %s", price, low, high, symbol)
	}
	return nil
}
