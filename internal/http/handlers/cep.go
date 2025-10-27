package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/cache"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"github.com/theretech/retech-core/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CEPHandler struct {
	db       *storage.Mongo
	redis    interface{} // interface{} para permitir nil (graceful degradation)
	settings *storage.SettingsRepo
}

func NewCEPHandler(db *storage.Mongo, redis interface{}, settings *storage.SettingsRepo) *CEPHandler {
	return &CEPHandler{
		db:       db,
		redis:    redis,
		settings: settings,
	}
}

// getTTL retorna o TTL configurado no admin/settings (ou 7 dias como padrão)
func (h *CEPHandler) getTTL(ctx *gin.Context) time.Duration {
	sysSettings, err := h.settings.Get(ctx.Request.Context())
	if err != nil || !sysSettings.Cache.Enabled {
		return 7 * 24 * time.Hour // Padrão: 7 dias
	}
	
	// Validar intervalo (1-365 dias)
	days := sysSettings.Cache.CEPTTLDays
	if days < 1 {
		days = 1
	}
	if days > 365 {
		days = 365
	}
	
	return time.Duration(days) * 24 * time.Hour
}

// CEPResponse representa o retorno da API de CEP
type CEPResponse struct {
	CEP         string  `json:"cep" bson:"cep"`
	Logradouro  string  `json:"logradouro" bson:"logradouro"`
	Complemento string  `json:"complemento,omitempty" bson:"complemento,omitempty"`
	Bairro      string  `json:"bairro" bson:"bairro"`
	Localidade  string  `json:"localidade" bson:"localidade"`
	UF          string  `json:"uf" bson:"uf"`
	IBGE        string  `json:"ibge,omitempty" bson:"ibge,omitempty"`
	DDD         string  `json:"ddd,omitempty" bson:"ddd,omitempty"`
	Latitude    float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Source      string  `json:"source" bson:"source"` // viacep, brasilapi, cache
	CachedAt    string  `json:"cachedAt,omitempty" bson:"cachedAt,omitempty"`
}

