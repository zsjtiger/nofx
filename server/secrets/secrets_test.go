package secrets

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestCryptoService builds a CryptoService from freshly generated keys,
// through the same environment variables the server reads at startup.
func newTestCryptoService(t *testing.T) *CryptoService {
	t.Helper()
	privatePEM, _, err := GenerateKeyPair()
	require.NoError(t, err)
	dataKey, err := GenerateDataKey()
	require.NoError(t, err)

	t.Setenv(EnvRSAPrivateKey, strings.ReplaceAll(privatePEM, "\n", "\\n"))
	t.Setenv(EnvDataEncryptionKey, dataKey)

	cs, err := NewCryptoService()
	require.NoError(t, err)
	return cs
}

func TestNewCryptoServiceRequiresBothKeys(t *testing.T) {
	t.Setenv(EnvRSAPrivateKey, "")
	t.Setenv(EnvDataEncryptionKey, "")
	_, err := NewCryptoService()
	assert.ErrorContains(t, err, "RSA private key")

	privatePEM, _, err := GenerateKeyPair()
	require.NoError(t, err)
	t.Setenv(EnvRSAPrivateKey, strings.ReplaceAll(privatePEM, "\n", "\\n"))
	_, err = NewCryptoService()
	assert.ErrorContains(t, err, "data encryption key")
}

func TestParseRSAPrivateKeyRejectsGarbage(t *testing.T) {
	_, err := ParseRSAPrivateKeyFromPEM([]byte("not a key"))
	assert.Error(t, err)
}

func TestStorageEncryptionRoundTrip(t *testing.T) {
	cs := newTestCryptoService(t)

	sealed, err := cs.EncryptForStorage("sk-live-abc123", "exchange", "user-1")
	require.NoError(t, err)
	assert.NotContains(t, sealed, "sk-live-abc123", "the stored value must not contain the secret")
	assert.True(t, cs.IsEncryptedStorageValue(sealed))

	plain, err := cs.DecryptFromStorage(sealed, "exchange", "user-1")
	require.NoError(t, err)
	assert.Equal(t, "sk-live-abc123", plain)
}

func TestStorageEncryptionUsesFreshNonce(t *testing.T) {
	cs := newTestCryptoService(t)
	a, err := cs.EncryptForStorage("same value")
	require.NoError(t, err)
	b, err := cs.EncryptForStorage("same value")
	require.NoError(t, err)
	assert.NotEqual(t, a, b)
}

func TestStorageEncryptionEmptyAndAlreadyEncrypted(t *testing.T) {
	cs := newTestCryptoService(t)

	empty, err := cs.EncryptForStorage("")
	require.NoError(t, err)
	assert.Equal(t, "", empty)

	sealed, err := cs.EncryptForStorage("secret")
	require.NoError(t, err)
	again, err := cs.EncryptForStorage(sealed)
	require.NoError(t, err)
	assert.Equal(t, sealed, again, "an encrypted value must not be encrypted twice")
}

func TestStorageDecryptionRejectsWrongContextAndTampering(t *testing.T) {
	cs := newTestCryptoService(t)
	sealed, err := cs.EncryptForStorage("secret", "exchange", "user-1")
	require.NoError(t, err)

	_, err = cs.DecryptFromStorage(sealed, "exchange", "user-2")
	assert.Error(t, err, "a value moved to another record must not decrypt")

	tampered := sealed[:len(sealed)-2] + "AA"
	_, err = cs.DecryptFromStorage(tampered, "exchange", "user-1")
	assert.Error(t, err)

	_, err = cs.DecryptFromStorage("plain text, never encrypted", "exchange", "user-1")
	assert.ErrorContains(t, err, "not encrypted")
}

func TestStorageEncryptionNeedsDataKey(t *testing.T) {
	cs := &CryptoService{}
	assert.False(t, cs.HasDataKey())
	_, err := cs.EncryptForStorage("secret")
	assert.Error(t, err)
}

func TestDecryptPayloadRejectsMissingOrStaleTimestamp(t *testing.T) {
	cs := newTestCryptoService(t)

	_, err := cs.DecryptPayload(&EncryptedPayload{})
	assert.ErrorContains(t, err, "missing timestamp")

	stale := &EncryptedPayload{TS: time.Now().Add(-10 * time.Minute).Unix()}
	_, err = cs.DecryptPayload(stale)
	assert.ErrorContains(t, err, "expired")

	future := &EncryptedPayload{TS: time.Now().Add(10 * time.Minute).Unix()}
	_, err = cs.DecryptPayload(future)
	assert.ErrorContains(t, err, "expired")
}
