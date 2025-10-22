# 🗺️ ROADMAP UNIFICADO - Retech Core API
## O Hub Definitivo de APIs Brasileiras

**Última atualização**: 2025-10-22  
**Versão atual**: 1.3.0  
**Status geral**: 🟡 Produção + Correções Críticas Necessárias

---

## 🎯 VISÃO ESTRATÉGICA

Transformar a Retech Core na **API definitiva de dados brasileiros**:
- **31+ APIs essenciais** em uma única plataforma
- **Performance < 100ms** em todas as respostas
- **Gratuito para começar** com planos pagos opcionais
- **DX (Developer Experience) excepcional**
- **Rate limiting justo e transparente**
- **Documentação completa e atualizada**

---

## 📊 STATUS ATUAL (v1.3.0)

### ✅ **IMPLEMENTADO E FUNCIONANDO**
- 🗺️ **API Geográfica** (27 estados + 5.570 municípios)
- 🔐 **Autenticação JWT** (login, register, refresh)
- 👨‍💼 **Admin Dashboard** (gestão de tenants + API keys)
- 👨‍💻 **Developer Portal** (self-service de API keys)
- 🏦 **Sistema de API Keys** (criação, rotação, revogação)
- 📊 **Activity Logs** (auditoria de ações)
- 🎨 **Landing Page** (com roadmap público)
- 🚀 **Deploy Produção** (Railway + domínios configurados)

### 🐛 **BUGS CRÍTICOS IDENTIFICADOS**
- ❌ **Rate Limiting NÃO FUNCIONA** (erro de sintaxe + lógica)
- ⚠️ **Dashboard Dev com dados mocados** (não consume APIs reais)
- ⚠️ **Documentação desatualizada** (OpenAPI + Redoc)
- ⚠️ **`/me/stats` não implementado** (backend faltante)
- ⚠️ **RequestsPerMinute não verificado** (só diário)

### 📈 **PROGRESSO POR FASE**

```
FASE 0: Fundação              ██████████ 100% ✅
FASE 1: Auth & Segurança      ████████░░  85% 🟡 (rate limit com bug)
FASE 2: Admin Dashboard       ██████████ 100% ✅
FASE 3: Developer Portal      ██████░░░░  70% 🟡 (dados mocados)
FASE 4: Correções Críticas    ░░░░░░░░░░   0% 🔴 (próximo)
FASE 5: APIs Essenciais       ░░░░░░░░░░   0% 🔴
FASE 6: Expansão              ░░░░░░░░░░   0% 🔴
FASE 7: Governo & Compliance  ░░░░░░░░░░   0% 🔴
FASE 8: Analytics Avançado    ░░░░░░░░░░   0% 🔴
FASE 9: Monetização           ░░░░░░░░░░   0% 🔴
```

**Progresso total**: 72% das features base + 10 bugs/pendências críticas

---

## 🚨 **FASE 4: CORREÇÕES CRÍTICAS** (ESTA SEMANA!)
**Status**: 🔴 URGENTE  
**Timeline**: 22-26 Out 2025  
**Prioridade**: 🔴 BLOQUEANTE

### 4.1 Corrigir Rate Limiting (CRÍTICO)
- [ ] **Fix sintaxe** - Linha 80 de `rate_limiter.go` (falta `{`)
- [ ] **Mover verificação** - Checar ANTES de incrementar contador
- [ ] **Implementar limite por minuto** - Atualmente só verifica diário
- [ ] **Reset automático** - Limpar contadores antigos
- [ ] **Testar rigorosamente**:
  - 10 req/dia, 1 req/min
  - Verificar 429 na 11ª request
  - Verificar headers `X-RateLimit-*`

**Arquivos**:
- `internal/middleware/rate_limiter.go`
- `internal/domain/rate_limit.go`
- Testes: criar `rate_limiter_test.go`

**Estimativa**: 1-2 dias

---

### 4.2 Implementar `/me/stats` (Backend)
- [ ] Criar handler `GetMyStats` em `handlers/tenant.go`
- [ ] Agregações MongoDB:
  - Requests hoje (count por data)
  - Requests mês atual
  - Limite atual do tenant
  - % usado
