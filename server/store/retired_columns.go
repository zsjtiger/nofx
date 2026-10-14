package store

import (
	"fmt"
	"time"

	"nofx/logger"
)

// retiredColumn is a column of the crypto-era schema that the stock version no
// longer reads or writes.
type retiredColumn struct {
	table  string
	column string
}

// retiredColumns lists the leverage and margin-mode columns. Dropping them is
// a one-way change, so a SQLite database is backed up first.
var retiredColumns = []retiredColumn{
	{"traders", "btc_eth_leverage"},
	{"traders", "altcoin_leverage"},
	{"traders", "is_cross_margin"},
	{"trader_positions", "leverage"},
	{"trader_orders", "leverage"},
}

// dropRetiredColumns removes the retired columns that still exist. It is safe to
// run on every start: once a column is gone it is skipped.
func (s *Store) dropRetiredColumns() error {
	m := s.gdb.Migrator()
	var pending []retiredColumn
	for _, rc := range retiredColumns {
		if m.HasColumn(rc.table, rc.column) {
			pending = append(pending, rc)
		}
	}
	if len(pending) == 0 {
		return nil
	}

	if s.gdb.Dialector.Name() == "sqlite" {
		backup, err := backupSQLite(s)
		if err != nil {
			return fmt.Errorf("failed to back up database before dropping retired columns: %w", err)
		}
		logger.Infof("🗄️ Backed up database to %s before dropping retired columns", backup)
	}

	for _, rc := range pending {
		if err := s.gdb.Exec(dropColumnSQL(s.gdb.Dialector.Name(), rc)).Error; err != nil {
			return fmt.Errorf("failed to drop %s.%s: %w", rc.table, rc.column, err)
		}
		logger.Infof("🗄️ Dropped retired column %s.%s", rc.table, rc.column)
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
