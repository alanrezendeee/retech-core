package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlaygroundRateLimiter struct {
	db       *mongo.Database
	settings *storage.SettingsRepo
}

func NewPlaygroundRateLimiter(db *mongo.Database, settings *storage.SettingsRepo) *PlaygroundRateLimiter {
	return &PlaygroundRateLimiter{
		db:       db,
		settings: settings,
	}
}

// PlaygroundRateLimit representa rate limit por IP + API Key demo
type PlaygroundRateLimit struct {
	ID             string    `bson:"_id" json:"id"`
	IPAddress      string    `bson:"ipAddress" json:"ipAddress"`
	APIKey         string    `bson:"apiKey" json:"apiKey"`
	Date           string    `bson:"date" json:"date"`     // YYYY-MM-DD
	Minute         string    `bson:"minute" json:"minute"` // YYYY-MM-DD HH:MM
	CountPerDay    int64     `bson:"countPerDay" json:"countPerDay"`
	CountPerMinute int64     `bson:"countPerMinute" json:"countPerMinute"`
	LastRequest    time.Time `bson:"lastRequest" json:"lastRequest"`
	UpdatedAt      time.Time `bson:"updatedAt" json:"updatedAt"`
}

// GlobalPlaygroundLimit representa limite global da API Key demo
type GlobalPlaygroundLimit struct {
	ID            string    `bson:"_id" json:"id"`
	APIKey        string    `bson:"apiKey" json:"apiKey"`
	Date          string    `bson:"date" json:"date"` // YYYY-MM-DD
	TotalRequests int64     `bson:"totalRequests" json:"totalRequests"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
}

// Middleware aplica rate limiting específico para playground
func (prl *PlaygroundRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ✅ 1. VERIFICAR SE É ROTA DO PLAYGROUND
		if !isPlaygroundRoute(c.Request.URL.Path) {
			c.Next()
			return
		}

		// ✅ 2. VERIFICAR SE PLAYGROUND ESTÁ HABILITADO E TEM API KEY VÁLIDA
		ctx := context.Background()
		settings, err := prl.settings.Get(ctx)
		if err != nil || settings == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"type":   "https://retech-core/errors/playground-disabled",
				"title":  "Playground Indisponível",
				"status": http.StatusServiceUnavailable,
				"detail": "Erro ao carregar configurações do playground",
			})
			c.Abort()
			return
		}

		// Verificar se playground está habilitado
		if !settings.Playground.Enabled {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"type":   "https://retech-core/errors/playground-disabled",
				"title":  "Playground Indisponível",
				"status": http.StatusServiceUnavailable,
				"detail": "O playground está temporariamente desabilitado",
			})
			c.Abort()
			return
		}

		// Verificar se tem API Key configurada
		if settings.Playground.APIKey == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"type":   "https://retech-core/errors/playground-not-configured",
				"title":  "Playground Não Configurado",
				"status": http.StatusServiceUnavailable,
				"detail": "API Key demo não configurada. Entre em contato com o administrador.",
			})
			c.Abort()
			return
		}

		// ✅ 3. EXTRAIR IP DO CLIENTE
		clientIP := getClientIP(c)
		fmt.Printf("🔒 [PLAYGROUND SECURITY] IP: %s | Path: %s\n", clientIP, c.Request.URL.Path)

		// ✅ 4. VERIFICAR SE É API KEY DEMO
		apiKeyValue, exists := c.Get("api_key")
		if !exists {
			fmt.Println("⚠️  [PLAYGROUND SECURITY] Nenhuma API key no contexto")
			c.JSON(http.StatusUnauthorized, gin.H{
				"type":   "https://retech-core/errors/unauthorized",
				"title":  "API Key Obrigatória",
				"status": http.StatusUnauthorized,
				"detail": "Esta rota requer uma API Key válida",
			})
			c.Abort()
			return
		}

		apiKey := apiKeyValue.(string)

		// ✅ 4. VERIFICAR SE É API KEY DEMO (começa com rtc_demo_)
		if !isDemoAPIKey(apiKey) {
			fmt.Println("✅ [PLAYGROUND SECURITY] API Key normal, passando...")
			c.Next()
			return
		}

		fmt.Printf("🔒 [PLAYGROUND SECURITY] API Key Demo detectada: %s...\n", apiKey[:20])

		// ✅ 5. USAR CONFIGURAÇÕES JÁ CARREGADAS (não buscar de novo)
		playgroundConfig := settings.Playground
		rateLimit := playgroundConfig.RateLimit

		// 🔍 DEBUG: Log rate limits configurados
		fmt.Printf("📊 [PLAYGROUND SECURITY] Rate Limits Configurados:\n")
		fmt.Printf("   - Requests/Dia: %d\n", rateLimit.RequestsPerDay)
		fmt.Printf("   - Requests/Min: %d\n", rateLimit.RequestsPerMinute)

		// ✅ 6. APLICAR RATE LIMITING POR IP
		if !prl.checkIPRateLimit(ctx, clientIP, apiKey, rateLimit, c) {
			return // Rate limit excedido, response já enviado
		}

		// ✅ 7. APLICAR RATE LIMITING GLOBAL DA API KEY DEMO
		if !prl.checkGlobalRateLimit(ctx, apiKey, rateLimit, c) {
			return // Rate limit global excedido, response já enviado
		}

		// ✅ 8. APLICAR THROTTLING (delay mínimo entre requests)
		prl.applyThrottling(ctx, clientIP, apiKey, c)

		// ✅ 9. INCREMENTAR CONTADORES
		prl.incrementCounters(ctx, clientIP, apiKey, c)

		fmt.Printf("✅ [PLAYGROUND SECURITY] Request permitida para IP %s\n", clientIP)
		c.Next()
	}
}

// checkIPRateLimit verifica limite por IP
func (prl *PlaygroundRateLimiter) checkIPRateLimit(ctx context.Context, clientIP, apiKey string, rateLimit domain.RateLimitConfig, c *gin.Context) bool {
	now := time.Now()
	// Usar timezone Brasília
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrasilia := now.In(loc)
	today := nowBrasilia.Format("2006-01-02")
	currentMinute := nowBrasilia.Format("2006-01-02 15:04")

	// Buscar rate limit por IP
	coll := prl.db.Collection("playground_rate_limits")

	var ipLimit PlaygroundRateLimit
	err := coll.FindOne(ctx, bson.M{
		"ipAddress": clientIP,
		"apiKey":    apiKey,
		"date":      today,
	}).Decode(&ipLimit)

	if err == mongo.ErrNoDocuments {
		// Primeira request do IP hoje
		ipLimit = PlaygroundRateLimit{
			ID:             fmt.Sprintf("%s_%s_%s", clientIP, apiKey, today),
			IPAddress:      clientIP,
			APIKey:         apiKey,
			Date:           today,
			Minute:         currentMinute,
			CountPerDay:    0,
			CountPerMinute: 0,
			LastRequest:    now,
			UpdatedAt:      now,
		}
	}

	// 🔍 DEBUG: Log estado atual
	fmt.Printf("📊 [PLAYGROUND SECURITY] IP: %s | Count: %d/%d (dia) | %d/%d (min)\n",
		clientIP, ipLimit.CountPerDay, rateLimit.RequestsPerDay,
		ipLimit.CountPerMinute, rateLimit.RequestsPerMinute)

	// Verificar limite diário por IP
	if ipLimit.CountPerDay >= rateLimit.RequestsPerDay {
		fmt.Printf("🚫 [PLAYGROUND SECURITY] Limite diário por IP excedido: %s (%d >= %d)\n",
			clientIP, ipLimit.CountPerDay, rateLimit.RequestsPerDay)

		c.Header("X-RateLimit-Limit-Day", fmt.Sprintf("%d", rateLimit.RequestsPerDay))
		c.Header("X-RateLimit-Remaining-Day", "0")
		c.Header("X-RateLimit-Reset-Day", getNextDayTimestampPlayground())

		c.JSON(http.StatusTooManyRequests, gin.H{
			"type":   "https://retech-core/errors/rate-limit-exceeded",
			"title":  "Rate Limit Exceeded",
			"status": http.StatusTooManyRequests,
			"detail": fmt.Sprintf("Limite de %d requests por dia por IP excedido. Tente novamente amanhã.", rateLimit.RequestsPerDay),
		})
		c.Abort()
		return false
	}

	// ✅ Verificar se mudou de minuto (reset automático)
	if ipLimit.Minute != currentMinute {
		// Novo minuto! Resetar contador
		ipLimit.CountPerMinute = 0
		ipLimit.Minute = currentMinute
		fmt.Printf("🔄 [PLAYGROUND SECURITY] Novo minuto detectado para IP %s: %s → %s (count resetado)\n",
			clientIP, ipLimit.Minute, currentMinute)
	}

	// Verificar limite por minuto por IP
	if ipLimit.CountPerMinute >= rateLimit.RequestsPerMinute {
		fmt.Printf("🚫 [PLAYGROUND SECURITY] Limite por minuto por IP excedido: %s (%d >= %d)\n",
			clientIP, ipLimit.CountPerMinute, rateLimit.RequestsPerMinute)

		c.Header("X-RateLimit-Limit-Minute", fmt.Sprintf("%d", rateLimit.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining-Minute", "0")
		c.Header("X-RateLimit-Reset-Minute", getNextMinuteTimestampPlayground())

		c.JSON(http.StatusTooManyRequests, gin.H{
			"type":   "https://retech-core/errors/rate-limit-exceeded",
			"title":  "Rate Limit Exceeded (Per Minute)",
			"status": http.StatusTooManyRequests,
			"detail": fmt.Sprintf("Limite de %d requests por minuto por IP excedido. Tente novamente em alguns segundos.", rateLimit.RequestsPerMinute),
		})
		c.Abort()
		return false
	}

	return true
}

// checkGlobalRateLimit verifica limite global da API Key demo
func (prl *PlaygroundRateLimiter) checkGlobalRateLimit(ctx context.Context, apiKey string, rateLimit domain.RateLimitConfig, c *gin.Context) bool {
	now := time.Now()
	// Usar timezone Brasília
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrasilia := now.In(loc)
	today := nowBrasilia.Format("2006-01-02")

	// Limite global: 10x o limite por IP (exemplo: 100 IPs × 10 req/dia = 1000 req/dia total)
	globalLimit := rateLimit.RequestsPerDay * 10

	coll := prl.db.Collection("playground_global_limits")

	var globalLimitRecord GlobalPlaygroundLimit
	err := coll.FindOne(ctx, bson.M{
		"apiKey": apiKey,
		"date":   today,
	}).Decode(&globalLimitRecord)

	if err == mongo.ErrNoDocuments {
		// Primeira request da API Key demo hoje
		globalLimitRecord = GlobalPlaygroundLimit{
			ID:            fmt.Sprintf("%s_%s", apiKey, today),
			APIKey:        apiKey,
			Date:          today,
			TotalRequests: 0,
			UpdatedAt:     now,
		}
	}

	// Verificar limite global
	if globalLimitRecord.TotalRequests >= globalLimit {
		fmt.Printf("🚫 [PLAYGROUND SECURITY] Limite global da API Key demo excedido: %d >= %d\n",
			globalLimitRecord.TotalRequests, globalLimit)

		c.JSON(http.StatusTooManyRequests, gin.H{
			"type":   "https://retech-core/errors/rate-limit-exceeded",
			"title":  "Rate Limit Exceeded (Global)",
			"status": http.StatusTooManyRequests,
			"detail": "Limite global do playground excedido. Tente novamente amanhã.",
		})
		c.Abort()
		return false
	}

	return true
}

// applyThrottling aplica delay mínimo entre requests (anti-spam)
func (prl *PlaygroundRateLimiter) applyThrottling(ctx context.Context, clientIP, apiKey string, c *gin.Context) {
	now := time.Now()

	// Buscar última request do IP
	coll := prl.db.Collection("playground_rate_limits")

	var lastRequest PlaygroundRateLimit
	err := coll.FindOne(ctx, bson.M{
		"ipAddress": clientIP,
		"apiKey":    apiKey,
	}).Decode(&lastRequest)

	if err == nil {
		// Verificar se passou tempo suficiente (2 segundos mínimo)
		timeSinceLastRequest := now.Sub(lastRequest.LastRequest)
		minInterval := 2 * time.Second

		if timeSinceLastRequest < minInterval {
			remainingTime := minInterval - timeSinceLastRequest
			fmt.Printf("⏱️  [PLAYGROUND SECURITY] Throttling aplicado: %v restante\n", remainingTime)

			c.Header("Retry-After", fmt.Sprintf("%.0f", remainingTime.Seconds()))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"type":   "https://retech-core/errors/rate-limit-exceeded",
				"title":  "Rate Limit Exceeded (Throttling)",
				"status": http.StatusTooManyRequests,
				"detail": fmt.Sprintf("Aguarde %.0f segundos antes de fazer outra requisição.", remainingTime.Seconds()),
			})
			c.Abort()
			return
		}
	}
}

// incrementCounters incrementa contadores de rate limit
func (prl *PlaygroundRateLimiter) incrementCounters(ctx context.Context, clientIP, apiKey string, c *gin.Context) {
	now := time.Now()
	// Usar timezone Brasília
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrasilia := now.In(loc)
	today := nowBrasilia.Format("2006-01-02")
	currentMinute := nowBrasilia.Format("2006-01-02 15:04")

	// Incrementar contador por IP
	collIP := prl.db.Collection("playground_rate_limits")

	// ✅ Buscar registro atual para verificar se mudou de minuto
	var currentRecord PlaygroundRateLimit
	err := collIP.FindOne(ctx, bson.M{
		"ipAddress": clientIP,
		"apiKey":    apiKey,
		"date":      today,
	}).Decode(&currentRecord)

	opts := options.Update().SetUpsert(true)

	// ✅ Se mudou de minuto, resetar countPerMinute
	var updateDoc bson.M
	if err == mongo.ErrNoDocuments || currentRecord.Minute != currentMinute {
		// Novo minuto ou primeiro registro do dia
		updateDoc = bson.M{
			"$inc": bson.M{
				"countPerDay": 1,
			},
			"$set": bson.M{
				"countPerMinute": 1, // ✅ Reseta para 1 no novo minuto
				"minute":         currentMinute,
				"lastRequest":    now,
				"updatedAt":      now,
			},
		}
	} else {
		// Mesmo minuto, incrementar normalmente
		updateDoc = bson.M{
			"$inc": bson.M{
				"countPerDay":    1,
				"countPerMinute": 1,
			},
			"$set": bson.M{
				"minute":      currentMinute,
				"lastRequest": now,
				"updatedAt":   now,
			},
		}
	}

	_, err = collIP.UpdateOne(ctx, bson.M{
		"ipAddress": clientIP,
		"apiKey":    apiKey,
		"date":      today,
	}, updateDoc, opts)

	if err != nil {
		fmt.Printf("⚠️  [PLAYGROUND SECURITY] Erro ao incrementar contador por IP: %v\n", err)
	}

	// Incrementar contador global
	collGlobal := prl.db.Collection("playground_global_limits")

	_, err = collGlobal.UpdateOne(ctx, bson.M{
		"apiKey": apiKey,
		"date":   today,
	}, bson.M{
		"$inc": bson.M{
			"totalRequests": 1,
		},
		"$set": bson.M{
			"updatedAt": now,
		},
	}, opts)

	if err != nil {
		fmt.Printf("⚠️  [PLAYGROUND SECURITY] Erro ao incrementar contador global: %v\n", err)
	}
}

// Helper functions

func isPlaygroundRoute(path string) bool {
	playgroundRoutes := []string{
		"/public/cep",
		"/public/cnpj",
		"/public/geo",
		// ❌ NÃO incluir /public/playground/status - é rota pública sem autenticação
	}

	for _, route := range playgroundRoutes {
		if path == route || (len(path) > len(route) && path[:len(route)] == route) {
			return true
		}
	}
	return false
}

func getClientIP(c *gin.Context) string {
	// Tentar diferentes headers para obter IP real
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For pode ter múltiplos IPs (proxy chain)
		// Pegar o primeiro (IP original)
		if idx := len(ip); idx > 0 {
			ip = ip[:idx]
		}
		return ip
	}

	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = c.GetHeader("CF-Connecting-IP") // Cloudflare
	if ip != "" {
		return ip
	}

	// Fallback para IP direto
	return c.ClientIP()
}

func isDemoAPIKey(apiKey string) bool {
	// API Key demo sempre começa com "rtc_demo_"
	return len(apiKey) > 9 && apiKey[:9] == "rtc_demo_"
}

func getNextDayTimestampPlayground() string {
	now := time.Now()
	// Usar timezone Brasília
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrasilia := now.In(loc)
	nextDay := nowBrasilia.AddDate(0, 0, 1)
	nextDayMidnight := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
	return fmt.Sprintf("%d", nextDayMidnight.Unix())
}

func getNextMinuteTimestampPlayground() string {
	now := time.Now()
	// Usar timezone Brasília
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	nowBrasilia := now.In(loc)
	nextMinute := nowBrasilia.Add(time.Minute)
	nextMinuteStart := time.Date(nextMinute.Year(), nextMinute.Month(), nextMinute.Day(), nextMinute.Hour(), nextMinute.Minute(), 0, 0, nextMinute.Location())
	return fmt.Sprintf("%d", nextMinuteStart.Unix())
}
