# 🆕 Fluxo de Registro de Novos Clientes

## 📋 Visão Geral

Quando um novo cliente se registra no sistema, o seguinte acontece automaticamente:

1. ✅ **Cria um Tenant** (empresa)
2. ✅ **Cria um User** (desenvolvedor principal)
3. ✅ **Autentica automaticamente** (JWT tokens)
4. ✅ **Redireciona para o dashboard do desenvolvedor**

---

## 🔄 Fluxo Completo

```
┌─────────────────────────────────────────────────────────┐
│  CLIENTE ACESSA /painel/register                        │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  PREENCHE FORMULÁRIO                                    │
│  ┌────────────────────────────────────────────────┐     │
│  │ Dados da Empresa:                              │     │
│  │ • Nome do Tenant: "Tech Solutions LTDA"       │     │
│  │ • Email da Empresa: empresa@techsolutions.com │     │
│  │ • Empresa: "Tech Solutions LTDA" (opcional)   │     │
│  │ • Finalidade: "Sistema de vendas" (opcional)  │     │
│  │                                                │     │
│  │ Dados do Desenvolvedor:                        │     │
│  │ • Nome: "João Silva"                          │     │
│  │ • Email: joao@techsolutions.com               │     │
│  │ • Senha: ********                             │     │
│  │ • Confirmar Senha: ********                   │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  BACKEND VALIDA                                         │
│  ┌────────────────────────────────────────────────┐     │
│  │ ❌ Email do desenvolvedor já existe?           │     │
│  │    → Retorna 409: "O email 'xxx' já está      │     │
│  │      cadastrado. Use outro ou faça login."     │     │
│  │                                                │     │
│  │ ❌ Email da empresa já existe?                 │     │
│  │    → Retorna 409: "Email da empresa já        │     │
│  │      cadastrado."                              │     │
│  │                                                │     │
│  │ ✅ Tudo OK?                                     │     │
│  │    → Continua...                               │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  CRIA TENANT (MongoDB)                                  │
│  ┌────────────────────────────────────────────────┐     │
│  │ {                                              │     │
│  │   "tenantId": "tenant-20251021180000",        │     │
│  │   "name": "Tech Solutions LTDA",              │     │
│  │   "email": "empresa@techsolutions.com",       │     │
│  │   "company": "Tech Solutions LTDA",           │     │
│  │   "purpose": "Sistema de vendas",             │     │
│  │   "active": true,                             │     │
│  │   "createdAt": "2025-10-21T18:00:00Z"         │     │
│  │ }                                              │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  CRIA USER (MongoDB)                                    │
│  ┌────────────────────────────────────────────────┐     │
│  │ {                                              │     │
│  │   "email": "joao@techsolutions.com",          │     │
│  │   "name": "João Silva",                       │     │
│  │   "password": "$2a$10$...",  // hash bcrypt   │     │
│  │   "role": "TENANT_USER",                      │     │
│  │   "tenantId": "tenant-20251021180000",        │     │
│  │   "active": true,                             │     │
│  │   "createdAt": "2025-10-21T18:00:00Z"         │     │
│  │ }                                              │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  GERA TOKENS JWT                                        │
│  ┌────────────────────────────────────────────────┐     │
│  │ • Access Token: eyJhbGc...  (15 min)          │     │
│  │ • Refresh Token: eyJhbGc... (7 dias)          │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  FRONTEND SALVA AUTENTICAÇÃO                            │
│  ┌────────────────────────────────────────────────┐     │
│  │ localStorage:                                  │     │
│  │ • accessToken                                  │     │
│  │ • refreshToken                                 │     │
│  │ • user (JSON)                                  │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  REDIRECIONA PARA /painel/dashboard                     │
│  ✅ Cliente já está logado automaticamente!             │
└─────────────────────────────────────────────────────────┘
```

---

## 🧪 Como Testar

### **1. Acesse a página de registro:**
```
http://localhost:3001/painel/register
```

### **2. Preencha o formulário:**

