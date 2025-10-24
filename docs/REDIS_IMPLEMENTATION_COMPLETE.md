# ⚡ Implementação Redis - Cache em 3 Camadas

## 📊 Resumo Executivo

Implementamos **Redis como cache L1** em todas as 3 APIs disponíveis (CEP, CNPJ, Geografia), criando uma arquitetura de **cache em 3 camadas** para atingir **latência <5ms** em produção.

---

## 🏗️ Arquitetura de Cache em 3 Camadas

```
┌─────────────────────────────────────────────────────┐
│                   CLIENT REQUEST                      │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│  ⚡ CAMADA 1: REDIS (L1 - Hot Cache)                 │
│  • Latência: <1ms                                     │
│  • TTL: 24 horas                                      │
│  • Objetivo: Atingir <5ms em produção                │
└─────────────────────────────────────────────────────┘
                         ↓ (cache miss)
┌─────────────────────────────────────────────────────┐
│  🗄️ CAMADA 2: MONGODB (L2 - Cold Cache)             │
│  • Latência: ~10-50ms                                 │
│  • TTL: 7 dias (CEP) / 30 dias (CNPJ)                │
│  • Promoção: Dados são promovidos para Redis         │
└─────────────────────────────────────────────────────┘
                         ↓ (cache miss)
┌─────────────────────────────────────────────────────┐
│  🌐 CAMADA 3: APIS EXTERNAS (Origin)                 │
│  • CEP: ViaCEP → Brasil API (fallback)               │
│  • CNPJ: Brasil API → ReceitaWS (fallback)           │
│  • Geo: MongoDB (dados fixos)                        │
│  • Latência: 100-300ms                                │
│  • Write-through: Salva em L1 (Redis) + L2 (Mongo)   │
└─────────────────────────────────────────────────────┘
```

---

## 🚀 APIs Implementadas

### 1. **API CEP** (`/cep/:codigo`)
- ⚡ **Redis L1**: 24h TTL, chave `cep:{codigo}`
- 🗄️ **MongoDB L2**: TTL configurável (padrão 7 dias)
- 🌐 **Origin**: ViaCEP → Brasil API (fallback)
- 🎯 **Latência esperada**: <1ms (cached) vs 100-200ms (origin)

### 2. **API CNPJ** (`/cnpj/:numero`)
- ⚡ **Redis L1**: 24h TTL, chave `cnpj:{numero}`
- 🗄️ **MongoDB L2**: TTL configurável (padrão 30 dias)
- 🌐 **Origin**: Brasil API → ReceitaWS (fallback)
- 🎯 **Latência esperada**: <1ms (cached) vs 200-400ms (origin)

### 3. **API Geografia** (`/geo/ufs`, `/geo/ufs/:sigla`, `/geo/municipios`)
- ⚡ **Redis L1**: 24h TTL
  - `geo:ufs:all` (lista completa)
  - `geo:uf:{sigla}` (estado individual)
- 🗄️ **MongoDB**: Dados fixos (5.570 municípios + 27 UFs)
- 🎯 **Latência esperada**: <1ms (cached) vs 10-30ms (MongoDB)

---

## 🔧 Detalhes Técnicos

### **Graceful Degradation**
```go
type CEPHandler struct {
    db       *storage.Mongo
    redis    interface{} // Permite nil sem quebrar
    settings *storage.SettingsRepo
}
```
- Se Redis falhar ou não estiver disponível, **não quebra** a aplicação
- Fallback automático para MongoDB (L2) → API Externa (L3)

### **Write-Through Cache**
Quando dados são buscados da API externa:
1. ✅ **Salva no Redis** (L1 - hot cache, 24h)
2. ✅ **Salva no MongoDB** (L2 - cold cache, 7-30 dias)
3. ✅ **Retorna ao cliente**

### **Cache Promotion**
Quando dados são encontrados no MongoDB mas não no Redis:
1. ✅ **Promove para Redis** (L1)
2. ✅ **Retorna ao cliente**
3. 🎯 **Próxima request será <1ms!**

