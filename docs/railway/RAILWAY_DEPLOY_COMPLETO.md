# 🚀 DEPLOY COMPLETO NO RAILWAY

## 🎯 Arquitetura Final

```
core.theretech.com.br/          → Frontend (Next.js)
core.theretech.com.br/admin     → Admin Dashboard  
core.theretech.com.br/painel    → Developer Portal
core.theretech.com.br/api       → Backend Go (API)
```

---

## 📋 CHECKLIST DE DEPLOY

### 1️⃣ BACKEND (retech-core) - JÁ CONFIGURADO ✅

**Status**: ✅ Já está no Railway com `core.theretech.com.br`

**Variáveis de ambiente necessárias:**
```bash
# MongoDB
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/
DB=retech_core

# JWT
JWT_ACCESS_SECRET=seu_jwt_access_secret_super_seguro
JWT_REFRESH_SECRET=seu_jwt_refresh_secret_super_seguro
JWT_ACCESS_TTL=900
JWT_REFRESH_TTL=604800

# Server
PORT=8080
ENV=production
```

**Para configurar no Railway:**
1. Acesse seu projeto no Railway
2. Clique em "Variables" 
3. Adicione todas as variáveis acima
4. **IMPORTANTE**: Mude `BACKEND_URL` para `http://localhost:8080` (interno)

---

### 2️⃣ FRONTEND (retech-core-admin) - NOVO DEPLOY

**Vamos criar um novo serviço no Railway:**

#### Passo 1: Conectar repositório
1. No Railway, clique "New Project"
2. "Deploy from GitHub repo"
3. Selecione `retech-core-admin`
4. Branch: `main`

#### Passo 2: Configurar variáveis
```bash
# Backend URL (INTERNO do Railway)
BACKEND_URL=http://retech-core:8080

# Next.js
NODE_ENV=production
NEXT_TELEMETRY_DISABLED=1
```

#### Passo 3: Configurar domínio
1. No serviço frontend, vá em "Settings"
2. "Domains" → "Custom Domain"
3. Adicione: `core.theretech.com.br`
4. Configure DNS no seu provedor

---

## 🔧 CONFIGURAÇÕES TÉCNICAS

### Backend (Go)
- **Porta**: 8080 (interno)
- **Health**: `/health`
- **Seeds**: Embedados na imagem
- **MongoDB**: Conectado via variável

### Frontend (Next.js)
- **Porta**: 3000 (interno)
- **Proxy**: `/api/*` → Backend
- **Build**: `output: 'standalone'`
- **Docker**: Multi-stage otimizado

---

## 🌐 CONFIGURAÇÃO DNS

**No seu provedor DNS:**
```
Tipo: CNAME
Nome: core
Valor: [URL do Railway]
TTL: 300
```

**Exemplo:**
```
core.theretech.com.br → your-app.railway.app
```

---

## 🧪 TESTE PÓS-DEPLOY

### 1. Testar Backend
```bash
curl https://core.theretech.com.br/api/health
# Deve retornar: {"status":"ok","timestamp":"..."}
```

### 2. Testar Frontend
```bash
# Acessar
https://core.theretech.com.br/
https://core.theretech.com.br/admin/login
https://core.theretech.com.br/painel/login
```

### 3. Testar Proxy
```bash
# Deve funcionar via frontend
curl https://core.theretech.com.br/api/health
```

---

## 👨‍💼 CRIAR SUPER ADMIN EM PRODUÇÃO

### Opção 1: Via API (Recomendado)
```bash
curl -X POST https://core.theretech.com.br/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@theretech.com.br",
    "name": "Super Admin",
    "password": "admin12345678",
    "role": "SUPER_ADMIN"
  }'
```

### Opção 2: Via Script
```bash
# No Railway, acesse o terminal do backend
./create-admin.sh
```

---

## 🚨 TROUBLESHOOTING

### Problema: Frontend não carrega
**Solução**: Verificar se `BACKEND_URL` está correto no frontend

### Problema: API retorna 404
**Solução**: Verificar se o proxy `/api/*` está funcionando

### Problema: Login não funciona
**Solução**: Verificar JWT secrets e MongoDB connection

### Problema: Seeds não aplicaram
**Solução**: Verificar se arquivos estão na imagem Docker

---

## 📊 MONITORAMENTO

### Logs do Backend
```bash
# No Railway, vá em "Deployments" → "View logs"
```

### Logs do Frontend  
```bash
# No Railway, vá em "Deployments" → "View logs"
```

### Health Checks
- Backend: `https://core.theretech.com.br/api/health`
- Frontend: `https://core.theretech.com.br/`

---

## 🎉 RESULTADO FINAL

Após o deploy, você terá:

✅ **Backend**: `core.theretech.com.br/api/*`  
✅ **Admin**: `core.theretech.com.br/admin/*`  
✅ **Portal**: `core.theretech.com.br/painel/*`  
✅ **Landing**: `core.theretech.com.br/`  

**Sistema 100% funcional em produção!** 🚀

---

## 🚀 PRÓXIMOS PASSOS

1. ✅ Deploy backend (já feito)
2. 🔄 Deploy frontend (novo serviço)
3. 🔄 Configurar DNS
4. 🔄 Testar tudo
5. 🔄 Criar super admin
6. 🎉 **Sistema no ar!**

**Vamos fazer isso!** 💪
