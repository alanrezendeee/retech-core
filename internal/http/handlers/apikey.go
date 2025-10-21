package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"github.com/theretech/retech-core/internal/utils"
)

type APIKeysHandler struct {
	Repo         *storage.APIKeysRepo
	TenantsRepo  *storage.TenantsRepo
	ActivityRepo *storage.ActivityLogsRepo
}

func NewAPIKeysHandler(r *storage.APIKeysRepo, t *storage.TenantsRepo, a *storage.ActivityLogsRepo) *APIKeysHandler {
	return &APIKeysHandler{
		Repo:         r,
		TenantsRepo:  t,
		ActivityRepo: a,
	}
}

type createKeyDTO struct {
	OwnerID string   `json:"ownerId" binding:"required"`
	Scopes  []string `json:"scopes"`
	Days    int      `json:"days"` // default from env
}

type rotateKeyDTO struct {
	KeyID string `json:"keyId" binding:"required"`
}

func (h *APIKeysHandler) Create(c *gin.Context) {
	var in createKeyDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Valida se o tenant existe
	tenant, err := h.TenantsRepo.ByTenantID(c, in.OwnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check tenant"})
		return
	}
	if tenant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	if !tenant.Active {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant is inactive"})
		return
	}

	if in.Days <= 0 {
		in.Days = envInt("APIKEY_TTL_DAYS", 90)
	}
	keyId := uuid.NewString()
	keySecret := randomBase32(32)
	secret := os.Getenv("APIKEY_HASH_SECRET")
	hash := hashKey(secret, keyId, keySecret)

	k := &domain.APIKey{
		KeyID:     keyId,
		KeyHash:   hash,
		OwnerID:   in.OwnerID,
		Scopes:    in.Scopes,
		ExpiresAt: time.Now().UTC().Add(time.Duration(in.Days) * 24 * time.Hour),
		Revoked:   false,
	}
	if err := h.Repo.Insert(c, k); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert apikey"})
		return
	}

	// Log da atividade
	utils.LogActivity(
		c,
		h.ActivityRepo,
		domain.ActivityTypeAPIKeyCreated,
		domain.ActionCreate,
		utils.BuildActorFromContext(c),
		domain.Resource{
			Type: domain.ResourceTypeAPIKey,
			ID:   keyId,
			Name: fmt.Sprintf("API Key para %s", tenant.Name),
		},
		map[string]interface{}{
			"tenantId":  in.OwnerID,
			"expiresAt": k.ExpiresAt,
			"scopes":    in.Scopes,
		},
	)

	// Mostramos a chave **apenas agora**:
	c.JSON(http.StatusCreated, gin.H{
		"api_key":   keyId + "." + keySecret,
		"expiresAt": k.ExpiresAt,
	})
}

func (h *APIKeysHandler) rotateExistingKey(c *gin.Context, in rotateKeyDTO) {
	// Buscar a API key existente
	existingKey, err := h.Repo.ByKeyIDAny(c, in.KeyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get existing key"})
		return
	}
	if existingKey == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "api key not found"})
		return
	}

	// Revogar a chave antiga
	if err := h.Repo.Revoke(c, in.KeyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "revoke old key"})
		return
	}

	// Criar nova chave com os mesmos dados do tenant
	keyId := uuid.NewString()
	keySecret := randomBase32(32)
	secret := os.Getenv("APIKEY_HASH_SECRET")
	hash := hashKey(secret, keyId, keySecret)

	// Calcular nova data de expiração (mesmo período da chave original)
	now := time.Now().UTC()
	days := int(existingKey.ExpiresAt.Sub(existingKey.CreatedAt).Hours() / 24)
	if days <= 0 {
		days = envInt("APIKEY_TTL_DAYS", 90)
	}

	newKey := &domain.APIKey{
		KeyID:     keyId,
		KeyHash:   hash,
		OwnerID:   existingKey.OwnerID,
		Scopes:    existingKey.Scopes,
		ExpiresAt: now.Add(time.Duration(days) * 24 * time.Hour),
		Revoked:   false,
	}

	if err := h.Repo.Insert(c, newKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert new key"})
		return
	}

	// Retornar a nova chave
	c.JSON(http.StatusCreated, gin.H{
		"api_key":   keyId + "." + keySecret,
		"expiresAt": newKey.ExpiresAt,
		"message":   "API key rotated successfully",
	})
}

