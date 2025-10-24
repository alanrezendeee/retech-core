# 🧪 Guia de Teste - Playground Configurável

## 📋 Pré-requisitos

- ✅ Backend deployado com últimas mudanças
- ✅ Frontend deployado com últimas mudanças
- ✅ Acesso admin ao `/admin/settings`
- ✅ DevTools aberto (Console + Network)

---

## 🎯 Teste 1: Verificar Status Atual

### **1.1 Via cURL:**
```bash
curl -s https://api-core.theretech.com.br/public/playground/status | jq
```

**Esperado:**
```json
{
  "enabled": true,
  "message": "Playground disponível",
  "apiKey": "rtc_demo_playground_2024",
  "allowedApis": ["cep", "cnpj", "geo"]
}
```

### **1.2 Via Browser:**
1. Abrir: `https://core.theretech.com.br/playground`
2. Abrir DevTools → Console
3. Ver log:
   ```
   ✅ API Key do playground carregada: rtc_demo_playground_2024
   🎮 Playground status: {enabled: true, apiKey: "...", allowedApis: Array(3)}
   ```

---

## 🚫 Teste 2: Desabilitar Playground

### **2.1 Admin desabilita:**
```
1. Login: https://core.theretech.com.br/admin/login
2. Ir: /admin/settings
3. Rolar até "Playground Público"
4. Toggle OFF "Habilitar Playground"
5. Clicar "Salvar Configurações"
6. ✅ Ver toast: "Configurações salvas com sucesso!"
```

### **2.2 Verificar no backend:**
```bash
curl -s https://api-core.theretech.com.br/public/playground/status | jq
```

**Esperado:**
```json
{
  "enabled": false,
  "message": "Playground temporariamente indisponível. Entre em contato para mais informações."
}
```

### **2.3 Verificar no Railway Logs:**
Procurar por:
```
🎮 Playground config recebido: enabled=false, ...
```

### **2.4 Usuário tenta acessar:**
```
1. Abrir nova aba anônima (Ctrl+Shift+N)
2. Ir: https://core.theretech.com.br/playground
3. Aguardar loading
4. ✅ Ver página de indisponível:
   ┌─────────────────────────────────────────┐
   │        [⚠️ ícone laranja]                │
   │   🚫 Playground Indisponível            │
   │   O playground está temporariamente     │
   │   desabilitado.                         │
   │   [🏠 Voltar para Home]                 │
   └─────────────────────────────────────────┘
```

### **2.5 Verificar no Console:**
```
🎮 Playground status: {enabled: false}
```

---

## ✅ Teste 3: Reabilitar Playground

### **3.1 Admin reabilita:**
```
1. Voltar ao /admin/settings
2. Toggle ON "Habilitar Playground"
3. Salvar
4. ✅ Toast de sucesso
```

### **3.2 Verificar backend:**
```bash
curl -s https://api-core.theretech.com.br/public/playground/status | jq
```

**Esperado:**
```json
{
  "enabled": true,
  "message": "Playground disponível",
  "apiKey": "rtc_demo_playground_2024",
  "allowedApis": ["cep", "cnpj", "geo"]
}
```

### **3.3 Usuário acessa novamente:**
```
1. Atualizar página do playground (F5)
2. ✅ Ver playground funcionando normalmente
3. Console log:
   🎮 Playground status: {enabled: true, ...}
```

---

## 🔑 Teste 4: Trocar API Key

### **4.1 Admin troca chave:**
```
1. /admin/settings
2. API Key Demo: [teste_nova_chave_123]
3. Salvar
```

### **4.2 Verificar backend:**
```bash
curl -s https://api-core.theretech.com.br/public/playground/status | jq
```

**Esperado:**
```json
{
  "enabled": true,
  "apiKey": "teste_nova_chave_123",  // ✅ Nova chave!
  ...
}
```

### **4.3 Usuário usa nova chave:**
```
1. Atualizar /playground (F5)
2. Console:
   ✅ API Key do playground carregada: teste_nova_chave_123
3. Testar uma API (ex: CEP 01310-100)
4. DevTools → Network → Selecionar request
5. Headers → Request Headers:
   X-API-Key: teste_nova_chave_123  // ✅ Nova chave!
6. ✅ Request funciona!
```

---

