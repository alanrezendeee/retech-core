#!/bin/bash

# Script para criar Super Admin via API
# Uso: ./scripts/create-admin-api.sh

echo "🔧 Criando Super Admin via API..."
echo ""

# Primeiro, deletar se já existir (via MongoDB)
docker exec build-mongo-1 mongosh retech_core --quiet --eval '
db.users.deleteOne({ email: "alanrezendeee@gmail.com" });
db.tenants.deleteOne({ email: "alanrezendeee@gmail.com" });
print("✅ Usuário anterior removido (se existia)");
'

echo ""

# Criar via endpoint de registro
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

# Verificar se funcionou
if echo "$RESPONSE" | grep -q "accessToken"; then
  echo "✅ Super Admin criado com sucesso!"
  echo ""
  echo "📧 Email: alanrezendeee@gmail.com"
  echo "🔑 Senha: admin123456"
  echo ""
  
  # Agora precisamos alterar o role para SUPER_ADMIN
  echo "🔄 Alterando role para SUPER_ADMIN..."
  docker exec build-mongo-1 mongosh retech_core --quiet --eval '
  db.users.updateOne(
    { email: "alanrezendeee@gmail.com" },
    { $set: { role: "SUPER_ADMIN" } }
  );
  print("✅ Role alterado para SUPER_ADMIN");
  '
  
  echo ""
  echo "🌐 Acesse: http://localhost:3001/admin/login"
  echo ""
else
  echo "❌ Erro ao criar Super Admin:"
  echo "$RESPONSE"
  echo ""
fi
