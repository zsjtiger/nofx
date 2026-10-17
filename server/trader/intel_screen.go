package trader

import (
	"fmt"
	"sort"
	"time"

	"crux/kernel"
	"crux/store"
)

const (
	defaultIntelScreenLookback = 72 * time.Hour
	maxIntelScreenLookback     = 168 * time.Hour
)

// intelScreenLookback returns the configured window, applying the default and cap.
func intelScreenLookback(cfg store.IntelScreenConfig) time.Duration {
	if cfg.LookbackHours <= 0 {
		return defaultIntelScreenLookback
	}
	d := time.Duration(cfg.LookbackHours) * time.Hour
	if d > maxIntelScreenLookback {
		return maxIntelScreenLookback
	}
	return d
}

// screenCandidatesByIntel filters and ranks candidates using their intel notes
// created within the screen's lookback window. A candidate is dropped when it has
// a high-risk note (if ExcludeHighRisk) or when its average sentiment is below
// MinSentiment. Candidates without notes are kept. The rest are ordered by average
// sentiment, highest first; candidates without notes count as 0 and keep their
// original order among equals. The dropped list gives a reason for each removal.
func screenCandidatesByIntel(candidates []kernel.CandidateCoin, notes map[string][]kernel.IntelNote, cfg store.IntelScreenConfig, now time.Time) ([]kernel.CandidateCoin, []string) {
	if !cfg.Enabled {
		return candidates, nil
	}
	cutoff := now.Add(-intelScreenLookback(cfg))

	type scored struct {
		coin  kernel.CandidateCoin
		score float64
		order int
	}
	var kept []scored
	var dropped []string

	for i, coin := range candidates {
		var inWindow []kernel.IntelNote
		for _, n := range notes[coin.Symbol] {
			if !n.Time.Before(cutoff) {
				inWindow = append(inWindow, n)
			}
		}

		if cfg.ExcludeHighRisk && hasHighRisk(inWindow) {
			dropped = append(dropped, fmt.Sprintf("%s: high-risk intel in window", coin.Symbol))
			continue
		}

		score := 0.0
		if len(inWindow) > 0 {
			sum := 0.0
			for _, n := range inWindow {
				sum += n.Sentiment
			}
			score = sum / float64(len(inWindow))
			if score < cfg.MinSentiment {
				dropped = append(dropped, fmt.Sprintf("%s: average sentiment %+.2f below %+.2f", coin.Symbol, score, cfg.MinSentiment))
				continue
			}
		}
		kept = append(kept, scored{coin: coin, score: score, order: i})
	}

	sort.SliceStable(kept, func(a, b int) bool {
		if kept[a].score != kept[b].score {
			return kept[a].score > kept[b].score
		}
		return kept[a].order < kept[b].order
	})

	out := make([]kernel.CandidateCoin, 0, len(kept))
	for _, s := range kept {
		out = append(out, s.coin)
	}
	return out, dropped
}

func hasHighRisk(notes []kernel.IntelNote) bool {
	for _, n := range notes {
		if n.Risk == store.IntelRiskHigh {
			return true
		}
	}
	return false
}
