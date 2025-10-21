package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/http/handlers"
	"github.com/theretech/retech-core/internal/middleware"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(
	log zerolog.Logger,
	db *mongo.Database,
	health *handlers.HealthHandler,
	apikeys *storage.APIKeysRepo,
	tenants *storage.TenantsRepo,
	users *storage.UsersRepo,
	estados *storage.EstadosRepo,
	municipios *storage.MunicipiosRepo,
	jwtService *auth.JWTService,
) *gin.Engine {
	r := gin.New()

	// Solução baseada no Stack Overflow: configurar para não reutilizar conexões
	r.Use(func(c *gin.Context) {
		fmt.Printf("Middleware global chamado para: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Header("Connection", "close")
		c.Next()
	})

	// Middlewares globais
	rateLimiter := middleware.NewRateLimiter(db, domain.GetDefaultRateLimit())
	usageLogger := middleware.NewUsageLogger(db)

	// Rotas públicas (sem autenticação)
	r.GET("/health", health.Get)
	r.GET("/version", handlers.Version)
	r.GET("/docs", handlers.DocsHTML)
	r.GET("/openapi.yaml", handlers.OpenAPIYAML)

	// Auth endpoints (públicos)
	authHandler := handlers.NewAuthHandler(users, tenants, apikeys, jwtService)
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
		auth.AuthAPIKey(apikeys),        // Requer API Key válida
		rateLimiter.Middleware(),        // Aplica rate limiting
		usageLogger.Middleware(),        // Loga uso
	)
	{
		geoGroup.GET("/ufs", geoHandler.ListUFs)
		geoGroup.GET("/ufs/:sigla", geoHandler.GetUF)
		geoGroup.GET("/municipios", geoHandler.ListMunicipios)
		geoGroup.GET("/municipios/:uf", geoHandler.ListMunicipiosByUF)
		geoGroup.GET("/municipios/id/:id", geoHandler.GetMunicipio)
	}

	// Admin endpoints (protegidos por JWT + role SUPER_ADMIN)
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.AuthJWT(jwtService), auth.RequireSuperAdmin())
	{
		// Tenants (admin only)
		tenantsHandler := handlers.NewTenantsHandler(tenants)
		adminGroup.GET("/tenants", tenantsHandler.List)
		adminGroup.GET("/tenants/:id", tenantsHandler.Get)
		adminGroup.POST("/tenants", tenantsHandler.Create)
		adminGroup.PUT("/tenants/:id", tenantsHandler.Update)
		adminGroup.DELETE("/tenants/:id", tenantsHandler.Delete)

		// API Keys (admin only)
		apikeysHandler := handlers.NewAPIKeysHandler(apikeys, tenants)
		adminGroup.GET("/apikeys", func(c *gin.Context) {
			// TODO: listar todas as API keys
			c.JSON(200, gin.H{"message": "list all api keys"})
		})
		adminGroup.POST("/apikeys", apikeysHandler.Create)
		adminGroup.POST("/apikeys/revoke", apikeysHandler.Revoke)

		// Analytics (admin only)
		adminGroup.GET("/stats", func(c *gin.Context) {
			// TODO: estatísticas globais
			c.JSON(200, gin.H{"message": "global stats"})
		})
		adminGroup.GET("/usage", func(c *gin.Context) {
			// TODO: uso da API
			c.JSON(200, gin.H{"message": "api usage"})
		})
	}

	// Tenant endpoints (protegidos por JWT + role TENANT_USER)
	meGroup := r.Group("/me")
	meGroup.Use(auth.AuthJWT(jwtService), auth.RequireTenantUser())
	{
		// Minhas API Keys
		meGroup.GET("/apikeys", func(c *gin.Context) {
			// TODO: listar minhas keys
			c.JSON(200, gin.H{"message": "my api keys"})
		})
		meGroup.POST("/apikeys", func(c *gin.Context) {
			// TODO: criar key
			c.JSON(200, gin.H{"message": "create key"})
		})
		meGroup.DELETE("/apikeys/:id", func(c *gin.Context) {
			// TODO: deletar key
			c.JSON(200, gin.H{"message": "delete key"})
		})

		// Meu uso
		meGroup.GET("/usage", func(c *gin.Context) {
			// TODO: meu uso
			c.JSON(200, gin.H{"message": "my usage"})
		})
	}

	return r
}
