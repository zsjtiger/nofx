package intel

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"crux/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeFetcher map[string]string

func (f fakeFetcher) Fetch(ctx context.Context, rawURL string) (string, error) {
	text, ok := f[rawURL]
	if !ok {
		return "", errors.New("not found")
	}
	return text, nil
}

func newTestIntelStore(t *testing.T) *store.IntelStore {
	t.Helper()
	st, err := store.New(filepath.Join(t.TempDir(), "intel.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })
	return st.Intel()
}

func TestServiceRunStoresSuccessfulPagesAndReportsFailures(t *testing.T) {
	st := newTestIntelStore(t)
	fetch := fakeFetcher{"https://a.example/1": "茅台年报利润增长"}
	complete := func(system, user string) (string, error) {
		return `{"sentiment": 0.6, "risk": "none", "summary": "利润增长"}`, nil
	}
	svc := NewService(fetch, complete, st)

	results, err := svc.Run(context.Background(), Job{
		UserID: "u1",
		Symbol: " 600519.sh ",
		URLs:   []string{"https://a.example/1", "https://a.example/missing"},
	})
	require.NoError(t, err)
	require.Len(t, results, 2)

	require.NotNil(t, results[0].Item)
	assert.Equal(t, "600519.SH", results[0].Item.Symbol)
	assert.Equal(t, "利润增长", results[0].Item.Summary)
	assert.Empty(t, results[0].Error)

	assert.Nil(t, results[1].Item)
	assert.Contains(t, results[1].Error, "not found")

	items, err := st.ListBySymbol("u1", "600519.SH", 10)
	require.NoError(t, err)
	assert.Len(t, items, 1)

	other, err := st.ListBySymbol("u2", "600519.SH", 10)
	require.NoError(t, err)
	assert.Empty(t, other, "items must be scoped to the user who ran the analysis")
}

func TestServiceRunValidatesJob(t *testing.T) {
	svc := NewService(fakeFetcher{}, nil, newTestIntelStore(t))

	_, err := svc.Run(context.Background(), Job{Symbol: "", URLs: []string{"https://a.example/1"}})
	assert.Error(t, err)

	_, err = svc.Run(context.Background(), Job{Symbol: "600519.SH"})
	assert.Error(t, err)

	tooMany := make([]string, MaxURLsPerRun+1)
	_, err = svc.Run(context.Background(), Job{Symbol: "600519.SH", URLs: tooMany})
	assert.Error(t, err)
}

func TestServiceRunStopsWhenCancelled(t *testing.T) {
	svc := NewService(fakeFetcher{}, nil, newTestIntelStore(t))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Run(ctx, Job{Symbol: "600519.SH", URLs: []string{"https://a.example/1"}})
	assert.ErrorIs(t, err, context.Canceled)
}
