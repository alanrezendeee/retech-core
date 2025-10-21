# 🔧 Variáveis de Ambiente - Retech Core

## 📋 Variáveis Obrigatórias

### **Backend (retech-core)**

```bash
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGODB_NAME=retech_core

# JWT
JWT_ACCESS_SECRET=your-super-secret-key-change-in-production
JWT_REFRESH_SECRET=your-super-secret-refresh-key-change-in-production
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# API Key
APIKEY_HASH_SECRET=your-apikey-hash-secret-change-in-production

# Server
HTTP_PORT=8080

# Environment (IMPORTANTE!)
ENV=production  # ← Define o ambiente (development, staging, production)
```

---

## 🎯 Variável ENV / NODE_ENV

A variável `ENV` (ou `NODE_ENV`) é usada para:

1. **Definir o ambiente da API** mostrado em `/admin/settings`
2. **Configurar logs** (mais verboso em dev, menos em prod)
3. **Ativar/desativar features** de desenvolvimento

### **Valores aceitos:**

```bash
ENV=development  # Desenvolvimento local
ENV=staging      # Ambiente de testes
ENV=production   # Produção (Railway)
```

---

## 🚀 Railway - Configuração de Produção

### **retech-core (Backend API)**

Adicione esta variável no Railway:

```bash
ENV=production
```

**Como adicionar:**

1. Acesse: https://railway.app
2. Vá para o serviço `retech-core`
3. Clique em **Variables**
4. Adicione:
   - **Key:** `ENV`
   - **Value:** `production`
5. Clique em **Add** e **Redeploy**

---

### **retech-core-admin (Frontend)**

Já configurado! O frontend detecta automaticamente o ambiente do backend.

---

## 📊 Como Verificar

### **1. Via API diretamente:**

```bash
curl https://api-core.theretech.com.br/admin/settings \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Resposta:
```json
{
  "api": {
    "version": "1.0.0",
    "environment": "production",  // ← Deve ser "production"!
    "maintenance": false
  }
}
```

---

### **2. Via Dashboard:**

1. Login em: https://core.theretech.com.br/admin/login
2. Vá em: **Settings** (menu lateral)
3. Veja a seção **"Informações da API"**:
   - **Versão:** v1.0.0
   - **Ambiente:** production ✅

---

## 🔄 Fluxo de Detecção

```go
// Backend: internal/domain/settings.go
env := os.Getenv("ENV")           // 1. Tenta ENV
if env == "" {
    env = os.Getenv("NODE_ENV")   // 2. Fallback para NODE_ENV
}
if env == "" {
    env = "development"           // 3. Fallback para development
}
```

**Prioridade:**
1. `ENV` (mais específico)
2. `NODE_ENV` (compatibilidade Node.js)
3. `"development"` (fallback seguro)

---

## 🐛 Troubleshooting

### **Problema: Aparece "development" em produção**

**Causa:** Variável `ENV` não está setada no Railway.

**Solução:**
1. Railway → Service `retech-core` → Variables
2. Adicionar `ENV=production`
3. Redeploy

---

### **Problema: Mudei ENV mas não atualizou**

**Causa:** Settings já existem no MongoDB com valor antigo.

**Solução 1 (Recomendado):** Resetar settings via API:
```bash
# Deletar settings existentes
curl -X DELETE https://api-core.theretech.com.br/admin/settings \
  -H "Authorization: Bearer YOUR_TOKEN"

# Backend vai recriar com ENV correto
```

**Solução 2:** Atualizar manualmente:
```bash
# Via MongoDB
docker exec -it retech-mongo mongosh
use retech_core
db.system_settings.updateOne(
  {},
  { $set: { "api.environment": "production" } }
)
```

---

### **Problema: ENV=production mas API não respeita**

**Causa:** Code ainda não foi rebuilded.

**Solução:**
```bash
# Railway: Force Redeploy
# Ou localmente:
cd /path/to/retech-core
go build ./cmd/api
```

---

## 📝 Exemplo Completo - Railway

### **Variáveis do retech-core:**

```bash
# Obrigatórias
MONGO_URI=mongodb://mongo.railway.internal:27017
MONGODB_NAME=retech_core
JWT_ACCESS_SECRET=super-secret-production-key-xyz123
JWT_REFRESH_SECRET=super-secret-refresh-production-key-abc456
APIKEY_HASH_SECRET=super-secret-apikey-hash-xyz789
HTTP_PORT=8080

# Ambiente (NOVA!)
ENV=production  # ← ADICIONAR ESTA!

# Opcionais (já têm fallbacks)
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
APIKEY_TTL_DAYS=90
```

---

## ✅ Checklist Pós-Deploy

Após adicionar `ENV=production` no Railway:

- [ ] Variável `ENV=production` adicionada no Railway
- [ ] Redeploy do backend concluído
- [ ] API respondendo (curl /health)
- [ ] Login funcionando
- [ ] Settings mostrando "production" na UI
- [ ] Logs menos verbosos (produção)

---

## 🎓 Boas Práticas

### **1. Nunca commitar secrets**

❌ **ERRADO:**
```bash
# .env (no git)
JWT_ACCESS_SECRET=my-secret-key
```

✅ **CORRETO:**
```bash
# .env.example (no git)
JWT_ACCESS_SECRET=your-secret-here

# Railway Variables (não vai pro git)
JWT_ACCESS_SECRET=actual-secret-production-key-xyz
```

---

### **2. Usar ENV específico para cada ambiente**

```bash
# Local
ENV=development

# Staging (futuro)
ENV=staging

# Produção
ENV=production
```

---

### **3. Documentar variáveis obrigatórias**

Manter este documento atualizado sempre que adicionar novas variáveis!

---

**Data:** 2025-10-21  
**Versão:** 1.0.0  
**Status:** ✅ Documentado

