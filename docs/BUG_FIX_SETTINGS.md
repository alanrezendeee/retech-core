# 🐛 Bug Fix: Erro ao Salvar Configurações

## 📋 Resumo do Problema

Ao tentar salvar configurações em `/admin/settings`, ocorriam dois erros:

1. **Frontend (React):**
   ```
   A component is changing an uncontrolled input to be controlled.
   ```

2. **Backend (MongoDB):**
   ```
   write exception: write errors: [Updating the path 'createdAt' would create a conflict at 'createdAt']
   ```

---

## 🔍 Diagnóstico

### **Erro 1: Uncontrolled → Controlled Input**

**Causa:**
- Campos `requestsPerDay` e `requestsPerMinute` iniciavam como `undefined`
- Ao digitar, mudavam para um valor numérico definido
- React não permite essa transição sem valor inicial

**Sintoma:**
```tsx
// ❌ ANTES (value pode ser undefined)
value={settings.defaultRateLimit.requestsPerDay}

// Quando o usuário digita, muda de undefined → 2000
```

---

### **Erro 2: Conflito no MongoDB**

**Causa:**
- Frontend enviava objeto completo incluindo `createdAt` e `updatedAt`
- Backend tentava fazer `$set` com esses campos
- MongoDB interpretava como tentativa de modificar `createdAt` via `$setOnInsert`
- Conflito: `$set` vs `$setOnInsert` no mesmo campo

**Sintoma:**
```json
{
  "detail": "write exception: write errors: [Updating the path 'createdAt' would create a conflict at 'createdAt']"
}
```

---

## ✅ Solução

### **Fix 1: Frontend - Garantir Valor Inicial**

**Arquivo:** `retech-core-admin/app/admin/settings/page.tsx`

```typescript
// ✅ DEPOIS (sempre tem valor numérico)
value={settings.defaultRateLimit.requestsPerDay || 1000}
onChange={(e) => handleInputChange(
  'defaultRateLimit', 
  'requestsPerDay', 
  parseInt(e.target.value) || 1000  // Fallback para evitar NaN
)}
```

**Resultado:**
- Input sempre tem valor numérico (nunca undefined)
- Transição suave entre valores
- Sem warnings no console

---

### **Fix 2: Frontend - Enviar Apenas Campos Editáveis**

**Arquivo:** `retech-core-admin/app/admin/settings/page.tsx`

```typescript
const handleSaveSettings = async () => {
  try {
    setIsSaving(true);
    
    // ✅ Enviar apenas os campos editáveis (sem id, createdAt, updatedAt)
    const payload = {
      defaultRateLimit: settings.defaultRateLimit,
      cors: settings.cors,
      jwt: settings.jwt,
      api: settings.api,
    };
    
    await api.put('/admin/settings', payload);
    toast.success('Configurações salvas com sucesso!');
    
    // Recarregar para pegar timestamps atualizados
    await loadSettings();
  } catch (error: any) {
    const errorMessage = error.response?.data?.detail || 'Erro ao salvar configurações';
    toast.error(errorMessage);
  } finally {
    setIsSaving(false);
  }
};
```

**Resultado:**
- Payload limpo sem campos de metadados
- Backend recebe apenas dados editáveis
- Sem conflito com campos do MongoDB

---

### **Fix 3: Backend - Separar INSERT de UPDATE**

**Arquivo:** `retech-core/internal/storage/settings_repo.go`

```go
func (r *SettingsRepo) Update(ctx context.Context, settings *domain.SystemSettings) error {
	now := time.Now().UTC()
	
	// Verificar se documento já existe
	var existing domain.SystemSettings
	err := r.col.FindOne(ctx, bson.M{"_id": SystemSettingsID}).Decode(&existing)
	
	if err == mongo.ErrNoDocuments {
		// ✅ Não existe - criar novo com todos os campos
		newSettings := &domain.SystemSettings{
			ID:               SystemSettingsID,
			DefaultRateLimit: settings.DefaultRateLimit,
			CORS:             settings.CORS,
			JWT:              settings.JWT,
			API:              settings.API,
			CreatedAt:        now,  // ← Define aqui
			UpdatedAt:        now,
		}
		_, err = r.col.InsertOne(ctx, newSettings)
		return err
	}
	
	if err != nil {
		return err
	}
	
	// ✅ Existe - atualizar apenas campos editáveis
	_, err = r.col.UpdateOne(
		ctx,
		bson.M{"_id": SystemSettingsID},
		bson.M{
			"$set": bson.M{
				"defaultRateLimit": settings.DefaultRateLimit,
				"cors":             settings.CORS,
				"jwt":              settings.JWT,
				"api":              settings.API,
				"updatedAt":        now,
				// ❌ NÃO inclui createdAt aqui!
			},
		},
	)

	return err
}
```

