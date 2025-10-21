# 🐛 Bug Fix: Erro 500 ao Registrar Novo Cliente

## 📋 Problema Reportado

**Sintoma:**
- Cliente tenta se registrar em `/painel/register`
- Backend cria o **Tenant** e o **User** com sucesso
- Mas retorna **erro 500** (Internal Server Error)
- Frontend exibe erro genérico

---

## 🔍 Diagnóstico

### **1. Logs do Backend:**

```bash
panic: interface conversion: interface {} is primitive.ObjectID, not string
goroutine 44 [running]:
github.com/theretech/retech-core/internal/storage.(*UsersRepo).Create
  /src/internal/storage/users_repo.go:41 +0x1d1
```

### **2. Linha do Erro:**

```go:41:internal/storage/users_repo.go
user.ID = result.InsertedID.(string)  // ❌ PANIC!
```

### **3. Causa Raiz:**

O MongoDB retorna `result.InsertedID` como **`primitive.ObjectID`**, mas o código tentava converter diretamente para **`string`** sem verificação.

```go
// ❌ ERRADO (causa panic)
user.ID = result.InsertedID.(string)

// ✅ CORRETO (conversão segura)
if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
    user.ID = oid.Hex()
}
```

---

## 🔧 Solução Implementada

### **1. Correção no `users_repo.go`:**

#### **Antes:**
```go
result, err := r.coll.InsertOne(ctx, user)
if err != nil {
    return err
}

user.ID = result.InsertedID.(string)  // ❌ PANIC!
return nil
```

#### **Depois:**
```go
result, err := r.coll.InsertOne(ctx, user)
if err != nil {
    return err
}

// Converter ObjectID para string com segurança
if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
    user.ID = oid.Hex()
}

return nil
```

---

### **2. Adicionar Import:**

```go
import (
    "context"
    "time"

    "github.com/theretech/retech-core/internal/domain"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"  // ✅ NOVO
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/crypto/bcrypt"
)
```

---

### **3. Melhorar Tratamento de Erros JWT:**

#### **Antes:**
```go
// Gerar tokens JWT
accessToken, _ := h.jwt.GenerateAccessToken(user)   // ❌ Ignora erro
refreshToken, _ := h.jwt.GenerateRefreshToken(user) // ❌ Ignora erro
```

#### **Depois:**
```go
// Gerar tokens JWT
accessToken, err := h.jwt.GenerateAccessToken(user)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "type":   "https://retech-core/errors/internal-error",
        "title":  "Erro ao gerar token",
        "status": http.StatusInternalServerError,
        "detail": "Erro ao gerar token de acesso",
    })
    return
}

refreshToken, err := h.jwt.GenerateRefreshToken(user)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "type":   "https://retech-core/errors/internal-error",
        "title":  "Erro ao gerar token",
        "status": http.StatusInternalServerError,
        "detail": "Erro ao gerar refresh token",
    })
    return
}
```

---

## 📊 Fluxo Corrigido

### **Antes (com bug):**

```
1. Cliente envia dados de registro
2. ✅ Tenant criado no MongoDB
3. ❌ User.Create() - PANIC ao converter ObjectID
4. ❌ Backend retorna 500
5. ❌ Frontend exibe erro genérico
```

### **Depois (corrigido):**

```
1. Cliente envia dados de registro
2. ✅ Tenant criado no MongoDB
3. ✅ User criado no MongoDB (ID convertido corretamente)
4. ✅ Tokens JWT gerados
5. ✅ Backend retorna 201 Created
6. ✅ Frontend autentica e redireciona para /painel/dashboard
```

---

## 🧪 Como Testar

### **1. Acesse:**
```
http://localhost:3001/painel/register
```

### **2. Preencha o formulário:**

```
Dados da Empresa:
- Nome do Tenant: Tech Solutions LTDA
- Email da Empresa: empresa@techsolutions.com

Dados do Desenvolvedor:
- Nome: João Silva
- Email: joao@techsolutions.com
- Senha: senha123456
- Confirmar Senha: senha123456
```

