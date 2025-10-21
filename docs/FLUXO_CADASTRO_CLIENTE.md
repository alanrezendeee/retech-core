# 🚀 Fluxo de Cadastro de Cliente/Desenvolvedor

## 📋 Cenário Atual

### Problema identificado:
- Cliente se cadastra no site (self-service)
- Precisa acessar o **Portal do Desenvolvedor** (`/painel`)
- Necessita de credenciais de acesso (email + senha)
- Deve ter **apenas permissões de desenvolvedor** (não admin)

---

## 🎯 Objetivo

Criar um fluxo onde:
1. Cliente se cadastra sozinho (ou admin cadastra manualmente)
2. Tenant é criado automaticamente
3. Usuário desenvolvedor é criado com senha padrão
4. Cliente recebe email com credenciais
5. Cliente pode alterar senha no primeiro login

---

## 💡 Sugestões de Implementação

### **OPÇÃO 1: Auto-Cadastro com Senha Padrão** ✅ (Recomendada)

#### Fluxo:
```
┌─────────────────────────────────────────────────────────────┐
│  1. Cliente acessa /painel/register                         │
│     - Preenche: Nome, Email, Empresa, Finalidade            │
│     - Cria senha customizada (min 8 caracteres)             │
│     - Aceita termos de uso                                  │
│                                                             │
│  2. Sistema cria automaticamente:                           │
│     ✓ Tenant (company, purpose, active: true)              │
│     ✓ User (role: TENANT_USER, password: bcrypt)           │
│     ✓ API Key inicial (scope: geo, expires: 90 dias)       │
│                                                             │
│  3. Email de boas-vindas enviado:                           │
│     - Confirmação de cadastro                               │
│     - Primeira API Key (para copiar)                        │
│     - Link para dashboard: /painel/dashboard                │
│     - Guia de início rápido                                 │
│                                                             │
│  4. Cliente faz login:                                      │
│     - Acessa /painel/login                                  │
│     - Usa email + senha que criou                           │
│     - Redireciona para /painel/dashboard                    │
└─────────────────────────────────────────────────────────────┘
```

#### Vantagens:
- ✅ Cliente cria sua própria senha (mais seguro)
- ✅ Não depende de email para obter senha
- ✅ Experiência moderna (similar ao GitHub, Vercel)
- ✅ Já implementado no endpoint `/auth/register`

#### Já está pronto:
```go
// retech-core/internal/http/handlers/auth.go
func (h *AuthHandler) Register(c *gin.Context) {
    // 1. Cria Tenant
    // 2. Cria User com senha bcrypt
    // 3. Cria API Key inicial
    // 4. Retorna JWT tokens
}
```

---

### **OPÇÃO 2: Admin Cria Cliente com Senha Padrão** 📧

#### Fluxo:
```
┌─────────────────────────────────────────────────────────────┐
│  1. Admin acessa /admin/tenants                             │
│     - Clica em "Novo Tenant"                                │
│     - Preenche: Nome, Email, Empresa, Finalidade            │
│     - Sistema gera senha padrão: "Retech@2025"              │
│     - Marca: "Forçar troca de senha no primeiro login"      │
│                                                             │
│  2. Sistema cria automaticamente:                           │
│     ✓ Tenant (company, purpose, active: true)              │
│     ✓ User (role: TENANT_USER, password: hash(padrão))     │
│     ✓ User (mustChangePassword: true)                      │
│     ✓ API Key inicial (scope: geo, expires: 90 dias)       │
│                                                             │
│  3. Email automático enviado ao cliente:                    │
│     - Assunto: "Bem-vindo à Retech Core API"               │
│     - Credenciais temporárias:                              │
│       • Email: cliente@empresa.com                          │
│       • Senha: Retech@2025                                  │
│     - Link: /painel/login                                   │
│     - Aviso: "Altere sua senha no primeiro acesso"         │
│                                                             │
│  4. Cliente faz primeiro login:                             │
│     - Usa credenciais do email                              │
│     - Sistema detecta mustChangePassword: true              │
│     - Força troca de senha antes de acessar dashboard       │
│     - Redireciona para /painel/change-password              │
└─────────────────────────────────────────────────────────────┘
```

#### Vantagens:
- ✅ Admin tem controle total sobre cadastros
- ✅ Pode validar clientes antes de criar
- ✅ Útil para vendas B2B (enterprise)

#### A implementar:
- Campo `mustChangePassword` no User
- Endpoint `/auth/change-password`
- Middleware para forçar troca de senha
- Sistema de envio de emails

---

### **OPÇÃO 3: Híbrida (Auto + Admin)** 🎯 (Mais Flexível)

#### Cenário 1: Auto-Cadastro
```
Cliente se registra em /painel/register
→ Cria própria senha
→ Acesso imediato
→ Tenant em status "pending" até aprovação do admin (opcional)
```

#### Cenário 2: Admin Cria
```
Admin cria em /admin/tenants
→ Sistema gera senha padrão
→ Envia email com credenciais
→ Força troca de senha no primeiro login
→ Tenant já ativo
```

---

## 🛠️ Implementação Recomendada (Curto Prazo)

### **FASE 4A: Melhorar Auto-Cadastro Existente** ⚡

#### 1. Adicionar ao endpoint `/auth/register`:
```go
// Após criar tenant, user e API key
// Enviar email de boas-vindas com primeira API key
emailService.SendWelcome(user.Email, tenant.Name, apiKey.Key)
```

#### 2. Melhorar tela `/painel/register`:
```tsx
// Adicionar:
- Validação de senha forte (min 8, maiúscula, número, especial)
- Checkbox "Li e aceito os termos de uso"
- Captcha (opcional, para evitar spam)
- Preview do plano Free (1000 req/dia)
```

