package intel

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// Completer sends a system and user prompt to an AI model and returns its reply.
type Completer func(systemPrompt, userPrompt string) (string, error)

// Finding is the structured result of analysing one page.
type Finding struct {
	Sentiment float64 `json:"sentiment"` // -1 (very negative) to 1 (very positive)
	Risk      string  `json:"risk"`      // none, low, high, or unknown
	Summary   string  `json:"summary"`
}

const (
	RiskNone    = "none"
	RiskLow     = "low"
	RiskHigh    = "high"
	RiskUnknown = "unknown"

	// MaxPageChars limits the page text sent to the model.
	MaxPageChars = 12000
)

const analysisSystemPrompt = `你是证券情报分析员。只依据用户提供的网页内容判断，不要引入网页之外的信息，也不要编造。
只输出一个 JSON 对象，不要输出其他文字：
{"sentiment": 数字，-1 到 1 之间，表示对该股票的情绪倾向，"risk": "none" | "low" | "high"，表示是否出现需要警惕的风险事件，"summary": "不超过80字的中文摘要"}
如果网页与该股票无关，sentiment 设为 0，risk 设为 "none"，summary 写明"与该股票无关"。`

// Analyze asks the model for a sentiment and risk assessment of pageText with
// respect to symbol.
func Analyze(complete Completer, symbol, pageText string) (Finding, error) {
	if strings.TrimSpace(pageText) == "" {
		return Finding{}, fmt.Errorf("empty page text")
	}
	text := truncateRunes(pageText, MaxPageChars)
	user := fmt.Sprintf("股票代码：%s\n\n网页内容：\n%s", symbol, text)

	reply, err := complete(analysisSystemPrompt, user)
	if err != nil {
		return Finding{}, fmt.Errorf("model call failed: %w", err)
	}
	return parseFinding(reply)
}

// parseFinding extracts the JSON object from a model reply, which may include
// surrounding prose or a code fence, and validates its fields.
func parseFinding(reply string) (Finding, error) {
	start := strings.Index(reply, "{")
	end := strings.LastIndex(reply, "}")
	if start < 0 || end <= start {
		return Finding{}, fmt.Errorf("model reply contains no JSON object")
	}
	var f Finding
	if err := json.Unmarshal([]byte(reply[start:end+1]), &f); err != nil {
		return Finding{}, fmt.Errorf("parse model reply: %w", err)
	}
	f.Summary = strings.TrimSpace(f.Summary)
	if f.Summary == "" {
		return Finding{}, fmt.Errorf("model reply has no summary")
	}
	f.Sentiment = math.Max(-1, math.Min(1, f.Sentiment))
	switch strings.ToLower(strings.TrimSpace(f.Risk)) {
	case RiskNone, RiskLow, RiskHigh:
		f.Risk = strings.ToLower(strings.TrimSpace(f.Risk))
	default:
		f.Risk = RiskUnknown
	}
	return f, nil
}

func truncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}
