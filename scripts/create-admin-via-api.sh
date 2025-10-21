#!/bin/bash

# 🚀 Script para criar super admin via API
# Execute este script para criar o primeiro super admin

echo "🚀 Criando super admin via API..."
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir com cor
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# URL da API
API_URL="https://api-core.theretech.com.br"

echo "🔧 Testando conexão com API..."
echo ""

# Testar health
print_info "Testando health check..."
HEALTH_RESPONSE=$(curl -s "$API_URL/health")
if [ $? -eq 0 ]; then
    print_status "API funcionando: $HEALTH_RESPONSE"
else
    print_error "API não está respondendo"
    exit 1
fi

echo ""
echo "👨‍💼 Criando tenant + super admin..."
echo ""

# Criar tenant + super admin
print_info "Criando tenant 'ThereTech' com super admin..."
TENANT_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantName": "ThereTech",
    "tenantEmail": "admin@theretech.com.br",
    "company": "ThereTech",
    "purpose": "Administração da plataforma",
    "userName": "Super Admin",
    "userEmail": "admin@theretech.com.br",
    "userPassword": "admin12345678"
  }')

if [ $? -eq 0 ]; then
    print_status "Tenant + Super admin criado: $TENANT_RESPONSE"
else
    print_error "Erro ao criar tenant + super admin"
    exit 1
fi

echo ""
echo "🧪 Testando login..."
echo ""

# Testar login
print_info "Testando login do super admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@theretech.com.br",
    "password": "admin12345678"
  }')

if [ $? -eq 0 ]; then
    print_status "Login funcionando: $LOGIN_RESPONSE"
else
    print_error "Erro no login"
    exit 1
fi

echo ""
print_status "Super admin criado com sucesso!"
echo ""
print_info "Para testar no frontend:"
print_info "1. Acesse: https://core.theretech.com.br/admin/login"
print_info "2. Email: admin@theretech.com.br"
print_info "3. Senha: admin12345678"
echo ""
print_warning "IMPORTANTE:"
print_info "O usuário criado será TENANT_USER, não SUPER_ADMIN"
print_info "Para ter acesso de admin, você precisará:"
print_info "1. Fazer login no frontend"
print_info "2. Acessar o banco diretamente"
print_info "3. Alterar o role de TENANT_USER para SUPER_ADMIN"
echo ""
print_status "Deploy completo! 🚀"
