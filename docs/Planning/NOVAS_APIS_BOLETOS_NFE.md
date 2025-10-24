# 🆕 NOVAS APIs: BOLETOS E NOTAS FISCAIS POR CNPJ

**Data:** 24 de Outubro de 2025  
**Status:** 📋 Planejamento e Pesquisa  
**Prioridade:** 🟡 Média-Alta (Fase 3-4)

---

## 🎯 FUNCIONALIDADES PROPOSTAS

### **1. Busca de Boletos por CNPJ** 📄
- Consultar boletos emitidos ou recebidos por uma empresa
- Histórico de títulos protestados
- Certidões negativas de débitos
- Status de inadimplência

### **2. Busca de Notas Fiscais por CNPJ** 🧾
- NF-e emitidas pela empresa
- NF-e recebidas pela empresa
- Dados fiscais e tributários
- Volume de vendas/compras

---

## 🔍 ANÁLISE DE FONTES DE DADOS

### **API 1: CONSULTA DE BOLETOS POR CNPJ**

#### **Possíveis Fontes:**

##### **1.1 Cartórios de Protesto (IEPTB)**
**Descrição:**
- Instituto de Estudos de Protesto de Títulos do Brasil
- Registro público de títulos protestados
- Dados de inadimplência e protestos

**Acesso:**
- ❌ **API pública:** NÃO existe API oficial gratuita
- 🟡 **API paga:** Alguns cartórios oferecem consultas pagas
- 🟡 **Web scraping:** Possível, mas complexo (cada estado tem sistema diferente)
- ✅ **Dados de 3ª:** Serasa, Boa Vista SCPC (pagos)

**Viabilidade:** 🟡 Média
- Requer parcerias ou web scraping
- Custo por consulta
- Complexidade técnica alta

**Estimativa de Implementação:**
- Tempo: 4-6 semanas
- Custo: R$ 0,50 - R$ 2,00 por consulta (APIs de 3ª)
- Complexidade: Alta

##### **1.2 Certidões Negativas de Débitos (CND)**
**Descrição:**
- Certidão de débitos trabalhistas (TST)
- Certidão de débitos federais (Receita Federal)
- Certidão de débitos estaduais (SEFAZ)
- Certidão de débitos municipais (Prefeituras)

**Acesso:**
- ✅ **Receita Federal:** Portal e-CAC (requer certificado digital)
- ✅ **TST:** Consulta pública com captcha
- 🟡 **SEFAZ:** Cada estado tem sistema próprio
- ❌ **Prefeituras:** Sem padrão nacional

**Viabilidade:** 🟢 Alta (para CNDs federais)
- TST tem consulta pública
- Receita Federal tem webservice
- Pode usar web scraping + OCR

**Estimativa de Implementação:**
- Tempo: 2-3 semanas
- Custo: Gratuito (scraping) ou pago (APIs 3ª)
- Complexidade: Média

##### **1.3 Registro de Títulos e Documentos (RTD)**
**Descrição:**
- Cartórios de registro de títulos
- Protestos de duplicatas, cheques, notas promissórias
- Dados públicos mas descentralizados

**Acesso:**
- ❌ **API pública:** Não existe
- 🟡 **Web scraping:** Cada cartório tem site próprio
- ❌ **Base centralizada:** Não existe no Brasil

**Viabilidade:** 🔴 Baixa
- Muito descentralizado
- Complexidade altíssima
- ROI baixo

---

### **API 2: CONSULTA DE NOTAS FISCAIS POR CNPJ**

#### **Possíveis Fontes:**

##### **2.1 Portal da Nota Fiscal Eletrônica (NF-e)**
**Descrição:**
- Sistema oficial da SEFAZ
- NF-e emitidas/recebidas
- Dados de operações fiscais

**Acesso:**
- 🔐 **Webservice SEFAZ:** Requer certificado digital A1/A3
- 🔐 **Portal da NF-e:** Login gov.br + certificado
- ❌ **API pública:** NÃO existe sem autenticação
- ❌ **Consulta por CNPJ de 3ª:** Ilegal (sigilo fiscal)

**Viabilidade:** 🔴 Muito Baixa
- **IMPEDIMENTO LEGAL:** Dados protegidos por sigilo fiscal
- Só o próprio CNPJ pode consultar suas NF-e
- Requer certificado digital
- Não é possível criar API pública

