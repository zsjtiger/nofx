package kernel

import (
	"fmt"
	"strings"

	"nofx/logger"
	"nofx/store"
	"nofx/trader/types"
)

// BuildEquitySystemPrompt is the decision prompt for stock markets. It states the
// market rules the code enforces, so the model does not propose trades that will be rejected.
func BuildEquitySystemPrompt(caps types.Capabilities, riskControl store.RiskControlConfig, accountEquity float64, customPrompt string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("You are a disciplined %s equity trader. Decide what to do with each stock in the universe, one decision per symbol.\n\n", marketLabel(caps)))

	sb.WriteString("# Market rules (enforced in code)\n")
	sb.WriteString(fmt.Sprintf("- Currency: %s. Session time zone: %s.\n", caps.Currency, caps.Timezone))
	sb.WriteString("- Data: daily bars only; the latest bar is the most recent completed session. No intraday signals.\n")
	if caps.AllowShort {
		sb.WriteString("- Short selling is allowed. Use `open_short` only for a confirmed downtrend.\n")
	} else {
		sb.WriteString("- Short selling is NOT allowed. Never output `open_short`.\n")
	}
	sb.WriteString("- Leverage is NOT allowed. Always set `leverage` to 1.\n")
	if caps.LotSize > 0 {
		sb.WriteString(fmt.Sprintf("- Buys are in whole lots of %d shares. Size positions so they round to a whole lot.\n", caps.LotSize))
	}
	if caps.TPlusOne {
		sb.WriteString("- T+1: shares bought today cannot be sold today. Do not propose `close_long` for a position opened in the current session.\n")
	}
	sb.WriteString("- Daily price limits apply. A buy at limit-up or a sell at limit-down will not fill; avoid those entries and exits.\n\n")

	sb.WriteString("# Risk limits\n")
	sb.WriteString(fmt.Sprintf("- Account equity: %.2f %s.\n", accountEquity, caps.Currency))
	maxPositions := riskControl.MaxPositions
	if maxPositions <= 0 {
		maxPositions = 3
	}
	sb.WriteString(fmt.Sprintf("- At most %d open positions at once.\n", maxPositions))
	ratio := riskControl.AltcoinMaxPositionValueRatio
	if ratio <= 0 {
		ratio = 1
	}
	sb.WriteString(fmt.Sprintf("- Single position value must not exceed %.2f %s (equity x %.2f).\n", accountEquity*ratio, caps.Currency, ratio))
	if riskControl.MinConfidence > 0 {
		sb.WriteString(fmt.Sprintf("- Only open positions with confidence of at least %d.\n", riskControl.MinConfidence))
	}
	sb.WriteString("- Every open must include a protective `stop_loss` below the entry for longs; above it for shorts.\n\n")

	if s := strings.TrimSpace(customPrompt); s != "" {
		sb.WriteString("# Additional instructions from the strategy owner\n")
		sb.WriteString(s)
		sb.WriteString("\n\n")
	}

	sb.WriteString("# Output\n")
	sb.WriteString("First explain your reasoning briefly. Then output a JSON array with one object per decision, inside a ```json code block:\n")
	sb.WriteString("```json\n[\n")
	sb.WriteString(fmt.Sprintf("  {\"symbol\": \"EXAMPLE\", \"action\": \"open_long\", \"leverage\": 1, \"position_size_usd\": 1000, \"stop_loss\": 0, \"take_profit\": 0, \"confidence\": 80, \"risk_usd\": 50, \"reasoning\": \"...\"}\n"))
	sb.WriteString("]\n```\n")
	sb.WriteString("- `action`: open_long | close_long | hold | wait")
	if caps.AllowShort {
		sb.WriteString(" | open_short | close_short")
	}
	sb.WriteString("\n- Required when opening: leverage (1), position_size_usd, stop_loss, take_profit, confidence, risk_usd\n")
	sb.WriteString("- Prices are in the market currency. Never copy the example values.\n")
	return sb.String()
}

func marketLabel(caps types.Capabilities) string {
	switch caps.Market {
	case types.MarketCNEquity:
		return "A-share"
	case types.MarketUSEquity:
		return "US"
	default:
		return "stock"
	}
}

// ValidateEquityDecisions enforces the market rules on the model's decisions and
// fails closed: a decision the market cannot take rejects the whole response.
// Leverage is clamped to 1 when the market does not allow it, as the crypto path does.
func ValidateEquityDecisions(decisions []Decision, caps types.Capabilities, accountEquity float64, riskControl store.RiskControlConfig) error {
	ratio := riskControl.AltcoinMaxPositionValueRatio
	if ratio <= 0 {
		ratio = 1
	}
	maxValue := accountEquity * ratio
	for i := range decisions {
		d := &decisions[i]
		if d.Action == "open_short" || d.Action == "close_short" {
			if !caps.AllowShort {
				return fmt.Errorf("decision #%d: %s is not allowed in %s", i+1, d.Action, caps.Market)
			}
		}
		if !caps.AllowLeverage && d.Leverage != 1 {
			if d.Leverage > 1 {
				logger.Infof("⚠️  [Equity] %s leverage %dx not allowed in %s, using 1x", d.Symbol, d.Leverage, caps.Market)
			}
			d.Leverage = 1
		}
		if d.Action != "open_long" && d.Action != "open_short" {
			continue
		}
		if d.StopLoss <= 0 {
			return fmt.Errorf("decision #%d: %s on %s needs a positive stop_loss", i+1, d.Action, d.Symbol)
		}
		if d.PositionSizeUSD <= 0 {
			return fmt.Errorf("decision #%d: position_size_usd must be positive for %s", i+1, d.Symbol)
		}
		if maxValue > 0 && d.PositionSizeUSD > maxValue+1e-9 {
			d.PositionSizeUSD = maxValue
			logger.Infof("⚠️  [Equity] %s position capped at %.2f %s", d.Symbol, maxValue, caps.Currency)
		}
	}
	return nil
}
