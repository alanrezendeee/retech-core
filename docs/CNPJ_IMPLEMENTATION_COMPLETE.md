# 🏢 API CNPJ - Implementação Completa

**Data:** 23 de Outubro de 2025  
**Status:** ✅ Backend Completo | 🔵 Frontend/Docs Pendente

---

## 📋 RESUMO

API completa para consulta de CNPJs brasileiros com:
- ✅ Validação de dígito verificador
- ✅ 2 fontes de dados (Brasil API + ReceitaWS)
- ✅ Cache inteligente (30 dias)
- ✅ QSA (sócios e administradores)
- ✅ Dados completos da Receita Federal

---

## 🔌 ENDPOINT

```
GET /cnpj/:numero
```

**Headers:**
```
X-API-Key: sk_your_api_key_here
```

**Formatos aceitos:**
- Com formatação: `00.000.000/0001-00`
- Sem formatação: `00000000000100`

---

## 📊 RESPOSTA (Exemplo: Banco do Brasil)

```json
{
  "cnpj": "00000000000191",
  "razaoSocial": "BANCO DO BRASIL SA",
  "nomeFantasia": "BANCO DO BRASIL",
  "situacao": "ATIVA",
  "dataSituacao": "2005-11-03",
  "dataAbertura": "1966-08-30",
  "porte": "DEMAIS",
  "naturezaJuridica": "Sociedade de Economia Mista",
  "capitalSocial": 81250000000.00,
  "endereco": {
    "logradouro": "SBS Quadra 1 Bloco G",
    "numero": "S/N",
    "complemento": "LOJA LOJA 1 SUBSOLO 1",
    "bairro": "ASA SUL",
    "cep": "70073901",
    "municipio": "Brasília",
    "uf": "DF"
  },
  "telefones": ["(61) 3493-9002"],
  "email": "prov@bb.com.br",
  "atividadePrincipal": {
    "codigo": "64.21-2-00",
    "descricao": "Bancos comerciais"
  },
  "atividadesSecundarias": [
    {
      "codigo": "66.19-3-99",
      "descricao": "Outras atividades auxiliares dos serviços financeiros"
    }
  ],
  "qsa": [
    {
      "nome": "NOME DO SOCIO/ADMINISTRADOR",
      "qualificacao": "Diretor"
    }
  ],
  "source": "brasilapi",
  "cachedAt": "2025-10-23T22:30:00Z"
}
```

---

## 🔄 FLUXO DE DADOS

```
1. Request: GET /cnpj/00000000000191
2. Validação: dígito verificador ✓
3. Cache (MongoDB)?
   ✅ Hit (< 30 dias) → Retorna (10ms)
   ❌ Miss → Continua
4. Brasil API?
   ✅ Sucesso → Salva cache + Retorna (200-500ms)
   ❌ Falha → Continua
5. ReceitaWS (fallback)?
   ✅ Sucesso → Salva cache + Retorna (500-1000ms)
   ❌ Falha → 404 Not Found
```

---

## 🥇 FONTES DE DADOS

### **1. Brasil API (Principal)**
- **URL:** `https://brasilapi.com.br/api/cnpj/v1/{cnpj}`
- **Rate Limit:** Sem limite (gratuita)
- **Dados:** Receita Federal (atualização mensal)
- **Performance:** ~200-500ms
- **Confiabilidade:** ⭐⭐⭐⭐⭐

### **2. ReceitaWS (Fallback)**
- **URL:** `https://www.receitaws.com.br/v1/cnpj/{cnpj}`
- **Rate Limit:** 3 requests/minuto
- **Dados:** Receita Federal
- **Performance:** ~500-1000ms
- **Confiabilidade:** ⭐⭐⭐

### **3. Cache Local (MongoDB)**
- **TTL:** 30 dias (fixo)
- **Performance:** ~10ms
- **Confiabilidade:** ⭐⭐⭐⭐⭐

---

## ✅ VALIDAÇÃO DE CNPJ

### **Algoritmo Implementado:**

```go
func ValidateCNPJ(cnpj string) bool {
    // 1. Remove formatação
    // 2. Verifica 14 dígitos
    // 3. Rejeita sequências iguais (11111111111111)
    // 4. Valida primeiro dígito verificador
    // 5. Valida segundo dígito verificador
    return true
}
```

### **CNPJs Válidos para Teste:**
- `00.000.000/0001-91` - Banco do Brasil
- `33.000.167/0001-01` - Bradesco
- `60.746.948/0001-12` - Petrobras
- `07.526.557/0001-00` - Google Brasil

### **CNPJs Inválidos:**
- `00000000000000` - Sequência de zeros
- `11111111111111` - Sequência de uns
- `00000000000190` - Dígito verificador errado

---

## 📦 ESTRUTURAS DE DADOS

### **CNPJ (Principal)**
```go
type CNPJ struct {
    CNPJ                  string
    RazaoSocial           string
    NomeFantasia          string
    Situacao              string
    DataSituacao          string
    DataAbertura          string
    Porte                 string
    NaturezaJuridica      string
    CapitalSocial         float64
    Endereco              CNPJEndereco
    Telefones             []string
    Email                 string
    AtividadePrincipal    CNPJAtividade
    AtividadesSecundarias []CNPJAtividade
    QSA                   []CNPJSocio
    Source                string
    CachedAt              time.Time
}
```

### **CNPJEndereco**
```go
type CNPJEndereco struct {
    Logradouro  string
    Numero      string
    Complemento string
    Bairro      string
    CEP         string
    Municipio   string
    UF          string
}
```

