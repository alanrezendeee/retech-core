package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/cache"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PenalHandler struct {
	db    *storage.Mongo
	redis interface{} // interface{} para permitir nil (graceful degradation)
}

func NewPenalHandler(db *storage.Mongo, redis interface{}) *PenalHandler {
	return &PenalHandler{
		db:    db,
		redis: redis,
	}
}

// ListArtigos retorna todos os artigos penais (para autocomplete)
// GET /penal/artigos
func (h *PenalHandler) ListArtigos(c *gin.Context) {
	ctx := c.Request.Context()
	query := strings.ToLower(strings.TrimSpace(c.Query("q")))
	tipo := c.Query("tipo") // "crime", "contravencao" ou vazio (todos)
	legislacao := c.Query("legislacao") // "CP", "LCP", etc

	// Criar chave de cache
	cacheKey := fmt.Sprintf("penal:artigos:%s:%s:%s", query, tipo, legislacao)
	if query == "" && tipo == "" && legislacao == "" {
		cacheKey = "penal:artigos:all"
	}

	// ⚡ CACHE REDIS (ultra-rápido, <1ms)
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			cachedJSON, err := redisClient.Get(ctx, cacheKey)
			if err == nil && cachedJSON != "" {
				c.Header("Content-Type", "application/json")
				c.String(http.StatusOK, cachedJSON)
				return // ⚡ <1ms!
			}
		}
	}

	// 🗄️ BUSCAR DO MONGODB
	collection := h.db.DB.Collection("penal_artigos")

	filter := bson.M{}
	
	// Filtro por tipo
	if tipo != "" {
		filter["tipo"] = tipo
	}
	
	// Filtro por legislação
	if legislacao != "" {
		filter["legislacao"] = legislacao
	}
	
	// Filtro por busca (texto)
	if query != "" {
		filter["busca"] = bson.M{"$regex": query, "$options": "i"}
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "artigo", Value: 1}, {Key: "paragrafo", Value: 1}}).
		SetLimit(100) // Limitar a 100 resultados para autocomplete

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/database-error",
			"title":  "Database Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao buscar artigos penais",
		})
		return
	}
	defer cursor.Close(ctx)

	var artigos []domain.ArtigoPenal
	if err := cursor.All(ctx, &artigos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/database-error",
			"title":  "Database Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao processar resultados",
		})
		return
	}

	// Converter para formato de resposta (autocomplete)
	results := make([]domain.PenalResponse, 0, len(artigos))
	for _, artigo := range artigos {
		results = append(results, domain.PenalResponse{
			Codigo:          artigo.Codigo,
			CodigoFormatado: artigo.CodigoFormatado,
			Descricao:       artigo.Descricao,
			Tipo:            artigo.Tipo,
			Legislacao:      artigo.Legislacao,
			LegislacaoNome:  artigo.LegislacaoNome,
		})
	}

	response := gin.H{
		"success": true,
		"code":    "OK",
		"data":    results,
		"meta": gin.H{
			"total": len(results),
			"query": query,
		},
	}

	// ⚡ Salvar no Redis (cache permanente para dados fixos - 365 dias)
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			// Cache permanente (365 dias) para dados fixos
			redisClient.Set(ctx, cacheKey, response, 365*24*time.Hour)
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetArtigo retorna um artigo específico por código
// GET /penal/artigos/:codigo
func (h *PenalHandler) GetArtigo(c *gin.Context) {
	ctx := c.Request.Context()
	codigo := c.Param("codigo")

	// Normalizar código (remover espaços, converter para minúsculo para busca)
	codigoNormalizado := strings.ToLower(strings.TrimSpace(codigo))

	// Criar chave de cache
	cacheKey := fmt.Sprintf("penal:artigo:%s", codigoNormalizado)

	// ⚡ CACHE REDIS
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			cachedJSON, err := redisClient.Get(ctx, cacheKey)
			if err == nil && cachedJSON != "" {
				c.Header("Content-Type", "application/json")
				c.String(http.StatusOK, cachedJSON)
				return
			}
		}
	}

	// 🗄️ BUSCAR DO MONGODB
	collection := h.db.DB.Collection("penal_artigos")

	// Buscar por código exato (sem regex para evitar matches parciais como 55 -> 155)
	filter := bson.M{"codigo": codigo}

	var artigo domain.ArtigoPenal
	err := collection.FindOne(ctx, filter).Decode(&artigo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"type":   "https://retech-core/errors/not-found",
				"title":  "Artigo Not Found",
				"status": http.StatusNotFound,
				"detail": fmt.Sprintf("Artigo %s não encontrado", codigo),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/database-error",
			"title":  "Database Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao buscar artigo",
		})
		return
	}

	response := gin.H{
		"success": true,
		"code":    "OK",
		"data":    artigo,
	}

	// ⚡ Salvar no Redis (cache permanente - 365 dias)
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			redisClient.Set(ctx, cacheKey, response, 365*24*time.Hour)
		}
	}

	c.JSON(http.StatusOK, response)
}

