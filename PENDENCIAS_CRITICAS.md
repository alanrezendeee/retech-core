# 🚨 PENDÊNCIAS CRÍTICAS - Retech Core API

**Data**: 2025-10-22  
**Prioridade**: 🔴 URGENTE  

---

## 🐛 **BUGS CRÍTICOS**

### 1. ❌ **Rate Limiting NÃO FUNCIONA** 
**Arquivo**: `internal/middleware/rate_limiter.go`  
**Linha**: 80  
**Status**: 🔴 CRÍTICO

**Problema**:
```go
// LINHA 80 - ERRO DE SINTAXE!
if rateLimit.Count >= config.RequestsPerDay
    c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerDay))
```

Falta `{` após a condição! O código NÃO compila corretamente.

**Problema adicional**:
- Verifica limite DEPOIS de incrementar (linha 96)
- `RequestsPerMinute` NÃO é implementado (só diário)
- Não há reset automático de contadores

**Fix necessário**:
```go
// Verificar limite ANTES de incrementar
if rateLimit.Count >= config.RequestsPerDay {
    c.JSON(http.StatusTooManyRequests, gin.H{
        "type":   "https://retech-core/errors/rate-limit-exceeded",
        "title":  "Rate Limit Exceeded",
        "status": http.StatusTooManyRequests,
        "detail": fmt.Sprintf("Limite de %d requests por dia excedido", config.RequestsPerDay),
    })
    c.Abort()
    return
}

// Verificar limite por minuto também
// TODO: Implementar
```

**Testes realizados pelo usuário**:
- ✅ Configurou tenant para 10 req/dia e 1 req/minuto
- ✅ Fez MAIS requisições que o permitido
- ❌ Sistema NÃO bloqueou

---

## 🔧 **FEATURES INCOMPLETAS**

### 2. ⚠️ **Dashboard do Desenvolvedor - Dados Mocados**
**Arquivo**: `retech-core-admin/app/painel/dashboard/page.tsx`  
**Status**: 🟡 PARCIAL

**Problema**:
- Cards do dashboard mostram dados MOCADOS
- Não consome `/me/usage` ou `/me/stats`
- Gráficos não exibem dados reais

**Endpoints necessários**:
```typescript
GET /me/stats        - KPIs do tenant (requests hoje/mês)
GET /me/usage        - Histórico de uso
GET /me/logs         - Últimas requisições (limitado)
```

**O que implementar**:
- [ ] Endpoint `/me/stats` no backend
- [ ] Endpoint `/me/usage` no backend
- [ ] Consumir dados reais no frontend
- [ ] Gráficos com recharts (dados reais)

---

### 3. ⚠️ **Uso da API - Developer Portal**
**Arquivo**: `retech-core-admin/app/painel/uso-api/page.tsx`  
**Status**: 🟡 PARCIAL

**Problema**:
- Página existe mas não mostra dados reais
- Falta integração com `/me/usage`
- Falta histórico de requisições

**O que implementar**:
- [ ] Gráfico de uso (últimos 7/30 dias)
- [ ] Endpoints mais usados
- [ ] Horários de pico
- [ ] Alertas quando próximo do limite
- [ ] Exportar relatórios

---

### 4. ❌ **Documentação Desatualizada**
**Arquivos**:
- `internal/docs/openapi.yaml`
- `internal/docs/redoc.html`
- Links na landing page

**Status**: 🔴 DESATUALIZADO

**Problemas**:
- OpenAPI não reflete endpoints atuais
- Falta documentação de:
  - `/admin/*` (stats, tenants, apikeys)
  - `/me/*` (stats, usage, logs)
  - `/auth/*` (login, register, refresh)
- Redoc não está acessível/atualizado
- Landing page pode ter link quebrado para docs

**O que fazer**:
- [ ] Atualizar `openapi.yaml` com TODOS endpoints
- [ ] Gerar Redoc atualizado
- [ ] Verificar link na landing page
- [ ] Criar página `/painel/docs` com getting started

---

### 5. ⚠️ **Rate Limiting por Minuto NÃO Implementado**
**Arquivo**: `internal/middleware/rate_limiter.go`  
**Status**: 🔴 NÃO IMPLEMENTADO

**Problema**:
- Configuração existe: `RequestsPerMinute`
- Admin pode configurar no UI
- Backend NÃO verifica este limite
- Apenas limite diário é checado

**O que implementar**:
- [ ] Adicionar collection `rate_limits_minute`
- [ ] Verificar timestamp da última request
- [ ] Calcular requests no minuto atual
- [ ] Retornar 429 se exceder
- [ ] Headers `X-RateLimit-Minute-Remaining`

---

### 6. ⚠️ **Reset Automático de Rate Limits**
**Status**: 🔴 NÃO IMPLEMENTADO

**Problema**:
- Contadores diários não resetam automaticamente
- Depende de nova request para criar novo registro
- Pode acumular dados antigos no banco

**O que implementar**:
- [ ] Cron job diário (meia-noite) para limpar
- [ ] Ou: lógica no middleware para verificar data
- [ ] Índice TTL no MongoDB (auto-delete após 30 dias)

---

## 📚 **DOCUMENTAÇÃO FALTANTE**

### 7. ⚠️ **Developer Portal - Getting Started**
**Arquivo**: `retech-core-admin/app/painel/docs/page.tsx` (NÃO EXISTE)  
**Status**: ❌ NÃO CRIADO

