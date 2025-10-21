#!/bin/bash

# Script para criar Super Admin no ambiente local
# Uso: ./scripts/create-admin.sh

set -e

echo "🔧 Criando Super Admin..."
echo ""

# Executar o comando Go dentro do container Docker
docker exec -it build-api-1 /app/retech-core create-admin \
  --email="alanrezendeee@gmail.com" \
  --name="Alan Rezende" \
  --password="admin123456"

echo ""
echo "✅ Super Admin criado com sucesso!"
echo ""
echo "📧 Email: alanrezendeee@gmail.com"
echo "🔑 Senha: admin123456"
echo ""
echo "🌐 Acesse: http://localhost:3001/admin/login"
echo ""

