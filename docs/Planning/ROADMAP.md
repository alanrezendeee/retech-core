# 🚀 ROADMAP RETECH CORE API

**Atualizado:** 24 de Outubro de 2025  
**Status:** Fase 1 Concluída ✅ | Fase 2 Em Andamento 🔵 (2/6 APIs - 33%)

---

## 📊 PROGRESSO GERAL

```
████████░░░░░░░░░░░░░ 29% (9/31 APIs)

Fase 1: ████████████ 100% ✅
Fase 2: ████░░░░░░░░  33% 🔵 (2/6)
Fase 3: ░░░░░░░░░░░░   0% ⚪
Fase 4: ░░░░░░░░░░░░   0% ⚪
```

**APIs Totais:** 31  
**Disponíveis:** 3 (Geografia + CEP + CNPJ) 🚀  
**Em Desenvolvimento:** 4 (Fase 2)  
**Planejadas:** 24

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
- [x] Indexação MongoDB (performance <100ms)

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

## 🟣 FASE 3 - EXPANSÃO (3-6 MESES)

**Meta:** 13 APIs | **Status:** Planejado

### **📋 Dados Cadastrais**
- [ ] **CPF**: Validação de dígitos + consulta Receita Federal
- [ ] **Email**: Validação real (não só regex)
- [ ] **Telefone**: Validação + operadora
- [ ] **Operadora**: Identificação + portabilidade

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
**Prioridade:** Média

---

## 🟡 FASE 4 - COMPLIANCE (6-9 MESES)

**Meta:** 5 APIs | **Status:** Planejado

### **🏛️ Dados Governamentais**
- [ ] **Judicial**: Processos públicos (PJe + TJs)
- [ ] **Transparência**: Licitações e convênios
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

✅ **Tudo em uma API** - Uma chave, 31+ endpoints  
✅ **Performance** - <100ms de resposta  
✅ **Confiável** - Cache inteligente + fallbacks  
✅ **Gratuito** - 1.000 requests/dia sem cartão  
✅ **Documentação** - Redoc + exemplos funcionais  
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

## 🆕 ATUALIZAÇÕES RECENTES (24/10/2025)

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

**🚀 Próxima sessão: Moedas API! Rumo às 31 APIs! Vamos nessa!**
