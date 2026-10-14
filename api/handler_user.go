package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"nofx/auth"
	"nofx/logger"
	"nofx/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// handleLogout Add current token to blacklist
func (s *Server) handleLogout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
		return
	}
	tokenString := parts[1]
	claims, err := auth.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	var exp time.Time
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Time
	} else {
		exp = time.Now().Add(24 * time.Hour)
	}
	auth.BlacklistToken(tokenString, exp)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// handleRegister Handle user registration request.
// handleRegister allows registration only when no users exist yet (first-time setup).
// This is a single-user system; subsequent registrations are permanently closed.
func (s *Server) handleRegister(c *gin.Context) {
	userCount, err := s.store.User().Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user count"})
		return
	}

	if userCount > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "System already initialized"})
		return
	}

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Lang     string `json:"lang"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		SafeBadRequest(c, "Invalid request parameters")
		return
	}

	lang := req.Lang
	if lang != "zh" && lang != "id" {
		lang = "en"
	}

	// Check if email already exists
	_, err = s.store.User().GetByEmail(req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Generate password hash
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password processing failed"})
		return
	}

	// Create user
	userID := uuid.New().String()
	user := &store.User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	err = s.store.User().Create(user)
	if err != nil {
		SafeInternalError(c, "Failed to create user", err)
		return
	}

	// NOTE: Orphan record adoption was removed for security reasons. Previously,
	// after a reset-account call, any new user would inherit the prior owner's
	// wallet keys and exchange API credentials — a catastrophic IDOR/takeover
	// path. Operators who need to migrate credentials across users must do so
	// explicitly via export/import, never via implicit adoption on registration.

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Initialize default model and exchange configs for user
	err = s.initUserDefaultConfigs(user.ID, lang)
	if err != nil {
		logger.Infof("Failed to initialize user default configs: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user_id": user.ID,
		"email":   user.Email,
		"message": "Registration successful",
	})
}

// dummyPasswordHash is a valid bcrypt hash of a throwaway value. It is compared
// against when the submitted email does not exist so that login takes roughly
// the same time whether or not the account exists — closing the timing side
// channel that would otherwise let an attacker enumerate valid emails (a fast
// "no such user" vs. a slow bcrypt compare). It is not a secret.
const dummyPasswordHash = "$2a$10$0iF0bCoQLJ6Ph1bF.MXwHOW.IMTxQjeEW.w38dctRQAB2kwB6ga1q"

// handleLogin Handle user login request
func (s *Server) handleLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		SafeBadRequest(c, "Invalid request parameters")
		return
	}

	// Get user information
	user, err := s.store.User().GetByEmail(req.Email)
	if err != nil {
		// Perform a dummy comparison so the response time does not reveal
		// whether the email exists (anti user-enumeration), then fail uniformly.
		auth.CheckPassword(req.Password, dummyPasswordHash)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password incorrect"})
		return
	}

	// Verify password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password incorrect"})
		return
	}

	// Issue token directly after password verification.
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user_id": user.ID,
		"email":   user.Email,
		"message": "Login successful",
	})
}

// handleChangePassword changes the password for the currently authenticated user.
func (s *Server) handleChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		SafeBadRequest(c, "new_password is required (min 8 chars)")
		return
	}
	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		SafeInternalError(c, "Password processing failed", err)
		return
	}
	if err := s.store.User().UpdatePassword(userID, hash); err != nil {
		SafeInternalError(c, "Failed to update password", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

// NOTE: Password and account recovery used to live here as the public,
// unauthenticated handlers handleResetPassword / handleResetAccount. They were
// removed because an unauthenticated recovery endpoint is a remotely
// exploitable auth-bypass on any public-facing deployment: the confirm phrase
// was embedded in the frontend (and echoed back by the API), so it was friction
// rather than authentication. Recovery now lives in the local CLI
// (`nofx reset-password` / `nofx reset-account`, see cli.go), which requires
// shell access to the host — something a remote attacker does not have.

// initUserDefaultConfigs Initialize default configs for new user
func (s *Server) initUserDefaultConfigs(userID string, lang string) error {
	if err := s.createDefaultStrategies(userID, lang); err != nil {
		logger.Warnf("Failed to create default strategies for user %s: %v", userID, err)
		// Non-fatal: user can create strategies manually
	}
	logger.Infof("✓ User %s registration completed with default strategies", userID)
	return nil
}

func (s *Server) createDefaultStrategies(userID string, lang string) error {
	type strategyDef struct {
		name        string
		description string
		isActive    bool
		applyConfig func(*store.StrategyConfig)
	}

	// Two presets, one per stock market. Each lists liquid large caps, so the
	// symbols are valid on the matching paper account.
	setEquityStrategy := func(symbols []string) func(*store.StrategyConfig) {
		return func(c *store.StrategyConfig) {
			c.CoinSource.SourceType = "static"
			c.CoinSource.StaticCoins = append([]string(nil), symbols...)
			c.CoinSource.UseAI500 = false
			c.CoinSource.UseOITop = false
			c.CoinSource.UseOILow = false
			c.RiskControl.MaxPositions = 3
			c.RiskControl.BTCETHMaxLeverage = 1
			c.RiskControl.AltcoinMaxLeverage = 1
			c.RiskControl.BTCETHMaxPositionValueRatio = 1.0
			c.RiskControl.AltcoinMaxPositionValueRatio = 1.0
			c.RiskControl.MaxMarginUsage = 1.0
			c.RiskControl.MinConfidence = 70
			c.RiskControl.MinRiskRewardRatio = 2.0
			c.Indicators.Klines.PrimaryTimeframe = "1d"
			c.Indicators.Klines.PrimaryCount = 60
			c.Indicators.Klines.LongerTimeframe = ""
			c.Indicators.Klines.LongerCount = 0
			c.Indicators.Klines.EnableMultiTimeframe = false
			c.Indicators.Klines.SelectedTimeframes = []string{"1d"}
			c.Indicators.EnableRawKlines = true
		}
	}

	usStrategy := strategyDef{
		name:        "NOFX US Stock Strategy",
		description: "Daily-bar AI strategy over large-cap US stocks on the Alpaca paper account. Long and short, no leverage.",
		isActive:    true,
		applyConfig: setEquityStrategy([]string{"AAPL", "MSFT", "NVDA", "AMZN", "GOOGL", "META", "TSLA"}),
	}
	cnStrategy := strategyDef{
		name:        "NOFX A-Share Strategy",
		description: "Daily-bar AI strategy over large-cap A-shares on the paper account. Long only, whole lots, T+1.",
		isActive:    false,
		applyConfig: setEquityStrategy([]string{"600519.SH", "000858.SZ", "601318.SH", "000333.SZ", "300750.SZ"}),
	}
	if lang == "zh" {
		usStrategy.name = "NOFX 美股策略"
		usStrategy.description = "基于日线的 AI 策略，覆盖纳斯达克与纽交所大盘股，使用 Alpaca 纸上账户，可做多做空，不加杠杆。"
		cnStrategy.name = "NOFX A 股策略"
		cnStrategy.description = "基于日线的 AI 策略，覆盖沪深大盘股，使用纸上账户，只做多，按整手交易，遵守 T+1。"
	}
	if lang == "id" {
		usStrategy.name = "Strategi Saham AS NOFX"
		usStrategy.description = "Strategi AI berbasis data harian untuk saham AS berkapitalisasi besar di akun paper Alpaca. Long dan short, tanpa leverage."
		cnStrategy.name = "Strategi Saham A-Share NOFX"
		cnStrategy.description = "Strategi AI berbasis data harian untuk saham A-share berkapitalisasi besar di akun paper. Hanya long, per lot 100 saham, aturan T+1."
	}

	definitions := []strategyDef{usStrategy, cnStrategy}

	// GetDefaultStrategyConfig only supports zh/en; map id -> en
	configLang := lang
	if lang == "id" {
		configLang = "en"
	}

	// Pre-build all strategy objects before opening the transaction
	var strategies []*store.Strategy
	for _, def := range definitions {
		config := store.GetDefaultStrategyConfig(configLang)
		def.applyConfig(&config)
		config.ClampLimits()

		strategy := &store.Strategy{
			ID:          uuid.New().String(),
			UserID:      userID,
			Name:        def.name,
			Description: def.description,
			IsActive:    def.isActive,
			IsDefault:   false,
		}
		if err := strategy.SetConfig(&config); err != nil {
			return fmt.Errorf("failed to set config for strategy %q: %w", def.name, err)
		}
		strategies = append(strategies, strategy)
	}

	legacyDefaultNames := []string{
		"NOFX Claw402 Auto Strategy", "Strategi Otomatis NOFX Claw402",
		"Balanced Strategy", "Steady Strategy", "Aggressive Strategy",
		"US Stock Trend Strategy", "US Stock Steady Strategy", "US Stock Breakout Strategy",
		"Balanced Strategy", "Conservative Strategy", "Aggressive Strategy",
		"US Stock Trend Strategy", "US Stock Steady Strategy", "US Stock Breakout Strategy",
		"Strategi Seimbang", "Strategi Konservatif", "Strategi Agresif",
		"Strategi Tren Saham AS", "Strategi Stabil Saham AS", "Strategi Breakout Saham AS",
	}

	return s.store.Transaction(func(tx *gorm.DB) error {
		// Remove obsolete built-in risk-profile presets for this user. If a trader still
		// references one of them, keep it to avoid breaking an existing running setup.
		deleteResult := tx.Where("user_id = ? AND name IN ? AND id NOT IN (SELECT strategy_id FROM traders WHERE user_id = ? AND strategy_id IS NOT NULL)", userID, legacyDefaultNames, userID).
			Delete(&store.Strategy{})
		if deleteResult.Error != nil {
			return fmt.Errorf("failed to remove legacy default strategies: %w", deleteResult.Error)
		}
		if deleteResult.RowsAffected > 0 {
			logger.Infof("  ✓ Removed %d legacy default strategy preset(s)", deleteResult.RowsAffected)
		}

		var activeCount int64
		if err := tx.Model(&store.Strategy{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&activeCount).Error; err != nil {
			return fmt.Errorf("failed to count active strategies: %w", err)
		}

		for _, strategy := range strategies {
			var existing store.Strategy
			query := tx.Where("user_id = ? AND name = ?", userID, strategy.Name).First(&existing)
			if query.Error == nil {
				config, err := existing.ParseConfig()
				if err != nil {
					return fmt.Errorf("failed to parse existing strategy %q: %w", strategy.Name, err)
				}
				if store.MigrateLegacyAutopilotRiskDefaults(config) {
					config.ClampLimits()
					if err := existing.SetConfig(config); err != nil {
						return fmt.Errorf("failed to serialize migrated strategy %q: %w", strategy.Name, err)
					}
					if err := tx.Model(&store.Strategy{}).
						Where("id = ? AND user_id = ?", existing.ID, userID).
						Updates(map[string]interface{}{"config": existing.Config, "updated_at": time.Now().UTC()}).Error; err != nil {
						return fmt.Errorf("failed to migrate strategy %q: %w", strategy.Name, err)
					}
					logger.Infof("  ✓ Migrated default strategy to eight-position Autopilot: %s", strategy.Name)
				}
				continue
			}
			if query.Error != gorm.ErrRecordNotFound {
				return fmt.Errorf("failed to check strategy %q: %w", strategy.Name, query.Error)
			}
			if activeCount > 0 {
				strategy.IsActive = false
			}
			if err := tx.Create(strategy).Error; err != nil {
				return fmt.Errorf("failed to create strategy %q: %w", strategy.Name, err)
			}
			if strategy.IsActive {
				activeCount++
			}
			logger.Infof("  ✓ Created default strategy: %s (active=%v)", strategy.Name, strategy.IsActive)
		}
		return nil
	})
}
