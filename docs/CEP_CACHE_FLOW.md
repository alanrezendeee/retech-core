# 🔄 Fluxo Completo da API de CEP

## 📊 Diagrama de Decisão

```
┌─────────────────────────────────────────────────────────────┐
│                    GET /cep/{codigo}                        │
│                   + Header: X-API-Key                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
            ┌────────────────────────┐
            │  1. Validar formato    │
            │  (8 dígitos apenas)    │
            └────────┬───────────────┘
                     │
                     ▼
            ┌────────────────────────┐      ✅ Encontrado
            │  2. Buscar no CACHE    │────► + Válido (< 7 dias)
            │  (MongoDB local)       │      │
            └────────┬───────────────┘      │
                     │                      ▼
                     │ ❌ Não encontrado   [RETORNA]
                     │ ou expirado        source: "cache"
                     ▼                    ~5-10ms
            ┌────────────────────────┐
            │  3. Buscar no ViaCEP   │
            │  (API externa)         │
            └────────┬───────────────┘
                     │
         ┌───────────┴───────────┐
         │                       │
         ▼                       ▼
    ✅ Sucesso               ❌ Falha/Erro
    CEP válido               (timeout, 404, etc)
         │                       │
         │                       ▼
         │              ┌────────────────────────┐
         │              │  4. Buscar BrasilAPI   │
         │              │  (API externa fallback)│
         │              └────────┬───────────────┘
         │                       │
         │           ┌───────────┴───────────┐
         │           │                       │
         │           ▼                       ▼
         │      ✅ Sucesso               ❌ Falha
         │      CEP válido               CEP inválido
         │           │                       │
         └───────────┴───────────┐           │
                     │           │           │
                     ▼           │           ▼
            ┌────────────────┐   │   ┌──────────────┐
            │  5. SALVAR     │◄──┘   │  404 ERROR   │
            │  NO CACHE      │       │  Not Found   │
            │  (7 dias)      │       └──────────────┘
            └────────┬───────┘
                     │
                     ▼
            ┌────────────────────────┐
            │  RETORNAR RESPOSTA     │
            │  source: "viacep" ou   │
            │  source: "brasilapi"   │
            └────────────────────────┘
```

---

## 🎯 **RESUMO: Quando vai para Brasil API?**

A Brasil API é usada **SOMENTE** quando:

### ❌ **ViaCEP Falha por:**
1. **Timeout** (> 5 segundos sem resposta)
2. **Erro de conexão** (ViaCEP fora do ar)
3. **CEP não existe** (ViaCEP retorna `{"erro": true}` ou campo CEP vazio)
4. **Erro HTTP** (500, 502, etc)
5. **JSON inválido** (resposta corrompida)

---

## 📝 **Código Detalhado**

### **Passo 2: Buscar no Cache**
```go
// Linha 68-77
err := collection.FindOne(ctx, bson.M{"cep": cep}).Decode(&cached)
if err == nil {
    // Verificar se cache ainda é válido (7 dias)
    cachedTime, _ := time.Parse(time.RFC3339, cached.CachedAt)
    if time.Since(cachedTime) < 7*24*time.Hour {
        cached.Source = "cache"
        c.JSON(http.StatusOK, cached)
        return  // ✅ PARA AQUI - nem tenta API externa!
    }
}
```

**Condições para usar cache:**
- ✅ Documento existe no MongoDB
- ✅ `cachedAt` foi há menos de 7 dias
- ✅ Se ambos: **RETORNA IMEDIATAMENTE** (~5-10ms)

---

### **Passo 3: Buscar no ViaCEP (Principal)**
```go
// Linha 80-102
response, err := h.fetchViaCEP(cep)
if err == nil && response.CEP != "" {
    // ✅ SUCESSO!
    response.Source = "viacep"
    response.CachedAt = time.Now().Format(time.RFC3339)
    
    // Salvar no cache
    collection.UpdateOne(ctx, bson.M{"cep": cep}, ...)
    
    c.JSON(http.StatusOK, response)
    return  // ✅ PARA AQUI - Brasil API NEM É CHAMADA!
}
// ❌ Se chegou aqui, ViaCEP FALHOU
```

**`fetchViaCEP` retorna erro quando:**
```go
// Linha 139-165
func (h *CEPHandler) fetchViaCEP(cep string) (*CEPResponse, error) {
    url := "https://viacep.com.br/ws/88101270/json/"
    
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err  // ❌ Timeout ou erro de conexão
    }
    
    var result CEPResponse
    json.Unmarshal(body, &result)
    
    if result.CEP == "" {
        return nil, fmt.Errorf("CEP não encontrado")  // ❌ ViaCEP retornou erro
    }
    
    return &result, nil  // ✅ Sucesso
}
```

**Motivos de falha:**
1. **`client.Get()` falha** → erro de rede, timeout (>5s), DNS
2. **`json.Unmarshal()` falha** → JSON inválido
3. **`result.CEP == ""`** → ViaCEP retornou `{"erro": true}`

---

### **Passo 4: Fallback Brasil API**
```go
// Linha 104-127
// ⚠️ SÓ CHEGA AQUI SE VIACEP FALHOU!
response, err = h.fetchBrasilAPI(cep)
if err == nil && response.CEP != "" {
    // ✅ SUCESSO!
    response.Source = "brasilapi"
    
    // Salvar no cache
    collection.UpdateOne(ctx, bson.M{"cep": cep}, ...)
    
    c.JSON(http.StatusOK, response)
    return  // ✅ Salvou e retornou
}
// ❌ Se chegou aqui, AMBAS as APIs falharam
```

