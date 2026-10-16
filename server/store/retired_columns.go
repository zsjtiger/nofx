package store

import (
	"fmt"
	"time"

	"crux/logger"
)

// retiredColumn is a column of the crypto-era schema that the stock version no
// longer reads or writes.
type retiredColumn struct {
	table  string
	column string
}

// retiredColumns lists the leverage and margin-mode columns of the crypto-era
// schema.
var retiredColumns = []retiredColumn{
	{"traders", "btc_eth_leverage"},
	{"traders", "altcoin_leverage"},
	{"traders", "is_cross_margin"},
	{"trader_positions", "leverage"},
	{"trader_orders", "leverage"},
}

// retiredTables lists tables of the crypto-era schema that nothing reads any
// more. telegram_configs holds the Telegram bot token, so it is removed with the
// bot code.
var retiredTables = []string{
	"telegram_configs",
}

// dropRetiredSchema removes the retired columns and tables that still exist. It
// is safe to run on every start: anything already gone is skipped. Both changes
// are one-way, so a SQLite database is backed up first.
func (s *Store) dropRetiredSchema() error {
	m := s.gdb.Migrator()
	var pendingColumns []retiredColumn
	for _, rc := range retiredColumns {
		if m.HasColumn(rc.table, rc.column) {
			pendingColumns = append(pendingColumns, rc)
		}
	}
	var pendingTables []string
	for _, table := range retiredTables {
		if m.HasTable(table) {
			pendingTables = append(pendingTables, table)
		}
	}
	if len(pendingColumns) == 0 && len(pendingTables) == 0 {
		return nil
	}

	if s.gdb.Dialector.Name() == "sqlite" {
		backup, err := backupSQLite(s)
		if err != nil {
			return fmt.Errorf("failed to back up database before dropping retired schema: %w", err)
		}
		logger.Infof("🗄️ Backed up database to %s before dropping retired schema", backup)
	}

	for _, rc := range pendingColumns {
		if err := s.gdb.Exec(dropColumnSQL(s.gdb.Dialector.Name(), rc)).Error; err != nil {
			return fmt.Errorf("failed to drop %s.%s: %w", rc.table, rc.column, err)
		}
		logger.Infof("🗄️ Dropped retired column %s.%s", rc.table, rc.column)
	}
	for _, table := range pendingTables {
		if err := s.gdb.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		logger.Infof("🗄️ Dropped retired table %s", table)
	}
	return nil
}

// dropColumnSQL returns the statement that drops a column. The bundled SQLite
// (3.35+) and PostgreSQL both support ALTER TABLE ... DROP COLUMN natively. GORM's
// SQLite migrator is not used because it needs a struct that still has the field.
func dropColumnSQL(dialect string, rc retiredColumn) string {
	if dialect == "postgres" {
		return fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s", rc.table, rc.column)
	}
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", rc.table, rc.column)
}

// backupSQLite writes a consistent copy of the open SQLite database next to it
// and returns the copy's path.
func backupSQLite(s *Store) (string, error) {
	var file string
	if err := s.gdb.Raw("SELECT file FROM pragma_database_list WHERE name = 'main'").Scan(&file).Error; err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("database is not file-backed")
	}
	backup := fmt.Sprintf("%s.bak-%s", file, time.Now().UTC().Format("20060102-150405"))
	if err := s.gdb.Exec("VACUUM INTO ?", backup).Error; err != nil {
		return "", err
	}
	return backup, nil
}
