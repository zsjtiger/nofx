package store

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDropRetiredColumnsRemovesOldLeverageAndMarginColumns(t *testing.T) {
	dir := t.TempDir()
	st, err := New(filepath.Join(dir, "data.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	// Re-create the crypto-era columns on the current schema.
	require.NoError(t, st.gdb.Exec("ALTER TABLE traders ADD COLUMN btc_eth_leverage INTEGER DEFAULT 5").Error)
	require.NoError(t, st.gdb.Exec("ALTER TABLE traders ADD COLUMN altcoin_leverage INTEGER DEFAULT 5").Error)
	require.NoError(t, st.gdb.Exec("ALTER TABLE traders ADD COLUMN is_cross_margin INTEGER DEFAULT 1").Error)
	require.NoError(t, st.gdb.Exec("ALTER TABLE trader_positions ADD COLUMN leverage INTEGER DEFAULT 1").Error)
	require.NoError(t, st.gdb.Exec("ALTER TABLE trader_orders ADD COLUMN leverage INTEGER DEFAULT 1").Error)
	require.NoError(t, st.gdb.Exec("INSERT INTO traders (id, user_id, name, ai_model_id, exchange_id, initial_balance, btc_eth_leverage, altcoin_leverage) VALUES ('t1','u1','keep','m','e',1000,7,8)").Error)

	require.NoError(t, st.dropRetiredSchema())

	for _, rc := range retiredColumns {
		require.False(t, st.gdb.Migrator().HasColumn(rc.table, rc.column), "%s.%s should be dropped", rc.table, rc.column)
	}
	// The rows that were there keep their data.
	var name string
	require.NoError(t, st.gdb.Raw("SELECT name FROM traders WHERE id = 't1'").Scan(&name).Error)
	require.Equal(t, "keep", name)

	// A backup was written before the change.
	matches, err := filepath.Glob(filepath.Join(dir, "data.db.bak-*"))
	require.NoError(t, err)
	require.Len(t, matches, 1)
}

func TestDropRetiredColumnsIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	st, err := New(filepath.Join(dir, "data.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	// A fresh schema has none of the retired columns: nothing to do, no backup.
	require.NoError(t, st.dropRetiredSchema())
	matches, err := filepath.Glob(filepath.Join(dir, "data.db.bak-*"))
	require.NoError(t, err)
	require.Empty(t, matches)
}

func TestDropRetiredSchemaRemovesTelegramTableWithBackup(t *testing.T) {
	dir := t.TempDir()
	st, err := New(filepath.Join(dir, "data.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	// Re-create the table an older install would still have, with a bot token in it.
	require.NoError(t, st.gdb.Exec("CREATE TABLE telegram_configs (id INTEGER PRIMARY KEY, bot_token TEXT)").Error)
	require.NoError(t, st.gdb.Exec("INSERT INTO telegram_configs (id, bot_token) VALUES (1, 'bot-secret')").Error)

	require.NoError(t, st.dropRetiredSchema())
	require.False(t, st.gdb.Migrator().HasTable("telegram_configs"))

	// The token is still recoverable from the backup taken before the change.
	matches, err := filepath.Glob(filepath.Join(dir, "data.db.bak-*"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	// Running again finds nothing to drop and does not write another backup.
	require.NoError(t, st.dropRetiredSchema())
	matches, err = filepath.Glob(filepath.Join(dir, "data.db.bak-*"))
	require.NoError(t, err)
	require.Len(t, matches, 1)
}