**Legislação:**
- Lei 5.172/66 (Código Tributário Nacional) - Art. 198: Sigilo Fiscal
- Não é permitido divulgar dados de 3ª sem autorização

**Alternativas:**
- ✅ Consulta por **chave de acesso da NF-e** (44 dígitos) - PÚBLICO
- ✅ QR Code da NF-e (validação) - PÚBLICO
- ❌ Listar NF-e por CNPJ - PRIVADO (ilegal)

##### **2.2 Consulta de NF-e por Chave de Acesso**
**Descrição:**
- Validação de NF-e por chave de 44 dígitos
- Dados: emitente, destinatário, valor, data
- Público e acessível

**Acesso:**
- ✅ **Webservice SEFAZ:** Público e gratuito
- ✅ **Portal Nacional NF-e:** http://www.nfe.fazenda.gov.br
- ✅ **QR Code:** Leitura e validação

**Viabilidade:** 🟢 ALTA
- API pública gratuita
- Sem autenticação necessária
- Dados completos da NF-e

**Estimativa de Implementação:**
- Tempo: 1-2 semanas
- Custo: Gratuito
- Complexidade: Média

**Endpoints Sugeridos:**
```
GET /nfe/consulta/:chave
POST /nfe/qrcode (valida QR Code)
GET /nfe/danfe/:chave (PDF da DANFE)
```

##### **2.3 Dados Fiscais Públicos (Transparência)**
**Descrição:**
- Portal da Transparência
- Compras governamentais
- Notas fiscais de órgãos públicos

**Acesso:**
- ✅ **Portal da Transparência:** API pública
- ✅ **ComprasNet:** Dados de licitações
- ✅ **PNCP:** Portal Nacional de Contratações Públicas

**Viabilidade:** 🟢 Alta
- APIs públicas e gratuitas
- Dados abertos
- Consultável por CNPJ fornecedor

**Estimativa de Implementação:**
- Tempo: 2-3 semanas
- Custo: Gratuito
- Complexidade: Média

---

## 💡 RECOMENDAÇÕES TÉCNICAS

### **✅ VIÁVEL E RECOMENDADO:**

#### **1. API de Validação de NF-e (por chave)**
```
GET /nfe/:chave
- Valida chave de 44 dígitos
- Retorna dados da NF-e (emitente, destinatário, valor, data)
- Fonte: Webservice SEFAZ (gratuito)
- Público e legal
```

**Casos de Uso:**
- E-commerce: Validar NF-e do fornecedor
- Contabilidade: Importar NF-e para sistemas
- Marketplaces: Verificar autenticidade de NF-e
- Auditoria: Validar documentos fiscais

**Complexidade:** Média  
**ROI:** Alto (muita demanda)  
**Legal:** ✅ Totalmente legal

---

#### **2. API de Certidões Negativas de Débitos**
```
GET /cnd/:cnpj
- Consulta CND federal (Receita)
- Consulta CNDT (débitos trabalhistas)
- Status: Regular/Irregular
- Fonte: TST + Receita Federal
```

**Casos de Uso:**
- Due diligence de fornecedores
- Pré-contratação (verificar regularidade)
- Licitações (exigência de CND)
- Crédito (análise de risco)

**Complexidade:** Média  
**ROI:** Alto  
**Legal:** ✅ Dados públicos

---

#### **3. API de Protestos (Registro Público)**
```
GET /protestos/:cnpj
- Títulos protestados
- Histórico de protestos
- Valores protestados
- Fonte: Cartórios (web scraping ou APIs pagas)
```

**Casos de Uso:**
- Análise de crédito
- Due diligence
- Compliance
- Risk assessment

**Complexidade:** Alta (descentralizado)  
**ROI:** Médio  
**Legal:** ✅ Dados públicos (mas dispersos)

---

#### **4. API de Compras Governamentais**
```
GET /compras-gov/:cnpj
- Licitações vencidas
- Contratos com governo
- Volume de vendas para setor público
- Fonte: Portal da Transparência (API pública)
```

**Casos de Uso:**
- Due diligence de fornecedores
- Análise de mercado
- Inteligência comercial
- Compliance

