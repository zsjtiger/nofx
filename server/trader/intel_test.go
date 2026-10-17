package trader

import (
	"testing"
	"time"

	"crux/kernel"
	"crux/store"

	"github.com/stretchr/testify/assert"
)

func ctxWithIntel(symbol string, notes ...kernel.IntelNote) *kernel.Context {
	return &kernel.Context{Intel: map[string][]kernel.IntelNote{symbol: notes}}
}

func TestIntelBlockReasonBlocksRecentHighRisk(t *testing.T) {
	ctx := ctxWithIntel("600519.SH", kernel.IntelNote{
		Time: time.Now().Add(-2 * time.Hour), Risk: store.IntelRiskHigh, Summary: "立案调查",
	})
	reason := intelBlockReason("600519.SH", ctx)
	assert.Contains(t, reason, "intel gate")
	assert.Contains(t, reason, "立案调查")
}

func TestIntelBlockReasonAllowsWhenNoteIsOldLowOrOtherSymbol(t *testing.T) {
	cases := map[string]*kernel.Context{
		"no intel":     {},
		"nil context":  nil,
		"expired high": ctxWithIntel("600519.SH", kernel.IntelNote{Time: time.Now().Add(-intelLookback - time.Hour), Risk: store.IntelRiskHigh}),
		"low risk":     ctxWithIntel("600519.SH", kernel.IntelNote{Time: time.Now(), Risk: store.IntelRiskLow}),
		"unknown risk": ctxWithIntel("600519.SH", kernel.IntelNote{Time: time.Now(), Risk: store.IntelRiskUnknown}),
		"other symbol": ctxWithIntel("000001.SZ", kernel.IntelNote{Time: time.Now(), Risk: store.IntelRiskHigh}),
	}
	for name, ctx := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Empty(t, intelBlockReason("600519.SH", ctx))
		})
	}
}
