# 🔒 Sistema de Scopes (Permissões de API)

**Data:** 23 de Outubro de 2025  
**Status:** ✅ Implementado e Documentado

---

## 📋 RESUMO

Sistema completo de permissões para API Keys, permitindo controle granular de acesso por API.

---

## 🎯 SCOPES DISPONÍVEIS

| Scope | Descrição | Endpoints Permitidos | Status |
|-------|-----------|---------------------|--------|
| **`geo`** | Dados geográficos | `/geo/*` | ✅ Ativo |
| **`cep`** | Consulta de CEP | `/cep/*` | ✅ Ativo |
| **`cnpj`** | Consulta de CNPJ | `/cnpj/*` | ✅ Ativo |
| **`all`** | Acesso total | Todas as rotas | ✅ Ativo |
| `cpf` | Validação de CPF | `/cpf/*` | ⏳ Futuro |
| `fipe` | Tabela FIPE | `/fipe/*` | ⏳ Futuro |
| `moedas` | Cotações | `/moedas/*` | ⏳ Futuro |
| `bancos` | Lista de bancos | `/bancos/*` | ⏳ Futuro |

---

## 🔐 COMO FUNCIONA

### **1. Criação de API Key**

Ao criar uma API Key, o admin seleciona os scopes:

```json
POST /admin/apikeys
{
  "ownerId": "tenant-abc",
  "scopes": ["geo", "cep"],  // ← Escolhe quais APIs o tenant pode acessar
  "days": 90
}
```

**Resultado**: API Key com acesso **APENAS** a `/geo/*` e `/cep/*`.

---

### **2. Validação Automática**

Quando uma request chega, o middleware `RequireScope` valida:

```go
// Em router.go
cepGroup.Use(
    auth.AuthAPIKey(apikeys),          // 1. Valida se key é válida
    auth.RequireScope(apikeys, "cep"), // 2. Valida se tem scope 'cep' ou 'all'
    rateLimiter.Middleware(),
    usageLogger.Middleware(),
)
```

**Fluxo:**
1. Request: `GET /cep/01310100` com `X-API-Key: sk_...`
2. `AuthAPIKey`: Verifica se key existe e não está revogada ✅
3. `RequireScope`: Busca scopes da key no MongoDB
4. Se tem `cep` ou `all` → ✅ Permite
5. Se não tem → ❌ 403 Forbidden

---

### **3. Scope `all` (Super Permissão)**

```json
{
  "scopes": ["all"]
}
```

**Acesso a:**
- ✅ `/geo/*`
- ✅ `/cep/*`
- ✅ `/cnpj/*`
- ✅ Todas as futuras APIs

**Uso recomendado:**
- Tenants premium
- Integração completa
- Desenvolvimento/testes

---

## ⚠️ RESPOSTAS DE ERRO

### **403 Forbidden (Sem Permissão)**

```json
{
  "type": "https://retech-core/errors/forbidden",
  "title": "Insufficient Permissions",
  "status": 403,
  "detail": "API Key não tem permissão para acessar este recurso. Scope necessário: cnpj",
  "meta": {
    "requiredScope": "cnpj",
    "yourScopes": ["geo", "cep"]
  }
}
```

**Causa:** Tentou acessar `/cnpj/*` mas a key só tem `["geo", "cep"]`.

---

### **400 Bad Request (Scope Inválido)**

```json
{
  "type": "https://retech-core/errors/validation",
  "title": "Scopes Inválidos",
  "status": 400,
  "detail": "scopes: Scope 'xyz' não reconhecido"
}
```

**Causa:** Tentou criar API Key com scope que não existe.

---

### **400 Bad Request (Scope Não Disponível)**

```json
{
  "type": "https://retech-core/errors/validation",
  "title": "Scopes Inválidos",
  "status": 400,
  "detail": "scopes: Scope 'cpf' ainda não está disponível"
}
```

**Causa:** Scope existe no sistema mas API ainda não foi implementada.

---

## 🛠️ IMPLEMENTAÇÃO TÉCNICA

### **Middleware `RequireScope`**

**Arquivo:** `internal/auth/scope_middleware.go`

```go
func RequireScope(repo *storage.APIKeysRepo, requiredScope string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Buscar API Key no MongoDB
        // 2. Se tem scope 'all' → permite
        // 3. Se tem scope específico → permite
        // 4. Caso contrário → 403 Forbidden
    }
}
```

**Função helper:**
```go
func hasScope(scopes []string, scope string) bool {
    for _, s := range scopes {
        if s == scope || s == "all" {
            return true
        }
    }
    return false
}
```

---

### **Validação na Criação**

**Arquivo:** `internal/auth/scope_middleware.go`

```go
func ValidateAPIKeyScopes(scopes []string) error {
    validScopes := map[string]bool{
        "geo":   true,  // ✅ Ativo
        "cep":   true,  // ✅ Ativo
        "cnpj":  true,  // ✅ Ativo
        "all":   true,  // ✅ Ativo
        "cpf":   false, // ⏳ Futuro
        "fipe":  false, // ⏳ Futuro
    }
    
    // Valida cada scope
    for _, scope := range scopes {
        if implemented, exists := validScopes[scope]; !exists {
            return error("Scope não reconhecido")
        } else if !implemented {
            return error("Scope ainda não disponível")
        }
    }
    return nil
}
```

---

## 🎨 FRONTEND (Admin UI)

**Arquivo:** `components/apikeys/apikey-drawer.tsx`

```typescript
const availableScopes = [
  { value: 'geo', label: '🗺️ GEO - Dados geográficos' },
  { value: 'cep', label: '📮 CEP - Consulta de endereços' },
  { value: 'cnpj', label: '🏢 CNPJ - Dados de empresas' },
  { value: 'all', label: '⭐ ALL - Acesso total' },
];
```