**Complexidade:** Média  
**ROI:** Médio-Alto  
**Legal:** ✅ Dados abertos

---

### **❌ NÃO VIÁVEL:**

#### **Busca de Boletos por CNPJ (genérico)**
**Motivos:**
- ❌ Boletos são dados privados (sigilo bancário)
- ❌ Não existe base pública centralizada
- ❌ Bancos não compartilham essas informações
- ❌ Seria ilegal disponibilizar

**Exceção:**
- ✅ Consulta de boleto específico (código de barras/linha digitável) - VIÁVEL
- ✅ Validação de boleto - VIÁVEL

#### **Busca de NF-e por CNPJ (listar todas)**
**Motivos:**
- ❌ Protegido por sigilo fiscal (CTN Art. 198)
- ❌ Só o próprio CNPJ pode acessar
- ❌ Requer certificado digital
- ❌ Seria crime divulgar (Lei 5.172/66)

**Exceção:**
- ✅ Validação de NF-e por chave (44 dígitos) - VIÁVEL E PÚBLICO

---

## 🎯 PROPOSTA FINAL (APIs Viáveis)

### **Fase 3: Validação e Compliance**

#### **✅ 1. API de Validação de NF-e**
```
Endpoints:
- GET /nfe/:chave (consulta por chave de 44 dígitos)
- POST /nfe/qrcode (valida QR Code)
- GET /nfe/danfe/:chave (PDF da DANFE)

Dados retornados:
- Emitente (CNPJ, razão social, endereço)
- Destinatário (CNPJ, nome)
- Valor total, ICMS, IPI
- Data de emissão
- Status (autorizada, cancelada, denegada)
- Chave de acesso

Fonte: Webservice SEFAZ (gratuito)
Cache: 30 dias (NF-e não muda)
Performance: ~500ms (SEFAZ lento)
```

#### **✅ 2. API de Certidões (CND/CNDT)**
```
Endpoints:
- GET /certidoes/:cnpj
- GET /certidoes/federal/:cnpj
- GET /certidoes/trabalhista/:cnpj

Dados retornados:
- Status: Regular/Irregular
- Data de emissão
- Validade
- Link para certidão PDF
- Débitos pendentes (se houver)

Fontes:
- Receita Federal (e-CAC)
- TST (Certidão Negativa Trabalhista)

Cache: 1 dia (atualiza frequente)
Performance: ~2s (scraping + OCR)
```

#### **✅ 3. API de Protestos**
```
Endpoints:
- GET /protestos/:cnpj
- GET /protestos/:cnpj/resumo

Dados retornados:
- Total de protestos
- Valor total protestado
- Último protesto (data)
- Cartório responsável
- Situação (protestado, pago, cancelado)

Fontes:
- Cartórios (scraping ou APIs pagas)
- Serasa (pago)
- Boa Vista SCPC (pago)

Cache: 7 dias
Performance: ~5s (scraping) / ~500ms (API paga)
Custo: R$ 0,50 - R$ 2,00 por consulta
```

#### **✅ 4. API de Compras Governamentais**
```
Endpoints:
- GET /compras-gov/:cnpj
- GET /compras-gov/:cnpj/licitacoes
- GET /compras-gov/:cnpj/contratos

Dados retornados:
- Licitações vencidas
- Contratos ativos
- Valor total contratado
- Órgãos contratantes
- Histórico de 5 anos

Fonte: Portal da Transparência + ComprasNet (APIs públicas)
Cache: 7 dias
Performance: ~1s
Custo: Gratuito
```

---

## 📊 MATRIZ DE VIABILIDADE

| API | Viabilidade | Custo | Complexidade | ROI | Legal | Prioridade |
|-----|-------------|-------|--------------|-----|-------|------------|
| **Validação NF-e** | 🟢 Alta | Grátis | Média | Alto | ✅ Sim | 🔥 Alta |
| **Certidões (CND)** | 🟢 Alta | Grátis | Média | Alto | ✅ Sim | 🔥 Alta |
| **Protestos** | 🟡 Média | R$ 0,50-2/req | Alta | Médio | ✅ Sim | 🟡 Média |
| **Compras Gov** | 🟢 Alta | Grátis | Baixa | Médio | ✅ Sim | 🟢 Média |
| ~~Listar NF-e por CNPJ~~ | 🔴 Impossível | - | - | - | ❌ Ilegal | ❌ Descartada |
| ~~Boletos por CNPJ~~ | 🔴 Impossível | - | - | - | ❌ Ilegal | ❌ Descartada |

