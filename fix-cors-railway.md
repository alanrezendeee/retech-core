# 🔧 Corrigir CORS no Railway - Guia Rápido

## 🚨 **Problema Identificado:**

Produção não está retornando headers CORS mesmo com configuração correta.

---

## ✅ **Solução 1: Verificar/Corrigir MongoDB**

### **Passo 1: Conectar no MongoDB do Railway**

```bash
# Obtenha a connection string do Railway:
# Railway → retech-core → MongoDB → Connect

# Exemplo:
mongodb://mongo:PASSWORD@HOST:PORT
```

### **Passo 2: Verificar Configuração Atual**

```javascript
// No mongosh ou MongoDB Compass:
use retech_core

db.system_settings.findOne({}, {cors: 1})
```

**Resultado atual (provavelmente):**
```javascript
{
  cors: {
    enabled: true,
    allowedOrigins: [
      "https://core.theretech.com.br, http://localhost:3000"  // ❌ ERRADO!
    ]
  }
}
```

### **Passo 3: CORRIGIR Configuração**

```javascript
db.system_settings.updateOne(
  {},
  {
    $set: {
      "cors.enabled": true,
      "cors.allowedOrigins": [
        "https://core.theretech.com.br",
        "http://localhost:3000",
        "http://localhost:3001",
        "http://localhost:3002",
        "http://localhost:3003"
      ]
    }
  }
)
```

**⚠️ IMPORTANTE:** Cada origin deve ser um **elemento separado** do array, NÃO uma string com vírgulas!

### **Passo 4: Verificar Correção**

```javascript
db.system_settings.findOne({}, {cors: 1})
```

**Resultado esperado:**
```javascript
{
  cors: {
    enabled: true,
    allowedOrigins: [
      "https://core.theretech.com.br",  // ✅ String separada
      "http://localhost:3000",          // ✅ String separada
      "http://localhost:3001",          // ✅ String separada
      "http://localhost:3002",
      "http://localhost:3003"
    ]
  }
}
```

---

## ✅ **Solução 2: Verificar Deploy**

### **Verificar versão do código:**

```bash
# Ver último commit deployado:
curl -s https://api-core.theretech.com.br/version | jq
```

**Deveria retornar algo como:**
```json
{
  "version": "1.0.0",
  "commit": "153ba52",  // ← Deve ser o último commit
  "buildTime": "..."
}
```

### **Se versão estiver desatualizada:**

1. Railway → retech-core → Deployments
2. Verificar se último deploy foi bem-sucedido
3. Se não, fazer redeploy manual

---

## ✅ **Solução 3: Limpar Cache do Settings**

Se o código estiver atualizado mas ainda não funciona, pode ser cache:

### **Opção A: Reiniciar Aplicação**

```bash
# Railway → retech-core → Settings → Restart
```

### **Opção B: Forçar Reload (via código)**

Se você tiver acesso SSH ou Railway CLI:

```bash
# Fazer uma request que force reload do cache
curl -X PUT https://api-core.theretech.com.br/admin/settings \
  -H "Authorization: Bearer SEU_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{...}' # Qualquer update força reload
```

---

## 🧪 **Testar Correção:**

### **Teste 1: curl com Origin**

```bash
curl -v -H "Origin: http://localhost:3000" \
  https://api-core.theretech.com.br/health 2>&1 | grep "Access-Control"
```

**✅ Deve retornar:**
```
< Access-Control-Allow-Origin: http://localhost:3000
< Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
< Access-Control-Allow-Headers: ...
```

### **Teste 2: OPTIONS Preflight**

```bash
curl -v -X OPTIONS \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  https://api-core.theretech.com.br/auth/login 2>&1 | grep "Access-Control"
```

**✅ Deve retornar headers CORS**

### **Teste 3: Do Frontend**

```javascript
// No console do browser (localhost:3000):
fetch('https://api-core.theretech.com.br/health', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(r => r.json())
.then(data => console.log('✅ CORS funcionando!', data))
.catch(err => console.error('❌ CORS bloqueado:', err));
```

---

## 🎯 **Checklist de Correção:**

- [ ] Conectei no MongoDB do Railway
- [ ] Verifiquei `cors.allowedOrigins` no MongoDB
- [ ] Corrigi array de origins (cada um separado, não string com vírgulas)
- [ ] Verifiquei que `cors.enabled` está `true`
- [ ] Reiniciei aplicação no Railway
- [ ] Testei com curl e vi headers CORS
- [ ] Testei do frontend e funcionou

---

## 🐛 **Se Ainda Não Funcionar:**

### **Debug Avançado:**

1. **Ver logs do Railway:**
   ```
   Railway → retech-core → Deployments → View Logs
   ```

2. **Procurar por:**
   ```
   [CORS] GET /health (Origin: http://localhost:3000)
   [CORS] Settings: CORS.Enabled=true, AllowedOrigins=[...]
   [CORS] ✅ Origin permitido - adicionando headers
   ```

3. **Se ver:**
   ```
   [CORS] ❌ Origin não está na lista permitida
   ```
   
   → Configuração do MongoDB está errada

4. **Se NÃO ver logs de CORS:**
   
   → Deploy não rodou ou código antigo está rodando

---

## 📝 **Formato Correto do allowedOrigins:**

### **❌ ERRADO:**

```json
{
  "allowedOrigins": [
    "https://site1.com, http://localhost:3000, http://localhost:3001"
  ]
}
```

**Problema:** É um array com **1 elemento** (string com vírgulas)

### **✅ CORRETO:**

```json
{
  "allowedOrigins": [
    "https://site1.com",
    "http://localhost:3000",
    "http://localhost:3001"
  ]
}
```

**Correto:** É um array com **3 elementos** (strings separadas)

---

## 💡 **Dica: Verificar no Admin UI**

Se o campo de texto no `/admin/settings` aceita "separado por vírgulas", pode ser que o **frontend está salvando errado**!

Veja o código do frontend que envia para o backend e verifique se está fazendo:

```javascript
// ❌ ERRADO:
const origins = "https://site1.com, http://localhost:3000";  // String
cors: {
  enabled: true,
  allowedOrigins: [origins]  // Array com 1 elemento
}

// ✅ CORRETO:
const origins = "https://site1.com, http://localhost:3000";
cors: {
  enabled: true,
  allowedOrigins: origins.split(',').map(o => o.trim())  // Split e trim
}
```

---

**Após seguir este guia, CORS deve funcionar! Se não funcionar, compartilhe:**
1. Output do `db.system_settings.findOne({}, {cors: 1})`
2. Logs do Railway
3. Output do `curl -v -H "Origin: ..." /health`

