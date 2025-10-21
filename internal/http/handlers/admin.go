package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
)

type AdminHandler struct {
	tenants    *storage.TenantsRepo
	apikeys    *storage.APIKeysRepo
	users      *storage.UsersRepo
	db         *storage.Mongo
}

func NewAdminHandler(tenants *storage.TenantsRepo, apikeys *storage.APIKeysRepo, users *storage.UsersRepo, m *storage.Mongo) *AdminHandler {
	return &AdminHandler{
		tenants: tenants,
		apikeys: apikeys,
		users:   users,
		db:      m,
	}
}

// GetStats retorna estatísticas globais do sistema
// GET /admin/stats
func (h *AdminHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	// Contar tenants
	totalTenants, _ := h.db.DB.Collection("tenants").CountDocuments(ctx, bson.M{})
	activeTenants, _ := h.db.DB.Collection("tenants").CountDocuments(ctx, bson.M{"active": true})

	// Contar API keys
	totalKeys, _ := h.db.DB.Collection("api_keys").CountDocuments(ctx, bson.M{})
	activeKeys, _ := h.db.DB.Collection("api_keys").CountDocuments(ctx, bson.M{"revoked": false})

	// Contar usuários
	totalUsers, _ := h.db.DB.Collection("users").CountDocuments(ctx, bson.M{})

	// Contar requests hoje
	today := time.Now().Format("2006-01-02")
	requestsToday, _ := h.db.DB.Collection("api_usage_logs").CountDocuments(ctx, bson.M{"date": today})

	// Contar requests mês
	startOfMonth := time.Now().Format("2006-01")
	requestsMonth, _ := h.db.DB.Collection("api_usage_logs").CountDocuments(ctx, bson.M{
		"date": bson.M{"$regex": "^" + startOfMonth},
	})

	c.JSON(http.StatusOK, gin.H{
		"totalTenants":   totalTenants,
		"activeTenants":  activeTenants,
		"totalAPIKeys":   totalKeys,
		"activeAPIKeys":  activeKeys,
		"totalUsers":     totalUsers,
		"requestsToday":  requestsToday,
		"requestsMonth":  requestsMonth,
		"systemUptime":   time.Since(time.Now()).String(), // TODO: real uptime
		"timestamp":      time.Now(),
	})
}

// ListAllAPIKeys lista todas as API keys do sistema
// GET /admin/apikeys
func (h *AdminHandler) ListAllAPIKeys(c *gin.Context) {
	ctx := c.Request.Context()

	cursor, err := h.db.DB.Collection("api_keys").Find(ctx, bson.M{})
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

// GetUsage retorna dados de uso da API
// GET /admin/usage
func (h *AdminHandler) GetUsage(c *gin.Context) {
	ctx := c.Request.Context()

	// Uso por dia (últimos 7 dias)
	pipeline := []bson.M{
		{"$match": bson.M{
			"timestamp": bson.M{"$gte": time.Now().AddDate(0, 0, -7)},
		}},
		{"$group": bson.M{
			"_id":   "$date",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"_id": 1}},
	}

	cursor, err := h.db.DB.Collection("api_usage_logs").Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var byDay []bson.M
	cursor.All(ctx, &byDay)

	// Uso por endpoint (top 10)
	pipelineEndpoints := []bson.M{
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

	// Uso por tenant (top 10)
	pipelineTenants := []bson.M{
		{"$group": bson.M{
			"_id":   "$tenantId",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor3, _ := h.db.DB.Collection("api_usage_logs").Aggregate(ctx, pipelineTenants)
	defer cursor3.Close(ctx)

	var byTenant []bson.M
	cursor3.All(ctx, &byTenant)

	c.JSON(http.StatusOK, gin.H{
		"byDay":      byDay,
		"byEndpoint": byEndpoint,
		"byTenant":   byTenant,
	})
}