### **3. Clique em "Criar Conta"**

### **4. Resultado Esperado:**

✅ **Status 201 Created**
✅ **Tenant criado com sucesso**
✅ **User criado com sucesso**
✅ **Tokens JWT gerados**
✅ **Redirecionado para `/painel/dashboard`**

---

## 🔍 Verificação no MongoDB

### **Verificar User criado:**

```bash
docker exec -it build-mongo-1 mongosh

use retech_core

db.users.find({ email: "joao@techsolutions.com" }).pretty()
```

**Resultado Esperado:**
```json
{
  "_id": ObjectId("671683e5..."),  // ✅ ObjectID do MongoDB
  "email": "joao@techsolutions.com",
  "name": "João Silva",
  "password": "$2a$10$...",
  "role": "TENANT_USER",
  "tenantId": "tenant-20251021180000",
  "active": true,
  "createdAt": ISODate("2025-10-21T18:00:00Z"),
  "updatedAt": ISODate("2025-10-21T18:00:00Z")
}
```

---

## 📝 Lições Aprendidas

### **1. Type Assertions em Go:**

Sempre use **type assertion com verificação**:

```go
// ❌ PERIGOSO (pode causar panic)
value := something.(string)

// ✅ SEGURO (verifica antes de converter)
if value, ok := something.(string); ok {
    // usar value
}
```

---

### **2. MongoDB ObjectID:**

O MongoDB sempre retorna `primitive.ObjectID`, não `string`:

```go
// ❌ ERRADO
user.ID = result.InsertedID.(string)

// ✅ CORRETO
if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
    user.ID = oid.Hex()  // Converte para string hexadecimal
}
```

---

### **3. Não Ignore Erros:**

```go
// ❌ RUIM (ignora possíveis erros)
token, _ := generateToken(user)

// ✅ BOM (trata erro adequadamente)
token, err := generateToken(user)
if err != nil {
    return fmt.Errorf("erro ao gerar token: %w", err)
}
```

---

## ✅ Checklist de Validação

Após a correção, verifique:

- [ ] Backend compila sem erros
- [ ] Registro cria Tenant no MongoDB
- [ ] Registro cria User no MongoDB
- [ ] User.ID é convertido corretamente
- [ ] Tokens JWT são gerados
- [ ] Frontend recebe 201 Created
- [ ] Frontend redireciona para `/painel/dashboard`
- [ ] Usuário fica autenticado

---

## 🚀 Deploy

### **1. Rebuild do Backend:**

```bash
cd /path/to/retech-core
docker-compose -f build/docker-compose.yml up --build -d
```

### **2. Verificar Logs:**

```bash
docker-compose -f build/docker-compose.yml logs api --tail=50
```

**Esperado:**
```
{"level":"info","time":"2025-10-21T18:00:00Z","message":"listening on :8080 (env=development)"}
```

---

## 📊 Arquivos Modificados

| Arquivo | Mudança |
|---------|---------|
| `internal/storage/users_repo.go` | Conversão segura de ObjectID para string |
| `internal/storage/users_repo.go` | Import do `bson/primitive` |
| `internal/http/handlers/auth.go` | Tratamento de erros JWT |

---

## 🎯 Resultado Final

| Item | Antes | Depois |
|------|-------|--------|
| **Registro** | ❌ Erro 500 | ✅ Sucesso 201 |
| **Tenant** | ✅ Criado | ✅ Criado |
| **User** | ⚠️ Criado mas com erro | ✅ Criado corretamente |
| **User.ID** | ❌ Panic | ✅ String hexadecimal |
| **JWT** | ⚠️ Erro ignorado | ✅ Tratado |
| **Frontend** | ❌ Erro genérico | ✅ Autenticado e redirecionado |

---

## 🐛 Outros Bugs Relacionados

Este fix também previne problemas similares em:
- Criação de Tenants (se houver conversão de ID)
- Criação de API Keys (se houver conversão de ID)
- Qualquer operação de Insert no MongoDB

---

**Bug corrigido em: 2025-10-21**
**Versão: 1.0.1**

