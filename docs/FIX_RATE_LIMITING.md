# 🔧 FIX: Rate Limiting - Retech Core API

**Data**: 2025-10-22  
**Versão**: 1.3.1  
**Status**: ✅ CORRIGIDO

---

## 🐛 **BUGS ENCONTRADOS**

### 1. Erro de Sintaxe (Linha 80)
**Problema**: Faltava `{` após `if` na linha 80  
**Impacto**: Código não executava o bloco condicional corretamente

```go
// ❌ ANTES (BUGADO)
if rateLimit.Count >= config.RequestsPerDay
    c.Header("X-RateLimit-Limit", ...)
```

```go
// ✅ DEPOIS (CORRIGIDO)
if rateLimitDaily.Count >= config.RequestsPerDay {
    c.Header("X-RateLimit-Limit-Day", ...)
    // ... resto do bloco
}
```

---

### 2. Verificação DEPOIS do Incremento
**Problema**: O código incrementava o contador **ANTES** de verificar o limite

```go
// ❌ ANTES (ERRADO)
// Verificar limite (linha 80)
if rateLimit.Count >= config.RequestsPerDay { ... }

// Incrementar contador (linha 96) ← ACONTECIA ANTES DA VERIFICAÇÃO!
rateLimit.Count++
```

**Resultado**: Permitia N+1 requests em vez de N

```go
// ✅ DEPOIS (CORRETO)
// 1. Buscar registro atual
// 2. VERIFICAR limite ANTES de incrementar
if rateLimitDaily.Count >= config.RequestsPerDay {
    // Bloquear com 429
    c.Abort()
    return
}

// 3. INCREMENTAR apenas se passou
rateLimitDaily.Count++
```

---

### 3. Limite por Minuto NÃO Implementado
**Problema**: Campo `RequestsPerMinute` existia mas nunca era verificado

```go
// ❌ ANTES
// Só verificava RequestsPerDay
if rateLimit.Count >= config.RequestsPerDay { ... }
```

```go
// ✅ DEPOIS
// Verifica AMBOS: diário E por minuto
if rateLimitDaily.Count >= config.RequestsPerDay { ... }
if rateLimitMinute.Count >= config.RequestsPerMinute { ... }
```

---

## ✅ **CORREÇÕES IMPLEMENTADAS**

### 1. Verificação Antes de Incrementar
```go
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// VERIFICAR LIMITE DIÁRIO
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
collDaily := rl.db.Collection("rate_limits")
var rateLimitDaily domain.RateLimit

err := collDaily.FindOne(ctx, bson.M{
    "apiKey": apiKey,
    "date":   today,
}).Decode(&rateLimitDaily)

// ✅ VERIFICAR **ANTES** DE INCREMENTAR!
if rateLimitDaily.Count >= config.RequestsPerDay {
    // Retornar 429
    c.JSON(http.StatusTooManyRequests, ...)
    c.Abort()
    return
}

// Só incrementa se passou na verificação
rateLimitDaily.Count++
```

---

### 2. Implementação de Limite por Minuto
```go
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// VERIFICAR LIMITE POR MINUTO
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
collMinute := rl.db.Collection("rate_limits_minute")
currentMinute := now.Format("2006-01-02 15:04") // YYYY-MM-DD HH:MM

var rateLimitMinute domain.RateLimit
err = collMinute.FindOne(ctx, bson.M{
    "apiKey": apiKey,
    "date":   currentMinute, // ← Chave por minuto!
}).Decode(&rateLimitMinute)

// ✅ VERIFICAR LIMITE POR MINUTO
if rateLimitMinute.Count >= config.RequestsPerMinute {
    // Retornar 429 (limite por minuto)
    c.JSON(http.StatusTooManyRequests, ...)
    c.Abort()
    return
}

// Incrementar contador por minuto
rateLimitMinute.Count++
```

---

