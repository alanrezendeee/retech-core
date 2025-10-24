package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"github.com/theretech/retech-core/internal/utils"
)

type SettingsHandler struct {
	settings     *storage.SettingsRepo
	activityRepo *storage.ActivityLogsRepo
}

func NewSettingsHandler(settings *storage.SettingsRepo, activityRepo *storage.ActivityLogsRepo) *SettingsHandler {
	return &SettingsHandler{
		settings:     settings,
		activityRepo: activityRepo,
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
		// Log detalhado do erro
		fmt.Printf("Erro ao fazer bind do JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Erro de validação",
			"status": http.StatusBadRequest,
			"detail": fmt.Sprintf("Erro ao processar JSON: %v", err),
		})
		return
	}

	// Log dos valores recebidos
	fmt.Printf("Settings recebidas: %+v\n", settings)
	fmt.Printf("📦 Cache config recebido: enabled=%v, cepTTLDays=%d, autoCleanup=%v\n", 
		settings.Cache.Enabled, settings.Cache.CEPTTLDays, settings.Cache.AutoCleanup)
	fmt.Printf("🎮 Playground config recebido: enabled=%v, apiKey=%s, reqPerDay=%d, reqPerMin=%d, allowedAPIs=%v\n",
		settings.Playground.Enabled, settings.Playground.APIKey, 
		settings.Playground.RateLimit.RequestsPerDay, settings.Playground.RateLimit.RequestsPerMinute,
		settings.Playground.AllowedAPIs)

	// ⚠️ IMPORTANTE: O campo environment SEMPRE vem da variável ENV
	// Não permitir que seja sobrescrito pelo frontend!
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Getenv("NODE_ENV")
	}
	if env == "" {
		env = "development"
	}
	settings.API.Environment = env
	fmt.Printf("✅ Environment forçado para: %s (da variável ENV)\n", env)

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
		fmt.Printf("Erro ao atualizar settings no MongoDB: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Erro ao salvar configurações",
			"status": http.StatusInternalServerError,
			"detail": fmt.Sprintf("Erro ao salvar no banco: %v", err),
		})
		return
	}

	// Log da atividade
	utils.LogActivity(
		c,
		h.activityRepo,
		domain.ActivityTypeSettingsUpdated,
		domain.ActionUpdate,
		utils.BuildActorFromContext(c),
		domain.Resource{
			Type: domain.ResourceTypeSettings,
			ID:   "system-settings",
			Name: "Configurações do Sistema",
		},
		map[string]interface{}{
			"defaultRateLimit": settings.DefaultRateLimit,
			"apiVersion":       settings.API.Version,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Configurações atualizadas com sucesso",
		"settings": settings,
	})
}

// GetPublicContact retorna informações públicas de contato (sem auth)
// GET /public/contact
func (h *SettingsHandler) GetPublicContact(c *gin.Context) {
	ctx := c.Request.Context()

	settings, err := h.settings.Get(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao carregar configurações de contato",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"whatsapp": settings.Contact.WhatsApp,
		"email":    settings.Contact.Email,
		"phone":    settings.Contact.Phone,
	})
}
