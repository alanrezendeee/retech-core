# 🗺️ ROADMAP - Retech Core API

Sistema completo de API como serviço (API-as-a-Service) com gerenciamento de tenants, API Keys, analytics e portais administrativos.

**Última atualização**: 2025-10-20 23:30  
**Versão atual**: 0.4.0  
**Status geral**: 🟢 Pronto para produção (features base)

---

## 📅 Histórico de Atualizações

### 2025-10-20 - Sessão Épica! 🎉
- ✅ **FASE 0** concluída (100%)
- ✅ **FASE 1** concluída (100%)
- ✅ Frontend criado (retech-admin)
- ✅ 40 arquivos criados
- ✅ ~4.500 linhas de código
- ✅ Autenticação JWT completa
- ✅ Rate limiting implementado
- ✅ Logging de uso implementado
- ✅ Landing page, login, registro, dashboards
- ✅ Arquitetura de domínio único
- ✅ Docker build corrigido (Go 1.24)

---

## 📊 Visão Geral do Progresso

```
FASE 0: Fundação           ██████████ 100% ✅
FASE 1: Auth & Segurança   ██████████ 100% ✅
FASE 2: Admin Dashboard    ████░░░░░░  40% 🟡
FASE 3: Developer Portal   ████░░░░░░  40% 🟡
FASE 4: Logs & Analytics   ░░░░░░░░░░   0% 🔴
FASE 5: Melhorias          ░░░░░░░░░░   0% 🔴
FASE 6: Monetização        ░░░░░░░░░░   0% 🔴
```

**Progresso total**: 47% (35/75 tarefas concluídas)

---

## 🏗️ Arquitetura de Deployment

### 🌐 Domínio Unificado: `core.theretech.com.br`

Toda a aplicação será servida através de um único domínio com paths diferentes:

```
https://core.theretech.com.br
├─ /api/*      → Backend Go (retech-core)     - API REST
├─ /admin/*    → Admin Dashboard (React)      - Super Admin UI
└─ /painel/*   → Developer Portal (React)     - Tenant UI
```

### 📦 Serviços

#### 1. Backend - retech-core (Go + Gin)
```
Port: 8080 (interno no Railway)
Paths públicos externos: /api/*

Endpoints:
├─ /auth/*   → Login, register, JWT
├─ /geo/*    → Dados geográficos (protegido)
├─ /admin/*  → Gestão (super admin)
├─ /me/*     → Tenant self-service
├─ /health   → Health check
└─ /docs     → Redoc
```

#### 2. Frontend - retech-admin (Next.js 14)
```
Port: 3000 (Railway)
Domínio público: core.theretech.com.br

Routes:
├─ /admin/*  → Dashboard admin (Shadcn/ui)
├─ /painel/* → Portal desenvolvedor
└─ /api/*    → Proxy → Go backend (rewrites)
```

#### 3. Database - MongoDB
```
Railway MongoDB ou MongoDB Atlas
Internal: mongodb://...
```

### 🔄 Fluxo de Requisições

```
Usuário
  │
  ├─ https://core.theretech.com.br/admin/dashboard
  │    └─→ Next.js renderiza admin UI
  │
  ├─ https://core.theretech.com.br/painel/apikeys
  │    └─→ Next.js renderiza portal UI
  │
  └─ https://core.theretech.com.br/api/geo/ufs
       └─→ Next.js rewrites → Go backend:8080/geo/ufs
```

### ⚙️ Configuração Next.js

```typescript
// next.config.ts
{
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: process.env.BACKEND_URL + '/:path*'
      }
    ]
  }
}
```

### 🎯 Vantagens desta Arquitetura

✅ **Domínio único** - Simples de lembrar  
✅ **CORS facilitado** - Tudo mesmo origin  
✅ **SSL único** - Railway gerencia automaticamente  
✅ **UX fluida** - Navegação entre seções  
✅ **Manutenção simples** - Um deploy, múltiplos serviços  
✅ **Escalável** - Fácil adicionar novos paths  
✅ **SEO friendly** - Paths semânticos  

