package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/storage"
)

func hashKey(secret, keyId, keySecret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(keyId + "." + keySecret))
	return hex.EncodeToString(h.Sum(nil))
}

// Espera header: X-API-Key: keyId.keySecret
func AuthAPIKey(repo *storage.APIKeysRepo) gin.HandlerFunc {
	secret := os.Getenv("APIKEY_HASH_SECRET")
	return func(c *gin.Context) {
		raw := c.GetHeader("X-API-Key")
		parts := strings.Split(raw, ".")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key format"})
			return
		}
		keyId, keySecret := parts[0], parts[1]
		k, err := repo.ByKeyID(c, keyId)
		if err != nil || k == nil || k.Revoked {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown api key"})
			return
		}
	if k.KeyHash != hashKey(secret, keyId, keySecret) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
		return
	}

	// ✅ Adicionar API key e tenant_id ao contexto para outros middlewares
	// NOTA: OwnerID na prática contém o TenantID (não o UserID)
	// TODO: Futuramente, migrar para usar UserID e buscar o tenant do usuário
	c.Set("api_key", raw)
	c.Set("tenant_id", k.OwnerID) // OwnerID = TenantID na implementação atual

	c.Next()
}
}
