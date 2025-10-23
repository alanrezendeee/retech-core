package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theretech/retech-core/internal/domain"
	"github.com/theretech/retech-core/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CNPJHandler struct {
	db       *storage.Mongo
	settings *storage.SettingsRepo
}

func NewCNPJHandler(db *storage.Mongo, settings *storage.SettingsRepo) *CNPJHandler {
	return &CNPJHandler{
		db:       db,
		settings: settings,
	}
}

// GetCNPJ busca informações de um CNPJ
// GET /cnpj/:numero
func (h *CNPJHandler) GetCNPJ(c *gin.Context) {
	cnpjParam := c.Param("numero")

	// Normalizar CNPJ (remover pontos, barras, traços)
	cnpj := domain.NormalizeCNPJ(cnpjParam)

	// Validar CNPJ
	if !domain.ValidateCNPJ(cnpj) {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://retech-core/errors/validation",
			"title":  "CNPJ Inválido",
			"status": http.StatusBadRequest,
			"detail": "CNPJ deve ter 14 dígitos válidos",
		})
		return
	}

	ctx := c.Request.Context()

	// Carregar configurações de cache
	settings, err := h.settings.Get(ctx)
	if err != nil {
		settings = domain.GetDefaultSettings() // Fallback para padrões
	}

	// 1. Tentar buscar no cache (se habilitado)
	collection := h.db.DB.Collection("cnpj_cache")

	if settings.Cache.Enabled {
		var cached domain.CNPJ
		err := collection.FindOne(ctx, bson.M{"cnpj": cnpj}).Decode(&cached)
		if err == nil {
			// Verificar se cache ainda é válido (usar TTL configurável)
			cacheTTL := time.Duration(settings.Cache.CNPJTTLDays) * 24 * time.Hour
			
			if time.Since(cached.CachedAt) < cacheTTL {
				cached.Source = "cache"
				c.JSON(http.StatusOK, cached)
				return
			}
		}
	}

	// 2. Buscar em Brasil API (fonte principal)
	cnpjData, err := h.fetchBrasilAPI(ctx, cnpj)
	if err == nil && cnpjData.CNPJ != "" {
		cnpjData.Source = "brasilapi"
		cnpjData.CachedAt = time.Now().UTC()
		
		// ✅ NORMALIZAR CNPJ para salvar sem formatação
		cnpjData.CNPJ = domain.NormalizeCNPJ(cnpjData.CNPJ)

		// Salvar no cache (se habilitado)
		if settings.Cache.Enabled {
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"cnpj": cnpj}, // cnpj já está normalizado
				bson.M{"$set": cnpjData},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				fmt.Printf("⚠️ Erro ao salvar CNPJ no cache: %v\n", err)
			}
		}

		c.JSON(http.StatusOK, cnpjData)
		return
	}

	// 3. Fallback: ReceitaWS
	cnpjData, err = h.fetchReceitaWS(ctx, cnpj)
	if err == nil && cnpjData.CNPJ != "" {
		cnpjData.Source = "receitaws"
		cnpjData.CachedAt = time.Now().UTC()
		
		// ✅ NORMALIZAR CNPJ para salvar sem formatação
		cnpjData.CNPJ = domain.NormalizeCNPJ(cnpjData.CNPJ)

		// Salvar no cache (se habilitado)
		if settings.Cache.Enabled {
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"cnpj": cnpj}, // cnpj já está normalizado
				bson.M{"$set": cnpjData},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				fmt.Printf("⚠️ Erro ao salvar CNPJ no cache: %v\n", err)
			}
		}

		c.JSON(http.StatusOK, cnpjData)
		return
	}

	// 4. CNPJ não encontrado em nenhuma fonte
	c.JSON(http.StatusNotFound, gin.H{
		"type":   "https://retech-core/errors/not-found",
		"title":  "CNPJ Not Found",
		"status": http.StatusNotFound,
		"detail": fmt.Sprintf("CNPJ %s não encontrado ou indisponível", cnpj),
	})
}

