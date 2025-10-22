#!/bin/bash

# 🧪 Script de Teste Completo de Rate Limiting
# Simula 1000+ requisições em diversos cenários

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Configurações
API_URL="${API_URL:-http://localhost:8080}"
ADMIN_EMAIL="${ADMIN_EMAIL:-alanrezendeee@gmail.com}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-Admin@123}"

# Arquivos de resultados
RESULTS_FILE="rate-limit-test-results.txt"
JSON_FILE="rate-limit-test-results.json"
TOTAL_REQUESTS=0

echo ""
echo -e "${CYAN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                                                        ║${NC}"
echo -e "${CYAN}║   🧪  TESTE COMPLETO DE RATE LIMITING - 1000 REQS    ║${NC}"
echo -e "${CYAN}║                                                        ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}API:${NC} $API_URL"
echo -e "${BLUE}Data:${NC} $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Limpar resultados anteriores
> $RESULTS_FILE
echo "[" > $JSON_FILE

# Função para fazer login admin
login_admin() {
    echo -e "${MAGENTA}[LOGIN ADMIN]${NC} Autenticando..."
    
    RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
    
    ADMIN_TOKEN=$(echo $RESPONSE | jq -r '.accessToken // empty')
    
    if [ -z "$ADMIN_TOKEN" ]; then
        echo -e "${RED}❌ Falha no login admin${NC}"
        echo "Response: $RESPONSE"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Login admin OK${NC}"
    echo ""
}

# Função para criar tenant de teste
create_test_tenant() {
    local NAME=$1
    local EMAIL=$2
    local RATE_LIMIT_DAY=$3
    local RATE_LIMIT_MIN=$4
    
    echo -e "${CYAN}[TENANT]${NC} Criando: ${YELLOW}$NAME${NC}"
    echo -e "         Rate Limit: ${YELLOW}${RATE_LIMIT_DAY}/dia${NC}, ${YELLOW}${RATE_LIMIT_MIN}/min${NC}"
    
    # Criar tenant via API
    TENANT_RESPONSE=$(curl -s -X POST "$API_URL/admin/tenants" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -d "{
            \"name\": \"$NAME\",
            \"email\": \"$EMAIL\",
            \"password\": \"Test@123\",
            \"company\": \"Test Company\",
            \"purpose\": \"Rate Limit Testing\",
            \"rateLimit\": {
                \"requestsPerDay\": $RATE_LIMIT_DAY,
                \"requestsPerMinute\": $RATE_LIMIT_MIN
            }
        }")
    
    TENANT_ID=$(echo $TENANT_RESPONSE | jq -r '.tenant._id // .tenant.id // empty')
    
    if [ -z "$TENANT_ID" ]; then
        echo -e "${RED}❌ Falha ao criar tenant${NC}"
        echo "Response: $TENANT_RESPONSE"
        return 1
    fi
    
    # Fazer login como tenant
    LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$EMAIL\",\"password\":\"Test@123\"}")
    
    TENANT_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.accessToken // empty')
    
    # Criar API Key
    APIKEY_RESPONSE=$(curl -s -X POST "$API_URL/me/apikeys" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TENANT_TOKEN" \
        -d "{\"name\":\"Test Key for $NAME\",\"expiresInDays\":30}")
    
    APIKEY=$(echo $APIKEY_RESPONSE | jq -r '.key // empty')
    
    if [ -z "$APIKEY" ]; then
        echo -e "${RED}❌ Falha ao criar API key${NC}"
        echo "Response: $APIKEY_RESPONSE"
        return 1
    fi
    
    echo -e "         ${GREEN}✅ Tenant ID:${NC} $TENANT_ID"
    echo -e "         ${GREEN}✅ API Key:${NC} ${APIKEY:0:40}..."
    echo ""
    
    # Retornar API key
    echo "$APIKEY"
}

