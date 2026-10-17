package trader

import (
	"fmt"
	"time"

	"crux/kernel"
	"crux/store"
)

const (
	// intelLookback is how far back intel is shown to the model and checked by the gate.
	intelLookback = 72 * time.Hour
	// intelNotesPerSymbol caps the notes sent to the model for one symbol.
	intelNotesPerSymbol = 3
)

// loadRecentIntel returns the trader owner's intel for each symbol within the
// lookback window. A read failure is logged and leaves that symbol without notes,
// so trading continues without intel rather than stopping.
func (at *AutoTrader) loadRecentIntel(symbols []string) map[string][]kernel.IntelNote {
	if at.store == nil || at.userID == "" || len(symbols) == 0 {
		return nil
	}
	since := time.Now().Add(-intelLookback)
	out := make(map[string][]kernel.IntelNote)
	for _, symbol := range symbols {
		items, err := at.store.Intel().ListSince(at.userID, symbol, since, intelNotesPerSymbol)
		if err != nil {
			at.logWarnf("⚠️ Failed to load intel for %s: %v", symbol, err)
			continue
		}
		for _, item := range items {
			out[symbol] = append(out[symbol], intelNoteFromStore(item))
		}
	}
	return out
}

func intelNoteFromStore(item store.IntelItem) kernel.IntelNote {
	return kernel.IntelNote{
		Time:      item.CreatedAt,
		Sentiment: item.Sentiment,
		Risk:      item.Risk,
		Summary:   item.Summary,
	}
}

// intelBlockReason blocks opening a new position in symbol when a recent
// analysed page reports a high risk event. Closing positions is never blocked
// by this check. It returns "" when opening is allowed.
func intelBlockReason(symbol string, ctx *kernel.Context) string {
	if ctx == nil {
		return ""
	}
	cutoff := time.Now().Add(-intelLookback)
	for _, note := range ctx.Intel[symbol] {
		if note.Risk == store.IntelRiskHigh && !note.Time.Before(cutoff) {
			return fmt.Sprintf("intel gate: %s has a high-risk event reported at %s (%s); no new positions for %s",
				symbol, note.Time.UTC().Format("2006-01-02 15:04 UTC"), note.Summary, intelLookback)
		}
	}
	return ""
}
