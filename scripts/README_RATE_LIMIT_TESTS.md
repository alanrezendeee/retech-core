# 🧪 Testes de Rate Limiting - Retech Core API

Este diretório contém scripts automatizados para testar o sistema de **Rate Limiting** da API Retech Core.

---

## 📋 **O QUE ESTES SCRIPTS FAZEM?**

### 1. `test-rate-limit.sh`
Script bash que:
- ✅ Faz login como admin
- ✅ Cria 3 tenants de teste com limites diferentes:
  - **Tenant 1**: 5 req/dia, 2 req/min (limite baixo)
  - **Tenant 2**: 50 req/dia, 10 req/min (limite médio)
  - **Tenant 3**: 1000 req/dia, 60 req/min (limite padrão)
- ✅ Gera API keys para cada tenant
- ✅ Executa testes fazendo múltiplas requisições
- ✅ Verifica se recebe `429 Too Many Requests` no momento certo
- ✅ Salva resultados em:
  - `rate-limit-test-results.txt` (texto)
  - `rate-limit-test-results.json` (JSON)

### 2. `generate-rate-limit-chart.js`
Script Node.js que:
- ✅ Lê os resultados do teste
- ✅ Gera visualização ASCII colorida no terminal
- ✅ Mostra gráficos de barras
- ✅ Calcula estatísticas (% sucesso, % bloqueado)
- ✅ Gera relatório HTML interativo com gráficos Chart.js

---

## 🚀 **COMO USAR**

### Pré-requisitos
- Bash (Linux/macOS)
- Node.js (para gerar gráficos)
- `curl` e `jq` instalados
- API rodando (local ou produção)
- Credenciais de admin

### Passo 1: Executar Testes

```bash
cd /path/to/retech-core

# Definir variáveis (opcional)
export API_URL="https://api-core.theretech.com.br"
export ADMIN_EMAIL="admin@theretech.com.br"
export ADMIN_PASSWORD="Admin@123"

# Executar testes
./scripts/test-rate-limit.sh
```

O script vai:
1. Criar 3 tenants automaticamente
2. Fazer dezenas de requisições
3. Monitorar as respostas (200, 429, etc.)
4. Salvar resultados

**Tempo estimado**: 2-3 minutos

---

### Passo 2: Visualizar Resultados

#### Opção A: Ver no terminal (texto simples)
```bash
cat rate-limit-test-results.txt
```

#### Opção B: Ver no terminal (gráficos ASCII coloridos)
```bash
node scripts/generate-rate-limit-chart.js
```

**Saída esperada**:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   📊 RELATÓRIO DE TESTES DE RATE LIMITING
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

RESUMO GERAL:
  Total de testes:    3
  ✅ Passaram:        3
  ❌ Falharam:        0
  Taxa de sucesso:   100.0%

DETALHES POR CENÁRIO:

1. Limite Baixo (5/dia) ✅
   Limite esperado: 5 requests
   Requests feitas: 10

   ✅ Sucesso (5):
      ██████████████████████████ 50.0%

   🚫 Rate Limited (5):
      ██████████████████████████ 50.0%

   🎯 Primeiro 429 na request: #6

   ANÁLISE:
   ✅ Rate limit funcionou corretamente!
      - 5 requests permitidas (≤ 5)
      - 5 requests bloqueadas
      - Bloqueio ocorreu exatamente após o limite
```

#### Opção C: Ver relatório HTML interativo
```bash
# Gerar relatório
node scripts/generate-rate-limit-chart.js

# Abrir no navegador (macOS)
open rate-limit-test-report.html

# Ou Linux
xdg-open rate-limit-test-report.html
```

O relatório HTML inclui:
- 📊 Gráficos interativos (Chart.js)
- 📈 Estatísticas detalhadas
- 🎨 Design moderno e responsivo
- 📱 Funciona em mobile

---

## 📊 **EXEMPLO DE RESULTADOS**

### ✅ Teste PASSOU
```
1. Limite Baixo (5/dia) ✅
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RESULTADOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Sucesso:        5
🚫 Rate Limited:   5
⚠️  Erros:          0

🎯 Primeiro 429 na request: #6

✅ TESTE PASSOU! Rate limit funcionou corretamente.
```

### ❌ Teste FALHOU
```
1. Limite Baixo (5/dia) ❌
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RESULTADOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Sucesso:        10
🚫 Rate Limited:   0
⚠️  Erros:          0

🎯 Nenhum 429 recebido

❌ TESTE FALHOU! Problema detectado:
   - Permitiu MAIS requests que o limite (10 > 5)
   - Nenhuma request foi bloqueada (429 não retornado)