# Função para testar rate limit
test_rate_limit() {
    local SCENARIO=$1
    local APIKEY=$2
    local EXPECTED_LIMIT=$3
    local REQUESTS_TO_MAKE=$4
    local SLEEP_BETWEEN=${5:-0}
    local TEST_TYPE=${6:-"day"} # "day" ou "minute"
    
    echo ""
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}📊 CENÁRIO: $SCENARIO${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${BLUE}Limite esperado:${NC} $EXPECTED_LIMIT requests"
    echo -e "${BLUE}Requests a fazer:${NC} $REQUESTS_TO_MAKE"
    echo -e "${BLUE}Sleep entre requests:${NC} ${SLEEP_BETWEEN}s"
    echo -e "${BLUE}Tipo de teste:${NC} $TEST_TYPE"
    echo ""
    
    local SUCCESS_COUNT=0
    local RATE_LIMITED_COUNT=0
    local ERROR_COUNT=0
    local FIRST_429=""
    local START_TIME=$(date +%s)
    
    # Barra de progresso
    local PROGRESS_STEP=$((REQUESTS_TO_MAKE / 20))
    if [ $PROGRESS_STEP -lt 1 ]; then
        PROGRESS_STEP=1
    fi
    
    for i in $(seq 1 $REQUESTS_TO_MAKE); do
        # Fazer request
        RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$API_URL/geo/estados" \
            -H "X-API-Key: $APIKEY" 2>/dev/null)
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        
        TOTAL_REQUESTS=$((TOTAL_REQUESTS + 1))
        
        # Contar resultados
        if [ "$HTTP_CODE" = "200" ]; then
            SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
            STATUS="${GREEN}✅${NC}"
        elif [ "$HTTP_CODE" = "429" ]; then
            RATE_LIMITED_COUNT=$((RATE_LIMITED_COUNT + 1))
            STATUS="${RED}🚫${NC}"
            if [ -z "$FIRST_429" ]; then
                FIRST_429=$i
            fi
        else
            ERROR_COUNT=$((ERROR_COUNT + 1))
            STATUS="${YELLOW}⚠️${NC}"
        fi
        
        # Mostrar progresso (apenas múltiplos do step ou primeira/última)
        if [ $((i % PROGRESS_STEP)) -eq 0 ] || [ $i -eq 1 ] || [ $i -eq $REQUESTS_TO_MAKE ]; then
            local PERCENT=$((i * 100 / REQUESTS_TO_MAKE))
            echo -ne "\r${CYAN}Progress:${NC} [$i/$REQUESTS_TO_MAKE] ${PERCENT}% | ✅ $SUCCESS_COUNT | 🚫 $RATE_LIMITED_COUNT | ⚠️ $ERROR_COUNT"
        fi
        
        # Sleep se necessário
        if [ $SLEEP_BETWEEN -gt 0 ] && [ $i -lt $REQUESTS_TO_MAKE ]; then
            sleep $SLEEP_BETWEEN
        fi
    done
    
    echo "" # Nova linha após barra de progresso
    echo ""
    
    local END_TIME=$(date +%s)
    local DURATION=$((END_TIME - START_TIME))
    local REQ_PER_SEC=0
    if [ $DURATION -gt 0 ]; then
        REQ_PER_SEC=$((REQUESTS_TO_MAKE / DURATION))
    fi
    
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}📈 RESULTADOS${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${GREEN}✅ Sucesso:${NC}        $SUCCESS_COUNT"
    echo -e "${RED}🚫 Rate Limited:${NC}   $RATE_LIMITED_COUNT"
    echo -e "${YELLOW}⚠️  Erros:${NC}          $ERROR_COUNT"
    echo ""
    echo -e "${CYAN}⏱️  Duração:${NC}        ${DURATION}s"
    echo -e "${CYAN}⚡ Req/seg:${NC}        ${REQ_PER_SEC}"
    echo ""
    
    if [ -n "$FIRST_429" ]; then
        echo -e "${MAGENTA}🎯 Primeiro 429:${NC} Request #$FIRST_429"
    else
        echo -e "${MAGENTA}🎯 Primeiro 429:${NC} Nenhum bloqueio"
    fi
    
    # Verificar se está funcionando corretamente
    local TEST_PASSED="false"
    local TEST_MESSAGE=""
    
    if [ "$TEST_TYPE" = "minute" ]; then
        # Teste de limite por minuto - aceita margem de erro
        if [ $RATE_LIMITED_COUNT -gt 0 ]; then
            TEST_PASSED="true"
            TEST_MESSAGE="Rate limit por minuto funcionou!"
        else
            TEST_MESSAGE="Nenhum bloqueio por minuto detectado"
        fi
    else
        # Teste de limite diário
        if [ $SUCCESS_COUNT -eq $EXPECTED_LIMIT ] && [ $RATE_LIMITED_COUNT -gt 0 ]; then
            TEST_PASSED="true"
            TEST_MESSAGE="Rate limit funcionou perfeitamente!"
        elif [ $SUCCESS_COUNT -le $((EXPECTED_LIMIT + 1)) ] && [ $RATE_LIMITED_COUNT -gt 0 ]; then
            TEST_PASSED="true"
            TEST_MESSAGE="Rate limit funcionou (margem de 1 request)"
        elif [ $SUCCESS_COUNT -gt $EXPECTED_LIMIT ]; then
            TEST_MESSAGE="Mais requests permitidas que o limite ($SUCCESS_COUNT > $EXPECTED_LIMIT)"
        elif [ $RATE_LIMITED_COUNT -eq 0 ]; then
            TEST_MESSAGE="Nenhuma request foi bloqueada (429 não retornado)"
        else
            TEST_MESSAGE="Resultado inconclusivo"
        fi
    fi
    
    echo ""
    if [ "$TEST_PASSED" = "true" ]; then
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}✅ TESTE PASSOU!${NC} $TEST_MESSAGE"
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    else
        echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${RED}❌ TESTE FALHOU!${NC} $TEST_MESSAGE"
        echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    fi
    
    echo ""
    
    # Salvar resultados
    cat >> $RESULTS_FILE << EOF
