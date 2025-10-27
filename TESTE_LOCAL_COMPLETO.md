# 🧪 GUIA DE TESTE LOCAL COMPLETO

## 📋 PRÉ-REQUISITOS

**Backend:**
```bash
cd /Users/alanleitederezende/source/theretech/projetos-producao/retech-core
docker-compose -f build/docker-compose.yml up --build -d
```
✅ Rodando em: http://localhost:8080

**Frontend:**
```bash
cd /Users/alanleitederezende/source/theretech/projetos-producao/retech-core-admin
yarn dev
```
✅ Rodando em: http://localhost:3001

---

## 🎯 TESTE 1: Gerar Nova API Key

### Passo 1: Configuração Inicial
1. Abrir: http://localhost:3001/admin/settings
2. Se não estiver logado, fazer login
3. Rolar até "Playground Público"
4. **Verificar:**
   - ✅ Playground: ON
   - ✅ API Key Demo: **VAZIO**
   - ✅ Requests/Dia: 4
   - ✅ Requests/Min: 1
   - ✅ APIs: CEP ✓, CNPJ ✓

### Passo 2: Gerar API Key
1. **Clicar:** Botão "🔑 Gerar Nova"
2. **Aguardar:** "Gerando..."
3. **Verificar toast:** "API Key gerada com sucesso! Scopes: cep, cnpj"
4. **Verificar campo:** API Key aparece (formato: `rtc_demo_playground_20251027.ABC123...`)

### Passo 3: Validar no MongoDB
```bash
docker exec build-mongo-1 mongosh retech_core --eval '
db.api_keys.find({keyId: {$regex: "^rtc_demo_playground"}}).pretty()
'
```

**Deve mostrar:**
```javascript
{
  keyId: "rtc_demo_playground_20251027.ABC123...",
  keyHash: "rtc_demo_playground_20251027.ABC123...",
  scopes: ["cep", "cnpj"],  // ✅ Baseado nos checkboxes!
  ownerId: "playground-public",
  expiresAt: ISODate("2035-10-27..."),
  revoked: false,
  createdAt: ISODate("2025-10-27...")
}
```

### Passo 4: Validar em Settings
```bash
docker exec build-mongo-1 mongosh retech_core --eval '
db.system_settings.findOne({}, {playground: 1})
'
```

**Deve mostrar:**
```javascript
{
  playground: {
    enabled: true,
    apiKey: "rtc_demo_playground_20251027.ABC123...",  // ✅ Mesma chave!
    rateLimit: {
      requestsPerDay: 4,
      requestsPerMinute: 1
    },
    allowedApis: ["cep", "cnpj"]
  }
}
```

---

## 🎯 TESTE 2: Testar no Playground

### Passo 1: Abrir Playground
1. Abrir: http://localhost:3001/playground
2. **Verificar:**
   - ✅ "Protegido por rate limiting e fingerprinting"
   - ✅ Apenas CEP e CNPJ aparecem (sem GEO!)
   - ✅ CEP selecionado automaticamente

### Passo 2: Fazer Requests
1. **Request 1:** CEP 01310-100 → Clicar "Testar API"
   - ✅ Deve funcionar
   - ✅ Retornar dados do CEP
   
2. **Request 2:** CNPJ 00000000000191 → Selecionar CNPJ → Testar
   - ✅ Deve funcionar
   - ✅ Retornar dados da empresa

### Passo 3: Verificar Logs do Backend
```bash
docker logs build-api-1 --tail 50 | grep "PLAYGROUND SECURITY"
```

**Deve mostrar:**
```
🔒 [PLAYGROUND SECURITY] IP: 172.20.0.1 | Path: /public/cep/01310100
🔒 [PLAYGROUND SECURITY] API Key Demo detectada: rtc_demo_playground_...
📊 [PLAYGROUND SECURITY] Rate Limits Configurados:
   - Requests/Dia: 4
   - Requests/Min: 1
📊 [PLAYGROUND SECURITY] IP: 172.20.0.1 | Count: 0/4 (dia) | 0/1 (min)
✅ [PLAYGROUND SECURITY] Request permitida para IP 172.20.0.1
```

---

## 🎯 TESTE 3: Testar Rate Limiting

### Passo 1: Fazer 4 Requests Rapidamente
1. CEP 01310-100 → Testar
2. **Aguardar 3 segundos** (throttling)
3. CEP 88015-100 → Testar
4. **Aguardar 3 segundos**
5. CEP 20040-020 → Testar
6. **Aguardar 3 segundos**
7. CEP 30140-071 → Testar

**Todas devem funcionar!** ✅

### Passo 2: Fazer 5ª Request (Deve Bloquear)
1. CEP 40010-000 → Testar
2. **Deve mostrar erro:**
   ```json
   {
     "type": "https://retech-core/errors/rate-limit-exceeded",
     "title": "Rate Limit Exceeded",
     "status": 429,
     "detail": "Limite de 4 requests por dia por IP excedido. Tente novamente amanhã."
   }
   ```

