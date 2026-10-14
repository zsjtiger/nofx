package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMarketForExchangeType(t *testing.T) {
	assert.Equal(t, "us_equity", MarketForExchangeType("alpaca_paper"))
	assert.Equal(t, "cn_equity", MarketForExchangeType(" CN_PAPER "))
	assert.Equal(t, "crypto", MarketForExchangeType("binance"))
	assert.Equal(t, "crypto", MarketForExchangeType(""))
}

func TestMissingRequiredFieldsForEquityExchanges(t *testing.T) {
	missing := MissingRequiredExchangeCredentialFields("alpaca_paper", "", "", "", "", "", "", "", "", "")
	assert.Equal(t, []string{"api_key", "secret_key"}, missing)
	assert.Empty(t, MissingRequiredExchangeCredentialFields("alpaca_paper", "k", "s", "", "", "", "", "", "", ""))

	assert.Equal(t, []string{"api_key (tushare token)"}, MissingRequiredExchangeCredentialFields("cn_paper", "", "", "", "", "", "", "", "", ""))
	assert.Empty(t, MissingRequiredExchangeCredentialFields("cn_paper", "tok", "", "", "", "", "", "", "", ""))
}

// TestEnsureMarketColumnMigratesOldTable covers a database created before the market column existed.
func TestEnsureMarketColumnMigratesOldTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.Exec(`CREATE TABLE exchanges (id TEXT PRIMARY KEY, exchange_type TEXT NOT NULL DEFAULT '')`).Error)
	require.NoError(t, db.Exec(`INSERT INTO exchanges (id, exchange_type) VALUES ('old-1', 'binance')`).Error)

	s := NewExchangeStore(db)
	require.False(t, db.Migrator().HasColumn(&Exchange{}, "Market"))
	require.NoError(t, s.ensureMarketColumn())
	require.True(t, db.Migrator().HasColumn(&Exchange{}, "Market"))

	var market string
	require.NoError(t, db.Raw(`SELECT market FROM exchanges WHERE id = 'old-1'`).Scan(&market).Error)
	assert.Equal(t, "crypto", market, "existing rows default to crypto")

	// Running again is a no-op.
	require.NoError(t, s.ensureMarketColumn())
}
