# 📚 FONTES DE DADOS - PESQUISA TÉCNICA

**Atualizado:** 23 de Outubro de 2025

---

## 🎯 CRITÉRIOS DE SELEÇÃO

1. **Gratuito** - Prioridade máxima
2. **API Oficial** - Melhor opção
3. **Scraping** - Última alternativa (com cache agressivo)
4. **Dados Públicos** - Preferir fontes governamentais
5. **Confiabilidade** - Uptime > 99%

---

## 📮 CEP

### **Fonte Principal: ViaCEP**
- **URL**: `https://viacep.com.br/ws/{cep}/json/`
- **Gratuito**: ✅ Sim
- **Limite**: Ilimitado
- **Dados**: Logradouro, bairro, cidade, UF, complemento
- **Coordenadas**: ❌ Não
- **Confiabilidade**: ⭐⭐⭐⭐⭐

### **Fallback: Brasil API**
- **URL**: `https://brasilapi.com.br/api/cep/v1/{cep}`
- **Gratuito**: ✅ Sim
- **Limite**: Ilimitado
- **Agregador**: Usa ViaCEP + outros
- **Confiabilidade**: ⭐⭐⭐⭐

### **Coordenadas: OpenStreetMap Nominatim**
- **URL**: `https://nominatim.openstreetmap.org/search?postalcode={cep}&country=BR&format=json`
- **Gratuito**: ✅ Sim
- **Limite**: 1 req/segundo
- **Cache**: Obrigatório (TOS)
- **Dados**: Lat/Long

### **Estratégia**
1. Buscar no cache local (7 dias)
2. ViaCEP (primário)
3. Brasil API (fallback)
4. Enriquecer com coordenadas (Nominatim)
5. Salvar no MongoDB

---

## 🏢 CNPJ

### **Fonte Principal: Receita Federal (Dados Públicos)**
- **URL**: `https://servicos.receita.fazenda.gov.br/Servicos/cnpjreva/Cnpjreva_Solicitacao.asp`
- **Gratuito**: ✅ Sim (público)
- **API Oficial**: ❌ Não (apenas consulta web)
- **Scraping**: ✅ Necessário
- **Dados**: Razão social, situação, sócios, atividades, endereço
- **Confiabilidade**: ⭐⭐⭐⭐⭐ (fonte oficial)

### **Fallback: Brasil API (CNPJ)**
- **URL**: `https://brasilapi.com.br/api/cnpj/v1/{cnpj}`
- **Gratuito**: ✅ Sim
- **Limite**: Rate limit não especificado
- **Dados**: Completos

### **Alternativa: ReceitaWS**
- **URL**: `https://www.receitaws.com.br/v1/cnpj/{cnpj}`
- **Gratuito**: ✅ 3 req/minuto
- **Plano Pago**: $29/mês (ilimitado)
- **Dados**: Completos

### **Estratégia**
1. Cache local (30 dias - dados mudam pouco)
2. Brasil API (primário)
3. ReceitaWS (fallback)
4. Scraping Receita Federal (último recurso)
5. Atualização sob demanda

---

## 💵 MOEDAS (COTAÇÕES)

### **Fonte: Banco Central do Brasil**
- **URL**: `https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata/CotacaoDolarPeriodo(dataInicial=@dataInicial,dataFinalCotacao=@dataFinalCotacao)`
- **Gratuito**: ✅ Sim
- **API Oficial**: ✅ Sim (OData)
- **Moedas**: USD, EUR, GBP, JPY, CHF, CAD, AUD
- **Histórico**: Sim (desde 1984)
- **Confiabilidade**: ⭐⭐⭐⭐⭐

### **Cripto: CoinGecko**
- **URL**: `https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=brl`
- **Gratuito**: ✅ 50 calls/minuto
- **Moedas**: BTC, ETH, etc.

### **Estratégia**
1. Cache (1 hora)
2. API Banco Central (moedas fiat)
3. CoinGecko (criptomoedas)
4. Endpoint único `/moedas/cotacao?moedas=USD,EUR,BTC`

---

## 🏦 BANCOS BRASILEIROS

