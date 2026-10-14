package store

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Hard limits that keep AI requests a sensible size.
const (
	MaxCandidateSymbols = 20
	MaxPositions        = 8
	MaxTimeframes       = 2
	MinKlineCount       = 10
	MaxKlineCount       = 120
	MinPositionValue    = 1.0
	MaxPositionValue    = 10.0 // max single position value, as a multiple of equity
	MinConfidence       = 50
	MaxConfidence       = 100
)

// Defaults for the stock strategy.
const (
	DefaultMaxPositions     = 3
	DefaultPositionValueCap = 1.0 // a single position may use at most 1x equity
)

// StrategyConfig is the decision configuration of a strategy. It applies to one
// stock market; the market is set by the trader's exchange account.
type StrategyConfig struct {
	Language       string                 `json:"language,omitempty"`
	CoinSource     CoinSourceConfig       `json:"-"`
	Indicators     IndicatorConfig        `json:"-"`
	CustomPrompt   string                 `json:"-"`
	RiskControl    RiskControlConfig      `json:"-"`
	PromptSections PromptSectionsConfig   `json:"-"`
	PublishConfig  *PublishStrategyConfig `json:"publish_config,omitempty"`
}

// AIStrategyConfig is the nested "ai_config" block of the stored JSON.
type AIStrategyConfig struct {
	CoinSource     CoinSourceConfig     `json:"coin_source"`
	Indicators     IndicatorConfig      `json:"indicators"`
	CustomPrompt   string               `json:"custom_prompt,omitempty"`
	RiskControl    RiskControlConfig    `json:"risk_control"`
	PromptSections PromptSectionsConfig `json:"prompt_sections,omitempty"`
}

// PublishStrategyConfig controls whether a strategy is listed publicly.
type PublishStrategyConfig struct {
	IsPublic      bool `json:"is_public"`
	ConfigVisible bool `json:"config_visible"`
}

// MarshalJSON writes the nested schema used by the API and the database.
func (c StrategyConfig) MarshalJSON() ([]byte, error) {
	out := struct {
		Language      string                 `json:"language,omitempty"`
		AIConfig      AIStrategyConfig       `json:"ai_config"`
		PublishConfig *PublishStrategyConfig `json:"publish_config,omitempty"`
	}{
		Language: c.Language,
		AIConfig: AIStrategyConfig{
			CoinSource:     c.CoinSource,
			Indicators:     c.Indicators,
			CustomPrompt:   c.CustomPrompt,
			RiskControl:    c.RiskControl,
			PromptSections: c.PromptSections,
		},
		PublishConfig: c.PublishConfig,
	}
	return json.Marshal(out)
}

// UnmarshalJSON accepts the nested schema and the older flat layout. Fields
// from retired strategy types are ignored.
func (c *StrategyConfig) UnmarshalJSON(data []byte) error {
	type rawStrategyConfig struct {
		Language       string                 `json:"language"`
		AIConfig       *AIStrategyConfig      `json:"ai_config"`
		PublishConfig  *PublishStrategyConfig `json:"publish_config"`
		CoinSource     *CoinSourceConfig      `json:"coin_source"`
		Indicators     *IndicatorConfig       `json:"indicators"`
		CustomPrompt   *string                `json:"custom_prompt"`
		RiskControl    *RiskControlConfig     `json:"risk_control"`
		PromptSections *PromptSectionsConfig  `json:"prompt_sections"`
	}
	var raw rawStrategyConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.Language = raw.Language
	c.PublishConfig = raw.PublishConfig
	if raw.AIConfig != nil {
		c.CoinSource = raw.AIConfig.CoinSource
		c.Indicators = raw.AIConfig.Indicators
		c.CustomPrompt = raw.AIConfig.CustomPrompt
		c.RiskControl = raw.AIConfig.RiskControl
		c.PromptSections = raw.AIConfig.PromptSections
		return nil
	}
	if raw.CoinSource != nil {
		c.CoinSource = *raw.CoinSource
	}
	if raw.Indicators != nil {
		c.Indicators = *raw.Indicators
	}
	if raw.CustomPrompt != nil {
		c.CustomPrompt = *raw.CustomPrompt
	}
	if raw.RiskControl != nil {
		c.RiskControl = *raw.RiskControl
	}
	if raw.PromptSections != nil {
		c.PromptSections = *raw.PromptSections
	}
	return nil
}

