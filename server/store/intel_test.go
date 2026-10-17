package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntelListSinceFiltersByUserSymbolAndTime(t *testing.T) {
	st, err := New(filepath.Join(t.TempDir(), "intel.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })
	intel := st.Intel()

	old := &IntelItem{UserID: "u1", Symbol: "600519.SH", SourceURL: "https://a", Risk: IntelRiskHigh}
	require.NoError(t, intel.Create(old))
	require.NoError(t, st.gdb.Model(old).Update("created_at", time.Now().Add(-100*time.Hour)).Error)
	recent := &IntelItem{UserID: "u1", Symbol: "600519.SH", SourceURL: "https://b", Risk: IntelRiskLow}
	require.NoError(t, intel.Create(recent))
	require.NoError(t, intel.Create(&IntelItem{UserID: "u2", Symbol: "600519.SH", SourceURL: "https://c", Risk: IntelRiskHigh}))

	items, err := intel.ListSince("u1", "600519.SH", time.Now().Add(-72*time.Hour), 10)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "https://b", items[0].SourceURL)
}
