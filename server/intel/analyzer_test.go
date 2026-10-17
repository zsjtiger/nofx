package intel

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFindingAcceptsJSONInsideProse(t *testing.T) {
	reply := "好的，结果如下：\n```json\n{\"sentiment\": 0.4, \"risk\": \"low\", \"summary\": \"业绩同比增长\"}\n```"
	f, err := parseFinding(reply)
	require.NoError(t, err)
	assert.InDelta(t, 0.4, f.Sentiment, 1e-9)
	assert.Equal(t, RiskLow, f.Risk)
	assert.Equal(t, "业绩同比增长", f.Summary)
}

func TestParseFindingClampsSentimentAndNormalizesRisk(t *testing.T) {
	f, err := parseFinding(`{"sentiment": 7, "risk": "EXTREME", "summary": "x"}`)
	require.NoError(t, err)
	assert.Equal(t, 1.0, f.Sentiment)
	assert.Equal(t, RiskUnknown, f.Risk)

	f, err = parseFinding(`{"sentiment": -9, "risk": " High ", "summary": "x"}`)
	require.NoError(t, err)
	assert.Equal(t, -1.0, f.Sentiment)
	assert.Equal(t, RiskHigh, f.Risk)
}

func TestParseFindingRejectsMalformedReplies(t *testing.T) {
	cases := map[string]string{
		"no json":         "无法判断",
		"broken json":     `{"sentiment": 0.1,`,
		"missing summary": `{"sentiment": 0.1, "risk": "none", "summary": "  "}`,
	}
	for name, reply := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := parseFinding(reply)
			assert.Error(t, err)
		})
	}
}

func TestAnalyzeSendsSymbolAndTruncatesPage(t *testing.T) {
	var userPrompt string
	complete := func(system, user string) (string, error) {
		userPrompt = user
		return `{"sentiment": 0, "risk": "none", "summary": "与该股票无关"}`, nil
	}
	long := strings.Repeat("字", MaxPageChars+500)

	f, err := Analyze(complete, "600519.SH", long)
	require.NoError(t, err)
	assert.Equal(t, RiskNone, f.Risk)
	assert.Contains(t, userPrompt, "600519.SH")
	assert.Equal(t, MaxPageChars, strings.Count(userPrompt, "字"))
}

func TestAnalyzeReportsModelFailureAndEmptyPage(t *testing.T) {
	failing := func(system, user string) (string, error) { return "", errors.New("quota") }
	_, err := Analyze(failing, "600519.SH", "some text")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "quota")

	_, err = Analyze(failing, "600519.SH", "   ")
	require.Error(t, err)
}