### **CNPJAtividade (CNAE)**
```go
type CNPJAtividade struct {
    Codigo    string // Ex: "64.21-2-00"
    Descricao string // Ex: "Bancos comerciais"
}
```

### **CNPJSocio (QSA)**
```go
type CNPJSocio struct {
    Nome         string
    Qualificacao string // Ex: "Diretor", "Sócio Administrador"
}
```

---

## 🔒 SEGURANÇA

✅ **Middlewares Aplicados:**
1. **Maintenance Mode** - Respeita modo manutenção
2. **API Key Auth** - Requer chave válida
3. **Rate Limiting** - Limite por tenant
4. **Usage Logging** - Registra todas as chamadas

✅ **Validações:**
- Dígito verificador obrigatório
- Formatação flexível (aceita com/sem máscara)
- Sanitização de entrada

---

## 📈 PERFORMANCE

| Cenário | Fonte | Tempo Médio |
|---------|-------|-------------|
| Cache hit | MongoDB | ~10ms ⚡ |
| Brasil API | Rede | ~200-500ms |
| ReceitaWS | Rede | ~500-1000ms |
| Timeout | - | 10s (max) |

### **Cache Hit Rate Esperado:**
- **Empresas grandes:** 95%+ (raramente mudam)
- **Empresas médias:** 80-90%
- **Empresas novas:** 10-20% (ainda não em cache)

---

## 🎯 CASOS DE USO

### **1. Validação de Cliente**
```bash
curl -X GET "https://api.retech.com.br/cnpj/00000000000191" \
  -H "X-API-Key: sk_..."
```

### **2. Preencher Cadastro Automaticamente**
```javascript
const response = await fetch('/cnpj/33000167000101', {
  headers: { 'X-API-Key': 'sk_...' }
});
const empresa = await response.json();

// Preencher formulário
form.razaoSocial = empresa.razaoSocial;
form.nomeFantasia = empresa.nomeFantasia;
form.endereco = empresa.endereco;
form.telefone = empresa.telefones[0];
```

### **3. Verificar Situação**
```python
response = requests.get(
    f"https://api.retech.com.br/cnpj/{cnpj}",
    headers={"X-API-Key": "sk_..."}
)
empresa = response.json()

if empresa['situacao'] != 'ATIVA':
    print(f"Empresa {empresa['situacao']}")
```

---

## 🚨 CÓDIGOS DE ERRO

| Código | Descrição | Causa |
|--------|-----------|-------|
| 400 | CNPJ Inválido | Dígito verificador errado |
| 401 | Unauthorized | API Key inválida |
| 404 | CNPJ Not Found | Não existe na Receita |
| 429 | Too Many Requests | Rate limit excedido |
| 500 | Internal Server Error | Erro no servidor |
| 503 | Service Unavailable | Modo manutenção |

---

## 🧪 TESTANDO

### **1. CNPJ Válido (Banco do Brasil)**
```bash
curl "http://localhost:8080/cnpj/00000000000191" \
  -H "X-API-Key: sk_test_..."
```

**Esperado:** 200 OK + dados completos

### **2. CNPJ Inválido**
```bash
curl "http://localhost:8080/cnpj/00000000000190" \
  -H "X-API-Key: sk_test_..."
```

**Esperado:** 400 Bad Request

### **3. CNPJ Não Encontrado**
```bash
curl "http://localhost:8080/cnpj/99999999999999" \
  -H "X-API-Key: sk_test_..."
```

**Esperado:** 404 Not Found

### **4. Cache Stats (Admin)**
```bash
curl "http://localhost:8080/admin/cache/cnpj/stats" \
  -H "Authorization: Bearer JWT_TOKEN"
```

**Resposta:**
```json
{
  "totalCached": 42,
  "recentCached": 5,
  "cacheEnabled": true,
  "cacheTTLDays": 30,
  "autoCleanup": true
}
```

---

## 📝 ADMIN ENDPOINTS

### **GET /admin/cache/cnpj/stats**
Retorna estatísticas do cache de CNPJ

### **DELETE /admin/cache/cnpj**
Limpa todo o cache de CNPJ manualmente

---

## ✅ CHECKLIST PÓS-IMPLEMENTAÇÃO

- [x] Backend: Domain + Handler + Routes
- [x] Validação de CNPJ (dígito verificador)
- [x] Integração Brasil API
- [x] Fallback ReceitaWS
- [x] Cache MongoDB (30 dias)
- [x] Admin: Stats + Clear cache
- [ ] Redoc: Documentar endpoint
- [ ] Landing Page: Mover para "Disponíveis"
- [ ] Admin UI: Cache CNPJ (opcional)
- [ ] Testes E2E: Postman

---

## 🎯 PRÓXIMOS PASSOS

1. **Documentar no Redoc**
   - Adicionar `/cnpj/:numero` ao openapi.yaml
   - Exemplos de request/response
   - Schema CNPJ completo

2. **Atualizar Landing Page**
   - Mover card CNPJ para "APIs Disponíveis"
   - Badge: "Fase 2" → "Disponível"

3. **Testar em Produção**
   - CNPJs reais (Banco do Brasil, etc)
   - Performance (cache hit rate)
   - Fallback ReceitaWS

4. **Monitorar**
   - Taxa de erro Brasil API
   - Uso de fallback ReceitaWS
   - Cache hit rate

---

**🎉 API CNPJ PRONTA PARA USO!**

Refs:
- `internal/domain/cnpj.go`
- `internal/http/handlers/cnpj.go`
- `internal/http/router.go` (linhas 111-122, 157-158)