```

---

## 🔍 **O QUE O TESTE VERIFICA?**

### ✅ Comportamento Esperado
1. **Limite diário respeitado**
   - Primeiras N requests retornam `200 OK`
   - Request N+1 retorna `429 Too Many Requests`

2. **Headers corretos**
   - `X-RateLimit-Limit`: limite configurado
   - `X-RateLimit-Remaining`: requests restantes
   - `X-RateLimit-Reset`: timestamp do reset

3. **Configuração por tenant**
   - Cada tenant tem seu limite próprio
   - Limites personalizados funcionam
   - Limite padrão (1000/dia) como fallback

### ❌ Problemas Detectados
- Requests acima do limite são permitidas
- Nenhum `429` é retornado
- `429` retornado na request errada
- Erros inesperados (500, 401, etc.)

---

## 🐛 **TROUBLESHOOTING**

### Erro: "Failed to create tenant"
**Causa**: Admin login falhou ou permissões insuficientes  
**Solução**: 
```bash
# Verificar credenciais
curl -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@theretech.com.br","password":"Admin@123"}'

# Verificar se é super admin
# Role deve ser "SUPER_ADMIN"
```

### Erro: "API Key não funciona"
**Causa**: API key não está sendo aceita  
**Solução**:
```bash
# Testar manualmente
curl -H "X-API-Key: rtc_..." https://api-core.theretech.com.br/geo/estados

# Verificar no banco
db.api_keys.find({keyId: "..."}).pretty()
```

### Todos os testes falharam
**Causa**: Rate limiting não está implementado ou tem bugs  
**Solução**: Ver `PENDENCIAS_CRITICAS.md` para correções necessárias

---

## 📝 **ARQUIVOS GERADOS**

Após executar os testes, você terá:

```
retech-core/
├── rate-limit-test-results.txt      # Resultados em texto
├── rate-limit-test-results.json     # Resultados em JSON
└── rate-limit-test-report.html      # Relatório visual HTML
```

**Não commitar** estes arquivos (já estão no `.gitignore`)

---

## 🎯 **CENÁRIOS DE TESTE**

### Cenário 1: Limite Baixo (5/dia)
- **Objetivo**: Testar bloqueio rápido
- **Config**: 5 req/dia, 2 req/min
- **Teste**: Fazer 10 requests
- **Esperado**: 5 sucesso + 5 bloqueadas

### Cenário 2: Limite Médio (50/dia)
- **Objetivo**: Testar limite realista
- **Config**: 50 req/dia, 10 req/min
- **Teste**: Fazer 60 requests
- **Esperado**: 50 sucesso + 10 bloqueadas

### Cenário 3: Limite Padrão (1000/dia)
- **Objetivo**: Verificar que limite alto funciona
- **Config**: 1000 req/dia, 60 req/min
- **Teste**: Fazer apenas 20 requests
- **Esperado**: Todas 20 com sucesso (bem abaixo do limite)

---

## 🔧 **CUSTOMIZAR TESTES**

### Adicionar novo cenário

Edite `test-rate-limit.sh`:

```bash
# Criar tenant com configuração customizada
APIKEY_4=$(create_test_tenant "Test-Custom" "test-custom-$(date +%s)@test.com" 100 5)

# Testar com parâmetros específicos
test_rate_limit "Meu Cenário Custom" "$APIKEY_4" 100 120 0
#                nome                 apikey         limite requests sleep
```

### Testar limite por minuto

```bash
# Fazer requests muito rápidas (sem sleep)
test_rate_limit "Teste por Minuto" "$APIKEY" 2 10 0
#                                             ^ limite/min
```

---

## 📊 **INTERPRETAR RESULTADOS**

### Taxa de Sucesso = 100%
✅ **EXCELENTE!** Rate limiting funcionando perfeitamente

### Taxa de Sucesso = 0%
❌ **CRÍTICO!** Rate limiting não está funcionando

### Taxa de Sucesso = 50-70%
🟡 **PARCIAL** - Alguns cenários funcionam, outros não

---

## 🚀 **CI/CD**

Para rodar em pipeline:

```yaml
# GitHub Actions / GitLab CI
test-rate-limit:
  script:
    - ./scripts/test-rate-limit.sh
    - node scripts/generate-rate-limit-chart.js
    - if [ $? -eq 0 ]; then echo "Rate limit OK"; else exit 1; fi
```

---

## 📞 **SUPORTE**

Dúvidas ou problemas com os testes?
- 📧 Email: contato@theretech.com.br
- 📚 Docs: Ver `PENDENCIAS_CRITICAS.md`

---

**Criado por**: The Retech Team  
**Última atualização**: 2025-10-22  
**Versão**: 1.0

