#!/bin/bash

# 🧪 Teste Simples de CORS (sem autenticação)
# Testa comportamento básico de CORS com rotas públicas

set -e

BASE_URL="http://localhost:8080"

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🧪 TESTE SIMPLES DE CORS${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

sleep 3  # Aguardar backend iniciar

# Teste 1: Verificar CORS com rota pública (health)
echo -e "${YELLOW}[1/4] ✅ Testando rota /health com origin permitido...${NC}"
RESPONSE=$(curl -s -H "Origin: http://localhost:3000" "$BASE_URL/health")

if echo "$RESPONSE" | grep -q '"status":"ok"'; then
    echo -e "${GREEN}✅ /health responde corretamente!${NC}"
else
    echo -e "${RED}❌ /health não respondeu: $RESPONSE${NC}"
    exit 1
fi
echo ""

# Teste 2: Verificar header CORS
echo -e "${YELLOW}[2/4] 🔍 Verificando headers CORS...${NC}"
HEADERS=$(curl -s -I -H "Origin: http://localhost:3000" "$BASE_URL/health" 2>&1)

if echo "$HEADERS" | grep -q "Access-Control-Allow-Origin"; then
    ORIGIN=$(echo "$HEADERS" | grep "Access-Control-Allow-Origin" | cut -d' ' -f2 | tr -d '\r')
    echo -e "${GREEN}✅ CORS header presente: $ORIGIN${NC}"
else
    echo -e "${YELLOW}⚠️  CORS header não presente (CORS pode estar desabilitado)${NC}"
fi
echo ""

# Teste 3: Testar origin NÃO permitido
echo -e "${YELLOW}[3/4] ❌ Testando origin NÃO permitido...${NC}"
RESPONSE=$(curl -s -H "Origin: https://malicious-site.com" "$BASE_URL/health")

if echo "$RESPONSE" | grep -q "CORS"; then
    echo -e "${GREEN}✅ Origin bloqueado com mensagem de erro!${NC}"
elif echo "$RESPONSE" | grep -q '"status":"ok"'; then
    # Verificar se header CORS foi adicionado
    HEADERS=$(curl -s -I -H "Origin: https://malicious-site.com" "$BASE_URL/health" 2>&1)
    if ! echo "$HEADERS" | grep -q "Access-Control-Allow-Origin: https://malicious-site.com"; then
        echo -e "${GREEN}✅ Origin bloqueado (sem CORS header)!${NC}"
    else
        echo -e "${RED}❌ Origin malicioso foi permitido!${NC}"
        exit 1
    fi
fi
echo ""

# Teste 4: Testar sem Origin header (deve sempre funcionar)
echo -e "${YELLOW}[4/4] 🚀 Testando SEM Origin header...${NC}"
RESPONSE=$(curl -s "$BASE_URL/health")

if echo "$RESPONSE" | grep -q '"status":"ok"'; then
    echo -e "${GREEN}✅ Request sem Origin funciona (correto)!${NC}"
else
    echo -e "${RED}❌ Request sem Origin falhou: $RESPONSE${NC}"
    exit 1
fi
echo ""

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ TESTES BÁSICOS PASSARAM!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${GREEN}📊 Resumo:${NC}"
echo -e "  ✅ Rota pública responde"
echo -e "  ✅ Headers CORS verificados"
echo -e "  ✅ Origins não permitidas bloqueadas"
echo -e "  ✅ Requests sem Origin funcionam"
echo ""

echo -e "${YELLOW}💡 Para teste completo com autenticação, veja GERENCIAR_CORS_SEM_INTERFACE.md${NC}"