// fetchBrasilAPI busca CNPJ na Brasil API
func (h *CNPJHandler) fetchBrasilAPI(ctx context.Context, cnpj string) (*domain.CNPJ, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cnpj/v1/%s", cnpj)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Brasil API retorna estrutura diferente, precisamos mapear
	var brasilAPIResp struct {
		CNPJ              string  `json:"cnpj"`
		RazaoSocial       string  `json:"razao_social"`
		NomeFantasia      string  `json:"nome_fantasia"`
		DescricaoSituacao string  `json:"descricao_situacao_cadastral"`
		DataSituacao      string  `json:"data_situacao_cadastral"`
		DataAbertura      string  `json:"data_inicio_atividade"`
		DescricaoPorte    string  `json:"porte"`
		NaturezaJuridica  string  `json:"descricao_natureza_juridica"`
		CapitalSocial     float64 `json:"capital_social"`
		Logradouro        string  `json:"logradouro"`
		Numero            string  `json:"numero"`
		Complemento       string  `json:"complemento"`
		Bairro            string  `json:"bairro"`
		CEP               string  `json:"cep"`
		Municipio         string  `json:"municipio"`
		UF                string  `json:"uf"`
		DDD1              string  `json:"ddd_telefone_1"`
		DDD2              string  `json:"ddd_telefone_2"`
		Email             string  `json:"email"`
		CNAEFiscal        struct {
			Codigo    string `json:"codigo"`
			Descricao string `json:"descricao"`
		} `json:"cnae_fiscal"`
		CNAEFiscalSecundarios []struct {
			Codigo    string `json:"codigo"`
			Descricao string `json:"descricao"`
		} `json:"cnaes_secundarios"`
		QSA []struct {
			Nome                 string `json:"nome_socio"`
			QualificacaoSocio    string `json:"qualificacao_socio"`
		} `json:"qsa"`
	}

	if err := json.Unmarshal(body, &brasilAPIResp); err != nil {
		return nil, err
	}

	// Mapear para nosso formato
	result := &domain.CNPJ{
		CNPJ:             brasilAPIResp.CNPJ,
		RazaoSocial:      brasilAPIResp.RazaoSocial,
		NomeFantasia:     brasilAPIResp.NomeFantasia,
		Situacao:         brasilAPIResp.DescricaoSituacao,
		DataSituacao:     brasilAPIResp.DataSituacao,
		DataAbertura:     brasilAPIResp.DataAbertura,
		Porte:            brasilAPIResp.DescricaoPorte,
		NaturezaJuridica: brasilAPIResp.NaturezaJuridica,
		CapitalSocial:    brasilAPIResp.CapitalSocial,
		Endereco: domain.CNPJEndereco{
			Logradouro:  brasilAPIResp.Logradouro,
			Numero:      brasilAPIResp.Numero,
			Complemento: brasilAPIResp.Complemento,
			Bairro:      brasilAPIResp.Bairro,
			CEP:         brasilAPIResp.CEP,
			Municipio:   brasilAPIResp.Municipio,
			UF:          brasilAPIResp.UF,
		},
		Email: brasilAPIResp.Email,
		AtividadePrincipal: domain.CNPJAtividade{
			Codigo:    brasilAPIResp.CNAEFiscal.Codigo,
			Descricao: brasilAPIResp.CNAEFiscal.Descricao,
		},
	}

	// Telefones
	if brasilAPIResp.DDD1 != "" {
		result.Telefones = append(result.Telefones, brasilAPIResp.DDD1)
	}
	if brasilAPIResp.DDD2 != "" {
		result.Telefones = append(result.Telefones, brasilAPIResp.DDD2)
	}

	// Atividades secundárias
	for _, cnae := range brasilAPIResp.CNAEFiscalSecundarios {
		result.AtividadesSecundarias = append(result.AtividadesSecundarias, domain.CNPJAtividade{
			Codigo:    cnae.Codigo,
			Descricao: cnae.Descricao,
		})
	}

	// QSA (sócios)
	for _, socio := range brasilAPIResp.QSA {
		result.QSA = append(result.QSA, domain.CNPJSocio{
			Nome:         socio.Nome,
			Qualificacao: socio.QualificacaoSocio,
		})
	}

	return result, nil
}

