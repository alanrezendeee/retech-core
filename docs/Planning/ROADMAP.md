# 🚀 ROADMAP RETECH CORE API

**Atualizado:** 31 de Outubro de 2025 🆕  
**Status:** Fase 1 Concluída ✅ | Fase 2 Em Andamento 🔵 (2/6 APIs - 33%) | Infraestrutura Avançada ✅

---

## 🔍 **ÍNDICE DE TAGS - BUSCA RÁPIDA**

**Ao implementar algo novo, busque por:**

### 🏗️ **Nova API Completa:**
- `#NOVA-API` - Checklist completo + backend + frontend
- `#CHECKLIST-COMPLETO` - Checklist detalhado de implementação (linha 1975)
- `#LEITURA-OBRIGATÓRIA` - Padrões arquiteturais obrigatórios (linha 8)

### 🔧 **Nova Funcionalidade/Endpoint:**
- `#FUNCIONALIDADE` - Checklist simplificado (linha 1657)
- `#ENDPOINT` - Como adicionar endpoint em API existente
- `#CHECKLIST-SIMPLIFICADO` - Checklist rápido para funcionalidades

### 📚 **Referência e Boas Práticas:**
- `#PADRÕES` - Padrões arquiteturais (linha 1598)
- `#LIÇÕES` - Lições aprendidas (linha 1454)
- `#DOCUMENTAÇÃO` - Como documentar no OpenAPI (linha 1936)
- `#REFERÊNCIA` - Documentação de referência

### 🎯 **Importante:**
- `#IMPORTANTE` - Critérios e decisões importantes (linha 2438)
- `#DECISÃO` - Como decidir entre Nova API vs Funcionalidade

### 🛠️ **Por Tipo:**
- `#BACKEND` - Checklist backend (linha 1571)
- `#FRONTEND` - Checklist frontend (linha 1975)
- `#OPENAPI` - Documentação Redoc

---

**🔎 Como usar:** `Ctrl+F` ou `Cmd+F` e busque pela tag (ex: `#NOVA-API`)

---

## 🔧 **PADRÃO ARQUITETURAL** 🆕
### 📖 `#LEITURA-OBRIGATÓRIA` `#NOVA-API` `#FUNCIONALIDADE`

> **⚠️ LEIA ANTES DE IMPLEMENTAR QUALQUER COISA!**

**REGRA:** ❌ **NUNCA hardcode URLs de APIs externas!**

**Usar:** ✅ Variáveis de ambiente (ENV) obrigatórias

**Template:** Ver seção "Checklist Nova API" (linha 2079) para implementação completa.

---

## ⚡ **RESUMO EXECUTIVO - OUT/2025**

### **🎉 Principais Conquistas:**
1. ✅ **3 APIs Completas** - CEP, CNPJ, Geografia (4 endpoints)
2. ✅ **Cache 3 Camadas** - Redis L1 (~1ms) + MongoDB L2 (~10ms) + API L3 (~200ms)
3. ✅ **Playground Público** - Teste sem cadastro + Browser fingerprinting
4. ✅ **Segurança Reforçada** - HMAC-SHA256 + Scopes + Rate limiting multi-camada
5. ✅ **SEO Completo** - Meta tags + Sitemap + Ferramentas públicas
6. ✅ **Admin Dashboard** - Controles independentes + Stats + Analytics (timezone BR)

### **🚀 Próximo Grande Passo:**
**Migração Oracle Cloud (São Paulo)** - Reduzir latência de 160ms → 5-15ms  
**Custo:** R$ 0/mês (Always Free Tier) | **Prazo:** 2-4 semanas

### **💰 Performance vs Custo:**
- **Atual:** Railway EUA (~160ms) - $5-10/mês
- **Futuro:** Oracle BR (~5-15ms) - **R$ 0/mês** ✅
- **Ganho:** 10-30x mais rápido + Gratuito!

---

## 📊 PROGRESSO GERAL

```
██████░░░░░░░░░░░░░░░ 27% (10/37 APIs)

Fase 1: ████████████ 100% ✅
Fase 2: █████░░░░░░░  43% 🔵 (3/7)
Fase 3: ░░░░░░░░░░░░   0% ⚪ (0/17)
Fase 4: ░░░░░░░░░░░░   0% ⚪ (0/7)
```

**APIs Totais:** 37 (+6 novas: NF-e, CND, Compras Gov, Protestos, Score, Artigos Penais) 🆕  
**Disponíveis:** 4 (Geografia + CEP + CNPJ + Artigos Penais) 🚀  
**Em Desenvolvimento:** 4 (Fase 2)  
**Planejadas:** 29

---

## ✅ FASE 1 - FUNDAÇÃO (CONCLUÍDA)

### **🎯 Infraestrutura Core**
- [x] Arquitetura Go + Gin + MongoDB
- [x] Autenticação JWT (SUPER_ADMIN + TENANT_USER)
- [x] API Key Management (criação, revogação, rotação)
- [x] Rate Limiting por tenant (daily + per-minute)
- [x] Usage Logging (API usage tracking)
- [x] Activity Logs (auditoria completa)
- [x] Maintenance Mode
- [x] Docker + Docker Compose
- [x] Deploy Railway (backend + frontend)

### **🎨 Frontend (Admin + Developer Portal)**
- [x] Admin Dashboard (gestão de tenants, API keys, settings)
- [x] Developer Portal (dashboard, usage, docs)
- [x] Landing Page com 31 APIs documentadas
- [x] Sistema de autenticação completo
- [x] Real-time metrics (auto-refresh 30s)

### **📚 Documentação**
- [x] OpenAPI 3.0 (Redoc)
- [x] URL dinâmica (dev/prod)
- [x] Documentação focada em desenvolvedores
- [x] Exemplos de código funcionais

### **🗺️ API: Geografia (DISPONÍVEL)**
- [x] `GET /geo/ufs` - Lista 27 estados
- [x] `GET /geo/ufs/:sigla` - Detalhes de estado
- [x] `GET /geo/municipios` - Lista 5.570 municípios
- [x] `GET /geo/municipios/:uf` - Municípios por UF
- [x] `GET /public/geo/*` - Endpoints públicos para playground
- [x] Dados do IBGE (completos)
- [x] **Cache Redis L1:**
  - [x] Estados (ufs:all) - cache permanente
  - [x] Municípios por UF - cache permanente
- [x] Performance: ~1ms (Redis) / ~160ms (primeira consulta)
- [x] Scope: `geo` (controle granular de acesso)

### **⚡ Performance & Cache (IMPLEMENTADO)** ✅ 🆕
- [x] **Redis Cache (L1 - Hot Cache):**
  - [x] Conexão Redis configurável via `REDIS_URL`
  - [x] Cache em memória (~1ms)
  - [x] Keys separadas por serviço (cep:*, cnpj:*, geo:*)
  - [x] TTL: 24h para hot data
  - [x] Graceful degradation (se cair, usa MongoDB)
  - [x] Monitoramento: stats de keys, memória usada
  - [x] Admin: Limpar cache Redis (all, cep, cnpj)
  
- [x] **MongoDB Cache (L2 - Persistent Cache):**
  - [x] Cache persistente (~10ms)
  - [x] TTL independente por serviço (CEP: 7 dias, CNPJ: 30 dias)
  - [x] Auto-cleanup via TTL Index
  - [x] Promoção automática para Redis quando hit
  - [x] Stats: total cached, recent 24h
  - [x] Admin: Controles independentes (CEP e CNPJ)

- [x] **MongoDB Indexes (Auto-criação):**
  - [x] API Keys: keyId (unique), ownerId, scopes
  - [x] Tenants: email (unique), active
  - [x] Users: email (unique), tenantId
  - [x] CEP Cache: cep (unique), cachedAt (TTL)
  - [x] CNPJ Cache: cnpj (unique), cachedAt (TTL)
  - [x] API Usage: date+apiKey, tenantId+date, endpoint+date
  - [x] Activity Logs: userId, createdAt (desc)
  - [x] Playground Rate Limits: ip+apiKey+date (unique)
  - [x] Indexes aplicados automaticamente no startup (migrations-like)

- [x] **Settings Cache (In-Memory):**
  - [x] Cache de configurações do sistema em memória
  - [x] Evita consultas frequentes ao MongoDB
  - [x] Invalidação automática ao salvar settings

### **🔒 Segurança Avançada (IMPLEMENTADO)** ✅ 🆕
- [x] **API Key Management:**
  - [x] Formato seguro: `keyId.keySecret`
  - [x] Hash HMAC-SHA256 com salt (`APIKEY_HASH_SECRET`)
  - [x] KeyHash armazenado (nunca a chave em texto)
  - [x] Validação obrigatória de `APIKEY_HASH_SECRET`
  - [x] Panic se variável não configurada (segurança)
  - [x] API Key oculta no frontend (🔒 •••••••)
  - [x] Criação, rotação e revogação de chaves

- [x] **Sistema de Scopes:**
  - [x] Scopes granulares: `cep`, `cnpj`, `geo`, `all`
  - [x] Middleware `RequireScope` em todas as rotas públicas
  - [x] Validação automática de permissões
  - [x] Retrocompatibilidade (aceita `geo:read` e `geo`)
  - [x] Frontend com seleção visual (checkboxes)

- [x] **Rate Limiting Multi-Camada:**
  - [x] **Por Tenant:** Daily + Per-minute (configurável)
  - [x] **Por IP (Playground):** 4/dia, 1/min (configurável)
  - [x] **Global (Playground):** Compartilhado entre todos os IPs
  - [x] Reset automático por minuto
  - [x] Logs detalhados de debug
  - [x] MongoDB indexes otimizados

- [x] **Playground Seguro:**
  - [x] API Key Demo gerenciável via admin/settings
  - [x] Scopes auto-rotacionados ao mudar checkboxes
  - [x] Rate limiting dedicado (IP + Global)
  - [x] Browser fingerprinting (coleta sem travar UI)
  - [x] Throttling (delay mínimo entre requests)
  - [x] Toggle enable/disable via admin
  - [x] Configuração de APIs permitidas

- [x] **CORS Dinâmico:**
  - [x] Enable/disable via admin/settings
  - [x] Origins configuráveis (textarea)
  - [x] Strict mode (sem exceções para localhost)
  - [x] Headers personalizados permitidos
  - [x] Respostas 204/200 para OPTIONS

- [x] **JWT Dinâmico:**
  - [x] Access TTL configurável (padrão: 15min)
  - [x] Refresh TTL configurável (padrão: 7 dias)
  - [x] Atualização em tempo real via admin/settings
  - [x] Roles: SUPER_ADMIN, TENANT_USER

### **📊 Analytics & Monitoring (IMPLEMENTADO)** ✅ 🆕
- [x] **Dashboard Admin:**
  - [x] Stats globais (tenants, API keys, users, requests)
  - [x] Gráficos de uso diário (últimos 30 dias)
  - [x] Métricas em tempo real (auto-refresh 30s)
  - [x] Top endpoints mais usados
  - [x] Timezone Brasília (America/Sao_Paulo)
  - [x] Formatação de datas pt-BR

- [x] **Usage Tracking:**
  - [x] Log de todas as requests (endpoint, tenant, timestamp)
  - [x] Agregação por dia, tenant, endpoint
  - [x] Rate limit tracking (daily + per-minute)
  - [x] MongoDB indexes para queries rápidas

- [x] **Activity Logs:**
  - [x] Auditoria completa de ações do sistema
  - [x] Login, criação de API Keys, updates de settings
  - [x] UserID, email, role, timestamp

- [x] **Redis Monitoring:**
  - [x] Total keys, keys por serviço (CEP, CNPJ, GEO)
  - [x] Memória usada (MB)
  - [x] Status de conexão (conectado/desconectado)
  - [x] Admin dashboard com stats em tempo real

---

## 🔵 FASE 2 - DADOS ESSENCIAIS (0-3 MESES)

**Meta:** 7 APIs | **Status:** 3/7 Concluídas (43%)

### **📮 CEP (DISPONÍVEL)** ✅
- [x] `GET /cep/:codigo` - Busca por CEP (CEP → Endereço)
- [x] `GET /cep/buscar` - Busca reversa (Endereço → CEP) 🆕
- [x] `GET /public/cep/:codigo` - Endpoint público para playground/ferramentas
- [x] `GET /public/cep/buscar` - Endpoint público para busca reversa 🆕
- [x] Integração: ViaCEP (gratuito)
- [x] Fallback: Brasil API
- [x] **Cache 3 Camadas:**
  - [x] Redis L1 (~1ms) - Hot cache em memória
  - [x] MongoDB L2 (~10ms) - Cache persistente
  - [x] API Externa L3 (~200ms) - ViaCEP/Brasil API
- [x] **Busca Reversa:** 🆕
  - [x] Parâmetros: uf, cidade, logradouro (mín. 3 caracteres)
  - [x] Retorna até 50 CEPs por busca
  - [x] Cache independente (search:uf:cidade:logradouro)
  - [x] Ferramenta pública: `/ferramentas/buscar-cep` 🆕
- [x] TTL configurável: 1-365 dias (padrão: 7 dias)
- [x] Coordenadas geográficas
- [x] Normalização automática (com/sem traço)
- [x] Performance: ~1ms (Redis) / ~10ms (MongoDB) / ~160ms (API)
- [x] **Admin Settings:**
  - [x] Toggle independente CEP (enable/disable)
  - [x] TTL dinâmico (1-365 dias)
  - [x] Auto-cleanup (MongoDB TTL Index)
  - [x] Stats em tempo real (total cached, recent 24h)
  - [x] Limpeza manual (botão destrutivo com confirmação)
- [x] Scope: `cep` (controle granular de acesso)
- [x] Graceful degradation (Redis cai → MongoDB funciona)