### **Fonte: Banco Central (STR - Sistema de Transferência de Reservas)**
- **URL**: `https://www.bcb.gov.br/pom/spb/estatistica/port/ASTR003.pdf` (lista manual)
- **API**: ❌ Não oficial
- **Dados Públicos**: ✅ Sim
- **Método**: Download + parse inicial (uma vez)
- **Dados**: COMPE, ISPB, nome, tipo

### **Brasil API**
- **URL**: `https://brasilapi.com.br/api/banks/v1`
- **Gratuito**: ✅ Sim
- **Dados**: Lista completa

### **Estratégia**
1. Carregar lista Brasil API (uma vez)
2. Salvar no MongoDB (permanente)
3. Endpoint: `GET /bancos` (cache local)
4. Atualização: Manual (trimestral)

---

## 🚗 FIPE (TABELA DE PREÇOS DE VEÍCULOS)

### **Fonte: FIPE API (via Denatran)**
- **URL**: `https://veiculos.fipe.org.br/api/veiculos/ConsultarValorComTodosParametros`
- **Gratuito**: ✅ Sim (não oficial mas público)
- **Dados**: Marcas, modelos, anos, preços
- **Tipos**: Carros, motos, caminhões

### **Brasil API**
- **URL**: `https://brasilapi.com.br/api/fipe/preco/v1/{fipeCode}`
- **Gratuito**: ✅ Sim

### **Estratégia**
1. Cache (7 dias - preços mudam mensalmente)
2. Brasil API (primário)
3. FIPE direto (fallback)
4. Endpoints:
   - `GET /fipe/marcas?tipo=carros`
   - `GET /fipe/modelos?marca={codigo}`
   - `GET /fipe/preco?codigo={fipe_code}`

---

## 📅 FERIADOS

### **Fonte: Arquivo Local (gerado programaticamente)**
- **Nacionais**: Lei 662/1949 + leis específicas
- **Estaduais**: Leis estaduais
- **Municipais**: ❌ Não (muito variável)
- **Método**: Cálculo + arquivo JSON

### **API Pública: Brasil API**
- **URL**: `https://brasilapi.com.br/api/feriados/v1/{ano}`
- **Gratuito**: ✅ Sim

### **Estratégia**
1. Gerar calendário para cada ano (script)
2. Salvar no MongoDB
3. Endpoint: `GET /feriados/:ano`
4. Tipos: nacional, estadual, ponto facultativo
5. Campos: data, nome, tipo, uf (se estadual)

---

## 👤 CPF

### **Validação de Dígitos**
- **Método**: Algoritmo local (não precisa API)
- **Cálculo**: Módulo 11

### **Consulta Situação (Receita Federal)**
- **URL**: `https://servicos.receita.fazenda.gov.br/Servicos/CPF/ConsultaSituacao/ConsultaPublica.asp`
- **Gratuito**: ✅ Sim (público)
- **API**: ❌ Não (scraping necessário)
- **Dados**: Situação cadastral (regular/pendente/suspenso)
- **Cache**: 30 dias

### **Estratégia**
1. Validação local (sempre)
2. Consulta Receita (se solicitado)
3. Cache agressivo (dados mudam pouco)

---

## ✉️ EMAIL (VALIDAÇÃO)

### **Validação Sintaxe**
- **Método**: Regex local
- **Gratuito**: ✅ Sim

### **Validação Real: Hunter.io**
- **URL**: `https://api.hunter.io/v2/email-verifier`
- **Gratuito**: ✅ 50 verificações/mês
- **Plano Pago**: $49/mês (1.000/mês)
- **Dados**: Existe, catch-all, disposable

### **Alternativa: ZeroBounce**
- **Gratuito**: ✅ 100 verificações/mês
- **Confiabilidade**: ⭐⭐⭐⭐

### **Estratégia**
1. Validação sintaxe (sempre)
2. Hunter.io (quota gratuita)
3. Cache (7 dias)
4. Endpoint premium (opcional)

---

## 📱 TELEFONE

### **Validação: Biblioteca libphonenumber**
- **Método**: Local (Google library)
- **Gratuito**: ✅ Sim (open source)
- **Dados**: Formato válido, tipo (fixo/móvel)

### **Operadora: Brasil API**
- **URL**: `https://brasilapi.com.br/api/ddd/v1/{ddd}`
- **Gratuito**: ✅ Sim
- **Dados**: Estado, cidades

