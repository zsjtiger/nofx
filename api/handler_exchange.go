package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"nofx/config"
	"nofx/crypto"
	"nofx/logger"
	"nofx/store"

	"github.com/gin-gonic/gin"
)

// ExchangeConfig is an exchange account as sent to the settings page.
type ExchangeConfig struct {
	ID           string `json:"id"`
	ExchangeType string `json:"exchange_type"` // "alpaca_paper" or "cn_paper"
	AccountName  string `json:"account_name"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	APIKey       string `json:"apiKey,omitempty"`
	SecretKey    string `json:"secretKey,omitempty"`
}

// SafeExchangeConfig is an exchange account without its secrets.
type SafeExchangeConfig struct {
	ID           string `json:"id"`
	ExchangeType string `json:"exchange_type"`
	AccountName  string `json:"account_name"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Market       string `json:"market"`
	Enabled      bool   `json:"enabled"`
	HasAPIKey    bool   `json:"has_api_key"`
	HasSecretKey bool   `json:"has_secret_key"`
}

func safeExchangeConfigFromStore(exchange *store.Exchange) SafeExchangeConfig {
	return SafeExchangeConfig{
		ID:           exchange.ID,
		ExchangeType: exchange.ExchangeType,
		AccountName:  exchange.AccountName,
		Name:         exchange.Name,
		Type:         exchange.Type,
		Market:       exchange.Market,
		Enabled:      exchange.Enabled,
		HasAPIKey:    exchange.APIKey != "",
		HasSecretKey: exchange.SecretKey != "",
	}
}

// ExchangeConfigUpdate is one exchange account's update payload. It is a named
// type so the log sanitizer in utils.go covers every sensitive field.
type ExchangeConfigUpdate struct {
	Enabled   bool   `json:"enabled"`
	APIKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

// UpdateExchangeConfigRequest updates one or more exchange accounts by UUID.
type UpdateExchangeConfigRequest struct {
	Exchanges map[string]ExchangeConfigUpdate `json:"exchanges"`
}

// CreateExchangeRequest creates an exchange account.
type CreateExchangeRequest struct {
	ExchangeType string `json:"exchange_type" binding:"required"`
	AccountName  string `json:"account_name"`
	APIKey       string `json:"api_key"`
	SecretKey    string `json:"secret_key"`
}

// decodeSecureBody reads a JSON body that is plain text or an encrypted payload,
// depending on whether transport encryption is enabled. It writes the error response and returns false on failure.
func (s *Server) decodeSecureBody(c *gin.Context, into interface{}) bool {
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return false
	}

	if !config.Get().TransportEncryption {
		if err := json.Unmarshal(bodyBytes, into); err != nil {
			logger.Infof("❌ Failed to parse plain JSON request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return false
		}
		return true
	}

	var encryptedPayload crypto.EncryptedPayload
	if err := json.Unmarshal(bodyBytes, &encryptedPayload); err != nil {
		logger.Infof("❌ Failed to parse encrypted payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format, encrypted transmission required"})
		return false
	}
	if encryptedPayload.WrappedKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "This endpoint only supports encrypted transmission, please use encrypted client",
			"code":    "ENCRYPTION_REQUIRED",
			"message": "Encrypted transmission is required for security reasons",
		})
		return false
	}
	decrypted, err := s.cryptoHandler.cryptoService.DecryptSensitiveData(&encryptedPayload)
	if err != nil {
		logger.Infof("❌ Failed to decrypt request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decrypt data"})
		return false
	}
	if err := json.Unmarshal([]byte(decrypted), into); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse decrypted data"})
		return false
	}
	return true
}

func (s *Server) handleGetExchangeConfigs(c *gin.Context) {
	userID := c.GetString("user_id")
	exchanges, err := s.store.Exchange().List(userID)
	if err != nil {
		SafeInternalError(c, "Failed to get exchange configs", err)
		return
	}

	safeExchanges := make([]SafeExchangeConfig, 0, len(exchanges))
	for _, exchange := range exchanges {
		if !store.IsVisibleExchange(exchange) {
			continue
		}
		safeExchanges = append(safeExchanges, safeExchangeConfigFromStore(exchange))
	}
	c.JSON(http.StatusOK, safeExchanges)
}

