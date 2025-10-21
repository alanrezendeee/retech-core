package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

type TenantsHandler struct {
	repo *storage.TenantsRepo
}

func NewTenantsHandler(repo *storage.TenantsRepo) *TenantsHandler {
	return &TenantsHandler{repo: repo}
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

	if err := h.repo.Update(ctx, tenantID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao atualizar tenant",
		})
		return
	}

	// Buscar tenant atualizado
	tenant, _ := h.repo.ByTenantID(ctx, tenantID)

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

	if err := h.repo.Delete(ctx, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/internal-error",
			"title":  "Internal Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao deletar tenant",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
