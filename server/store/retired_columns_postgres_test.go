package store

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDropRetiredColumnsOnPostgres runs the retired-column migration against a
// real PostgreSQL database. It is skipped unless CRUX_TEST_POSTGRES=1 is set,
// and it uses the DB_* variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD,
// DB_NAME, DB_SSLMODE). Point DB_NAME at a throwaway database: the test adds
// and drops columns on the traders, trader_positions and trader_orders tables.
func TestDropRetiredColumnsOnPostgres(t *testing.T) {
	if os.Getenv("CRUX_TEST_POSTGRES") != "1" {
		t.Skip("set CRUX_TEST_POSTGRES=1 and the DB_* variables to run against PostgreSQL")
	}
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	require.NoError(t, err)
	st, err := NewWithConfig(DBConfig{
		Type:     DBTypePostgres,
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   getEnv("DB_NAME", "crux_test"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	// Re-create the crypto-era columns on the current schema.
	for _, stmt := range []string{
		"ALTER TABLE traders ADD COLUMN IF NOT EXISTS btc_eth_leverage INTEGER DEFAULT 5",
		"ALTER TABLE traders ADD COLUMN IF NOT EXISTS altcoin_leverage INTEGER DEFAULT 5",
		"ALTER TABLE traders ADD COLUMN IF NOT EXISTS is_cross_margin INTEGER DEFAULT 1",
		"ALTER TABLE trader_positions ADD COLUMN IF NOT EXISTS leverage INTEGER DEFAULT 1",
		"ALTER TABLE trader_orders ADD COLUMN IF NOT EXISTS leverage INTEGER DEFAULT 1",
	} {
		require.NoError(t, st.gdb.Exec(stmt).Error, stmt)
	}

	require.NoError(t, st.dropRetiredSchema())

	for _, rc := range retiredColumns {
		require.False(t, st.gdb.Migrator().HasColumn(rc.table, rc.column), "%s.%s should be dropped", rc.table, rc.column)
	}

	// The Telegram table from the removed bot is dropped as well.
	require.NoError(t, st.gdb.Exec("CREATE TABLE IF NOT EXISTS telegram_configs (id BIGINT PRIMARY KEY, bot_token TEXT)").Error)
	require.NoError(t, st.dropRetiredSchema())
	require.False(t, st.gdb.Migrator().HasTable("telegram_configs"))

	// A second run finds nothing to drop and succeeds.
	require.NoError(t, st.dropRetiredSchema())
}
