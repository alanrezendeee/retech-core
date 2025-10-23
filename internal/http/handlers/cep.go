package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CEPHandler struct {
	db *storage.Mongo
}

func NewCEPHandler(db *storage.Mongo) *CEPHandler {
	return &CEPHandler{db: db}
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

	// 1. Tentar buscar no cache (7 dias)
	var cached CEPResponse
	collection := h.db.DB.Collection("cep_cache")

	err := collection.FindOne(ctx, bson.M{"cep": cep}).Decode(&cached)
	if err == nil {
		// Verificar se cache ainda é válido (7 dias)
		cachedTime, _ := time.Parse(time.RFC3339, cached.CachedAt)
		if time.Since(cachedTime) < 7*24*time.Hour {
			cached.Source = "cache"
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	// 2. Buscar em ViaCEP (fonte principal)
	response, err := h.fetchViaCEP(cep)
	if err == nil && response.CEP != "" {
		response.Source = "viacep"
		response.CachedAt = time.Now().Format(time.RFC3339)

		// Salvar no cache
		collection.FindOneAndUpdate(
			ctx,
			bson.M{"cep": cep},
			bson.M{"$set": response},
			nil,
		)

		c.JSON(http.StatusOK, response)
		return
	}

	// 3. Fallback: Brasil API
	response, err = h.fetchBrasilAPI(cep)
	if err == nil && response.CEP != "" {
		response.Source = "brasilapi"
		response.CachedAt = time.Now().Format(time.RFC3339)

		// Salvar no cache
		collection.FindOneAndUpdate(
			ctx,
			bson.M{"cep": cep},
			bson.M{"$set": response},
			nil,
		)

		c.JSON(http.StatusOK, response)
		return
	}

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

	// Consultas hoje
	today := time.Now().Format("2006-01-02")
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
