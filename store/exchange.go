package store

import (
	"fmt"
	"nofx/logger"
	"nofx/secrets"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExchangeStore exchange storage
type ExchangeStore struct {
	db *gorm.DB
}

// Exchange is a brokerage account used for paper trading. The market is fixed by the exchange type.
type Exchange struct {
	ID           string                  `gorm:"primaryKey" json:"id"`
	ExchangeType string                  `gorm:"column:exchange_type;not null;default:''" json:"exchange_type"`
	AccountName  string                  `gorm:"column:account_name;not null;default:''" json:"account_name"`
	UserID       string                  `gorm:"column:user_id;not null;default:default;index" json:"user_id"`
	Name         string                  `gorm:"not null" json:"name"`
	Type         string                  `gorm:"not null" json:"type"`
	Enabled      bool                    `gorm:"default:false" json:"enabled"`
	APIKey       secrets.EncryptedString `gorm:"column:api_key;default:''" json:"apiKey"`
	SecretKey    secrets.EncryptedString `gorm:"column:secret_key;default:''" json:"secretKey"`
	Market       string                  `gorm:"column:market;not null;default:''" json:"market"` // "us_equity" or "cn_equity"
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
}

func (Exchange) TableName() string { return "exchanges" }

// NewExchangeStore creates a new ExchangeStore
func NewExchangeStore(db *gorm.DB) *ExchangeStore {
	return &ExchangeStore{db: db}
}

func (s *ExchangeStore) initTables() error {
	// For PostgreSQL with existing table, skip AutoMigrate and run the column migrations only.
	if s.db.Dialector.Name() == "postgres" {
		var tableExists int64
		s.db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'exchanges'`).Scan(&tableExists)
		if tableExists > 0 {
			if err := s.ensureMarketColumn(); err != nil {
				logger.Warnf("Exchange market column migration warning: %v", err)
			}
			if err := s.cleanupIncompleteExchangeConfigs(); err != nil {
				logger.Warnf("Exchange cleanup migration warning: %v", err)
			}
			return nil
		}
	}

	if err := s.db.AutoMigrate(&Exchange{}); err != nil {
		return err
	}
	if err := s.ensureMarketColumn(); err != nil {
		logger.Warnf("Exchange market column migration warning: %v", err)
	}
	s.db.Model(&Exchange{}).Where("account_name = '' OR account_name IS NULL").Update("account_name", "Default")
	if err := s.cleanupIncompleteExchangeConfigs(); err != nil {
		logger.Warnf("Exchange cleanup migration warning: %v", err)
	}
	return nil
}

// ensureMarketColumn adds the market column to databases created before per-market
// support, and fills it from the exchange type for existing rows.
func (s *ExchangeStore) ensureMarketColumn() error {
	if !s.db.Migrator().HasColumn(&Exchange{}, "Market") {
		if err := s.db.Migrator().AddColumn(&Exchange{}, "Market"); err != nil {
			return err
		}
	}
	for _, typ := range []string{"alpaca_paper", "cn_paper"} {
		if err := s.db.Model(&Exchange{}).
			Where("exchange_type = ? AND (market = '' OR market IS NULL)", typ).
			UpdateColumn("market", MarketForExchangeType(typ)).Error; err != nil {
			return err
		}
	}
	return nil
}

// cleanupIncompleteExchangeConfigs removes accounts that lack credentials and
// enables complete ones. Accounts of retired exchange types are removed as well.
func (s *ExchangeStore) cleanupIncompleteExchangeConfigs() error {
	var exchanges []Exchange
	if err := s.db.Find(&exchanges).Error; err != nil {
		return err
	}
	for _, exchange := range exchanges {
		missing := MissingRequiredExchangeCredentialFields(exchange.ExchangeType, string(exchange.APIKey), string(exchange.SecretKey))
		if len(missing) > 0 {
			if err := s.db.Delete(&Exchange{}, "id = ? AND user_id = ?", exchange.ID, exchange.UserID).Error; err != nil {
				return err
			}
			logger.Infof("🧹 Removed incomplete exchange config during migration: id=%s user=%s missing=%s", exchange.ID, exchange.UserID, strings.Join(missing, ","))
			continue
		}
		if !exchange.Enabled {
			if err := s.db.Model(&Exchange{}).Where("id = ? AND user_id = ?", exchange.ID, exchange.UserID).Update("enabled", true).Error; err != nil {
				return err
			}
			logger.Infof("🧹 Enabled complete exchange config during migration: id=%s user=%s", exchange.ID, exchange.UserID)
		}
	}
	return nil
}

// List gets user's exchange list
func (s *ExchangeStore) List(userID string) ([]*Exchange, error) {
	var exchanges []*Exchange
	err := s.db.Where("user_id = ?", userID).Order("exchange_type, account_name").Find(&exchanges).Error
	if err != nil {
		return nil, err
	}
	return exchanges, nil
}

// GetByID gets a specific exchange by UUID
func (s *ExchangeStore) GetByID(userID, id string) (*Exchange, error) {
	var exchange Exchange
	err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&exchange).Error
	if err != nil {
		return nil, err
	}
	return &exchange, nil
}

// getExchangeNameAndType returns the display name and type for an exchange type
func getExchangeNameAndType(exchangeType string) (name string, typ string) {
	switch exchangeType {
	case "alpaca_paper":
		return "Alpaca Paper (US Equities)", "paper"
	case "cn_paper":
		return "A-Share Paper (CN)", "paper"
	default:
		return exchangeType + " Exchange", "paper"
	}
}

// Create adds an exchange account. It returns the new account's UUID.
func (s *ExchangeStore) Create(userID, exchangeType, accountName string, enabled bool, apiKey, secretKey string) (string, error) {
	if missing := MissingRequiredExchangeCredentialFields(exchangeType, apiKey, secretKey); len(missing) > 0 {
		return "", fmt.Errorf("missing required exchange fields: %s", strings.Join(missing, ", "))
	}

	id := uuid.New().String()
	name, typ := getExchangeNameAndType(exchangeType)
	if accountName == "" {
		accountName = "Default"
	}

	logger.Debugf("🔧 ExchangeStore.Create: userID=%s, exchangeType=%s, accountName=%s, id=%s",
		userID, exchangeType, accountName, id)

	exchange := &Exchange{
		ID:           id,
		ExchangeType: exchangeType,
		AccountName:  accountName,
		UserID:       userID,
		Name:         name,
		Type:         typ,
		Enabled:      true,
		APIKey:       secrets.EncryptedString(apiKey),
		SecretKey:    secrets.EncryptedString(secretKey),
		Market:       MarketForExchangeType(exchangeType),
	}
	if err := s.db.Create(exchange).Error; err != nil {
		return "", err
	}
	return id, nil
}

// Update changes an account's credentials. An empty key keeps the stored value.
func (s *ExchangeStore) Update(userID, id string, apiKey, secretKey string) error {
	logger.Debugf("🔧 ExchangeStore.Update: userID=%s, id=%s", userID, id)

	updates := map[string]interface{}{
		"enabled":    true,
		"updated_at": time.Now().UTC(),
	}
	if apiKey != "" {
		updates["api_key"] = secrets.EncryptedString(apiKey)
	}
	if secretKey != "" {
		updates["secret_key"] = secrets.EncryptedString(secretKey)
	}

	result := s.db.Model(&Exchange{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("exchange not found: id=%s, userID=%s", id, userID)
	}
	return nil
}

// UpdateAccountName updates the account name for an exchange
func (s *ExchangeStore) UpdateAccountName(userID, id, accountName string) error {
	result := s.db.Model(&Exchange{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"account_name": accountName, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("exchange not found: id=%s, userID=%s", id, userID)
	}
	return nil
}

// Delete removes an exchange account
func (s *ExchangeStore) Delete(userID, id string) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Exchange{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("exchange not found: id=%s, userID=%s", id, userID)
	}
	return nil
}