### **🏢 CNPJ (DISPONÍVEL)** ✅
- [x] `GET /cnpj/:numero` - Consulta CNPJ
- [x] `GET /public/cnpj/:numero` - Endpoint público para playground/ferramentas
- [x] Fonte: Brasil API (gratuita, Receita Federal)
- [x] Fallback: ReceitaWS
- [x] **Cache 3 Camadas:**
  - [x] Redis L1 (~1ms) - Hot cache em memória
  - [x] MongoDB L2 (~10ms) - Cache persistente
  - [x] API Externa L3 (~200ms) - Brasil API
- [x] TTL configurável: 1-365 dias (padrão: 30 dias)
- [x] Validação: Dígito verificador + normalização
- [x] Dados: razão social, nome fantasia, situação
- [x] QSA: Quadro de sócios e administradores
- [x] CNAEs: Atividade principal + secundárias
- [x] Endereço completo + contatos
- [x] Performance: ~1ms (Redis) / ~10ms (MongoDB) / ~200ms (Brasil API)
- [x] **Admin Settings:**
  - [x] Toggle independente CNPJ (enable/disable)
  - [x] TTL dinâmico (1-365 dias)
  - [x] Auto-cleanup (MongoDB TTL Index)
  - [x] Stats em tempo real (total cached, recent 24h)
  - [x] Limpeza manual (botão destrutivo com confirmação)
- [x] Scope: `cnpj` (controle granular de acesso)
- [x] Graceful degradation (Redis cai → MongoDB funciona)

### **💵 Moedas**
- [ ] `GET /moedas/cotacao` - Cotações em tempo real
- [ ] Fonte: Banco Central API (gratuita)
- [ ] Moedas: USD, EUR, BTC
- [ ] Cache: 1 hora
- [ ] Histórico: últimos 30 dias

### **🏦 Bancos**
- [ ] `GET /bancos` - Lista bancos brasileiros
- [ ] `GET /bancos/:codigo` - Busca por código
- [ ] Fonte: Banco Central (dados públicos STR)
- [ ] Dados: COMPE, ISPB, nome completo
- [ ] Cache local permanente (atualização mensal)

### **🚗 FIPE**
- [ ] `GET /fipe/marcas` - Marcas de veículos
- [ ] `GET /fipe/veiculos/:codigo` - Preço FIPE
- [ ] Fonte: FIPE API (gratuita via Denatran)
- [ ] Cache: 7 dias
- [ ] Filtros: marca, modelo, ano

### **📅 Feriados**
- [ ] `GET /feriados/:ano` - Feriados nacionais
- [ ] `GET /feriados/:uf/:ano` - Feriados estaduais
- [ ] Fonte: Arquivo local + leis federais
- [ ] Cache permanente (gerado por ano)
- [ ] Tipos: nacional, estadual, municipal, ponto facultativo

### **⚖️ Artigos Penais (DISPONÍVEL)** ✅ 🆕
- [x] `GET /penal/artigos` - Lista artigos penais (autocomplete)
- [x] `GET /penal/artigos/:codigo` - Busca artigo específico
- [x] `GET /penal/search` - Busca por texto (descrição)
- [x] `GET /public/penal/*` - Endpoints públicos para playground
- [x] Fonte: Dados fixos (Código Penal + Lei de Contravenções)
- [x] **Cache 3 Camadas:**
  - [x] Redis L1 (~1ms) - Hot cache em memória (365 dias)
  - [x] MongoDB L2 (~10ms) - Cache persistente (permanente)
  - [x] Dados fixos (não muda)
- [x] Estrutura hierárquica: Artigo → Parágrafo → Inciso → Alínea
- [x] Tipos: Crime, Contravenção
- [x] Legislações: CP (Código Penal), LCP (Lei de Contravenções)
- [x] **Seed inicial:** `seeds/penal.json` com artigos mais comuns
- [x] Performance: ~1ms (Redis) / ~10ms (MongoDB)
- [x] **Uso:** Autocomplete/Select2 para seleção de crimes
- [x] Scope: `penal` (controle granular de acesso)
- [x] Índices MongoDB otimizados (codigo, busca, tipo, legislacao)

### **📱 Telefone (PLANEJADO)** 🆕

#### **Objetivo:**
Validação de telefones brasileiros com diferencial único: **WhatsApp Verification real** via Evolution API auto-hospedada (custo ZERO).

#### **Features Planejadas:**

**1. WhatsApp Validator** ✅ (Diferencial Competitivo)
- [ ] `GET /phone/:numero` - Validação completa de telefone
- [ ] **WhatsApp Verification:** Consulta REAL na rede WhatsApp (via Evolution API)
  - Custo: R$ 0 (Evolution auto-hospedada)
  - Confiabilidade: 100% (verificação real, não heurística)
  - Retorna: `{ exists: true/false, jid: "number@s.whatsapp.net" }`
- [ ] **Validação de Formato:** Regras ANATEL (95-98% preciso)
  - 11 dígitos → móvel (9º dígito obrigatório)
  - 10 dígitos → fixo (primeiro dígito 2-5)
  - DDDs válidos (11-99, exceto inválidos)
- [ ] **Tipo:** Mobile ou Landline (99%+ preciso)
  - Baseado em 9º dígito (Lei 12.249/2010)
  - Sem exceções conhecidas
- [ ] **DDD → Localização:** Estado e cidades possíveis (100% preciso)
  - Integração: BrasilAPI (`GET /ddd/v1/:ddd`)
  - Fallback: Tabela local
  - Cache: Permanente (DDDs não mudam)
- [ ] **Cache 3 Camadas:**
  - Redis L1 (~1ms)
  - MongoDB L2 (~10ms)
  - Evolution API L3 (~100-500ms)
- [ ] **Scope:** `phone`

**Response Exemplo:**
```json
{
  "numero": "5548988612609",
  "valido": true,
  "tipo": "mobile",
  "ddd": "48",
  "estado": "SC",
  "cidades_possiveis": ["Florianópolis", "São José", "Palhoça"],
  "whatsapp": {
    "existe": true,
    "jid": "5548988612609@s.whatsapp.net",
    "verificado_em": "2025-10-28T22:00:00Z"
  },
  "observacoes": {
    "formato": "Validado segundo regras ANATEL",
    "tipo": "Baseado em 9º dígito obrigatório",
    "localizacao": "DDD pode abranger múltiplas cidades",
    "whatsapp": "Verificação real na rede WhatsApp"
  }
}
```

**2. WhatsApp OTP** 🔥 (Inovação - Custo Zero)
- [ ] `POST /phone/otp/send` - Enviar código OTP via WhatsApp
- [ ] `POST /phone/otp/verify` - Validar código OTP

**Fluxo WhatsApp OTP:**
```
┌─────────────┐
│   Dev App   │
└──────┬──────┘
       │
       │ POST /phone/otp/send
       │ {
       │   "numero": "5548988612609",
       │   "ttl": 300,          // Segundos (opcional, padrão: 300)
       │   "digits": 6,         // Tamanho código (opcional, padrão: 6)
       │   "template": "custom" // Template (opcional, padrão: "default")
       │ }
       ▼
┌─────────────────────┐
│  Retech Core API    │
│                     │
│ 1. Valida formato   │
│ 2. Checa WhatsApp   │ ← Evolution API (verificar se número existe)
│ 3. Verifica quota   │ ← Limites por plano (100/1k/10k OTPs/mês)
│ 4. Rate limit       │ ← Máx 3 OTPs/5min por número (anti-spam)
│ 5. Gera OTP         │ ← 4-8 dígitos aleatório
│ 6. Salva Redis      │
│    Key: phone:otp:{numero}
│    TTL: Configurável (padrão: 5 min)
│    Data: {
│      code: "123456",
│      used: false,
│      attempts: 0,
│      tenant: "tenant_id",
│      created_at: timestamp
│    }
│ 7. Envia WhatsApp   │ ← Evolution API (custo R$ 0!)
│    Template:
│    "🔐 Seu código {APP_NAME}:\n\n*{OTP}*\n\nVálido por {TTL} minutos."
└──────┬──────────────┘
       │
       │ Response:
       │ {
       │   "enviado": true,
       │   "numero": "5548988612609",
       │   "metodo": "whatsapp",
       │   "expira_em": "2025-10-28T23:05:00Z",
       │   "tentativas_restantes": 3
       │ }
       ▼
┌─────────────────────┐
│   WhatsApp User     │
│                     │
│ 📱 Recebe mensagem: │
│                     │
│ 🔐 Seu código       │
│ MeuApp:             │
│                     │
│ *123456*            │
│                     │
│ Válido por 5        │
│ minutos.            │
└──────┬──────────────┘
       │
       │ Usuário digita código no app
       ▼
┌─────────────────────┐
│   Dev App           │
└──────┬──────────────┘
       │
       │ POST /phone/otp/verify
       │ {
       │   "numero": "5548988612609",
       │   "code": "123456"
       │ }
       ▼
┌─────────────────────┐
│  Retech Core API    │
│                     │
│ 1. Busca Redis      │ ← GET phone:otp:{numero}
│ 2. Valida código    │ ← Compara code
│ 3. Checa expirado   │ ← TTL Redis
│ 4. Checa usado      │ ← used == false
│ 5. Incrementa       │ ← attempts++ (máx 5)
│    tentativas       │
│ 6. Se correto:      │
│    - Marca usado    │ ← used = true
│    - Deleta Redis   │ ← DEL phone:otp:{numero}
│    - Log sucesso    │
│    - Webhook (opt)  │ ← POST {dev_webhook_url}
│ 7. Se errado:       │
│    - Retorna erro   │
│    - Mantém OTP     │
│ 8. Response         │
└─────────────────────┘
       │
       │ Response (sucesso):
       │ {
       │   "valido": true,
       │   "numero": "5548988612609",
       │   "verificado_em": "2025-10-28T22:45:00Z"
       │ }
       │
       │ Response (erro):
       │ {
       │   "valido": false,
       │   "erro": "Código incorreto",
       │   "tentativas_restantes": 2
       │ }
```

**Features WhatsApp OTP:**
- [x] **Custo ZERO** (Evolution API auto-hospedada)
- [x] **Taxa de abertura 98%** (vs 20% SMS)
- [x] **Customização:** Templates configuráveis por tenant
- [x] **Segurança:**
  - Rate limit: 3 OTPs/5min por número (anti-spam)
  - Rate limit: Quota mensal por plano (100/1k/10k)
  - Máx 5 tentativas de verificação por OTP
  - Código expira (TTL configurável: 5-30 min)
  - Marca como usado (não reutilizável)
- [x] **Webhook:** Notificação quando OTP validado (opcional)
- [x] **Logs:** Auditoria completa (envio, tentativas, verificação)

**Configurações (Painel do Dev):**
```json
{
  "otp": {
    "ttl": 300,              // Segundos (5 min padrão)
    "digits": 6,             // Tamanho do código
    "max_attempts": 5,       // Tentativas de verificação
    "rate_limit_window": 300, // Janela rate limit (5 min)
    "rate_limit_max": 3,     // Máx OTPs na janela
    "template": "default",   // ou "custom"
    "custom_template": "Seu código é: {OTP}",
    "webhook_url": "https://seuapp.com/otp/verified", // opcional
    "app_name": "Meu App"    // Nome no template
  }
}
```

**Quotas por Plano:**
```
FREE:     100 OTPs/mês
BASIC:    1.000 OTPs/mês  (R$ 29/mês)
PRO:      10.000 OTPs/mês (R$ 99/mês)
BUSINESS: Ilimitado       (R$ 299/mês)
```

**Endpoints Completos:**
```
GET  /phone/:numero           - Validar + WhatsApp check
POST /phone/otp/send          - Enviar OTP via WhatsApp
POST /phone/otp/verify        - Verificar código OTP
GET  /phone/otp/status/:numero - Status do OTP (dev only)
```

#### **💡 Análise do Fluxo (Opinião Técnica):**

**✅ PONTOS FORTES:**
- Fluxo simples e direto (dev-friendly)
- Expiração configurável (flexível)
- Validação de uso único (segurança)
- WhatsApp (alta taxa de abertura)
- Custo ZERO (Evolution própria)

**⚠️ MELHORIAS SUGERIDAS:**

1. **Rate Limiting Duplo:**
   - Por número: 3 OTPs/5min (evita spam ao usuário)
   - Por tenant: Quota mensal (evita abuso do serviço)

2. **Tentativas Limitadas:**
   - Máx 5 tentativas de verificação por OTP
   - Após 5 falhas, bloquear e exigir novo OTP

3. **Webhook de Confirmação:**
   - Dev pode receber POST quando OTP validado
   - Payload: `{ numero, verificado_em, tenant_id }`
   - Evita polling constante

4. **Templates Customizáveis:**
   - Variáveis: `{APP_NAME}`, `{OTP}`, `{TTL}`
   - Exemplo: "Seu código {APP_NAME} é: {OTP}"
   - Configurável no painel do dev

5. **Múltiplos Tamanhos de OTP:**
   - Configurável: 4, 6, 8 dígitos
   - Padrão: 6 dígitos
   - Ajustável por nível de segurança

6. **Logs de Auditoria:**
   - Registrar envio, tentativas, verificação
   - Dashboard: quantos OTPs enviados/verificados
   - Alertas: quota próxima do limite

**🚨 RISCOS E MITIGAÇÕES:**

**Risco 1: Banimento WhatsApp**
- Problema: WhatsApp pode banir número por spam
- Limite: ~1.000 msgs/dia por número
- Solução:
  - Usar múltiplas instâncias Evolution (rotação)
  - Rate limit de 500 OTPs/dia por instância
  - Monitoramento de health (QR Code válido?)
  - Alertas de desconexão

**Risco 2: Confiabilidade Evolution**
- Problema: Evolution depende de conexão WhatsApp estável
- Solução:
  - Health check a cada 5 min
  - Reconectar automaticamente se cair
  - Fallback opcional para SMS (se dev configurar gateway próprio)