## 🎚️ Teste 5: Ajustar Rate Limits

### **5.1 Admin ajusta:**
```
1. /admin/settings
2. Requests por Dia: [50]
3. Requests por Minuto: [5]
4. Salvar
```

### **5.2 Verificar logs Railway:**
```
🎮 Playground config recebido: enabled=true, apiKey=..., reqPerDay=50, reqPerMin=5, ...
```

### **5.3 Testar limite:**
```
1. Fazer 6 requests rápidas no playground
2. 6ª request: ❌ 429 Too Many Requests
```

---

## 🔧 Teste 6: Selecionar APIs Disponíveis

### **6.1 Admin desabilita CNPJ:**
```
1. /admin/settings
2. APIs Disponíveis:
   ☑ CEP
   ☐ CNPJ  (desmarcar)
   ☑ GEO
3. Salvar
```

### **6.2 Verificar backend:**
```bash
curl -s https://api-core.theretech.com.br/public/playground/status | jq '.allowedApis'
```

**Esperado:**
```json
["cep", "geo"]  // ✅ Sem "cnpj"
```

### **6.3 Usuário tenta usar CNPJ:**
```
1. Playground → Tab CNPJ
2. Tentar consultar
3. ❌ 403 Forbidden ou esconder tab (futuro)
```

---

## 🐛 Troubleshooting

### **Problema: Admin salva mas nada muda**

**Verificar:**
1. Railway logs → Ver se recebeu:
   ```
   🎮 Playground config recebido: enabled=..., apiKey=..., ...
   ```

2. Se NÃO aparece → Problema no frontend (payload errado)
3. Se aparece → Problema no save do MongoDB

**Solução:**
```bash
# Ver settings salvos no MongoDB
mongo
> use retech
> db.system_settings.findOne()
```

---

### **Problema: Playground sempre habilitado mesmo desabilitando**

**Verificar:**
1. Browser cache:
   ```bash
   # Hard refresh
   Ctrl+Shift+R (Windows/Linux)
   Cmd+Shift+R (Mac)
   ```

2. Backend retorna enabled=false?
   ```bash
   curl https://api-core.theretech.com.br/public/playground/status
   ```

3. Console logs:
   ```
   🎮 Playground status: {enabled: ???}
   ```

**Solução:**
- Fix implementado: `cache: 'no-store'` no fetch

---

### **Problema: API Key antiga ainda funciona**

**Explicação:**
- Tenant demo usa a chave do settings
- Se chave mudar no settings, tenant demo é atualizado no próximo restart

**Forçar update:**
```bash
# Restart do serviço no Railway
# Ou aguardar próximo deploy
```

---

## ✅ Checklist de Teste Completo

### **Backend:**
- [ ] Endpoint `/public/playground/status` responde
- [ ] Retorna `enabled`, `apiKey`, `allowedApis`
- [ ] Muda quando admin salva settings
- [ ] Logs aparecem no Railway

### **Frontend (Admin):**
- [ ] Toggle ON/OFF funciona
- [ ] Input de API Key salva
- [ ] Rate limits salvam
- [ ] Checkboxes de APIs salvam
- [ ] Toast de sucesso aparece

### **Frontend (Playground):**
- [ ] Carrega status ao abrir
- [ ] Mostra "indisponível" se disabled
- [ ] Mostra playground se enabled
- [ ] Usa API Key do backend
- [ ] Console logs aparecem
- [ ] Não usa cache (sempre fresh)

### **Integração End-to-End:**
- [ ] Admin desabilita → Usuário vê indisponível
- [ ] Admin habilita → Usuário vê playground
- [ ] Admin troca chave → Usuário usa nova chave
- [ ] Admin ajusta limits → Limites aplicados
- [ ] Admin remove API → Usuário não acessa

---

## 📊 Status Esperado

```
✅ CORS dinâmico
✅ API Key gerenciável
✅ TTL Redis configurável
✅ Playground ON/OFF
✅ Rate limits editáveis
✅ Seletor de APIs
✅ Cache desabilitado (sempre fresh)
✅ Logs de debug
```

---

## 🎉 Sucesso!

Se TODOS os testes passarem:
- ✅ Sistema 100% funcional
- ✅ Admin tem controle total
- ✅ Mudanças aplicam instantaneamente
- ✅ Zero código necessário