// SearchArtigos busca artigos por texto (descrição)
// GET /penal/search?q=texto
func (h *PenalHandler) SearchArtigos(c *gin.Context) {
	ctx := c.Request.Context()
	query := strings.TrimSpace(c.Query("q"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation",
			"title":  "Invalid Query",
			"status": http.StatusBadRequest,
			"detail": "Parâmetro 'q' é obrigatório",
		})
		return
	}

	// Criar chave de cache
	queryLower := strings.ToLower(query)
	cacheKey := fmt.Sprintf("penal:search:%s", queryLower)

	// ⚡ CACHE REDIS
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			cachedJSON, err := redisClient.Get(ctx, cacheKey)
			if err == nil && cachedJSON != "" {
				c.Header("Content-Type", "application/json")
				c.String(http.StatusOK, cachedJSON)
				return
			}
		}
	}

	// 🗄️ BUSCAR DO MONGODB
	collection := h.db.DB.Collection("penal_artigos")

	// Busca em múltiplos campos
	filter := bson.M{
		"$or": []bson.M{
			{"descricao": bson.M{"$regex": query, "$options": "i"}},
			{"textoCompleto": bson.M{"$regex": query, "$options": "i"}},
			{"busca": bson.M{"$regex": queryLower, "$options": "i"}},
		},
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "artigo", Value: 1}}).
		SetLimit(50)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/database-error",
			"title":  "Database Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao buscar artigos",
		})
		return
	}
	defer cursor.Close(ctx)

	var artigos []domain.ArtigoPenal
	if err := cursor.All(ctx, &artigos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://retech-core/errors/database-error",
			"title":  "Database Error",
			"status": http.StatusInternalServerError,
			"detail": "Erro ao processar resultados",
		})
		return
	}

	// Converter para formato de resposta
	results := make([]domain.PenalResponse, 0, len(artigos))
	for _, artigo := range artigos {
		results = append(results, domain.PenalResponse{
			Codigo:          artigo.Codigo,
			CodigoFormatado: artigo.CodigoFormatado,
			Descricao:       artigo.Descricao,
			Tipo:            artigo.Tipo,
			Legislacao:      artigo.Legislacao,
			LegislacaoNome:  artigo.LegislacaoNome,
		})
	}

	response := gin.H{
		"success": true,
		"code":    "OK",
		"data":    results,
		"meta": gin.H{
			"total": len(results),
			"query": query,
		},
	}

	// ⚡ Salvar no Redis (cache 24h para buscas)
	if h.redis != nil {
		if redisClient, ok := h.redis.(*cache.RedisClient); ok {
			redisClient.Set(ctx, cacheKey, response, 24*time.Hour)
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetCacheStats retorna estatísticas do cache de Artigos Penais
// GET /admin/cache/penal/stats
func (h *PenalHandler) GetCacheStats(c *gin.Context) {
	ctx := c.Request.Context()
	collection := h.db.DB.Collection("penal_artigos")

	// Total de artigos penais no banco (seed permanente)
	totalCached, _ := collection.CountDocuments(ctx, bson.M{})

	// Artigos adicionados nas últimas 24h (caso tenha novos artigos)
	yesterday := time.Now().Add(-24 * time.Hour)
	recentCached, _ := collection.CountDocuments(ctx, bson.M{
		"createdAt": bson.M{"$gte": yesterday},
	})

	c.JSON(http.StatusOK, gin.H{
		"totalCached":  totalCached,
		"recentCached": recentCached, // últimas 24h
		"cacheEnabled": true,          // Sempre habilitado (dados fixos)
		"cacheTTLDays": 365,           // Cache permanente (1 ano)
		"autoCleanup":  false,         // Não limpa automaticamente (dados fixos)
	})
}