**Risco 3: LGPD/Compliance**
- Problema: WhatsApp Business Terms + LGPD
- Solução:
  - Opt-in obrigatório (documentar no cadastro)
  - Permitir opt-out
  - Logs de consentimento
  - Não enviar marketing (só OTP)

#### **🔧 Implementação Técnica:**

**Arquivos Principais:**
```
Backend:
- internal/http/handlers/phone.go          (handler principal)
- internal/services/evolution_client.go    (client Evolution API)
- internal/services/otp_service.go         (lógica OTP)
- internal/http/router.go                  (rotas)
- internal/domain/settings.go              (config OTP)
- internal/auth/scope_middleware.go        (scope "phone")

Frontend:
- app/ferramentas/validar-telefone/page.tsx (ferramenta pública)
- app/painel/docs/page.tsx                  (docs dev)
- components/otp-config.tsx                 (config painel dev)

Docs:
- internal/docs/openapi.yaml               (Redoc)
```

**Dependencies:**
```go
// Evolution API Client
type EvolutionClient struct {
    BaseURL  string
    APIKey   string
    Instance string
}

// OTP Service
type OTPService struct {
    Redis    *redis.Client
    Evolution *EvolutionClient
    Config   OTPConfig
}
```

**Tempo Estimado:**
- WhatsApp Validator: 3-4 horas
- WhatsApp OTP: 5-6 horas
- Testes + Docs: 2-3 horas
- **Total: 10-13 horas** (~2 dias)

#### **📊 Diferencial Competitivo:**

**O que concorrentes oferecem:**
- Twilio: SMS ($0,08/msg) + WhatsApp Business API ($$$)
- Zenvia: SMS (R$ 0,10/msg) + WhatsApp caro
- NumVerify: Validação básica (sem WhatsApp)
- AbstractAPI: Validação básica (sem WhatsApp)

**O que NÓS oferecemos:**
- ✅ WhatsApp Validator (100% preciso, custo R$ 0)
- ✅ WhatsApp OTP (98% abertura, custo R$ 0)
- ✅ Validação formato ANATEL (95-98% preciso)
- ✅ DDD → Cidades (100% preciso, BrasilAPI)
- ✅ Tipo mobile/fixo (99%+ preciso)
- ✅ Cache 3 camadas (performance)
- ✅ Planos acessíveis (R$ 29-299/mês vs $100+/mês)

**🔥 Diferencial ÚNICO:**
> "Única API brasileira com WhatsApp Verification real e OTP por WhatsApp sem custo adicional por mensagem!"

#### **❌ O que NÃO vamos implementar (e por quê):**

**1. Operadora Exata:**
- Problema: Portabilidade invalida heurística
- Precisão: ~60% (muito baixa)
- Solução real: API paga (R$ 0,01/req) ou base ANATEL (80%)
- Decisão: **NÃO implementar agora**

**2. HLR Lookup (número ativo?):**
- Problema: Requer acesso a operadoras
- Custo: R$ 0,01-0,05/consulta
- Decisão: **Avaliar demanda futura**

**3. SMS OTP:**
- Problema: Custo alto (R$ 0,10-0,20/msg)
- Concorrência: Twilio/Zenvia já fazem
- Decisão: **Apenas WhatsApp** (diferencial)

#### **Fontes de Dados:**

**✅ CONFIÁVEIS (100%):**
- BrasilAPI (DDD → Cidades)
- Evolution API (WhatsApp verification)
- Regras ANATEL (formato, tipo)

**⚠️ PARCIAIS (80%):**
- Base ANATEL prefixos (sem portabilidade)

**❌ NÃO CONFIÁVEIS (60%):**
- Heurística operadora (último dígito)
- Tabelas desatualizadas

**Decisão:** Usar apenas fontes 100% confiáveis!

#### **🎯 Status:**
- [ ] Planejado
- [ ] Documentado (este arquivo)
- [ ] Aguardando implementação

**Prazo:** 2-3 dias após aprovação  
**Prioridade:** Média-Alta (diferencial único)

---

**Prazo:** 3 meses  
**Prioridade:** Alta (APIs mais demandadas)

---

## 🟣 FASE 3 - EXPANSÃO E COMPLIANCE (3-6 MESES)

**Meta:** 17 APIs | **Status:** Planejado

### **📋 Dados Cadastrais**
- [ ] **CPF**: Validação de dígitos + consulta Receita Federal
- [ ] **Email**: Validação real (não só regex)
- [ ] **Telefone**: Validação + operadora
- [ ] **Operadora**: Identificação + portabilidade

### **🧾 Validação Fiscal e Compliance** 🆕

#### **APIs Públicas (qualquer CNPJ):**
- [ ] **NF-e Validation**: Consulta NF-e por chave de 44 dígitos
  - Fonte: Webservice SEFAZ (gratuito e público)
  - Dados: emitente, destinatário, valor, status
  - Cache: 30 dias (NF-e não muda)
  - Performance: ~500ms
  - **Casos de uso:** Validação de fornecedores, e-commerce, contabilidade

- [ ] **Certidões (CND/CNDT)**: Consulta certidões negativas de débitos
  - CND Federal (Receita Federal)
  - CNDT (Débitos Trabalhistas - TST)
  - Status: Regular/Irregular
  - Fonte: TST + Receita (gratuito via scraping)
  - Cache: 1 dia
  - **Casos de uso:** Due diligence, pré-contratação, licitações

- [ ] **Compras Governamentais**: Licitações e contratos por CNPJ
  - Fonte: Portal da Transparência + ComprasNet (APIs públicas)
  - Dados: licitações vencidas, contratos, valores
  - Cache: 7 dias
  - Custo: Gratuito
  - **Casos de uso:** Inteligência comercial, due diligence

#### **Dados do Próprio Cliente (com autorização):** 🔥 🆕
- [ ] **Meus Documentos Fiscais**: NF-e auto-sync do cliente
  - Cliente envia certificado digital A1
  - Sync automático diário (e-CAC/SEFAZ)
  - NF-e emitidas + recebidas (últimos 12 meses)
  - Download XML/PDF (DANFE)
  - Analytics: volume, valor, top fornecedores/clientes
  - **Diferencial:** Dashboard fiscal unificado
  - **Plano:** Business (R$ 99/mês)

- [ ] **Meus Boletos**: Open Finance integrado 🏦
  - Cliente autoriza via OAuth 2.0
  - Conexão com múltiplos bancos
  - Boletos a pagar + a receber
  - Alertas de vencimento
  - Projeção de cash flow
  - **Diferencial:** Dashboard financeiro unificado
  - **Plano:** Enterprise (R$ 299/mês)
  - **Prazo:** Requer homologação BACEN (3-4 meses)

### **🗺️ Geografia Avançada**
- [ ] **Bairros**: Lista por cidade
- [ ] **Coordenadas**: Lat/Long por CEP
- [ ] **Distância**: Cálculo entre CEPs

### **💰 Financeiro**
- [ ] **SELIC/CDI/IPCA**: Taxas oficiais Banco Central
- [ ] **Indicadores**: Histórico e projeções

### **🚚 Logística**
- [ ] **Frete**: Cálculo Correios + transportadoras
- [ ] **Rastreamento**: Código Correios
- [ ] **Veículos**: Consulta por placa (DENATRAN)

**Prazo:** 3 meses  
**Prioridade:** Alta (NF-e e Certidões) / Média (demais)

---

## 🟡 FASE 4 - DADOS AVANÇADOS E COMPLIANCE (6-9 MESES)

**Meta:** 7 APIs | **Status:** Planejado

### **⚖️ Compliance e Risco** 🆕
- [ ] **Protestos**: Títulos protestados por CNPJ
  - Fonte: Serasa (pago R$ 1,50/req) OU Web scraping cartórios (gratuito)
  - Dados: total protestos, valores, datas, cartórios
  - Cache: 7 dias
  - **Decisão:** Avaliar demanda antes de contratar API paga
  - **Casos de uso:** Análise de crédito, due diligence, risk assessment

- [ ] **Score de Crédito**: Análise de risco empresarial (futuro)
  - Agregação de dados: CNPJ, Certidões, Protestos, Compras Gov
  - Score proprietário (0-1000)
  - Indicadores de risco

### **🏛️ Dados Governamentais**
- [ ] **Judicial**: Processos públicos (PJe + TJs)
- [ ] **CEIS/CNEP**: Empresas inidôneas
- [ ] **Simples Nacional**: Consulta optantes
- [ ] **PEP**: Pessoas Politicamente Expostas

**Fonte:** Portais públicos do governo  
**Método:** Scraping + cache local  
**Prazo:** 3 meses  
**Prioridade:** Baixa (nicho específico)

---

## 🚀 **INFRAESTRUTURA & PERFORMANCE (PLANEJADO)** 🆕

### **Migração Oracle Cloud (São Paulo)** 🎯
**Objetivo:** Reduzir latência de 160ms para 5-15ms

**Status:** Pesquisa concluída ✅ | Automação planejada 📝

#### **Por que Oracle Cloud?**
- ✅ **Região São Paulo disponível** (sa-saopaulo-1)
- ✅ **Always Free Tier** (GRATUITO para sempre)
- ✅ **Latência:** 5-15ms (vs 150-200ms Railway EUA)
- ✅ **Recursos generosos:** 4 vCPUs ARM + 24GB RAM
- ✅ **200GB Storage gratuito**
- ✅ **10TB bandwidth/mês**
- ✅ **Load Balancer incluído**

#### **Automação Via CLI** 🔧
- [ ] **Script de Provisionamento** (`scripts/oracle/01-provision.sh`)
  - [ ] Criar VM (4 cores ARM, 24GB RAM, região SP)
  - [ ] Configurar VCN + Subnet + Security Lists
  - [ ] Instalar Docker + Docker Compose
  - [ ] Setup usuário deploy + SSH
  - [ ] Verificação de custos (R$ 0,00 se free tier)

- [ ] **Script de Deploy** (`scripts/oracle/02-deploy.sh`)
  - [ ] Subir Redis (via Docker Hub)
  - [ ] Subir MongoDB (via Docker Hub)
  - [ ] Subir aplicação backend
  - [ ] Configurar variáveis de ambiente
  - [ ] Setup volumes persistentes (200GB)

- [ ] **CI/CD Automático** (`.github/workflows/deploy-oci.yml`)
  - [ ] Integração com GitHub (branch: main)
  - [ ] Deploy automático em push
  - [ ] Rollback automático em erro
  - [ ] Notificações Slack/Discord

- [ ] **Monitoramento** (`scripts/oracle/monitoring-setup.sh`)
  - [ ] Logs centralizados (OCI Logging)
  - [ ] Métricas (CPU, RAM, Disco)
  - [ ] Alertas (CPU >80%, RAM >90%, Disco >85%)
  - [ ] Dashboard de saúde

- [ ] **Escalabilidade** (`scripts/oracle/scale-up.sh`)
  - [ ] Aumentar vCPU via script
  - [ ] Aumentar RAM via script
  - [ ] Adicionar storage via script
  - [ ] Load balancer setup

#### **Comparação Railway vs Oracle**
| Item | Railway (Atual) | Oracle Free | Oracle Pago |
|------|----------------|-------------|-------------|
| **Região** | EUA (us-west) | BR (São Paulo) | BR (São Paulo) |
| **Latência** | 150-200ms | 5-15ms 🚀 | 5-15ms 🚀 |
| **vCPU** | ~0.5 | 4 cores ARM | 4-64 cores |
| **RAM** | ~512MB | 24GB | 24-512GB |
| **Storage** | ~1GB | 200GB | Ilimitado |
| **Bandwidth** | ~100GB | 10TB | Ilimitado |
| **Custo** | ~$5-10/mês | **R$ 0/mês** ✅ | R$ 130-500/mês |

#### **Resultado Esperado**
```
Performance com cache Redis + Servidor BR:
├─ 1ª request: ~50ms (vs 200ms atual) → 4x mais rápido
├─ 2ª+ request: ~1-5ms (vs 160ms atual) → 32x mais rápido
└─ Competitivo com Brasil API (31ms)
```

**Prazo:** 1-2 semanas  
**Prioridade:** Alta (diferencial competitivo)  
**Documentação:** `/docs/ORACLE_CLOUD_RESEARCH.md`

---

## ⚪ FUTURO (9-12 MESES)

**Meta:** 6 APIs | **Status:** Backlog

### **Recursos Avançados**
- [ ] **Ruas**: Autocomplete de endereços
- [ ] **Demografia**: População, IDH, PIB
- [ ] **NF-e**: Validação de chave
- [ ] **Inscrição Estadual**: Validação por UF
- [ ] **Boletos**: Geração código de barras
- [ ] **Pix**: QR Code estático
- [ ] **Dias Úteis**: Cálculo entre datas
- [ ] **Fusos**: Horários por cidade

**Prazo:** 3 meses  
**Prioridade:** Baixa

---

## 📊 RESUMO POR CATEGORIA

| Categoria | Total | Disponível | Fase 2 | Fase 3 | Fase 4 | Futuro |
|-----------|-------|-----------|--------|--------|--------|--------|
| **📋 Cadastrais** | 7 | 0 | 1 | 4 | 0 | 2 |
| **🗺️ Geografia** | 6 | 1 | 1 | 3 | 0 | 1 |
| **💰 Financeiro** | 5 | 0 | 2 | 1 | 0 | 2 |
| **🚚 Logística** | 4 | 0 | 1 | 3 | 0 | 0 |
| **🔧 Utilidades** | 4 | 0 | 1 | 1 | 0 | 2 |
| **🏛️ Governo** | 5 | 0 | 0 | 0 | 5 | 0 |
| **TOTAL** | **31** | **1** | **6** | **13** | **5** | **6** |

---

## 🎯 PRÓXIMOS 30 DIAS

