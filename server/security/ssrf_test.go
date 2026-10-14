package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateURLRejectsUnsafeTargets(t *testing.T) {
	cases := map[string]string{
		"empty":               "",
		"file scheme":         "file:///etc/passwd",
		"ftp scheme":          "ftp://example.com/x",
		"localhost":           "http://localhost:8080/api",
		"loopback IP":         "http://127.0.0.1/admin",
		"IPv6 loopback":       "http://[::1]/admin",
		"private 10/8":        "http://10.0.0.5/internal",
		"private 192.168":     "http://192.168.1.10/",
		"cloud metadata":      "http://169.254.169.254/latest/meta-data",
		"metadata host":       "http://metadata.google.internal/",
		"unspecified 0.0.0.0": "http://0.0.0.0/",
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			err := ValidateURL(raw)
			require.Error(t, err)
			assert.IsType(t, &SSRFError{}, err)
		})
	}
}

func TestValidateURLAllowsPublicIPLiterals(t *testing.T) {
	// IP literals resolve without DNS, so these checks are deterministic offline.
	assert.NoError(t, ValidateURL("https://8.8.8.8/dns-query"))
	assert.NoError(t, ValidateURL("https://1.1.1.1/"))
}

func TestSafeGetBlocksLoopbackServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// httptest listens on 127.0.0.1, which the SSRF guard must refuse.
	resp, err := SafeGet(server.URL, 2*time.Second)
	if resp != nil {
		resp.Body.Close()
	}
	assert.Error(t, err)
}
