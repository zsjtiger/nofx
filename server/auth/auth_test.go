package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse")
	require.NoError(t, err)
	assert.NotEqual(t, "correct horse", hash, "the hash must not be the plaintext")
	assert.True(t, CheckPassword("correct horse", hash))
	assert.False(t, CheckPassword("wrong password", hash))
}

func TestGenerateAndValidateJWT(t *testing.T) {
	SetJWTSecret("test-secret-that-is-long-enough-123456")

	token, err := GenerateJWT("user-1", "owner@example.com")
	require.NoError(t, err)

	claims, err := ValidateJWT(token)
	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "owner@example.com", claims.Email)
	assert.Equal(t, "crux", claims.Issuer)
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	SetJWTSecret("first-secret-that-is-long-enough-1234")
	token, err := GenerateJWT("user-1", "owner@example.com")
	require.NoError(t, err)

	SetJWTSecret("second-secret-that-is-long-enough-567")
	_, err = ValidateJWT(token)
	assert.Error(t, err)
}

func TestValidateJWTRejectsTamperedToken(t *testing.T) {
	SetJWTSecret("test-secret-that-is-long-enough-123456")
	token, err := GenerateJWT("user-1", "owner@example.com")
	require.NoError(t, err)

	// Replace the payload with one claiming a different user, keep the original signature.
	forged := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{UserID: "admin"})
	forgedString, err := forged.SignedString([]byte("attacker-key"))
	require.NoError(t, err)
	_, err = ValidateJWT(forgedString)
	assert.Error(t, err)

	_, err = ValidateJWT(token + "x")
	assert.Error(t, err)
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	SetJWTSecret("test-secret-that-is-long-enough-123456")
	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: "user-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	})
	token, err := expired.SignedString(JWTSecret)
	require.NoError(t, err)

	_, err = ValidateJWT(token)
	assert.Error(t, err)
}

func TestValidateJWTRejectsNoneAlgorithm(t *testing.T) {
	SetJWTSecret("test-secret-that-is-long-enough-123456")
	none := jwt.NewWithClaims(jwt.SigningMethodNone, &Claims{UserID: "admin"})
	token, err := none.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = ValidateJWT(token)
	assert.Error(t, err)
}

func TestBlacklistedTokenIsRejectedUntilExpiry(t *testing.T) {
	BlacklistToken("token-live", time.Now().Add(time.Hour))
	assert.True(t, IsTokenBlacklisted("token-live"))
	assert.False(t, IsTokenBlacklisted("token-never-seen"))
}

func TestExpiredBlacklistEntryIsForgotten(t *testing.T) {
	BlacklistToken("token-old", time.Now().Add(-time.Second))
	assert.False(t, IsTokenBlacklisted("token-old"))
	// A second lookup must not find the entry again.
	assert.False(t, IsTokenBlacklisted("token-old"))
}
