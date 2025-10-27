package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/cache"
	"go.mongodb.org/mongo-driver/mongo"
)

var startTime = time.Now()

type HealthHandler struct {
	Mongo *mongo.Client
	Redis interface{} // Pode ser *cache.RedisClient ou nil
}

func NewHealthHandler(m *mongo.Client, redis interface{}) *HealthHandler {
	return &HealthHandler{
		Mongo: m,
		Redis: redis,
	}
}

func (h *HealthHandler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	// ✅ Verificar MongoDB
	mongoStatus := false
	err := h.Mongo.Ping(ctx, nil)
	if err == nil {
		mongoStatus = true
	}

	// ✅ Verificar Redis (se disponível)
	redisStatus := false
	if h.Redis != nil {
		if redisClient, ok := h.Redis.(*cache.RedisClient); ok {
			// Tentar fazer um ping no Redis
			_, err := redisClient.Get(ctx, "health:check")
			// Se não der erro de conexão, Redis está up
			// (erro de "key not found" é OK, significa que conectou)
			if err == nil || err.Error() == "redis: nil" {
				redisStatus = true
			}
		}
	}

	// ✅ Status geral
	overallStatus := "ok"
	if !mongoStatus {
		overallStatus = "degraded" // MongoDB down é crítico
	}
	// Redis down não é crítico (graceful degradation)

	// ✅ Uptime
	uptime := time.Since(startTime)
	uptimeStr := fmt.Sprintf("%dd %dh %dm",
		int(uptime.Hours())/24,
		int(uptime.Hours())%24,
		int(uptime.Minutes())%60,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":    overallStatus,
		"version":   "1.0.0",
		"uptime":    uptimeStr,
		"timestamp": time.Now().UTC(),
		"services": gin.H{
			"mongodb": mongoStatus,
			"redis":   redisStatus,
		},
	})
}
