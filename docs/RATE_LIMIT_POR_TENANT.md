# 🛡️ Rate Limiting Personalizado por Tenant

## 📋 Visão Geral

Agora cada tenant pode ter seu próprio limite de requisições, sobrescrevendo o limite padrão do sistema.

---

## 🎯 Como Funciona

### **1. Limite Padrão (DEFAULT):**
```
Definido em: /admin/settings
Aplicado quando: tenant.rateLimit == null
Valores padrão: 1.000 req/dia, 60 req/minuto
```

### **2. Limite Personalizado (POR TENANT):**
```
Definido em: /admin/tenants (criar/editar)
Aplicado quando: tenant.rateLimit != null
Valores: Definidos pelo admin para cada tenant
```

---

## 🎨 Interface de Configuração

### **Localização:**
```
/admin/tenants
→ Criar Novo Tenant / Editar Tenant
→ Seção "Rate Limiting"
```

### **Visual:**

```
┌────────────────────────────────────────────┐
│  [Switch] Limite Personalizado             │
│  ○ Usando limite padrão do sistema         │
│     (1.000/dia)                             │
└────────────────────────────────────────────┘

Quando ativado:

┌────────────────────────────────────────────┐
│  [Switch] Limite Personalizado             │
│  ● Usando limite customizado               │
│                                             │
│  ┌────────────────────────────────────┐   │
│  │ ℹ️ Configure os limites específicos │   │
│  │   para este tenant                  │   │
│  │                                      │   │
│  │ Requests por Dia: [____]            │   │
│  │ Máximo de requisições por dia       │   │
│  │                                      │   │
│  │ Requests por Minuto: [____]         │   │
│  │ Máximo de requisições por minuto    │   │
│  └────────────────────────────────────┘   │
└────────────────────────────────────────────┘
```

---

## 🔄 Fluxo de Decisão do Middleware

```
┌─────────────────────────────────────────┐
│  Request com API Key                    │
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│  Identificar Tenant (via API Key)       │
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│  Buscar Tenant no MongoDB               │
└─────────────────────────────────────────┘
              │
              ▼
        ┌─────────┐
        │ tenant? │
        └─────────┘
         /       \
       SIM       NÃO
        │         │
        ▼         ▼
  ┌──────────┐   ┌──────────────────┐
  │rateLimit?│   │Usa DEFAULT       │
  └──────────┘   │do sistema        │
    /     \      └──────────────────┘
  SIM     NÃO
   │       │
   ▼       ▼
┌─────┐  ┌────────┐
│Usa  │  │Usa     │
│CUSTOM│  │DEFAULT │
└─────┘  └────────┘
```

---

## 📊 Exemplos de Uso

### **Exemplo 1: Tenant Free (sem personalização)**

```json
{
  "tenantId": "tenant-free-001",
  "name": "Startup ABC",
  "email": "contato@startup.com",
  "active": true,
  "rateLimit": null  // ← Usa DEFAULT
}
```

**Resultado:**
- ✅ 1.000 requisições/dia
- ✅ 60 requisições/minuto

---

### **Exemplo 2: Tenant Premium (personalizado)**

```json
{
  "tenantId": "tenant-premium-001",
  "name": "Empresa XYZ",
  "email": "contato@empresa.com",
  "active": true,
  "rateLimit": {
    "requestsPerDay": 100000,
    "requestsPerMinute": 1000
  }
}
```

**Resultado:**
- ✅ 100.000 requisições/dia
- ✅ 1.000 requisições/minuto

---

### **Exemplo 3: Tenant Enterprise (limite alto)**

```json
{
  "tenantId": "tenant-enterprise-001",
  "name": "Corporação ABC",
  "email": "contato@corp.com",
  "active": true,
  "rateLimit": {
    "requestsPerDay": 1000000,
    "requestsPerMinute": 5000
  }
}
```

**Resultado:**
- ✅ 1.000.000 requisições/dia
- ✅ 5.000 requisições/minuto

---

## 🎯 Casos de Uso

### **1. Plano Free:**
```
Switch: Desativado
→ Usa limite padrão (1k/dia)
```

### **2. Plano Pro:**
```
Switch: Ativado
requestsPerDay: 10000
requestsPerMinute: 200
→ 10x mais que o free
```

### **3. Plano Enterprise:**
```
Switch: Ativado
requestsPerDay: 100000
requestsPerMinute: 1000
→ 100x mais que o free
```

### **4. Tenant Beta/Teste:**
```
Switch: Ativado
requestsPerDay: 500
requestsPerMinute: 10
→ Limite mais baixo para testes
```

---

## 🧪 Como Testar

### **1. Criar Tenant com Limite Personalizado:**

```bash
# Login como admin
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alanrezendeee@gmail.com","password":"admin123456"}' \
  | jq -r '.accessToken')

# Criar tenant com rate limit personalizado
curl -s -X POST http://localhost:8080/admin/tenants \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Empresa Premium",
    "email": "premium@empresa.com",
    "company": "Premium LTDA",
    "active": true,
    "rateLimit": {
      "requestsPerDay": 5000,
      "requestsPerMinute": 100
    }
  }' \
  | jq '.'
```