- [ ] Retornar JSON:
  ```json
  {
    "requestsToday": 42,
    "requestsThisMonth": 1205,
    "limitDaily": 1000,
    "limitMonthly": 30000,
    "percentUsedDaily": 4.2,
    "percentUsedMonthly": 40.2
  }
  ```
- [ ] Adicionar rota `GET /me/stats` no router
- [ ] Testar com Postman

**Arquivos**:
- `internal/http/handlers/tenant.go`
- `internal/http/router.go`

**Estimativa**: 1 dia

---

### 4.3 Implementar `/me/usage` (Backend)
- [ ] Criar handler `GetMyUsage`
- [ ] Agregações:
  - Últimos 30 dias (array)
  - Por endpoint (top 10)
  - Por hora do dia (0-23)
- [ ] Cache de 5 minutos (opcional)
- [ ] Rota `GET /me/usage`

**Estimativa**: 1 dia

---

### 4.4 Dashboard Developer - Dados Reais (Frontend)
- [ ] Atualizar `app/painel/dashboard/page.tsx`
- [ ] Consumir `/me/stats`
- [ ] Gráfico recharts com `/me/usage`
- [ ] Loading states
- [ ] Error handling

**Estimativa**: 1 dia

---

### 4.5 Atualizar Documentação
- [ ] **OpenAPI (`openapi.yaml`)**:
  - Adicionar `/auth/*` endpoints
  - Adicionar `/admin/*` endpoints
  - Adicionar `/me/*` endpoints
  - Exemplos de requests/responses
  - Rate limit headers documentados
  
- [ ] **Redoc**:
  - Gerar HTML atualizado
  - Verificar `/docs` está acessível
  - Testar exemplos funcionam

- [ ] **Links**:
  - Landing page → `/docs`
  - Dashboard admin → "API Docs"
  - Dashboard painel → "Documentação"

**Estimativa**: 1-2 dias

---

### 4.6 Página de Docs no Painel
- [ ] Criar `app/painel/docs/page.tsx`
- [ ] Seções:
  - **Getting Started** (primeiros passos)
  - **Autenticação** (como usar API key)
  - **Endpoints** (lista com exemplos)
  - **Rate Limits** (explicação)
  - **Erros** (códigos e significados)
  - **SDKs** (futuro: links)
  - **FAQ**
- [ ] Syntax highlighting (prismjs ou shiki)
- [ ] Copyable code snippets
- [ ] Link para Redoc completo

**Estimativa**: 1 dia

---

**Total Fase 4**: ~6-8 dias de trabalho

---

## 📦 **FASE 5: APIs ESSENCIAIS** (Jan-Mar 2026)
**Status**: 🔴 Planejado  
**Timeline**: 3 meses (2-3 APIs/mês)  
**Prioridade**: 🔴 ALTA

### Dados Cadastrais
- [ ] **CEP** - Busca de endereços (ViaCEP + Correios)
- [ ] **CNPJ** - Consulta de empresas (ReceitaWS + RF)

### Dados Financeiros
- [ ] **Cotação de Moedas** - USD, EUR, BTC em tempo real
- [ ] **Bancos Brasileiros** - Lista com COMPE/ISPB

### Transporte
- [ ] **Tabela FIPE** - Preços de veículos

### Utilidades
- [ ] **Feriados Nacionais** - Calendário completo

**Entregáveis**: 6 APIs novas + testes + docs

---

## 📦 **FASE 6: EXPANSÃO** (Abr-Jun 2026)
**Status**: 🔴 Planejado  
**Timeline**: 3 meses (2 APIs/mês)  
**Prioridade**: 🟡 MÉDIA

### Dados Cadastrais (expansão)
- [ ] Validação de CPF
- [ ] Validação de Email
- [ ] Validação de Telefone

### Geografia (expansão)
- [ ] Bairros por Cidade
- [ ] Coordenadas de CEPs
- [ ] Dados Demográficos IBGE

### Finanças (expansão)
- [ ] SELIC, CDI, IPCA

### Transporte (expansão)
- [ ] Cálculo de Frete
- [ ] Rastreamento Correios
- [ ] Consulta de Veículos (DENATRAN)