### **Semana 1-2: CEP + CNPJ**
1. Integrar ViaCEP
2. Implementar fallback Brasil API
3. Scraping Receita Federal (CNPJ)
4. Testes de carga

### **Semana 3: Moedas + Bancos**
1. Integrar API Banco Central
2. Carregar lista de bancos (STR)
3. Endpoints + cache

### **Semana 4: FIPE + Feriados**
1. Integrar FIPE API
2. Gerar calendário de feriados
3. Documentação atualizada
4. Deploy e testes

---

## 📈 CRONOGRAMA VISUAL

```
2025
│
├─ Out/Nov/Dez ████████ Fase 2: Dados Essenciais
│  └─ 6 APIs (CEP, CNPJ, Moedas, Bancos, FIPE, Feriados)
│
├─ Jan/Fev/Mar ████████ Fase 3: Expansão
│  └─ 13 APIs (CPF, Email, Telefone, Bairros, Frete, etc.)
│
├─ Abr/Mai/Jun ████████ Fase 4: Compliance
│  └─ 5 APIs (Judicial, Transparência, CEIS, Simples, PEP)
│
└─ Jul/Ago/Set ████████ Futuro
   └─ 6 APIs (Ruas, Demografia, NF-e, Boletos, Pix, etc.)
```

---

## 💡 OBSERVAÇÕES

### **Dados Gratuitos**
- Priorizar APIs públicas governamentais
- Uso de Brasil API como aggregator
- Banco Central, IBGE, FIPE (todas gratuitas)

### **Scraping**
- Apenas quando não há API oficial
- Respeitar robots.txt
- Cache agressivo (reduzir requests)
- Exemplos: CNPJ (Receita Federal), Judicial (PJe)

### **Cache Strategy**
- **Estático** (Geografia, Bancos): Permanente
- **Semi-estático** (CEP, FIPE): 7 dias
- **Dinâmico** (Moedas): 1 hora
- **Alta rotatividade** (Judicial): 24 horas

### **Escalabilidade**
- MongoDB indexado para performance
- Redis implementado para cache distribuído ✅
- Rate limiting por tenant (já implementado)
- CDN para assets estáticos

---

## 🎮 **SEO & CONVERSÃO (IMPLEMENTADO)** ✅ 🆕

### **Landing Page & Marketing**
- [x] Landing page com hero otimizado
- [x] Cards de APIs em destaque (CEP, CNPJ, GEO)
- [x] Roadmap visual com 36 APIs
- [x] Meta tags avançadas (Open Graph, Twitter Cards)
- [x] Schema.org (Organization, Product, BreadcrumbList)
- [x] Sitemap dinâmico (100+ páginas)
- [x] Robots.txt otimizado
- [x] 14 keywords estratégicas

### **Playground Interativo** ✅
- [x] **Página:** `/playground`
- [x] Teste CEP, CNPJ e GEO **sem cadastro**
- [x] Código copy-paste (JavaScript, Python, PHP, cURL)
- [x] Response time display (~1-200ms)
- [x] Rotas públicas seguras (`/public/*`)
- [x] Seleção automática da primeira API disponível
- [x] Browser fingerprinting para segurança
- [x] Rate limiting por IP
- [x] **Diferencial:** NENHUM concorrente brasileiro tem isso
- [x] **Conversão esperada:** 10-15%

### **Ferramentas Públicas** ✅
1. **CEP Checker** (`/ferramentas/consultar-cep`)
   - [x] Consulta gratuita e sem cadastro
   - [x] Share links funcionais
   - [x] Usa mesma API Key demo do playground
   - [x] Validação de scopes
   - [x] Target: 18.000 buscas/mês

2. **Buscar CEP** (`/ferramentas/buscar-cep`) 🆕
   - [x] Busca reversa: encontra CEP pelo endereço
   - [x] Parâmetros: UF, Cidade, Logradouro
   - [x] Retorna até 50 CEPs por busca
   - [x] Grid responsivo de resultados
   - [x] Badge "NOVO" na página inicial
   - [x] Usa mesma API Key demo do playground
   - [x] Target: 15.000 buscas/mês

3. **CNPJ Validator** (`/ferramentas/validar-cnpj`)
   - [x] Validação em tempo real
   - [x] Dados da Receita Federal
   - [x] Usa mesma API Key demo do playground
   - [x] Validação de scopes
   - [x] Target: 12.000 buscas/mês

### **Landing Pages de APIs** ✅
- [x] **`/apis/cep`** - Hero + Features + Código + Comparação + FAQ
- [x] Tabela comparativa (Retech vs ViaCEP vs Brasil API)
- [x] Casos de uso (E-commerce, Marketplaces, Cadastros)
- [x] FAQ com Accordions (5 perguntas)
- [x] CTAs estratégicos
- [x] Tempos de resposta realistas (~160ms)

### **SEO Técnico** ✅
- [x] Metadata dinâmica por página
- [x] Canonical URLs
- [x] Alt text em imagens
- [x] Semantic HTML
- [x] Acessibilidade (ARIA)
- [x] Performance otimizada (Next.js 15)

---

## 📝 ÚLTIMAS ATUALIZAÇÕES (Out/2025)

### **📅 29 de Outubro de 2025** 🆕

#### **⚖️ API de Artigos Penais (DISPONÍVEL)** ✅
- **Nova API:** Artigos penais brasileiros (Código Penal + Lei de Contravenções)
- **Endpoints:**
  - `GET /penal/artigos` - Lista artigos (autocomplete/select2)
  - `GET /penal/artigos/:codigo` - Busca artigo específico
  - `GET /penal/search?q=texto` - Busca por texto
- **Rotas públicas:** `/public/penal/*` para playground
- **Estrutura hierárquica:** Artigo → Parágrafo → Inciso → Alínea
- **Tipos:** Crime, Contravenção
- **Legislações:** CP (Código Penal), LCP (Lei de Contravenções)
- **Cache 3 camadas:**
  - Redis L1 (~1ms) - 365 dias (dados fixos)
  - MongoDB L2 (~10ms) - Permanente
  - Dados fixos (não muda)
- **Seed inicial:** `seeds/penal.json` com artigos mais comuns
- **Uso principal:** Autocomplete/Select2 para seleção de crimes
- **Scope:** `penal` (controle granular de acesso)
- **Índices MongoDB:** codigo (unique), busca (text), tipo, legislacao
- **Performance:** ~1ms (Redis) / ~10ms (MongoDB)

### **📅 28 de Outubro de 2025** 🆕

#### **🔍 Busca Reversa de CEP (Endereço → CEP)** ✅
- **Novo endpoint:** `GET /cep/buscar?uf=SP&cidade=São+Paulo&logradouro=Paulista`
- **Endpoint público:** `GET /public/cep/buscar` (para ferramentas/playground)
- **Integração:** ViaCEP (busca por endereço)
- **Cache 3 camadas:**
  - Redis L1 (~1ms)
  - MongoDB L2 (~10ms)  
  - ViaCEP L3 (~100ms)
- **Retorno:** Array de até 50 CEPs por busca
- **Validações:**
  - UF: 2 caracteres
  - Cidade e Logradouro: mínimo 3 caracteres
- **Features:**
  - Cache normalizado (search:UF:cidade:logradouro)
  - Promoção automática Redis → MongoDB
  - TTL configurável (mesmo do CEP normal)
  - Graceful degradation
- **Frontend:**
  - Nova ferramenta: `/ferramentas/buscar-cep`
  - Grid responsivo de resultados (até 50 cards)
  - Badge "NOVO" na landing page
  - Botão copiar CEP
  - Integrado com playground
- **Performance:**
  - 1ª busca: ~100ms (ViaCEP)
  - 2ª+ busca: ~1-10ms (cache)
- **Use cases:**
  - Autocomplete de endereços
  - Validação de formulários
  - Preenchimento automático de CEP
  - Busca quando usuário não sabe o CEP

### **📅 27 de Outubro de 2025** 🆕

#### **🔴 Redis Cache (L1) - Sistema Completo** ✅
- Cache em memória para máxima performance (~1ms)
- Conexão via `REDIS_URL` (Railway/Oracle)
- Graceful degradation (se cair, usa MongoDB)
- Keys separadas: `cep:*`, `cnpj:*`, `geo:*`
- TTL: 24h para hot data
- **Admin Dashboard:**
  - Stats: total keys, keys por serviço (CEP, CNPJ, GEO)
  - Memória usada (MB)
  - Status de conexão (🟢/🔴)
  - Botão "Limpar Todo Redis" com confirmação
  - Explicação do fluxo L1→L2→L3

#### **🔧 Cache Independente (CEP + CNPJ)** ✅
- Controles 100% independentes por serviço
- Cada serviço tem seu próprio:
  - Toggle enable/disable
  - TTL (1-365 dias)
  - AutoCleanup (MongoDB TTL Index)
  - Stats em tempo real
  - Botão de limpeza manual
- **Cards renomeados:**
  - "MongoDB Cache - CEP (L2)"
  - "MongoDB Cache - CNPJ (L2)"
- Migração automática de estrutura antiga
- Tudo salvando corretamente ✅

#### **🔒 Segurança API Key Reforçada** ✅
- API Key oculta no frontend (🔒 •••••••)
- Removidos fallbacks inseguros
- `APIKEY_HASH_SECRET` obrigatório
- Panic se variável não configurada
- Secret forte gerado (256 bits base64)
- Botões "Gerar Nova" e "Rotacionar"
- Auto-rotação ao mudar scopes

#### **📊 Analytics com Timezone Brasil** ✅
- Timezone: America/Sao_Paulo (todas as datas)
- Requests "Hoje" vs "Ontem" corretos
- Formatação pt-BR (27 de outubro de 2025)
- Gráficos com dias mais recentes primeiro
- Sem mais dados de datas futuras

#### **🎮 Playground Multi-Camada** ✅
- Rate limiting por IP (configurável)
- Rate limiting global (shared)
- Browser fingerprinting (WebGL, Canvas, Audio)
- Validação de scopes (cep, cnpj, geo)
- Seleção automática da primeira API
- Ferramentas integradas (mesmo API Key demo)

#### **🔍 Pesquisa Oracle Cloud** ✅
- Região São Paulo disponível
- Always Free Tier mapeado
- Automação via OCI CLI planejada
- Scripts de provisionamento desenhados
- Estimativa: R$ 0/mês (free) ou R$ 130-500/mês (expansão)
- Latência esperada: 5-15ms (vs 160ms atual)

#### **🏥 Health Check Melhorado** ✅
- Status REAL de MongoDB e Redis
- Uptime desde startup
- Versão da API
- Auto-refresh 30s no frontend
- Estados visuais (🟢 Operacional, 🟡 Degradado, 🔴 Indisponível)
- Graceful degradation (Redis down não afeta status geral)

#### **🎨 UX Final** ✅
- Performance corrigida em todas as páginas
- Env `NEXT_PUBLIC_DOCS_URL` (links dinâmicos)
- Hero "The Retech Core"
- Rodapé completo (Alan Rezende, Florianópolis, WhatsApp)
- 6 páginas novas:
  - `/precos` - Planos + Status da plataforma
  - `/sobre` - História + Fundador + Missão
  - `/contato` - Formulário → WhatsApp
  - `/status` - Health check real (30s refresh)
  - `/legal/termos` - LGPD compliant
  - `/legal/privacidade` - LGPD compliant

---

### **📅 24 de Outubro de 2025**

#### **✅ API CEP Implementada**
- Endpoint `/cep/:codigo` funcional
- Cache com ViaCEP + Brasil API (fallback)
- Performance: 95% das requests em <10ms (cache)
- Normalização automática de formato

#### **✅ API CNPJ Implementada**
- Endpoint `/cnpj/:numero` funcional
- Brasil API + ReceitaWS (fallback)
- Validação de dígito verificador
- QSA (Quadro de Sócios e Administradores)
- CNAEs completos (principal + secundários)
- Endereço + contatos + capital social
- Cache 30 dias (otimizado para empresas)
- Performance: ~10ms (cache) / ~200ms (Brasil API)

#### **✅ Sistema de Cache Configurável**
- Admin pode ajustar TTL (1-365 dias)
- Toggle enable/disable por serviço
- Limpeza automática (MongoDB TTL Index)
- Limpeza manual com AlertDialog
- Stats em tempo real (total cached, recent 24h)
- Suporte para CEP e CNPJ

#### **✅ Settings Completas**
- Contato/Vendas (WhatsApp dinâmico)
- Cache configurável por API
- Migration automática de schemas antigos
- Bug fix: contact e cache agora salvam corretamente
- Todas as configs persistem entre reloads

#### **✅ Melhorias de UX**
- AlertDialog para confirmações críticas
- Auto-refresh de métricas
- Tratamento de erros aprimorado
- Feedback visual em todas as ações
- Landing page com 3 APIs em destaque

---

## 🎁 DIFERENCIAIS COMPETITIVOS

✅ **Hub Completo** - 36 APIs essenciais, 1 chave única  
✅ **Ultra Performance** - ~1ms (Redis L1) / ~10ms (MongoDB L2) / ~160ms (API L3)  
✅ **Confiável** - 3 camadas de cache + fallback automático + graceful degradation + 99.9% uptime  
✅ **Gratuito** - 1.000 requests/dia sem cartão de crédito  
✅ **Profissional** - Dashboard completo + Redoc + Analytics em tempo real + Timezone BR  
✅ **Transparente** - Dashboard com métricas em tempo real + Activity logs  
✅ **Configurável** - Admin controla cache (independente), rate limits, CORS, JWT, playground  
✅ **Seguro** - API Keys com HMAC-SHA256 + Scopes granulares + Rate limiting multi-camada  
✅ **Playground Público** - Teste sem cadastro + Browser fingerprinting + Ferramentas gratuitas  
✅ **SEO Otimizado** - Meta tags + Sitemap + Schema.org + 14 keywords estratégicas  
✅ **Oracle Cloud Ready** - Migração planejada para latência <15ms (servidor BR)

