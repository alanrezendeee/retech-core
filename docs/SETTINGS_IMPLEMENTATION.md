# ⚙️ Implementação do Sistema de Configurações

## 📋 Visão Geral

Sistema completo de configurações globais para o Retech Core API, incluindo:
- ✅ Configurações DEFAULT de rate limiting (aplicadas a novos tenants)
- ✅ Configurações de CORS
- ✅ Configurações de JWT
- ✅ Informações da API (versão, ambiente, modo de manutenção)
- ✅ Rate limiting personalizado por tenant
- ✅ Tela de admin modernizada

---

## 🏗️ Arquitetura

### **1. Backend (Go)**

#### **Domain Models:**
```
internal/domain/settings.go
├── SystemSettings          // Configurações globais do sistema
├── CORSConfig             // Configuração de CORS
├── JWTConfig              // Configuração de JWT
└── APIConfig              // Informações da API
```

#### **Repository:**
```
internal/storage/settings_repo.go
├── Get()                  // Buscar configurações (ou padrão)
├── Update()               // Atualizar configurações (upsert)
└── Ensure()               // Garantir configurações padrão existam
```

#### **Handlers:**
```
internal/http/handlers/settings.go
├── GET /admin/settings    // Retorna configurações atuais
└── PUT /admin/settings    // Atualiza configurações
```

#### **Middleware:**
```
internal/middleware/rate_limiter.go
└── getRateLimitConfig()   // Busca config por tenant ou usa padrão
```

---

### **2. Frontend (React/Next.js)**

```
app/admin/settings/page.tsx
└── Tela modernizada com:
    ├── Header com ícones e gradientes
    ├── Card informativo
    ├── Grid 2 colunas (responsive)
    ├── Cards para cada seção de config
    └── Validações e feedback visual
```

---

## 🔑 Conceitos-Chave

### **1. Rate Limiting: DEFAULT vs POR TENANT**

#### **DEFAULT (Global):**
- Definido em `/admin/settings`
- Aplicado automaticamente a **novos tenants**
- Não afeta tenants já existentes
- Exemplo: 1.000 req/dia, 60 req/minuto

#### **POR TENANT (Personalizado):**
- Opcional: campo `rateLimit` no Tenant
- Se `null`, usa o DEFAULT do sistema
- Se definido, sobrescreve o DEFAULT
- Editável por tenant (futuro: `/admin/tenants/:id/edit`)

---

### **2. Fluxo do Rate Limiting:**

```
┌────────────────────────────────────────────────────────┐
│  Request com API Key                                   │
└────────────────────────────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────────────────┐
│  Middleware: AuthAPIKey                                │
│  ├── Valida API Key                                    │
│  ├── Adiciona tenant_id no contexto                    │
│  └── Adiciona api_key no contexto                      │
└────────────────────────────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────────────────┐
│  Middleware: RateLimiter                               │
│  ├── Busca tenant pelo tenant_id                       │
│  ├── Verifica se tenant tem rateLimit personalizado    │
│  ├── Se SIM: usa config do tenant                      │
│  ├── Se NÃO: busca SystemSettings.defaultRateLimit    │
│  ├── Aplica limite                                     │
│  └── Retorna 429 se exceder                            │
└────────────────────────────────────────────────────────┘
                      │
                      ▼
┌────────────────────────────────────────────────────────┐
│  Request processado normalmente                        │
└────────────────────────────────────────────────────────┘
```

---

## 📊 Estrutura de Dados

### **SystemSettings (MongoDB):**
```json
{
  "_id": "system-settings-singleton",
  "defaultRateLimit": {
    "requestsPerDay": 1000,
    "requestsPerMinute": 60
  },
  "cors": {
    "enabled": true,
    "allowedOrigins": [
      "https://core.theretech.com.br",
      "http://localhost:3000",
      "http://localhost:3001"
    ]
  },
  "jwt": {
    "accessTokenTTL": 900,      // 15 minutos
    "refreshTokenTTL": 604800   // 7 dias
  },
  "api": {
    "version": "1.0.0",
    "environment": "development",
    "maintenance": false
  },
  "createdAt": "2025-10-21T18:00:00Z",
  "updatedAt": "2025-10-21T18:30:00Z"
}
```

### **Tenant com Rate Limit Personalizado:**
```json
{
  "_id": ObjectId("..."),
  "tenantId": "tenant-enterprise-001",
  "name": "Empresa Premium",
  "email": "contato@premium.com",
  "active": true,
  "rateLimit": {              // ← PERSONALIZADO!
    "requestsPerDay": 100000,  // 100k req/dia
    "requestsPerMinute": 1000  // 1k req/minuto
  },
  "createdAt": "2025-10-21T18:00:00Z",
  "updatedAt": "2025-10-21T18:30:00Z"
}
```

### **Tenant sem Rate Limit Personalizado:**
```json
{
  "_id": ObjectId("..."),
  "tenantId": "tenant-free-001",
  "name": "Startup Free",
  "email": "contato@startup.com",
  "active": true,
  "rateLimit": null,  // ← Usa DEFAULT do sistema
  "createdAt": "2025-10-21T18:00:00Z",
  "updatedAt": "2025-10-21T18:30:00Z"
}
```

---

## 🎨 Tela de Settings (/admin/settings)

### **Melhorias Implementadas:**

#### **1. Header Modernizado:**
- ✅ Ícone de configurações com gradiente azul/roxo
- ✅ Título grande e descrição
- ✅ Botão "Salvar" com gradiente e ícone
- ✅ Loading state no botão

#### **2. Card Informativo:**
- ✅ Gradiente de fundo (azul/roxo)
- ✅ Ícone de informação
- ✅ Explicação clara sobre cada configuração
- ✅ Destaque para conceito de DEFAULT vs POR TENANT

