#!/bin/bash

# 🧪 Script de Teste Completo de CORS
# Testa todos os cenários: salvar, desabilitar, habilitar, origins permitidas/negadas

set -e

BASE_URL="http://localhost:8080"
ADMIN_EMAIL="admin@theretech.com.br"
ADMIN_PASSWORD="Admin@123"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🧪 TESTE COMPLETO DE CORS${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Função para fazer login e obter token
get_auth_token() {
    echo -e "${YELLOW}[1/9] 🔐 Fazendo login...${NC}"
    
    TOKEN=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}" \
        | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)
    
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}❌ Erro ao fazer login. Verifique as credenciais.${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Login realizado com sucesso!${NC}"
    echo ""
}

# Teste 1: Salvar configuração com CORS habilitado
test_save_cors_enabled() {
    echo -e "${YELLOW}[2/9] 💾 Salvando configuração com CORS.Enabled=true...${NC}"
    
    RESPONSE=$(curl -s -X PUT "$BASE_URL/admin/settings" \
        -H "Origin: http://localhost:3000" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "defaultRateLimit": {
                "requestsPerDay": 1000,
                "requestsPerMinute": 60
            },
            "cors": {
                "enabled": true,
                "allowedOrigins": [
                    "https://core.theretech.com.br",
                    "http://localhost:3000",
                    "http://localhost:3001"
                ]
            },
            "jwt": {
                "accessTokenTTL": 900,
                "refreshTokenTTL": 604800
            },
            "api": {
                "version": "1.0.0",
                "environment": "development",
                "maintenance": false
            },
            "contact": {
                "whatsapp": "48999616679",
                "email": "suporte@theretech.com.br",
                "phone": "+55 48 99961-6679"
            },
            "cache": {
                "enabled": true,
                "cepTtlDays": 7,
                "cnpjTtlDays": 30,
                "maxSizeMb": 100,
                "autoCleanup": true
            },
            "playground": {
                "enabled": true,
                "apiKey": "rtc_demo_playground_2024",
                "rateLimit": {
                    "requestsPerDay": 100,
                    "requestsPerMinute": 10
                },
                "allowedApis": ["cep", "cnpj", "geo"]
            }
        }')
    
    if echo "$RESPONSE" | grep -q "sucesso"; then
        echo -e "${GREEN}✅ Configuração salva com sucesso!${NC}"
    else
        echo -e "${RED}❌ Erro ao salvar: $RESPONSE${NC}"
        exit 1
    fi
    
    sleep 1
    echo ""
}

# Teste 2: Verificar CORS com origin permitido
test_cors_allowed_origin() {
    echo -e "${YELLOW}[3/9] ✅ Testando origin PERMITIDO (localhost:3000)...${NC}"
    
    HEADERS=$(curl -s -I -H "Origin: http://localhost:3000" "$BASE_URL/health" 2>&1)
    
    if echo "$HEADERS" | grep -q "Access-Control-Allow-Origin: http://localhost:3000"; then
        echo -e "${GREEN}✅ CORS permitido para localhost:3000!${NC}"
    else
        echo -e "${RED}❌ CORS NÃO foi adicionado para localhost:3000${NC}"
        echo "$HEADERS"
        exit 1
    fi
    
    echo ""
}

# Teste 3: Verificar CORS com origin NÃO permitido
test_cors_denied_origin() {
    echo -e "${YELLOW}[4/9] ❌ Testando origin NÃO PERMITIDO (malicious-site.com)...${NC}"
    
    RESPONSE=$(curl -s -H "Origin: https://malicious-site.com" "$BASE_URL/admin/settings")
    
    if echo "$RESPONSE" | grep -q "Origin.*não está na lista"; then
        echo -e "${GREEN}✅ Origin bloqueado corretamente! Mensagem de erro presente.${NC}"
    else
        # Verificar se não tem header CORS (que também é correto)
        HEADERS=$(curl -s -I -H "Origin: https://malicious-site.com" "$BASE_URL/health" 2>&1)
        if ! echo "$HEADERS" | grep -q "Access-Control-Allow-Origin: https://malicious-site.com"; then
            echo -e "${GREEN}✅ Origin bloqueado corretamente! (sem CORS header)${NC}"
        else
            echo -e "${RED}❌ Origin NÃO foi bloqueado!${NC}"
            echo "$RESPONSE"
            exit 1
        fi
    fi
    
    echo ""
}

