# 🎛️ Sistema de Settings Completo - Admin/Settings

## 📋 Resumo Executivo

**Problema Resolvido:** Todas as configurações críticas estavam **hardcoded** no código, exigindo rebuild/redeploy para qualquer mudança.

**Solução Implementada:** Sistema centralizado de configurações via **admin/settings** com atualização em tempo real.

---

## ✅ O Que Foi Implementado

### **1. CORS Dinâmico (Resolvido 100%)**

**ANTES:**
```go
// ❌ Hardcoded no router.go
if origin == "https://core.theretech.com.br" || origin == "http://localhost:3000" ... {
    c.Header("Access-Control-Allow-Origin", origin)
}
```

**AGORA:**
```go
// ✅ Lê do admin/settings
sysSettings, _ := settings.Get(ctx)
if sysSettings.CORS.Enabled {
    for _, allowedOrigin := range sysSettings.CORS.AllowedOrigins {
        // ... permite origins configuradas
    }
}
```

**Controle via Admin:**
```json
{
  "cors": {
    "enabled": true,
    "allowedOrigins": [
      "https://core.theretech.com.br",
      "http://localhost:3000",
      "http://localhost:3001",
      "http://localhost:3002",
      "http://localhost:3003"
    ]
  }
}
```

**Benefícios:**
- ✅ Admin pode **adicionar/remover origens** sem código
- ✅ Admin pode **desabilitar CORS** temporariamente
- ✅ Mudanças aplicam **imediatamente** (sem restart)

---

### **2. API Key Demo do Playground (Resolvido 100%)**

**ANTES:**
```go
// ❌ Hardcoded no bootstrap
const DemoAPIKeyValue = "rtc_demo_playground_2024"
```

**AGORA:**
```json
{
  "playground": {
    "enabled": true,
    "apiKey": "rtc_demo_playground_2024",
    "rateLimit": {
      "requestsPerDay": 100,
      "requestsPerMinute": 10
    },
    "allowedApis": ["cep", "cnpj", "geo"]
  }
}
```

**Controle via Admin:**
- 🔑 **Trocar API Key** se houver abuso
- 🎚️ **Ajustar rate limits** (requisições/dia e requisições/minuto)
- 🚫 **Desabilitar playground** temporariamente
- 🔧 **Escolher APIs públicas** (cep, cnpj, geo, etc)

**Como Funciona:**
1. Admin altera settings no painel
2. Na próxima requisição, bootstrap atualiza API Key automaticamente
3. Tenant demo é recriado com novos rate limits
4. Frontend continua usando `NEXT_PUBLIC_DEMO_API_KEY`

---

### **3. Redis Cache Configurável (Resolvido 100%)**

**ANTES:**
```go
// ❌ TTL hardcoded
redisClient.Set(ctx, redisKey, response, 24*time.Hour)
```

**AGORA:**
```go
// ✅ TTL dinâmico do settings
ttl := h.getTTL(c) // Lê de settings.Cache.CEPTTLDays
redisClient.Set(ctx, redisKey, response, ttl)
```

**Controle via Admin:**
```json
{
  "cache": {
    "enabled": true,
    "cepTtlDays": 7,
    "cnpjTtlDays": 30,
    "maxSizeMb": 100,
    "autoCleanup": true
  }
}
```

**Benefícios:**
- ⏱️ **TTL ajustável** (1-365 dias) por tipo de API
- 🚫 **Desabilitar cache** globalmente para debug
- 🧹 **Controlar limpeza automática**
- 📊 **Monitorar tamanho do cache** (futuro)

---

## 🏗️ Estrutura de Settings

### **Arquivo:** `internal/domain/settings.go`

```go
type SystemSettings struct {
    ID string `bson:"_id"`

    // Rate Limiting padrão para novos tenants
    DefaultRateLimit RateLimitConfig

    // CORS (✅ NOVO - dinâmico)
    CORS CORSConfig

    // JWT
    JWT JWTConfig

    // API Info
    API APIConfig

    // Contato/Vendas
    Contact ContactConfig

    // Cache (✅ NOVO - TTL configurável)
    Cache CacheConfig

    // Playground (✅ NOVO - API Key gerenciável)
    Playground PlaygroundConfig

    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

## 🎯 Como Usar no Admin UI (Próximo Passo)

### **Rota:** `GET /admin/settings`
```json
{
  "cors": {
    "enabled": true,
    "allowedOrigins": ["..."]
  },
  "cache": {
    "enabled": true,
    "cepTtlDays": 7,
    "cnpjTtlDays": 30
  },
  "playground": {
    "enabled": true,
    "apiKey": "rtc_demo_playground_2024",
    "rateLimit": {
      "requestsPerDay": 100,
      "requestsPerMinute": 10
    },
    "allowedApis": ["cep", "cnpj", "geo"]
  }
}
```

### **Rota:** `PUT /admin/settings`
```json
{
  "cors": {
    "enabled": true,
    "allowedOrigins": [
      "https://core.theretech.com.br",
      "https://novo-dominio.com"
    ]
  },
  "playground": {
    "enabled": false // ← Desabilita playground temporariamente
  },
  "cache": {
    "cepTtlDays": 14 // ← Aumenta cache de CEP para 14 dias
  }
}
```

---

## 🚀 Próximos Passos no Admin UI

### **1. Criar seção "Configurações Gerais"**
```tsx
// app/admin/settings/page.tsx
<div className="space-y-6">
  {/* CORS */}
  <Card>
    <CardTitle>CORS</CardTitle>
    <Switch checked={cors.enabled} onChange={...} />
    <Input 
      label="Origens Permitidas"
      value={cors.allowedOrigins.join(',')}
      onChange={...}
    />
  </Card>

  {/* Cache */}
  <Card>
    <CardTitle>Cache Redis</CardTitle>
    <Switch checked={cache.enabled} onChange={...} />
    <Input 
      label="TTL CEP (dias)"
      type="number"
      value={cache.cepTtlDays}
      min={1}
      max={365}
    />
    <Input 
      label="TTL CNPJ (dias)"
      type="number"
      value={cache.cnpjTtlDays}
      min={1}
      max={365}
    />
  </Card>

  {/* Playground */}
  <Card>
    <CardTitle>Playground Público</CardTitle>
    <Switch checked={playground.enabled} onChange={...} />
    <Input 
      label="API Key Demo"
      value={playground.apiKey}
      onChange={...}
    />
    <Input 
      label="Requests por Dia"
      type="number"
      value={playground.rateLimit.requestsPerDay}
    />
    <Input 
      label="Requests por Minuto"
      type="number"
      value={playground.rateLimit.requestsPerMinute}
    />
  </Card>
