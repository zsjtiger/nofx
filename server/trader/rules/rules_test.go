package rules

import (
	"testing"
	"time"

	"crux/kernel"
)

func TestActionPredicatesAndSides(t *testing.T) {
	cases := []struct {
		action    string
		isOpen    bool
		isClose   bool
		openSide  string
		closeSide string
	}{
		{"open_long", true, false, "long", ""},
		{"  OPEN_SHORT ", true, false, "short", ""},
		{"close_long", false, true, "", "long"},
		{"Close_Short", false, true, "", "short"},
		{"hold", false, false, "", ""},
		{"", false, false, "", ""},
	}
	for _, c := range cases {
		if got := IsOpenAction(c.action); got != c.isOpen {
			t.Errorf("IsOpenAction(%q) = %v, want %v", c.action, got, c.isOpen)
		}
		if got := IsCloseAction(c.action); got != c.isClose {
			t.Errorf("IsCloseAction(%q) = %v, want %v", c.action, got, c.isClose)
		}
		if got := OpenActionSide(c.action); got != c.openSide {
			t.Errorf("OpenActionSide(%q) = %q, want %q", c.action, got, c.openSide)
		}
		if got := CloseActionSide(c.action); got != c.closeSide {
			t.Errorf("CloseActionSide(%q) = %q, want %q", c.action, got, c.closeSide)
		}
	}
}

func TestValidateProtectionPrices(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		market  float64
		stop    float64
		tp      float64
		wantErr bool
	}{
		{"long ok", "open_long", 100, 95, 110, false},
		{"long stop above market", "open_long", 100, 101, 110, true},
		{"long take profit below market", "open_long", 100, 95, 99, true},
		{"short ok", "open_short", 100, 105, 90, false},
		{"short stop below market", "open_short", 100, 99, 90, true},
		{"short take profit above market", "open_short", 100, 105, 101, true},
		{"short take profit zero", "open_short", 100, 105, 0, true},
		{"non-positive market price", "open_long", 0, 95, 110, true},
		{"non-positive stop", "open_long", 100, 0, 110, true},
		{"unsupported action", "close_long", 100, 95, 110, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProtectionPrices(tt.action, tt.market, tt.stop, tt.tp)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateProtectionPrices() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanonicalPositionSymbol(t *testing.T) {
	if got := CanonicalPositionSymbol("  aapl "); got != "AAPL" {
		t.Errorf("got %q, want AAPL", got)
	}
}

func TestFindContextPosition(t *testing.T) {
	ctx := &kernel.Context{Positions: []kernel.PositionInfo{
		{Symbol: "AAPL", Side: "long"},
		{Symbol: "MSFT", Side: "short"},
	}}
	if pos := FindContextPosition(ctx, "AAPL", "LONG"); pos == nil || pos.Symbol != "AAPL" {
		t.Errorf("expected AAPL long, got %+v", pos)
	}
	if pos := FindContextPosition(ctx, "AAPL", "short"); pos != nil {
		t.Errorf("AAPL short should not match, got %+v", pos)
	}
	if pos := FindAnyContextPosition(ctx, "MSFT"); pos == nil || pos.Side != "short" {
		t.Errorf("expected MSFT short, got %+v", pos)
	}
	if FindContextPosition(nil, "AAPL", "long") != nil || FindAnyContextPosition(nil, "AAPL") != nil {
		t.Error("nil context should return nil")
	}
}

func TestRoundDuration(t *testing.T) {
	cases := map[time.Duration]string{
		30 * time.Second:                "0m",
		90 * time.Second:                "2m0s",
		10*time.Minute + 20*time.Second: "10m0s",
		2 * time.Hour:                   "2h0m0s",
	}
	for d, want := range cases {
		if got := RoundDuration(d); got != want {
			t.Errorf("RoundDuration(%v) = %q, want %q", d, got, want)
		}
	}
}