# Teste 4: Desabilitar CORS
test_disable_cors() {
    echo -e "${YELLOW}[5/9] 🔒 Desabilitando CORS...${NC}"
    
    RESPONSE=$(curl -s -X PUT "$BASE_URL/admin/settings" \
        -H "Origin: http://localhost:3000" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "defaultRateLimit": {
                "requestsPerDay": 1000,
                "requestsPerMinute": 60
            },
            "cors": {
                "enabled": false,
                "allowedOrigins": [
                    "https://core.theretech.com.br",
                    "http://localhost:3000"
                ]
            },
            "jwt": {
                "accessTokenTTL": 900,
                "refreshTokenTTL": 604800
            },
            "api": {
                "version": "1.0.0",
                "environment": "development",
                "maintenance": false
            },
            "contact": {
                "whatsapp": "48999616679",
                "email": "suporte@theretech.com.br",
                "phone": "+55 48 99961-6679"
            },
            "cache": {
                "enabled": true,
                "cepTtlDays": 7,
                "cnpjTtlDays": 30,
                "maxSizeMb": 100,
                "autoCleanup": true
            },
            "playground": {
                "enabled": true,
                "apiKey": "rtc_demo_playground_2024",
                "rateLimit": {
                    "requestsPerDay": 100,
                    "requestsPerMinute": 10
                },
                "allowedApis": ["cep", "cnpj", "geo"]
            }
        }')
    
    if echo "$RESPONSE" | grep -q "sucesso"; then
        echo -e "${GREEN}✅ CORS desabilitado com sucesso!${NC}"
    else
        echo -e "${RED}❌ Erro ao desabilitar: $RESPONSE${NC}"
        exit 1
    fi
    
    sleep 1
    echo ""
}

# Teste 5: Verificar que CORS bloqueado após desabilitar
test_cors_blocked_after_disable() {
    echo -e "${YELLOW}[6/9] 🚫 Testando se CORS está bloqueado após desabilitar...${NC}"
    
    RESPONSE=$(curl -s -H "Origin: http://localhost:3001" "$BASE_URL/health")
    
    if echo "$RESPONSE" | grep -q "CORS está desabilitado"; then
        echo -e "${GREEN}✅ CORS bloqueado corretamente! Mensagem de erro presente.${NC}"
    else
        echo -e "${RED}❌ CORS deveria estar bloqueado mas não está!${NC}"
        echo "$RESPONSE"
        exit 1
    fi
    
    echo ""
}

# Teste 6: Verificar que /admin/settings AINDA funciona de localhost (exceção)
test_admin_settings_exception() {
    echo -e "${YELLOW}[7/9] ⚙️ Testando exceção: /admin/settings de localhost...${NC}"
    
    RESPONSE=$(curl -s -H "Origin: http://localhost:3000" -H "Authorization: Bearer $TOKEN" "$BASE_URL/admin/settings")
    
    if echo "$RESPONSE" | grep -q "cors.*enabled.*false" || echo "$RESPONSE" | grep -q "data"; then
        echo -e "${GREEN}✅ /admin/settings funciona de localhost (exceção ativa)!${NC}"
    else
        echo -e "${RED}❌ /admin/settings deveria funcionar de localhost!${NC}"
        echo "$RESPONSE"
        exit 1
    fi
    
    echo ""
}

