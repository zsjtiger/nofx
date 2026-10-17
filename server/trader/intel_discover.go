package trader

import (
	"sort"
	"strings"
	"time"

	"crux/kernel"
	"crux/logger"
	"crux/store"
	"crux/trader/types"
)

const (
	defaultDiscoverNew = 5
	maxDiscoverNew     = 10
)

// applyIntelToCandidates screens the static candidate list with intel and adds
// symbols discovered from intel, as the strategy configures. It returns the list
// unchanged when neither feature is enabled.
func (at *AutoTrader) applyIntelToCandidates(coins []kernel.CandidateCoin) []kernel.CandidateCoin {
	cfg := at.strategyEngineConfig()
	if cfg == nil {
		return coins
	}
	screen := cfg.CoinSource.IntelScreen
	discovery := cfg.CoinSource.IntelDiscovery
	if !screen.Enabled && !discovery.Enabled {
		return coins
	}
	now := time.Now()

	if screen.Enabled {
		symbols := make([]string, 0, len(coins))
		for _, coin := range coins {
			symbols = append(symbols, coin.Symbol)
		}
		notes := at.loadRecentIntel(symbols, intelScreenLookback(screen))
		total := len(coins)
		var dropped []string
		coins, dropped = screenCandidatesByIntel(coins, notes, screen, now)
		for _, reason := range dropped {
			at.logInfof("🔎 intel screen dropped %s", reason)
		}
		logger.Infof("📋 [%s] Intel screen kept %d of %d candidates", at.name, len(coins), total)
	}

	if discovery.Enabled {
		found := at.discoverIntelCandidates(coins, cfg.CoinSource.ExcludedCoins, discovery, now)
		for _, c := range found {
			at.logInfof("🧭 intel discovery added %s", c.Symbol)
		}
		coins = append(coins, found...)
	}
	return coins
}

// discoverIntelCandidates returns symbols mentioned favourably in recent intel
// that are not already candidates. The list is capped so the total stays within
// store.MaxCandidateSymbols and MaxNew.
func (at *AutoTrader) discoverIntelCandidates(existing []kernel.CandidateCoin, excluded []string, cfg store.IntelDiscoveryConfig, now time.Time) []kernel.CandidateCoin {
	if at.store == nil || at.userID == "" {
		return nil
	}
	maxNew := cfg.MaxNew
	if maxNew <= 0 {
		maxNew = defaultDiscoverNew
	}
	if maxNew > maxDiscoverNew {
		maxNew = maxDiscoverNew
	}
	if room := store.MaxCandidateSymbols - len(existing); maxNew > room {
		maxNew = room
	}
	if maxNew <= 0 {
		return nil
	}

	market := at.marketCapabilities().Market
	exclude := make(map[string]bool, len(existing)+len(excluded))
	for _, coin := range existing {
		exclude[coin.Symbol] = true
	}
	for _, raw := range excluded {
		if sym, err := types.NormalizeSymbol(market, raw); err == nil {
			exclude[sym] = true
		}
	}

	since := now.Add(-intelDiscoveryLookback(cfg))
	items, err := at.store.Intel().ListSinceAll(at.userID, since, 1000)
	if err != nil {
		at.logWarnf("⚠️ Failed to load intel for discovery: %v", err)
		return nil
	}

	symbols := rankDiscoveredSymbols(items, market, cfg.MinSentiment, exclude, since)
	if len(symbols) > maxNew {
		symbols = symbols[:maxNew]
	}
	out := make([]kernel.CandidateCoin, 0, len(symbols))
	for _, sym := range symbols {
		out = append(out, kernel.CandidateCoin{Symbol: sym, Sources: []string{"intel"}})
	}
	return out
}

// rankDiscoveredSymbols scores each stock mentioned in recent intel. A page's
// subject and the stocks it mentions take the page's sentiment. Pages with a
// high-risk flag are skipped. A stock is kept when its average sentiment is at
// least minSentiment. Results are ordered by average sentiment, then mention count,
// then symbol, so the output is deterministic.
func rankDiscoveredSymbols(items []store.IntelItem, market types.Market, minSentiment float64, exclude map[string]bool, since time.Time) []string {
	type score struct {
		sum      float64
		mentions int
	}
	scores := make(map[string]*score)

	for _, item := range items {
		if item.CreatedAt.Before(since) || item.Risk == store.IntelRiskHigh {
			continue
		}
		mentioned := append([]string{item.Symbol}, splitSymbols(item.Symbols)...)
		seen := make(map[string]bool, len(mentioned))
		for _, raw := range mentioned {
			sym, err := types.NormalizeSymbol(market, raw)
			if err != nil || exclude[sym] || seen[sym] {
				continue
			}
			seen[sym] = true
			s := scores[sym]
			if s == nil {
				s = &score{}
				scores[sym] = s
			}
			s.sum += item.Sentiment
			s.mentions++
		}
	}

	type ranked struct {
		symbol string
		avg    float64
		n      int
	}
	var list []ranked
	for sym, s := range scores {
		avg := s.sum / float64(s.mentions)
		if avg < minSentiment {
			continue
		}
		list = append(list, ranked{symbol: sym, avg: avg, n: s.mentions})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].avg != list[j].avg {
			return list[i].avg > list[j].avg
		}
		if list[i].n != list[j].n {
			return list[i].n > list[j].n
		}
		return list[i].symbol < list[j].symbol
	})

	out := make([]string, 0, len(list))
	for _, r := range list {
		out = append(out, r.symbol)
	}
	return out
}

func splitSymbols(csv string) []string {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	return strings.Split(csv, ",")
}

// intelDiscoveryLookback returns the mention window, applying the default and cap.
func intelDiscoveryLookback(cfg store.IntelDiscoveryConfig) time.Duration {
	if cfg.LookbackHours <= 0 {
		return defaultIntelScreenLookback
	}
	d := time.Duration(cfg.LookbackHours) * time.Hour
	if d > maxIntelScreenLookback {
		return maxIntelScreenLookback
	}
	return d
}

// intelContextLookback is the window of intel shown to the model: the longest
// window any enabled intel feature uses, and at least the open-position gate's window.
func (at *AutoTrader) intelContextLookback() time.Duration {
	window := intelLookback
	cfg := at.strategyEngineConfig()
	if cfg == nil {
		return window
	}
	if cfg.CoinSource.IntelScreen.Enabled {
		if d := intelScreenLookback(cfg.CoinSource.IntelScreen); d > window {
			window = d
		}
	}
	if cfg.CoinSource.IntelDiscovery.Enabled {
		if d := intelDiscoveryLookback(cfg.CoinSource.IntelDiscovery); d > window {
			window = d
		}
	}
	return window
}
