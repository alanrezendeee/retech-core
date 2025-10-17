package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

type TenantsHandler struct {
	Repo *storage.TenantsRepo
}

func NewTenantsHandler(r *storage.TenantsRepo) *TenantsHandler {
	return &TenantsHandler{Repo: r}
}

type createTenantDTO struct {
	TenantID string `json:"tenantId" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Create cria um novo tenant
func (h *TenantsHandler) Create(c *gin.Context) {
	var in createTenantDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifica se já existe
	existing, err := h.Repo.ByTenantID(c, in.TenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check tenant"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "tenant already exists"})
		return
	}

	tenant := &domain.Tenant{
		TenantID: in.TenantID,
		Name:     in.Name,
		Email:    in.Email,
		Active:   true,
	}

	if err := h.Repo.Insert(c, tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert tenant"})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

// Get retorna um tenant por ID
func (h *TenantsHandler) Get(c *gin.Context) {
	tenantID := c.Param("id")

	tenant, err := h.Repo.ByTenantID(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get tenant"})
		return
	}
	if tenant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// List retorna todos os tenants
func (h *TenantsHandler) List(c *gin.Context) {
	tenants, err := h.Repo.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list tenants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tenants": tenants})
}

type updateTenantDTO struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Active *bool   `json:"active"`
}

// Update atualiza um tenant
func (h *TenantsHandler) Update(c *gin.Context) {
	tenantID := c.Param("id")

	var in updateTenantDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifica se existe
	tenant, err := h.Repo.ByTenantID(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get tenant"})
		return
	}
	if tenant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	// Monta o map de updates
	updates := make(map[string]interface{})
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Email != nil {
		updates["email"] = *in.Email
	}
	if in.Active != nil {
		updates["active"] = *in.Active
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := h.Repo.Update(c, tenantID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update tenant"})
		return
	}

	// Retorna o tenant atualizado
	updated, _ := h.Repo.ByTenantID(c, tenantID)
	c.JSON(http.StatusOK, updated)
}

// Delete remove um tenant
func (h *TenantsHandler) Delete(c *gin.Context) {
	tenantID := c.Param("id")

	// Verifica se existe
	tenant, err := h.Repo.ByTenantID(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get tenant"})
		return
	}
	if tenant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	if err := h.Repo.Delete(c, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete tenant"})
		return
	}

	c.Status(http.StatusNoContent)
}

