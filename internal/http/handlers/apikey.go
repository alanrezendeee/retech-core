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
	"github.com/theretech/retech-core/internal/auth"
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

	// Validar scopes
	if err := auth.ValidateAPIKeyScopes(in.Scopes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation",
			"title":  "Scopes Inválidos",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
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

// Rotate rotaciona uma API key (revoga antiga e cria nova)
// POST /admin/apikeys/rotate
func (h *APIKeysHandler) Rotate(c *gin.Context) {
	var req struct {
		KeyID string `json:"keyId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Erro de validação",
			"status": http.StatusBadRequest,
			"detail": "keyId é obrigatório",
		})
		return
	}

	// Buscar a API key existente
	existingKey, err := h.Repo.ByKeyIDAny(c, req.KeyID)
	if err != nil {
		fmt.Printf("❌ Erro ao buscar API key: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro interno",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao buscar API key",
		})
		return
	}
	if existingKey == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://retech-core/errors/not-found",
			"title":  "API Key não encontrada",
			"status": http.StatusNotFound,
			"detail": "API key não existe",
		})
		return
	}

	// Buscar tenant (para log)
	tenant, _ := h.TenantsRepo.ByTenantID(c, existingKey.OwnerID)
	tenantName := "Unknown"
	if tenant != nil {
		tenantName = tenant.Name
	}

	// Revogar a chave antiga
	if err := h.Repo.Revoke(c, req.KeyID); err != nil {
		fmt.Printf("❌ Erro ao revogar API key: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro interno",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao revogar chave antiga",
		})
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
		fmt.Printf("❌ Erro ao criar nova API key: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro interno",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao criar nova chave",
		})
		return
	}

	// Log da atividade
	utils.LogActivity(
		c,
		h.ActivityRepo,
		"apikey.rotated",
		domain.ActionUpdate,
		utils.BuildActorFromContext(c),
		domain.Resource{
			Type: domain.ResourceTypeAPIKey,
			ID:   keyId,
			Name: fmt.Sprintf("API Key de %s (rotacionada)", tenantName),
		},
		map[string]interface{}{
			"tenantId":  existingKey.OwnerID,
			"oldKeyId":  req.KeyID,
			"newKeyId":  keyId,
			"expiresAt": newKey.ExpiresAt,
		},
	)

	fmt.Printf("✅ API Key rotacionada com sucesso!\n")
	fmt.Printf("   Old keyId: %s\n", req.KeyID)
	fmt.Printf("   New keyId: %s\n", keyId)
	fmt.Printf("   Tenant: %s\n", tenantName)

	// Retornar a nova chave (apenas uma vez!)
	c.JSON(http.StatusCreated, gin.H{
		"api_key":   keyId + "." + keySecret,
		"expiresAt": newKey.ExpiresAt,
		"message":   "API key rotacionada com sucesso",
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
