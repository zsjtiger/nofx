package trader

import (
	"testing"
	"time"

	"crux/kernel"
	"crux/store"

	"github.com/stretchr/testify/assert"
)

func symbolsOf(coins []kernel.CandidateCoin) []string {
	out := make([]string, 0, len(coins))
	for _, c := range coins {
		out = append(out, c.Symbol)
	}
	return out
}

func candidates(symbols ...string) []kernel.CandidateCoin {
	out := make([]kernel.CandidateCoin, 0, len(symbols))
	for _, s := range symbols {
		out = append(out, kernel.CandidateCoin{Symbol: s})
	}
	return out
}

func note(hoursAgo float64, sentiment float64, risk string) kernel.IntelNote {
	return kernel.IntelNote{
		Time:      time.Now().Add(-time.Duration(hoursAgo * float64(time.Hour))),
		Sentiment: sentiment,
		Risk:      risk,
	}
}

func TestScreenDisabledLeavesCandidatesUntouched(t *testing.T) {
	in := candidates("A", "B", "C")
	notes := map[string][]kernel.IntelNote{"A": {note(1, -1, store.IntelRiskHigh)}}
	out, dropped := screenCandidatesByIntel(in, notes, store.IntelScreenConfig{Enabled: false, ExcludeHighRisk: true}, time.Now())
	assert.Equal(t, []string{"A", "B", "C"}, symbolsOf(out))
	assert.Empty(t, dropped)
}

func TestScreenExcludesHighRiskAndLowSentiment(t *testing.T) {
	in := candidates("A", "B", "C", "D")
	notes := map[string][]kernel.IntelNote{
		"A": {note(1, 0.2, store.IntelRiskHigh)}, // high risk -> dropped
		"B": {note(1, -0.6, store.IntelRiskLow)}, // sentiment below threshold -> dropped
		"C": {note(1, 0.5, store.IntelRiskLow)},  // kept
	}
	cfg := store.IntelScreenConfig{Enabled: true, ExcludeHighRisk: true, MinSentiment: -0.3}
	out, dropped := screenCandidatesByIntel(in, notes, cfg, time.Now())
	assert.Equal(t, []string{"C", "D"}, symbolsOf(out), "candidates without intel must be kept")
	assert.Len(t, dropped, 2)
}

func TestScreenRanksBySentimentWithUnknownsAsZeroAndStable(t *testing.T) {
	in := candidates("A", "B", "C", "D")
	notes := map[string][]kernel.IntelNote{
		"A": {note(1, -0.2, store.IntelRiskNone)},
		"C": {note(1, 0.8, store.IntelRiskNone)},
		"D": {note(2, 0.4, store.IntelRiskNone), note(3, 0.0, store.IntelRiskNone)}, // average 0.2
	}
	out, _ := screenCandidatesByIntel(in, notes, store.IntelScreenConfig{Enabled: true, MinSentiment: -1}, time.Now())
	// C 0.8, D 0.2, B 0 (no intel, original order), A -0.2
	assert.Equal(t, []string{"C", "D", "B", "A"}, symbolsOf(out))
}

func TestScreenIgnoresNotesOutsideLookback(t *testing.T) {
	in := candidates("A")
	notes := map[string][]kernel.IntelNote{"A": {note(100, -1, store.IntelRiskHigh)}}
	cfg := store.IntelScreenConfig{Enabled: true, ExcludeHighRisk: true, LookbackHours: 72}
	out, dropped := screenCandidatesByIntel(in, notes, cfg, time.Now())
	assert.Equal(t, []string{"A"}, symbolsOf(out))
	assert.Empty(t, dropped)

	// A longer window brings the same note back into play.
	cfg.LookbackHours = 168
	out, _ = screenCandidatesByIntel(in, notes, cfg, time.Now())
	assert.Empty(t, out)
}

func TestIntelScreenLookbackDefaultsAndCaps(t *testing.T) {
	assert.Equal(t, 72*time.Hour, intelScreenLookback(store.IntelScreenConfig{}))
	assert.Equal(t, 24*time.Hour, intelScreenLookback(store.IntelScreenConfig{LookbackHours: 24}))
	assert.Equal(t, 168*time.Hour, intelScreenLookback(store.IntelScreenConfig{LookbackHours: 9999}))
}