**UI:**
- ☑️ Checkboxes para cada scope
- ✅ Validação: pelo menos 1 scope obrigatório
- 🎨 Emojis para facilitar identificação

---

## 📊 EXEMPLOS PRÁTICOS

### **Exemplo 1: Tenant só precisa de Geografia**

```json
// Criar API Key
POST /admin/apikeys
{
  "ownerId": "tenant-123",
  "scopes": ["geo"],
  "days": 90
}

// ✅ Pode acessar:
GET /geo/ufs
GET /geo/municipios

// ❌ NÃO pode acessar:
GET /cep/01310100 → 403 Forbidden
GET /cnpj/00000000000191 → 403 Forbidden
```

---

### **Exemplo 2: Tenant precisa de CEP e CNPJ**

```json
{
  "scopes": ["cep", "cnpj"]
}

// ✅ Pode:
GET /cep/01310100 → 200 OK
GET /cnpj/00000000000191 → 200 OK

// ❌ NÃO pode:
GET /geo/ufs → 403 Forbidden
```

---

### **Exemplo 3: Tenant Premium (all)**

```json
{
  "scopes": ["all"]
}

// ✅ Pode TUDO:
GET /geo/ufs → 200 OK
GET /cep/01310100 → 200 OK
GET /cnpj/00000000000191 → 200 OK
GET /moedas/cotacao → 200 OK (quando implementado)
```

---

## 🔄 CHECKLIST: ADICIONAR NOVO SCOPE

Quando implementar uma nova API, atualizar:

### **1. Backend**

- [ ] `internal/auth/scope_middleware.go`
  ```go
  validScopes := map[string]bool{
      "nova_api": true, // ← Adicionar aqui
  }
  ```

- [ ] `internal/http/router.go`
  ```go
  novaAPIGroup.Use(
      auth.AuthAPIKey(apikeys),
      auth.RequireScope(apikeys, "nova_api"), // ← Adicionar aqui
  )
  ```

### **2. Frontend**

- [ ] `components/apikeys/apikey-drawer.tsx`
  ```typescript
  { value: 'nova_api', label: '🆕 Nova API - Descrição' }
  ```

### **3. Documentação**

- [ ] Atualizar `docs/SCOPES_SYSTEM.md`
- [ ] Adicionar em `CHECKLIST_POS_IMPLEMENTACAO.md`

---

## 🧪 TESTES

### **Teste 1: Criar key com scope específico**
```bash
# Criar key com apenas 'geo'
POST /admin/apikeys
{ "scopes": ["geo"], "ownerId": "tenant-123" }

# Testar acesso
GET /geo/ufs → 200 OK ✅
GET /cep/01310100 → 403 Forbidden ❌
```

### **Teste 2: Criar key com 'all'**
```bash
POST /admin/apikeys
{ "scopes": ["all"], "ownerId": "tenant-123" }

# Testar acesso
GET /geo/ufs → 200 OK ✅
GET /cep/01310100 → 200 OK ✅
GET /cnpj/00000000000191 → 200 OK ✅
```

### **Teste 3: Scope inválido**
```bash
POST /admin/apikeys
{ "scopes": ["xyz"], "ownerId": "tenant-123" }

# Resposta:
400 Bad Request
{ "detail": "Scope 'xyz' não reconhecido" }
```

### **Teste 4: Scope futuro**
```bash
POST /admin/apikeys
{ "scopes": ["cpf"], "ownerId": "tenant-123" }

# Resposta:
400 Bad Request
{ "detail": "Scope 'cpf' ainda não está disponível" }
```

---

## 📈 BENEFÍCIOS

✅ **Segurança**: Tenants só acessam APIs autorizadas  
✅ **Flexibilidade**: Admin controla permissões granulares  
✅ **Escalabilidade**: Fácil adicionar novos scopes  
✅ **Monetização**: Scopes podem definir planos (Free, Pro, Business)  
✅ **Auditoria**: Logs mostram quais APIs cada tenant usa  
✅ **UX**: Interface visual intuitiva com checkboxes

---

## 🎁 PLANOS FUTUROS (Sugestão)

| Plano | Scopes Inclusos | Requests/Dia |
|-------|----------------|--------------|
| **Free** | `geo`, `cep` | 1.000 |
| **Pro** | `geo`, `cep`, `cnpj`, `cpf` | 10.000 |
| **Business** | `all` | 100.000 |
| **Enterprise** | `all` + features premium | Ilimitado |

---

## 📝 MIGRATION GUIDE

### **Migrar Keys Antigas (só tinham 'geo')**

```go
// Atualizar todas as keys antigas para incluir novos scopes
db.api_keys.updateMany(
  { scopes: ["geo:read"] },  // Keys antigas
  { $set: { scopes: ["all"] } }  // Dar acesso total
)
```

**Ou manualmente no admin:**
1. Listar todas as API Keys
2. Editar cada uma
3. Adicionar scopes: `cep`, `cnpj`
4. Salvar

---

## ✅ CHECKLIST IMPLEMENTAÇÃO

- [x] Middleware `RequireScope` criado
- [x] `ValidateAPIKeyScopes` na criação
- [x] Router com scope para geo, cep, cnpj
- [x] Frontend com 4 scopes (geo, cep, cnpj, all)
- [x] Documentação completa
- [x] Mensagens de erro amigáveis
- [x] hasScope() suporta 'all'
- [x] ValidationError no domain

---

**🎉 Sistema de Scopes 100% Funcional!**

**Próximo:** Quando implementar CPF, Moedas, etc, basta adicionar nos 3 lugares (middleware, router, frontend).

