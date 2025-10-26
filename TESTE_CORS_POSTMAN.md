# 🧪 Guia Completo de Teste de CORS via Postman

## 📋 **O que é CORS e como funciona no Retech Core**

### **Arquitetura de CORS:**

```
┌──────────────────────────────────────────────┐
│  REQUEST INCOMING                            │
└──────────────────┬───────────────────────────┘
                   │
                   ▼
         ┌─────────────────┐
         │  É rota pública? │
         │  /public/*       │
         │  /health, /docs  │
         └────┬────────┬────┘
              │        │
         SIM ▼        │ NÃO
    ┌──────────────┐  │
    │ CORS: *      │  │
    │ (irrestrito) │  │
    │ ✅ Permitido │  │
    └──────────────┘  │
                      ▼
              ┌──────────────────┐
              │ Origin=localhost?│
              └────┬───────┬─────┘
                   │       │
              SIM ▼       │ NÃO
         ✅ Permitido     │
         (dev mode)       │
                          ▼
                  ┌──────────────────┐
                  │ Verifica Settings│
                  │ CORS.Enabled?    │
                  └────┬───────┬─────┘
                       │       │
                  SIM ▼       │ NÃO
             ┌──────────────┐ │
             │ Origin está  │ │
             │ na lista?    │ │
             └──┬──────┬───┘ │
                │      │     │
           SIM ▼      │ NÃO │
        ✅ Permitido  │     │
                      ▼     ▼
                  ❌ Bloqueado
```

---

## 🔧 **Passo 1: Importar Collection no Postman**

1. **Abra o Postman**
2. **Import → File → `POSTMAN_CORS_TEST.json`**
3. **Configure as variáveis de ambiente:**

### **Variáveis:**

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `BASE_URL` | URL do backend | `http://localhost:8080` (local) ou `https://api-core.theretech.com.br` (prod) |
| `AUTH_TOKEN` | JWT token do admin | Obter após login em `/auth/login` |
| `API_KEY` | API Key válida | Criar em `/admin/apikeys` |

---

## 🎯 **Passo 2: Obter Token de Autenticação**

### **Request: Login Admin**

```bash
POST {{BASE_URL}}/auth/login
Content-Type: application/json

{
  "email": "admin@theretech.com.br",
  "password": "sua_senha_aqui"
}
```

**Resposta:**
```json
{
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "user": {
    "id": "...",
    "email": "admin@theretech.com.br",
    "role": "SUPER_ADMIN"
  }
}
```

**👉 Copie o `accessToken` e cole na variável `{{AUTH_TOKEN}}`**

---

## 🧪 **Passo 3: Executar Testes de CORS**

### **Teste 1: Health Check (Rota Pública)**

```bash
GET {{BASE_URL}}/health
Origin: http://localhost:3000
```

**✅ Resultado Esperado:**
```
Status: 200 OK
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization, X-Requested-With, X-API-Key
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 86400
```

---

### **Teste 2: OPTIONS Preflight (Rota Pública)**

```bash
OPTIONS {{BASE_URL}}/public/playground/status
Origin: http://localhost:3000
Access-Control-Request-Method: GET
Access-Control-Request-Headers: Content-Type, Authorization
```

**✅ Resultado Esperado:**
```
Status: 204 No Content
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization, X-Requested-With, X-API-Key
Access-Control-Max-Age: 86400
```

---

### **Teste 3: Playground Status (Rota Pública)**

```bash
GET {{BASE_URL}}/public/playground/status
Origin: http://localhost:3000
```

**✅ Resultado Esperado:**
```
Status: 200 OK
Access-Control-Allow-Origin: http://localhost:3000 (ou *)
Content-Type: application/json

{
  "enabled": true,
  "apiKey": "rtc_demo_playground_2024",
  "allowedApis": ["cep", "cnpj", "geo"]
}
```

---

### **Teste 4: Admin Settings GET (localhost permitido)**

```bash
GET {{BASE_URL}}/admin/settings
Origin: http://localhost:3000
Authorization: Bearer {{AUTH_TOKEN}}
```

**✅ Resultado Esperado:**
```
Status: 200 OK
Access-Control-Allow-Origin: http://localhost:3000
Content-Type: application/json

{
  "data": {
    "defaultRateLimit": {...},
    "cors": {...},
    ...
  }
}
```

---

### **Teste 5: Admin Settings GET (Origin NÃO permitido)**

```bash
GET {{BASE_URL}}/admin/settings
Origin: https://malicious-site.com
Authorization: Bearer {{AUTH_TOKEN}}
```

**✅ Resultado Esperado:**
```
Status: 200 OK (auth válido)
❌ NO Access-Control-Allow-Origin header
(ou header presente mas com valor diferente de https://malicious-site.com)
```

**👉 Browser bloqueará o response devido à falta de CORS header**

---

### **Teste 6: CEP API (Rota Protegida com API Key)**

```bash
GET {{BASE_URL}}/cep/01310100
Origin: http://localhost:3000
X-API-Key: {{API_KEY}}
```

**✅ Resultado Esperado:**
```
Status: 200 OK
Access-Control-Allow-Origin: http://localhost:3000
Content-Type: application/json

{
  "cep": "01310-100",
  "logradouro": "Avenida Paulista",
  ...
}
```

---

### **Teste 7: Admin Settings PUT (Salvar Configurações)**

