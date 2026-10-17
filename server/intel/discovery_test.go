package intel

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifiedSymbolsKeepsOnlyCodesInPageText(t *testing.T) {
	page := "贵州茅台与五粮液（000858）同属白酒板块，AAPL 也被提及。"
	got := verifiedSymbols([]string{"000858.SZ", "601318.SH", "AAPL", "", "MSFT"}, page)
	assert.Equal(t, []string{"000858.SZ", "AAPL"}, got, "made-up codes must be dropped")
}

func TestParseFindingCleansReportedSymbols(t *testing.T) {
	reply := `{"sentiment": 0.1, "risk": "none", "summary": "行业提及", "symbols": [" 600519.sh ", "600519.SH", "", "a", "b", "c", "d", "e", "f"]}`
	f, err := parseFinding(reply)
	require.NoError(t, err)
	assert.Equal(t, []string{"600519.SH", "A", "B", "C", "D"}, f.Symbols, "trimmed, upper-cased, de-duplicated and capped")
}

func TestAnalysisPromptAsksForMentionedSymbols(t *testing.T) {
	var system string
	complete := func(s, u string) (string, error) {
		system = s
		return `{"sentiment": 0, "risk": "none", "summary": "x", "symbols": []}`, nil
	}
	_, err := Analyze(complete, "600519.SH", strings.Repeat("正文", 10))
	require.NoError(t, err)
	assert.Contains(t, system, `"symbols"`)
}