---

### **2. Verificar Tenant Criado:**

```bash
curl -s http://localhost:8080/admin/tenants \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.tenants[] | select(.email == "premium@empresa.com")'
```

**Resposta esperada:**
```json
{
  "tenantId": "tenant-...",
  "name": "Empresa Premium",
  "email": "premium@empresa.com",
  "active": true,
  "rateLimit": {
    "requestsPerDay": 5000,
    "requestsPerMinute": 100
  }
}
```

---

### **3. Testar Rate Limiting:**

1. Criar API Key para o tenant
2. Fazer requisições até atingir o limite (5.000)
3. Verificar header `X-RateLimit-Remaining`
4. Ao exceder, receber 429

```bash
# Fazer requisição com API key do tenant
curl -i http://localhost:8080/geo/ufs \
  -H "X-API-Key: sua-api-key-aqui"

# Verificar headers:
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4999
X-RateLimit-Reset: 1729641600
```

---

### **4. Atualizar Rate Limit de Tenant Existente:**

```bash
# Alterar limite
curl -s -X PUT http://localhost:8080/admin/tenants/tenant-... \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rateLimit": {
      "requestsPerDay": 10000,
      "requestsPerMinute": 200
    }
  }' \
  | jq '.'
```

---

### **5. Remover Limite Personalizado (voltar ao DEFAULT):**

```bash
# Setar rateLimit como null
curl -s -X PUT http://localhost:8080/admin/tenants/tenant-... \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rateLimit": null
  }' \
  | jq '.'
```

---

## 🎨 Interface Visual

### **No Drawer de Edição:**

```
╔══════════════════════════════════════════════╗
║  Editar Tenant                               ║
╠══════════════════════════════════════════════╣
║                                              ║
║  Nome: [Empresa Premium___________________]  ║
║  Email: [premium@empresa.com_____________]   ║
║  Empresa: [Premium LTDA_________________]   ║
║                                              ║
║  ─────────── Rate Limiting ────────────      ║
║                                              ║
║  ┌──────────────────────────────────────┐   ║
║  │ 🛡️ Limite Personalizado     [ON] ●  │   ║
║  │ Usando limite customizado            │   ║
║  └──────────────────────────────────────┘   ║
║                                              ║
║  ┌──────────────────────────────────────┐   ║
║  │ ℹ️ Configure os limites específicos  │   ║
║  │                                       │   ║
║  │ Requests por Dia                     │   ║
║  │ [5000__________________________]     │   ║
║  │ Máximo de requisições por dia        │   ║
║  │                                       │   ║
║  │ Requests por Minuto                  │   ║
║  │ [100___________________________]     │   ║
║  │ Máximo de requisições por minuto     │   ║
║  └──────────────────────────────────────┘   ║
║                                              ║
║  [Cancelar]  [Atualizar com gradiente]      ║
╚══════════════════════════════════════════════╝
```

---

## 📝 Validações

### **Frontend:**
- requestsPerDay: 1 a 1.000.000
- requestsPerMinute: 1 a 10.000
- Inputs numéricos com min/max

### **Backend:**
- Aceita `null` (remove limite personalizado)
- Aceita objeto com `requestsPerDay` e `requestsPerMinute`
- Valida tipos (int64)

---

## 🔍 Middleware - Lógica Interna

```go
// middleware/rate_limiter.go

func (rl *RateLimiter) getRateLimitConfig(tenantID string) RateLimitConfig {
    // 1. Buscar tenant
    tenant, err := rl.tenants.ByTenantID(ctx, tenantID)
    
    // 2. Se tenant tem rateLimit, usar ele
    if err == nil && tenant != nil && tenant.RateLimit != nil {
        return *tenant.RateLimit  // ← PERSONALIZADO
    }
    
    // 3. Senão, buscar DEFAULT do sistema
    settings, err := rl.settings.Get(ctx)
    if err == nil && settings != nil {
        return settings.DefaultRateLimit  // ← DEFAULT
    }
    
    // 4. Fallback hardcoded
    return RateLimitConfig{
        RequestsPerDay: 1000,
        RequestsPerMinute: 60,
    }
}
```

---

## ✅ Checklist de Implementação

- [x] Interface de switch no TenantDrawer
- [x] Inputs para requestsPerDay e requestsPerMinute
- [x] Estado local para customRateLimit
- [x] Salvar rateLimit no tenant (null ou objeto)
- [x] Backend aceita campo rateLimit no Tenant
- [x] Middleware verifica tenant.RateLimit
- [x] Fallback para SystemSettings.defaultRateLimit
- [x] Testes de criação e edição
- [x] Documentação completa

---

## 🎉 Resultado Final

Agora você pode:

1. ✅ **Definir limite padrão** em `/admin/settings`
2. ✅ **Personalizar por tenant** em `/admin/tenants`
3. ✅ **Ativar/desativar** limite personalizado com switch
4. ✅ **Configurar valores** específicos por tenant
5. ✅ **Middleware aplica** automaticamente o limite correto
6. ✅ **Headers retornam** o limite aplicado (`X-RateLimit-Limit`)

---

**Sistema completo de Rate Limiting implementado! 🚀**