### 🔐 Configuração de Rate Limit Admin

O admin terá interface para configurar:
- Requests por dia (free tier)
- Requests por minuto
- Limites customizados por tenant
- Whitelist de IPs
- Blacklist de IPs

---

## 🎯 FASE 0: Fundação ⚙️

**Objetivo**: Estrutura básica do sistema (backend + frontend)  
**Status**: 🟢 100% CONCLUÍDO ✅  
**Prioridade**: 🔴 CRÍTICA  
**Concluído em**: 2025-10-20

### Backend (retech-core)

- [x] Setup projeto Go + MongoDB
- [x] Sistema de Tenants (CRUD)
- [x] Sistema de API Keys (CRUD)
- [x] Endpoints GEO (estados, municípios)
- [x] Migrations e seeds automáticos
- [x] Deploy Railway configurado
- [x] Documentação OpenAPI/Redoc
- [x] **Sistema de Users (auth)** ✅
- [x] **Roles e permissões** ✅

### Frontend (retech-admin)

- [x] **Criar projeto Next.js 14 (App Router)** ✅
- [x] **Setup Shadcn/ui + Tailwind CSS** ✅
- [x] **Estrutura de pastas (admin + painel)** ✅
- [x] **Configurar Axios/cliente HTTP** ✅
- [x] **Setup variáveis de ambiente** ✅
- [x] **Configurar basePath e rewrites** ✅
  - `/admin/*` → Admin dashboard
  - `/painel/*` → Developer portal
  - `/api/*` → Proxy para Go backend
- [x] **Layout compartilhado entre admin/painel** ✅
- [x] **Sistema de rotas com middleware** ✅

**Entregue**: 2025-10-20 ✅

---

## 🔐 FASE 1: Autenticação e Segurança

**Objetivo**: Sistema completo de auth + proteção da API  
**Status**: 🟢 100% CONCLUÍDO ✅  
**Prioridade**: 🔴 CRÍTICA  
**Concluído em**: 2025-10-20

### Backend (retech-core)

#### 1.1 Sistema de Usuários ✅
- [x] Model `User` (id, email, password, role, tenantId)
- [x] Repository `UsersRepo`
- [x] Hash de senha (bcrypt)
- [x] Validações de email/senha

**Arquivos**: `internal/domain/user.go`, `internal/storage/users_repo.go`

#### 1.2 Autenticação JWT ✅
- [x] Endpoint `POST /auth/login`
- [x] Endpoint `POST /auth/register` (self-service)
- [x] Endpoint `GET /auth/me`
- [x] Endpoint `POST /auth/refresh`
- [x] Geração de JWT (access + refresh token)
- [x] Middleware JWT (`AuthJWT`)
- [x] Roles: `SUPER_ADMIN`, `TENANT_USER`

**Arquivos**: `internal/auth/jwt.go`, `internal/auth/jwt_middleware.go`, `internal/http/handlers/auth.go`

**Configuração**:
- Access token: 15 minutos
- Refresh token: 7 dias
- Algoritmo: HS256

#### 1.3 Proteção de Endpoints ✅
- [x] Aplicar `AuthAPIKey` em `/geo/*`
- [x] Retornar 401 sem API Key válida
- [x] Manter `/health`, `/version`, `/docs` públicos
- [x] Endpoints admin requerem `SUPER_ADMIN`
- [x] Endpoints `/me/*` requerem `TENANT_USER`

**Arquivos**: `internal/http/router.go`

#### 1.4 Rate Limiting ✅
- [x] Model `RateLimit` (key, count, window)
- [x] Middleware de rate limiting
- [x] Configuração por plano:
  - Free: 1000 req/dia, 100 req/minuto ✅
  - Pro: 10k req/dia (configurável)
- [x] Headers `X-RateLimit-*`
- [x] Retornar 429 quando exceder
- [x] Resetar contadores diariamente

