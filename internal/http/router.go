package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/theretech/retech-core/internal/auth"
	"github.com/theretech/retech-core/internal/http/handlers"
	"github.com/theretech/retech-core/internal/storage"
)

func NewRouter(log zerolog.Logger, health *handlers.HealthHandler,
	j *auth.JWTService, users *storage.UsersRepo, tokens *storage.TokensRepo, apikeys *storage.APIKeysRepo,
) *gin.Engine {
	r := gin.New()
	// middlewares...
	// rotas básicas
	r.GET("/health", health.Get)
	r.GET("/version", handlers.Version)
	r.GET("/docs", handlers.DocsHTML)
	r.GET("/openapi.yaml", handlers.OpenAPIYAML)

	// Auth
	ah := handlers.NewAuthHandler(users, tokens, j)
	r.POST("/auth/register", ah.Register)   // desabilita em prod
	r.POST("/auth/login", ah.Login)
	r.POST("/auth/refresh", ah.Refresh)
	r.POST("/auth/logout", ah.Logout)

	// Exemplo protegido (JWT)
	r.GET("/me", auth.AuthJWT(j), handlers.Me)

	// API Keys
	k := handlers.NewAPIKeysHandler(apikeys)
	r.POST("/apikeys", k.Create)
	r.POST("/apikeys/rotate", k.Rotate)
	r.POST("/apikeys/revoke", k.Revoke)

	// Exemplo protegido por API Key
	r.GET("/protected/apikey", auth.AuthAPIKey(apikeys), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// Exemplo protegido por Role
	r.GET("/admin/only", auth.AuthJWT(j), auth.RequireRoles("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"admin": true})
	})

	return r
}
