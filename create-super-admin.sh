#!/bin/bash

# 🚀 Script para criar super admin diretamente no banco
# Execute este script para criar o primeiro super admin

echo "🚀 Criando super admin no banco de produção..."
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

# Verificar se Railway CLI está instalado
if ! command -v railway &> /dev/null; then
    print_error "Railway CLI não encontrado!"
    print_info "Instale com: npm install -g @railway/cli"
    exit 1
fi

# Verificar se está logado
if ! railway whoami &> /dev/null; then
    print_error "Não está logado no Railway!"
    print_info "Execute: railway login"
    exit 1
fi

echo "🔧 Executando script de criação de super admin..."
echo ""

# Executar script via Railway
print_info "Criando super admin..."
railway run --service retech-core ./create-admin.sh

if [ $? -eq 0 ]; then
    print_status "Super admin criado com sucesso!"
else
    print_error "Erro ao criar super admin"
    exit 1
fi

echo ""
print_status "Super admin criado!"
echo ""
print_info "Para testar:"
print_info "1. Acesse: https://core.theretech.com.br/admin/login"
print_info "2. Use as credenciais que você criou"
print_info "3. Teste todas as funcionalidades"
echo ""
print_status "Deploy completo! 🚀"
