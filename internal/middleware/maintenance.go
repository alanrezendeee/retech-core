package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/storage"
)

type MaintenanceMiddleware struct {
	settings *storage.SettingsRepo
}

func NewMaintenanceMiddleware(settings *storage.SettingsRepo) *MaintenanceMiddleware {
	return &MaintenanceMiddleware{
		settings: settings,
	}
}

// Middleware verifica se a API está em modo de manutenção
func (m *MaintenanceMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Buscar configurações
		settings, err := m.settings.Get(ctx)
		if err != nil {
			// Se erro ao buscar settings, permite a request continuar
			c.Next()
			return
		}

		// Verificar se está em manutenção
		if settings != nil && settings.API.Maintenance {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"type":   "https://retech-core/errors/maintenance",
				"title":  "API em Manutenção",
				"status": http.StatusServiceUnavailable,
				"detail": "A API está temporariamente indisponível para manutenção. Por favor, tente novamente mais tarde.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