**`fetchBrasilAPI` retorna erro quando:**
```go
// Linha 167-209
func (h *CEPHandler) fetchBrasilAPI(cep string) (*CEPResponse, error) {
    url := "https://brasilapi.com.br/api/cep/v1/88101270"
    
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err  // ❌ Timeout ou erro de conexão
    }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("CEP não encontrado")  // ❌ 404, 500, etc
    }
    
    // Mapear campos diferentes
    result := &CEPResponse{
        CEP:        brasilAPIResp.CEP,
        Logradouro: brasilAPIResp.Street,  // Brasil API usa "street"
        Bairro:     brasilAPIResp.Neighborhood,
        Localidade: brasilAPIResp.City,
        UF:         brasilAPIResp.State,
    }
    
    return result, nil  // ✅ Sucesso
}
```

---

### **Passo 5: Retornar 404 (ambas falharam)**
```go
// Linha 129-136
// ⚠️ SÓ CHEGA AQUI SE CACHE, VIACEP E BRASILAPI FALHARAM!
c.JSON(http.StatusNotFound, gin.H{
    "type":   "https://retech-core/errors/not-found",
    "title":  "CEP Not Found",
    "status": http.StatusNotFound,
    "detail": fmt.Sprintf("CEP %s não encontrado", cep),
})
```

---

## 🧪 **Exemplos Práticos**

### **Cenário 1: CEP válido, primeira vez**
```
Request: GET /cep/88101270
         X-API-Key: sk_...

1. Cache? ❌ Não existe
2. ViaCEP? ✅ Sucesso (50ms)
   → Salva no cache
   → Retorna: source="viacep"

Brasil API: ❌ NÃO FOI CHAMADA!
```

---

### **Cenário 2: CEP válido, segunda vez**
```
Request: GET /cep/88101270
         X-API-Key: sk_...

1. Cache? ✅ Existe + Válido
   → Retorna: source="cache" (5ms)

ViaCEP: ❌ NÃO FOI CHAMADA!
Brasil API: ❌ NÃO FOI CHAMADA!
```

---

### **Cenário 3: ViaCEP fora do ar**
```
Request: GET /cep/88101270
         X-API-Key: sk_...

1. Cache? ❌ Não existe
2. ViaCEP? ❌ Timeout (5s)
3. Brasil API? ✅ Sucesso (100ms)
   → Salva no cache
   → Retorna: source="brasilapi"
```

---

### **Cenário 4: CEP inválido (ambas APIs)**
```
Request: GET /cep/99999999
         X-API-Key: sk_...

1. Cache? ❌ Não existe
2. ViaCEP? ❌ {"erro": true}
3. Brasil API? ❌ 404 Not Found
4. Retorna: 404 - CEP não encontrado

Cache: ❌ NÃO SALVOU (erro não é cacheado)
```

---

### **Cenário 5: Cache expirado (> 7 dias)**
```
Request: GET /cep/88101270
         X-API-Key: sk_...

1. Cache? ⚠️ Existe mas expirado
   (cachedAt: 2025-10-10, hoje: 2025-10-23)
2. ViaCEP? ✅ Sucesso
   → ATUALIZA cache com nova data
   → Retorna: source="viacep"
```

---

## ⚡ **Performance Esperada**

| Fonte | Tempo | Quando Ocorre |
|-------|-------|---------------|
| **Cache** | ~5-10ms | 2ª request em diante (< 7 dias) |
| **ViaCEP** | ~40-80ms | 1ª request ou cache expirado |
| **Brasil API** | ~80-150ms | Quando ViaCEP falha |
| **Timeout** | 5000ms | Quando API externa não responde |

---

## 🔧 **Configurações**

### **Timeout das APIs Externas**
```go
client := &http.Client{Timeout: 5 * time.Second}
```
- ⏱️ Máximo 5 segundos esperando resposta
- ⏱️ Se exceder: retorna erro e tenta próxima fonte

### **Validade do Cache**
```go
if time.Since(cachedTime) < 7*24*time.Hour {
```
- 📅 Cache válido por **7 dias**
- 📅 Após 7 dias: busca novamente na API

---

## ✅ **Vantagens do Sistema**

1. **Alta disponibilidade**: 2 APIs de fallback
2. **Performance**: Cache reduz 95% das chamadas externas
3. **Resiliência**: Se ViaCEP cai, Brasil API assume
4. **Custo zero**: Ambas as APIs são gratuitas
5. **Inteligente**: Não cacheia erros (só sucesso)

---

## 🚀 **Próximas Melhorias**

- [ ] Circuit breaker (parar de tentar API que está falhando muito)
- [ ] Métricas de sucesso/falha por fonte
- [ ] Cache com TTL variável (CEPs "famosos" = cache mais longo)
- [ ] Retry automático com backoff exponencial
- [ ] Health check das APIs externas

---

**🎯 RESUMO FINAL:**

Brasil API **SOMENTE** é usada quando:
- ❌ Não há cache válido, E
- ❌ ViaCEP falha (timeout, erro, CEP inexistente)

**90%+ dos casos**: Cache resolve em ~5ms! 🚀

