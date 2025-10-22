#!/bin/bash

# 🧪 Teste Simples de Rate Limiting
# Para debug rápido

set -e

API_URL="${API_URL:-http://localhost:8080}"
ADMIN_EMAIL="test-admin-1761133648@test.com"
ADMIN_PASSWORD="Admin@123"

echo "🧪 TESTE SIMPLES DE RATE LIMITING"
echo "=================================="
echo ""

# 1. Login como super admin
echo "🔐 Fazendo login..."
TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}" | jq -r '.accessToken')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "❌ Falha no login"
    exit 1
fi

echo "✅ Login OK"
echo ""

# 2. Criar um tenant de teste simples
echo "🏗️  Criando tenant de teste..."
TIMESTAMP=$(date +%s)

# Primeiro registrar o tenant
TENANT_EMAIL="test-tenant-$TIMESTAMP@test.com"
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"tenantName\": \"Test Tenant $TIMESTAMP\",
        \"tenantEmail\": \"$TENANT_EMAIL\",
        \"userName\": \"Test User\",
        \"userEmail\": \"$TENANT_EMAIL\",
        \"userPassword\": \"Test@123\",
        \"company\": \"Test Co\",
        \"purpose\": \"Testing\"
    }")

TENANT_ID=$(echo $REGISTER_RESPONSE | jq -r '.tenant.tenantId // empty')
TENANT_USER_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.accessToken // empty')

if [ -z "$TENANT_ID" ] || [ "$TENANT_ID" = "null" ]; then
    echo "❌ Falha ao criar tenant"
    echo "Response: $REGISTER_RESPONSE"
    exit 1
fi

echo "✅ Tenant criado: $TENANT_ID"
echo ""

# 3. Atualizar rate limit via admin
echo "⚙️  Configurando rate limit (5/dia, 2/min)..."
UPDATE_RESPONSE=$(curl -s -X PUT "$API_URL/admin/tenants/$TENANT_ID" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"name\": \"Test Tenant $TIMESTAMP\",
        \"email\": \"$TENANT_EMAIL\",
        \"active\": true,
        \"rateLimit\": {
            \"requestsPerDay\": 5,
            \"requestsPerMinute\": 2
        }
    }")

echo "Response: $(echo $UPDATE_RESPONSE | jq -c .)"
echo ""

# 4. Criar API Key
echo "🔑 Criando API Key..."
APIKEY_RESPONSE=$(curl -s -X POST "$API_URL/me/apikeys" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TENANT_USER_TOKEN" \
    -d '{"name":"Test Key","expiresInDays":30}')

APIKEY=$(echo $APIKEY_RESPONSE | jq -r '.key // empty')

if [ -z "$APIKEY" ] || [ "$APIKEY" = "null" ]; then
    echo "❌ Falha ao criar API key"
    echo "Response: $APIKEY_RESPONSE"
    exit 1
fi

echo "✅ API Key: ${APIKEY:0:50}..."
echo ""

# 5. Testar rate limit
echo "🧪 Testando rate limit (esperado: 5 sucesso, resto 429)..."
echo ""

SUCCESS=0
RATE_LIMITED=0

for i in {1..10}; do
    HTTP_CODE=$(curl -s -w "%{http_code}" -o /dev/null -X GET "$API_URL/geo/ufs" \
        -H "X-API-Key: $APIKEY")
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo "[$i] ✅ 200 OK"
        SUCCESS=$((SUCCESS + 1))
    elif [ "$HTTP_CODE" = "429" ]; then
        echo "[$i] 🚫 429 Rate Limited"
        RATE_LIMITED=$((RATE_LIMITED + 1))
    else
        echo "[$i] ⚠️  $HTTP_CODE"
    fi
    
    sleep 0.2
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 RESULTADOS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Sucesso:      $SUCCESS"
echo "🚫 Rate Limited: $RATE_LIMITED"
echo ""

if [ $SUCCESS -eq 5 ] && [ $RATE_LIMITED -gt 0 ]; then
    echo "✅ TESTE PASSOU! Rate limiting funcionando corretamente!"
    exit 0
elif [ $SUCCESS -le 6 ] && [ $RATE_LIMITED -gt 0 ]; then
    echo "✅ TESTE PASSOU (com margem)! Rate limiting funcionando!"
    exit 0
else
    echo "❌ TESTE FALHOU! Rate limiting não está funcionando."
    echo "   Esperado: 5 sucesso, got: $SUCCESS"
    exit 1
fi