---

## 🚀 ROADMAP DE IMPLEMENTAÇÃO

### **Fase 3 (3-6 meses):**

#### **✅ Prioridade ALTA:**
1. **API de Validação de NF-e** (2 semanas)
   - Integração com webservice SEFAZ
   - Validação de chave de 44 dígitos
   - QR Code validation
   - Cache de 30 dias

2. **API de Certidões** (3 semanas)
   - CND Federal (Receita)
   - CNDT (TST)
   - Web scraping + OCR
   - Cache de 1 dia

#### **🟡 Prioridade MÉDIA:**
3. **API de Compras Governamentais** (2 semanas)
   - Portal da Transparência
   - ComprasNet
   - API pública gratuita
   - Cache de 7 dias

#### **🟡 Prioridade MÉDIA-BAIXA:**
4. **API de Protestos** (4-6 semanas)
   - Integração com Serasa (pago) OU
   - Web scraping de cartórios (gratuito mas complexo)
   - Cache de 7 dias
   - Decisão: custo x benefício

---

## 💰 ANÁLISE DE CUSTO-BENEFÍCIO

### **APIs Gratuitas (Recomendadas):**

| API | Fonte | Custo | Valor para Cliente |
|-----|-------|-------|-------------------|
| **NF-e Validation** | SEFAZ | R$ 0 | Alto (validação de fornecedores) |
| **Certidões** | TST + Receita | R$ 0 | Alto (compliance) |
| **Compras Gov** | Transparência | R$ 0 | Médio (inteligência comercial) |

**Total de custo:** R$ 0  
**Valor agregado:** Alto  
**ROI:** Excelente

### **APIs Pagas (Avaliar):**

| API | Fonte | Custo/req | Volume Estimado | Custo Mensal |
|-----|-------|-----------|----------------|--------------|
| **Protestos** | Serasa | R$ 1,50 | 1.000 req/mês | R$ 1.500 |
| **Protestos** | Boa Vista | R$ 2,00 | 1.000 req/mês | R$ 2.000 |

**Estratégia:**
1. **Começar com web scraping** (custo zero)
2. **Se houver demanda**, migrar para API paga
3. **Repassar custo** ao cliente (R$ 2,00-5,00 por consulta)

---

## 🛠️ IMPLEMENTAÇÃO TÉCNICA

### **1. API de Validação de NF-e**

**Webservice SEFAZ:**
```go
// Consulta NF-e por chave de 44 dígitos
type NFeConsultaRequest struct {
    Chave string `json:"chave"` // 44 dígitos
    UF    string `json:"uf"`    // Estado emissor
}

// Webservice SOAP da SEFAZ
func ConsultarNFe(chave string) (*NFeResponse, error) {
    // 1. Extrair UF da chave (posições 0-1)
    uf := chave[0:2]
    
    // 2. Montar SOAP request
    soapBody := fmt.Sprintf(`
        <nfeDadosMsg xmlns="http://www.portalfiscal.inf.br/nfe">
            <consSitNFe versao="4.00" xmlns="http://www.portalfiscal.inf.br/nfe">
                <tpAmb>1</tpAmb>
                <xServ>CONSULTAR</xServ>
                <chNFe>%s</chNFe>
            </consSitNFe>
        </nfeDadosMsg>
    `, chave)
    
    // 3. Enviar para webservice da SEFAZ
    url := fmt.Sprintf("https://nfe.sefaz.%s.gov.br/ws/NFeConsulta4", uf)
    
    // 4. Parsear resposta XML
    // 5. Retornar JSON
}
```

**Cache Strategy:**
- 30 dias (NF-e não muda após autorização)
- MongoDB collection: `nfe_cache`

**Rate Limiting:**
- SEFAZ limita: ~1 req/segundo por IP
- Nossa API: 100 req/minuto por tenant

---

### **2. API de Certidões**