```bash
PUT {{BASE_URL}}/admin/settings
Origin: http://localhost:3000
Authorization: Bearer {{AUTH_TOKEN}}
Content-Type: application/json

{
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
  ...
}
```

**✅ Resultado Esperado:**
```
Status: 200 OK
Access-Control-Allow-Origin: http://localhost:3000
Content-Type: application/json

{
  "message": "Configurações atualizadas com sucesso",
  "settings": {...}
}
```

---

## 📊 **Tabela de Comportamento Esperado:**

| Rota | Origin | CORS Enabled | Na Lista? | Resultado |
|------|--------|--------------|-----------|-----------|
| `/health` | qualquer | N/A | N/A | ✅ CORS permitido |
| `/public/*` | qualquer | N/A | N/A | ✅ CORS permitido |
| `/admin/settings` | `localhost:*` | qualquer | N/A | ✅ CORS permitido (dev mode) |
| `/admin/settings` | `core.theretech.com.br` | ✅ true | ✅ sim | ✅ CORS permitido |
| `/admin/settings` | `core.theretech.com.br` | ✅ true | ❌ não | ❌ CORS bloqueado |
| `/admin/settings` | `malicious.com` | ✅ true | ❌ não | ❌ CORS bloqueado |
| `/admin/settings` | qualquer | ❌ false | N/A | ❌ CORS bloqueado (exceto localhost) |
| `/cep/:codigo` | `localhost:*` | qualquer | N/A | ✅ CORS permitido (dev mode) |
| `/cep/:codigo` | `core.theretech.com.br` | ✅ true | ✅ sim | ✅ CORS permitido |

---

## 🔍 **Como Debugar Problemas de CORS:**

### **1. Verificar Logs do Backend:**

```bash
# Local (Docker)
docker-compose -f build/docker-compose.yml logs api --tail=50

# Procure por:
[CORS] GET /admin/settings (Origin: http://localhost:3000)
[CORS] Settings OK, CORS.Enabled=true
[CORS] Origin na lista: true
```

---

### **2. Verificar Headers no Postman:**

1. **Clique na aba "Headers" na resposta**
2. **Procure por:**
   - `Access-Control-Allow-Origin`
   - `Access-Control-Allow-Methods`
   - `Access-Control-Allow-Headers`
   - `Access-Control-Allow-Credentials`

---

### **3. Verificar Configuração no MongoDB:**

```bash
# No terminal local:
docker exec -it build-mongo-1 mongosh retech_core

# No mongosh:
db.system_settings.findOne()
```

**Verifique:**
```javascript
{
  cors: {
    enabled: true,
    allowedOrigins: [
      "https://core.theretech.com.br",
      "http://localhost:3000",
      ...
    ]
  }
}
```

---

## 🚨 **Troubleshooting:**

### **Problema: "CORS blocked" no navegador mas Postman funciona**

**Causa:** Postman não aplica política CORS (é um cliente HTTP, não um browser)

**Solução:** 
1. Verificar se o header `Access-Control-Allow-Origin` está presente na resposta
2. Se estiver presente com valor correto, o problema pode ser cache do browser
3. Limpar cache do browser (Cmd+Shift+Delete no Chrome)

---

### **Problema: Admin Settings não salva (erro no frontend)**

**Causa:** CORS bloqueando a resposta do PUT

**Solução:**
1. Verificar no Postman se PUT funciona
2. Se sim, problema é no frontend (cache ou URL errada)
3. Se não, problema é no backend (verificar logs)

---

### **Problema: `Access-Control-Allow-Origin` não aparece**

**Causa:** Backend não está adicionando headers CORS

**Possíveis razões:**
1. Origin não está na lista permitida
2. CORS.Enabled=false e origin não é localhost
3. Erro ao buscar settings do MongoDB

**Solução:**
```bash
# Verificar logs:
[CORS] GET /admin/settings (Origin: http://localhost:3000)
[CORS] Settings OK, CORS.Enabled=false
[CORS] CORS desabilitado, mas permitindo localhost  # ✅ Deve aparecer

# Se não aparecer, há um bug no código
```

---

## ✅ **Checklist Final:**

- [ ] Collection importada no Postman
- [ ] Variáveis configuradas (`BASE_URL`, `AUTH_TOKEN`, `API_KEY`)
- [ ] Teste 1: Health Check - ✅ CORS headers presentes
- [ ] Teste 2: OPTIONS Preflight - ✅ Status 204
- [ ] Teste 3: Playground Status - ✅ CORS permitido
- [ ] Teste 4: Admin Settings GET (localhost) - ✅ CORS permitido
- [ ] Teste 5: Admin Settings GET (origin inválido) - ❌ CORS bloqueado
- [ ] Teste 6: CEP API - ✅ CORS permitido
- [ ] Teste 7: Admin Settings PUT - ✅ Salva com sucesso

---

## 🎯 **Próximos Passos:**

1. ✅ Rodar todos os testes no Postman
2. ✅ Verificar que rotas públicas SEMPRE têm CORS
3. ✅ Verificar que localhost SEMPRE é permitido (dev mode)
4. ✅ Verificar que origins da lista são permitidos quando CORS.Enabled=true
5. ✅ Verificar que origins fora da lista são bloqueados

---

**Documentação criada em:** 2024-10-26  
**Versão:** 1.0.0  
**Autor:** Retech Core Team

