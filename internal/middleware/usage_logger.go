package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/utils"
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

		endpoint := c.Request.URL.Path

		// Extrair nome da API do endpoint
		apiName := extractAPIName(endpoint)

		// ✅ Usar timezone de Brasília para data
		nowBrasilia := utils.GetBrasiliaTime()

		log := domain.APIUsageLog{
			APIKey:       apiKeyStr,
			TenantID:     tenantIDStr, // ✅ Usar a variável validada
			APIName:      apiName,     // ✅ Nome da API (cep, geografia, etc)
			Endpoint:     endpoint,
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			ResponseTime: responseTime,
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			Timestamp:    time.Now(),                       // ✅ UTC para ordenação
			Date:         nowBrasilia.Format("2006-01-02"), // ✅ Brasília para agrupamento
			Hour:         nowBrasilia.Hour(),               // ✅ Brasília para hora local
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

// extractAPIName extrai o nome da API a partir do endpoint
// Exemplos:
// /geo/ufs → "geografia"
// /cep/01310-100 → "cep"
// /cnpj/12345678000190 → "cnpj"
func extractAPIName(endpoint string) string {
	// Remover query string se houver
	if idx := strings.Index(endpoint, "?"); idx != -1 {
		endpoint = endpoint[:idx]
	}

	// Split por /
	parts := strings.Split(strings.Trim(endpoint, "/"), "/")
	if len(parts) == 0 {
		return "unknown"
	}

	// Primeiro segmento é o nome da API
	apiName := parts[0]

	// Mapear alguns casos especiais
	switch apiName {
	case "geo":
		return "geografia"
	case "cep":
		return "cep"
	case "cnpj":
		return "cnpj"
	case "cpf":
		return "cpf"
	case "fipe":
		return "fipe"
	case "moedas":
		return "moedas"
	case "bancos":
		return "bancos"
	case "feriados":
		return "feriados"
	default:
		return apiName
	}
}
