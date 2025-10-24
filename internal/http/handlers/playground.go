package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/storage"
)

type PlaygroundHandler struct {
	settings *storage.SettingsRepo
}

func NewPlaygroundHandler(settings *storage.SettingsRepo) *PlaygroundHandler {
	return &PlaygroundHandler{
		settings: settings,
	}
}

// GetStatus retorna se o playground está habilitado
// GET /public/playground/status
func (h *PlaygroundHandler) GetStatus(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Buscar configurações
	sysSettings, err := h.settings.Get(ctx)
	if err != nil {
		// Se falhar ao buscar settings, assume que está habilitado (graceful degradation)
		c.JSON(http.StatusOK, gin.H{
			"enabled": true,
			"message": "Playground disponível",
		})
		return
	}

	// Verificar se playground está habilitado
	if !sysSettings.Playground.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"enabled": false,
			"message": "Playground temporariamente indisponível. Entre em contato para mais informações.",
		})
		return
	}

	// Playground habilitado
	c.JSON(http.StatusOK, gin.H{
		"enabled": true,
		"message": "Playground disponível",
		"allowedApis": sysSettings.Playground.AllowedAPIs,
	})
}

