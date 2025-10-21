# 🔄 API Key Rotation - Documentação

**Data:** 2025-10-21  
**Feature:** Rotação de API Keys (Admin e Tenant)

---

## 🎯 **O que é Rotação de API Key?**

Rotação de API Key é uma **prática de segurança** onde uma chave existente é **revogada** e uma **nova chave é gerada automaticamente** com os mesmos privilégios.

### **Por que rotacionar?**

1. **🔒 Segurança preventiva:** Mesmo sem vazamento conhecido, é boa prática rotacionar periodicamente
2. **🚨 Resposta a incidentes:** Se uma key foi potencialmente comprometida
3. **👥 Mudança de equipe:** Quando desenvolvedores saem da empresa
4. **🔄 Renovação:** Evitar que keys expirem inesperadamente

### **O que acontece na rotação:**

```
Key Antiga: abc-123.SECRET-OLD
    ↓
1. Revoga key antiga (revoked: true)
2. Gera nova key (uuid.SECRET-NEW)
3. Mantém mesmo tenant, scopes, validade
    ↓
Key Nova: xyz-789.SECRET-NEW ✅
```

---

## 📡 **Endpoints Implementados**

### **1. Admin - Rotacionar API Key** (qualquer tenant)

```
POST /admin/apikeys/rotate
```

**Autenticação:** JWT Bearer Token (SUPER_ADMIN)

**Request Body:**
```json
{
  "keyId": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (201 Created):**
```json
{
  "api_key": "abc123-new-uuid.NEWSECRET32CHARS",
  "expiresAt": "2026-01-20T10:00:00Z",
  "message": "API key rotacionada com sucesso"
}
```

**Casos de erro:**
- `400`: keyId não fornecido
- `404`: API key não encontrada
- `500`: Erro ao revogar ou criar nova key

---

### **2. Tenant - Rotacionar Minha API Key**

```
POST /me/apikeys/:id/rotate
```

**Autenticação:** JWT Bearer Token (TENANT_USER)

**URL Params:**
- `:id` - keyId da API key a ser rotacionada

**Response (201 Created):**
```json
{
  "key": "abc123-new-uuid.NEWSECRET32CHARS",
  "expiresAt": "2026-01-20T10:00:00Z",
  "message": "API key rotacionada com sucesso"
}
```

**Casos de erro:**
- `401`: Token inválido ou ausente
- `403`: Key não pertence ao tenant
- `404`: API key não encontrada
- `500`: Erro ao revogar ou criar nova key

---

## 🔐 **Diferenças entre Admin e Tenant**

| Aspecto | Admin (`/admin/apikeys/rotate`) | Tenant (`/me/apikeys/:id/rotate`) |
|---------|----------------------------------|-----------------------------------|
| **Permissão** | SUPER_ADMIN | TENANT_USER |
| **Escopo** | Qualquer API key | Apenas keys do próprio tenant |
| **Body** | `{"keyId": "..."}` | Sem body (ID na URL) |
| **Ownership Check** | ❌ Não verifica | ✅ Verifica se pertence ao tenant |
| **Activity Log** | ✅ Sim (`apikey.rotated`) | ❌ Não (por enquanto) |

---

## 🧪 **Como Testar**

### **Teste 1: Admin rotaciona key**

```bash
# 1. Login como admin
curl -X POST https://api-core.theretech.com.br/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@retech.com.br",
    "password": "admin123456"
  }'

# Salve o accessToken

# 2. Listar todas as API keys
curl https://api-core.theretech.com.br/admin/apikeys \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"

# Copie um keyId

# 3. Rotacionar
curl -X POST https://api-core.theretech.com.br/admin/apikeys/rotate \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "keyId": "550e8400-e29b-41d4-a716-446655440000"
  }'

# Response:
{
  "api_key": "abc123-new.NEWSECRET",
  "expiresAt": "2026-01-20T10:00:00Z",
  "message": "API key rotacionada com sucesso"
}

# 4. Verificar que a antiga foi revogada
curl https://api-core.theretech.com.br/geo/ufs \
  -H "X-API-Key: 550e8400-e29b-41d4-a716-446655440000.OLDSECRET"

# Response: {"error": "unknown api key"} ✅

# 5. Testar com a nova
curl https://api-core.theretech.com.br/geo/ufs \
  -H "X-API-Key: abc123-new.NEWSECRET"

