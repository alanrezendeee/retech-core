package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/http/handlers"
	"github.com/theretech/retech-core/internal/storage"
)

func NewRouter(log zerolog.Logger, health *handlers.HealthHandler, apikeys *storage.APIKeysRepo, tenants *storage.TenantsRepo,
) *gin.Engine {
	r := gin.New()

	// Solução baseada no Stack Overflow: configurar para não reutilizar conexões
	r.Use(func(c *gin.Context) {
		fmt.Printf("Middleware global chamado para: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Header("Connection", "close")
		c.Next()
	})

	// rotas básicas
	r.GET("/health", health.Get)
	r.GET("/version", handlers.Version)
	r.GET("/docs", handlers.DocsHTML)
	r.GET("/openapi.yaml", handlers.OpenAPIYAML)

	// Tenants
	t := handlers.NewTenantsHandler(tenants)
	r.POST("/tenants", t.Create)
	r.GET("/tenants", t.List)
	r.GET("/tenants/:id", t.Get)
	r.PUT("/tenants/:id", t.Update)
	r.DELETE("/tenants/:id", t.Delete)

	// API Keys
	k := handlers.NewAPIKeysHandler(apikeys, tenants)
	r.POST("/apikeys", k.Create)
	r.POST("/apikeys/revoke", k.Revoke)

	// Rota para rotação (solução alternativa)
	r.POST("/apikeys/refresh", k.RotateTest)

	// Rotate com nome completamente diferente
	r.POST("/rotate-key", k.Rotate)
	r.POST("/rotate-new", k.RotateNew)
	r.POST("/test-rotate-completely-different", k.Rotate)
	r.POST("/test-rotate-handler", k.RotateTest)

	// Teste GET simples
	r.GET("/test-rotate", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "GET rotate test"})
	})

	// Teste GET com handler
	r.GET("/test-rotate-handler-get", k.RotateTest)

	// Teste POST simples
	r.POST("/test-rotate-post", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POST rotate test"})
	})

	// Exemplo protegido por API Key
	r.GET("/protected/apikey", auth.AuthAPIKey(apikeys), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	return r
}