**Resultado:**
- INSERT: Define `createdAt` + `updatedAt`
- UPDATE: Atualiza apenas campos editáveis + `updatedAt`
- Sem conflito entre `$set` e `$setOnInsert`
- MongoDB feliz 😊

---

### **Fix 4: Backend - Adicionar Logs de Debug**

**Arquivo:** `retech-core/internal/storage/settings_repo.go`

```go
fmt.Printf("📝 Atualizando settings: %+v\n", settings)

if err == mongo.ErrNoDocuments {
	fmt.Println("ℹ️ Documento não existe, criando novo...")
	// ... INSERT ...
	fmt.Println("✅ Documento criado com sucesso!")
} else {
	fmt.Println("ℹ️ Documento existe, atualizando...")
	// ... UPDATE ...
	fmt.Println("✅ Documento atualizado com sucesso!")
}
```

**Resultado:**
- Facilita debug em desenvolvimento
- Confirma fluxo correto (INSERT vs UPDATE)
- Identifica problemas rapidamente

---

## 🧪 Teste de Validação

```bash
#!/bin/bash

# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alanrezendeee@gmail.com","password":"admin123456"}' \
  | jq -r '.accessToken')

# 2. Atualizar configurações
curl -s -X PUT http://localhost:8080/admin/settings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "defaultRateLimit": {
      "RequestsPerDay": 2000,
      "RequestsPerMinute": 120
    },
    "cors": {
      "enabled": true,
      "allowedOrigins": ["https://core.theretech.com.br"]
    },
    "jwt": {
      "accessTokenTTL": 900,
      "refreshTokenTTL": 604800
    },
    "api": {
      "version": "1.0.0",
      "environment": "development",
      "maintenance": false
    }
  }' | jq '.'

# 3. Verificar se salvou
curl -s http://localhost:8080/admin/settings \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.defaultRateLimit'

# Resultado esperado:
# {
#   "RequestsPerDay": 2000,
#   "RequestsPerMinute": 120
# }
```

---

## 📊 Antes vs Depois

| Aspecto | ❌ Antes | ✅ Depois |
|---------|---------|----------|
| **Frontend Input** | `value={undefined}` → Warning | `value={1000}` → OK |
| **Payload** | Inclui `createdAt`, `updatedAt` | Apenas campos editáveis |
| **Backend INSERT** | `$set` + `$setOnInsert` juntos | INSERT com campos completos |
| **Backend UPDATE** | Conflito no `createdAt` | Atualiza apenas editáveis |
| **Erro MongoDB** | 500 (conflict) | 200 (success) |
| **Logs** | Nenhum | Logs detalhados |

---

## ✅ Checklist de Verificação

- [x] Input `requestsPerDay` sempre tem valor numérico
- [x] Input `requestsPerMinute` sempre tem valor numérico
- [x] Frontend envia apenas campos editáveis
- [x] Backend separa lógica de INSERT e UPDATE
- [x] MongoDB não reclama de conflito em `createdAt`
- [x] Configurações são salvas corretamente
- [x] Timestamps (`createdAt`, `updatedAt`) funcionam
- [x] Logs de debug ajudam no troubleshooting

---

## 🎯 Resultado Final

```
🔐 Fazendo login...
✅ Token obtido!

📋 Obtendo configurações atuais...
{
  "RequestsPerDay": 1000,
  "RequestsPerMinute": 60
}

💾 Atualizando configurações...
{
  "message": "Configurações atualizadas com sucesso",
  "settings": {
    "defaultRateLimit": {
      "RequestsPerDay": 2000,
      "RequestsPerMinute": 120
    },
    ...
  }
}

🔍 Verificando se salvou...
{
  "RequestsPerDay": 2000,  ✅
  "RequestsPerMinute": 120 ✅
}
```

---

**Sistema de configurações 100% funcional! 🚀**

