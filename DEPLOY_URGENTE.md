# 🚨 DEPLOY URGENTE - ROTAS PÚBLICAS

**Data:** 24 de Outubro de 2025  
**Prioridade:** 🔴 ALTA  
**Status:** ⏳ Aguardando deploy

---

## 🐛 PROBLEMA ATUAL

### **Erro em Produção:**
```json
{
  "error": "invalid api key format"
}
```

**Afeta:**
- ❌ https://core.theretech.com.br/playground
- ❌ https://core.theretech.com.br/ferramentas/consultar-cep
- ❌ https://core.theretech.com.br/ferramentas/validar-cnpj

**Causa:**
- Frontend em produção já está usando rotas `/public/*`
- Backend em produção NÃO tem essas rotas ainda
- Resultado: 401 Unauthorized ou erro de formato

---

## ✅ SOLUÇÃO

### **Backend precisa de deploy!**

**Arquivo modificado:**
- `internal/http/router.go`

**Mudança:**
```go
// Rotas públicas adicionadas (LINHAS 71-82)
publicGroup := r.Group("/public")
{
    publicGroup.GET("/cep/:codigo", cepHandler.GetCEP)
    publicGroup.GET("/cnpj/:numero", cnpjHandler.GetCNPJ)
    publicGroup.GET("/geo/ufs", geoHandler.ListUFs)
    publicGroup.GET("/geo/ufs/:sigla", geoHandler.GetUF)
}
```

**O que isso faz:**
- Cria rotas **SEM autenticação**
- Playground funciona sem API Key
- Ferramentas públicas funcionam
- Conversão muito maior

---

## 🚀 COMO FAZER DEPLOY

### **Opção 1: Railway (Automático)**

```bash
cd retech-core

# Commitar
git add internal/http/router.go docs/SEO_STRATEGY.md
git commit -m "feat(api): rotas públicas para playground"
git push

# Railway faz deploy automático
```

### **Opção 2: Docker Local (Teste)**

```bash
cd retech-core/build
docker-compose up -d --build api

# Testar
curl http://localhost:8080/public/cep/88101270
curl http://localhost:8080/public/cnpj/00000000000191
```

---

## ✅ VERIFICAÇÃO PÓS-DEPLOY

### **Testar em produção:**

```bash
# CEP
curl https://api-core.theretech.com.br/public/cep/01310100

# Deve retornar:
{
  "cep": "01310100",
  "logradouro": "Avenida Paulista",
  ...
}

# CNPJ
curl https://api-core.theretech.com.br/public/cnpj/00000000000191

# Deve retornar:
{
  "cnpj": "00000000000191",
  "razaoSocial": "BANCO DO BRASIL SA",
  ...
}
```

### **Testar no navegador:**
- https://core.theretech.com.br/playground
- Escolher CEP
- Digitar: 01310-100
- Clicar "Testar API"
- ✅ Deve funcionar!

---

## ⚠️ IMPORTANTE

**Ordem de deploy:**
1. ✅ Frontend já foi deployado (Railway automático)
2. ⏳ **Backend precisa de deploy AGORA**
3. ✅ Depois tudo funciona

**Sem o deploy do backend:**
- ❌ Playground quebrado
- ❌ Ferramentas quebradas
- ❌ Usuários não conseguem testar
- ❌ Conversão ZERO

**Com o deploy do backend:**
- ✅ Playground funciona
- ✅ Ferramentas funcionam
- ✅ Usuários testam
- ✅ Conversão 10-15%

---

## 🎯 COMMIT SUGERIDO

```bash
cd retech-core

git add internal/http/router.go docs/SEO_STRATEGY.md

git commit -m "feat(api): rotas públicas para playground e ferramentas

PROBLEMA:
- Playground em produção quebrado
- Ferramentas públicas retornando erro
- 'invalid api key format'

CAUSA:
- Frontend usando /public/* (sem API Key)
- Backend não tinha essas rotas

SOLUÇÃO:
✅ Criado grupo /public/* sem autenticação
- GET /public/cep/:codigo
- GET /public/cnpj/:numero  
- GET /public/geo/ufs
- GET /public/geo/ufs/:sigla

BENEFÍCIOS:
- Playground funciona sem API Key
- Ferramentas públicas funcionam
- Barreira de entrada ZERO
- Conversão muito maior

TESTADO LOCAL:
✅ curl http://localhost:8080/public/cep/88101270
✅ curl http://localhost:8080/public/cnpj/00000000000191

STATUS: Pronto para produção"

git push
```

---

## 📋 CHECKLIST

- [x] Código commitado localmente
- [x] Testado em ambiente local
- [x] Build passando
- [ ] **Push para GitHub** ← FAZER AGORA
- [ ] **Aguardar deploy Railway** (2-3 minutos)
- [ ] **Testar em produção**
- [ ] **Verificar playground funcionando**

---

**⏰ URGENTE: Faça o push do backend AGORA para o playground voltar a funcionar!**

**Após o deploy, TUDO vai funcionar perfeitamente! ✅**

