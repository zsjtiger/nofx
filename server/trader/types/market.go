package types

import (
	"fmt"
	"regexp"
	"strings"
)

// Market identifies the asset class a trader operates in. It decides which
// order types, symbol formats, lot sizes and settlement rules apply.
type Market string

const (
	// MarketUSEquity is US stocks (Eastern Time, USD).
	MarketUSEquity Market = "us_equity"
	// MarketCNEquity is A-shares on Shanghai and Shenzhen (Beijing Time, CNY).
	MarketCNEquity Market = "cn_equity"
)

// Capabilities describes what a market allows. Risk and decision code reads
// these instead of hard-coding per-exchange behaviour.
type Capabilities struct {
	Market     Market
	Currency   string
	Timezone   string // IANA name, e.g. "Asia/Shanghai"
	AllowShort bool   // false: opening a short position is rejected
	LotSize    int    // minimum buy increment in shares; 0 means no lot rule
	TPlusOne   bool   // shares bought today cannot be sold today
}

var capabilityTable = map[Market]Capabilities{
	MarketUSEquity: {
		Market:     MarketUSEquity,
		Currency:   "USD",
		Timezone:   "America/New_York",
		AllowShort: true,
	},
	MarketCNEquity: {
		Market:     MarketCNEquity,
		Currency:   "CNY",
		Timezone:   "Asia/Shanghai",
		AllowShort: false,
		LotSize:    100,
		TPlusOne:   true,
	},
}

// ParseMarket validates a market name from configuration.
func ParseMarket(s string) (Market, error) {
	m := Market(strings.ToLower(strings.TrimSpace(s)))
	if _, ok := capabilityTable[m]; !ok {
		return "", fmt.Errorf("unknown market %q", s)
	}
	return m, nil
}

// CapabilitiesOf returns the capability set for a market.
func CapabilitiesOf(m Market) (Capabilities, error) {
	c, ok := capabilityTable[m]
	if !ok {
		return Capabilities{}, fmt.Errorf("unknown market %q", m)
	}
	return c, nil
}

var (
	usSymbolPattern = regexp.MustCompile(`^[A-Z]{1,5}(\.[A-Z])?$`)
	cnCodePattern   = regexp.MustCompile(`^[0-9]{6}$`)
)

// NormalizeSymbol converts user input to the canonical symbol form for a market.
//   - us_equity: upper-case ticker, e.g. "aapl" -> "AAPL", "brk.b" -> "BRK.B"
//   - cn_equity: six digits with an exchange suffix, e.g. "600519" -> "600519.SH",
//     "000001.sz" -> "000001.SZ". The suffix is inferred from the code prefix when omitted.
func NormalizeSymbol(m Market, raw string) (string, error) {
	s := strings.ToUpper(strings.TrimSpace(raw))
	if s == "" {
		return "", fmt.Errorf("empty symbol")
	}
	switch m {
	case MarketUSEquity:
		if !usSymbolPattern.MatchString(s) {
			return "", fmt.Errorf("invalid US ticker %q", raw)
		}
		return s, nil
	case MarketCNEquity:
		return normalizeCNSymbol(s, raw)
	default:
		return "", fmt.Errorf("unknown market %q", m)
	}
}

func normalizeCNSymbol(s, raw string) (string, error) {
	code, suffix, hasSuffix := strings.Cut(s, ".")
	if !cnCodePattern.MatchString(code) {
		return "", fmt.Errorf("invalid A-share code %q", raw)
	}
	if hasSuffix {
		switch suffix {
		case "SH", "SZ", "BJ":
			return code + "." + suffix, nil
		default:
			return "", fmt.Errorf("invalid A-share exchange suffix %q", raw)
		}
	}
	switch {
	case strings.HasPrefix(code, "6"): // Shanghai main board and STAR market
		return code + ".SH", nil
	case strings.HasPrefix(code, "0"), strings.HasPrefix(code, "3"): // Shenzhen main board and ChiNext
		return code + ".SZ", nil
	case strings.HasPrefix(code, "4"), strings.HasPrefix(code, "8"): // Beijing Stock Exchange
		return code + ".BJ", nil
	default:
		return "", fmt.Errorf("cannot infer exchange for A-share code %q", raw)
	}
}

// ValidateOrderQuantity checks a quantity against the market's lot rule.
// Buys must be whole lots. Sells may be an odd lot only when it closes the
// whole position, which the caller signals with closesPosition.
func ValidateOrderQuantity(c Capabilities, qty float64, isBuy, closesPosition bool) error {
	if qty <= 0 {
		return fmt.Errorf("quantity must be positive, got %v", qty)
	}
	if c.LotSize == 0 {
		return nil
	}
	lot := float64(c.LotSize)
	whole := int64(qty)
	if float64(whole) != qty {
		return fmt.Errorf("quantity %v must be a whole number of shares", qty)
	}
	if whole%int64(c.LotSize) == 0 {
		return nil
	}
	if !isBuy && closesPosition {
		return nil
	}
	return fmt.Errorf("quantity %v must be a multiple of %v shares", qty, lot)
}

// ValidateSideAllowed rejects actions the market does not permit.
func ValidateSideAllowed(c Capabilities, isShort bool) error {
	if isShort && !c.AllowShort {
		return fmt.Errorf("short selling is not allowed in %s", c.Market)
	}
	return nil
}

// ResolveMarket maps a configured market name to a Market. An empty value is
// an error: every trader must name its market.
func ResolveMarket(s string) (Market, error) {
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("market is required")
	}
	return ParseMarket(s)
}

// FloorToLot rounds a buy quantity down to the market's lot size. Markets
// without a lot rule are returned unchanged.
func FloorToLot(c Capabilities, qty float64) float64 {
	if c.LotSize <= 0 || qty <= 0 {
		return qty
	}
	lot := float64(c.LotSize)
	return float64(int64(qty/lot)) * lot
}