**O que criar**:
- [ ] Página `/painel/docs`
- [ ] Getting Started (primeiros passos)
- [ ] Exemplos de código (cURL, JS, Python, Go)
- [ ] Lista de endpoints disponíveis
- [ ] Limites do plano atual
- [ ] FAQ
- [ ] Link para Redoc completo

---

### 8. ⚠️ **Links para Documentação**
**Status**: 🟡 PARCIAL

**Verificar e adicionar**:
- [ ] Landing page → link para docs
- [ ] Dashboard admin → link para docs
- [ ] Dashboard painel → link para docs
- [ ] Email de boas-vindas → link para docs
- [ ] README.md → link para docs públicas

---

## 🎯 **ENDPOINTS FALTANTES (Backend)**

### `/me/*` - Tenant Self-Service

- [ ] `GET /me/stats` - KPIs do tenant
  ```json
  {
    "requestsToday": 42,
    "requestsThisMonth": 1205,
    "limitDaily": 1000,
    "limitMonthly": 30000,
    "percentUsed": 40.2
  }
  ```

- [ ] `GET /me/usage` - Histórico de uso
  ```json
  {
    "daily": [...],  // Últimos 30 dias
    "byEndpoint": [...],
    "byHour": [...]
  }
  ```

- [ ] `GET /me/logs` - Últimas requisições (limitado a 100)
  ```json
  {
    "logs": [
      {
        "timestamp": "2025-10-22T10:30:00Z",
        "endpoint": "/geo/estados",
        "method": "GET",
        "statusCode": 200,
        "responseTime": 45
      }
    ]
  }
  ```

---

## 🔄 **MELHORIAS NECESSÁRIAS**

### 9. ⚠️ **Activity Logs - Deep Links**
**Status**: 🟡 PARCIAL

**Problema atual**:
- Activity logs existem
- Deep links podem não funcionar 100%
- Ex: clicar em "Tenant criado" deveria ir para `/admin/tenants/:id`

**Verificar**:
- [ ] Deep links funcionam
- [ ] Links abrem em nova aba ou mesma?
- [ ] Loading state ao clicar
- [ ] Erro 404 se recurso deletado

---

### 10. ⚠️ **Analytics Avançado**
**Status**: 🔴 NÃO IMPLEMENTADO

**O que adicionar**:
- [ ] Gráficos de comparação de períodos
- [ ] Heatmap de uso por horário
- [ ] Geolocalização de requests (por IP)
- [ ] Top endpoints por tenant
- [ ] Exportar relatórios CSV/PDF

---

## 📋 **CHECKLIST DE PRIORIDADE**

### 🔴 **URGENTE (Esta Semana)**
- [ ] **1. Corrigir Rate Limiting** (bug crítico)
  - [ ] Fix sintaxe linha 80
  - [ ] Mover verificação ANTES de incrementar
  - [ ] Implementar limite por minuto
  - [ ] Testar com 10 req/dia e 1 req/min
  
- [ ] **2. Implementar `/me/stats`** (dashboard dev)
  - [ ] Criar handler
  - [ ] Agregações MongoDB
  - [ ] Testar no frontend

- [ ] **3. Atualizar Documentação**
  - [ ] `openapi.yaml` completo
  - [ ] Gerar Redoc atualizado
  - [ ] Verificar links

### 🟡 **IMPORTANTE (Próxima Semana)**
- [ ] **4. Dashboard Developer com dados reais**
  - [ ] `/me/usage` backend
  - [ ] Gráficos recharts
  - [ ] Integrar frontend

- [ ] **5. Página de Docs no Painel**
  - [ ] Criar `/painel/docs`
  - [ ] Getting started
  - [ ] Exemplos de código

### 🟢 **DESEJÁVEL (Próximo Sprint)**
- [ ] **6. Analytics Avançado**
- [ ] **7. Reset automático de rate limits**
- [ ] **8. Geolocalização de requests**

---

## 🧪 **COMO TESTAR RATE LIMITING**

### Teste Manual (depois do fix)

1. **Configurar tenant**:
```bash
# Via Admin UI
- Ir em /admin/tenants
- Editar tenant
- Rate Limit Personalizado: Ativo
- Requests por Dia: 5
- Requests por Minuto: 2
```

2. **Fazer requisições**:
```bash
# Com API key do tenant
curl -H "X-API-Key: rtc_..." https://api-core.theretech.com.br/geo/estados

# Repetir 6x rapidamente
# Na 6ª deve retornar 429
```

3. **Verificar headers**:
```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1729641600
```

4. **Verificar MongoDB**:
```javascript
db.rate_limits.find({}).pretty()
```

---

## 📊 **ROADMAP UNIFICADO**

Vou criar `ROADMAP_UNIFIED.md` com:
- ✅ Tudo do ROADMAP.md antigo
- ✅ Novas 31 APIs do ROADMAP_V2.md
- ✅ Estas pendências críticas
- ✅ Timeline realista

---

## 🤝 **PRÓXIMAS AÇÕES**

1. ✅ Criar este documento (PENDENCIAS_CRITICAS.md)
2. 🔄 Corrigir Rate Limiting (PR separado)
3. 🔄 Implementar `/me/stats` (PR separado)
4. 🔄 Atualizar documentação OpenAPI
5. 🔄 Criar ROADMAP_UNIFIED.md

---

**Mantido por**: The Retech Team  
**Última atualização**: 2025-10-22  
**Revisão**: Sempre que completar uma pendência