### Passo 3: Verificar Headers
Abrir DevTools → Network → Selecionar request bloqueada:

**Headers de resposta devem ter:**
```
X-RateLimit-Limit-Day: 4
X-RateLimit-Remaining-Day: 0
X-RateLimit-Reset-Day: 1730073600
```

---

## 🎯 TESTE 4: Rotacionar API Key

### Passo 1: Voltar em Settings
1. Abrir: http://localhost:3001/admin/settings
2. **Verificar:**
   - ✅ API Key preenchida
   - ✅ Botão "🔄 Rotacionar" aparece

### Passo 2: Rotacionar
1. **Clicar:** Botão "🔄 Rotacionar"
2. **Aguardar:** "Rotacionando..."
3. **Verificar toasts:**
   - "API Key rotacionada! Scopes: cep, cnpj"
   - "Chave antiga desativada automaticamente"
4. **Verificar campo:** API Key mudou

### Passo 3: Validar Rotação no MongoDB
```bash
docker exec build-mongo-1 mongosh retech_core --eval '
db.api_keys.find({ownerId: "playground-public"}).sort({createdAt: -1}).limit(2).pretty()
'
```

**Deve mostrar:**
```javascript
// Chave NOVA (revoked: false)
{
  keyId: "rtc_demo_playground_20251027.XYZ789...",
  revoked: false,  // ✅ Ativa
  ...
}

// Chave ANTIGA (revoked: true)
{
  keyId: "rtc_demo_playground_20251027.ABC123...",
  revoked: true,  // ✅ Desativada
  ...
}
```

---

## 🎯 TESTE 5: Mudar Scopes e Regenerar

### Passo 1: Alterar APIs Permitidas
1. Em /admin/settings
2. **Desmarcar:** CNPJ
3. **Deixar apenas:** CEP ✓
4. **Clicar:** Salvar

### Passo 2: Rotacionar API Key
1. **Clicar:** "🔄 Rotacionar"
2. **Verificar toast:** "Scopes: cep" (só CEP!)

### Passo 3: Testar no Playground
1. Abrir: http://localhost:3001/playground
2. **Verificar:**
   - ✅ Apenas CEP aparece
   - ❌ CNPJ não aparece mais
   - ❌ GEO não aparece

### Passo 4: Testar CNPJ (Deve Bloquear)
Como só CEP tem scope, tentar usar CNPJ diretamente:

```bash
curl -X GET http://localhost:8080/public/cnpj/00000000000191 \
  -H "X-API-Key: $(docker exec build-mongo-1 mongosh retech_core --quiet --eval 'db.system_settings.findOne({},{playground:1}).playground.apiKey')"
```

**Deve retornar:**
```json
{
  "type": "https://retech-core/errors/insufficient-scope",
  "title": "Insufficient Scope",
  "status": 403,
  "detail": "Esta API Key não tem permissão para acessar cnpj"
}
```

---

## 🎯 TESTE 6: Throttling (2 segundos)

### Passo 1: Requests Rápidas
1. CEP 01310-100 → Testar
2. **Imediatamente** (sem esperar) → Testar novamente

**Deve mostrar:**
```json
{
  "type": "https://retech-core/errors/rate-limit-exceeded",
  "title": "Rate Limit Exceeded (Throttling)",
  "status": 429,
  "detail": "Aguarde 2 segundos antes de fazer outra requisição."
}
```

---

## ✅ CHECKLIST COMPLETO

- [ ] Backend rodando (http://localhost:8080/health)
- [ ] Frontend rodando (http://localhost:3001)
- [ ] Settings carregam sem API Key
- [ ] Botão "Gerar Nova" aparece quando vazio
- [ ] Clicar "Gerar Nova" cria API Key no MongoDB
- [ ] API Key aparece no campo após gerar
- [ ] Scopes corretos (baseados em checkboxes)
- [ ] Playground mostra apenas APIs permitidas
- [ ] Rate limiting funciona (4/dia, 1/min)
- [ ] Throttling funciona (2 segundos)
- [ ] Botão "Rotacionar" aparece quando tem API Key
- [ ] Rotacionar desativa antiga e cria nova
- [ ] Mudar scopes e rotacionar atualiza permissões

---

## 🐛 SE DER ERRO

**Logs do backend:**
```bash
docker logs build-api-1 --tail 100 -f
```

**Logs do frontend:**
```bash
tail -f /tmp/frontend-test.log
```

**MongoDB:**
```bash
docker exec -it build-mongo-1 mongosh retech_core
```

---

**PRONTO PARA TESTAR! 🚀**

Abra http://localhost:3001/admin/settings e me diga o que aparece!

