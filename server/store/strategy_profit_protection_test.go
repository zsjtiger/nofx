package store

import "testing"

func TestProfitProtectionDefaultsOffAndCanBeTurnedOff(t *testing.T) {
	base := GetDefaultStrategyConfig("en")
	if base.RiskControl.ProfitProtection {
		t.Fatal("default strategy must not enable profit protection")
	}

	on, err := MergeStrategyConfig(base, map[string]any{
		"ai_config": map[string]any{"risk_control": map[string]any{"profit_protection": true}},
	})
	if err != nil {
		t.Fatalf("merge on: %v", err)
	}
	if !on.RiskControl.ProfitProtection {
		t.Fatal("patch should enable profit protection")
	}

	off, err := MergeStrategyConfig(on, map[string]any{
		"ai_config": map[string]any{"risk_control": map[string]any{"profit_protection": false}},
	})
	if err != nil {
		t.Fatalf("merge off: %v", err)
	}
	if off.RiskControl.ProfitProtection {
		t.Fatal("patch with false should turn profit protection off")
	}
}