func (s *Server) handleUpdateExchangeConfigs(c *gin.Context) {
	userID := c.GetString("user_id")
	var req UpdateExchangeConfigRequest
	if !s.decodeSecureBody(c, &req) {
		return
	}

	tradersToReload := make(map[string]bool)
	for exchangeID, exchangeData := range req.Exchanges {
		existing, err := s.store.Exchange().GetByID(userID, exchangeID)
		if err != nil {
			SafeInternalError(c, fmt.Sprintf("Load exchange %s", exchangeID), err)
			return
		}

		effectiveAPIKey := strings.TrimSpace(exchangeData.APIKey)
		if effectiveAPIKey == "" {
			effectiveAPIKey = strings.TrimSpace(string(existing.APIKey))
		}
		effectiveSecretKey := strings.TrimSpace(exchangeData.SecretKey)
		if effectiveSecretKey == "" {
			effectiveSecretKey = strings.TrimSpace(string(existing.SecretKey))
		}
		if missing := store.MissingRequiredExchangeCredentialFields(existing.ExchangeType, effectiveAPIKey, effectiveSecretKey); len(missing) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          fmt.Sprintf("Missing required exchange fields: %s", strings.Join(missing, ", ")),
				"missing_fields": missing,
			})
			return
		}

		traders, _ := s.store.Trader().ListByExchangeID(userID, exchangeID)
		for _, t := range traders {
			tradersToReload[t.ID] = true
		}

		if err := s.store.Exchange().Update(userID, exchangeID, exchangeData.APIKey, exchangeData.SecretKey); err != nil {
			SafeInternalError(c, fmt.Sprintf("Update exchange %s", exchangeID), err)
			return
		}
	}
	s.exchangeAccountStateCache.Invalidate(userID)

	for traderID := range tradersToReload {
		logger.Infof("🔄 Removing trader %s from memory to reload with new exchange config", traderID)
		s.traderManager.RemoveTrader(traderID)
	}
	if err := s.traderManager.LoadUserTradersFromStore(s.store, userID); err != nil {
		logger.Infof("⚠️ Failed to reload user traders into memory: %v", err)
	}

	logger.Infof("✓ Exchange config updated: %+v", SanitizeExchangeConfigForLog(req.Exchanges))
	c.JSON(http.StatusOK, gin.H{"message": "Exchange configuration updated"})
}

func (s *Server) handleCreateExchange(c *gin.Context) {
	userID := c.GetString("user_id")
	var req CreateExchangeRequest
	if !s.decodeSecureBody(c, &req) {
		return
	}

	if store.MarketForExchangeType(req.ExchangeType) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid exchange type: %s", req.ExchangeType)})
		return
	}
	if missing := store.MissingRequiredExchangeCredentialFields(req.ExchangeType, req.APIKey, req.SecretKey); len(missing) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          fmt.Sprintf("Missing required exchange fields: %s", strings.Join(missing, ", ")),
			"missing_fields": missing,
		})
		return
	}

	id, err := s.store.Exchange().Create(userID, req.ExchangeType, req.AccountName, true, req.APIKey, req.SecretKey)
	if err != nil {
		logger.Infof("❌ Failed to create exchange account: %v", err)
		SafeInternalError(c, "Failed to create exchange account", err)
		return
	}
	s.exchangeAccountStateCache.Invalidate(userID)
	logger.Infof("✓ Created exchange account: type=%s, name=%s, id=%s", req.ExchangeType, req.AccountName, id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Exchange account created",
		"id":      id,
	})
}

func (s *Server) handleDeleteExchange(c *gin.Context) {
	userID := c.GetString("user_id")
	exchangeID := c.Param("id")
	if exchangeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exchange ID is required"})
		return
	}

	traders, err := s.store.Trader().ListByExchangeID(userID, exchangeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check traders"})
		return
	}
	if len(traders) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot delete: %d trader(s) still use this exchange account. Delete them first.", len(traders)),
		})
		return
	}

	if err := s.store.Exchange().Delete(userID, exchangeID); err != nil {
		SafeInternalError(c, "Failed to delete exchange account", err)
		return
	}
	s.exchangeAccountStateCache.Invalidate(userID)
	c.JSON(http.StatusOK, gin.H{"message": "Exchange account deleted"})
}

// handleGetSupportedExchanges lists the exchange types a user can connect. They are templates, not accounts.
func (s *Server) handleGetSupportedExchanges(c *gin.Context) {
	c.JSON(http.StatusOK, []ExchangeConfig{
		{ExchangeType: "alpaca_paper", Name: "Alpaca Paper (US Equities)", Type: "paper"},
		{ExchangeType: "cn_paper", Name: "A-Share Paper (CN)", Type: "paper"},
	})
}