**Arquivos**: `internal/domain/rate_limit.go`, `internal/middleware/rate_limiter.go`

#### 1.5 Logging de Uso ✅
- [x] Collection `api_usage_logs`
- [x] Middleware para registrar:
  - Timestamp ✅
  - API Key ✅
  - Tenant ID ✅
  - Endpoint ✅
  - Status code ✅
  - Response time ✅
  - IP address ✅
  - User agent ✅

**Arquivos**: `internal/domain/api_usage_log.go`, `internal/middleware/usage_logger.go`

### Frontend (retech-admin)

#### 1.6 Páginas de Auth ✅
- [x] Tela de login admin (`/admin/login`)
- [x] Tela de login painel (`/painel/login`)
- [x] Tela de registro (`/painel/register`)
- [x] Form validation (React Hook Form + Zod)
- [x] Loading states
- [x] Error handling

**Arquivos**: `app/admin/login/page.tsx`, `app/painel/login/page.tsx`, `app/painel/register/page.tsx`

#### 1.7 Proteção de Rotas ✅
- [x] Auth store (Zustand)
- [x] Redirect para `/login` se não autenticado
- [x] Verificar role (admin vs tenant)
- [x] Guardar JWT em localStorage
- [x] Auto-refresh de token via interceptor

**Arquivos**: `lib/stores/auth-store.ts`, `lib/api/client.ts`, `lib/api/auth.ts`

**Entregue**: 2025-10-20 ✅

---

## 👨‍💼 FASE 2: Admin Dashboard

**Objetivo**: Interface completa para super admin gerenciar tudo  
**Status**: 🔴 Não iniciado (0%)  
**Prioridade**: 🔴 ALTA

### Backend (retech-core)

#### 2.1 Endpoints Admin
- [ ] `GET /admin/stats` - KPIs globais
- [ ] `GET /admin/tenants` - Lista todos tenants
- [ ] `GET /admin/tenants/:id` - Detalhes do tenant
- [ ] `PUT /admin/tenants/:id` - Editar tenant
- [ ] `PUT /admin/tenants/:id/status` - Ativar/desativar
- [ ] `DELETE /admin/tenants/:id` - Deletar tenant
- [ ] `GET /admin/apikeys` - Todas API Keys
- [ ] `POST /admin/apikeys` - Criar key para tenant
- [ ] `DELETE /admin/apikeys/:id` - Revogar key
- [ ] `GET /admin/usage` - Analytics global
- [ ] `GET /admin/usage/by-tenant` - Uso por tenant
- [ ] `GET /admin/usage/by-endpoint` - Uso por endpoint

### Frontend (retech-admin)

#### 2.2 Layout Admin
- [ ] Sidebar com navegação
- [ ] Header com user menu
- [ ] Breadcrumbs
- [ ] Notifications center
- [ ] Search global

#### 2.3 Dashboard Overview
- [ ] Cards com KPIs:
  - Total de tenants
  - Total de API Keys ativas
  - Requests hoje/mês
  - Uptime
- [ ] Gráfico de requests (7/30 dias)
- [ ] Top 10 tenants por uso
- [ ] Últimos tenants registrados
- [ ] Alertas de sistema

#### 2.4 Gestão de Tenants
- [ ] Lista de tenants (tabela)
- [ ] Filtros: status, data criação, uso
- [ ] Busca por nome/email
- [ ] Paginação
- [ ] Ações: editar, ativar/desativar, deletar
- [ ] Modal de criação manual
- [ ] Modal de edição
- [ ] Modal de confirmação (delete)

#### 2.5 Detalhes do Tenant
- [ ] Informações completas
- [ ] Status e histórico
- [ ] Lista de API Keys do tenant
- [ ] Gráfico de uso do tenant
- [ ] Logs recentes do tenant
- [ ] Ações rápidas:
  - Criar nova key
  - Revogar keys
  - Ativar/desativar tenant
  - Enviar email

