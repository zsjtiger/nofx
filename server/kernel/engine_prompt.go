package kernel

import (
	"fmt"
	"strings"

	"crux/market"
)

// BuildSystemPrompt returns the decision prompt for the engine's market. It
// returns an empty string when no stock market is configured; the caller checks
// capabilities first and reports the error.
func (e *StrategyEngine) BuildSystemPrompt(accountEquity float64, _ string) string {
	caps, err := e.capabilities()
	if err != nil {
		return ""
	}
	return BuildEquitySystemPrompt(caps, e.config.RiskControl, accountEquity, e.config.CustomPrompt)
}

// BuildUserPrompt describes the current cycle: account, positions, recent
// results and the market data for every candidate symbol.
func (e *StrategyEngine) BuildUserPrompt(ctx *Context) string {
	var sb strings.Builder
	caps, _ := e.capabilities()

	sb.WriteString(fmt.Sprintf("Time: %s | Cycle: #%d | Runtime: %d min\n\n", ctx.CurrentTime, ctx.CallCount, ctx.RuntimeMinutes))

	sb.WriteString("## Account\n")
	sb.WriteString(fmt.Sprintf("- Equity: %.2f %s\n", ctx.Account.TotalEquity, caps.Currency))
	sb.WriteString(fmt.Sprintf("- Available cash: %.2f %s\n", ctx.Account.AvailableBalance, caps.Currency))
	sb.WriteString(fmt.Sprintf("- Unrealized P&L: %+.2f %s (%+.2f%% of equity)\n", ctx.Account.UnrealizedPnL, caps.Currency, pctOf(ctx.Account.UnrealizedPnL, ctx.Account.TotalEquity)))
	sb.WriteString(fmt.Sprintf("- Open positions: %d\n\n", ctx.Account.PositionCount))

	sb.WriteString("## Current positions\n")
	if len(ctx.Positions) == 0 {
		sb.WriteString("None.\n\n")
	} else {
		for i, pos := range ctx.Positions {
			sb.WriteString(formatPositionInfo(i+1, pos, caps.Currency))
		}
		sb.WriteString("\n")
	}

	if ctx.TradingStats != nil {
		s := ctx.TradingStats
		sb.WriteString("## Closed-trade statistics\n")
		sb.WriteString(fmt.Sprintf("- Trades: %d | Win rate: %.1f%% | Profit factor: %.2f | Total P&L: %+.2f %s | Max drawdown: %.2f%%\n\n",
			s.TotalTrades, s.WinRate, s.ProfitFactor, s.TotalPnL, caps.Currency, s.MaxDrawdownPct))
	}

	if len(ctx.RecentOrders) > 0 {
		sb.WriteString("## Recently closed trades\n")
		for _, o := range ctx.RecentOrders {
			sb.WriteString(fmt.Sprintf("- %s %s: %.2f -> %.2f, P&L %+.2f (%+.2f%%), held %s\n",
				o.Symbol, o.Side, o.EntryPrice, o.ExitPrice, o.RealizedPnL, o.PnLPct, o.HoldDuration))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Candidate symbols (daily bars, oldest to latest)\n\n")
	for _, coin := range ctx.CandidateCoins {
		data, ok := ctx.MarketDataMap[coin.Symbol]
		if !ok || data == nil {
			sb.WriteString(fmt.Sprintf("### %s\nNo market data available this cycle; do not trade it.\n\n", coin.Symbol))
			continue
		}
		sb.WriteString(fmt.Sprintf("### %s\n", coin.Symbol))
		sb.WriteString(market.Format(data))
		sb.WriteString("\n")
	}

	sb.WriteString("Decide for each symbol above and for each open position. Reply in the format the system prompt specifies.\n")
	return sb.String()
}

func formatPositionInfo(index int, pos PositionInfo, currency string) string {
	return fmt.Sprintf("%d. %s %s | qty %.4f | entry %.2f | mark %.2f | unrealized %+.2f %s (%+.2f%%) | peak %+.2f%%\n",
		index, pos.Symbol, strings.ToUpper(pos.Side), pos.Quantity, pos.EntryPrice, pos.MarkPrice,
		pos.UnrealizedPnL, currency, pos.UnrealizedPnLPct, pos.PeakPnLPct)
}

func pctOf(value, base float64) float64 {
	if base == 0 {
		return 0
	}
	return value / base * 100
}
