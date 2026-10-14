package agent

import "testing"

func TestAllowlistPermitsCurrentEndpoints(t *testing.T) {
	allowed := []struct{ method, path string }{
		{"GET", "/api/equity/klines"},
		{"GET", "/api/equity-history"},
		{"GET", "/api/traders/tr1/config"},
		{"GET", "/api/strategies/s1"},
	}
	for _, c := range allowed {
		if !isPathAllowed(c.method, c.path) {
			t.Errorf("%s %s should be allowed", c.method, c.path)
		}
	}
}

func TestAllowlistRejectsRemovedEndpoints(t *testing.T) {
	denied := []struct{ method, path string }{
		{"GET", "/api/traders/tr1/grid-risk"},
		{"GET", "/api/klines"},
		{"GET", "/api/symbols"},
		{"GET", "/api/ai-costs"},
		{"POST", "/api/equity/klines"},
	}
	for _, c := range denied {
		if isPathAllowed(c.method, c.path) {
			t.Errorf("%s %s should be rejected", c.method, c.path)
		}
	}
}
