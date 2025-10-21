package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

type SettingsHandler struct {
	settings *storage.SettingsRepo
}

func NewSettingsHandler(settings *storage.SettingsRepo) *SettingsHandler {
	return &SettingsHandler{
		settings: settings,
	}
}

// GetSettings retorna as configurações do sistema
// GET /admin/settings
func (h *SettingsHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	
	settings, err := h.settings.Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao carregar configurações",
			"status": http.StatusInternalServerError,
			"detail": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, settings)
}

// UpdateSettings atualiza as configurações do sistema
// PUT /admin/settings
func (h *SettingsHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	
	var settings domain.SystemSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Erro de validação",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}
	
	// Validações
	if settings.DefaultRateLimit.RequestsPerDay < 1 || settings.DefaultRateLimit.RequestsPerDay > 1000000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Erro de validação",
			"status": http.StatusBadRequest,
			"detail": "RequestsPerDay deve estar entre 1 e 1.000.000",
		})
		return
	}
	
	if settings.DefaultRateLimit.RequestsPerMinute < 1 || settings.DefaultRateLimit.RequestsPerMinute > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Erro de validação",
			"status": http.StatusBadRequest,
			"detail": "RequestsPerMinute deve estar entre 1 e 10.000",
		})
		return
	}
	
	if settings.JWT.AccessTokenTTL < 60 || settings.JWT.AccessTokenTTL > 3600 {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Erro de validação",
			"status": http.StatusBadRequest,
			"detail": "AccessTokenTTL deve estar entre 60 e 3600 segundos",
		})
		return
	}
	
	if settings.JWT.RefreshTokenTTL < 3600 || settings.JWT.RefreshTokenTTL > 2592000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Erro de validação",
			"status": http.StatusBadRequest,
			"detail": "RefreshTokenTTL deve estar entre 3600 e 2.592.000 segundos",
		})
		return
	}
	
	if err := h.settings.Update(ctx, &settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao salvar configurações",
			"status": http.StatusInternalServerError,
			"detail": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Configurações atualizadas com sucesso",
		"settings": settings,
	})
}

