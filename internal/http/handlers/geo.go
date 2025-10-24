package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/cache"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type GeoHandler struct {
	estados    *storage.EstadosRepo
	municipios *storage.MunicipiosRepo
	redis      interface{} // interface{} para permitir nil (graceful degradation)
}

func NewGeoHandler(estados *storage.EstadosRepo, municipios *storage.MunicipiosRepo, redis interface{}) *GeoHandler {
	return &GeoHandler{
		estados:    estados,
		municipios: municipios,
		redis:      redis,
	}
}

// Response padrão de sucesso
type SuccessResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Response padrão de erro (RFC 7807)
type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	TraceID  string `json:"traceId,omitempty"`
}

// ListUFs retorna todos os estados
// GET /geo/ufs
func (h *GeoHandler) ListUFs(c *gin.Context) {
	ctx := c.Request.Context()
	query := strings.ToLower(c.Query("q"))

	// ⚡ CACHE REDIS (apenas sem filtro)
	if query == "" && h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			redisKey := "geo:ufs:all"
			cachedJSON, err := redisClient.Get(ctx, redisKey)
			if err == nil && cachedJSON != "" {
				c.Header("Content-Type", "application/json")
				c.String(http.StatusOK, cachedJSON)
				return // ⚡ <1ms!
			}
		}
	}

	// 🗄️ BUSCAR DO MONGODB
	estados, err := h.estados.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Type:     "https://retech-core/errors/database-error",
			Title:    "Database Error",
			Status:   http.StatusInternalServerError,
			Detail:   "Erro ao buscar estados",
			Instance: c.Request.URL.Path,
		})
		return
	}

	// Aplica filtro opcional (client-side)
	if query != "" {
		filtered := []interface{}{}
		for _, e := range estados {
			if strings.Contains(strings.ToLower(e.Nome), query) ||
				strings.Contains(strings.ToLower(e.Sigla), query) {
				filtered = append(filtered, e)
			}
		}
		c.JSON(http.StatusOK, SuccessResponse{
			Success: true,
			Code:    "OK",
			Data:    filtered,
		})
		return
	}

	response := SuccessResponse{
		Success: true,
		Code:    "OK",
		Data:    estados,
	}

	// ✅ SALVAR NO REDIS (cache longo, dados fixos)
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			redisKey := "geo:ufs:all"
			if err := redisClient.Set(ctx, redisKey, response, 24*time.Hour); err != nil {
				fmt.Printf("⚠️ Erro ao salvar no Redis: %v\n", err)
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetUF retorna um estado pela sigla
// GET /geo/ufs/:sigla
func (h *GeoHandler) GetUF(c *gin.Context) {
	ctx := c.Request.Context()
	sigla := strings.ToUpper(c.Param("sigla"))

	// ⚡ CACHE REDIS
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			redisKey := fmt.Sprintf("geo:uf:%s", sigla)
			cachedJSON, err := redisClient.Get(ctx, redisKey)
			if err == nil && cachedJSON != "" {
				c.Header("Content-Type", "application/json")
				c.String(http.StatusOK, cachedJSON)
				return // ⚡ <1ms!
			}
		}
	}

	// 🗄️ BUSCAR DO MONGODB
	estado, err := h.estados.FindBySigla(ctx, sigla)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Type:     "https://retech-core/errors/not-found",
				Title:    "Not Found",
				Status:   http.StatusNotFound,
				Detail:   "Estado não encontrado",
				Instance: c.Request.URL.Path,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Type:     "https://retech-core/errors/database-error",
			Title:    "Database Error",
			Status:   http.StatusInternalServerError,
			Detail:   "Erro ao buscar estado",
			Instance: c.Request.URL.Path,
		})
		return
	}

	response := SuccessResponse{
		Success: true,
		Code:    "OK",
		Data:    estado,
	}

	// ✅ SALVAR NO REDIS (cache longo, dados fixos)
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			redisKey := fmt.Sprintf("geo:uf:%s", sigla)
			if err := redisClient.Set(ctx, redisKey, response, 24*time.Hour); err != nil {
				fmt.Printf("⚠️ Erro ao salvar no Redis: %v\n", err)
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

// ListMunicipios retorna todos os municípios (ou filtra por UF)
// GET /geo/municipios
// GET /geo/municipios?uf=PE
func (h *GeoHandler) ListMunicipios(c *gin.Context) {
	ctx := c.Request.Context()
	uf := c.Query("uf")
	query := c.Query("q")

	// Se tem busca por nome
	if query != "" {
		municipios, err := h.municipios.Search(ctx, query, uf)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Type:     "https://retech-core/errors/database-error",
				Title:    "Database Error",
				Status:   http.StatusInternalServerError,
				Detail:   "Erro ao buscar municípios",
				Instance: c.Request.URL.Path,
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{
			Success: true,
			Code:    "OK",
			Data:    municipios,
		})
		return
	}

	// Se filtra por UF
	if uf != "" {
		municipios, err := h.municipios.FindByUF(ctx, uf)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Type:     "https://retech-core/errors/database-error",
				Title:    "Database Error",
				Status:   http.StatusInternalServerError,
				Detail:   "Erro ao buscar municípios",
				Instance: c.Request.URL.Path,
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{
			Success: true,
			Code:    "OK",
			Data:    municipios,
		})
		return
	}

	// Retorna todos (cuidado: pode ser muito grande)
	municipios, err := h.municipios.FindAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Type:     "https://retech-core/errors/database-error",
			Title:    "Database Error",
			Status:   http.StatusInternalServerError,
			Detail:   "Erro ao buscar municípios",
			Instance: c.Request.URL.Path,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Code:    "OK",
		Data:    municipios,
	})
}

// ListMunicipiosByUF retorna municípios de um estado específico
// GET /geo/municipios/:uf
func (h *GeoHandler) ListMunicipiosByUF(c *gin.Context) {
	ctx := c.Request.Context()
	uf := strings.ToUpper(c.Param("uf"))

	municipios, err := h.municipios.FindByUF(ctx, uf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Type:     "https://retech-core/errors/database-error",
			Title:    "Database Error",
			Status:   http.StatusInternalServerError,
			Detail:   "Erro ao buscar municípios",
			Instance: c.Request.URL.Path,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Code:    "OK",
		Data:    municipios,
	})
}

// GetMunicipio retorna um município pelo ID do IBGE
// GET /geo/municipios/id/:id
func (h *GeoHandler) GetMunicipio(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Type:     "https://retech-core/errors/validation-error",
			Title:    "Validation Error",
			Status:   http.StatusBadRequest,
			Detail:   "ID inválido",
			Instance: c.Request.URL.Path,
		})
		return
	}

	municipio, err := h.municipios.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Type:     "https://retech-core/errors/not-found",
				Title:    "Not Found",
				Status:   http.StatusNotFound,
				Detail:   "Município não encontrado",
				Instance: c.Request.URL.Path,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Type:     "https://retech-core/errors/database-error",
			Title:    "Database Error",
			Status:   http.StatusInternalServerError,
			Detail:   "Erro ao buscar município",
			Instance: c.Request.URL.Path,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Code:    "OK",
		Data:    municipio,
	})
}

