#!/bin/bash

# Script para criar primeiro super admin
# Usa o MongoDB do Docker

echo "🔐 Criar Super Admin - Retech Core"
echo "==================================="
echo ""

# Configurar variáveis para o MongoDB do Docker
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DB="retech_core"

echo "📡 Conectando ao MongoDB (Docker)..."
echo "   URI: $MONGO_URI"
echo "   DB: $MONGO_DB"
echo ""

# Executar script
go run scripts/create-admin.go

