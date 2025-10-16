package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"net/http"
	"strconv"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

type APIKeysHandler struct {
	Repo *storage.APIKeysRepo
}

func NewAPIKeysHandler(r *storage.APIKeysRepo) *APIKeysHandler { return &APIKeysHandler{Repo: r} }

type createKeyDTO struct {
	OwnerID string   `json:"ownerId" binding:"required"`
	Roles   []string `json:"roles"`
	Scopes  []string `json:"scopes"`
	Days    int      `json:"days"` // default from env
}

func (h *APIKeysHandler) Create(c *gin.Context) {
	var in createKeyDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if in.Days <= 0 { in.Days = envInt("APIKEY_TTL_DAYS", 90) }
	keyId := uuid.NewString()
	keySecret := randomBase32(32)
	secret := os.Getenv("APIKEY_HASH_SECRET")
	hash := hashKey(secret, keyId, keySecret)

	k := &domain.APIKey{
		KeyID: keyId,
		KeyHash: hash,
		OwnerID: in.OwnerID,
		Roles: in.Roles,
		Scopes: in.Scopes,
		ExpiresAt: time.Now().UTC().Add(time.Duration(in.Days)*24*time.Hour),
		Revoked: false,
	}
	if err := h.Repo.Insert(c, k); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"insert apikey"}); return
	}
	// Mostramos a chave **apenas agora**:
	c.JSON(http.StatusCreated, gin.H{
		"api_key": keyId + "." + keySecret,
		"expiresAt": k.ExpiresAt,
	})
}

type rotateDTO struct{ KeyID string `json:"keyId" binding:"required"` }
func (h *APIKeysHandler) Rotate(c *gin.Context) {
	var in rotateDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	// revoga a antiga e cria nova
	_ = h.Repo.Revoke(c, in.KeyID)
	// reusa Create flow (ou duplicar código)
	c.Set("rotating", true)
	h.Create(c)
}

type revokeDTO struct{ KeyID string `json:"keyId" binding:"required"` }
func (h *APIKeysHandler) Revoke(c *gin.Context) {
	var in revokeDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := h.Repo.Revoke(c, in.KeyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"revoke apikey"}); return
	}
	c.Status(http.StatusNoContent)
}

func hashKey(secret, keyId, keySecret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(keyId + "." + keySecret))
	return hex.EncodeToString(h.Sum(nil))
}
func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil { return i }
	}
	return def
}
func randomBase32(n int) string {
	const a = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	sb := strings.Builder{}
	for i := 0; i < n; i++ { sb.WriteByte(a[rand.Intn(len(a))]) }
	return sb.String()
}