### Utilidades (expansão)
- [ ] Operadora de Telefone
- [ ] Dias Úteis
- [ ] Fusos Horários

**Entregáveis**: 12 APIs novas

---

## 📦 **FASE 7: GOVERNO & COMPLIANCE** (Jul-Set 2026)
**Status**: 🔴 Planejado  
**Timeline**: 3 meses (1-2 APIs/mês)  
**Prioridade**: 🟢 BAIXA

- [ ] Dados Judiciais (PJe + TJs)
- [ ] Portal da Transparência
- [ ] CEIS/CNEP
- [ ] Simples Nacional
- [ ] PEP

**Entregáveis**: 5 APIs novas

---

## 📦 **FASE 8: ANALYTICS AVANÇADO** (Out-Dez 2026)
**Status**: 🔴 Planejado  
**Timeline**: 3 meses  
**Prioridade**: 🟢 BAIXA

### Backend
- [ ] Agregações otimizadas (cache Redis)
- [ ] Geolocalização por IP
- [ ] Alertas configuráveis
- [ ] Webhooks de eventos
- [ ] Exportar relatórios (CSV/PDF)

### Frontend
- [ ] Dashboards customizáveis
- [ ] Filtros avançados
- [ ] Comparação de períodos
- [ ] Heatmaps de uso
- [ ] Logs em tempo real (WebSocket)

---

## 📦 **FASE 9: MONETIZAÇÃO** (2027+)
**Status**: 🔴 Futuro  
**Prioridade**: 🟢 BAIXA

### Planos Pagos
- [ ] Model `Plan` (Free, Pro, Business)
- [ ] Integração Stripe
- [ ] Checkout de planos
- [ ] Upgrade/downgrade
- [ ] Invoices automáticos

### Features Premium
- [ ] SLA garantido
- [ ] Webhooks customizados
- [ ] White label
- [ ] Suporte prioritário

---

## 📋 **RESUMO DAS 31 APIs PLANEJADAS**

| Categoria | Fase 5 | Fase 6 | Fase 7 | Fase 8 | **Total** |
|-----------|--------|--------|--------|--------|-----------|
| 📋 Dados Cadastrais | 2 | 5 | 0 | 2 | **9** |
| 🗺️ Geografia | 0 | 3 | 0 | 2 | **5** |
| 💰 Finanças | 2 | 1 | 0 | 2 | **5** |
| 🚚 Transporte | 1 | 3 | 0 | 0 | **4** |
| 🔧 Utilidades | 1 | 3 | 0 | 0 | **4** |
| 🏛️ Governo | 0 | 0 | 5 | 0 | **5** |
| **TOTAL** | **6** | **15** | **5** | **6** | **32** |

*Nota: +1 API já disponível (Geografia)*

---

## 🏗️ ARQUITETURA TÉCNICA

### Stack Atual
- **Backend**: Go 1.24 + Gin + MongoDB
- **Frontend**: Next.js 14 + Shadcn/ui + Tailwind
- **Deploy**: Railway (auto-deploy)
- **Domínios**:
  - `core.theretech.com.br` → Frontend
  - `api-core.theretech.com.br` → Backend

### Infraestrutura Planejada (Fase 5+)
- [ ] **Redis** - Cache de APIs externas + rate limiting por minuto
- [ ] **Bull Queue** - Jobs assíncronos (scraping, agregações)
- [ ] **Circuit Breaker** - Proteção contra falhas de APIs externas
- [ ] **CDN** - Cloudflare para assets estáticos
- [ ] **Sentry** - Error tracking
- [ ] **Uptime Robot** - Monitoring 24/7

---

## 🎯 MÉTRICAS DE SUCESSO

### Fase 4 (Esta Semana)
- [ ] Rate limiting 100% funcional
- [ ] Dashboard dev com dados reais
- [ ] Documentação atualizada
- [ ] 0 bugs críticos

### Fase 5 (3 meses)
- [ ] 6 novas APIs funcionando
- [ ] 500 desenvolvedores cadastrados
- [ ] 100K requests/dia
- [ ] Latência < 150ms
- [ ] Uptime > 99.5%

