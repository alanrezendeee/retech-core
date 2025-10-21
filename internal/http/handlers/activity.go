package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
)

// ActivityHandler gerencia endpoints de activity logs
type ActivityHandler struct {
	repo *storage.ActivityLogsRepo
}

// NewActivityHandler cria um novo handler de activity logs
func NewActivityHandler(repo *storage.ActivityLogsRepo) *ActivityHandler {
	return &ActivityHandler{
		repo: repo,
	}
}

// GetRecent retorna as atividades mais recentes
// GET /admin/activity?limit=20
func (h *ActivityHandler) GetRecent(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse limit (padrão 20, máximo 100)
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	// Buscar atividades
	activities, err := h.repo.Recent(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "internal_server_error",
			"detail": "Erro ao buscar atividades",
		})
		return
	}

	// Se vazio, retornar array vazio ao invés de null
	if activities == nil {
		activities = []*domain.ActivityLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      len(activities),
	})
}

// GetByUser retorna atividades de um usuário específico
// GET /admin/activity/user/:userId?limit=20
func (h *ActivityHandler) GetByUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("userId")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_request",
			"detail": "userId é obrigatório",
		})
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	activities, err := h.repo.ByUser(ctx, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "internal_server_error",
			"detail": "Erro ao buscar atividades do usuário",
		})
		return
	}

	if activities == nil {
		activities = []*domain.ActivityLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      len(activities),
	})
}

// GetByType retorna atividades de um tipo específico
// GET /admin/activity/type/:type?limit=20
func (h *ActivityHandler) GetByType(c *gin.Context) {
	ctx := c.Request.Context()
	eventType := c.Param("type")

	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_request",
			"detail": "type é obrigatório",
		})
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	activities, err := h.repo.ByType(ctx, eventType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "internal_server_error",
			"detail": "Erro ao buscar atividades por tipo",
		})
		return
	}

	if activities == nil {
		activities = []*domain.ActivityLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      len(activities),
	})
}

// GetByResource retorna atividades relacionadas a um recurso
// GET /admin/activity/resource/:type/:id?limit=20
func (h *ActivityHandler) GetByResource(c *gin.Context) {
	ctx := c.Request.Context()
	resourceType := c.Param("type")
	resourceID := c.Param("id")

	if resourceType == "" || resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_request",
			"detail": "type e id são obrigatórios",
		})
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	activities, err := h.repo.ByResource(ctx, resourceType, resourceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "internal_server_error",
			"detail": "Erro ao buscar atividades do recurso",
		})
		return
	}

	if activities == nil {
		activities = []*domain.ActivityLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      len(activities),
	})
}