# Teste 7: Re-habilitar CORS
test_reenable_cors() {
    echo -e "${YELLOW}[8/9] ✅ Re-habilitando CORS...${NC}"
    
    RESPONSE=$(curl -s -X PUT "$BASE_URL/admin/settings" \
        -H "Origin: http://localhost:3000" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "defaultRateLimit": {
                "requestsPerDay": 1000,
                "requestsPerMinute": 60
            },
            "cors": {
                "enabled": true,
                "allowedOrigins": [
                    "https://core.theretech.com.br",
                    "http://localhost:3000",
                    "http://localhost:3001",
                    "http://localhost:3002"
                ]
            },
            "jwt": {
                "accessTokenTTL": 900,
                "refreshTokenTTL": 604800
            },
            "api": {
                "version": "1.0.0",
                "environment": "development",
                "maintenance": false
            },
            "contact": {
                "whatsapp": "48999616679",
                "email": "suporte@theretech.com.br",
                "phone": "+55 48 99961-6679"
            },
            "cache": {
                "enabled": true,
                "cepTtlDays": 7,
                "cnpjTtlDays": 30,
                "maxSizeMb": 100,
                "autoCleanup": true
            },
            "playground": {
                "enabled": true,
                "apiKey": "rtc_demo_playground_2024",
                "rateLimit": {
                    "requestsPerDay": 100,
                    "requestsPerMinute": 10
                },
                "allowedApis": ["cep", "cnpj", "geo"]
            }
        }')
    
    if echo "$RESPONSE" | grep -q "sucesso"; then
        echo -e "${GREEN}✅ CORS re-habilitado com sucesso!${NC}"
    else
        echo -e "${RED}❌ Erro ao re-habilitar: $RESPONSE${NC}"
        exit 1
    fi
    
    sleep 1
    echo ""
}

# Teste 8: Verificar lista de origins permitidas
test_allowed_origins_list() {
    echo -e "${YELLOW}[9/9] 📋 Testando lista de origins permitidas...${NC}"
    
    # Testar origin na lista
    HEADERS1=$(curl -s -I -H "Origin: http://localhost:3002" "$BASE_URL/health" 2>&1)
    
    if echo "$HEADERS1" | grep -q "Access-Control-Allow-Origin: http://localhost:3002"; then
        echo -e "${GREEN}✅ localhost:3002 (na lista) permitido!${NC}"
    else
        echo -e "${RED}❌ localhost:3002 deveria ser permitido!${NC}"
        exit 1
    fi
    
    # Testar origin NÃO na lista
    HEADERS2=$(curl -s -I -H "Origin: http://localhost:3003" "$BASE_URL/health" 2>&1)
    
    if ! echo "$HEADERS2" | grep -q "Access-Control-Allow-Origin: http://localhost:3003"; then
        echo -e "${GREEN}✅ localhost:3003 (NÃO na lista) bloqueado!${NC}"
    else
        echo -e "${RED}❌ localhost:3003 NÃO deveria ser permitido!${NC}"
        exit 1
    fi
    
    echo ""
}

# Executar todos os testes
main() {
    sleep 3  # Aguardar backend iniciar
    
    get_auth_token
    test_save_cors_enabled
    test_cors_allowed_origin
    test_cors_denied_origin
    test_disable_cors
    test_cors_blocked_after_disable
    test_admin_settings_exception
    test_reenable_cors
    test_allowed_origins_list
    
    echo -e "${BLUE}========================================${NC}"
    echo -e "${GREEN}✅ TODOS OS TESTES PASSARAM!${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
    echo -e "${GREEN}📊 Resumo:${NC}"
    echo -e "  ✅ Salvar configuração: OK"
    echo -e "  ✅ CORS habilitado + origins permitidas: OK"
    echo -e "  ✅ Origins não permitidas bloqueadas: OK"
    echo -e "  ✅ Desabilitar CORS: OK"
    echo -e "  ✅ Requests bloqueadas após desabilitar: OK"
    echo -e "  ✅ /admin/settings funciona de localhost: OK"
    echo -e "  ✅ Re-habilitar CORS: OK"
    echo -e "  ✅ Lista de origins respeitada: OK"
    echo ""
}

main

