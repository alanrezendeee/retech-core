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
	"go.mongodb.org/mongo-driver/bson"
)

type TenantHandler struct {
	apikeys *storage.APIKeysRepo
	users   *storage.UsersRepo
	db      *storage.Mongo
}

func NewTenantHandler(apikeys *storage.APIKeysRepo, users *storage.UsersRepo, m *storage.Mongo) *TenantHandler {
	return &TenantHandler{
		apikeys: apikeys,
		users:   users,
		db:      m,
	}
}

// ListMyAPIKeys lista as API keys do tenant logado
// GET /me/apikeys
func (h *TenantHandler) ListMyAPIKeys(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Tenant ID não encontrado",
		})
		return
	}

	ctx := c.Request.Context()
	
	cursor, err := h.db.DB.Collection("api_keys").Find(ctx, bson.M{"ownerId": tenantID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao buscar API keys",
		})
		return
	}
	defer cursor.Close(ctx)

	var keys []bson.M
	if err := cursor.All(ctx, &keys); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao processar API keys",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apikeys": keys,
		"total":   len(keys),
	})
}

// CreateAPIKey cria uma nova API key para o tenant logado
// POST /me/apikeys
func (h *TenantHandler) CreateAPIKey(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Tenant ID não encontrado",
		})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	
	// ✅ Gerar API key REAL (mesmo algoritmo do admin)
	keyId := uuid.NewString()
	keySecret := randomBase32Tenant(32)
	secret := os.Getenv("APIKEY_HASH_SECRET")
	hash := hashKeyTenant(secret, keyId, keySecret)
	
	// Validade padrão: 90 dias
	days := envIntTenant("APIKEY_TTL_DAYS", 90)
	now := time.Now().UTC()

	// Usar domain.APIKey para garantir consistência
	k := &domain.APIKey{
		KeyID:     keyId,
		KeyHash:   hash,
		OwnerID:   tenantID,
		Scopes:    []string{"geo:read"},
		ExpiresAt: now.Add(time.Duration(days) * 24 * time.Hour),
		Revoked:   false,
		CreatedAt: now,
	}

	if err := h.apikeys.Insert(ctx, k); err != nil {
		fmt.Printf("❌ Erro ao inserir API key: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao criar API key",
		})
		return
	}

	fmt.Printf("✅ API Key criada com sucesso para tenant %s\n", tenantID)
	fmt.Printf("   keyId: %s\n", keyId)
	fmt.Printf("   Full key: %s.%s\n", keyId, keySecret)

	// ⚠️ IMPORTANTE: Retornar a chave COMPLETA apenas agora!
	c.JSON(http.StatusCreated, gin.H{
		"key":       keyId + "." + keySecret, // ← Chave completa!
		"expiresAt": k.ExpiresAt,
		"name":      req.Name,
	})
}

// DeleteAPIKey deleta uma API key do tenant logado
// DELETE /me/apikeys/:id
func (h *TenantHandler) DeleteAPIKey(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	keyID := c.Param("id")

	ctx := c.Request.Context()

	// Verificar se a key pertence ao tenant
	var apikey bson.M
	err := h.db.DB.Collection("api_keys").FindOne(ctx, bson.M{
		"keyId":   keyID,
		"ownerId": tenantID,
	}).Decode(&apikey)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://retech-core/errors/not-found",
			"title":  "Not Found",
			"status": http.StatusNotFound,
			"detail": "API key não encontrada",
		})
		return
	}

	// Revogar (soft delete)
	_, err = h.db.DB.Collection("api_keys").UpdateOne(ctx, bson.M{"keyId": keyID}, bson.M{
		"$set": bson.M{"revoked": true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao revogar API key",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key revogada com sucesso",
	})
}

// GetMyUsage retorna uso da API do tenant logado
// GET /me/usage
func (h *TenantHandler) GetMyUsage(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://retech-core/errors/unauthorized",
			"title":  "Unauthorized",
			"status": http.StatusUnauthorized,
			"detail": "Tenant ID não encontrado",
		})
		return
	}

	ctx := c.Request.Context()

	// Total de requests
	totalRequests, _ := h.db.DB.Collection("api_usage_logs").CountDocuments(ctx, bson.M{"tenantId": tenantID})

	// Requests hoje
	today := time.Now().Format("2006-01-02")
	requestsToday, _ := h.db.DB.Collection("api_usage_logs").CountDocuments(ctx, bson.M{
		"tenantId": tenantID,
		"date":     today,
	})

	// Requests mês
	startOfMonth := time.Now().Format("2006-01")
	requestsMonth, _ := h.db.DB.Collection("api_usage_logs").CountDocuments(ctx, bson.M{
		"tenantId": tenantID,
		"date":     bson.M{"$regex": "^" + startOfMonth},
	})

	// Limite diário (free tier)
	dailyLimit := int64(1000)
	remaining := dailyLimit - requestsToday
	if remaining < 0 {
		remaining = 0
	}

	// Por dia (últimos 7 dias)
	pipeline := []bson.M{
		{"$match": bson.M{
			"tenantId":  tenantID,
			"timestamp": bson.M{"$gte": time.Now().AddDate(0, 0, -7)},
		}},
		{"$group": bson.M{
			"_id":   "$date",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"_id": 1}},
	}

	cursor, _ := h.db.DB.Collection("api_usage_logs").Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)

	var byDay []bson.M
	cursor.All(ctx, &byDay)

	// Por endpoint
	pipelineEndpoints := []bson.M{
		{"$match": bson.M{"tenantId": tenantID}},
		{"$group": bson.M{
			"_id":   "$endpoint",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor2, _ := h.db.DB.Collection("api_usage_logs").Aggregate(ctx, pipelineEndpoints)
	defer cursor2.Close(ctx)

	var byEndpoint []bson.M
	cursor2.All(ctx, &byEndpoint)

	c.JSON(http.StatusOK, gin.H{
		"totalRequests":  totalRequests,
		"requestsToday":  requestsToday,
		"requestsMonth":  requestsMonth,
		"dailyLimit":     dailyLimit,
		"remaining":      remaining,
		"percentageUsed": float64(requestsToday) / float64(dailyLimit) * 100,
		"byDay":          byDay,
		"byEndpoint":     byEndpoint,
	})
}

// ========================================
// Funções auxiliares (mesmas do apikey.go)
// ========================================

func hashKeyTenant(secret, keyId, keySecret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(keyId + "." + keySecret))
	return hex.EncodeToString(h.Sum(nil))
}

func envIntTenant(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func randomBase32Tenant(n int) string {
	const a = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		sb.WriteByte(a[rand.Intn(len(a))])
	}
	return sb.String()
}
