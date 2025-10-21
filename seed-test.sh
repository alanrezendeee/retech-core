#!/bin/bash

# Script para criar dados de teste completos

export MONGO_URI="mongodb://localhost:27017"
export MONGO_DB="retech_core"

echo "🌱 Criando dados de teste completos..."
echo ""

go run scripts/seed-test-data.go

