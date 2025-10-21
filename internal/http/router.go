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

	// CORS para permitir requisições do frontend
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "https://core.theretech.com.br" || origin == "http://localhost:3000" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
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

	// Rotas públicas (sem autenticação)
	r.GET("/health", health.Get)
	r.GET("/version", handlers.Version)
	r.GET("/docs", handlers.DocsHTML)
	r.GET("/openapi.yaml", handlers.OpenAPIYAML)

	// Auth endpoints (públicos)
	authHandler := handlers.NewAuthHandler(users, tenants, apikeys, activityLogs, jwtService)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.GET("/me", auth.AuthJWT(jwtService), authHandler.Me)
	}

	// GEO endpoints (protegidos por API Key + rate limit + logging)
	geoHandler := handlers.NewGeoHandler(estados, municipios)
	geoGroup := r.Group("/geo")
	geoGroup.Use(
		auth.AuthAPIKey(apikeys), // Requer API Key válida
		rateLimiter.Middleware(), // Aplica rate limiting
		usageLogger.Middleware(), // Loga uso
	)
	{
		geoGroup.GET("/ufs", geoHandler.ListUFs)
		geoGroup.GET("/ufs/:sigla", geoHandler.GetUF)
		geoGroup.GET("/municipios", geoHandler.ListMunicipios)
		geoGroup.GET("/municipios/:uf", geoHandler.ListMunicipiosByUF)
		geoGroup.GET("/municipios/id/:id", geoHandler.GetMunicipio)
	}

	// Admin endpoints (protegidos por JWT + role SUPER_ADMIN)
	adminHandler := handlers.NewAdminHandler(tenants, apikeys, users, m)
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.AuthJWT(jwtService), auth.RequireSuperAdmin())
	{
		// Tenants (admin only)
		tenantsHandler := handlers.NewTenantsHandler(tenants, activityLogs)
		adminGroup.GET("/tenants", tenantsHandler.List)
		adminGroup.GET("/tenants/:id", tenantsHandler.Get)
		adminGroup.POST("/tenants", tenantsHandler.Create)
		adminGroup.PUT("/tenants/:id", tenantsHandler.Update)
		adminGroup.DELETE("/tenants/:id", tenantsHandler.Delete)

		// API Keys (admin only)
		apikeysHandler := handlers.NewAPIKeysHandler(apikeys, tenants, activityLogs)
		adminGroup.GET("/apikeys", adminHandler.ListAllAPIKeys)
		adminGroup.POST("/apikeys", apikeysHandler.Create)
		adminGroup.POST("/apikeys/revoke", apikeysHandler.Revoke)

		// Analytics (admin only)
		adminGroup.GET("/stats", adminHandler.GetStats)
		adminGroup.GET("/usage", adminHandler.GetUsage)

		// Settings (admin only)
		settingsHandler := handlers.NewSettingsHandler(settings, activityLogs)
		adminGroup.GET("/settings", settingsHandler.Get)
		adminGroup.PUT("/settings", settingsHandler.Update)

		// Activity Logs (admin only)
		activityHandler := handlers.NewActivityHandler(activityLogs)
		adminGroup.GET("/activity", activityHandler.GetRecent)
		adminGroup.GET("/activity/user/:userId", activityHandler.GetByUser)
		adminGroup.GET("/activity/type/:type", activityHandler.GetByType)
		adminGroup.GET("/activity/resource/:type/:id", activityHandler.GetByResource)
	}

	// Tenant endpoints (protegidos por JWT + role TENANT_USER)
	tenantHandler := handlers.NewTenantHandler(apikeys, users, m)
	meGroup := r.Group("/me")
	meGroup.Use(auth.AuthJWT(jwtService), auth.RequireTenantUser())
	{
		// Minhas API Keys
		meGroup.GET("/apikeys", tenantHandler.ListMyAPIKeys)
		meGroup.POST("/apikeys", tenantHandler.CreateAPIKey)
		meGroup.DELETE("/apikeys/:id", tenantHandler.DeleteAPIKey)

		// Meu uso
		meGroup.GET("/usage", tenantHandler.GetMyUsage)
	}

	return r
}
