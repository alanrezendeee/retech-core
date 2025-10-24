# ⚡ Como Funciona o Redis - Explicação Completa

## 🎯 Resumo Simples

**Local (Dev):**
```
Backend (porta 8080) → Redis (porta 6379) → MongoDB
         ↑                    ↑
    docker-compose.yml    docker-compose.yml
```

**Produção (Railway):**
```
Backend Railway → REDIS_URL (env) → Redis Railway (instância dedicada) → MongoDB Railway
```

---

## 🏗️ Arquitetura Detalhada

### **1. Ambiente Local (Desenvolvimento)**

#### **docker-compose.yml**
```yaml
services:
  # Backend Go
  api:
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis://redis:6379  # ← Aponta para serviço Redis
  
  # Redis
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"  # ← Porta do Redis
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
```

**Como funciona:**
1. Backend roda na **porta 8080**
2. Redis roda na **porta 6379** (serviço separado)
3. Backend conecta ao Redis via hostname `redis:6379` (Docker DNS interno)
4. Se Redis falhar → Backend continua funcionando (graceful degradation)

---

### **2. Ambiente de Produção (Railway)**

#### **Serviços no Railway:**
```
┌─────────────────────────────────────────────┐
│  Serviço 1: retech-core (Backend Go)       │
│  • Variável: REDIS_URL                      │
│  • Valor: redis://default:SENHA@HOST:PORT  │
└─────────────────────────────────────────────┘
                    ↓ conecta
┌─────────────────────────────────────────────┐
│  Serviço 2: Redis (Instância Dedicada)     │
│  • Expõe REDIS_URL automaticamente          │
│  • Memória: 256MB                           │
│  • Policy: allkeys-lru (auto-evict)         │
└─────────────────────────────────────────────┘
```

**Como funciona:**
1. Você cria **2 serviços** no Railway:
   - `retech-core` (backend)
   - `Redis` (instância dedicada)

2. Railway gera automaticamente a variável `REDIS_URL` no serviço Redis

3. Você **copia** `REDIS_URL` para o serviço `retech-core`

4. Backend lê a variável e conecta ao Redis:
```go
redisURL := os.Getenv("REDIS_URL")
// redisURL = "redis://default:abc123@redis.railway.internal:6379"

client, err := cache.NewRedisClient(redisURL, "", 0, log)
```

---

## 🔄 Fluxo de Cache em 3 Camadas

### **Exemplo: Request para `/cep/01310100`**

```
┌──────────────────────────────────────────────────┐
│ 1. REQUEST CHEGA                                  │
│    GET /cep/01310100 + API Key                   │
└──────────────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────────────┐
│ ⚡ CAMADA 1: REDIS (L1 - Hot Cache)              │
│    • Verifica: key "cep:01310100" existe?        │
│    • Se SIM → Retorna <1ms ⚡                     │
│    • Se NÃO → Vai para Camada 2                  │
└──────────────────────────────────────────────────┘
                    ↓ (cache miss)
┌──────────────────────────────────────────────────┐
│ 🗄️ CAMADA 2: MONGODB (L2 - Cold Cache)          │
│    • Busca na collection "cep_cache"             │
│    • Se ENCONTROU e ainda válido (< 7 dias):    │
│      1. ✅ Retorna ao cliente (~10ms)            │
│      2. ✅ PROMOVE para Redis (write-back)       │
│    • Se NÃO encontrou → Vai para Camada 3       │
└──────────────────────────────────────────────────┘
                    ↓ (cache miss)
┌──────────────────────────────────────────────────┐
│ 🌐 CAMADA 3: VIACEP (API Externa)               │
│    • Busca CEP na API ViaCEP (~100-200ms)       │
│    • Quando retorna:                             │
│      1. ✅ Salva no Redis (L1, TTL 24h)          │
│      2. ✅ Salva no MongoDB (L2, TTL 7 dias)     │
│      3. ✅ Retorna ao cliente                    │
└──────────────────────────────────────────────────┘
```

---

## 💾 Estratégias de Cache

### **Write-Through (Usado)**
Quando busca da API externa, salva em **L1 E L2** simultaneamente:
```go
// Buscar de ViaCEP
response := fetchViaCEP(cep)

// Salvar no Redis (L1)
redis.Set("cep:"+cep, response, 24*time.Hour)

// Salvar no MongoDB (L2)
mongo.UpdateOne(cep, response, upsert=true)

// Retornar
return response
```