# Response: [lista de UFs] ✅
```

---

### **Teste 2: Tenant rotaciona própria key**

```bash
# 1. Registrar como tenant
curl -X POST https://api-core.theretech.com.br/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenantName": "Startup Teste",
    "tenantEmail": "contato@startup.com",
    "company": "Startup LTDA",
    "purpose": "Testes",
    "userName": "João Dev",
    "userEmail": "joao@startup.com",
    "userPassword": "senha123456"
  }'

# Salve o accessToken

# 2. Criar API Key
curl -X POST https://api-core.theretech.com.br/me/apikeys \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Minha Key"}'

# Response:
{
  "key": "old-key-id.OLDSECRET",
  "expiresAt": "2026-01-20T10:00:00Z",
  "name": "Minha Key"
}

# 3. Listar minhas keys (pegar o keyId)
curl https://api-core.theretech.com.br/me/apikeys \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"

# Copie o keyId

# 4. Rotacionar
curl -X POST https://api-core.theretech.com.br/me/apikeys/old-key-id/rotate \
  -H "Authorization: Bearer SEU_ACCESS_TOKEN"

# Response:
{
  "key": "new-key-id.NEWSECRET",
  "expiresAt": "2026-01-20T10:00:00Z",
  "message": "API key rotacionada com sucesso"
}

# 5. Verificar que a antiga não funciona
curl https://api-core.theretech.com.br/geo/ufs \
  -H "X-API-Key: old-key-id.OLDSECRET"

# Response: {"error": "unknown api key"} ✅

# 6. Testar com a nova
curl https://api-core.theretech.com.br/geo/ufs \
  -H "X-API-Key: new-key-id.NEWSECRET"

# Response: [lista de UFs] ✅
```

---

### **Teste 3: Tenant tenta rotacionar key de outro tenant (deve falhar)**

```bash
# 1. Login como Tenant A
curl -X POST https://api-core.theretech.com.br/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "tenantA@example.com",
    "password": "senha123"
  }'

# 2. Tentar rotacionar key do Tenant B
curl -X POST https://api-core.theretech.com.br/me/apikeys/keyId-do-tenant-B/rotate \
  -H "Authorization: Bearer TOKEN_TENANT_A"

# Response (403 Forbidden):
{
  "type": "https://retech-core/errors/forbidden",
  "title": "Acesso negado",
  "status": 403,
  "detail": "Você não tem permissão para rotacionar esta API key"
}

# ✅ Segurança funcionando!
```

---

## 📊 **Activity Logs**

### **Admin Rotation:**

Quando um admin rotaciona uma key, é criado um log de atividade:

```json
{
  "timestamp": "2025-10-21T20:00:00Z",
  "type": "apikey.rotated",
  "action": "update",
  "actor": {
    "userId": "admin-123",
    "email": "admin@retech.com.br",
    "name": "Admin User",
    "role": "SUPER_ADMIN"
  },
  "resource": {
    "type": "apikey",
    "id": "new-key-id",
    "name": "API Key de Startup (rotacionada)"
  },
  "metadata": {
    "tenantId": "tenant-456",
    "oldKeyId": "old-key-id",
    "newKeyId": "new-key-id",
    "expiresAt": "2026-01-20T10:00:00Z"
  },
  "ip": "192.168.1.100",
  "userAgent": "curl/7.x"
}
```

**Visualizar logs:**
```bash
curl https://api-core.theretech.com.br/admin/activity/type/apikey.rotated \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## 🔒 **Segurança**

### **Validações implementadas:**

#### **Admin (`/admin/apikeys/rotate`):**
1. ✅ Requer autenticação JWT
2. ✅ Requer role SUPER_ADMIN
3. ✅ Valida que keyId existe
4. ✅ Revoga antiga antes de criar nova
5. ✅ Gera nova key com HMAC-SHA256 seguro
6. ✅ Logs de auditoria

#### **Tenant (`/me/apikeys/:id/rotate`):**
1. ✅ Requer autenticação JWT
2. ✅ Requer role TENANT_USER
3. ✅ Valida que keyId existe
4. ✅ **Verifica ownership** (key pertence ao tenant)
5. ✅ Retorna 403 se não for o dono
6. ✅ Revoga antiga antes de criar nova
7. ✅ Gera nova key com HMAC-SHA256 seguro