// CoinSourceConfig is the stock universe. Only a static symbol list is supported.
type CoinSourceConfig struct {
	SourceType    string   `json:"source_type"`
	StaticCoins   []string `json:"static_coins,omitempty"`
	ExcludedCoins []string `json:"excluded_coins,omitempty"`
}

// IndicatorConfig selects what the AI sees for each symbol.
type IndicatorConfig struct {
	Klines              KlineConfig          `json:"klines"`
	EnableRawKlines     bool                 `json:"enable_raw_klines"`
	EnableEMA           bool                 `json:"enable_ema"`
	EnableMACD          bool                 `json:"enable_macd"`
	EnableRSI           bool                 `json:"enable_rsi"`
	EnableATR           bool                 `json:"enable_atr"`
	EnableBOLL          bool                 `json:"enable_boll"`
	EnableVolume        bool                 `json:"enable_volume"`
	EMAPeriods          []int                `json:"ema_periods,omitempty"`
	RSIPeriods          []int                `json:"rsi_periods,omitempty"`
	ATRPeriods          []int                `json:"atr_periods,omitempty"`
	BOLLPeriods         []int                `json:"boll_periods,omitempty"`
	ExternalDataSources []ExternalDataSource `json:"external_data_sources,omitempty"`
}

// KlineConfig selects the bars per symbol. Stock data is daily, so the default is 1d.
type KlineConfig struct {
	PrimaryTimeframe     string   `json:"primary_timeframe"`
	PrimaryCount         int      `json:"primary_count"`
	LongerTimeframe      string   `json:"longer_timeframe,omitempty"`
	LongerCount          int      `json:"longer_count,omitempty"`
	EnableMultiTimeframe bool     `json:"enable_multi_timeframe"`
	SelectedTimeframes   []string `json:"selected_timeframes,omitempty"`
}

// ExternalDataSource is a user-supplied HTTP JSON feed added to the prompt.
type ExternalDataSource struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers,omitempty"`
	DataPath    string            `json:"data_path,omitempty"`
	RefreshSecs int               `json:"refresh_secs,omitempty"`
}

// RiskControlConfig holds the limits enforced in code. Stock markets take no leverage.
type RiskControlConfig struct {
	MaxPositions                 int     `json:"max_positions"`
	AltcoinMaxPositionValueRatio float64 `json:"altcoin_max_position_value_ratio"` // single position value cap, as a multiple of equity
	MaxMarginUsage               float64 `json:"max_margin_usage"`
	MinPositionSize              float64 `json:"min_position_size"`
	MinRiskRewardRatio           float64 `json:"min_risk_reward_ratio"`
	MinConfidence                int     `json:"min_confidence"`
}

// PromptSectionsConfig overrides sections of the built-in prompt.
type PromptSectionsConfig struct {
	RoleDefinition   string `json:"role_definition,omitempty"`
	TradingFrequency string `json:"trading_frequency,omitempty"`
	EntryStandards   string `json:"entry_standards,omitempty"`
	DecisionProcess  string `json:"decision_process,omitempty"`
}

// ClampLimits brings numeric settings into the ranges the decision engine supports.
func (c *StrategyConfig) ClampLimits() {
	rc := &c.RiskControl
	rc.MaxPositions = clampInt(rc.MaxPositions, 1, MaxPositions, DefaultMaxPositions)
	rc.AltcoinMaxPositionValueRatio = clampFloat(rc.AltcoinMaxPositionValueRatio, 0.1, MaxPositionValue, DefaultPositionValueCap)
	rc.MaxMarginUsage = clampFloat(rc.MaxMarginUsage, 0.1, 1.0, 1.0)
	rc.MinRiskRewardRatio = clampFloat(rc.MinRiskRewardRatio, 1.0, 10.0, 2.0)
	rc.MinConfidence = clampInt(rc.MinConfidence, MinConfidence, MaxConfidence, 70)

	kl := &c.Indicators.Klines
	kl.PrimaryCount = clampInt(kl.PrimaryCount, MinKlineCount, MaxKlineCount, 60)
	kl.LongerCount = clampInt(kl.LongerCount, 0, MaxKlineCount, 0)
	if len(kl.SelectedTimeframes) > MaxTimeframes {
		kl.SelectedTimeframes = kl.SelectedTimeframes[:MaxTimeframes]
	}
	if len(c.CoinSource.StaticCoins) > MaxCandidateSymbols {
		c.CoinSource.StaticCoins = c.CoinSource.StaticCoins[:MaxCandidateSymbols]
	}
}