### Fase 6 (6 meses)
- [ ] 15+ APIs totais
- [ ] 2.000 desenvolvedores
- [ ] 500K requests/dia
- [ ] Latência < 100ms
- [ ] Uptime > 99.9%

### Fase 7-9 (12 meses)
- [ ] 30+ APIs totais
- [ ] 10.000 desenvolvedores
- [ ] 5M requests/dia
- [ ] 200 clientes pagos
- [ ] NPS > 50

---

## 📈 DIFERENCIAIS COMPETITIVOS

### vs Sites Governamentais (IBGE, Receita, etc.)
✅ **100x mais rápido** (<100ms vs segundos)  
✅ **API moderna** (REST + JSON)  
✅ **Disponibilidade 99.9%** (vs downtime frequente)  
✅ **Rate limit justo** (1000 req/dia grátis)  
✅ **Documentação completa** (vs inexistente)

### vs APIs separadas (ViaCEP, ReceitaWS, etc.)
✅ **Tudo em uma API** (não precisa integrar 10+)  
✅ **Dashboard unificado** (um lugar para tudo)  
✅ **Uma API key** (simplifica segurança)  
✅ **Suporte técnico** (vs comunidade)  
✅ **SLA garantido** (planos pagos)

### vs Soluções Internacionais (Stripe, Plaid)
✅ **Foco no Brasil** (dados locais únicos)  
✅ **Preço acessível** (R$ vs USD)  
✅ **LGPD compliance** (desde o início)  
✅ **Suporte em PT-BR** (horário comercial BR)  
✅ **Pagamento em Real** (sem IOF)

---

## 🚀 PRÓXIMAS AÇÕES (Esta Semana)

### Segunda (22 Out)
- [x] ✅ Criar `PENDENCIAS_CRITICAS.md`
- [x] ✅ Criar `ROADMAP_UNIFIED.md`
- [ ] 🔄 Corrigir bug do rate limiting

### Terça (23 Out)
- [ ] Implementar `/me/stats` (backend)
- [ ] Testar com Postman
- [ ] Atualizar collection

### Quarta (24 Out)
- [ ] Implementar `/me/usage` (backend)
- [ ] Dashboard dev dados reais (frontend)
- [ ] Testar integração

### Quinta (25 Out)
- [ ] Atualizar `openapi.yaml`
- [ ] Gerar Redoc novo
- [ ] Verificar links

### Sexta (26 Out)
- [ ] Criar página `/painel/docs`
- [ ] Testes finais
- [ ] Deploy de tudo

---

## 📚 DOCUMENTAÇÃO

### ✅ Criados
- [x] README.md
- [x] ROADMAP.md (antigo)
- [x] ROADMAP_V2.md
- [x] ROADMAP_UNIFIED.md (este)
- [x] PENDENCIAS_CRITICAS.md
- [x] RAILWAY_DEPLOY.md
- [x] APIKEY_ROTATE.md

### ⏳ Próximos
- [ ] API_REFERENCE.md (completo)
- [ ] INTEGRATION_GUIDE.md
- [ ] BEST_PRACTICES.md
- [ ] FAQ.md
- [ ] SECURITY.md

---

## 🤝 CONTRIBUINDO

Tem alguma sugestão de API essencial? Entre em contato!

**Email**: contato@theretech.com.br  
**Website**: https://core.theretech.com.br  
**Docs**: https://api-core.theretech.com.br/docs  

---

## 📞 SUPORTE

- **Bugs críticos**: Criar issue no GitHub
- **Dúvidas técnicas**: contato@theretech.com.br
- **Feature requests**: Via landing page ou email

---

**Legenda de Status**:
- 🟢 Completo
- 🟡 Parcial/Com bugs
- 🔴 Não iniciado
- ❌ Bloqueante

**Legenda de Prioridade**:
- 🔴 CRÍTICA (bloqueante)
- 🟡 ALTA (importante)
- 🟢 MÉDIA (desejável)
- ⚪ BAIXA (futuro)

---

**Última atualização**: 2025-10-22 23:00  
**Próxima revisão**: Ao completar Fase 4  
**Mantido por**: The Retech Team  
**Versão**: 1.0 (unificado)