// fetchReceitaWS busca CNPJ na ReceitaWS (fallback)
func (h *CNPJHandler) fetchReceitaWS(ctx context.Context, cnpj string) (*domain.CNPJ, error) {
	url := fmt.Sprintf("https://www.receitaws.com.br/v1/cnpj/%s", cnpj)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// ReceitaWS retorna estrutura diferente
	var receitaResp struct {
		Status         string  `json:"status"`
		CNPJ           string  `json:"cnpj"`
		Nome           string  `json:"nome"`
		Fantasia       string  `json:"fantasia"`
		Situacao       string  `json:"situacao"`
		DataSituacao   string  `json:"data_situacao"`
		Abertura       string  `json:"abertura"`
		Porte          string  `json:"porte"`
		Natureza       string  `json:"natureza_juridica"`
		CapitalSocial  string  `json:"capital_social"`
		Logradouro     string  `json:"logradouro"`
		Numero         string  `json:"numero"`
		Complemento    string  `json:"complemento"`
		Bairro         string  `json:"bairro"`
		CEP            string  `json:"cep"`
		Municipio      string  `json:"municipio"`
		UF             string  `json:"uf"`
		Telefone       string  `json:"telefone"`
		Email          string  `json:"email"`
		AtividadePrincipal []struct {
			Code string `json:"code"`
			Text string `json:"text"`
		} `json:"atividade_principal"`
		AtividadesSecundarias []struct {
			Code string `json:"code"`
			Text string `json:"text"`
		} `json:"atividades_secundarias"`
		QSA []struct {
			Nome string `json:"nome"`
			Qual string `json:"qual"`
		} `json:"qsa"`
	}

	if err := json.Unmarshal(body, &receitaResp); err != nil {
		return nil, err
	}

	// Se status é ERROR, CNPJ não existe
	if receitaResp.Status == "ERROR" {
		return nil, fmt.Errorf("cnpj not found")
	}

	// Converter capital social (string → float64)
	var capitalSocial float64
	fmt.Sscanf(receitaResp.CapitalSocial, "%f", &capitalSocial)

	result := &domain.CNPJ{
		CNPJ:             receitaResp.CNPJ,
		RazaoSocial:      receitaResp.Nome,
		NomeFantasia:     receitaResp.Fantasia,
		Situacao:         receitaResp.Situacao,
		DataSituacao:     receitaResp.DataSituacao,
		DataAbertura:     receitaResp.Abertura,
		Porte:            receitaResp.Porte,
		NaturezaJuridica: receitaResp.Natureza,
		CapitalSocial:    capitalSocial,
		Endereco: domain.CNPJEndereco{
			Logradouro:  receitaResp.Logradouro,
			Numero:      receitaResp.Numero,
			Complemento: receitaResp.Complemento,
			Bairro:      receitaResp.Bairro,
			CEP:         receitaResp.CEP,
			Municipio:   receitaResp.Municipio,
			UF:          receitaResp.UF,
		},
		Email: receitaResp.Email,
	}

	// Telefone
	if receitaResp.Telefone != "" {
		result.Telefones = []string{receitaResp.Telefone}
	}

	// Atividade principal
	if len(receitaResp.AtividadePrincipal) > 0 {
		result.AtividadePrincipal = domain.CNPJAtividade{
			Codigo:    receitaResp.AtividadePrincipal[0].Code,
			Descricao: receitaResp.AtividadePrincipal[0].Text,
		}
	}

	// Atividades secundárias
	for _, ativ := range receitaResp.AtividadesSecundarias {
		result.AtividadesSecundarias = append(result.AtividadesSecundarias, domain.CNPJAtividade{
			Codigo:    ativ.Code,
			Descricao: ativ.Text,
		})
	}

	// QSA
	for _, socio := range receitaResp.QSA {
		result.QSA = append(result.QSA, domain.CNPJSocio{
			Nome:         socio.Nome,
			Qualificacao: socio.Qual,
		})
	}

	return result, nil
}

// GetCacheStats retorna estatísticas do cache de CNPJ
// GET /admin/cache/cnpj/stats
func (h *CNPJHandler) GetCacheStats(c *gin.Context) {
	ctx := c.Request.Context()
	collection := h.db.DB.Collection("cnpj_cache")

	// Total de CNPJs no cache
	totalCached, _ := collection.CountDocuments(ctx, bson.M{})

	// CNPJs adicionados nas últimas 24h
	yesterday := time.Now().Add(-24 * time.Hour)
	recentCached, _ := collection.CountDocuments(ctx, bson.M{
		"cachedAt": bson.M{"$gte": yesterday},
	})

	// Carregar configurações
	settings, err := h.settings.Get(ctx)
	if err != nil {
		settings = domain.GetDefaultSettings()
	}

	c.JSON(http.StatusOK, gin.H{
		"totalCached":   totalCached,
		"recentCached":  recentCached, // últimas 24h
		"cacheEnabled":  settings.Cache.Enabled,
		"cacheTTLDays":  settings.Cache.CNPJTTLDays, // ✅ TTL configurável
		"autoCleanup":   settings.Cache.AutoCleanup,
	})
}

// ClearCache limpa o cache de CNPJ manualmente
// DELETE /admin/cache/cnpj
func (h *CNPJHandler) ClearCache(c *gin.Context) {
	ctx := c.Request.Context()
	collection := h.db.DB.Collection("cnpj_cache")

	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao limpar cache de CNPJ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Cache de CNPJ limpo com sucesso",
		"deletedCount": result.DeletedCount,
	})
}