#### 2.6 Gestão de API Keys
- [ ] Lista global de keys (tabela)
- [ ] Filtros por tenant, status
- [ ] Busca por key/tenant
- [ ] Ver último uso
- [ ] Revogar em massa
- [ ] Exportar lista

#### 2.7 Analytics Admin
- [ ] Gráficos de uso:
  - Timeline (requests ao longo do tempo)
  - Por endpoint
  - Por tenant
  - Por status code
- [ ] Filtros de data
- [ ] Métricas de performance:
  - Latência média
  - Taxa de erro
  - Endpoints mais usados
- [ ] Exportar relatórios (CSV)

**Prazo**: 2025-10-27 (6 dias)  
**Dependências**: FASE 1

---

## 👨‍💻 FASE 3: Developer Portal

**Objetivo**: Portal self-service para desenvolvedores (tenants)  
**Status**: 🔴 Não iniciado (0%)  
**Prioridade**: 🟡 MÉDIA

### Backend (retech-core)

#### 3.1 Endpoints do Tenant
- [ ] `GET /me` - Dados do tenant logado
- [ ] `PUT /me` - Atualizar perfil
- [ ] `GET /me/usage` - Uso atual do tenant
- [ ] `GET /me/usage/history` - Histórico de uso
- [ ] `GET /me/apikeys` - Lista keys do tenant
- [ ] `POST /me/apikeys` - Criar nova key
- [ ] `PUT /me/apikeys/:id` - Renomear key
- [ ] `POST /me/apikeys/:id/rotate` - Rotacionar key
- [ ] `DELETE /me/apikeys/:id` - Deletar key
- [ ] `GET /me/logs` - Logs do tenant (limitado)

### Frontend (retech-admin)

#### 3.2 Dashboard do Desenvolvedor
- [ ] Overview cards:
  - Requests hoje
  - Requests mês
  - Limite restante
  - Status da conta
- [ ] Gráfico de uso (últimos 7 dias)
- [ ] Quick actions:
  - Criar API Key
  - Ver docs
  - Ver limites

#### 3.3 Gestão de API Keys
- [ ] Lista de keys do tenant
- [ ] Criar nova key (modal)
  - Nome da key
  - Mostra key APENAS uma vez
  - Copiar para clipboard
- [ ] Rotacionar key
- [ ] Revogar key
- [ ] Ver último uso
- [ ] Status (ativa/inativa)

#### 3.4 Uso e Analytics
- [ ] Gráfico de requests
- [ ] Endpoints mais usados
- [ ] Horários de pico
- [ ] Histórico mensal
- [ ] Alertas quando próximo do limite

#### 3.5 Documentação
- [ ] Getting Started
- [ ] Exemplos de código (cURL, JavaScript, Python, Go)
- [ ] Lista de endpoints disponíveis
- [ ] Limites do plano atual
- [ ] FAQ
- [ ] Status da API

#### 3.6 Configurações
- [ ] Editar perfil
- [ ] Alterar senha
- [ ] Preferências de notificação
- [ ] Webhooks (futuro)

**Prazo**: 2025-10-30 (5 dias)  
**Dependências**: FASE 2

---

## 📊 FASE 4: Logging e Monitoramento

**Objetivo**: Sistema completo de analytics e observabilidade  
**Status**: 🔴 Não iniciado (0%)  
**Prioridade**: 🟢 BAIXA

### Backend (retech-core)

#### 4.1 Sistema de Logs
- [ ] Collection `api_usage_logs` otimizada
- [ ] Agregações por hora/dia/mês
- [ ] Índices para queries rápidas
- [ ] Retention policy (ex: 90 dias)
- [ ] Arquivamento de logs antigos

#### 4.2 Analytics Engine
- [ ] Calcular métricas agregadas
- [ ] Cache de métricas (Redis?)
- [ ] Background jobs para agregação
- [ ] Endpoints de analytics com filtros avançados

