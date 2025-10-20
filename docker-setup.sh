#!/bin/bash

# Script de setup para Docker - Retech Core
# Verifica pré-requisitos e inicia a aplicação

set -e

echo "🚀 Retech Core - Docker Setup"
echo "================================"
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para print colorido
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# 1. Verificar se estamos no diretório correto
echo "📂 Verificando diretório..."
if [ ! -f "go.mod" ]; then
    print_error "Execute este script da raiz do projeto retech-core"
    exit 1
fi
print_success "Diretório correto"

# 2. Verificar Docker
echo ""
echo "🐳 Verificando Docker..."
if ! command -v docker &> /dev/null; then
    print_error "Docker não está instalado"
    echo "   Instale em: https://docs.docker.com/get-docker/"
    exit 1
fi
print_success "Docker instalado: $(docker --version)"

# 3. Verificar Docker Compose
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose não está instalado"
    exit 1
fi
print_success "Docker Compose instalado: $(docker-compose --version)"

# 4. Verificar se Docker está rodando
if ! docker info &> /dev/null; then
    print_error "Docker não está rodando"
    echo "   Inicie o Docker Desktop ou o daemon do Docker"
    exit 1
fi
print_success "Docker está rodando"

# 5. Verificar arquivos de seed
echo ""
echo "📦 Verificando arquivos de seed..."

if [ ! -d "seeds" ]; then
    print_warning "Diretório seeds/ não existe. Criando..."
    mkdir -p seeds
fi

SEEDS_OK=true

if [ ! -f "seeds/estados.json" ]; then
    print_error "seeds/estados.json não encontrado"
    
    # Tentar copiar de Downloads
    if [ -f "$HOME/Downloads/estados.json" ]; then
        print_warning "Encontrado em ~/Downloads. Copiando..."
        cp "$HOME/Downloads/estados.json" seeds/
        print_success "estados.json copiado"
    else
        SEEDS_OK=false
    fi
else
    print_success "estados.json encontrado"
fi

if [ ! -f "seeds/municipios.json" ]; then
    print_error "seeds/municipios.json não encontrado"
    
    # Tentar copiar de Downloads
    if [ -f "$HOME/Downloads/municipios.json" ]; then
        print_warning "Encontrado em ~/Downloads. Copiando..."
        cp "$HOME/Downloads/municipios.json" seeds/
        print_success "municipios.json copiado"
    else
        SEEDS_OK=false
    fi
else
    print_success "municipios.json encontrado"
fi

if [ "$SEEDS_OK" = false ]; then
    echo ""
    print_error "Arquivos de seed não encontrados!"
    echo ""
    echo "   Coloque os arquivos JSON no diretório seeds/:"
    echo "   - seeds/estados.json"
    echo "   - seeds/municipios.json"
    echo ""
    echo "   Você pode baixá-los de:"
    echo "   - https://servicodados.ibge.gov.br/api/v1/localidades/estados"
    echo "   - https://servicodados.ibge.gov.br/api/v1/localidades/municipios"
    exit 1
fi

# 6. Verificar portas
echo ""
echo "🔌 Verificando portas..."

check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        return 1
    else
        return 0
    fi
}

if ! check_port 8080; then
    print_warning "Porta 8080 já está em uso"
    echo "   Você pode parar o serviço ou mudar a porta no docker-compose.yml"
    read -p "   Continuar mesmo assim? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    print_success "Porta 8080 disponível"
fi

if ! check_port 27017; then
    print_warning "Porta 27017 já está em uso (MongoDB)"
    read -p "   Continuar mesmo assim? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    print_success "Porta 27017 disponível"
fi

# 7. Verificar .env (opcional)
echo ""
echo "⚙️  Verificando configuração..."
if [ ! -f ".env" ]; then
    print_warning "Arquivo .env não encontrado (opcional)"
    if [ -f "env.example" ]; then
        echo "   Deseja criar .env baseado em env.example? (y/N)"
        read -p "   " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            cp env.example .env
            print_success ".env criado"
        fi
    fi
else
    print_success "Arquivo .env encontrado"
fi

# 8. Perguntar sobre rebuild
echo ""
echo "🔨 Build da aplicação..."
read -p "   Fazer rebuild das imagens? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    BUILD_FLAG="--build"
else
    BUILD_FLAG=""
fi

# 9. Iniciar containers
echo ""
echo "🎬 Iniciando containers..."
cd build

# Parar containers anteriores se existirem
if docker-compose ps -q | grep -q .; then
    print_warning "Parando containers anteriores..."
    docker-compose down
fi

# Subir containers
echo ""
echo "   Executando: docker-compose up $BUILD_FLAG"
echo "   (Pressione Ctrl+C para parar)"
echo ""

docker-compose up $BUILD_FLAG

# Este ponto só será alcançado se o usuário parar com Ctrl+C
echo ""
echo "🛑 Parando containers..."
docker-compose down

echo ""
print_success "Setup concluído!"

