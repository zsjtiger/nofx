package store

import (
	"time"

	"gorm.io/gorm"
)

// IntelItem is one analysed web page about a stock. Items are scoped to the
// user who ran the analysis.
type IntelItem struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"column:user_id;not null;index" json:"user_id"`
	Symbol    string    `gorm:"column:symbol;not null;index" json:"symbol"`
	SourceURL string    `gorm:"column:source_url;not null" json:"source_url"`
	Summary   string    `gorm:"column:summary;not null;default:''" json:"summary"`
	Sentiment float64   `gorm:"column:sentiment;not null;default:0" json:"sentiment"`
	Risk      string    `gorm:"column:risk;not null;default:'unknown'" json:"risk"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName returns the table name for IntelItem
func (IntelItem) TableName() string {
	return "intel_items"
}

// IntelStore persists intelligence analysis results.
type IntelStore struct {
	db *gorm.DB
}

// NewIntelStore creates an intel store.
func NewIntelStore(db *gorm.DB) *IntelStore {
	return &IntelStore{db: db}
}

func (s *IntelStore) initTables() error {
	return s.db.AutoMigrate(&IntelItem{})
}

// Create stores one analysed item.
func (s *IntelStore) Create(item *IntelItem) error {
	return s.db.Create(item).Error
}

// ListBySymbol returns the newest items for a symbol, newest first.
func (s *IntelStore) ListBySymbol(userID, symbol string, limit int) ([]IntelItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var items []IntelItem
	err := s.db.Where("user_id = ? AND symbol = ?", userID, symbol).
		Order("created_at DESC, id DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}