// Rotate - Nova implementação do zero
type rotateRequest struct {
	KeyID string `json:"keyId" binding:"required"`
}

func (h *APIKeysHandler) Rotate(c *gin.Context) {
	// Debug: verificar se o request está chegando
	fmt.Printf("Rotate handler chamado!\n")

	// Debug: verificar headers
	fmt.Printf("Headers: %v\n", c.Request.Header)

	// Debug: verificar body
	body, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Erro ao ler body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}
	fmt.Printf("Body: %s\n", string(body))

	// Resposta simples
	c.JSON(http.StatusOK, gin.H{
		"message": "Rotate endpoint working - debug mode",
		"status":  "success",
	})
}

func (h *APIKeysHandler) RotateNew(c *gin.Context) {
	// Handler completamente novo
	c.String(http.StatusOK, "RotateNew endpoint working")
}

func (h *APIKeysHandler) RotateTest(c *gin.Context) {
	// Implementação completa de rotação (solução alternativa)
	var req rotateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// Buscar a API key existente
	existingKey, err := h.Repo.ByKeyIDAny(c, req.KeyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get existing key"})
		return
	}
	if existingKey == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "api key not found"})
		return
	}

	// Revogar a chave antiga
	if err := h.Repo.Revoke(c, req.KeyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "revoke old key"})
		return
	}

	// Criar nova chave com os mesmos dados do tenant
	keyId := uuid.NewString()
	keySecret := randomBase32(32)
	secret := os.Getenv("APIKEY_HASH_SECRET")
	hash := hashKey(secret, keyId, keySecret)

	// Calcular nova data de expiração (mesmo período da chave original)
	now := time.Now().UTC()
	days := int(existingKey.ExpiresAt.Sub(existingKey.CreatedAt).Hours() / 24)
	if days <= 0 {
		days = envInt("APIKEY_TTL_DAYS", 90)
	}

	newKey := &domain.APIKey{
		KeyID:     keyId,
		KeyHash:   hash,
		OwnerID:   existingKey.OwnerID,
		Scopes:    existingKey.Scopes,
		ExpiresAt: now.Add(time.Duration(days) * 24 * time.Hour),
		Revoked:   false,
	}

	if err := h.Repo.Insert(c, newKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert new key"})
		return
	}

	// Retornar a nova chave
	c.JSON(http.StatusCreated, gin.H{
		"api_key":   keyId + "." + keySecret,
		"expiresAt": newKey.ExpiresAt,
		"message":   "API key rotated successfully",
	})
}

type revokeDTO struct {
	KeyID string `json:"keyId" binding:"required"`
}

func (h *APIKeysHandler) Revoke(c *gin.Context) {
	var in revokeDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar API Key antes de revogar (para log)
	apiKey, _ := h.Repo.ByKeyIDAny(c, in.KeyID)

	if err := h.Repo.Revoke(c, in.KeyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "revoke apikey"})
		return
	}

	// Log da atividade
	if apiKey != nil {
		tenant, _ := h.TenantsRepo.ByTenantID(c, apiKey.OwnerID)
		tenantName := "Unknown"
		if tenant != nil {
			tenantName = tenant.Name
		}

		utils.LogActivity(
			c,
			h.ActivityRepo,
			domain.ActivityTypeAPIKeyRevoked,
			domain.ActionRevoke,
			utils.BuildActorFromContext(c),
			domain.Resource{
				Type: domain.ResourceTypeAPIKey,
				ID:   in.KeyID,
				Name: fmt.Sprintf("API Key de %s", tenantName),
			},
			map[string]interface{}{
				"tenantId": apiKey.OwnerID,
			},
		)
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
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
func randomBase32(n int) string {
	const a = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteByte(a[rand.Intn(len(a))])
	}
	return sb.String()
}
