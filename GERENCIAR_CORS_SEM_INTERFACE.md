# 🔧 Como Gerenciar CORS Quando Está Desabilitado

## 🚨 **Problema:**

Se você **desabilitar o CORS** via interface web, você **não conseguirá mais usar a interface web** para re-habilitar, porque o CORS está bloqueando as requests!

---

## ✅ **Solução 1: MongoDB Direto (Recomendado)**

### **Passo 1: Conectar no MongoDB**

```bash
# Local (Docker)
docker exec -it build-mongo-1 mongosh retech_core

# Produção (Railway)
# Use o Railway CLI ou MongoDB Compass com a connection string
```

### **Passo 2: Habilitar CORS**

```javascript
// No mongosh:
db.system_settings.updateOne(
  {},
  {
    $set: {
      "cors.enabled": true,
      "cors.allowedOrigins": [
        "https://core.theretech.com.br",
        "http://localhost:3000",
        "http://localhost:3001"
      ]
    }
  }
)
```

### **Passo 3: Verificar**

```javascript
db.system_settings.findOne({ }, { cors: 1 })
```

**Resultado esperado:**
```javascript
{
  _id: ObjectId("..."),
  cors: {
    enabled: true,
    allowedOrigins: [
      'https://core.theretech.com.br',
      'http://localhost:3000',
      'http://localhost:3001'
    ]
  }
}
```

---

## ✅ **Solução 2: curl sem Origin Header**

Se você fizer request **sem** o header `Origin`, o CORS **não é verificado**!

```bash
# Habilitar CORS via curl (sem Origin header)
curl -X PUT http://localhost:8080/admin/settings \
  -H "Authorization: Bearer SEU_TOKEN_JWT" \
  -H "Content-Type: application/json" \
  -d '{
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
      "accessTokenTTL": 900,
      "refreshTokenTTL": 604800
    },
    "api": {
      "version": "1.0.0",
      "environment": "development",
      "maintenance": false
    },
    "contact": {
      "whatsapp": "48999616679",
      "email": "suporte@theretech.com.br",
      "phone": "+55 48 99961-6679"
    },
    "cache": {
      "enabled": true,
      "cepTtlDays": 7,
      "cnpjTtlDays": 30,
      "maxSizeMb": 100,
      "autoCleanup": true
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
  }'
```

**Por que funciona?**
- Requests **sem** `Origin` header não são consideradas cross-origin
- CORS só se aplica a requests com `Origin` header (feitas pelo browser)
- curl, Postman, etc não enviam `Origin` por padrão

---

## ✅ **Solução 3: Postman (Disable CORS)**

No Postman, você pode:

1. **Desabilitar interceptor de CORS:**
   - Settings → General → Request Settings → "SSL certificate verification" OFF

2. **Ou simplesmente:** Postman **não aplica** política de CORS (ele é um cliente HTTP, não um browser)

---

## 📋 **Workflow Recomendado:**

### **Cenário 1: Desenvolvimento Local**

```
1. CORS.Enabled=true
2. AllowedOrigins=[http://localhost:3000, http://localhost:3001, ...]
3. Trabalhar normalmente via interface web
```

### **Cenário 2: Preciso Desabilitar CORS Temporariamente**

```
1. Desabilitar via interface web
2. Se precisar re-habilitar:
   a) Usar MongoDB direto (mongosh)
   b) Ou usar curl sem Origin header
   c) Ou usar Postman
3. Re-habilitar CORS
4. Voltar a usar interface web
```

### **Cenário 3: Produção**

```
1. CORS.Enabled=true SEMPRE
2. AllowedOrigins=[https://core.theretech.com.br]
3. Se precisar emergência:
   a) Conectar no MongoDB do Railway
   b) Ou usar Railway CLI com curl
```

---

## 🎯 **Por Que Não Tem Exceção para Localhost?**

### **❌ Problemas da Exceção:**

1. **Inconsistência:** Local funciona, produção quebra
2. **Segurança falsa:** Parece que está funcionando, mas não está
3. **Testes inválidos:** Não consegue testar o comportamento real
4. **Violação do princípio:** "Desabilitado" deveria ser desabilitado para TODOS

### **✅ Vantagens de NÃO Ter Exceção:**

1. **Consistência:** Comportamento igual em dev e prod
2. **Testes reais:** Você testa exatamente o que vai para produção
3. **Segurança:** Sem brechas "especiais"
4. **Clareza:** CORS desabilitado = desabilitado para TODOS

---

## 🧪 **Como Testar CORS Desabilitado:**

### **1. Desabilitar CORS:**

```bash
# Via MongoDB:
db.system_settings.updateOne({}, {$set: {"cors.enabled": false}})
```

### **2. Tentar Request com Origin:**

```bash
curl -H "Origin: http://localhost:3000" http://localhost:8080/health
```

**Resultado Esperado:**
```json
{
  "type": "https://retech-core/errors/cors-disabled",
  "title": "CORS Desabilitado",
  "status": 403,
  "detail": "CORS está desabilitado. Origin 'http://localhost:3000' não permitido. Configure em /admin/settings."
}
```

### **3. Tentar Request SEM Origin:**

```bash
curl http://localhost:8080/health
```

**Resultado Esperado:**
```json
{
  "status": "ok",
  "mongo": "ok",
  "success": true
}
```

✅ **Funciona porque não tem Origin header!**

---

## 📊 **Tabela de Comportamento:**

| Método | Origin Header | CORS.Enabled | Resultado |
|--------|---------------|--------------|-----------|
| **Browser** | ✅ Sempre | true | ✅ Se na lista |
| **Browser** | ✅ Sempre | false | ❌ 403 + erro |
| **curl** | ❌ Não (padrão) | true | ✅ Funciona |
| **curl** | ❌ Não (padrão) | false | ✅ Funciona |
| **curl** | ✅ Manual (-H) | true | ✅ Se na lista |
| **curl** | ✅ Manual (-H) | false | ❌ 403 + erro |
| **Postman** | ❌ Não (padrão) | true | ✅ Funciona |
| **Postman** | ❌ Não (padrão) | false | ✅ Funciona |

---

## 💡 **Dica Pro:**

Se você quiser **sempre** ter acesso ao admin/settings, mantenha CORS **sempre habilitado** e apenas controle a lista de origins:

```javascript
// CORS sempre ligado
cors: {
  enabled: true,
  allowedOrigins: [
    "https://core.theretech.com.br",  // Produção
    "http://localhost:3000",          // Dev local
    // Adicione/remova origins conforme necessário
  ]
}
```

---

**Resumo:** Sem exceções = comportamento consistente e testável! 🎯

