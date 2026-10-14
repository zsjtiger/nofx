package kernel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"nofx/logger"
	"nofx/market"
	"nofx/security"
	"nofx/store"
	"nofx/trader/types"
)

// StrategyEngine decides for one stock market. The market is set by WithMarket
// when the trader is created, and it selects the rules, prompt and data source.
type StrategyEngine struct {
	config *store.StrategyConfig
	market types.Market
}

// NewStrategyEngine creates a strategy engine. Call WithMarket before deciding.
func NewStrategyEngine(config *store.StrategyConfig) *StrategyEngine {
	return &StrategyEngine{config: config}
}

// WithMarket sets the market the engine decides for and returns the engine.
func (e *StrategyEngine) WithMarket(m types.Market) *StrategyEngine {
	e.market = m
	return e
}

// capabilities returns the rules for the engine's market. An engine without a
// valid stock market refuses to decide.
func (e *StrategyEngine) capabilities() (types.Capabilities, error) {
	if e.market == "" {
		return types.Capabilities{}, fmt.Errorf("strategy engine has no market configured")
	}
	return types.CapabilitiesOf(e.market)
}

// GetRiskControlConfig returns the risk control configuration.
func (e *StrategyEngine) GetRiskControlConfig() store.RiskControlConfig {
	return e.config.RiskControl
}

// GetLanguage returns the language from config or falls back to auto-detection.
func (e *StrategyEngine) GetLanguage() Language {
	switch e.config.Language {
	case "zh":
		return LangChinese
	case "en":
		return LangEnglish
	default:
		return detectLanguage(e.config.PromptSections.RoleDefinition)
	}
}

// GetConfig returns the complete strategy configuration.
func (e *StrategyEngine) GetConfig() *store.StrategyConfig {
	return e.config
}

// GetCandidateCoins returns the strategy's stock universe. Only the static
// list is supported; index-based universes are not implemented yet.
func (e *StrategyEngine) GetCandidateCoins() ([]CandidateCoin, error) {
	caps, err := e.capabilities()
	if err != nil {
		return nil, err
	}
	coinSource := e.config.CoinSource
	if coinSource.SourceType != "static" {
		return nil, fmt.Errorf("coin source %q is not available for %s; use a static symbol list", coinSource.SourceType, caps.Market)
	}

	candidates := make([]CandidateCoin, 0, len(coinSource.StaticCoins))
	seen := make(map[string]bool, len(coinSource.StaticCoins))
	for _, raw := range coinSource.StaticCoins {
		symbol, err := types.NormalizeSymbol(caps.Market, raw)
		if err != nil {
			logger.Warnf("⚠️  Skipping invalid symbol %q: %v", raw, err)
			continue
		}
		if seen[symbol] {
			continue
		}
		seen[symbol] = true
		candidates = append(candidates, CandidateCoin{Symbol: symbol, Sources: []string{"static"}})
	}
	return e.filterExcludedCoins(candidates), nil
}

func (e *StrategyEngine) filterExcludedCoins(candidates []CandidateCoin) []CandidateCoin {
	if len(e.config.CoinSource.ExcludedCoins) == 0 {
		return candidates
	}
	caps, err := e.capabilities()
	if err != nil {
		return candidates
	}
	excluded := make(map[string]bool, len(e.config.CoinSource.ExcludedCoins))
	for _, coin := range e.config.CoinSource.ExcludedCoins {
		if symbol, err := types.NormalizeSymbol(caps.Market, coin); err == nil {
			excluded[symbol] = true
		}
	}

	filtered := make([]CandidateCoin, 0, len(candidates))
	for _, c := range candidates {
		if excluded[c.Symbol] {
			logger.Infof("🚫 Excluded symbol: %s", c.Symbol)
			continue
		}
		filtered = append(filtered, c)
	}
	return filtered
}

// FetchMarketData returns daily-bar market data for a symbol in the engine's market.
func (e *StrategyEngine) FetchMarketData(symbol string) (*market.Data, error) {
	if _, err := e.capabilities(); err != nil {
		return nil, err
	}
	return market.GetForMarket(symbol, e.market)
}

// FetchExternalData fetches the strategy's configured external data sources.
func (e *StrategyEngine) FetchExternalData() (map[string]interface{}, error) {
	externalData := make(map[string]interface{})
	for _, source := range e.config.Indicators.ExternalDataSources {
		data, err := e.fetchSingleExternalSource(source)
		if err != nil {
			logger.Infof("⚠️  Failed to fetch external data source [%s]: %v", source.Name, err)
			continue
		}
		externalData[source.Name] = data
	}
	return externalData, nil
}

func (e *StrategyEngine) fetchSingleExternalSource(source store.ExternalDataSource) (interface{}, error) {
	// SSRF protection: validate the URL before making the request.
	if err := security.ValidateURL(source.URL); err != nil {
		return nil, fmt.Errorf("external source URL validation failed: %w", err)
	}

	timeout := time.Duration(source.RefreshSecs) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	client := security.SafeHTTPClient(timeout)

	req, err := http.NewRequest(source.Method, source.URL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range source.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if source.DataPath != "" {
		result = extractJSONPath(result, source.DataPath)
	}
	return result, nil
}

func extractJSONPath(data interface{}, path string) interface{} {
	current := data
	for _, part := range strings.Split(path, ".") {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = m[part]
	}
	return current
}

// detectLanguage returns LangChinese if the text contains Chinese characters, otherwise LangEnglish.
func detectLanguage(text string) Language {
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			return LangChinese
		}
	}
	return LangEnglish
}