---

## 🎯 PRÓXIMOS PASSOS

### **Imediato (Esta Semana)** 🔥
1. **Deploy com `APIKEY_HASH_SECRET`** (Segurança)
   - Adicionar variável no Railway
   - Testar playground em produção
   - Confirmar API Keys demo funcionando

2. **Testes de Carga**
   - Validar Redis cache em produção
   - Monitorar hit rate (L1, L2, L3)
   - Otimizar TTLs se necessário

### **Curto Prazo (2-4 Semanas)** 🚀
1. **Migração Oracle Cloud** (Performance)
   - Criar conta Oracle Cloud (região SP)
   - Desenvolver scripts de automação
   - Setup CI/CD com GitHub Actions
   - Testar latência (meta: <15ms)
   - Migração gradual (staging → produção)

2. **Moedas API** (Fase 2)
   - Banco Central API
   - Cotações real-time
   - Histórico 30 dias

3. **Bancos API** (Fase 2)
   - Dados STR Banco Central
   - Cache permanente
   - Lista completa COMPE/ISPB

### **Médio Prazo (1-3 Meses)** 📊
1. **FIPE API** (Fase 2)
   - Tabela FIPE
   - Preços de veículos
   - Cache 7 dias

2. **Feriados API** (Fase 2)
   - Nacionais + Estaduais
   - Ponto facultativo

3. **NF-e Validation** (Fase 3 - Alta demanda)
   - Webservice SEFAZ
   - Validação de chave 44 dígitos

---

## 🆕 ATUALIZAÇÕES RECENTES

### **📅 24 de Outubro de 2025 - Manhã**

### **🏢 API CNPJ Implementada** ✅
- GET /cnpj/:numero completo
- Brasil API + ReceitaWS fallback
- Cache 30 dias configurável
- Validação CNPJ + normalização
- Admin settings completo
- 100% funcional e testado

### **🔒 Sistema de Scopes Completo** ✅
- Scopes granulares: `geo`, `cep`, `cnpj`, `all`
- Proteção em todas as rotas públicas
- Validação automática de permissões
- Frontend com seleção visual (checkboxes)
- Backend retrocompatível (aceita `geo:read` e `geo`)
- Documentação completa em `docs/SCOPES_ANALYSIS.md`

### **⚙️ Admin Settings Aprimorado** ✅
- Cache CEP: TTL configurável (1-365 dias)
- Cache CNPJ: TTL configurável (1-365 dias)
- Stats de cache em tempo real
- Limpeza manual com AlertDialog
- Validação de inputs (onBlur)
- Auto-cleanup via TTL index MongoDB

### **🐛 Correções de Bugs** ✅
- TTL inputs agora aceitam campo vazio durante digitação
- Cache normalizando CEP/CNPJ antes de salvar
- Settings salvando `contact` e `cache` corretamente
- Upsert habilitado em cache (cria se não existir)

### **📚 Documentação** ✅
- `SCOPES_ANALYSIS.md` - Análise completa do sistema
- `SCOPES_SYSTEM.md` - Guia de uso atualizado
- `CHECKLIST_POS_IMPLEMENTACAO.md` - Processo padronizado
- Redoc atualizado com CNPJ

---

### **📅 24 de Outubro de 2025 - Noite/Madrugada** 🌙

### **🚀 Estratégia Completa de SEO Implementada** ✅
- Meta tags avançadas (Open Graph, Twitter Cards, Schema.org)
- Sitemap dinâmico com 100+ páginas
- Robots.txt otimizado
- 14 keywords estratégicas
- Build passando em produção

### **🎮 API Playground Interativo** ✅
- Teste CEP, CNPJ e Geografia **sem cadastro**
- Código copy-paste (JavaScript, Python, PHP, cURL)
- Response time display (~5-200ms)
- Rotas públicas (`/public/*`) implementadas
- **Diferencial:** NENHUM concorrente brasileiro tem isso
- **Conversão esperada:** 10-15%

### **🔧 Ferramentas Públicas (2)** ✅
1. **CEP Checker** (`/ferramentas/consultar-cep`)
   - Target: 18.000 buscas/mês
   - Consulta gratuita e ilimitada
   - Share links funcionais

2. **CNPJ Validator** (`/ferramentas/validar-cnpj`)
   - Target: 12.000 buscas/mês
   - Validação em tempo real
   - Dados da Receita Federal

### **📄 Landing Page API CEP** ✅
- Hero + Features + Código + Comparação
- Tabela comparativa (Retech vs ViaCEP vs Brasil API)
- Casos de uso (E-commerce, Marketplaces, Cadastros, Análise)
- FAQ com Accordions (5 perguntas)
- CTAs estratégicos

### **🆕 Novas APIs Planejadas** ✅
- **NF-e Validation** (Fase 3 - Alta prioridade)
- **Certidões CND/CNDT** (Fase 3 - Alta prioridade)
- **Compras Governamentais** (Fase 3 - Média prioridade)
- **Protestos** (Fase 4 - Avaliar demanda)
- **Score de Crédito** (Fase 4 - Futuro)
- Documentação: `NOVAS_APIS_BOLETOS_NFE.md`

### **📊 Impacto SEO** ✅
- 3.000+ linhas de código
- 18 arquivos criados
- Keywords-alvo: 50k+ buscas/mês
- Expectativa: 5.000+ visitas/mês (mês 3)

---

## 📚 **DOCUMENTAÇÃO TÉCNICA**

### **Arquitetura & Infraestrutura**
- `/docs/ORACLE_CLOUD_RESEARCH.md` - Pesquisa completa Oracle Cloud + Scripts CLI
- `/build/docker-compose.yml` - Compose local (MongoDB + Redis + API)
- `/internal/bootstrap/indexes.go` - Indexes MongoDB (auto-criação no startup)

### **Cache & Performance**
- **Redis L1:** Implementado em `internal/cache/redis_client.go`
- **MongoDB L2:** Implementado nos handlers (CEP, CNPJ)
- **Graceful Degradation:** Redis cai → MongoDB funciona
- **Admin Dashboard:** Cards independentes (Redis, CEP, CNPJ)
- **Migração Automática:** Estrutura antiga → nova (transparente)

### **Segurança**
- **API Keys:** HMAC-SHA256 com `APIKEY_HASH_SECRET` obrigatório
- **Scopes:** Granulares (`cep`, `cnpj`, `geo`, `all`)
- **Rate Limiting:** Multi-camada (Tenant, IP, Global)
- **Browser Fingerprinting:** WebGL, Canvas, Audio (sem travar UI)
- **CORS:** Configurável via admin (strict mode)
- **JWT:** TTL dinâmico (Access: 15min, Refresh: 7 dias)

### **Analytics & Logs**
- **Timezone:** America/Sao_Paulo (todas as datas)
- **Activity Logs:** Auditoria completa (login, API keys, settings)
- **Usage Tracking:** Por dia, tenant, endpoint
- **Metrics Dashboard:** Auto-refresh 30s, gráficos 30 dias

### **Frontend**
- **Admin Portal:** Dashboard, Tenants, API Keys, Settings, Analytics
- **Developer Portal:** Dashboard, Usage, API Keys, Docs
- **Public Pages:** Landing, Playground, Ferramentas, APIs
- **SEO:** Meta tags, Sitemap, Schema.org, Open Graph

---

## 🔧 **VARIÁVEIS DE AMBIENTE NECESSÁRIAS**

### **Backend (Go)**
```bash
# Obrigatórias
MONGO_URI=mongodb://localhost:27017/retech
APIKEY_HASH_SECRET=9gJlYXwSR1kfAv8Dh4mHRb/WGJKpLV5v+NYDsNFWTJ8=  # ✅ NOVO
JWT_ACCESS_SECRET=seu-secret-access
JWT_REFRESH_SECRET=seu-secret-refresh

# Opcionais (com fallback)
PORT=8080
ENV=production  # ou development
REDIS_URL=redis://localhost:6379  # Se não tiver, usa graceful degradation
API_BASE_URL=https://core.theretech.com.br
CORS_ENABLE=true
JWT_ACCESS_TTL=900   # 15 minutos (sobrescrito por admin/settings)
JWT_REFRESH_TTL=604800  # 7 dias (sobrescrito por admin/settings)
```

### **Frontend (Next.js)**
```bash
# Obrigatórias
NEXT_PUBLIC_API_URL=https://core.theretech.com.br

# Opcionais
NODE_ENV=production
```

---

## 💡 **LIÇÕES APRENDIDAS**
### 📚 `#REFERÊNCIA` `#BOAS-PRÁTICAS` `#LIÇÕES`

> **💡 Leia para evitar erros comuns e aplicar melhores práticas!**

### **Performance**
✅ **Redis é essencial:** Reduz latência de 160ms → 1ms (160x)  
✅ **MongoDB Indexes:** Auto-criação no startup evita erros  
✅ **Graceful Degradation:** Sistema funciona mesmo se Redis cair  
✅ **Timezone Matters:** Usar America/Sao_Paulo evita bugs de data

### **Segurança**
✅ **Nunca usar fallbacks inseguros:** Panic se variável crítica não existir  
✅ **Ocultar secrets no frontend:** Usuário não precisa ver API Key completa  
✅ **Scopes granulares:** Melhor que `all` ou nada  
✅ **Rate limiting por camada:** IP + Tenant + Global

### **UX & DX**
✅ **Feedback visual:** Usuário precisa saber que ação foi bem-sucedida  
✅ **Confirmações críticas:** AlertDialog antes de deletar  
✅ **Logs detalhados:** Console logs ajudam muito no debug  
✅ **Playground público:** Maior diferencial competitivo (conversão)

### **DevOps**
✅ **Docker Compose:** Facilita dev e deploy  
✅ **Variáveis de ambiente:** Separar por ambiente (.env.local, .env.production)  
✅ **CI/CD:** GitHub Actions + SSH deploy  
✅ **Migrations:** Auto-aplicar indexes/schemas no startup

---

## 🎯 **METAS 2025-2026**

### **Q4 2025 (Out-Dez)**
- [x] Fase 1 completa (3 APIs)
- [ ] Migração Oracle Cloud (latência <15ms)
- [ ] Fase 2 completa (6 APIs)
- [ ] 10.000 requests/mês

### **Q1 2026 (Jan-Mar)**
- [ ] Fase 3 completa (17 APIs)
- [ ] NF-e + Certidões implementadas
- [ ] 50.000 requests/mês
- [ ] 100 tenants ativos

### **Q2 2026 (Abr-Jun)**
- [ ] Fase 4 completa (7 APIs)
- [ ] Open Finance integrado (Boletos)
- [ ] 200.000 requests/mês
- [ ] 500 tenants ativos

### **Q3-Q4 2026 (Jul-Dez)**
- [ ] 36 APIs completas
- [ ] 1M requests/mês
- [ ] 2.000 tenants ativos
- [ ] Break-even financeiro

---

## 📊 **KPIs DE SUCESSO**

### **Performance**
- ✅ Latência P50: <160ms (atual)
- 🎯 Latência P50: <15ms (Oracle Cloud)
- 🎯 Cache Hit Rate: >80%
- ✅ Uptime: 99.9%

### **Adoção**
- ✅ APIs disponíveis: 3/36 (8%)
- 🎯 APIs disponíveis: 36/36 (100%)
- ✅ Playground público: Implementado
- 🎯 Conversão playground: 10-15%

### **Financeiro**
- ✅ Custo infraestrutura: $5-10/mês (Railway)
- 🎯 Custo infraestrutura: R$ 0/mês (Oracle Free Tier)
- 🎯 MRR (Monthly Recurring Revenue): R$ 5.000/mês (Q2 2026)

### **Qualidade**
- ✅ Cobertura de testes: 0% (TODO)
- 🎯 Cobertura de testes: 80%
- ✅ Documentação: OpenAPI 3.0
- ✅ Admin Dashboard: Completo

---

## 🔗 **LINKS IMPORTANTES**

### **Produção**
- **Frontend:** https://core.theretech.com.br
- **Backend API:** https://core.theretech.com.br (via Railway)
- **Docs:** https://core.theretech.com.br/docs
- **Playground:** https://core.theretech.com.br/playground

### **Repositórios**
- **Backend:** github.com/theretech/retech-core
- **Frontend:** github.com/theretech/retech-core-admin

### **Infraestrutura**
- **Railway:** dashboard.railway.app
- **Cloudflare:** dash.cloudflare.com (DNS)
- **Oracle Cloud:** cloud.oracle.com (futuro)