#### **Dados da Empresa:**
```
Nome do Tenant: Tech Solutions LTDA
Email da Empresa: empresa@techsolutions.com
Empresa (opcional): Tech Solutions LTDA
Finalidade (opcional): Sistema de vendas
```

#### **Dados do Desenvolvedor:**
```
Nome: João Silva
Email: joao@techsolutions.com
Senha: senha123456
Confirmar Senha: senha123456
```

### **3. Clique em "Criar Conta"**

### **4. Você será redirecionado para:**
```
http://localhost:3001/painel/dashboard
```
✅ Já autenticado como `TENANT_USER`!

---

## ⚠️ Erros Comuns

### **Erro 409: "Email já em uso"**

#### **Causa:**
O email do **desenvolvedor** já foi usado em outro cadastro.

#### **Solução:**
Use um email diferente:
```
❌ joao@techsolutions.com  (já usado)
✅ joao2@techsolutions.com (novo)
✅ dev@novaempresa.com     (novo)
```

---

### **Erro 409: "Email da empresa já em uso"**

#### **Causa:**
O email da **empresa** já foi usado em outro tenant.

#### **Solução:**
Use um email diferente para a empresa:
```
❌ empresa@techsolutions.com  (já usado)
✅ contato@techsolutions.com  (novo)
✅ empresa@outraempresa.com   (novo)
```

---

## 🔐 Validações Implementadas

### **Backend (`/auth/register`):**

1. ✅ **Validação de campos obrigatórios**
   - Nome do Tenant
   - Email da Empresa
   - Nome do Desenvolvedor
   - Email do Desenvolvedor
   - Senha

2. ✅ **Email do desenvolvedor único**
   - Retorna 409 se já existe
   - Mensagem: `"O email 'xxx' já está cadastrado. Por favor, use outro email ou faça login."`

3. ✅ **Email da empresa único**
   - Retorna 409 se já existe
   - Mensagem: `"O email da empresa 'xxx' já está cadastrado."`

4. ✅ **Senha forte** (mínimo 8 caracteres)
   - Hash bcrypt automático

5. ✅ **Tenant ativo por padrão**
   - `active: true`

6. ✅ **User ativo por padrão**
   - `active: true`

7. ✅ **Role automático**
   - `role: TENANT_USER`

---

## 📊 Estrutura dos Dados Criados

### **Tenant (MongoDB - `tenants` collection):**
```json
{
  "_id": ObjectId("..."),
  "tenantId": "tenant-20251021180000",
  "name": "Tech Solutions LTDA",
  "email": "empresa@techsolutions.com",
  "company": "Tech Solutions LTDA",
  "purpose": "Sistema de vendas",
  "active": true,
  "createdAt": ISODate("2025-10-21T18:00:00Z"),
  "updatedAt": ISODate("2025-10-21T18:00:00Z")
}
```

### **User (MongoDB - `users` collection):**
```json
{
  "_id": ObjectId("..."),
  "email": "joao@techsolutions.com",
  "name": "João Silva",
  "password": "$2a$10$XxXxXxXxXxXxXxXxXxXxXxXx",
  "role": "TENANT_USER",
  "tenantId": "tenant-20251021180000",
  "active": true,
  "createdAt": ISODate("2025-10-21T18:00:00Z"),
  "updatedAt": ISODate("2025-10-21T18:00:00Z")
}
```

---

## 🎯 Próximos Passos Após Registro

Depois de se registrar, o desenvolvedor pode:

1. ✅ **Acessar `/painel/dashboard`**
   - Ver visão geral da conta

2. ✅ **Criar API Keys em `/painel/apikeys`**
   - Gerar chaves para acessar a API

3. ✅ **Ver uso em `/painel/usage`**
   - Acompanhar consumo da API

4. ✅ **Acessar docs em `/painel/docs`**
   - Ler documentação da API

---

## 🔄 Diferença: Registro vs Login

### **Registro (`/painel/register`):**
- ✅ Cria novo Tenant
- ✅ Cria novo User
- ✅ Autentica automaticamente
- ✅ Redireciona para `/painel/dashboard`