// GET /cep/:codigo
// Consulta CEP com cache, ViaCEP como principal e Brasil API como fallback
func (h *CEPHandler) GetCEP(c *gin.Context) {
	// ⏱️ Iniciar medição de tempo do servidor (SEM rede)
	startTime := time.Now()

	cep := c.Param("codigo")

	// Limpar CEP (remover pontos e traços)
	cep = strings.ReplaceAll(cep, "-", "")
	cep = strings.ReplaceAll(cep, ".", "")

	// Validar formato
	if len(cep) != 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation",
			"title":  "Invalid CEP",
			"status": http.StatusBadRequest,
			"detail": "CEP deve ter 8 dígitos",
		})
		return
	}

	ctx := c.Request.Context()

	// ⏱️ Adicionar header com tempo de processamento do servidor (ao final da função)
	defer func() {
		serverTime := time.Since(startTime)
		c.Header("X-Server-Time-Ms", fmt.Sprintf("%.2f", float64(serverTime.Microseconds())/1000.0))
		fmt.Printf("⏱️ [CEP:%s] Tempo de processamento do servidor: %.2fms\n", cep, float64(serverTime.Microseconds())/1000.0)
	}()

	// Carregar configurações de cache
	settings, err := h.settings.Get(ctx)
	if err != nil {
		settings = domain.GetDefaultSettings() // Fallback para padrões
	}

	// ⚡ CAMADA 1: REDIS (ultra-rápido, <1ms)
	if h.redis != nil && settings.Cache.Enabled {
		redisClient, ok := h.redis.(*cache.RedisClient)
		if ok {
			redisKey := fmt.Sprintf("cep:%s", cep)
			cachedJSON, err := redisClient.Get(ctx, redisKey)
			if err == nil && cachedJSON != "" {
				var cached CEPResponse
				if json.Unmarshal([]byte(cachedJSON), &cached) == nil {
					cached.Source = "redis-cache"
					fmt.Printf("✅ [CEP:%s] CACHE HIT → Redis L1 (ultra-rápido)\n", cep)
					c.JSON(http.StatusOK, cached)
					return // ⚡ <1ms!
				}
			}
			fmt.Printf("⚠️ [CEP:%s] CACHE MISS → Redis L1 (tentando L2...)\n", cep)
		}
	} else {
		if h.redis == nil {
			fmt.Printf("⚠️ [CEP:%s] Redis não disponível (graceful degradation)\n", cep)
		}
		if !settings.Cache.Enabled {
			fmt.Printf("⚠️ [CEP:%s] Cache desabilitado nas configurações\n", cep)
		}
	}

	// 🗄️ CAMADA 2: MONGODB (backup, ~10ms)
	collection := h.db.DB.Collection("cep_cache")

	if settings.Cache.Enabled {
		var cached CEPResponse
		err := collection.FindOne(ctx, bson.M{"cep": cep}).Decode(&cached)
		if err == nil {
			// Verificar se cache ainda é válido (usar TTL dinâmico)
			cacheTTL := time.Duration(settings.Cache.CEPTTLDays) * 24 * time.Hour
			cachedTime, _ := time.Parse(time.RFC3339, cached.CachedAt)
			if time.Since(cachedTime) < cacheTTL {
				fmt.Printf("✅ [CEP:%s] CACHE HIT → MongoDB L2 (válido, promovendo para Redis...)\n", cep)

				// ✅ Promover para Redis (para próximas requests)
				if h.redis != nil {
					if redisClient, ok := h.redis.(*cache.RedisClient); ok {
						redisKey := fmt.Sprintf("cep:%s", cep)
						ttl := h.getTTL(c)
						if err := redisClient.Set(ctx, redisKey, cached, ttl); err == nil {
							fmt.Printf("✅ [CEP:%s] Promovido para Redis L1 com sucesso (TTL: %v)\n", cep, ttl)
						} else {
							fmt.Printf("⚠️ [CEP:%s] Erro ao promover para Redis: %v\n", cep, err)
						}
					}
				}
				cached.Source = "mongodb-cache"
				c.JSON(http.StatusOK, cached)
				return // ~10ms
			} else {
				fmt.Printf("⚠️ [CEP:%s] CACHE EXPIRADO → MongoDB L2 (TTL: %v, tentando APIs...)\n", cep, time.Since(cachedTime))
			}
		} else {
			fmt.Printf("⚠️ [CEP:%s] CACHE MISS → MongoDB L2 (tentando APIs externas...)\n", cep)
		}
	}

	// 🌐 CAMADA 3: VIACEP (API Externa, ~100ms)
	fmt.Printf("🌐 [CEP:%s] Buscando em ViaCEP (API externa)...\n", cep)
	response, err := h.fetchViaCEP(cep)
	if err == nil && response.CEP != "" {
		response.Source = "viacep"
		response.CachedAt = time.Now().Format(time.RFC3339)

		fmt.Printf("✅ [CEP:%s] SUCESSO → ViaCEP (salvando em caches...)\n", cep)

		// ✅ NORMALIZAR CEP para salvar sem traço no cache
		response.CEP = strings.ReplaceAll(response.CEP, "-", "")
		response.CEP = strings.ReplaceAll(response.CEP, ".", "")

		// Salvar em AMBAS camadas de cache (se habilitado)
		if settings.Cache.Enabled {
			// ⚡ Salvar no Redis (L1 - hot cache, 24h)
			if h.redis != nil {
				if redisClient, ok := h.redis.(*cache.RedisClient); ok {
					redisKey := fmt.Sprintf("cep:%s", cep)
					ttl := h.getTTL(c)
					if err := redisClient.Set(ctx, redisKey, response, ttl); err != nil {
						fmt.Printf("⚠️ [CEP:%s] Erro ao salvar no Redis: %v\n", cep, err)
					} else {
						fmt.Printf("✅ [CEP:%s] Salvo no Redis L1 (TTL: %v)\n", cep, ttl)
					}
				}
			}

			// 🗄️ Salvar no MongoDB (L2 - cold cache, 7 dias)
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"cep": cep},
				bson.M{"$set": response},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				fmt.Printf("⚠️ [CEP:%s] Erro ao salvar no MongoDB: %v\n", cep, err)
			} else {
				fmt.Printf("✅ [CEP:%s] Salvo no MongoDB L2 (TTL: %d dias)\n", cep, settings.Cache.CEPTTLDays)
			}
		}

		c.JSON(http.StatusOK, response)
		return
	}

	fmt.Printf("⚠️ [CEP:%s] ERRO em ViaCEP: %v (tentando Brasil API...)\n", cep, err)

	// 🌐 CAMADA 3 (Fallback): BRASIL API (~150ms)
	fmt.Printf("🌐 [CEP:%s] Buscando em Brasil API (fallback)...\n", cep)
	response, err = h.fetchBrasilAPI(cep)
	if err == nil && response.CEP != "" {
		response.Source = "brasilapi"
		response.CachedAt = time.Now().Format(time.RFC3339)

		fmt.Printf("✅ [CEP:%s] SUCESSO → Brasil API (salvando em caches...)\n", cep)

		// ✅ NORMALIZAR CEP para salvar sem traço no cache
		response.CEP = strings.ReplaceAll(response.CEP, "-", "")
		response.CEP = strings.ReplaceAll(response.CEP, ".", "")

		// Salvar em AMBAS camadas de cache (se habilitado)
		if settings.Cache.Enabled {
			// ⚡ Salvar no Redis (L1 - hot cache, 24h)
			if h.redis != nil {
				if redisClient, ok := h.redis.(*cache.RedisClient); ok {
					redisKey := fmt.Sprintf("cep:%s", cep)
					ttl := h.getTTL(c)
					if err := redisClient.Set(ctx, redisKey, response, ttl); err != nil {
						fmt.Printf("⚠️ [CEP:%s] Erro ao salvar no Redis: %v\n", cep, err)
					} else {
						fmt.Printf("✅ [CEP:%s] Salvo no Redis L1 (TTL: %v)\n", cep, ttl)
					}
				}
			}

			// 🗄️ Salvar no MongoDB (L2 - cold cache, 7 dias)
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"cep": cep},
				bson.M{"$set": response},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				fmt.Printf("⚠️ [CEP:%s] Erro ao salvar no MongoDB: %v\n", cep, err)
			} else {
				fmt.Printf("✅ [CEP:%s] Salvo no MongoDB L2 (TTL: %d dias)\n", cep, settings.Cache.CEPTTLDays)
			}
		}

		c.JSON(http.StatusOK, response)
		return
	}

	fmt.Printf("❌ [CEP:%s] ERRO em Brasil API: %v (nenhuma fonte disponível)\n", cep, err)

	// 4. CEP não encontrado
	c.JSON(http.StatusNotFound, gin.H{
		"type":   "https://retech-core/errors/not-found",
		"title":  "CEP Not Found",
		"status": http.StatusNotFound,
		"detail": fmt.Sprintf("CEP %s não encontrado", cep),
	})
}

