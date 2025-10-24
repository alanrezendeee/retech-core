# 🎮 Sistema de Playground Configurável

## 📋 Resumo

O playground é **100% gerenciável** via `admin/settings`, permitindo:
- ✅ Habilitar/Desabilitar com toggle
- ✅ Trocar API Key sem código
- ✅ Ajustar rate limits
- ✅ Escolher APIs disponíveis

---

## 🔄 Fluxo Completo

### **1. Admin Configura (admin/settings)**

```
Admin → Login → /admin/settings → Seção "Playground Público"

Campos Disponíveis:
┌─────────────────────────────────────────┐
│ ☑ Habilitar Playground                  │
│                                         │
│ API Key Demo:                           │
│ [rtc_demo_playground_2024]              │
│                                         │
│ Rate Limits:                            │
│ Requests/Dia:    [100]                  │
│ Requests/Minuto: [10]                   │
│                                         │
│ APIs Disponíveis:                       │
│ ☑ CEP  ☑ CNPJ  ☑ GEO  □ IBGE          │
└─────────────────────────────────────────┘

[Salvar Configurações] ✅
```

### **2. Backend Salva no MongoDB**

```json
// Collection: system_settings
{
  "_id": "system-settings-singleton",
  "playground": {
    "enabled": true,
    "apiKey": "rtc_demo_playground_2024",
    "rateLimit": {
      "RequestsPerDay": 100,
      "RequestsPerMinute": 10
    },
    "allowedApis": ["cep", "cnpj", "geo"]
  },
  "updatedAt": "2025-10-24T..."
}
```

### **3. Usuário Acessa /playground**

```
1. Página carrega
   ↓
2. useEffect() dispara
   ↓
3. Faz request: GET /public/playground/status
   ↓
4. Backend lê settings do MongoDB
   ↓
5. Backend retorna:
   {
     "enabled": true,
     "apiKey": "rtc_demo_playground_2024",
     "allowedApis": ["cep", "cnpj", "geo"]
   }
   ↓
6. Frontend atualiza states:
   - setIsPlaygroundEnabled(true)
   - setDemoApiKey("rtc_demo_playground_2024")
   ↓
7. Playground renderiza normalmente
   ↓
8. Usuário testa API
   ↓
9. Request usa header: X-API-Key: rtc_demo_playground_2024
   ↓
10. ✅ Funciona!
```

---

## 🛡️ Graceful Degradation

### **Cenário 1: Settings Vazio (Primeira Vez)**

Se o campo `playground` não existe no MongoDB:

```go
// Backend (playground.go)
if apiKey == "" {
    apiKey = "rtc_demo_playground_2024"
    enabled = true  // ✅ Assume habilitado
}
```

**Resultado:** Playground funciona out-of-the-box! 🎉

---

### **Cenário 2: API Falhar**

Se `/public/playground/status` falhar:

```tsx
// Frontend (playground/page.tsx)
catch (error) {
  setIsPlaygroundEnabled(true); // ✅ Assume habilitado
}
```

**Resultado:** Playground continua funcionando! 🎉

---

## 🚫 Desabilitar Playground

### **Admin:**
```
/admin/settings → Toggle OFF "Habilitar Playground" → Salvar
```

### **Backend:**
```json
{
  "playground": {
    "enabled": false,  // ❌ Desabilitado
    ...
  }
}
```

### **Usuário vê:**
```
┌─────────────────────────────────────────┐
│        [⚠️ ícone laranja]                │
│                                         │
│   🚫 Playground Indisponível            │
│                                         │
│   O playground está temporariamente     │
│   desabilitado.                         │
│                                         │
│   [🏠 Voltar para Home]                 │
│   [Login Admin]                         │
└─────────────────────────────────────────┘
```

---

## 🔑 Trocar API Key (Abuso Detectado)

### **Problema:**
```
Alguém está abusando da API Key demo!
Rate limits estão sendo atingidos.
```

### **Solução Rápida:**

1. **Admin:**
   ```
   /admin/settings → API Key Demo: [nova_chave_segura_xyz]
   → Salvar (2 segundos)
   ```

2. **Backend:**
   ```json
   {
     "playground": {
       "apiKey": "nova_chave_segura_xyz"  // ✅ Atualizado
     }
   }
   ```