### 3. Headers Separados para Diário e Por Minuto
```go
// Headers diários
c.Header("X-RateLimit-Limit-Day", fmt.Sprintf("%d", config.RequestsPerDay))
c.Header("X-RateLimit-Remaining-Day", fmt.Sprintf("%d", remainingDay))
c.Header("X-RateLimit-Reset-Day", getNextDayTimestamp())

// Headers por minuto
c.Header("X-RateLimit-Limit-Minute", fmt.Sprintf("%d", config.RequestsPerMinute))
c.Header("X-RateLimit-Remaining-Minute", fmt.Sprintf("%d", remainingMinute))
c.Header("X-RateLimit-Reset-Minute", getNextMinuteTimestamp())
```

---

### 4. Logs Detalhados
```go
fmt.Printf("🔍 Rate Limit Config para tenant %s: %d/dia, %d/min\n", 
    tenantID, config.RequestsPerDay, config.RequestsPerMinute)

// Se bloqueado:
fmt.Printf("🚫 Rate Limit DIÁRIO excedido: %d >= %d\n", 
    rateLimitDaily.Count, config.RequestsPerDay)

// Se permitido:
fmt.Printf("✅ Request permitida. Restante: %d/dia, %d/min\n", 
    remainingDay, remainingMinute)
```

---

## 📊 **EXEMPLO DE FUNCIONAMENTO**

### Cenário: Tenant com limite 5/dia e 2/min

```
Request #1:  ✅ 200 OK  (0→1/dia, 0→1/min)
Request #2:  ✅ 200 OK  (1→2/dia, 1→2/min)
Request #3:  🚫 429    (2/min atingido!)
             └─ "Limite de 2 requests por minuto excedido"

... espera 1 minuto ...

Request #4:  ✅ 200 OK  (2→3/dia, 0→1/min novo minuto)
Request #5:  ✅ 200 OK  (3→4/dia, 1→2/min)
Request #6:  ✅ 200 OK  (4→5/dia, 2→3/min) ← mas só 2/min!
             └─ BLOQUEADO primeiro pelo limite/min (2/min)
             
Request #6 (retry após 1 min):  
             ✅ 200 OK  (4→5/dia, 0→1/min)
             
Request #7:  🚫 429    (5/dia atingido!)
             └─ "Limite de 5 requests por dia excedido"
```

---

## 🧪 **COMO TESTAR**

### Teste Manual Rápido

```bash
# 1. Configurar tenant com limite baixo
# Via Admin UI: 5/dia, 2/min

# 2. Fazer 3 requests rápidas (< 1 minuto)
for i in {1..3}; do
  curl -i -H "X-API-Key: rtc_..." \
    https://api-core.theretech.com.br/geo/estados
  echo "Request $i"
done

# Esperado:
# Request 1: 200 OK
# Request 2: 200 OK
# Request 3: 429 Rate Limit Exceeded (Per Minute)
```

### Teste Automatizado Completo

```bash
# Usar script de teste
./scripts/test-rate-limit.sh

# Ver resultados com gráficos
node scripts/generate-rate-limit-chart.js

# Abrir relatório HTML
open rate-limit-test-report.html
```

---

## 📋 **COLLECTIONS MONGODB**

### `rate_limits` (Diário)
```javascript
{
  "_id": ObjectId("..."),
  "apiKey": "rtc_abc123...",
  "date": "2025-10-22",           // YYYY-MM-DD
  "count": 42,
  "lastReset": ISODate("..."),
  "updatedAt": ISODate("...")
}
```

### `rate_limits_minute` (Por Minuto)
```javascript
{
  "_id": ObjectId("..."),
  "apiKey": "rtc_abc123...",
  "date": "2025-10-22 14:35",     // YYYY-MM-DD HH:MM
  "count": 2,
  "lastReset": ISODate("..."),
  "updatedAt": ISODate("...")
}
```

**TTL Index**: Registros de minuto expiram automaticamente após 2 minutos

