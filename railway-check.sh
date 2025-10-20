#!/bin/bash

# Script de verificação pré-deploy para Railway
# Verifica se tudo está pronto para deploy

set -e

echo "🚂 Railway Deploy - Pre-flight Check"
echo "====================================="
echo ""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

ERRORS=0

# 1. Verificar estrutura
echo "📂 Verificando estrutura de arquivos..."

if [ ! -f "Dockerfile.railway" ]; then
    print_error "Dockerfile.railway não encontrado"
    ((ERRORS++))
else
    print_success "Dockerfile.railway encontrado"
fi

if [ ! -f "railway.json" ] && [ ! -f "railway.toml" ]; then
    print_error "Configuração Railway não encontrada (railway.json ou railway.toml)"
    ((ERRORS++))
else
    print_success "Configuração Railway encontrada"
fi

# 2. Verificar seeds
echo ""
echo "📦 Verificando seeds..."

if [ ! -d "seeds" ]; then
    print_error "Diretório seeds/ não existe"
    ((ERRORS++))
else
    if [ ! -f "seeds/estados.json" ]; then
        print_error "seeds/estados.json não encontrado"
        ((ERRORS++))
    else
        ESTADOS_SIZE=$(wc -c < seeds/estados.json)
        if [ "$ESTADOS_SIZE" -lt 1000 ]; then
            print_error "seeds/estados.json muito pequeno ($ESTADOS_SIZE bytes)"
            ((ERRORS++))
        else
            print_success "estados.json encontrado ($(numfmt --to=iec-i --suffix=B $ESTADOS_SIZE))"
        fi
    fi
    
    if [ ! -f "seeds/municipios.json" ]; then
        print_error "seeds/municipios.json não encontrado"
        ((ERRORS++))
    else
        MUNICIPIOS_SIZE=$(wc -c < seeds/municipios.json)
        if [ "$MUNICIPIOS_SIZE" -lt 100000 ]; then
            print_error "seeds/municipios.json muito pequeno ($MUNICIPIOS_SIZE bytes)"
            ((ERRORS++))
        else
            print_success "municipios.json encontrado ($(numfmt --to=iec-i --suffix=B $MUNICIPIOS_SIZE))"
        fi
    fi
fi

# 3. Verificar se seeds estão no git
echo ""
echo "🔍 Verificando Git..."

if git ls-files --error-unmatch seeds/estados.json >/dev/null 2>&1; then
    print_success "estados.json está no Git"
else
    print_warning "estados.json NÃO está no Git"
    echo "   Execute: git add seeds/estados.json"
fi

if git ls-files --error-unmatch seeds/municipios.json >/dev/null 2>&1; then
    print_success "municipios.json está no Git"
else
    print_warning "municipios.json NÃO está no Git"
    echo "   Execute: git add seeds/municipios.json"
fi

# 4. Verificar Go
echo ""
echo "🔧 Verificando build..."

if ! command -v go &> /dev/null; then
    print_warning "Go não instalado (ok se usar só Docker)"
else
    print_success "Go instalado: $(go version)"
    
    echo "   Testando build..."
    if go build -o /tmp/retech-core-test ./cmd/api 2>/dev/null; then
        print_success "Build Go funciona"
        rm -f /tmp/retech-core-test
    else
        print_error "Build Go falhou"
        ((ERRORS++))
    fi
fi

# 5. Verificar Dockerfile
echo ""
echo "🐳 Verificando Dockerfile..."

if command -v docker &> /dev/null; then
    print_success "Docker instalado"
    
    # Teste de build (opcional)
    echo "   Deseja testar build do Dockerfile.railway? (y/N)"
    read -p "   " -n 1 -r -t 5
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "   Buildando... (pode demorar)"
        if docker build -f Dockerfile.railway -t retech-core-railway-test . >/dev/null 2>&1; then
            print_success "Docker build passou"
            docker rmi retech-core-railway-test >/dev/null 2>&1 || true
        else
            print_error "Docker build falhou"
            echo "   Execute manualmente: docker build -f Dockerfile.railway -t test ."
            ((ERRORS++))
        fi
    else
        print_warning "Build do Docker pulado"
    fi
else
    print_warning "Docker não instalado (ok para Railway)"
fi

# 6. Verificar branch
echo ""
echo "🌿 Verificando Git branch..."

CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")
if [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ]; then
    print_success "Branch: $CURRENT_BRANCH"
else
    print_warning "Branch atual: $CURRENT_BRANCH (Railway usa 'main' por padrão)"
fi

# 7. Verificar mudanças não commitadas
if [ -n "$(git status --porcelain)" ]; then
    print_warning "Existem mudanças não commitadas"
    echo "   Execute: git status"
else
    print_success "Sem mudanças pendentes"
fi

# 8. Resumo
echo ""
echo "======================================="
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ TUDO PRONTO PARA DEPLOY!${NC}"
    echo ""
    echo "Próximos passos:"
    echo "1. Commit e push:"
    echo "   git add ."
    echo "   git commit -m 'chore: preparar para Railway'"
    echo "   git push origin main"
    echo ""
    echo "2. Configure no Railway:"
    echo "   - Criar projeto MongoDB"
    echo "   - Deploy do GitHub repo"
    echo "   - Configurar variáveis de ambiente"
    echo ""
    echo "3. Consulte: RAILWAY_DEPLOY.md"
else
    echo -e "${RED}❌ $ERRORS ERRO(S) ENCONTRADO(S)${NC}"
    echo ""
    echo "Corrija os erros acima antes de fazer deploy."
    exit 1
fi

echo "======================================="