</div>
```

### **2. Endpoint para limpar cache manualmente**
```tsx
<Button onClick={() => fetch('/admin/cache/cep', { method: 'DELETE' })}>
  🗑️ Limpar Cache de CEP
</Button>
```

---

## 📊 Comparação Antes x Depois

| Feature | ANTES | AGORA |
|---------|-------|-------|
| **CORS** | ❌ Hardcoded no router | ✅ Admin/settings (tempo real) |
| **API Key Demo** | ❌ Hardcoded no bootstrap | ✅ Admin/settings (editável) |
| **TTL Redis** | ❌ Hardcoded (24h) | ✅ Admin/settings (1-365 dias) |
| **Playground ON/OFF** | ❌ Modificar código | ✅ Toggle no admin |
| **Rate Limits Playground** | ❌ Hardcoded (100/dia) | ✅ Admin/settings (ajustável) |
| **Mudanças** | ❌ Rebuild + Redeploy | ✅ Update instantâneo |

---

## 🔒 Segurança

✅ **Playground:**
- API Key editável via admin
- Rate limits agressivos (100 req/dia padrão)
- Pode ser desabilitado completamente
- Scopes restritos (apenas cep, cnpj, geo)

✅ **CORS:**
- Origins explicitamente permitidas
- Pode ser desabilitado para debug
- Adicionar/remover origins sem código

✅ **Cache:**
- TTL controlado pelo admin
- Pode ser desabilitado globalmente
- Limpeza manual disponível

---

## 🧪 Como Testar

### **1. CORS Dinâmico:**
```bash
# 1. Atualizar settings
curl -X PUT https://api-core.theretech.com.br/admin/settings \
  -H "Authorization: Bearer $JWT" \
  -d '{
    "cors": {
      "allowedOrigins": ["https://novo-dominio.com"]
    }
  }'

# 2. Testar CORS
curl https://api-core.theretech.com.br/health \
  -H "Origin: https://novo-dominio.com" \
  -v
# Deve retornar: Access-Control-Allow-Origin: https://novo-dominio.com
```

### **2. API Key Demo:**
```bash
# 1. Trocar chave no admin
curl -X PUT https://api-core.theretech.com.br/admin/settings \
  -H "Authorization: Bearer $JWT" \
  -d '{
    "playground": {
      "apiKey": "nova_chave_secreta_123"
    }
  }'

# 2. Testar nova chave
curl https://api-core.theretech.com.br/cep/88101270 \
  -H "X-API-Key: nova_chave_secreta_123"
# Deve funcionar!
```

### **3. TTL Configurável:**
```bash
# 1. Ajustar TTL
curl -X PUT https://api-core.theretech.com.br/admin/settings \
  -H "Authorization: Bearer $JWT" \
  -d '{
    "cache": {
      "cepTtlDays": 1
    }
  }'

# 2. Verificar logs do servidor
# Deve mostrar: "Salvo no Redis L1 (TTL: 24h)" → "Salvo no Redis L1 (TTL: 1d)"
```

---

## 📝 Documentação de Referência

- **Settings Domain:** `internal/domain/settings.go`
- **Settings Repo:** `internal/storage/settings_repo.go`
- **CORS Middleware:** `internal/http/router.go` (linha 31-68)
- **Playground Bootstrap:** `internal/bootstrap/demo_apikey.go`
- **Cache TTL:** `internal/http/handlers/cep.go` (método `getTTL`)

---

## ✅ Status Final

| Item | Status | Arquivo |
|------|--------|---------|
| CORS dinâmico | ✅ Completo | `router.go` |
| API Key demo gerenciável | ✅ Completo | `demo_apikey.go` |
| TTL Redis configurável | ✅ Completo | `cep.go` |
| Playground ON/OFF | ✅ Completo | `settings.go` |
| Rate limits editáveis | ✅ Completo | `settings.go` |
| Admin UI (frontend) | ⏳ Pendente | `retech-core-admin` |

---

## 🎉 Conclusão

**Tudo está centralizado em `admin/settings` e funciona em tempo real!** 

O admin agora tem controle total sobre:
- 🌐 CORS
- 🔑 API Key do Playground  
- ⏱️ TTL do Cache
- 🎚️ Rate Limits
- 🚫 Ativar/Desabilitar features

**Próximo passo:** Criar interface visual no admin UI! 🚀