#### 4.3 Alertas
- [ ] Sistema de alertas configurável
- [ ] Alertas de uso excessivo
- [ ] Alertas de erro (taxa alta)
- [ ] Alertas de abuso/padrões suspeitos
- [ ] Notificações (email/webhook)

### Frontend (retech-admin)

#### 4.4 Analytics Avançado
- [ ] Dashboards customizáveis
- [ ] Filtros avançados (data, tenant, endpoint)
- [ ] Comparação de períodos
- [ ] Heatmaps de uso
- [ ] Geolocalização de requests
- [ ] Exportar relatórios PDF

#### 4.5 Logs em Tempo Real
- [ ] Stream de logs (WebSocket)
- [ ] Filtros em tempo real
- [ ] Search/highlight
- [ ] Pausar/retomar stream

#### 4.6 Alertas UI
- [ ] Configurar alertas
- [ ] Ver histórico de alertas
- [ ] Notifications center
- [ ] Email/SMS configuration

**Prazo**: 2025-11-03 (4 dias)  
**Dependências**: FASE 3

---

## ✨ FASE 5: Melhorias e Polimento

**Objetivo**: UX refinada, emails, testes  
**Status**: 🔴 Não iniciado (0%)  
**Prioridade**: 🟢 BAIXA

### Backend (retech-core)

#### 5.1 Sistema de Emails
- [ ] Serviço de email (SendGrid/Resend)
- [ ] Templates de email
- [ ] Email de boas-vindas
- [ ] Email de confirmação
- [ ] Email de API Key criada
- [ ] Email de alerta de uso
- [ ] Email de recuperação de senha

#### 5.2 Segurança Adicional
- [ ] Rate limiting por IP
- [ ] Detecção de abuso
- [ ] Blacklist/Whitelist de IPs
- [ ] 2FA (TOTP)
- [ ] Audit log completo

#### 5.3 Testes
- [ ] Testes unitários (backend)
- [ ] Testes de integração
- [ ] Testes E2E (frontend)
- [ ] CI/CD completo

### Frontend (retech-admin)

#### 5.4 UX/UI
- [ ] Dark mode
- [ ] Responsivo mobile
- [ ] Animações e transições
- [ ] Loading skeletons
- [ ] Empty states
- [ ] Error boundaries
- [ ] Acessibilidade (a11y)

#### 5.5 Features
- [ ] Landing page pública
- [ ] Página de status (`status.retech.com`)
- [ ] Blog/Changelog
- [ ] Termos de uso
- [ ] Política de privacidade
- [ ] Tutoriais em vídeo

**Prazo**: 2025-11-10 (7 dias)  
**Dependências**: FASE 4

---

## 💰 FASE 6: Monetização (Futuro)

**Objetivo**: Planos pagos e billing  
**Status**: 🔴 Planejado (0%)  
**Prioridade**: 🟢 FUTURO

### Backend

#### 6.1 Sistema de Planos
- [ ] Model `Plan` (name, limits, price)
- [ ] Free, Pro, Business tiers
- [ ] Limites configuráveis por plano
- [ ] Upgrade/downgrade de plano
- [ ] Proration de valores

#### 6.2 Billing
- [ ] Integração Stripe
- [ ] Checkout de planos
- [ ] Gerenciar cartões
- [ ] Invoices automáticos
- [ ] Histórico de pagamentos
- [ ] Cancelamento de plano

#### 6.3 Features Premium
- [ ] Webhooks customizados
- [ ] SLA garantido
- [ ] Suporte prioritário
- [ ] White label
- [ ] API para parceiros

### Frontend

#### 6.4 Billing UI
- [ ] Página de planos
- [ ] Comparação de planos
- [ ] Checkout Stripe
- [ ] Gerenciar assinatura
- [ ] Invoices e recibos
- [ ] Upgrade/downgrade UI

**Prazo**: TBD (a definir)  
**Dependências**: FASE 5

---

## 📈 Métricas de Sucesso