### **Login (`/painel/login`):**
- ❌ NÃO cria Tenant
- ❌ NÃO cria User
- ✅ Autentica usuário existente
- ✅ Redireciona para `/painel/dashboard`

---

## 🛡️ Segurança

### **Senhas:**
- ✅ Hash bcrypt (custo 10)
- ✅ Nunca retornadas na API
- ✅ Validação de força no frontend

### **Tokens JWT:**
- ✅ Access Token: 15 minutos
- ✅ Refresh Token: 7 dias
- ✅ Auto-renovação no frontend

### **Emails:**
- ✅ Validação de unicidade
- ✅ Validação de formato
- ✅ Case-sensitive

---

## 🧪 Exemplos de Teste

### **Cenário 1: Registro bem-sucedido**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenantName": "Tech Solutions LTDA",
    "tenantEmail": "empresa@techsolutions.com",
    "company": "Tech Solutions LTDA",
    "purpose": "Sistema de vendas",
    "userName": "João Silva",
    "userEmail": "joao@techsolutions.com",
    "userPassword": "senha123456"
  }'
```

**Resposta (201 Created):**
```json
{
  "tenant": {
    "tenantId": "tenant-20251021180000",
    "name": "Tech Solutions LTDA",
    "email": "empresa@techsolutions.com",
    "active": true
  },
  "user": {
    "email": "joao@techsolutions.com",
    "name": "João Silva",
    "role": "TENANT_USER",
    "active": true
  },
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "apiKey": ""
}
```

---

### **Cenário 2: Email já cadastrado**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenantName": "Nova Empresa",
    "tenantEmail": "nova@empresa.com",
    "userName": "Maria",
    "userEmail": "joao@techsolutions.com",  # Email duplicado!
    "userPassword": "senha123456"
  }'
```

**Resposta (409 Conflict):**
```json
{
  "type": "https://retech-core/errors/conflict",
  "title": "Email já em uso",
  "status": 409,
  "detail": "O email 'joao@techsolutions.com' já está cadastrado. Por favor, use outro email ou faça login."
}
```

---

## ✅ Checklist de Validações

Antes de criar conta, o sistema verifica:

- [ ] Email do desenvolvedor não existe?
- [ ] Email da empresa não existe?
- [ ] Senha tem no mínimo 8 caracteres?
- [ ] Campos obrigatórios preenchidos?
- [ ] Email tem formato válido?

Se **TUDO OK** ✅ → Cria Tenant + User + Autentica

Se **ALGUM ERRO** ❌ → Retorna mensagem clara

---

## 🎨 Visual do Formulário

O formulário de registro está dividido em duas seções:

### **1. Dados da Empresa (Tenant):**
```
┌──────────────────────────────────────┐
│ 🏢 Dados da Empresa                  │
├──────────────────────────────────────┤
│ Nome do Tenant: [____________]       │
│ Email da Empresa: [__________]       │
│ Empresa (opcional): [________]       │
│ Finalidade (opcional): [_____]       │
└──────────────────────────────────────┘
```

### **2. Dados do Desenvolvedor (User):**
```
┌──────────────────────────────────────┐
│ 👤 Dados do Desenvolvedor            │
├──────────────────────────────────────┤
│ Nome: [______________]               │
│ Email: [_____________]               │
│ Senha: [*************]               │
│ Confirmar Senha: [***********]       │
└──────────────────────────────────────┘
```

---

## 📝 Resumo Final

| Item | Descrição |
|------|-----------|
| **Endpoint** | `POST /auth/register` |
| **Acesso Frontend** | `/painel/register` |
| **Cria** | Tenant + User |
| **Autentica** | Sim (automático) |
| **Redireciona** | `/painel/dashboard` |
| **Validações** | Email único, senha forte |
| **Mensagens** | Claras e em português |
| **Segurança** | Hash bcrypt, JWT tokens |

---

**Agora você pode testar o registro com diferentes emails! 🎉**