// NormalizeProductSchema fills defaults and canonicalises the symbols of a stored or submitted config.
// Symbols are normalised per market at trade time; here only the source type and blanks are fixed.
func (c *StrategyConfig) NormalizeProductSchema() {
	c.CoinSource.SourceType = "static"
	c.CoinSource.StaticCoins = normalizeSymbols(c.CoinSource.StaticCoins)
	c.CoinSource.ExcludedCoins = normalizeSymbols(c.CoinSource.ExcludedCoins)
	if strings.TrimSpace(c.Indicators.Klines.PrimaryTimeframe) == "" {
		c.Indicators.Klines.PrimaryTimeframe = "1d"
	}
	c.Indicators.Klines.SelectedTimeframes = normalizeTimeframes(c.Indicators.Klines.SelectedTimeframes)
	if len(c.Indicators.Klines.SelectedTimeframes) == 0 {
		c.Indicators.Klines.SelectedTimeframes = []string{c.Indicators.Klines.PrimaryTimeframe}
	}
	if c.Language != "zh" && c.Language != "en" {
		c.Language = "en"
	}
}

// normalizeSymbols trims, upper-cases and removes duplicates and blanks.
func normalizeSymbols(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		s := strings.ToUpper(strings.TrimSpace(v))
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

// normalizeTimeframes keeps the supported stock timeframes in their given order.
func normalizeTimeframes(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]bool, len(values))
	for _, v := range values {
		tf := strings.ToLower(strings.TrimSpace(v))
		if tf == "" || seen[tf] {
			continue
		}
		seen[tf] = true
		out = append(out, tf)
	}
	return out
}

