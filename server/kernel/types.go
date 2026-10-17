package kernel

import (
	"time"

	"crux/market"
)

// Language is the language a strategy's text is written in.
type Language string

const (
	LangChinese Language = "zh-CN"
	LangEnglish Language = "en-US"
)

// PositionInfo is one open position as shown to the model.
type PositionInfo struct {
	Symbol           string  `json:"symbol"`
	Side             string  `json:"side"` // "long" or "short"
	EntryPrice       float64 `json:"entry_price"`
	MarkPrice        float64 `json:"mark_price"`
	Quantity         float64 `json:"quantity"`
	UnrealizedPnL    float64 `json:"unrealized_pnl"`
	UnrealizedPnLPct float64 `json:"unrealized_pnl_pct"`
	PeakPnLPct       float64 `json:"peak_pnl_pct"` // Historical peak profit percentage
	MarginUsed       float64 `json:"margin_used"`
	UpdateTime       int64   `json:"update_time"` // Position update timestamp (milliseconds)
}

// AccountInfo is the account summary shown to the model.
type AccountInfo struct {
	TotalEquity      float64 `json:"total_equity"`
	AvailableBalance float64 `json:"available_balance"`
	UnrealizedPnL    float64 `json:"unrealized_pnl"`
	TotalPnL         float64 `json:"total_pnl"`
	TotalPnLPct      float64 `json:"total_pnl_pct"`
	MarginUsed       float64 `json:"margin_used"`
	MarginUsedPct    float64 `json:"margin_used_pct"`
	PositionCount    int     `json:"position_count"`
}

// CandidateCoin is a symbol in the strategy's universe.
type CandidateCoin struct {
	Symbol  string   `json:"symbol"`
	Sources []string `json:"sources"`
}

// TradingStats are closed-trade statistics shown to the model.
type TradingStats struct {
	TotalTrades    int     `json:"total_trades"`
	WinRate        float64 `json:"win_rate"`
	ProfitFactor   float64 `json:"profit_factor"`
	SharpeRatio    float64 `json:"sharpe_ratio"`
	TotalPnL       float64 `json:"total_pnl"`
	AvgWin         float64 `json:"avg_win"`
	AvgLoss        float64 `json:"avg_loss"`
	MaxDrawdownPct float64 `json:"max_drawdown_pct"`
}

// RecentOrder is a recently closed trade shown to the model.
type RecentOrder struct {
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`
	EntryPrice   float64 `json:"entry_price"`
	ExitPrice    float64 `json:"exit_price"`
	RealizedPnL  float64 `json:"realized_pnl"`
	PnLPct       float64 `json:"pnl_pct"`
	EntryTime    string  `json:"entry_time"`
	ExitTime     string  `json:"exit_time"`
	HoldDuration string  `json:"hold_duration"`
}

// Context is everything the decision step knows about one cycle.
type Context struct {
	CurrentTime    string                  `json:"current_time"`
	RuntimeMinutes int                     `json:"runtime_minutes"`
	CallCount      int                     `json:"call_count"`
	Account        AccountInfo             `json:"account"`
	Positions      []PositionInfo          `json:"positions"`
	CandidateCoins []CandidateCoin         `json:"candidate_coins"`
	PromptVariant  string                  `json:"prompt_variant,omitempty"`
	TradingStats   *TradingStats           `json:"trading_stats,omitempty"`
	RecentOrders   []RecentOrder           `json:"recent_orders,omitempty"`
	MarketDataMap  map[string]*market.Data `json:"-"`
	// Intel holds recent analysed web pages per symbol, newest first.
	Intel         map[string][]IntelNote             `json:"-"`
	MultiTFMarket map[string]map[string]*market.Data `json:"-"`
	Timeframes    []string                           `json:"-"`
}

// IntelNote is one analysed web page about a candidate symbol.
type IntelNote struct {
	Time      time.Time
	Sentiment float64 // -1 to 1
	Risk      string  // none, low, high, or unknown
	Summary   string
}

// Decision is one trading decision from the model.
type Decision struct {
	Symbol string `json:"symbol"`
	Action string `json:"action"` // "open_long", "open_short", "close_long", "close_short", "hold", "wait"

	// Opening position parameters
	PositionSizeUSD float64 `json:"position_size_usd,omitempty"`
	StopLoss        float64 `json:"stop_loss,omitempty"`
	TakeProfit      float64 `json:"take_profit,omitempty"`

	// Common parameters
	Confidence int     `json:"confidence,omitempty"` // Confidence level (0-100)
	RiskUSD    float64 `json:"risk_usd,omitempty"`   // Maximum loss in account currency
	Reasoning  string  `json:"reasoning"`
}

// FullDecision is the model's complete answer for one cycle, including the prompts sent.
type FullDecision struct {
	SystemPrompt        string     `json:"system_prompt"`
	UserPrompt          string     `json:"user_prompt"`
	CoTTrace            string     `json:"cot_trace"`
	Decisions           []Decision `json:"decisions"`
	RawResponse         string     `json:"raw_response"`
	Timestamp           time.Time  `json:"timestamp"`
	AIRequestDurationMs int64      `json:"ai_request_duration_ms,omitempty"`
}