### **O que NÃO pode:**
- ❌ Tenant não pode rotacionar keys de outros tenants
- ❌ Não é possível recuperar a key antiga após rotação
- ❌ Não é possível rotacionar uma key já revogada

---

## 🎨 **Implementação no Frontend**

### **Admin (`/admin/apikeys`):**

```tsx
// Componente AdminAPIKeysPage
const handleRotate = async (keyId: string) => {
  if (!confirm('Rotacionar esta API Key? A chave antiga será revogada.')) {
    return;
  }

  try {
    setRotating(keyId);
    const response = await rotateAPIKey(keyId); // POST /admin/apikeys/rotate
    
    // Mostrar a nova key (apenas uma vez!)
    setNewKey(response.api_key);
    
    toast.success('API Key rotacionada com sucesso!');
    loadAPIKeys(); // Recarregar lista
  } catch (error) {
    toast.error('Erro ao rotacionar API Key');
  } finally {
    setRotating(null);
  }
};

// Botão na tabela
<Button onClick={() => handleRotate(apikey.keyId)}>
  <RefreshCw className="w-4 h-4 mr-2" />
  Rotacionar
</Button>
```

### **Tenant (`/painel/apikeys`):**

```tsx
// Componente PainelAPIKeysPage
const handleRotate = async (keyId: string) => {
  if (!confirm('Rotacionar esta API Key? A chave antiga será revogada.')) {
    return;
  }

  try {
    setRotating(keyId);
    const response = await rotateMyAPIKey(keyId); // POST /me/apikeys/:id/rotate
    
    // Mostrar a nova key (apenas uma vez!)
    setNewKey(response.key);
    
    toast.success('API Key rotacionada com sucesso!');
    loadMyAPIKeys(); // Recarregar lista
  } catch (error) {
    if (error.response?.status === 403) {
      toast.error('Você não tem permissão para rotacionar esta key');
    } else {
      toast.error('Erro ao rotacionar API Key');
    }
  } finally {
    setRotating(null);
  }
};
```

### **API Client (`lib/api/admin.ts` e `lib/api/tenant.ts`):**

```typescript
// Admin
export const rotateAPIKey = async (keyId: string) => {
  const response = await api.post('/admin/apikeys/rotate', { keyId });
  return response.data;
};

// Tenant
export const rotateMyAPIKey = async (keyId: string) => {
  const response = await api.post(`/me/apikeys/${keyId}/rotate`);
  return response.data;
};
```

---

## 📋 **Checklist de Implementação**

### **Backend:**
- [x] Handler `Rotate` em `apikey.go` (admin)
- [x] Handler `RotateAPIKey` em `tenant.go` (tenant)
- [x] Rota `POST /admin/apikeys/rotate`
- [x] Rota `POST /me/apikeys/:id/rotate`
- [x] Validação de ownership para tenant
- [x] Activity log para admin
- [x] Tratamento de erros (400, 403, 404, 500)
- [x] Código compila sem erros
- [ ] Deploy em produção

### **Frontend:**
- [ ] Botão "Rotacionar" em `/admin/apikeys`
- [ ] Botão "Rotacionar" em `/painel/apikeys`
- [ ] Modal de confirmação
- [ ] Mostrar nova key (apenas uma vez!)
- [ ] AlertDialog do Shadcn/ui para confirmação
- [ ] Toast notifications (sucesso/erro)
- [ ] Recarregar lista após rotação
- [ ] Ícone `RefreshCw` do Lucide React

---

## 🚀 **Próximos Passos**

1. **Frontend:**
   - Implementar botões de rotação
   - Testar fluxo completo
   - Garantir UX elegante

2. **Documentação:**
   - Atualizar collection Postman
   - Adicionar exemplos ao guia

3. **Melhorias futuras:**
   - Activity logs para tenant rotation
   - Notificação por email após rotação
   - Rotação automática programada
   - Dashboard de segurança

---

## 📚 **Documentação Relacionada**

- [FIX_APIKEY_TENANT_INVALID.md](./FIX_APIKEY_TENANT_INVALID.md) - Fix da geração de API Keys
- [POSTMAN_GUIDE.md](./POSTMAN_GUIDE.md) - Collection Postman
- [ACTIVITY_LOGS_IMPLEMENTATION.md](./ACTIVITY_LOGS_IMPLEMENTATION.md) - Sistema de logs

---

**Status:** ✅ Backend implementado  
**Próximo:** Implementar frontend (botões e modals)

