# 🧪 TESTE LOCAL - ROTEIRO COMPLETO

## ✅ AMBIENTE CONFIGURADO

- ✅ **Backend:** http://localhost:8080 (Docker Compose)
- ✅ **Frontend:** http://localhost:3001 (yarn dev)
- ✅ **MongoDB:** Configurado com playground vazio
- ✅ **CORS:** Permitindo localhost:3001

---

## 🚨 **CORREÇÕES APLICADAS (NÃO COMMITADAS):**

### **Backend:**
1. ✅ Rotas `/public/*` agora **REQUEREM API Key**
2. ✅ Endpoints `/admin/playground/apikey/generate`
3. ✅ Endpoints `/admin/playground/apikey/rotate`
4. ✅ Scopes baseados em `allowedAPIs`

### **Frontend:**
1. ✅ Sem defaults hardcoded (`rtc_demo_playground_2024`)
2. ✅ API Key vazia por padrão
3. ✅ Botões Gerar Nova / Rotacionar
4. ✅ **Auto-rotaciona ao mudar checkboxes de APIs**

---

## 🧪 **TESTE 1: Admin Settings (Inicial)**

### Abrir:
```
http://localhost:3001/admin/settings
```

### Deve Mostrar:
- ✅ Playground: **OFF** (desabilitado)
- ✅ API Key Demo: **VAZIO** (placeholder: "Clique em 'Gerar Nova'")
- ✅ Botão: **"🔑 Gerar Nova"**
- ✅ Requests/Dia: vazio ou 4
- ✅ Requests/Min: vazio ou 1
- ✅ APIs: Nenhuma marcada ou CEP/CNPJ marcados

---

## 🧪 **TESTE 2: Playground SEM API Key**

### Passo 1: Abrir Playground
```
http://localhost:3001/playground
```

### Deve Mostrar:
- ⚠️ Playground habilitado MAS sem APIs
- ⚠️ "Playground sem APIs Configuradas"
- ⚠️ Ou tela vazia/erro de API Key

**MOTIVO:** Playground está OFF ou sem API Key válida

---

## 🧪 **TESTE 3: Habilitar e Gerar API Key**

### Passo 1: Voltar em Settings
1. Ligar: **Playground ON**
2. Requests/Dia: **4**
3. Requests/Min: **1**
4. Marcar: **CEP ✓, CNPJ ✓** (sem GEO)

### Passo 2: Gerar API Key
1. **Clicar:** "🔑 Gerar Nova"
2. **Aguardar:** Toast "API Key gerada com sucesso! Scopes: cep, cnpj"
3. **Verificar:** Campo preenchido com `rtc_demo_playground_20251027.ABC...`
4. **Botão muda:** "🔄 Rotacionar" aparece

### Passo 3: Salvar Configurações
1. **Clicar:** "💾 Salvar Configurações"
2. **Aguardar:** Toast "Configurações salvas com sucesso!"

---

## 🧪 **TESTE 4: Playground COM API Key**

### Passo 1: Abrir Playground
```
http://localhost:3001/playground
```

### Deve Mostrar:
- ✅ "Protegido por rate limiting e fingerprinting"
- ✅ **Apenas CEP e CNPJ** (sem GEO)
- ✅ CEP selecionado automaticamente

### Passo 2: Testar CEP
1. CEP: **01310-100**
2. **Clicar:** "Testar API"
3. **Deve retornar:** Dados do CEP ✅

### Passo 3: Testar CNPJ
1. **Selecionar:** CNPJ
2. CNPJ: **00000000000191**
3. **Clicar:** "Testar API"
4. **Deve retornar:** Dados da empresa ✅

### Passo 4: Tentar GEO (Deve Falhar)
- ❌ GEO não aparece (correto!)
- Se tentar via curl deve dar erro de scope

---

## 🧪 **TESTE 5: Auto-Rotação ao Mudar Scopes**

### Passo 1: Voltar em Settings
1. **Marcar:** GEO ✓ (agora CEP ✓, CNPJ ✓, GEO ✓)

### Deve Acontecer AUTOMATICAMENTE:
1. ⏳ Toast: "Atualizando scopes da API Key..."
2. 🔄 API Key é rotacionada automaticamente
3. ✅ Toast: "Scopes atualizados: cep, cnpj, geo"
4. 🔑 Campo API Key muda (nova chave)
5. 📝 Sem precisar clicar "Salvar"!

