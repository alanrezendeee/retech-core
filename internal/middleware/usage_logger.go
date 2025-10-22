package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsageLogger struct {
	db *mongo.Database
}

func NewUsageLogger(db *mongo.Database) *UsageLogger {
	return &UsageLogger{db: db}
}

// Middleware loga cada request da API
func (ul *UsageLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capturar tempo de início
		startTime := time.Now()

		// Processar request
		c.Next()

		// Capturar tempo de fim
		responseTime := time.Since(startTime).Milliseconds()

		// Extrair informações
		apiKey, _ := c.Get("api_key")
		tenantID, _ := c.Get("tenant_id") // ✅ Corrigido: snake_case (não camelCase)

		// Se não tem API key, não loga (rota pública ou erro de auth)
		apiKeyStr, ok := apiKey.(string)
		if !ok || apiKeyStr == "" {
			return
		}

		tenantIDStr, ok := tenantID.(string)
		if !ok || tenantIDStr == "" {
			return // Sem tenant_id, não loga
		}

		now := time.Now()
		log := domain.APIUsageLog{
			APIKey:       apiKeyStr,
			TenantID:     tenantIDStr, // ✅ Usar a variável validada
			Endpoint:     c.Request.URL.Path,
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ResponseTime: responseTime,
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			Timestamp:    now,
			Date:         now.Format("2006-01-02"),
			Hour:         now.Hour(),
		}

		// Inserir em background (não bloquear response)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			coll := ul.db.Collection("api_usage_logs")
			_, _ = coll.InsertOne(ctx, log)
		}()
	}
}
