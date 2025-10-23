package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"github.com/theretech/retech-core/internal/utils"
)

type TenantsHandler struct {
	repo         *storage.TenantsRepo
	activityRepo *storage.ActivityLogsRepo
	settingsRepo *storage.SettingsRepo
}

func NewTenantsHandler(repo *storage.TenantsRepo, activityRepo *storage.ActivityLogsRepo, settingsRepo *storage.SettingsRepo) *TenantsHandler {
	return &TenantsHandler{
		repo:         repo,
		activityRepo: activityRepo,
		settingsRepo: settingsRepo,
	}
}

// List retorna todos os tenants
func (h *TenantsHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	tenants, err := h.repo.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao buscar tenants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenants": tenants,
		"total":   len(tenants),
	})
}

// Get busca um tenant por ID (tenantId)
func (h *TenantsHandler) Get(c *gin.Context) {
	tenantID := c.Param("id")
	ctx := c.Request.Context()

	tenant, err := h.repo.ByTenantID(ctx, tenantID)
	if err != nil || tenant == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://retech-core/errors/not-found",
			"title":  "Not Found",
			"status": http.StatusNotFound,
			"detail": "Tenant não encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant": tenant,
	})
}

// Create cria um novo tenant
func (h *TenantsHandler) Create(c *gin.Context) {
	var tenant domain.Tenant
	if err := c.ShouldBindJSON(&tenant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Gerar tenantId se não foi fornecido
	if tenant.TenantID == "" {
		tenant.TenantID = fmt.Sprintf("tenant-%d", time.Now().Unix())
	}

	if err := h.repo.Insert(ctx, &tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao criar tenant",
		})
		return
	}

	// Log da atividade
	utils.LogActivity(
		c,
		h.activityRepo,
		domain.ActivityTypeTenantCreated,
		domain.ActionCreate,
		utils.BuildActorFromContext(c),
		domain.Resource{
			Type: domain.ResourceTypeTenant,
			ID:   tenant.TenantID,
			Name: tenant.Name,
		},
		map[string]interface{}{
			"email":   tenant.Email,
			"company": tenant.Company,
			"active":  tenant.Active,
		},
	)

	c.JSON(http.StatusCreated, gin.H{
		"tenant": tenant,
	})
}

// Update atualiza um tenant
func (h *TenantsHandler) Update(c *gin.Context) {
	tenantID := c.Param("id")
	ctx := c.Request.Context()

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation-error",
			"title":  "Validation Error",
			"status": http.StatusBadRequest,
			"detail": err.Error(),
		})
		return
	}

	// 🔍 DEBUG: Ver o que está vindo do frontend
	fmt.Printf("🔍 [Update Tenant] TenantID: %s\n", tenantID)
	fmt.Printf("🔍 [Update Tenant] Updates recebidos: %+v\n", updates)

	// ✅ Lógica de rate limit:
	// - Se rateLimit é null → buscar padrão e salvar
	// - Se rateLimit existe → usar o que veio do frontend
	if rl, ok := updates["rateLimit"]; ok {
		if rl == nil {
			// Frontend enviou null → usar padrão
			fmt.Printf("⚠️ [Update Tenant] rateLimit é NULL, buscando padrão...\n")
			settings, err := h.settingsRepo.Get(ctx)
			if err == nil && settings != nil {
				updates["rateLimit"] = map[string]interface{}{
					"RequestsPerDay":    settings.DefaultRateLimit.RequestsPerDay,
					"RequestsPerMinute": settings.DefaultRateLimit.RequestsPerMinute,
				}
				fmt.Printf("✅ [Update Tenant] Aplicando rate limit padrão: %+v\n", updates["rateLimit"])
			} else {
				// Fallback
				updates["rateLimit"] = map[string]interface{}{
					"RequestsPerDay":    int64(1000),
					"RequestsPerMinute": int64(60),
				}
				fmt.Printf("⚠️ [Update Tenant] Usando fallback: 1000/dia, 60/min\n")
			}
		} else {
			fmt.Printf("✅ [Update Tenant] rateLimit personalizado: %+v (tipo: %T)\n", rl, rl)
		}
	} else {
		fmt.Printf("⚠️ [Update Tenant] rateLimit não enviado, mantendo valor atual\n")
	}

	if err := h.repo.Update(ctx, tenantID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao atualizar tenant",
		})
		return
	}

	fmt.Printf("✅ [Update Tenant] Update executado com sucesso!\n")

	// Buscar tenant atualizado
	tenant, _ := h.repo.ByTenantID(ctx, tenantID)

	fmt.Printf("🔍 [Update Tenant] Tenant após update: %+v\n", tenant)
	if tenant != nil && tenant.RateLimit != nil {
		fmt.Printf("✅ [Update Tenant] RateLimit salvo: %+v\n", tenant.RateLimit)
	} else if tenant != nil {
		fmt.Printf("⚠️ [Update Tenant] RateLimit é NULL após salvar!\n")
	}

	// Log da atividade
	if tenant != nil {
		// Determinar tipo específico de ação
		activityType := domain.ActivityTypeTenantUpdated
		action := domain.ActionUpdate

		if active, ok := updates["active"].(bool); ok {
			if active {
				activityType = domain.ActivityTypeTenantActivated
				action = domain.ActionActivate
			} else {
				activityType = domain.ActivityTypeTenantDeactivated
				action = domain.ActionDeactivate
			}
		}

		utils.LogActivity(
			c,
			h.activityRepo,
			activityType,
			action,
			utils.BuildActorFromContext(c),
			domain.Resource{
				Type: domain.ResourceTypeTenant,
				ID:   tenant.TenantID,
				Name: tenant.Name,
			},
			updates,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant": tenant,
	})
}

// Delete deleta um tenant
func (h *TenantsHandler) Delete(c *gin.Context) {
	tenantID := c.Param("id")
	ctx := c.Request.Context()

	// Proteger tenant do SUPER_ADMIN
	if tenantID == "tenant-20251021145821" {
		c.JSON(http.StatusForbidden, gin.H{
			"type":   "https://retech-core/errors/forbidden",
			"title":  "Forbidden",
			"status": http.StatusForbidden,
			"detail": "Não é possível deletar o tenant do SUPER_ADMIN",
		})
		return
	}

	// Buscar tenant antes de deletar (para log)
	tenant, _ := h.repo.ByTenantID(ctx, tenantID)

	if err := h.repo.Delete(ctx, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao deletar tenant",
		})
		return
	}

	// Log da atividade
	if tenant != nil {
		utils.LogActivity(
			c,
			h.activityRepo,
			domain.ActivityTypeTenantDeleted,
			domain.ActionDelete,
			utils.BuildActorFromContext(c),
			domain.Resource{
				Type: domain.ResourceTypeTenant,
				ID:   tenant.TenantID,
				Name: tenant.Name,
			},
			map[string]interface{}{
				"email":   tenant.Email,
				"company": tenant.Company,
			},
		)
	}

	c.Status(http.StatusNoContent)
}
