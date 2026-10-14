package api

import (
	"crux/config"
	"crux/secrets"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EncryptionHandler Encryption API handler
type EncryptionHandler struct {
	cryptoService *secrets.CryptoService
}

// NewEncryptionHandler Creates encryption handler
func NewEncryptionHandler(cryptoService *secrets.CryptoService) *EncryptionHandler {
	return &EncryptionHandler{
		cryptoService: cryptoService,
	}
}

// ==================== Crypto Config Endpoint ====================

// HandleGetCryptoConfig Get crypto configuration
func (h *EncryptionHandler) HandleGetCryptoConfig(c *gin.Context) {
	cfg := config.Get()
	c.JSON(http.StatusOK, gin.H{
		"transport_encryption": cfg.TransportEncryption,
	})
}

// ==================== Public Key Endpoint ====================

// HandleGetPublicKey Get server public key
func (h *EncryptionHandler) HandleGetPublicKey(c *gin.Context) {
	cfg := config.Get()
	if !cfg.TransportEncryption {
		c.JSON(http.StatusOK, gin.H{
			"public_key":           "",
			"algorithm":            "",
			"transport_encryption": false,
		})
		return
	}

	publicKey := h.cryptoService.GetPublicKeyPEM()
	c.JSON(http.StatusOK, gin.H{
		"public_key":           publicKey,
		"algorithm":            "RSA-OAEP-2048",
		"transport_encryption": true,
	})
}

// ==================== Encrypted Data Decryption ====================
//
// SECURITY: there is deliberately NO public decrypt endpoint. Transport
// encryption is one-directional — clients encrypt sensitive fields to the
// server's RSA public key and the authenticated config-update handlers
// (handleUpdateModelConfigs / handleUpdateExchangeConfigs / handleCreateExchange)
// decrypt them server-side via cryptoService.DecryptSensitiveData. Exposing a
// generic decrypt route would turn the server into a decryption oracle that any
// unauthenticated caller could use to recover the plaintext of a captured
// ciphertext, defeating the entire transport-encryption layer.
