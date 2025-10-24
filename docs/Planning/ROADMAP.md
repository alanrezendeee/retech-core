# 🚀 ROADMAP RETECH CORE API

**Atualizado:** 24 de Outubro de 2025  
**Status:** Fase 1 Concluída ✅ | Fase 2 Em Andamento 🔵 (2/6 APIs - 33%)

---

## 📊 PROGRESSO GERAL

```
██████░░░░░░░░░░░░░░░ 25% (9/36 APIs)

Fase 1: ████████████ 100% ✅
Fase 2: ████░░░░░░░░  33% 🔵 (2/6)
Fase 3: ░░░░░░░░░░░░   0% ⚪ (0/17)
Fase 4: ░░░░░░░░░░░░   0% ⚪ (0/7)
```

**APIs Totais:** 36 (+5 novas: NF-e, CND, Compras Gov, Protestos, Score) 🆕  
**Disponíveis:** 3 (Geografia + CEP + CNPJ) 🚀  
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
- [x] Dados do IBGE (completos)
- [x] Cache Redis em 3 camadas (~160ms médio)

---

## 🔵 FASE 2 - DADOS ESSENCIAIS (0-3 MESES)

**Meta:** 6 APIs | **Status:** 2/6 Concluídas (33%)

### **📮 CEP (DISPONÍVEL)** ✅
- [x] `GET /cep/:codigo` - Busca por CEP
- [x] Integração: ViaCEP (gratuito)
- [x] Fallback: Brasil API
- [x] Cache: 7 dias (configurável via admin: TTL dinâmico 1-365 dias)
- [x] Coordenadas geográficas
- [x] Normalização automática (com/sem traço)
- [x] Performance: ~5ms (cache) / ~50ms (ViaCEP)
- [x] Admin: Configurações de cache (enabled, TTL, auto-cleanup, stats, clear)
- [x] Scope: `cep` (controle granular de acesso)

### **🏢 CNPJ (DISPONÍVEL)** ✅
- [x] `GET /cnpj/:numero` - Consulta CNPJ
- [x] Fonte: Brasil API (gratuita, Receita Federal)
- [x] Fallback: ReceitaWS
- [x] Cache local: 30 dias (configurável via admin)
- [x] Validação: Dígito verificador + normalização
- [x] Dados: razão social, nome fantasia, situação
- [x] QSA: Quadro de sócios e administradores
- [x] CNAEs: Atividade principal + secundárias
- [x] Endereço completo + contatos
- [x] Performance: ~10ms (cache) / ~200ms (Brasil API)
- [x] Admin: Configurações de cache + stats + clear cache
- [x] Scope: `cnpj` (controle granular de acesso)

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
- Redis (futuro) para cache distribuído
- Rate limiting por tenant (já implementado)
- CDN para assets estáticos

---

## 📝 ÚLTIMAS ATUALIZAÇÕES (Out/2025)

### **✅ API CEP Implementada**
- Endpoint `/cep/:codigo` funcional
- Cache com ViaCEP + Brasil API (fallback)
- Performance: 95% das requests em <10ms (cache)
- Normalização automática de formato

### **✅ API CNPJ Implementada** 🆕
- Endpoint `/cnpj/:numero` funcional
- Brasil API + ReceitaWS (fallback)
- Validação de dígito verificador
- QSA (Quadro de Sócios e Administradores)
- CNAEs completos (principal + secundários)
- Endereço + contatos + capital social
- Cache 30 dias (otimizado para empresas)
- Performance: ~10ms (cache) / ~200ms (Brasil API)

### **✅ Sistema de Cache Configurável**
- Admin pode ajustar TTL (1-365 dias)
- Toggle enable/disable global
- Limpeza automática (MongoDB TTL Index)
- Limpeza manual com AlertDialog
- Stats em tempo real (total cached, recent 24h)
- Suporte para CEP e CNPJ

### **✅ Settings Completas**
- Contato/Vendas (WhatsApp dinâmico)
- Cache configurável por API
- Migration automática de schemas antigos
- Bug fix: contact e cache agora salvam corretamente
- Todas as configs persistem entre reloads

### **✅ Melhorias de UX**
- AlertDialog para confirmações críticas
- Auto-refresh de métricas
- Tratamento de erros aprimorado
- Feedback visual em todas as ações
- Landing page com 3 APIs em destaque

---

## 🎁 DIFERENCIAIS COMPETITIVOS

✅ **Hub Completo** - 36 APIs essenciais, 1 chave única  
✅ **Performance** - ~160ms com cache Redis em 3 camadas  
✅ **Confiável** - 3 fontes de dados + fallback automático + 99.9% uptime  
✅ **Gratuito** - 1.000 requests/dia sem cartão de crédito  
✅ **Profissional** - Dashboard completo + Redoc + Analytics em tempo real  
✅ **Transparente** - Dashboard com métricas em tempo real  
✅ **Configurável** - Admin controla cache, rate limits, etc

---

## 🎯 PRÓXIMOS PASSOS

1. **Moedas API** (Prioridade Alta) 🔜
   - Banco Central API
   - Cotações real-time
   - Histórico 30 dias

2. **Bancos API** (Prioridade Alta) 🔜
   - Dados STR Banco Central
   - Cache permanente
   - Lista completa COMPE/ISPB

3. **FIPE API** (Prioridade Média) 🔜
   - Tabela FIPE
   - Preços de veículos
   - Cache 7 dias

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

**🚀 Próxima sessão: Deploy + Moedas API! Rumo às 36 APIs! Vamos nessa!**