### **Monitoramento**
- **Railway Logs:** Via dashboard web
- **MongoDB:** Compass local (mongodb://localhost:27017)
- **Redis:** RedisInsight ou CLI

---

## 🤝 **CONTRIBUINDO**

### **Processo de Desenvolvimento**
1. Criar branch: `feature/nome-da-api`
2. Implementar handler: `internal/http/handlers/nome.go`
3. Adicionar rota: `internal/http/router.go`
4. Criar testes (TODO)
5. Atualizar documentação
6. PR para `main`

### **Checklist Nova API**
### 🏗️ `#NOVA-API` `#BACKEND` `#CHECKLIST-BACKEND`

> **⚠️ CHECKLIST PARA BACKEND! Para frontend completo, veja linha 1975.**

- [ ] Handler com cache 3 camadas (Redis + MongoDB + Externa)
- [ ] Validação de input
- [ ] Normalização de dados
- [ ] **Scope específico (`nome`)** 🆕
  - [ ] Backend: `scope_middleware.go` → `validScopes`
  - [ ] Frontend - API Keys: `apikey-drawer.tsx` → `availableScopes`
  - [ ] Frontend - Settings: `admin/settings/page.tsx` → `allowedApis`
- [ ] Rate limiting
- [ ] Admin settings (toggle, TTL, stats, clear)
- [ ] **Seed/Migration** (se dados fixos) 🆕
  - [ ] Criar migration em `bootstrap/migrations.go`
  - [ ] ⚠️ Seed roda automaticamente no startup (não precisa fazer nada manual)
- [ ] Testes unitários
- [ ] Documentação OpenAPI
- [ ] **Landing page completa da API** (`/apis/nome`)
  - [ ] `app/apis/[nome]/page.tsx` - Hero + Features + Code + FAQ
  - [ ] `app/apis/[nome]/layout.tsx` - Metadata SEO completa
  - [ ] Breadcrumbs Schema para SEO
  - [ ] Adicionar no sitemap
- [ ] Card na landing principal (`app/page.tsx`)
- [ ] Adicionar no playground (se aplicável)

### **🔧 PADRÕES ARQUITETURAIS (Implementado 29/10/2025)** 🆕
### 📖 `#LEITURA-OBRIGATÓRIA` `#PADRÕES` `#NOVA-API`

> **⚠️ LEIA ANTES de implementar handlers com APIs externas!**

#### **PADRÃO 1: APIs Externas**
❌ **Nunca** hardcode URLs de APIs externas  
✅ **Sempre** usar ENV vars obrigatórias

**Implementar:**
1. Criar `Get{NOME}PrimaryURL()` em `config/apis.go` com panic se ENV vazia
2. Adicionar validação em `config.ValidateExternalAPIsConfig()`
3. Documentar ENVs em `env.example`
4. Usar `config.Get{NOME}PrimaryURL()` no handler, nunca string literal
5. Configurar ENVs no Railway antes do deploy

**Sequência de fallback:** Cache Redis → Cache MongoDB → Primary URL (ENV) → Fallback URL (ENV) → Erro 404

#### **PADRÃO 2: Cache Multi-Camada**
Toda API externa deve ter cache duplo: Redis (L1 - 24h) + MongoDB (L2 - 7-30 dias)

**Implementar:**
1. Buscar Redis primeiro (return se hit)
2. Buscar MongoDB segundo (promover para Redis se hit)
3. Buscar API externa (salvar em ambos se sucesso)
4. Retornar erro 404 apenas se todas falharem

#### **PADRÃO 3: Scopes Granulares**
Cada API deve ter scope específico (`cep`, `cnpj`, `geo`, etc), nunca genérico.

**Implementar:**
1. Adicionar scope em `auth/scope_middleware.go` no map `validScopes`
2. Aplicar middleware `RequireScope()` nas rotas
3. Documentar scope necessário no OpenAPI

#### **PADRÃO 4: Configurações Dinâmicas**
Settings editáveis via Admin devem estar em MongoDB (`system_settings`), não em ENV.

**Usar ENV para:** URLs de APIs, secrets, infra  
**Usar MongoDB para:** Rate limits, cache TTL, features toggles

#### **PADRÃO 5: Fail-Fast**
Validar configurações críticas no startup (`main.go`), não na primeira request.

**Implementar:** Chamar `config.Validate{Feature}Config()` antes de `config.Load()`

---

## 📞 **CONTATO & SUPORTE**

**WhatsApp:** +55 48 99961-6679  
**Email:** suporte@theretech.com.br  
**Documentação:** https://core.theretech.com.br/docs

---

**🚀 ROADMAP EM CONSTANTE EVOLUÇÃO!**

**Última atualização:** 31 de Outubro de 2025 🆕  
**Revisão Maior:** Consolidado "Checklist Completo: Nova API" com todas as lições de Artigos Penais  
**Próxima revisão:** 15 de Novembro de 2025 (após migração Oracle Cloud)

**Juntos, construindo o futuro das APIs brasileiras! 🇧🇷**

---

## 📋 **CHECKLIST: NOVA FUNCIONALIDADE/ENDPOINT** 
### 🔧 `#FUNCIONALIDADE` `#ENDPOINT` `#CHECKLIST-SIMPLIFICADO`

> **⚠️ Para NOVA API COMPLETA, use o "CHECKLIST COMPLETO: NOVA API - IMPLEMENTAÇÃO FRONTEND" acima!**

**Após implementar uma nova funcionalidade ou endpoint em API existente:**

> 💡 **Baseado na implementação da "Busca Reversa de CEP"**  
> 💡 **Este checklist é para funcionalidades, não APIs novas!**  
> ⚠️ **USE ESTE CHECKLIST ao adicionar endpoint em API existente!**

### **📝 O Que Fazer Após Implementar:**

1. **Atualizar Redoc (OpenAPI)**
   - Arquivo: `internal/docs/openapi.yaml`
   - Adicionar endpoint com descrição, parâmetros, responses e exemplos
   - **⚠️ Documentar tratamento de dados:**
     - Acentos: aceita ou precisa remover?
     - Case: maiúscula, minúscula ou tanto faz?
     - Encoding: automático ou dev precisa fazer?
     - Formato: com/sem traço, pontos, etc
     - Adicionar exemplos com múltiplos formatos

2. **Atualizar Documentação do Painel**
   - Arquivo: `internal/http/handlers/tenant.go` (função `GetMyConfig`)
   - Adicionar endpoint na lista da categoria correspondente
   - Incluir emoji 🆕 se for funcionalidade recente
   - Descrição clara e objetiva (uma linha)

3. **Verificar Analytics/Logging (Automático)**
   - ✅ Middleware `UsageLogger` deve estar aplicado na rota (verificar router.go)
   - ✅ Logs salvam automaticamente em `api_usage_logs`
   - ✅ Analytics agrupa por `apiName` (extraído do primeiro segmento da URL)
   - ✅ Dashboard `/admin/analytics` mostra automaticamente
   - ⚠️ **NADA precisa fazer** se middleware está aplicado!

4. **Playground - Avaliar se Faz Sentido Adicionar**
   
   **⚠️ NEM TUDO vai para o playground!**
   
   **✅ ADICIONAR no playground SE:**
   - Funcionalidade CORE da API
   - Input simples (1-2 campos max)
   - Desenvolvedores vão querer **testar o código**
   - Gerar código automático é útil
   - Exemplo: consulta CEP por código, busca CNPJ, listar UFs
   
   **❌ NÃO ADICIONAR no playground SE:**
   - Já existe ferramenta dedicada funcional
   - Input muito complexo (3+ campos)
   - Foco é usuário final, não desenvolvedor
   - Já tem SEO próprio (ferramenta pública)
   - Exemplo: busca reversa CEP (3 inputs + ferramenta própria)
   
   **🎯 Regra de ouro:**
   > "Playground é para devs testarem e copiarem código. Ferramenta é para usuários resolverem problemas."
   
   **📋 Exemplos de Decisões:**
   
   | Funcionalidade | Playground? | Ferramenta? | Motivo |
   |----------------|-------------|-------------|--------|
   | Consulta CEP por código | ✅ Sim | ✅ Sim | Core + simples (1 input) |
   | Busca reversa CEP | ❌ Não | ✅ Sim | 3 inputs + foco SEO |
   | Consulta CNPJ | ✅ Sim | ✅ Sim | Core + simples (1 input) |
   | Lista UFs | ✅ Sim | ❌ Não | Sem input + útil para devs |
   | Cotação moedas | ✅ Sim | ❌ Não | Simples + devs precisam testar |
   | Cálculo de frete | ❌ Não | ✅ Sim | 5+ inputs + foco usuário final |
   
   **🔄 Fluxo de Decisão:**
   ```
   Nova funcionalidade implementada
            ↓
   Quantos inputs? → 1-2 → Público-alvo? → Devs → ✅ PLAYGROUND + Ferramenta
            ↓                            → Usuários → ✅ Ferramenta
            ↓
   Quantos inputs? → 3+ → ✅ Apenas FERRAMENTA (não playground)
   ```

5. **Atualizar ROADMAP**
   - Marcar endpoints como [x] concluído
   - Adicionar na seção "Últimas Atualizações" com data
   - ⚠️ Verificar se altera contador (Nova API vs Funcionalidade)

6. **Verificar Console (Frontend)** 🆕
   - ⚠️ **Verificar erros de React Keys duplicadas**
   - **Problema comum:** Duas chaves com mesmo `keyId` ou `key` não única
   - **Solução:** Sempre usar `key` única em `.map()`:
     - Preferir `item.id` se existir
     - Se não, usar combinação: `${item.id || item.keyId}-${index}`
     - Para arrays aninhados: `${parentId}-child-${index}-${childValue}`
   - **Locais críticos:**
     - Tabelas: `TableRow key={item.id || \`${item.keyId}-${index}\`}`
     - Lists: `div key={item.id || index}`
     - Badges dentro de array: `key={\`${parentId}-badge-${index}\`}`
   - **Exemplo corrigido:**
     ```typescript
     // ❌ ERRADO: Pode ter duplicatas
     {items.map((item) => <Row key={item.keyId} />)}
     
     // ✅ CORRETO: Sempre único
     {items.map((item, index) => <Row key={item.id || `${item.keyId}-${index}`} />)}
     
     // ✅ CORRETO: Arrays aninhados
     {items.map((item, i) => (
       <div key={item.id || i}>
         {item.scopes.map((scope, j) => (
           <Badge key={`${item.id}-scope-${j}-${scope}`} />
         ))}
       </div>
     ))}
     ```
   - **Verificar no console do navegador:**
     - Abrir DevTools (F12)
     - Verificar Console tab
     - Procurar erros: "Encountered two children with the same key"
     - Corrigir todos antes de commit

7. **Testar Tudo**
   - Backend: endpoint funcionando, cache L1/L2/L3, validações
   - Frontend: ferramenta pública, playground (se foi adicionado)
   - Docs: Redoc e Painel Docs mostrando endpoint
   - Analytics: fazer 2-3 requests e verificar em `/admin/analytics`
   - Mobile: responsividade
   - **Console:** Verificar que não há erros de keys duplicadas 🆕

8. **Verificar Segurança**
   - API Key obrigatória
   - Scope correto aplicado
   - Rate limiting funcionando
   - Logs de usage salvando

9. **Performance**
    - Cache hit após 2ª request
    - Response time adequado
    - Graceful degradation (se Redis cair)

10. **Melhorias no Código (se aplicável)**
    - URL Encoding: usar `url.PathEscape()` ou `url.QueryEscape()` para parâmetros
    - Validação: normalizar entrada antes de processar
    - Tratamento: aceitar diferentes formatos (com/sem acentos, formatação, etc)

11. **Verificar SEO (Pós-Deploy)**
    - **Títulos únicos:** Cada página deve ter title diferente
      - Criar `layout.tsx` em cada pasta se necessário
      - Formato: `[Função] - [Seção] | Retech Core`
      - Ex: `Login - Portal do Desenvolvedor | Retech Core API`
    - **Meta descriptions únicas:** Cada página deve ter description específica
      - Descrever o propósito exato da página
      - Incluir keywords relevantes
    - **Verificar 404s:**
      - Testar todos os links internos
      - Criar redirects se necessário (`/termos` → `/legal/termos`)
    - **Robots.txt:** Verificar se permite crawling
    - **Sitemap:** Verificar se todas as páginas públicas estão incluídas
    - **Ferramenta:** Usar Google Search Console ou Ahrefs Site Audit

12. **Commit e Deploy**
   - Build sem erros (Go + Next.js)
   - Commit com mensagem clara
   - Deploy (Railway auto-deploy)
   - Smoke test em produção

---

### **📦 Arquivos Principais (Funcionalidades):**

**Backend:**
- `internal/http/handlers/[nome].go` - Handler (novo endpoint)
- `internal/http/handlers/tenant.go` - GetMyConfig (docs painel)
- `internal/http/router.go` - Rotas (public + protected)

**Frontend (se aplicável):**
- `app/ferramentas/[nome]/page.tsx` - Ferramenta pública (se criou)
- `app/playground/page.tsx` - Playground (se adicionou)

**Documentação:**
- `internal/docs/openapi.yaml` - Redoc (novo endpoint)
- `docs/Planning/ROADMAP.md` - Última atualização

---

### **📝 Exemplo Real - Busca Reversa de CEP (Funcionalidade):**

**Backend (3 arquivos):**
- Handler (+255 linhas)
- GetMyConfig (+6 linhas)
- Router (+12 linhas)

**Frontend (1 arquivo novo):**
- Ferramenta pública (+250 linhas)

**Documentação (2 arquivos):**
- OpenAPI (+220 linhas)
- ROADMAP (atualização)

**Total:** 6 arquivos | ~750 linhas | ~4 horas ⏱️

> **💡 Para uma NOVA API completa, veja o exemplo da "API de Artigos Penais" no checklist acima!**

---

### **⚠️ IMPORTANTE - Nova API vs Funcionalidade:**

**Determinar se altera contador da landing page (9/36 APIs):**

✅ **NOVA API** (atualizar contador):
- Serviço completamente novo
- Fonte de dados distinta
- Scope próprio

❌ **FUNCIONALIDADE** (não altera contador):
- Novo endpoint em API existente
- Mesma fonte de dados
- Mesmo scope

**Exemplo:** Busca reversa CEP = Funcionalidade (não altera 9/36)

---

### **🔗 Arquivos de Referência:**
- Handler: `internal/http/handlers/cep.go` (linha 274)
- Router: `internal/http/router.go` (linhas 156-162, 229)
- Tenant: `internal/http/handlers/tenant.go` (linha 371-376)
- OpenAPI: `internal/docs/openapi.yaml` (linhas 196-377)
- Ferramenta: `app/ferramentas/buscar-cep/page.tsx`
- Landing: `app/page.tsx` (linhas 399-425)
- UsageLogger: `internal/middleware/usage_logger.go` (extrai apiName automaticamente)

---

### **📊 Como Verificar se Analytics Está Funcionando:**

1. **Fazer algumas requests** para o novo endpoint:
```bash
   curl "http://localhost:8080/cep/buscar?uf=SP&cidade=Sao+Paulo&logradouro=Paulista" \
     -H "X-API-Key: sua_api_key"
   ```

2. **Acessar dashboard admin:**
   ```
   http://localhost:3002/admin/analytics
   ```

3. **Verificar:**
   - ✅ Total de requests aumentou
   - ✅ API "CEP" aparece com mais requests
   - ✅ Endpoint `/cep/buscar` aparece no "Top Endpoints"
   - ✅ Response time está sendo medido

4. **⚠️ Nota importante:**
   - `/cep/buscar` e `/cep/:codigo` são contados juntos como API "cep"
   - Mas aparecem separados em "Top Endpoints"
   - Isso é o comportamento esperado!

5. **O que você verá no analytics:**
   ```
   📊 Breakdown por API:
   - CEP: 150 requests (inclui /cep/:codigo + /cep/buscar)
   - CNPJ: 80 requests
   - Geografia: 45 requests
   
   📈 Top Endpoints:
   - /cep/:codigo - 95 requests
   - /cep/buscar - 55 requests  ← NOVO!
   - /cnpj/:numero - 80 requests
```

---

### **📝 Boas Práticas de Documentação:**
### 📚 `#REFERÊNCIA` `#DOCUMENTAÇÃO` `#OPENAPI`

> **💡 Use estas práticas ao documentar endpoints no OpenAPI!**

**Sempre documente para o desenvolvedor:**

1. **Tratamento de Acentos:**
   ```yaml
   description: |
     **✅ Aceita acentos:** "São Paulo", "João Pessoa"
     - Com acentos: cidade=São Paulo (recomendado)
     - Sem acentos: cidade=Sao Paulo (funciona, menos preciso)
   ```

2. **Case Sensitivity:**
   ```yaml
   description: |
     **Case:** Maiúsculas/minúsculas não importam
     - ✅ "são paulo" = "São Paulo" = "SÃO PAULO"
     - ⚠️ UF deve ser MAIÚSCULO: "SP" (não "sp")
   ```

3. **Formato de Entrada:**
   ```yaml
   description: |
     **Formato aceito:**
     - Com formatação: 00.000.000/0001-91
     - Sem formatação: 00000000000191
     - Ambos funcionam! A API normaliza automaticamente.
   ```

4. **Encoding:**
   ```yaml
   description: |
     **Encoding:** Automático pelo backend
     - Espaços: use + ou %20
     - Acentos: enviados diretamente
     - Caracteres especiais: URL encoded automaticamente
   ```

5. **Exemplos Práticos:**
   - Sempre incluir 2-3 exemplos com diferentes formatos
   - Mostrar caso típico + caso com acentos + caso URL encoded
   - Indicar qual é recomendado (⭐)

**Exemplo Completo (Busca Reversa CEP):**
- ✅ 3 exemplos de cURL (com acentos, sem acentos, encoded)
- ✅ Dicas de formatação (acentos, espaços, case)
- ✅ Indicação de recomendado
- ✅ Avisos sobre precisão

---

## 📋 **CHECKLIST COMPLETO: NOVA API - IMPLEMENTAÇÃO FRONTEND** 🆕
### 🏗️ `#NOVA-API` `#CHECKLIST-COMPLETO` `#FRONTEND`

> **💡 Consolidado das lições aprendidas durante implementação de Artigos Penais**  
> **Data:** 31 de Outubro de 2025  
> **Objetivo:** Garantir que NADA seja esquecido ao criar uma nova API
> 
> ⚠️ **USE ESTE CHECKLIST ao implementar uma NOVA API!**

### **🎯 Como Usar Este Checklist:**

1. ✅ **Após implementar backend** (handlers, rotas, seeds, cache)
2. ✅ **Execute TODOS os itens** desta lista
3. ✅ **Marque cada item** como concluído
4. ✅ **Teste visual completo** (Seção 10)
5. ✅ **Verifique console browser** (zero erros)

---

### **1️⃣ LANDING PAGE DA API** (`/apis/[nome]`)

#### **1.1. Criar Arquivos Obrigatórios:**
- [ ] **`app/apis/[nome]/page.tsx`** - Landing page completa
- [ ] **`app/apis/[nome]/layout.tsx`** - Metadata SEO completa

#### **1.2. Estrutura da Landing Page:**
- [ ] **Hero Section:**
  - [ ] Título chamativo com nome da API
  - [ ] Descrição clara e objetiva (1-2 frases)
  - [ ] Badge de status ("Disponível", "NOVO", etc)
  - [ ] 2 CTAs: "Consultar Grátis" + "Criar Conta (1.000 req/dia)"
  - [ ] **NÃO adicionar "Testar no Playground"** (evitar redundância)

- [ ] **Features (3-4 cards):**
  - [ ] Ícones representativos
  - [ ] Títulos curtos
  - [ ] Descrições objetivas
  - [ ] Destaque dos diferenciais

- [ ] **Exemplos de Código:**
  - [ ] Tabs: JavaScript, Python, PHP
  - [ ] **URLs usando `apiBaseURL`** (NUNCA hardcoded!)
  - [ ] Exemplos práticos e copy-paste ready
  - [ ] Headers com API Key demo

- [ ] **Casos de Uso (4 cards):**
  - [ ] Título do caso de uso
  - [ ] Descrição de aplicação real
  - [ ] Público-alvo claro

- [ ] **FAQ (Accordion com 4-6 perguntas):**
  - [ ] Linguagem **comercial**, não técnica
  - [ ] Evitar termos como "seed", "migration", "TTL"
  - [ ] Focar em benefícios para o cliente
  - [ ] Exemplo: "Os dados são atualizados?" → Falar de equipe, não de tecnologia

- [ ] **Seção "Ferramenta Grátis"** (se aplicável):
  - [ ] Card grande linkando para `/ferramentas/[nome]`
  - [ ] Badge "✨ Teste Grátis"
  - [ ] CTA: "Consultar Agora →"
  - [ ] Texto: "⚡ Cache permanente • ✅ 100% Gratuito • 🎁 Sem cadastro"

- [ ] **CTA Final:**
  - [ ] Card com gradiente
  - [ ] "Pronto para começar?"
  - [ ] Botão: "Começar Grátis Agora"

#### **1.3. Metadata SEO (layout.tsx):**
- [ ] `title`: "API de [Nome] - [Descrição] | Retech Core"
- [ ] `description`: Com keywords + benefícios (160 caracteres)
- [ ] `keywords`: Array com 10-15 termos relevantes
- [ ] `openGraph`: title, description, type, url, images
- [ ] `twitter`: card, title, description, images
- [ ] Canonical URL correto

#### **1.4. SEO Schemas:**
- [ ] `<BreadcrumbSchema>` com navegação
- [ ] `<FAQSchema>` se houver FAQ

#### **1.5. Verificações de Conteúdo:**
- [ ] **URLs Dinâmicas:** TODAS usando `const apiBaseURL = process.env.NEXT_PUBLIC_API_URL`
- [ ] **Tempos Realistas:** Evitar `~0.5ms`, `~1ms` (considerar latência de rede)
  - [ ] Usar: `< 5ms` (Redis), `< 20ms` (MongoDB), `~160ms` (API externa)
- [ ] **Transparência de Limites:**
  - [ ] NUNCA "sem limites" ou "ilimitado" para ferramentas gratuitas
  - [ ] Usar: "Gratuito", "Sem cadastro"
  - [ ] Limites reais: "100 req/dia", "10 req/min"
- [ ] **Terminologia Moderna:**
  - [ ] "Autocomplete" ou "componentes de busca" (NÃO "Select2")
- [ ] **Linguagem Cliente-Friendly:**
  - [ ] FAQs sem jargão técnico
  - [ ] Explicações comerciais, não de implementação

#### **1.6. Referências:**
- Ver: `app/apis/cep/page.tsx`
- Ver: `app/apis/penal/page.tsx`

---

### **2️⃣ FERRAMENTA GRÁTIS** (`/ferramentas/[nome]`) - SE APLICÁVEL

- [ ] **Criar:** `app/ferramentas/[nome]/page.tsx`
- [ ] **Criar:** `app/ferramentas/[nome]/layout.tsx` (SEO)
- [ ] **Input de consulta** + botão
- [ ] **Exibir resultado completo** com formatação
- [ ] **Suporte a URL parameter** (`?parametro=valor`) para pré-preencher
- [ ] **SEO Schemas:** FAQSchema, BreadcrumbSchema
- [ ] **Usar API Key demo** (`demoApiKey` do playground)
- [ ] **Textos honestos:**
  - [ ] "Gratuito" ✅
  - [ ] "Sem cadastro" ✅
  - [ ] "Sem limites" ❌ (NUNCA!)
- [ ] **Exibir tempo de resposta real** (medido no client-side com `performance.now()`)
- [ ] **Mensagens de erro** claras
- [ ] **Ícones e badges** para visual

---

### **3️⃣ HOME PAGE** (`/`)

#### **3.1. Seção "APIs Disponíveis":**
- [ ] **Adicionar novo card:**
  - [ ] Emoji representativo
  - [ ] Título da API
  - [ ] Descrição curta (1-2 linhas)
  - [ ] Badge "Disponível" ou "NOVO"
  - [ ] Link para `/apis/[nome]`
- [ ] **Atualizar contagem total:**
  - [ ] Ex: "36 APIs" → "37 APIs"
  - [ ] Atualizar em TODOS os lugares (hero, stats, etc)

#### **3.2. Seção "Comece a Usar Imediatamente":**
- [ ] **Avaliar onde adicionar:**
  - [ ] Card para `/apis/[nome]` (documentação) OU
  - [ ] Card para `/ferramentas/[nome]` (ferramenta grátis)
- [ ] **Layout responsivo:**
  - [ ] Grid 2 colunas se houver 2+ cards relacionados
  - [ ] Card limpo e elegante
- [ ] **Link correto:**
  - [ ] Ferramenta: `<Link href="/ferramentas/[nome]">`
  - [ ] Documentação: `<Link href="/apis/[nome]">`

#### **3.3. Seção "Planos e Preços":**
- [ ] **USAR COMPONENTE:** `<PricingPlans variant="landing" />`
- [ ] **NÃO duplicar código** dos planos inline!
- [ ] **Verificar se componente está atualizado:**
  - [ ] Plano Free: "4 APIs (CEP, CNPJ, Geografia, Penal)"
  - [ ] Plano Pro: "Todas as APIs (sem premium)"
  - [ ] Plano Business/Enterprise: "Todas as APIs + Premium"

#### **3.4. Seção "Integração em Minutos":**
- [ ] **Exemplos de código usando `apiBaseURL`:**
  - [ ] JavaScript: `${apiBaseURL}/[endpoint]`
  - [ ] Python: `${apiBaseURL}/[endpoint]`
  - [ ] PHP: `${apiBaseURL}/[endpoint]`
- [ ] **Botão "Ver Documentação Completa":**
  - [ ] `href={`${apiBaseURL}/docs`}`

#### **3.5. Footer:**
- [ ] **Link "Documentação":**
  - [ ] `href={`${apiBaseURL}/docs`}`
  - [ ] NUNCA hardcoded!

---

### **4️⃣ PÁGINAS ADMINISTRATIVAS**

#### **4.1. Página de Preços** (`/precos`):
- [ ] **USAR COMPONENTE:** `<PricingPlans variant="page" />`
- [ ] **NÃO duplicar código** dos planos inline!
- [ ] **Seção "Status da Plataforma":**
  - [ ] Adicionar nova API à lista (ex: "✅ API Penal: Operacional")

#### **4.2. Página de Status** (`/status`):
- [ ] **Atualizar contagem:**
  - [ ] Ex: "3/36" → "4/37"
- [ ] **Seção "APIs Principais":**
  - [ ] Atualizar texto: "CEP, CNPJ, Geografia" → "CEP, CNPJ, Geografia, Penal"

#### **4.3. Página de Registro** (`/painel/register`):
- [ ] **Seção "Plano Free incluso":**
  - [ ] "1.000 requisições/dia" ✅
  - [ ] "4 APIs (CEP, CNPJ, Geografia, Penal)" ✅
  - [ ] "Cache 3 camadas + Dashboard" ✅
  - [ ] "Documentação completa" ✅
  - [ ] "Sem cartão de crédito" ✅

#### **4.4. Admin Settings** (`/admin/settings`):
- [ ] **Playground - allowedApis:**
  - [ ] Adicionar scope: `allowedApis: ['cep', 'cnpj', 'geo', 'penal']`
- [ ] **Playground - Lista de APIs:**
  - [ ] Adicionar no `.map()`: `{['cep', 'cnpj', 'geo', 'penal'].map(...)}`
- [ ] **Cache Config:**
  - [ ] Adicionar `penal?: { enabled, ttlDays, autoCleanup }` (se aplicável)
  - [ ] Normalização de dados ao carregar
  - [ ] Valores padrão no estado inicial
- [ ] **Redis Stats:**
  - [ ] Adicionar contagem: `penalKeys` em `loadRedisStats()`
  - [ ] Exibir: "└─ ⚖️ PENAL: {penalKeys} keys"
- [ ] **Cache Stats Endpoint:**
  - [ ] Se tiver cache, carregar stats: `loadPenalCacheStats()`
- [ ] **Limpeza Redis:**
  - [ ] Incluir nova API na confirmação de limpeza

#### **4.5. API Keys Drawer** (`components/apikeys/apikey-drawer.tsx`):
- [ ] **Adicionar scope:**
  - [ ] Array `availableScopes`:
    ```typescript
    { value: 'penal', label: '⚖️ PENAL - Artigos penais (Código Penal)' }
    ```

---

### **5️⃣ PÁGINAS LEGAIS**

#### **5.1. Termos de Uso** (`/legal/termos`):
- [ ] **Seção "Uso Permitido":**
  - [ ] Adicionar API: "Consultas a dados públicos (CEP, CNPJ, Geografia, Artigos Penais)"
- [ ] **Seção "Fontes de Dados":**
  - [ ] Adicionar fonte (se aplicável): "Legislação Brasileira"
- [ ] **Seção "Dados Atualizados":**
  - [ ] Explicar política de cache/TTL
  - [ ] Ex: "Artigos Penais são permanentes (alterados apenas quando legislação muda)"
- [ ] **Seção "Dados Públicos":**
  - [ ] Adicionar à lista: "Artigos Penais"

#### **5.2. Política de Privacidade** (`/legal/privacidade`):
- [ ] **Seção "Dados NÃO Coletados":**
  - [ ] "CEPs, CNPJs, ou Artigos Penais que você consultou"
- [ ] **Seção "Como Funciona o Cache":**
  - [ ] Explicar TTL da nova API
  - [ ] Ex: "7 dias para CEP, 30 dias para CNPJ, permanente para Artigos Penais"

---

### **6️⃣ PLAYGROUND**

- [ ] **Avaliar se faz sentido adicionar:**
  - [ ] ✅ Adicionar SE: funcionalidade CORE, input simples (1-2 campos), foco em devs
  - [ ] ❌ NÃO adicionar SE: já tem ferramenta dedicada, input complexo (3+ campos), foco em usuários finais
- [ ] **Se adicionar:**
  - [ ] Adicionar API à lista com endpoints de exemplo
  - [ ] Garantir que `?api=penal` pré-seleciona a API
  - [ ] **Playground deve defaultar para CEP** se nenhum parâmetro fornecido

---

### **7️⃣ DOCUMENTAÇÃO (REDOC)**

#### **7.1. OpenAPI YAML** (`internal/docs/openapi.yaml`):
- [ ] **Criar nova tag:**
  - [ ] Nome da API com emoji
- [ ] **Documentar TODOS os endpoints:**
  - [ ] `GET`, `POST`, etc
  - [ ] Parâmetros com tipos e validações
  - [ ] Responses (200, 400, 404, 500)
  - [ ] Exemplos de request/response
- [ ] **Schemas (request/response):**
  - [ ] Definir na seção `components/schemas`
  - [ ] Referenciar: `$ref: '#/components/schemas/ArtigoPenal'`
- [ ] **Seção "Uso Típico (Autocomplete)":**
  - [ ] **SEMPRE usar `__API_BASE_URL__`** (nunca hardcode!)
  - [ ] Exemplo de fetch/axios
- [ ] **⚠️ CUIDADOS:**
  - [ ] **Evitar acentos** em `description` (podem quebrar parser)
  - [ ] **Enum multi-line:**
    ```yaml
    enum:
      - crime
      - contravencao
    ```
  - [ ] Não usar enum inline: `enum: [crime, contravencao]` ❌

---

### **8️⃣ SEO E SITEMAP**

#### **8.1. Sitemap** (`app/sitemap.ts`):
- [ ] **Adicionar URL `/apis/[nome]`:**
  - [ ] `url: 'https://core.theretech.com.br/apis/penal'`
  - [ ] `changeFrequency: 'monthly'`
  - [ ] `priority: 0.8`
- [ ] **Adicionar ferramenta** (se criou):
  - [ ] `url: 'https://core.theretech.com.br/ferramentas/penal'`
  - [ ] `changeFrequency: 'weekly'`
  - [ ] `priority: 0.7`

---

### **9️⃣ BACKEND COMPLEMENTOS**

#### **9.1. Cache Stats Endpoint:**
- [ ] **Criar:** `GET /admin/cache/[nome]/stats`
- [ ] **Handler:** `[nome]Handler.GetCacheStats()`
- [ ] **Retornar:**
  - [ ] `totalCached` (total de registros)
  - [ ] `recentCached` (últimas 24h)
  - [ ] `cacheEnabled` (bool)
  - [ ] `cacheTTLDays` (número)
  - [ ] `autoCleanup` (bool)

#### **9.2. Redis Monitoring:**
- [ ] **Arquivo:** `internal/http/handlers/redis_stats.go`
- [ ] **GetStats:**
  - [ ] Adicionar: `penalKeys, _ := redisClient.Keys(ctx, "penal:*")`
  - [ ] Incluir em `totalKeys`
  - [ ] Retornar: `"penalKeys": len(penalKeys)`

---

### **🔟 CHECKLIST PÓS-IMPLEMENTAÇÃO**

#### **10.1. Teste Visual Completo:**
- [ ] `http://localhost:3000/` (home - TODAS as seções!)
- [ ] `http://localhost:3000/apis/[nome]` (landing page)
- [ ] `http://localhost:3000/ferramentas/[nome]` (ferramenta grátis)
- [ ] `http://localhost:3000/playground?api=[nome]` (se adicionou)
- [ ] `http://localhost:3000/precos` (planos atualizados)
- [ ] `http://localhost:3000/status` (contagem atualizada)
- [ ] `http://localhost:3000/painel/register` (plano free)
- [ ] `http://localhost:3000/admin/settings` (scope + cache)
- [ ] `http://localhost:3000/admin/apikeys` (scope disponível)
- [ ] `http://localhost:3000/legal/termos` (API mencionada)
- [ ] `http://localhost:3000/legal/privacidade` (API mencionada)
- [ ] `http://localhost:8080/docs` (Redoc atualizado)

#### **10.2. Teste Console Browser:**
- [ ] **Abrir DevTools (F12)**
- [ ] **Verificar Console tab**
- [ ] **Zero erros de React keys:**
  - [ ] "Encountered two children with the same key" ❌
  - [ ] Se houver, corrigir com `key` única:
    ```typescript
    // ❌ ERRADO
    {items.map(item => <Row key={item.keyId} />)}
    
    // ✅ CORRETO
    {items.map((item, i) => <Row key={item.id || `${item.keyId}-${i}`} />)}
    ```
- [ ] **Zero warnings desnecessários**
- [ ] **Teste em múltiplas páginas** (home, admin, painel)

#### **10.3. Teste de URLs Dinâmicas:**
- [ ] **Todas as URLs refletem `NEXT_PUBLIC_API_URL`:**
  - [ ] Exemplos de código na landing page
  - [ ] Botões "Ver Documentação"
  - [ ] Links no footer
  - [ ] Ferramenta grátis (fetch calls)
- [ ] **Testar mudando env:** `NEXT_PUBLIC_API_URL=http://localhost:8080`
- [ ] **Verificar que NENHUMA URL está hardcoded**

#### **10.4. Teste Funcional:**
- [ ] **Backend:** Endpoint funcionando, cache L1/L2/L3
- [ ] **Frontend:** Ferramenta pública funcional
- [ ] **Playground:** Seleção e teste da API (se adicionou)
- [ ] **Docs:** Redoc mostrando endpoint
- [ ] **Analytics:** Fazer 2-3 requests e verificar em `/admin/analytics`
- [ ] **Mobile:** Responsividade

#### **10.5. Verificar Segurança:**
- [ ] API Key obrigatória
- [ ] Scope correto aplicado
- [ ] Rate limiting funcionando
- [ ] Logs de usage salvando

#### **10.6. Performance:**
- [ ] Cache hit após 2ª request
- [ ] Response time adequado
- [ ] Graceful degradation (se Redis cair)

---

### **📦 ARQUIVOS PRINCIPAIS A MODIFICAR**

#### **Frontend:**
- `app/apis/[nome]/page.tsx` ← NOVA API (landing)
- `app/apis/[nome]/layout.tsx` ← NOVA API (SEO)
- `app/ferramentas/[nome]/page.tsx` ← Ferramenta grátis
- `app/page.tsx` ← Home (5+ seções!)
- `app/precos/page.tsx` ← Preços (usar componente)
- `app/status/page.tsx` ← Status (contagem)
- `app/painel/register/page.tsx` ← Registro (plano free)
- `app/admin/settings/page.tsx` ← Admin (scopes, cache, Redis)
- `app/admin/apikeys/page.tsx` ← Verificar keys (se houver erros)
- `components/apikeys/apikey-drawer.tsx` ← Scopes
- `components/pricing/PricingPlans.tsx` ← Planos (verificar se atualizado)
- `app/legal/termos/page.tsx` ← Termos
- `app/legal/privacidade/page.tsx` ← Privacidade
- `app/sitemap.ts` ← Sitemap

#### **Backend:**
- `internal/http/handlers/[nome].go` ← Handler principal
- `internal/http/handlers/redis_stats.go` ← Redis monitoring
- `internal/http/router.go` ← Rotas (public + protected + admin)
- `internal/auth/scope_middleware.go` ← Scopes
- `internal/domain/settings.go` ← Cache config (se aplicável)

#### **Documentação:**
- `internal/docs/openapi.yaml` ← Redoc
- `docs/Planning/ROADMAP.md` ← Este arquivo

---

### **⚠️ ARMADILHAS COMUNS**

#### **❌ O QUE EVITAR:**

1. **URLs Hardcoded:**
   - ❌ `https://api-core.theretech.com.br`
   - ✅ `process.env.NEXT_PUBLIC_API_URL`

2. **Tempos Irrealistas:**
   - ❌ "~0.5ms", "~1ms" (ignora latência de rede)
   - ✅ "< 5ms" (Redis), "< 20ms" (MongoDB)

3. **Promessas Falsas:**
   - ❌ "Sem limites", "Ilimitado" (para ferramentas gratuitas)
   - ✅ "Gratuito", "Sem cadastro"

4. **Terminologia Desatualizada:**
   - ❌ "Select2" (jQuery plugin antigo)
   - ✅ "Autocomplete", "Componentes de busca"

5. **Jargão Técnico em FAQs:**
   - ❌ "Seed", "Migration", "TTL", "Cache L1/L2"
   - ✅ "Nossa equipe atualiza", "Dados ficam disponíveis"

6. **Duplicação de Código:**
   - ❌ Copiar/colar cards de planos em múltiplas páginas
   - ✅ Usar `<PricingPlans />` component

7. **React Keys Duplicadas:**
   - ❌ `key={item.keyId}` (pode ter duplicatas)
   - ✅ `key={item.id || `${item.keyId}-${index}`}`

8. **Acentos em OpenAPI:**
   - ❌ `description: "Código do artigo"`
   - ✅ `description: "Codigo do artigo"` (sem acento)

9. **Esquecer de Atualizar:**
   - ❌ Atualizar apenas home, esquecer preços/status/termos
   - ✅ Seguir checklist completo (todas as 10 seções!)

---

### **✅ CONCLUSÃO**

Este checklist consolida **TUDO** que aprendemos durante a implementação da API de Artigos Penais. 

**Tempo estimado:** 4-6 horas para frontend completo de uma nova API.

**Resultado esperado:** Zero esquecimentos, zero erros, 100% consistente!

---

**Última atualização:** 31 de Outubro de 2025  
**Revisado e aprovado:** ✅

---


## ⚠️ **CRITÉRIO DE CONTAGEM: NOVA API vs FUNCIONALIDADE**
### 🎯 `#IMPORTANTE` `#NOVA-API` `#FUNCIONALIDADE` `#DECISÃO`

> **⚠️ LEIA ANTES de atualizar contador na landing page!**

**Use este guia para decidir se atualiza o contador da landing page:**

### **✅ CONTA como "NOVA API" (atualizar 9/36 → 10/36):**
1. **Serviço completamente novo** com fonte de dados distinta
2. **Scope próprio** (novo escopo de permissão)
3. **Domínio diferente** (ex: após CEP/CNPJ, adicionar Moedas)
4. **Collection MongoDB separada** para cache principal
5. **Documentação independente** no Redoc

**Exemplos:**
- ✅ CPF (após ter CEP/CNPJ)
- ✅ Moedas (após ter CEP/CNPJ/Geografia)
- ✅ FIPE (após ter Moedas)
- ✅ Feriados (após ter FIPE)

### **❌ NÃO CONTA como nova API (manter 9/36):**
1. **Novo endpoint** na mesma API
2. **Variação de busca** (ex: busca reversa)
3. **Filtro adicional** em API existente
4. **Formato alternativo** de resposta
5. **Mesmo domínio** e scope

**Exemplos:**
- ❌ Busca reversa CEP (já temos CEP)
- ❌ CNPJ por nome fantasia (já temos CNPJ)
- ❌ Geografia com filtro adicional (já temos Geografia)
- ❌ CEP com coordenadas (já temos CEP)

### **📊 Impacto na Landing Page:**

**Se for NOVA API:**
```
Antes: 25% (9/36 APIs)
Depois: 27% (10/36 APIs)
```

**Se for FUNCIONALIDADE:**
```
Antes: 25% (9/36 APIs)
Depois: 25% (9/36 APIs) ← NÃO MUDA!
```

**O que atualizar quando for FUNCIONALIDADE:**
- ✅ Seção da API no ROADMAP (adicionar novo endpoint)
- ✅ Card da API na landing (adicionar recurso)
- ✅ Documentação Redoc (novo path)
- ✅ Última atualização no ROADMAP
- ❌ Contador de APIs (mantém igual!)
- ❌ Barra de progresso (mantém igual!)

### **🎯 Regra de Ouro:**

> **"Se usa o mesmo scope e mesma fonte de dados, é FUNCIONALIDADE, não API nova!"**

**Em caso de dúvida:**
- Pergunte: "Um desenvolvedor precisaria de 2 API Keys diferentes?"
- Se NÃO → É funcionalidade
- Se SIM → É API nova

---



**✅ CHECKLIST SIMPLIFICADO PRONTO PARA USO!**

---

**Última atualização:** 31 de Outubro de 2025 🆕  
**Revisão Maior:** Consolidado "Checklist Completo: Nova API"  
**Próxima revisão:** 15 de Novembro de 2025

**Juntos, construindo o futuro das APIs brasileiras! 🇧🇷**