### **Portabilidade: Não há API pública gratuita**
- **Soluções pagas**: TotalVoice, Nexmo
- **Estratégia**: Apenas validação + DDD

---

## 🏘️ BAIRROS

### **Fonte: OpenStreetMap + IBGE**
- **Método**: Scraping + download datasets
- **Gratuito**: ✅ Sim (dados abertos)
- **Cobertura**: Parcial (cidades grandes)
- **Cache**: Permanente (atualização semestral)

---

## 📦 FRETE (CORREIOS)

### **Fonte: Correios**
- **URL**: `http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx`
- **Gratuito**: ❌ Requer contrato
- **Alternativa**: Melhor Envio API
- **Gratuito Melhor Envio**: ✅ Sim (integração)

---

## 📍 RASTREAMENTO (CORREIOS)

### **Fonte: Correios SRO**
- **URL**: `https://proxyapp.correios.com.br/v1/sro-rastro/{codigo}`
- **Gratuito**: ❌ Requer token
- **Scraping**: Possível (com cautela)
- **Cache**: 1 hora

---

## 🚙 VEÍCULOS (PLACA)

### **Fonte: DENATRAN**
- **API Oficial**: ❌ Não pública
- **Soluções pagas**: Olho no Carro, Consulta Placa
- **Estratégia**: Futuro (quando viável)

---

## ⚖️ JUDICIAL

### **Fonte: Tribunais (PJe + e-SAJ)**
- **APIs**: ❌ Não públicas
- **Scraping**: Complexo (cada TJ diferente)
- **Estratégia**: Fase 4 (estudo de viabilidade)

---

## 📊 TRANSPARÊNCIA

### **Fonte: Portal da Transparência**
- **URL**: `http://www.portaltransparencia.gov.br/api-de-dados`
- **API**: ✅ Sim (pública)
- **Gratuito**: ✅ Sim
- **Dados**: Licitações, convênios, gastos

---

## 🚫 CEIS/CNEP

### **Fonte: Portal da Transparência**
- **URL**: `http://www.portaltransparencia.gov.br/sancoes/ceis`
- **API**: ✅ Sim
- **Gratuito**: ✅ Sim
- **Dados**: Empresas/pessoas punidas

---

## 📈 SELIC/CDI/IPCA

### **Fonte: Banco Central**
- **URL**: `https://api.bcb.gov.br/dados/serie/bcdata.sgs.{codigo}/dados`
- **Códigos**:
  - SELIC: 432
  - CDI: 12
  - IPCA: 433
- **Gratuito**: ✅ Sim
- **API Oficial**: ✅ Sim

---

## 🎯 SUMÁRIO DE VIABILIDADE

| API | Fonte | Tipo | Custo | Viabilidade | Prioridade |
|-----|-------|------|-------|-------------|-----------|
| CEP | ViaCEP | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Alta |
| CNPJ | Brasil API | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Alta |
| Moedas | Banco Central | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Alta |
| Bancos | Brasil API | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Alta |
| FIPE | Brasil API | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Alta |
| Feriados | Local | Arquivo | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Alta |
| CPF | Receita | Scraping | 🟢 Grátis | ⭐⭐⭐⭐ | Média |
| Email | Hunter.io | API | 🟡 Quota | ⭐⭐⭐ | Média |
| Telefone | libphonenumber | Local | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Média |
| Bairros | OSM | Scraping | 🟢 Grátis | ⭐⭐⭐ | Baixa |
| Frete | Melhor Envio | API | 🟢 Grátis | ⭐⭐⭐⭐ | Média |
| Rastreamento | Correios | Scraping | 🟡 Complexo | ⭐⭐ | Baixa |
| Veículos | - | - | 🔴 Pago | ⭐ | Baixa |
| Judicial | TJs | Scraping | 🔴 Complexo | ⭐ | Baixa |
| Transparência | Portal | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Média |
| CEIS | Portal | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Média |
| SELIC/CDI | Banco Central | API | 🟢 Grátis | ⭐⭐⭐⭐⭐ | Média |

---

**✅ Conclusão**: 22/31 APIs são viáveis com fontes gratuitas! 🎉