---

## 📈 Melhorias de Performance Esperadas

| API   | Antes (MongoDB) | Depois (Redis) | Melhoria |
|-------|----------------|----------------|----------|
| CEP   | ~160ms         | **<5ms**       | **97%** ↓ |
| CNPJ  | ~180ms         | **<5ms**       | **97%** ↓ |
| Geo   | ~30ms          | **<2ms**       | **93%** ↓ |

---

## 🔌 Infraestrutura Railway

### **Redis Service**
- ✅ **Instância dedicada** já provisionada no Railway
- ✅ **Variável de ambiente**: `REDIS_URL`
- ✅ **Conexão automática** via `internal/cache/redis_client.go`
- ✅ **Fallback**: Se Redis falhar, usa MongoDB

### **Configuração**
```go
// cmd/api/main.go
redisURL := os.Getenv("REDIS_URL")
redisClient, err := cache.NewRedisClient(redisURL)
if err != nil {
    log.Printf("⚠️ Redis não disponível, usando apenas MongoDB: %v", err)
    redisClient = nil
}
```

---

## 📊 Monitoramento

### **Logs de Debug**
```bash
# No servidor, verificar se Redis está conectado
grep "Redis" logs.txt

# Sucesso:
⚡ Redis conectado com sucesso!

# Fallback:
⚠️ Redis não disponível, usando apenas MongoDB
```

### **Verificar Cache Hit Rate**
```bash
# Testar CEP
curl https://api-core.theretech.com.br/cep/88111477 -H "x-api-key: YOUR_KEY"
# 1ª request: "source": "viacep" (origin)
# 2ª request: "source": "redis-cache" (⚡ <1ms!)
```

---

## ✅ Checklist de Deploy

- [x] Código implementado em CEP, CNPJ, Geografia
- [x] Compilação Go bem-sucedida
- [x] Graceful degradation testado
- [x] Instância Redis provisionada no Railway
- [x] Variável `REDIS_URL` configurada
- [ ] Deploy em produção
- [ ] Teste de latência (<5ms esperado)
- [ ] Monitoramento de cache hit rate
- [ ] Atualizar documentação do Redoc (adicionar `source` field)

---

## 🎯 Próximos Passos

1. **Deploy**: Fazer push do código e rebuild no Railway
2. **Teste**: Verificar latência real em produção
3. **Monitoramento**: Acompanhar cache hit rate e performance
4. **Otimização**: Ajustar TTLs se necessário
5. **Documentação**: Atualizar Redoc com novo campo `source`

---

## 📚 Arquivos Modificados

### Backend (Go)
- ✅ `cmd/api/main.go` - Inicializa Redis client
- ✅ `internal/cache/redis_client.go` - Novo cliente Redis
- ✅ `internal/http/router.go` - Passa Redis para handlers
- ✅ `internal/http/handlers/cep.go` - Cache L1 (Redis) + L2 (Mongo)
- ✅ `internal/http/handlers/cnpj.go` - Cache L1 (Redis) + L2 (Mongo)
- ✅ `internal/http/handlers/geo.go` - Cache L1 (Redis) para geo dados
- ✅ `internal/cache/settings_cache.go` - Fix return type
- ✅ `go.mod` - Adiciona `github.com/redis/go-redis/v9`

### Documentação
- ✅ `docs/REDIS_IMPLEMENTATION.md` - Estratégia inicial
- ✅ `docs/REDIS_IMPLEMENTATION_COMPLETE.md` - Este documento

---

## 💡 Observações Finais

- **Zero Breaking Changes**: Todas as APIs continuam funcionando mesmo sem Redis
- **Performance**: Esperamos reduzir latência de ~160ms para **<5ms**
- **Escalabilidade**: Redis suporta milhões de requests/segundo
- **Custo**: Instância Redis Railway ~$5/mês

**Status**: ✅ PRONTO PARA DEPLOY! 🚀