### Técnicas
- [ ] 99.9% uptime
- [ ] Latência média < 100ms
- [ ] Taxa de erro < 0.1%
- [ ] Cobertura de testes > 80%
- [ ] Build time < 5min

### Negócio
- [ ] 100 tenants registrados
- [ ] 1M requests/dia
- [ ] 10% conversão free → paid
- [ ] NPS > 50
- [ ] Churn < 5%/mês

---

## 🚀 Deploy e Infraestrutura

### Ambientes
- [ ] Development (local)
- [ ] Staging (Railway)
- [ ] Production (Railway/AWS)

### CI/CD
- [ ] GitHub Actions
- [ ] Deploy automático (main → prod)
- [ ] Preview deploys (PRs)
- [ ] Rollback automático

### Monitoramento
- [ ] Railway metrics
- [ ] Sentry (error tracking)
- [ ] Uptime monitoring
- [ ] Performance monitoring

---

## 📚 Documentação

- [x] README.md (geral)
- [x] RAILWAY_DEPLOY.md
- [x] CHANGELOG.md
- [ ] **API_REFERENCE.md** (completo)
- [ ] **CONTRIBUTING.md**
- [ ] **ARCHITECTURE.md**
- [ ] **SECURITY.md**
- [ ] **FAQ.md**

---

## 🎯 Próximas Ações (Sprint Atual)

**Sprint**: 2025-10-20 a 2025-10-27  
**Foco**: FASE 0 + FASE 1

### Esta semana (até 23/10):
1. ✅ Criar ROADMAP.md
2. 🔄 Criar projeto retech-admin (Next.js)
3. 🔄 Sistema de Users (backend)
4. 🔄 Autenticação JWT
5. 🔄 Telas de login/registro

### Próxima semana (24-27/10):
1. Proteger endpoints GEO
2. Rate limiting
3. Logging de uso
4. Admin dashboard (básico)

---

## 🤝 Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md) (quando criado)

---

## 📞 Contato

- **Projeto**: Retech Core API
- **Maintainer**: The Retech Team
- **Status**: Em desenvolvimento ativo

---

**Legenda de Status**:
- 🔴 Não iniciado
- 🟡 Em andamento
- 🟢 Concluído
- ⏸️ Pausado
- ❌ Cancelado

**Legenda de Prioridade**:
- 🔴 CRÍTICA (bloqueante)
- 🟡 ALTA (importante)
- 🟢 MÉDIA (desejável)
- ⚪ BAIXA (nice to have)

---

---

## 📝 Como Usar Este Documento

### Este é o documento CENTRAL do projeto

- ✅ **Único documento de progresso** - Toda evolução é registrada aqui
- ✅ **Fonte única da verdade** - Status real do projeto
- ✅ **Atualizado continuamente** - Marque tarefas conforme completa
- ✅ **Histórico visível** - Git mostra evolução ao longo do tempo

### Como atualizar

1. **Ao completar uma tarefa**: Marque `[ ]` → `[x]`
2. **Ao iniciar uma fase**: Mude status para 🟡
3. **Ao completar uma fase**: Mude status para 🟢 e data
4. **Adicione no histórico** no topo do documento
5. **Atualize percentual** de progresso total

### Estrutura mantida

```
📅 Histórico de Atualizações (topo)
📊 Visão Geral do Progresso
🏗️ Arquitetura de Deployment
🎯 FASE 0, 1, 2, 3, 4, 5, 6 (com checkboxes)
📈 Métricas de Sucesso
🚀 Deploy e Infraestrutura
📚 Documentação
📞 Contato
📝 Como Usar Este Documento
```

### Outros documentos

- `README.md` - Introdução e overview
- `CHANGELOG.md` - Releases e versões
- `docs/INDEX.md` - Índice de toda documentação
- Guias específicos (Railway, Docker, etc)

---

**Última atualização**: 2025-10-20 23:30  
**Próxima revisão**: Ao completar FASE 2  
**Mantido por**: The Retech Team