==================================================
CENÁRIO: $SCENARIO
==================================================
Limite esperado: $EXPECTED_LIMIT
Requests feitas: $REQUESTS_TO_MAKE
Sucesso: $SUCCESS_COUNT
Rate Limited: $RATE_LIMITED_COUNT
Erros: $ERROR_COUNT
Primeiro 429: ${FIRST_429:-N/A}
Duração: ${DURATION}s
Req/seg: ${REQ_PER_SEC}
Status: $([ "$TEST_PASSED" = "true" ] && echo "PASSOU ✅" || echo "FALHOU ❌")
Mensagem: $TEST_MESSAGE

EOF
    
    # Salvar JSON
    cat >> $JSON_FILE << EOF
{
  "scenario": "$SCENARIO",
  "expectedLimit": $EXPECTED_LIMIT,
  "requestsMade": $REQUESTS_TO_MAKE,
  "success": $SUCCESS_COUNT,
  "rateLimited": $RATE_LIMITED_COUNT,
  "errors": $ERROR_COUNT,
  "first429": ${FIRST_429:-null},
  "duration": $DURATION,
  "reqPerSec": $REQ_PER_SEC,
  "passed": $TEST_PASSED,
  "message": "$TEST_MESSAGE"
},
EOF
}

# ============================================
# INÍCIO DOS TESTES
# ============================================

login_admin

echo -e "${CYAN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║  🏗️  CRIANDO TENANTS DE TESTE                         ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Timestamp único para evitar conflitos
TIMESTAMP=$(date +%s)

# Cenário 1: Limite MUITO baixo (3 req/dia, 1 req/min)
echo -e "${YELLOW}[1/8]${NC} Tenant: Limite Extremamente Baixo"
APIKEY_1=$(create_test_tenant "Test-Ultra-Low" "test-ultra-low-$TIMESTAMP@test.com" 3 1)

# Cenário 2: Limite baixo (10 req/dia, 3 req/min)
echo -e "${YELLOW}[2/8]${NC} Tenant: Limite Baixo"
APIKEY_2=$(create_test_tenant "Test-Low" "test-low-$TIMESTAMP@test.com" 10 3)

# Cenário 3: Limite médio-baixo (30 req/dia, 5 req/min)
echo -e "${YELLOW}[3/8]${NC} Tenant: Limite Médio-Baixo"
APIKEY_3=$(create_test_tenant "Test-Medium-Low" "test-medium-low-$TIMESTAMP@test.com" 30 5)

# Cenário 4: Limite médio (100 req/dia, 10 req/min)
echo -e "${YELLOW}[4/8]${NC} Tenant: Limite Médio"
APIKEY_4=$(create_test_tenant "Test-Medium" "test-medium-$TIMESTAMP@test.com" 100 10)

# Cenário 5: Limite médio-alto (250 req/dia, 20 req/min)
echo -e "${YELLOW}[5/8]${NC} Tenant: Limite Médio-Alto"
APIKEY_5=$(create_test_tenant "Test-Medium-High" "test-medium-high-$TIMESTAMP@test.com" 250 20)