3. **Próximo usuário:**
   ```
   Acessa /playground
   → Carrega nova chave: nova_chave_segura_xyz
   → Usa nas requisições
   → ✅ Funciona!
   ```

4. **Abusador:**
   ```
   Ainda usa chave antiga: rtc_demo_playground_2024
   → Backend rejeita: 401 Unauthorized
   → ❌ Bloqueado!
   ```

---

## 📊 Endpoints

### **GET /public/playground/status**

**Público** (sem autenticação)

**Response (enabled=true):**
```json
{
  "enabled": true,
  "message": "Playground disponível",
  "apiKey": "rtc_demo_playground_2024",
  "allowedApis": ["cep", "cnpj", "geo"]
}
```

**Response (enabled=false):**
```json
{
  "enabled": false,
  "message": "Playground temporariamente indisponível. Entre em contato para mais informações."
}
```

---

## 🧪 Testes

### **Teste 1: Verificar Status**

```bash
curl https://api-core.theretech.com.br/public/playground/status | jq
```

**Esperado:**
```json
{
  "enabled": true,
  "apiKey": "rtc_demo_playground_2024",
  "allowedApis": ["cep", "cnpj", "geo"]
}
```

---

### **Teste 2: Desabilitar Playground**

1. Login admin
2. `/admin/settings` → Toggle OFF
3. Salvar
4. Testar:
   ```bash
   curl https://api-core.theretech.com.br/public/playground/status | jq
   ```

**Esperado:**
```json
{
  "enabled": false,
  "message": "Playground temporariamente indisponível..."
}
```

5. Abrir `/playground` em aba anônima
6. ✅ Ver mensagem de indisponível

---

### **Teste 3: Trocar API Key**

1. Login admin
2. `/admin/settings` → API Key: `teste_nova_chave`
3. Salvar
4. Abrir DevTools Console no `/playground`
5. Atualizar página
6. ✅ Ver log: `✅ API Key do playground carregada: teste_nova_chave`
7. Testar uma API (ex: CEP)
8. Ver na aba Network: `X-API-Key: teste_nova_chave`

---

## 🔒 Segurança

### **Rate Limits:**

```json
{
  "playground": {
    "rateLimit": {
      "RequestsPerDay": 100,      // Limite global diário
      "RequestsPerMinute": 10     // Limite por minuto
    }
  }
}
```

**Aplicação:**
- Rate limit por **tenantID** (API Key)
- Rate limit adicional por **IP** (prevenir abuso)

---

### **API Key Exposta?**

**Não é problema!** 🛡️

- API Key é **pública por design** (para playground)
- Rate limits **agressivos** (100 req/dia, 10 req/min)
- Admin pode **trocar instantaneamente**
- Scopes **restritos** (apenas cep, cnpj, geo)

---

## 📂 Arquivos Modificados

### **Backend:**
- `internal/http/handlers/playground.go` (novo)
- `internal/http/router.go` (rota adicionada)
- `internal/domain/settings.go` (PlaygroundConfig)
- `internal/storage/settings_repo.go` (save playground)
- `internal/bootstrap/demo_apikey.go` (auto-create)

### **Frontend:**
- `app/admin/settings/page.tsx` (UI de configuração)
- `app/playground/page.tsx` (verifica status e usa API Key dinâmica)

---

## ✅ Status

| Feature | Backend | Frontend | Status |
|---------|---------|----------|--------|
| Toggle ON/OFF | ✅ | ✅ | **Completo** |
| API Key dinâmica | ✅ | ✅ | **Completo** |
| Rate limits editáveis | ✅ | ✅ | **Completo** |
| Seletor de APIs | ✅ | ✅ | **Completo** |
| Página indisponível | ✅ | ✅ | **Completo** |
| Graceful degradation | ✅ | ✅ | **Completo** |

---

## 🎉 Conclusão

**Tudo funciona via interface! Zero código!** 🚀

Admin tem controle total:
- 🔄 Ativar/Desativar
- 🔑 Trocar API Key
- 🎚️ Ajustar limites
- 🔧 Escolher APIs

**Mudanças aplicam instantaneamente!** ⚡