**TST (Web Scraping):**
```go
func ConsultarCNDT(cnpj string) (*CertidaoResponse, error) {
    // 1. Acessar https://www.tst.jus.br/certidao
    // 2. Resolver captcha (2captcha ou anti-captcha)
    // 3. Submeter CNPJ
    // 4. Parsear resposta HTML
    // 5. OCR do PDF (se necessário)
    // 6. Retornar JSON
}
```

**Receita Federal (API ou Scraping):**
```go
func ConsultarCNDFederal(cnpj string) (*CertidaoResponse, error) {
    // Opção 1: Webservice Receita (requer certificado)
    // Opção 2: Web scraping (público)
}
```

---

## 📝 DOCUMENTAÇÃO DE ENDPOINTS

### **API de NF-e:**

```yaml
/nfe/{chave}:
  get:
    summary: Consulta NF-e por chave de acesso
    parameters:
      - name: chave
        in: path
        required: true
        schema:
          type: string
          pattern: '^[0-9]{44}$'
    responses:
      200:
        description: NF-e encontrada
        content:
          application/json:
            schema:
              type: object
              properties:
                chave:
                  type: string
                numero:
                  type: string
                serie:
                  type: string
                dataEmissao:
                  type: string
                emitente:
                  type: object
                  properties:
                    cnpj: string
                    razaoSocial: string
                    uf: string
                destinatario:
                  type: object
                valorTotal:
                  type: number
                situacao:
                  type: string
                  enum: [AUTORIZADA, CANCELADA, DENEGADA]
```

---

## 🎯 CONCLUSÃO E RECOMENDAÇÕES

### **✅ IMPLEMENTAR (Fase 3):**

1. **API de Validação de NF-e** 🔥
   - Viabilidade: ALTA
   - Custo: ZERO
   - Legal: SIM
   - Demanda: ALTA

2. **API de Certidões (CND/CNDT)** 🔥
   - Viabilidade: ALTA
   - Custo: ZERO
   - Legal: SIM
   - Demanda: ALTA

3. **API de Compras Governamentais** ✅
   - Viabilidade: ALTA
   - Custo: ZERO
   - Legal: SIM
   - Demanda: MÉDIA

### **🟡 AVALIAR DEPOIS (Fase 4):**

4. **API de Protestos**
   - Avaliar demanda
   - Testar web scraping primeiro
   - Se houver demanda, contratar API paga

### **❌ NÃO IMPLEMENTAR:**

- ~~Busca de boletos por CNPJ~~ (ilegal - sigilo bancário)
- ~~Listar NF-e por CNPJ~~ (ilegal - sigilo fiscal)

---

## 📋 CHECKLIST DE IMPLEMENTAÇÃO

### **Para cada nova API:**

- [ ] Pesquisar fonte de dados oficial
- [ ] Verificar legalidade (sigilo fiscal/bancário)
- [ ] Testar webservice/scraping
- [ ] Implementar handler Go
- [ ] Adicionar cache MongoDB
- [ ] Configurar rate limiting
- [ ] Criar scope (permissão)
- [ ] Atualizar Redoc
- [ ] Atualizar landing page
- [ ] Criar landing page dedicada
- [ ] Adicionar ao playground
- [ ] Criar ferramenta pública (se aplicável)
- [ ] Documentar em FONTES_DE_DADOS.md
- [ ] Atualizar ROADMAP.md

---

## 📚 FONTES E REFERÊNCIAS

### **Legislação:**
- Lei 5.172/66 (Código Tributário Nacional) - Sigilo Fiscal
- Lei Complementar 105/2001 - Sigilo Bancário
- Lei 12.527/2011 - Lei de Acesso à Informação

### **Webservices Oficiais:**
- SEFAZ NF-e: http://www.nfe.fazenda.gov.br/portal/webServices.aspx
- Portal da Transparência: https://portaldatransparencia.gov.br/api-de-dados
- TST Certidões: https://www.tst.jus.br/certidao
- ComprasNet: https://www.gov.br/compras/pt-br

### **APIs de Terceiros (Pagas):**
- Serasa Experian: https://www.serasaexperian.com.br/api
- Boa Vista SCPC: https://www.boavistaservicos.com.br
- Serpro: https://www.serpro.gov.br/links-fixos-superiores/api-serpro

---

**🎉 Análise completa realizada!**

**Próximo passo:** Atualizar ROADMAP.md com as APIs viáveis! 🚀

