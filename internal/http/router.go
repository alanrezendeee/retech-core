package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/http/handlers"
	"github.com/theretech/retech-core/internal/middleware"
	"github.com/theretech/retech-core/internal/storage"
)

func NewRouter(
	log zerolog.Logger,
	m *storage.Mongo,
	redisClient interface{}, // interface{} para permitir nil (graceful degradation)
	health *handlers.HealthHandler,
	apikeys *storage.APIKeysRepo,
	tenants *storage.TenantsRepo,
	users *storage.UsersRepo,
	estados *storage.EstadosRepo,
	municipios *storage.MunicipiosRepo,
	settings *storage.SettingsRepo,
	activityLogs *storage.ActivityLogsRepo,
	jwtService *auth.JWTService,
) *gin.Engine {
	r := gin.New()

	// 🌐 CORS DINÂMICO (segue EXATAMENTE admin/settings)
	r.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		origin := c.Request.Header.Get("Origin")
		method := c.Request.Method
		path := c.Request.URL.Path

		// 🔍 DEBUG: Log de todas as requests
		fmt.Printf("[CORS] %s %s (Origin: %s)\n", method, path, origin)

		// 🔒 Buscar settings (sem fallbacks ou exceções)
		sysSettings, err := settings.Get(ctx)
		if err != nil {
			fmt.Printf("[CORS] ❌ Erro ao buscar settings: %v - SEM headers CORS\n", err)
			
			// ✅ BEST PRACTICE: Não bloquear, apenas não adicionar headers CORS
			// Browser bloqueará por falta dos headers
			if method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
			return
		}

		fmt.Printf("[CORS] Settings: CORS.Enabled=%v, AllowedOrigins=%v\n", 
			sysSettings.CORS.Enabled, sysSettings.CORS.AllowedOrigins)

		// ❌ Se CORS desabilitado, não adicionar headers (browser bloqueará)
		if !sysSettings.CORS.Enabled {
			fmt.Printf("[CORS] ❌ CORS desabilitado - não adicionando headers\n")
			
			// ✅ BEST PRACTICE: Responder OPTIONS com 204, mas SEM headers CORS
			if method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			// Para requests normais, processar normalmente mas sem headers CORS
			c.Next()
			return
		}

		// ✅ CORS habilitado: verificar se origin está na lista
		allowed := false
		for _, allowedOrigin := range sysSettings.CORS.AllowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if !allowed && origin != "" {
			fmt.Printf("[CORS] ❌ Origin '%s' não está na lista permitida: %v\n", origin, sysSettings.CORS.AllowedOrigins)
			
			// ✅ BEST PRACTICE: Não bloquear, apenas não adicionar headers CORS
			if method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			// Para requests normais, processar normalmente mas sem headers CORS
			c.Next()
			return
		}

		// ✅ Origin permitido: adicionar headers CORS
		fmt.Printf("[CORS] ✅ Origin permitido - adicionando headers\n")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-API-Key")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Responder preflight requests
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Solução baseada no Stack Overflow: configurar para não reutilizar conexões
	r.Use(func(c *gin.Context) {
		fmt.Printf("Middleware global chamado para: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Header("Connection", "close")
		c.Next()
	})

	// Middlewares globais
	rateLimiter := middleware.NewRateLimiter(m.DB, tenants, settings)
	usageLogger := middleware.NewUsageLogger(m.DB)
	maintenanceMiddleware := middleware.NewMaintenanceMiddleware(settings)

	// Rotas públicas (sem autenticação e sem manutenção)
	r.GET("/health", health.Get)
	r.GET("/version", handlers.Version)
	r.GET("/docs", handlers.DocsHTML)
	r.GET("/openapi.yaml", handlers.OpenAPIYAML)

	// Public endpoints
	publicSettingsHandler := handlers.NewSettingsHandler(settings, activityLogs)
	r.GET("/public/contact", publicSettingsHandler.GetPublicContact)

	// Playground status (público, sem autenticação)
	playgroundHandler := handlers.NewPlaygroundHandler(settings)
	r.GET("/public/playground/status", playgroundHandler.GetStatus)

	// Public playground/tools endpoints (sem API Key, rate limit por IP)
	cepHandler := handlers.NewCEPHandler(m, redisClient, settings)
	cnpjHandler := handlers.NewCNPJHandler(m, redisClient, settings)
	geoHandler := handlers.NewGeoHandler(estados, municipios, redisClient)

	// 🔒 ROTAS PÚBLICAS DESABILITADAS (usar API Key Demo no playground)
	// Motivo: Prevenir abuso. Playground usa API Key "rtc_demo_playground" com rate limit agressivo.
	// publicGroup := r.Group("/public")
	// {
	// 	publicGroup.GET("/cep/:codigo", cepHandler.GetCEP)
	// 	publicGroup.GET("/cnpj/:numero", cnpjHandler.GetCNPJ)
	// 	publicGroup.GET("/geo/ufs", geoHandler.ListUFs)
	// 	publicGroup.GET("/geo/ufs/:sigla", geoHandler.GetUF)
	// }

	// Auth endpoints (públicos)
	authHandler := handlers.NewAuthHandler(users, tenants, apikeys, activityLogs, settings, jwtService)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.GET("/me", auth.AuthJWT(jwtService), authHandler.Me)
	}

	// GEO endpoints (protegidos por API Key + rate limit + logging + manutenção + scopes)
	geoGroup := r.Group("/geo")
	geoGroup.Use(
		maintenanceMiddleware.Middleware(), // Verifica manutenção
		auth.AuthAPIKey(apikeys),           // Requer API Key válida
		auth.RequireScope(apikeys, "geo"),  // ✅ Verifica scope 'geo' ou 'all'
		rateLimiter.Middleware(),           // Aplica rate limiting
		usageLogger.Middleware(),           // Loga uso
	)
	{
		geoGroup.GET("/ufs", geoHandler.ListUFs)
		geoGroup.GET("/ufs/:sigla", geoHandler.GetUF)
		geoGroup.GET("/municipios", geoHandler.ListMunicipios)
		geoGroup.GET("/municipios/:uf", geoHandler.ListMunicipiosByUF)
		geoGroup.GET("/municipios/id/:id", geoHandler.GetMunicipio)
	}

	// CEP endpoints (protegidos por API Key + rate limit + logging + manutenção + scopes)
	cepGroup := r.Group("/cep")
	cepGroup.Use(
		maintenanceMiddleware.Middleware(), // Verifica manutenção
		auth.AuthAPIKey(apikeys),           // Requer API Key válida
		auth.RequireScope(apikeys, "cep"),  // ✅ Verifica scope 'cep' ou 'all'
		rateLimiter.Middleware(),           // Aplica rate limiting
		usageLogger.Middleware(),           // Loga uso
	)
	{
		cepGroup.GET("/:codigo", cepHandler.GetCEP)
	}

	// CNPJ endpoints (protegidos por API Key + rate limit + logging + manutenção + scopes)
	cnpjGroup := r.Group("/cnpj")
	cnpjGroup.Use(
		maintenanceMiddleware.Middleware(), // Verifica manutenção
		auth.AuthAPIKey(apikeys),           // Requer API Key válida
		auth.RequireScope(apikeys, "cnpj"), // ✅ Verifica scope 'cnpj' ou 'all'
		rateLimiter.Middleware(),           // Aplica rate limiting
		usageLogger.Middleware(),           // Loga uso
	)
	{
		cnpjGroup.GET("/:numero", cnpjHandler.GetCNPJ)
	}

	// Admin endpoints (protegidos por JWT + role SUPER_ADMIN)
	adminHandler := handlers.NewAdminHandler(tenants, apikeys, users, m)
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.AuthJWT(jwtService), auth.RequireSuperAdmin())
	{
		// Tenants (admin only)
		tenantsHandler := handlers.NewTenantsHandler(tenants, activityLogs, settings)
		adminGroup.GET("/tenants", tenantsHandler.List)
		adminGroup.GET("/tenants/:id", tenantsHandler.Get)
		adminGroup.POST("/tenants", tenantsHandler.Create)
		adminGroup.PUT("/tenants/:id", tenantsHandler.Update)
		adminGroup.DELETE("/tenants/:id", tenantsHandler.Delete)

		// API Keys (admin only)
		apikeysHandler := handlers.NewAPIKeysHandler(apikeys, tenants, activityLogs)
		adminGroup.GET("/apikeys", adminHandler.ListAllAPIKeys)
		adminGroup.POST("/apikeys", apikeysHandler.Create)
		adminGroup.POST("/apikeys/rotate", apikeysHandler.Rotate)
		adminGroup.POST("/apikeys/revoke", apikeysHandler.Revoke)

		// Analytics (admin only)
		adminGroup.GET("/stats", adminHandler.GetStats)
		adminGroup.GET("/usage", adminHandler.GetUsage)
		adminGroup.GET("/analytics", adminHandler.GetAnalytics) // ✅ NOVO: Analytics detalhado com breakdown por API

		// Settings (admin only)
		settingsHandler := handlers.NewSettingsHandler(settings, activityLogs)
		adminGroup.GET("/settings", settingsHandler.Get)
		adminGroup.PUT("/settings", settingsHandler.Update)

		// Cache Management (admin only)
		adminGroup.GET("/cache/cep/stats", cepHandler.GetCacheStats)
		adminGroup.DELETE("/cache/cep", cepHandler.ClearCache)
		adminGroup.GET("/cache/cnpj/stats", cnpjHandler.GetCacheStats)
		adminGroup.DELETE("/cache/cnpj", cnpjHandler.ClearCache)

		// Activity Logs (admin only)
		activityHandler := handlers.NewActivityHandler(activityLogs)
		adminGroup.GET("/activity", activityHandler.GetRecent)
		adminGroup.GET("/activity/user/:userId", activityHandler.GetByUser)
		adminGroup.GET("/activity/type/:type", activityHandler.GetByType)
		adminGroup.GET("/activity/resource/:type/:id", activityHandler.GetByResource)
	}

	// Tenant endpoints (protegidos por JWT + role TENANT_USER)
	tenantHandler := handlers.NewTenantHandler(apikeys, users, tenants, m)
	meGroup := r.Group("/me")
	meGroup.Use(auth.AuthJWT(jwtService), auth.RequireTenantUser())
	{
		// Minhas API Keys
		meGroup.GET("/apikeys", tenantHandler.ListMyAPIKeys)
		meGroup.POST("/apikeys", tenantHandler.CreateAPIKey)
		meGroup.POST("/apikeys/:id/rotate", tenantHandler.RotateAPIKey)
		meGroup.DELETE("/apikeys/:id", tenantHandler.DeleteAPIKey)

		// Meu uso
		meGroup.GET("/stats", tenantHandler.GetMyStats)   // Métricas rápidas para dashboard
		meGroup.GET("/usage", tenantHandler.GetMyUsage)   // Uso detalhado com gráficos
		meGroup.GET("/config", tenantHandler.GetMyConfig) // Configurações para docs
	}

	return r
}
