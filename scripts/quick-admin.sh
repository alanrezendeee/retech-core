#!/bin/bash

# Script RÁPIDO para criar Super Admin
# Uso: ./scripts/quick-admin.sh

set -e

echo "🔧 Criando Super Admin..."
echo ""

# Verificar se os containers estão rodando
if ! docker ps | grep -q "build-mongo-1"; then
    echo "⚠️  Container MongoDB não está rodando!"
    echo "Execute: cd /path/to/retech-core && docker-compose -f build/docker-compose.yml up -d"
    exit 1
fi

# Deletar usuário anterior (se existir)
echo "🗑️  Removendo usuário anterior (se existir)..."
docker exec build-mongo-1 mongosh retech_core --quiet --eval '
db.users.deleteOne({ email: "alanrezendeee@gmail.com" });
db.tenants.deleteOne({ email: "alanrezendeee@gmail.com" });
' > /dev/null 2>&1

# Criar via API
echo "📝 Criando novo usuário..."
RESPONSE=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenantName": "Super Admin",
    "tenantEmail": "alanrezendeee@gmail.com",
    "company": "Theretech",
    "purpose": "Administração do sistema",
    "userName": "Alan Rezende",
    "userEmail": "alanrezendeee@gmail.com",
    "userPassword": "admin123456"
  }')

# Verificar se deu certo
if echo "$RESPONSE" | grep -q "accessToken\|tenant"; then
  echo "✅ Usuário criado!"
  
  # Alterar role para SUPER_ADMIN
  echo "🔄 Alterando para SUPER_ADMIN..."
  docker exec build-mongo-1 mongosh retech_core --quiet --eval '
  var result = db.users.updateOne(
    { email: "alanrezendeee@gmail.com" },
    { $set: { role: "SUPER_ADMIN" } }
  );
  if (result.modifiedCount > 0) {
    print("OK");
  }
  ' > /dev/null 2>&1
  
  echo ""
  echo "✅ SUPER ADMIN CRIADO COM SUCESSO!"
  echo ""
  echo "┌─────────────────────────────────────────┐"
  echo "│  📧 Email: alanrezendeee@gmail.com      │"
  echo "│  🔑 Senha: admin123456                  │"
  echo "│  🌐 URL: http://localhost:3001/admin/login │"
  echo "└─────────────────────────────────────────┘"
  echo ""
else
  echo "❌ Erro ao criar usuário"
  echo "Resposta da API:"
  echo "$RESPONSE" | head -20
  echo ""
  echo "💡 Dica: Verifique se a API está rodando em http://localhost:8080"
fi