#### **3. Grid Responsivo (2 colunas):**
```
Desktop:
┌──────────────────┬──────────────────┐
│ Rate Limit       │ CORS             │
├──────────────────┼──────────────────┤
│ JWT              │ API Info         │
└──────────────────┴──────────────────┘

Mobile:
┌──────────────────┐
│ Rate Limit       │
├──────────────────┤
│ CORS             │
├──────────────────┤
│ JWT              │
├──────────────────┤
│ API Info         │
└──────────────────┘
```

#### **4. Cards Individuais:**

**Rate Limiting DEFAULT:**
- ✅ Ícone roxo (Shield)
- ✅ Alerta amarelo explicando que é apenas para novos tenants
- ✅ Inputs numéricos com hints
- ✅ Validação de limites

**CORS:**
- ✅ Ícone verde (Globe)
- ✅ Switch para ativar/desativar
- ✅ Textarea para múltiplas origens
- ✅ Hint de formatação

**JWT:**
- ✅ Ícone azul (Key)
- ✅ Inputs para Access e Refresh Token TTL
- ✅ Conversão segundos → minutos/dias nos hints

**API Info:**
- ✅ Ícone laranja (Database)
- ✅ Badges para versão e ambiente
- ✅ Switch para modo de manutenção
- ✅ Alerta vermelho quando em manutenção

---

## 🔒 Segurança e Validações

### **Backend:**
```go
// Validações em handlers/settings.go
RequestsPerDay:      1 a 1.000.000
RequestsPerMinute:   1 a 10.000
AccessTokenTTL:      60 a 3.600 segundos
RefreshTokenTTL:     3.600 a 2.592.000 segundos
```

### **Frontend:**
```typescript
// Validações em page.tsx
- Inputs numéricos com min/max
- Feedback visual de erros
- Toast de sucesso/erro
- Loading states
```

---

## 🧪 Testes

### **1. Testar GET /admin/settings:**
```bash
# Login como admin
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alanrezendeee@gmail.com","password":"admin123456"}' \
  | jq -r '.accessToken')

# Buscar configurações
curl -s http://localhost:8080/admin/settings \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.'
```

**Resposta esperada:**
```json
{
  "id": "system-settings-singleton",
  "defaultRateLimit": {
    "requestsPerDay": 1000,
    "requestsPerMinute": 60
  },
  "cors": {
    "enabled": true,
    "allowedOrigins": [...]
  },
  ...
}
```

---

### **2. Testar PUT /admin/settings:**
```bash
curl -s -X PUT http://localhost:8080/admin/settings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "defaultRateLimit": {
      "requestsPerDay": 2000,
      "requestsPerMinute": 100
    },
    "cors": {
      "enabled": true,
      "allowedOrigins": ["https://core.theretech.com.br"]
    },
    "jwt": {
      "accessTokenTTL": 900,
      "refreshTokenTTL": 604800
    },
    "api": {
      "version": "1.0.0",
      "environment": "production",
      "maintenance": false
    }
  }' \
  | jq '.'
```

---

### **3. Testar Rate Limiting com Tenant Personalizado:**

**Criar tenant com limite personalizado:**
```bash
curl -s -X PUT http://localhost:8080/admin/tenants/tenant-001 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rateLimit": {
      "requestsPerDay": 5000,
      "requestsPerMinute": 200
    }
  }' \
  | jq '.'
```

**Verificar se middleware usa config personalizada:**
```bash
# Fazer requisições com API key do tenant-001
# Verificar headers:
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4999
```

---

## 🚀 Próximos Passos

### **1. Tela de Edição de Tenant (settings-6):**
```
/admin/tenants/:id/edit
└── Adicionar seção "Rate Limiting Personalizado"
    ├── Switch: "Usar limite personalizado?"
    ├── Se SIM: mostrar inputs
    ├── Se NÃO: usar DEFAULT do sistema
    └── Salvar no campo tenant.rateLimit
```

### **2. Melhorias Futuras:**
- [ ] Logs de auditoria (quem mudou o quê)
- [ ] Histórico de configurações
- [ ] Validação de CORS em tempo real
- [ ] Preview de JWT TTL (ex: "15 min" em vez de 900)
- [ ] Importar/Exportar configurações
- [ ] Configurações por ambiente (dev/prod)

---

## ✅ Checklist de Implementação

- [x] Domain model `SystemSettings`
- [x] Repository `SettingsRepo`
- [x] Handlers GET/PUT `/admin/settings`
- [x] Endpoints registrados no router
- [x] Middleware de rate limiting atualizado
- [x] Busca config por tenant ou DEFAULT
- [x] Tela frontend modernizada
- [x] Grid responsivo
- [x] Cards com ícones e gradientes
- [x] Validações backend e frontend
- [x] Toast de feedback
- [x] Loading states
- [x] Campo `rateLimit` no Tenant
- [x] Documentação completa
- [ ] Tela de edição de tenant (rate limit)
- [ ] Testes automatizados

---

## 📝 Resumo

**O que mudou:**

1. **Backend:**
   - SystemSettings armazenado no MongoDB (singleton)
   - Configurações padrão garantidas no startup
   - Endpoints protegidos por SUPER_ADMIN
   - Middleware de rate limiting inteligente

2. **Frontend:**
   - Tela moderna e responsiva
   - Explicações claras
   - Feedback visual excelente
   - API client integrado

3. **Rate Limiting:**
   - DEFAULT: usado para novos tenants
   - POR TENANT: opcional, sobrescreve DEFAULT
   - Middleware busca automaticamente

**Status:** ✅ Funcional e pronto para uso!