---

## 🔍 **VERIFICAR NO BANCO**

```javascript
// Ver rate limits diários
db.rate_limits.find({}).sort({updatedAt: -1}).limit(10).pretty()

// Ver rate limits por minuto
db.rate_limits_minute.find({}).sort({updatedAt: -1}).limit(10).pretty()

// Ver por API key específica
db.rate_limits.find({apiKey: "rtc_..."}).pretty()

// Limpar contadores (dev/test)
db.rate_limits.deleteMany({})
db.rate_limits_minute.deleteMany({})
```

---

## 🎯 **COMPORTAMENTO ESPERADO**

### ✅ Funcionando Corretamente
- Primeiras N requests retornam `200 OK`
- Request N+1 retorna `429 Too Many Requests`
- Headers `X-RateLimit-*` corretos
- Limites diário E por minuto respeitados
- Bloqueio acontece exatamente após o limite

### ❌ Ainda Com Problemas (Se...)
- Requests acima do limite são permitidas
- Nenhum `429` é retornado
- Headers ausentes ou incorretos
- Limite por minuto não funciona
- Bloqueio acontece na request errada

---

## 📝 **LOGS ESPERADOS**

### Request Permitida
```
🔍 Rate Limit Config para tenant 67183abc: 5/dia, 2/min
✅ Request permitida. Restante: 4/dia, 1/min
```

### Request Bloqueada (Diário)
```
🔍 Rate Limit Config para tenant 67183abc: 5/dia, 2/min
🚫 Rate Limit DIÁRIO excedido: 5 >= 5
```

### Request Bloqueada (Por Minuto)
```
🔍 Rate Limit Config para tenant 67183abc: 5/dia, 2/min
🚫 Rate Limit POR MINUTO excedido: 2 >= 2
```

---

## 🚀 **DEPLOY**

### Produção
```bash
# 1. Push para GitHub
git add internal/middleware/rate_limiter.go
git commit -m "fix: rate limiting - verificação antes de incrementar + limite por minuto"
git push origin main

# 2. Railway faz deploy automático

# 3. Verificar logs
railway logs --tail

# 4. Testar em produção
./scripts/test-rate-limit.sh
```

### Local
```bash
# 1. Rebuild
go build -o bin/retech-core ./cmd/api

# 2. Rodar
./bin/retech-core

# 3. Testar
./scripts/test-rate-limit.sh
```

---

## ✅ **CHECKLIST DE VERIFICAÇÃO**

Após o fix, verificar:

- [ ] ✅ Código compila sem erros
- [ ] ✅ Limite diário funciona (5 requests = 429 na 6ª)
- [ ] ✅ Limite por minuto funciona (2 requests/min = 429 na 3ª)
- [ ] ✅ Headers `X-RateLimit-*-Day` presentes
- [ ] ✅ Headers `X-RateLimit-*-Minute` presentes
- [ ] ✅ Logs mostram configuração e bloqueios
- [ ] ✅ MongoDB tem registros em ambas collections
- [ ] ✅ Teste automatizado passa 100%
- [ ] ✅ Produção funciona corretamente

---

## 📚 **REFERÊNCIAS**

- `PENDENCIAS_CRITICAS.md` - Lista completa de bugs
- `scripts/test-rate-limit.sh` - Script de teste
- `scripts/generate-rate-limit-chart.js` - Visualização
- `ROADMAP_UNIFIED.md` - Fase 4 concluída

---

**Arquivo modificado**: `internal/middleware/rate_limiter.go`  
**Linhas modificadas**: 31-244 (reescrita completa do middleware)  
**Collections criadas**: `rate_limits_minute` (nova)  
**Testes**: `scripts/test-rate-limit.sh`  

---

**Mantido por**: The Retech Team  
**Última atualização**: 2025-10-22  
**Status**: ✅ PRONTO PARA PRODUÇÃO

