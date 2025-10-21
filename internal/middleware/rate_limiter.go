package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RateLimiter struct {
	db     *mongo.Database
	config domain.RateLimitConfig
}

func NewRateLimiter(db *mongo.Database, config domain.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		db:     db,
		config: config,
	}
}

// Middleware aplica rate limiting baseado em API Key
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair API Key do contexto (já foi validada pelo middleware de API Key)
		apiKeyValue, exists := c.Get("api_key")
		if !exists {
			// Se não tem API key, deixa passar (rota pública)
			c.Next()
			return
		}

		apiKey := apiKeyValue.(string)
		today := time.Now().Format("2006-01-02")

		// Buscar ou criar registro de rate limit
		coll := rl.db.Collection("rate_limits")
		ctx := context.Background()

		var rateLimit domain.RateLimit
		err := coll.FindOne(ctx, bson.M{
			"apiKey": apiKey,
			"date":   today,
		}).Decode(&rateLimit)

		if err == mongo.ErrNoDocuments {
			// Criar novo registro
			rateLimit = domain.RateLimit{
				APIKey:    apiKey,
				Date:      today,
				Count:     0,
				LastReset: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		// Verificar limite diário
		if rateLimit.Count >= rl.config.RequestsPerDay {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.RequestsPerDay))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", getNextDayTimestamp())

			c.JSON(http.StatusTooManyRequests, gin.H{
				"type":   "https://retech-core/errors/rate-limit-exceeded",
				"title":  "Rate Limit Exceeded",
				"status": http.StatusTooManyRequests,
				"detail": fmt.Sprintf("Limite de %d requests por dia excedido", rl.config.RequestsPerDay),
			})
			c.Abort()
			return
		}

		// Incrementar contador
		rateLimit.Count++
		rateLimit.UpdatedAt = time.Now()

		// Atualizar ou inserir
		opts := options.Update().SetUpsert(true)
		_, err = coll.UpdateOne(ctx, bson.M{
			"apiKey": apiKey,
			"date":   today,
		}, bson.M{
			"$set": rateLimit,
		}, opts)

		if err != nil {
			// Log erro mas não bloqueia
			fmt.Printf("Erro ao atualizar rate limit: %v\n", err)
		}

		// Adicionar headers de rate limit
		remaining := rl.config.RequestsPerDay - rateLimit.Count
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.config.RequestsPerDay))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", getNextDayTimestamp())

		c.Next()
	}
}

// getNextDayTimestamp retorna timestamp Unix do próximo dia (meia-noite)
func getNextDayTimestamp() string {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return fmt.Sprintf("%d", tomorrow.Unix())
}

