package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRiskControlReadsLegacyPositionRatioKey(t *testing.T) {
	var rc RiskControlConfig
	rc.MaxPositionValueRatio = 1
	require.NoError(t, json.Unmarshal([]byte(`{"altcoin_max_position_value_ratio": 0.5}`), &rc))
	assert.Equal(t, 0.5, rc.MaxPositionValueRatio)
}

func TestRiskControlPrefersCurrentPositionRatioKey(t *testing.T) {
	var rc RiskControlConfig
	require.NoError(t, json.Unmarshal([]byte(`{"max_position_value_ratio": 0.7, "altcoin_max_position_value_ratio": 0.5}`), &rc))
	assert.Equal(t, 0.7, rc.MaxPositionValueRatio)
}

func TestRiskControlWritesOnlyCurrentKey(t *testing.T) {
	out, err := json.Marshal(RiskControlConfig{MaxPositionValueRatio: 0.9})
	require.NoError(t, err)
	assert.Contains(t, string(out), `"max_position_value_ratio":0.9`)
	assert.NotContains(t, string(out), "altcoin")
}
