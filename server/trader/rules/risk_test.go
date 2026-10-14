package rules

import "testing"

func TestValidateRiskReward(t *testing.T) {
	cases := []struct {
		name          string
		price, sl, tp float64
		minRatio      float64
		wantErr       bool
	}{
		{"meets minimum long", 100, 95, 110, 2, false},
		{"below minimum long", 100, 95, 107, 2, true},
		{"meets minimum short", 100, 105, 90, 2, false},
		{"below minimum short", 100, 105, 99, 2, true},
		{"disabled when zero", 100, 95, 100.1, 0, false},
		{"stop equals price", 100, 100, 110, 2, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateRiskReward(c.price, c.sl, c.tp, c.minRatio)
			if (err != nil) != c.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func TestExposureRoom(t *testing.T) {
	if got := ExposureRoom(40, 100, 1); got != 60 {
		t.Errorf("room = %v, want 60", got)
	}
	if got := ExposureRoom(150, 100, 1); got != 0 {
		t.Errorf("over-exposed room = %v, want 0", got)
	}
	if got := ExposureRoom(0, 100, 0); got != 100 {
		t.Errorf("zero maxUsage should default to 1, got %v", got)
	}
	if got := ExposureRoom(0, 100, 2); got != 200 {
		t.Errorf("maxUsage 2 room = %v, want 200", got)
	}
}