### Passo 2: Verificar Playground
1. Abrir: http://localhost:3001/playground
2. **Deve mostrar:** CEP, CNPJ, **E GEO** ✅

---

## 🧪 **TESTE 6: Rate Limiting por IP**

### Configuração Atual:
- Requests/Dia: **4**
- Requests/Min: **1**

### Teste:
1. CEP 01310-100 → Testar → ✅ Sucesso (1/4)
2. **Aguardar 3 segundos** (throttling)
3. CEP 88015-100 → Testar → ✅ Sucesso (2/4)
4. **Aguardar 3 segundos**
5. CEP 20040-020 → Testar → ✅ Sucesso (3/4)
6. **Aguardar 3 segundos**
7. CEP 30140-071 → Testar → ✅ Sucesso (4/4)
8. **Aguardar 3 segundos**
9. CEP 40010-000 → Testar → ❌ **BLOQUEADO!**

### Mensagem Esperada:
```json
{
  "type": "https://retech-core/errors/rate-limit-exceeded",
  "title": "Rate Limit Exceeded",
  "status": 429,
  "detail": "Limite de 4 requests por dia por IP excedido. Tente novamente amanhã."
}
```

### Logs do Backend:
```bash
docker logs build-api-1 --tail 20 | grep "PLAYGROUND SECURITY"
```

**Deve mostrar:**
```
📊 [PLAYGROUND SECURITY] Rate Limits Configurados:
   - Requests/Dia: 4
   - Requests/Min: 1
📊 [PLAYGROUND SECURITY] IP: 172.20.0.1 | Count: 3/4 (dia) | 0/1 (min)
✅ [PLAYGROUND SECURITY] Request permitida para IP 172.20.0.1
...
🚫 [PLAYGROUND SECURITY] Limite diário por IP excedido: 172.20.0.1 (4 >= 4)
```

---

## 🧪 **TESTE 7: Throttling (2 segundos)**

### Teste:
1. CEP 01310-100 → Testar
2. **Imediatamente** (sem esperar) → Testar novamente

### Deve Bloquear:
```json
{
  "type": "https://retech-core/errors/rate-limit-exceeded",
  "title": "Rate Limit Exceeded (Throttling)",
  "status": 429,
  "detail": "Aguarde 2 segundos antes de fazer outra requisição."
}
```

---

## 🧪 **TESTE 8: Request SEM API Key (Deve Bloquear)**

### Teste via curl:
```bash
curl -X GET http://localhost:8080/public/cep/01310100
```

### Deve Retornar:
```json
{
  "type": "https://retech-core/errors/unauthorized",
  "title": "Unauthorized",
  "status": 401,
  "detail": "API Key não fornecida"
}
```

**✅ CORRETO!** Agora rotas públicas exigem API Key!

---

## 📊 **VALIDAÇÕES NO MONGODB**

### Ver API Keys Demo:
```bash
docker exec build-mongo-1 mongosh retech_core --eval '
db.api_keys.find({ownerId: "playground-public"}).pretty()
'
```

### Ver Settings do Playground:
```bash
docker exec build-mongo-1 mongosh retech_core --eval '
db.system_settings.findOne({}, {playground: 1})
'
```

### Ver Rate Limits por IP:
```bash
docker exec build-mongo-1 mongosh retech_core --eval '
db.playground_rate_limits.find().sort({updatedAt: -1}).limit(5).pretty()
'
```

---

## ✅ **CHECKLIST FINAL**

- [ ] Backend rodando sem erros
- [ ] Frontend rodando sem erros
- [ ] Admin/settings carrega com API Key vazia
- [ ] Botão "Gerar Nova" funciona
- [ ] API Key criada no MongoDB
- [ ] Scopes corretos (baseados em checkboxes)
- [ ] Playground mostra apenas APIs permitidas
- [ ] Requests funcionam COM API Key
- [ ] Requests falham SEM API Key
- [ ] Rate limiting funciona (4/dia, 1/min)
- [ ] Throttling funciona (2 segundos)
- [ ] Auto-rotação ao mudar checkboxes
- [ ] Botão "Rotacionar" funciona manualmente

---

## 🎯 **COMEÇE POR AQUI:**

**1. Abrir:** http://localhost:3001/admin/settings  
**2. Scroll até:** Playground Público  
**3. Me diga:** O que você vê?

---

**PRONTO PARA TESTAR! 🚀**