#### 3. Email de boas-vindas:
```
Assunto: 🚀 Bem-vindo à Retech Core API!

Olá [Nome],

Sua conta foi criada com sucesso!

📧 Email: [email]
🔑 Sua primeira API Key: rk_live_xxxxxxxxxxxxx

🚀 Próximos passos:
1. Acesse: https://core.theretech.com.br/painel/login
2. Faça login com suas credenciais
3. Veja a documentação em /painel/docs
4. Faça sua primeira requisição!

📚 Documentação: https://docs.theretech.com.br
💬 Suporte: suporte@theretech.com.br

--
Equipe Retech
```

---

### **FASE 4B: Admin Cria Tenant com Senha Padrão** 🔐

#### 1. Modificar `TenantDrawer`:
```tsx
// Adicionar toggle:
"Enviar credenciais por email"
  [x] Gerar senha padrão e enviar email
  [ ] Apenas criar tenant (sem usuário)
```

#### 2. Endpoint novo no backend:
```go
// POST /admin/tenants/create-with-user
{
  "tenantName": "Empresa X",
  "tenantEmail": "contato@empresax.com",
  "company": "Empresa X LTDA",
  "purpose": "Integração mobile",
  "userName": "João Silva",
  "userEmail": "joao@empresax.com",
  "sendEmail": true  // Envia credenciais por email
}

// Sistema:
// 1. Cria tenant
// 2. Cria user com senha padrão: "Retech@" + ano + número aleatório
// 3. Cria API key
// 4. Envia email se sendEmail: true
```

#### 3. Email com credenciais:
```
Assunto: 🎉 Sua conta Retech Core API foi criada

Olá [Nome],

Uma conta foi criada para você na Retech Core API.

📧 Email: joao@empresax.com
🔑 Senha temporária: Retech@2025#7392

⚠️ IMPORTANTE: Altere sua senha no primeiro acesso!

🔗 Acessar: https://core.theretech.com.br/painel/login

Após fazer login, você será direcionado para alterar sua senha.

--
Equipe Retech
```

---

## 🔒 Segurança

### Senha Padrão:
```javascript
// Formato sugerido:
"Retech@" + ano + "#" + random(4 dígitos)

// Exemplos:
"Retech@2025#8472"
"Retech@2025#1938"
"Retech@2025#5621"

// Critérios:
✓ Min 8 caracteres
✓ Letra maiúscula
✓ Letra minúscula
✓ Número
✓ Caracter especial
✓ Único (por causa do random)
```

### Força de troca de senha:
```go
// Middleware
func RequirePasswordChange() gin.HandlerFunc {
  return func(c *gin.Context) {
    user := getUserFromContext(c)
    
    if user.MustChangePassword {
      // Se não estiver na rota de change-password
      if c.Request.URL.Path != "/auth/change-password" {
        c.JSON(403, gin.H{
          "error": "password_change_required",
          "message": "Você precisa alterar sua senha antes de continuar",
          "redirectTo": "/painel/change-password"
        })
        c.Abort()
        return
      }
    }
    
    c.Next()
  }
}
```

---

## 📊 Comparação das Opções

| Critério | Auto-Cadastro | Admin Cria | Híbrida |
|----------|---------------|------------|---------|
| **Velocidade de setup** | ⚡ Imediato | 🐌 Depende do admin | ⚡🐌 Flexível |
| **Controle do admin** | ⭐ Baixo | ⭐⭐⭐ Alto | ⭐⭐ Médio |
| **Experiência do cliente** | 😃 Excelente | 😐 Regular | 😃 Boa |
| **Segurança** | ✅ Cliente cria senha | ⚠️ Senha enviada por email | ✅ Flexível |
| **Complexidade dev** | ✅ Simples | ⚠️ Média | 🔴 Alta |
| **Uso ideal** | B2C, Self-service | B2B, Enterprise | Ambos |

---

## 🎯 Recomendação Final

### **Implementar em 2 fases:**

#### **FASE 4.1: Curto Prazo (1-2 dias)** ✅
1. Melhorar `/painel/register` existente
2. Adicionar validação de senha forte
3. Implementar email de boas-vindas (mesmo que mock)
4. Testar fluxo completo end-to-end

#### **FASE 4.2: Médio Prazo (1 semana)** 🔄
1. Adicionar opção no `/admin/tenants` para criar com usuário
2. Implementar geração de senha padrão
3. Adicionar campo `mustChangePassword` no User
4. Criar endpoint `/auth/change-password`
5. Criar página `/painel/change-password`
6. Implementar envio real de emails (SendGrid, AWS SES, etc)

---

## 📝 Próximos Passos

**Escolha uma opção:**

**A)** Focar no **Auto-Cadastro** (mais rápido, já funciona)
- Melhorar UX do `/painel/register`
- Adicionar email de boas-vindas
- Pronto para produção em 1-2 dias

**B)** Implementar **Admin Cria Tenant** (mais controle)
- Modificar drawer de tenant
- Adicionar geração de senha
- Implementar troca forçada de senha
- Pronto em ~1 semana

**C)** Fazer **Híbrida** (mais completo)
- Combinar A + B
- Máxima flexibilidade
- Pronto em ~2 semanas

---

**Qual opção faz mais sentido para o seu negócio?** 🤔

Considerando que você mencionou:
- "Cliente entra em contato comigo ou acessa o site"
- "Se cadastra como cliente"

→ Sugiro começar com **OPÇÃO A (Auto-Cadastro)** pois:
- ✅ Já está funcionando
- ✅ Experiência moderna
- ✅ Menos dependência sua para cadastros
- ✅ Escalável (self-service)

E depois adicionar **OPÇÃO B** para casos B2B/Enterprise onde você quer aprovar manualmente.

