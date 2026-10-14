package cnpaper

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "paper", "trader-1.json")
	quotes := fakeQuoter{"600519.SH": normalQuote}

	first, err := NewPersistent(1_000_000, quotes, path)
	require.NoError(t, err)
	first.now = func() time.Time { return morning }
	_, err = first.OpenLong("600519.SH", 100)
	require.NoError(t, err)
	require.NoError(t, first.SetStopLoss("600519.SH", "LONG", 100, 95))

	// A new process reads the same file and sees the same cash, position and protection.
	second, err := NewPersistent(999, quotes, path)
	require.NoError(t, err, "initial cash is ignored when a ledger exists")
	second.now = func() time.Time { return morning }

	bal, err := second.GetBalance()
	require.NoError(t, err)
	assert.InDelta(t, 1_000_000-10_000-5, bal["totalWalletBalance"], 0.001)

	pos, err := second.GetPositions()
	require.NoError(t, err)
	require.Len(t, pos, 1)
	assert.Equal(t, 100.0, pos[0]["positionAmt"])

	orders, err := second.GetOpenOrders("600519.SH")
	require.NoError(t, err)
	require.Len(t, orders, 1, "stop-loss is restored")
}

func TestClosedTradesPersistAndIDsContinue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.json")
	quotes := fakeQuoter{"600519.SH": normalQuote}

	first, err := NewPersistent(1_000_000, quotes, path)
	require.NoError(t, err)
	first.now = func() time.Time { return morning }
	_, err = first.OpenLong("600519.SH", 100)
	require.NoError(t, err)

	second, err := NewPersistent(1_000_000, quotes, path)
	require.NoError(t, err)
	second.now = func() time.Time { return nextDay }
	_, err = second.CloseLong("600519.SH", 0)
	require.NoError(t, err)

	closed, err := second.GetClosedPnL(time.Time{}, 10)
	require.NoError(t, err)
	assert.Len(t, closed, 1)

	third, err := NewPersistent(1, quotes, path)
	require.NoError(t, err)
	closedAgain, _ := third.GetClosedPnL(time.Time{}, 10)
	assert.Len(t, closedAgain, 1, "closed trades are restored")
	assert.Equal(t, "cnpaper-3", third.newOrderID(), "order numbering continues after restart")
}

func TestNoStatePathKeepsMemoryOnly(t *testing.T) {
	tr, err := NewPersistent(1000, fakeQuoter{}, "")
	require.NoError(t, err)
	assert.NoError(t, tr.save(), "saving without a path is a no-op")
}

func TestCorruptLedgerIsAnError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	require.NoError(t, writeFile(path, "{not json"))
	_, err := NewPersistent(1000, fakeQuoter{}, path)
	assert.ErrorContains(t, err, "parse paper ledger")
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600)
}
