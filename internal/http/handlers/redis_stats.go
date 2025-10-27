package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/cache"
)

type RedisStatsHandler struct {
	redis interface{}
}

func NewRedisStatsHandler(redis interface{}) *RedisStatsHandler {
	return &RedisStatsHandler{
		redis: redis,
	}
}

// GetStats retorna estatísticas do Redis
// GET /admin/cache/redis/stats
func (h *RedisStatsHandler) GetStats(c *gin.Context) {
	ctx := context.Background()

	if h.redis == nil {
		c.JSON(http.StatusOK, gin.H{
			"connected":    false,
			"totalKeys":    0,
			"cepKeys":      0,
			"cnpjKeys":     0,
			"geoKeys":      0,
			"memoryUsedMB": 0,
			"message":      "Redis não está disponível",
		})
		return
	}

	redisClient, ok := h.redis.(*cache.RedisClient)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"connected":    false,
			"totalKeys":    0,
			"cepKeys":      0,
			"cnpjKeys":     0,
			"geoKeys":      0,
			"memoryUsedMB": 0,
			"message":      "Redis client inválido",
		})
		return
	}

	// Contar keys por tipo
	cepKeys, _ := redisClient.Keys(ctx, "cep:*")
	cnpjKeys, _ := redisClient.Keys(ctx, "cnpj:*")
	geoKeys, _ := redisClient.Keys(ctx, "geo:*")

	totalKeys := len(cepKeys) + len(cnpjKeys) + len(geoKeys)

	// Pegar info de memória
	info, err := redisClient.Info(ctx, "memory")
	memoryUsedMB := 0.0

	if err == nil {
		// Extrair used_memory_human do INFO
		// Exemplo: used_memory_human:1.23M
		var usedMemoryBytes int64
		fmt.Sscanf(info, "# Memory\nused_memory:%d", &usedMemoryBytes)
		memoryUsedMB = float64(usedMemoryBytes) / 1024.0 / 1024.0
	}

	c.JSON(http.StatusOK, gin.H{
		"connected":    true,
		"totalKeys":    totalKeys,
		"cepKeys":      len(cepKeys),
		"cnpjKeys":     len(cnpjKeys),
		"geoKeys":      len(geoKeys),
		"memoryUsedMB": fmt.Sprintf("%.2f", memoryUsedMB),
		"message":      "Redis operacional",
	})
}

// ClearAll limpa TODOS os caches do Redis
// DELETE /admin/cache/redis
func (h *RedisStatsHandler) ClearAll(c *gin.Context) {
	ctx := context.Background()

	if h.redis == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Redis não está disponível",
		})
		return
	}

	redisClient, ok := h.redis.(*cache.RedisClient)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Redis client inválido",
		})
		return
	}

	// Limpar todos os caches
	err := redisClient.FlushDB(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao limpar Redis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Todos os caches do Redis foram limpos",
		"deletedCount": "all",
	})
}

// ClearCEP limpa apenas cache de CEP do Redis
// DELETE /admin/cache/redis/cep
func (h *RedisStatsHandler) ClearCEP(c *gin.Context) {
	ctx := context.Background()

	if h.redis == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Redis não está disponível",
		})
		return
	}

	redisClient, ok := h.redis.(*cache.RedisClient)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Redis client inválido",
		})
		return
	}

	// Buscar todas as keys de CEP
	keys, err := redisClient.Keys(ctx, "cep:*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar keys de CEP",
		})
		return
	}

	// Deletar todas
	if len(keys) > 0 {
		err = redisClient.Del(ctx, keys...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao deletar keys de CEP",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Cache de CEP do Redis foi limpo",
		"deletedCount": len(keys),
	})
}

// ClearCNPJ limpa apenas cache de CNPJ do Redis
// DELETE /admin/cache/redis/cnpj
func (h *RedisStatsHandler) ClearCNPJ(c *gin.Context) {
	ctx := context.Background()

	if h.redis == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Redis não está disponível",
		})
		return
	}

	redisClient, ok := h.redis.(*cache.RedisClient)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Redis client inválido",
		})
		return
	}

	// Buscar todas as keys de CNPJ
	keys, err := redisClient.Keys(ctx, "cnpj:*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar keys de CNPJ",
		})
		return
	}

	// Deletar todas
	if len(keys) > 0 {
		err = redisClient.Del(ctx, keys...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao deletar keys de CNPJ",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Cache de CNPJ do Redis foi limpo",
		"deletedCount": len(keys),
	})
}