### **Cache Promotion (Usado)**
Quando encontra no MongoDB mas não no Redis:
```go
// Buscar no MongoDB
cached := mongo.FindOne("cep", cep)

if found && stillValid {
    // Promover para Redis (próxima request será <1ms)
    redis.Set("cep:"+cep, cached, 24*time.Hour)
    return cached
}
```

---

## 🔧 Graceful Degradation

### **O que acontece se Redis falhar?**

```go
// cmd/api/main.go
var redisClient interface{}
redisURL := os.Getenv("REDIS_URL")

if redisURL != "" {
    client, err := cache.NewRedisClient(redisURL, "", 0, log)
    if err != nil {
        // ⚠️ Redis falhou → Continua sem ele
        log.Warn().Msg("Redis não disponível, usando MongoDB")
        redisClient = nil
    } else {
        // ✅ Redis conectado
        log.Info().Msg("Redis conectado!")
        redisClient = client
    }
} else {
    // ⚠️ REDIS_URL não configurado
    log.Warn().Msg("REDIS_URL não configurado")
    redisClient = nil
}

// Passa redisClient (pode ser nil) para os handlers
router := nethttp.NewRouter(log, m, redisClient, ...)
```

### **Handlers lidam com Redis nil:**

```go
// internal/http/handlers/cep.go
func (h *CEPHandler) GetCEP(c *gin.Context) {
    // ⚡ Tentar Redis (se disponível)
    if h.redis != nil {
        if redisClient, ok := h.redis.(*cache.RedisClient); ok {
            cached, err := redisClient.Get(ctx, "cep:"+cep)
            if err == nil {
                return cached // ⚡ <1ms
            }
        }
    }
    
    // 🗄️ Fallback: MongoDB
    cached := mongo.FindOne("cep", cep)
    if found {
        return cached // ~10ms
    }
    
    // 🌐 Fallback: API Externa
    response := fetchViaCEP(cep)
    return response // ~150ms
}
```

**Resultado:**
- **Com Redis**: <5ms (L1)
- **Sem Redis**: ~10-50ms (L2 - MongoDB)
- **Sem Cache**: ~150ms (L3 - API Externa)

**NUNCA QUEBRA!** ✅

---

## 📊 Comparação: Local vs Produção

| Aspecto | Local (Docker Compose) | Produção (Railway) |
|---------|------------------------|-------------------|
| **Redis** | Container local | Instância dedicada Railway |
| **Conexão** | `redis://redis:6379` | `redis://...@redis.railway.internal:6379` |
| **Variável** | Hardcoded no docker-compose | `REDIS_URL` do Railway |
| **Memória** | 256MB (limitado) | 256MB (Railway) |
| **Failover** | Graceful degradation | Graceful degradation |
| **Performance** | <1ms (local) | <5ms (network latency) |

---

## 🎯 Por Que 3 Camadas?

### **Por que não só Redis?**
- Redis é **volátil** (dados podem ser evitados se memória cheia)
- MongoDB é **persistente** (backup de longo prazo)

### **Por que não só MongoDB?**
- MongoDB é **lento** (~10-50ms)
- Redis é **ultra-rápido** (<1ms)

### **Por que não só API Externa?**
- API Externa é **muito lenta** (~100-300ms)
- API Externa pode ter **rate limit** ou **downtime**

---

## ✅ Checklist de Configuração

### **Local (Dev)**
- [x] `docker-compose.yml` tem serviço `redis`
- [x] Backend aponta para `redis://redis:6379`
- [x] `docker-compose up` inicia Redis e Backend

### **Produção (Railway)**
- [x] Criar serviço Redis dedicado no Railway
- [x] Copiar `REDIS_URL` do serviço Redis
- [x] Adicionar `REDIS_URL` no serviço `retech-core`
- [x] Deploy do backend (rebuild automático)
- [ ] Verificar logs: `"⚡ Redis conectado - cache ultra-rápido habilitado!"`

---

## 🎉 Conclusão

**Você entendeu 100% correto!**

✅ **Local**: Docker Compose gerencia Redis na porta 6379  
✅ **Produção**: Railway gerencia Redis via `REDIS_URL`  
✅ **Graceful**: Se falhar, continua funcionando (só mais lento)  
✅ **Performance**: <5ms com Redis vs ~160ms sem Redis  

🚀 **Resultado Final**: Sistema robusto, rápido e escalável!

