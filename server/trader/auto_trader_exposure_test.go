package trader

import (
	"testing"

	"crux/store"
)

func testTraderWithRisk(rc store.RiskControlConfig) *AutoTrader {
	cfg := &store.StrategyConfig{RiskControl: rc}
	return &AutoTrader{config: AutoTraderConfig{StrategyConfig: cfg}}
}

func openPositions(values ...float64) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(values))
	for _, v := range values {
		out = append(out, map[string]interface{}{"positionAmt": v, "markPrice": 1.0})
	}
	return out
}

func TestEnforceExposureLimitCapsToRoom(t *testing.T) {
	at := testTraderWithRisk(store.RiskControlConfig{MaxMarginUsage: 1})
	// 600 already open against 1000 equity leaves 400 of room.
	size, err := at.enforceExposureLimit(openPositions(600), 900, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if size != 400 {
		t.Fatalf("size = %v, want 400", size)
	}
}

func TestEnforceExposureLimitKeepsSizeWithinRoom(t *testing.T) {
	at := testTraderWithRisk(store.RiskControlConfig{MaxMarginUsage: 1})
	size, err := at.enforceExposureLimit(openPositions(100), 300, 1000)
	if err != nil || size != 300 {
		t.Fatalf("got (%v, %v), want (300, nil)", size, err)
	}
}

func TestEnforceExposureLimitRejectsWhenNoRoom(t *testing.T) {
	at := testTraderWithRisk(store.RiskControlConfig{MaxMarginUsage: 1})
	if _, err := at.enforceExposureLimit(openPositions(1000), 500, 1000); err == nil {
		t.Fatal("expected rejection when exposure is already at the limit")
	}
}

func TestEnforceExposureLimitCountsShortsByAbsoluteValue(t *testing.T) {
	at := testTraderWithRisk(store.RiskControlConfig{MaxMarginUsage: 1})
	// A short of -600 still uses 600 of exposure.
	size, err := at.enforceExposureLimit(openPositions(-600), 900, 1000)
	if err != nil || size != 400 {
		t.Fatalf("got (%v, %v), want (400, nil)", size, err)
	}
}

func TestMinPositionSizeFallsBackToStoreMinimum(t *testing.T) {
	at := testTraderWithRisk(store.RiskControlConfig{})
	if got := at.minPositionSize(); got != store.MinPositionValue {
		t.Fatalf("minPositionSize = %v, want %v", got, store.MinPositionValue)
	}
	at = testTraderWithRisk(store.RiskControlConfig{MinPositionSize: 50})
	if got := at.minPositionSize(); got != 50 {
		t.Fatalf("minPositionSize = %v, want 50", got)
	}
}

func TestProfitProtectionIsOffByDefault(t *testing.T) {
	if testTraderWithRisk(store.RiskControlConfig{}).profitProtectionEnabled() {
		t.Fatal("profit protection must be off unless the strategy enables it")
	}
	if !testTraderWithRisk(store.RiskControlConfig{ProfitProtection: true}).profitProtectionEnabled() {
		t.Fatal("profit protection should be on when the strategy enables it")
	}
	if (&AutoTrader{}).profitProtectionEnabled() {
		t.Fatal("a trader without a strategy must not run profit protection")
	}
}