func clampInt(v, lo, hi, fallback int) int {
	if v == 0 {
		return fallback
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func clampFloat(v, lo, hi, fallback float64) float64 {
	if v == 0 {
		return fallback
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// MergeStrategyConfig applies a partial JSON patch on top of a base config.
// Nested objects merge key by key; other values are replaced.
func MergeStrategyConfig(base StrategyConfig, patch map[string]any) (StrategyConfig, error) {
	raw, err := json.Marshal(base)
	if err != nil {
		return base, err
	}
	var merged map[string]any
	if err := json.Unmarshal(raw, &merged); err != nil {
		return base, err
	}
	mergeJSONMaps(merged, patch)
	out, err := json.Marshal(merged)
	if err != nil {
		return base, err
	}
	var result StrategyConfig
	if err := json.Unmarshal(out, &result); err != nil {
		return base, err
	}
	return result, nil
}

// GetDefaultStrategyConfig returns the starting configuration for a new stock strategy.
func GetDefaultStrategyConfig(lang string) StrategyConfig {
	normalizedLang := "en"
	if lang == "zh" {
		normalizedLang = "zh"
	}

	config := StrategyConfig{
		Language: normalizedLang,
		CoinSource: CoinSourceConfig{
			SourceType:  "static",
			StaticCoins: []string{"AAPL", "MSFT", "NVDA"},
		},
		Indicators: IndicatorConfig{
			Klines: KlineConfig{
				PrimaryTimeframe:   "1d",
				PrimaryCount:       60,
				SelectedTimeframes: []string{"1d"},
			},
			EnableRawKlines: true,
			EMAPeriods:      []int{20, 50},
			RSIPeriods:      []int{14},
			ATRPeriods:      []int{14},
			BOLLPeriods:     []int{20},
		},
		RiskControl: RiskControlConfig{
			MaxPositions:                 DefaultMaxPositions,
			AltcoinMaxPositionValueRatio: DefaultPositionValueCap,
			MaxMarginUsage:               1.0,
			MinPositionSize:              MinPositionValue,
			MinRiskRewardRatio:           2.0,
			MinConfidence:                70,
		},
	}

	if normalizedLang == "zh" {
		config.PromptSections = PromptSectionsConfig{
			RoleDefinition:   "# 你是 NOFX 股票自动交易员\n\n根据每个标的的日线数据与风控规则做决策。只交易候选列表中的股票。",
			TradingFrequency: "# 交易频率\n\n- 日线级别决策，不追求日内短线。\n- 趋势未改变时继续持有。",
			EntryStandards:   "# 入场标准\n\n只在趋势与量能配合、且止损位置清晰时开仓。置信度低于门槛时观望。",
			DecisionProcess:  "# 决策流程\n\n1. 查看账户与持仓。\n2. 对已有持仓判断是否止盈、止损或继续持有。\n3. 对候选标的判断是否开仓。\n4. 严格按 JSON 输出。",
		}
	} else {
		config.PromptSections = PromptSectionsConfig{
			RoleDefinition:   "# You are the NOFX stock auto-trader\n\nDecide from each symbol's daily bars and the risk rules. Trade only the candidate symbols.",
			TradingFrequency: "# Trading Frequency\n\n- Decide on daily bars; this is not an intraday strategy.\n- Keep holding while the trend is intact.",
			EntryStandards:   "# Entry Standards\n\nOpen only when trend and volume agree and the stop-loss level is clear. Wait when confidence is below the minimum.",
			DecisionProcess:  "# Decision Process\n\n1. Review the account and open positions.\n2. For each position, decide to take profit, cut the loss, or hold.\n3. For each candidate, decide whether to open.\n4. Output strict JSON.",
		}
	}

	return config
}

// EstimateTokens estimates the token count of a strategy's prompts for each model
// context limit. It is a pure computation based on config fields and makes no network calls.
func (c *StrategyConfig) EstimateTokens() TokenEstimate {
	breakdown := TokenBreakdown{}

	baseChars := 3000
	if c.Language == "zh" {
		baseChars = 2400
	}
	baseChars += len(c.PromptSections.RoleDefinition) + len(c.PromptSections.TradingFrequency) +
		len(c.PromptSections.EntryStandards) + len(c.PromptSections.DecisionProcess) + len(c.CustomPrompt)
	if c.Language == "zh" {
		breakdown.SystemPrompt = baseChars / 2
	} else {
		breakdown.SystemPrompt = baseChars / 4
	}
	breakdown.FixedOverhead = 200

	symbols := len(c.CoinSource.StaticCoins)
	timeframes := len(c.Indicators.Klines.SelectedTimeframes)
	if timeframes == 0 {
		timeframes = 1
	}
	// Each bar is roughly 60 characters of text once formatted.
	barTokens := c.Indicators.Klines.PrimaryCount * 60 / 4
	breakdown.MarketData = symbols * timeframes * barTokens
	if breakdown.MarketData == 0 {
		breakdown.MarketData = barTokens
	}

	total := breakdown.SystemPrompt + breakdown.MarketData + breakdown.FixedOverhead + breakdown.RankingData + breakdown.QuantData

	var limits []ModelLimit
	for _, name := range []string{"deepseek", "openai", "claude", "qwen", "gemini", "grok", "kimi", "minimax"} {
		limit := GetContextLimit(name)
		usage := 0
		if limit > 0 {
			usage = total * 100 / limit
		}
		level := "ok"
		switch {
		case usage >= 90:
			level = "danger"
		case usage >= 70:
			level = "warning"
		}
		limits = append(limits, ModelLimit{Name: name, ContextLimit: limit, UsagePct: usage, Level: level})
	}

	var suggestions []string
	if total > GetContextLimit("deepseek") {
		suggestions = append(suggestions, "Reduce the number of symbols or the bar count.")
	}
	return TokenEstimate{Total: total, Breakdown: breakdown, ModelLimits: limits, Suggestions: suggestions}
}

// StrategyStore strategy storage
type StrategyStore struct {
	db *gorm.DB
}

// Strategy strategy configuration
type Strategy struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	UserID        string    `gorm:"column:user_id;not null;default:'';index" json:"user_id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `gorm:"default:''" json:"description"`
	IsActive      bool      `gorm:"column:is_active;default:false;index" json:"is_active"`
	IsDefault     bool      `gorm:"column:is_default;default:false" json:"is_default"`
	IsPublic      bool      `gorm:"column:is_public;default:false;index" json:"is_public"`    // whether visible in strategy market
	ConfigVisible bool      `gorm:"column:config_visible;default:true" json:"config_visible"` // whether config details are visible
	Config        string    `gorm:"not null;default:'{}'" json:"config"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Strategy) TableName() string { return "strategies" }

func NewStrategyStore(db *gorm.DB) *StrategyStore {
	return &StrategyStore{db: db}
}

func (s *StrategyStore) initTables() error {
	// AutoMigrate will add missing columns without dropping existing data
	return s.db.AutoMigrate(&Strategy{})
}

func (s *StrategyStore) initDefaultData() error {
	// No longer pre-populate strategies - create on demand when user configures
	return nil
}

// GetDefaultStrategyConfig returns the default strategy configuration for the given language

// Create create a strategy
func (s *StrategyStore) Create(strategy *Strategy) error {
	return s.db.Create(strategy).Error
}

// Update update a strategy
func (s *StrategyStore) Update(strategy *Strategy) error {
	return s.db.Model(&Strategy{}).
		Where("id = ? AND user_id = ?", strategy.ID, strategy.UserID).
		Updates(map[string]interface{}{
			"name":           strategy.Name,
			"description":    strategy.Description,
			"config":         strategy.Config,
			"is_public":      strategy.IsPublic,
			"config_visible": strategy.ConfigVisible,
			"updated_at":     time.Now().UTC(),
		}).Error
}

// Delete delete a strategy
func (s *StrategyStore) Delete(userID, id string) error {
	// do not allow deleting system default strategy
	var st Strategy
	if err := s.db.Where("id = ?", id).First(&st).Error; err == nil {
		if st.IsDefault {
			return fmt.Errorf("cannot delete system default strategy")
		}
		if st.IsActive {
			return fmt.Errorf("cannot delete active strategy")
		}
	}

	// Check if any trader references this strategy
	var count int64
	if err := s.db.Model(&Trader{}).
		Where("user_id = ? AND strategy_id = ?", userID, id).
		Count(&count).Error; err == nil && count > 0 {
		return fmt.Errorf("cannot delete strategy in use by %d trader(s) - reassign those traders first", count)
	}

	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Strategy{}).Error
}

// List get user's strategy list
func (s *StrategyStore) List(userID string) ([]*Strategy, error) {
	var strategies []*Strategy
	err := s.db.Where("user_id = ? OR is_default = ?", userID, true).
		Order("is_default DESC, created_at DESC").
		Find(&strategies).Error
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

// ListPublic get all public strategies for the strategy market
func (s *StrategyStore) ListPublic() ([]*Strategy, error) {
	var strategies []*Strategy
	err := s.db.Where("is_public = ?", true).
		Order("created_at DESC").
		Find(&strategies).Error
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

// Get get a single strategy
func (s *StrategyStore) Get(userID, id string) (*Strategy, error) {
	var st Strategy
	err := s.db.Where("id = ? AND (user_id = ? OR is_default = ?)", id, userID, true).
		First(&st).Error
	if err != nil {
		return nil, err
	}
	return &st, nil
}

// GetActive get user's currently active strategy
func (s *StrategyStore) GetActive(userID string) (*Strategy, error) {
	var st Strategy
	err := s.db.Where("user_id = ? AND is_active = ?", userID, true).First(&st).Error
	if err == gorm.ErrRecordNotFound {
		// no active strategy, return system default strategy
		return s.GetDefault()
	}
	if err != nil {
		return nil, err
	}
	return &st, nil
}

// GetDefault get system default strategy
func (s *StrategyStore) GetDefault() (*Strategy, error) {
	var st Strategy
	err := s.db.Where("is_default = ?", true).First(&st).Error
	if err != nil {
		return nil, err
	}
	return &st, nil
}

// SetActive set active strategy (will first deactivate other strategies)
func (s *StrategyStore) SetActive(userID, strategyID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// first deactivate all strategies for the user
		if err := tx.Model(&Strategy{}).Where("user_id = ?", userID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// activate specified strategy
		return tx.Model(&Strategy{}).
			Where("id = ? AND (user_id = ? OR is_default = ?)", strategyID, userID, true).
			Update("is_active", true).Error
	})
}

// Duplicate duplicate a strategy (used to create custom strategy based on default strategy)
func (s *StrategyStore) Duplicate(userID, sourceID, newID, newName string) error {
	// get source strategy
	source, err := s.Get(userID, sourceID)
	if err != nil {
		return fmt.Errorf("failed to get source strategy: %w", err)
	}

	// create new strategy
	newStrategy := &Strategy{
		ID:          newID,
		UserID:      userID,
		Name:        newName,
		Description: "Created based on [" + source.Name + "]",
		IsActive:    false,
		IsDefault:   false,
		Config:      source.Config,
	}

	return s.Create(newStrategy)
}

// ParseConfig parse strategy configuration JSON
func (s *Strategy) ParseConfig() (*StrategyConfig, error) {
	var config StrategyConfig
	if err := json.Unmarshal([]byte(s.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to parse strategy configuration: %w", err)
	}
	return &config, nil
}

// SetConfig set strategy configuration
func (s *Strategy) SetConfig(config *StrategyConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to serialize strategy configuration: %w", err)
	}
	s.Config = string(data)
	return nil
}

// ============================================================================
// Token Estimation
// ============================================================================

// TokenEstimate holds the result of token estimation

type TokenEstimate struct {
	Total       int            `json:"total"`
	Breakdown   TokenBreakdown `json:"breakdown"`
	ModelLimits []ModelLimit   `json:"model_limits"`
	Suggestions []string       `json:"suggestions"`
}

// TokenBreakdown shows estimated tokens per component
type TokenBreakdown struct {
	SystemPrompt  int `json:"system_prompt"`
	MarketData    int `json:"market_data"`
	RankingData   int `json:"ranking_data"`
	QuantData     int `json:"quant_data"`
	FixedOverhead int `json:"fixed_overhead"`
}

// ModelLimit shows token usage against a specific model's context limit
type ModelLimit struct {
	Name         string `json:"name"`
	ContextLimit int    `json:"context_limit"`
	UsagePct     int    `json:"usage_pct"`
	Level        string `json:"level"` // "ok" | "warning" | "danger"
}

// Context window sizes (tokens) for each model family
const (
	contextLimitDeepSeek = 131_072   // 128K
	contextLimitOpenAI   = 128_000   // 128K
	contextLimitClaude   = 200_000   // 200K
	contextLimitQwen     = 131_072   // 128K
	contextLimitGemini   = 1_000_000 // 1M
	contextLimitGrok     = 131_072   // 128K
	contextLimitKimi     = 131_072   // 128K
	contextLimitMinimax  = 1_000_000 // 1M
)

// ModelContextLimits maps provider names to their context window sizes (in tokens)
var ModelContextLimits = map[string]int{
	"deepseek": contextLimitDeepSeek,
	"openai":   contextLimitOpenAI,
	"claude":   contextLimitClaude,
	"qwen":     contextLimitQwen,
	"gemini":   contextLimitGemini,
	"grok":     contextLimitGrok,
	"kimi":     contextLimitKimi,
	"minimax":  contextLimitMinimax,
}

// GetContextLimit returns the context limit for a given provider
func GetContextLimit(provider string) int {
	if limit, ok := ModelContextLimits[provider]; ok {
		return limit
	}
	return contextLimitDeepSeek // safe default
}

func mergeJSONMaps(dst, src map[string]any) {
	for key, srcVal := range src {
		srcMap, srcIsMap := srcVal.(map[string]any)
		dstMap, dstIsMap := dst[key].(map[string]any)
		if srcIsMap && dstIsMap {
			mergeJSONMaps(dstMap, srcMap)
			continue
		}
		dst[key] = srcVal
	}
}

func StrategyClampWarnings(before, after StrategyConfig, lang string) []string {
	if lang != "zh" {
		lang = "en"
	}
	warnings := make([]string, 0, 8)
	appendInt := func(labelZH, labelEN string, from, to int) {
		if from == to {
			return
		}
		if lang == "zh" {
			warnings = append(warnings, fmt.Sprintf("%s adjusted from %d to %d", labelZH, from, to))
			return
		}
		warnings = append(warnings, fmt.Sprintf("%s adjusted from %d to %d", labelEN, from, to))
	}
	appendFloat := func(labelZH, labelEN string, from, to float64) {
		if from == to {
			return
		}
		if lang == "zh" {
			warnings = append(warnings, fmt.Sprintf("%s adjusted from %.2f to %.2f", labelZH, from, to))
			return
		}
		warnings = append(warnings, fmt.Sprintf("%s adjusted from %.2f to %.2f", labelEN, from, to))
	}

	appendInt("Max Positions", "max_positions", before.RiskControl.MaxPositions, after.RiskControl.MaxPositions)
	appendFloat("Max Position Value Ratio", "altcoin_max_position_value_ratio", before.RiskControl.AltcoinMaxPositionValueRatio, after.RiskControl.AltcoinMaxPositionValueRatio)
	appendFloat("Min Risk/Reward Ratio", "min_risk_reward_ratio", before.RiskControl.MinRiskRewardRatio, after.RiskControl.MinRiskRewardRatio)
	appendFloat("Max Margin Usage", "max_margin_usage", before.RiskControl.MaxMarginUsage, after.RiskControl.MaxMarginUsage)
	appendFloat("Min Position Size", "min_position_size", before.RiskControl.MinPositionSize, after.RiskControl.MinPositionSize)
	appendInt("Min Confidence", "min_confidence", before.RiskControl.MinConfidence, after.RiskControl.MinConfidence)
	return warnings
}