// fetchViaCEP busca CEP no ViaCEP
func (h *CEPHandler) fetchViaCEP(cep string) (*CEPResponse, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CEPResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// ViaCEP retorna {"erro": true} quando CEP não existe
	if result.CEP == "" {
		return nil, fmt.Errorf("CEP não encontrado")
	}

	return &result, nil
}

// fetchBrasilAPI busca CEP no Brasil API
func (h *CEPHandler) fetchBrasilAPI(cep string) (*CEPResponse, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CEP não encontrado")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Brasil API tem campos diferentes, precisamos mapear
	var brasilAPIResp struct {
		CEP          string `json:"cep"`
		State        string `json:"state"`
		City         string `json:"city"`
		Neighborhood string `json:"neighborhood"`
		Street       string `json:"street"`
	}

	if err := json.Unmarshal(body, &brasilAPIResp); err != nil {
		return nil, err
	}

	// Mapear para nosso formato
	result := &CEPResponse{
		CEP:        brasilAPIResp.CEP,
		Logradouro: brasilAPIResp.Street,
		Bairro:     brasilAPIResp.Neighborhood,
		Localidade: brasilAPIResp.City,
		UF:         brasilAPIResp.State,
	}

	return result, nil
}

// GetStats retorna estatísticas da API de CEP (para analytics)
func (h *CEPHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()
	collection := h.db.DB.Collection("api_usage_logs")

	// Total de consultas CEP
	total, _ := collection.CountDocuments(ctx, bson.M{
		"api_name": "cep",
	})

	// Consultas hoje (timezone Brasília)
	today := utils.GetTodayBrasilia()
	today_count, _ := collection.CountDocuments(ctx, bson.M{
		"api_name": "cep",
		"date":     today,
	})

	// Tempo médio de resposta
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"api_name": "cep"}}},
		{{Key: "$group", Value: bson.M{
			"_id":             nil,
			"avgResponseTime": bson.M{"$avg": "$responseTime"},
		}}},
	}

	cursor, _ := collection.Aggregate(ctx, pipeline)
	var avgResult []struct {
		AvgResponseTime float64 `bson:"avgResponseTime"`
	}
	cursor.All(ctx, &avgResult)

	avgTime := 0.0
	if len(avgResult) > 0 {
		avgTime = avgResult[0].AvgResponseTime
	}

	c.JSON(http.StatusOK, gin.H{
		"api":             "cep",
		"totalRequests":   total,
		"requestsToday":   today_count,
		"avgResponseTime": avgTime,
	})
}

// GetCacheStats retorna estatísticas do cache de CEP
// GET /admin/cache/cep/stats
func (h *CEPHandler) GetCacheStats(c *gin.Context) {
	ctx := c.Request.Context()
	collection := h.db.DB.Collection("cep_cache")

	// Total de CEPs no cache
	totalCached, _ := collection.CountDocuments(ctx, bson.M{})

	// CEPs adicionados nas últimas 24h
	yesterday := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	recentCached, _ := collection.CountDocuments(ctx, bson.M{
		"cachedAt": bson.M{"$gte": yesterday},
	})

	// Carregar configurações
	settings, err := h.settings.Get(ctx)
	if err != nil {
		settings = domain.GetDefaultSettings()
	}

	c.JSON(http.StatusOK, gin.H{
		"totalCached":  totalCached,
		"recentCached": recentCached, // últimas 24h
		"cacheEnabled": settings.Cache.Enabled,
		"cacheTTLDays": settings.Cache.CEPTTLDays,
		"autoCleanup":  settings.Cache.AutoCleanup,
	})
}

// ClearCache limpa o cache de CEP manualmente
// DELETE /admin/cache/cep
func (h *CEPHandler) ClearCache(c *gin.Context) {
	ctx := c.Request.Context()
	collection := h.db.DB.Collection("cep_cache")

	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao limpar cache",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Cache limpo com sucesso",
		"deletedCount": result.DeletedCount,
	})
}