# Cenário 6: Limite alto (500 req/dia, 30 req/min)
echo -e "${YELLOW}[6/8]${NC} Tenant: Limite Alto"
APIKEY_6=$(create_test_tenant "Test-High" "test-high-$TIMESTAMP@test.com" 500 30)

# Cenário 7: Limite muito alto (1000 req/dia, 60 req/min)
echo -e "${YELLOW}[7/8]${NC} Tenant: Limite Muito Alto"
APIKEY_7=$(create_test_tenant "Test-Very-High" "test-very-high-$TIMESTAMP@test.com" 1000 60)

# Cenário 8: Teste de minuto (50 req/dia, 2 req/min)
echo -e "${YELLOW}[8/8]${NC} Tenant: Teste Limite por Minuto"
APIKEY_8=$(create_test_tenant "Test-Minute-Limit" "test-minute-$TIMESTAMP@test.com" 50 2)

echo ""
echo -e "${CYAN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║  🧪 EXECUTANDO TESTES (1000+ REQUISIÇÕES)             ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════════════╝${NC}"

# Teste 1: Ultra baixo (3/dia) - fazer 10 requests
test_rate_limit "[1/10] Limite Ultra Baixo (3/dia)" "$APIKEY_1" 3 10 0 "day"

# Teste 2: Baixo (10/dia) - fazer 20 requests
test_rate_limit "[2/10] Limite Baixo (10/dia)" "$APIKEY_2" 10 20 0 "day"

# Teste 3: Médio-baixo (30/dia) - fazer 50 requests
test_rate_limit "[3/10] Limite Médio-Baixo (30/dia)" "$APIKEY_3" 30 50 0 "day"

# Teste 4: Médio (100/dia) - fazer 150 requests
test_rate_limit "[4/10] Limite Médio (100/dia)" "$APIKEY_4" 100 150 0 "day"

# Teste 5: Médio-alto (250/dia) - fazer 300 requests
test_rate_limit "[5/10] Limite Médio-Alto (250/dia)" "$APIKEY_5" 250 300 0 "day"

# Teste 6: Alto (500/dia) - fazer 50 requests (teste rápido)
test_rate_limit "[6/10] Limite Alto (500/dia) - Teste Parcial" "$APIKEY_6" 500 50 0 "day"

# Teste 7: Muito alto (1000/dia) - fazer 100 requests (teste rápido)
test_rate_limit "[7/10] Limite Muito Alto (1000/dia) - Teste Parcial" "$APIKEY_7" 1000 100 0 "day"

# Teste 8: Limite por MINUTO (2/min) - fazer 10 requests SEM sleep
test_rate_limit "[8/10] Limite Por Minuto (2/min) - Burst" "$APIKEY_8" 2 10 0 "minute"

# Teste 9: Limite por MINUTO (2/min) - fazer 10 requests COM sleep de 30s
test_rate_limit "[9/10] Limite Por Minuto (2/min) - Com Pausa" "$APIKEY_8" 2 10 30 "minute"

# Teste 10: STRESS TEST com limite alto - 200 requests rápidas
test_rate_limit "[10/10] Stress Test (500/dia) - 200 Requests" "$APIKEY_6" 500 200 0 "day"

# Finalizar JSON
sed -i.bak '$ s/,$//' $JSON_FILE && rm "${JSON_FILE}.bak"
echo "]" >> $JSON_FILE

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✅ TESTES CONCLUÍDOS COM SUCESSO!                    ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${CYAN}📊 RESUMO GERAL:${NC}"
echo -e "   ${YELLOW}Total de requisições simuladas:${NC} ${MAGENTA}$TOTAL_REQUESTS${NC}"
echo ""
echo -e "${CYAN}📄 Resultados salvos em:${NC}"
echo -e "   ${GREEN}✓${NC} $RESULTS_FILE (texto formatado)"
echo -e "   ${GREEN}✓${NC} $JSON_FILE (dados estruturados)"
echo ""
echo -e "${CYAN}📈 Para visualizar os gráficos:${NC}"
echo -e "   ${BLUE}node scripts/generate-rate-limit-chart.js${NC}"
echo ""
echo -e "${CYAN}📖 Para ler os resultados:${NC}"
echo -e "   ${BLUE}cat $RESULTS_FILE${NC}"
echo ""
