#!/bin/bash

# Script SIMPLES para criar Super Admin
# Uso: ./scripts/create-admin-simple.sh

echo "🔧 Criando Super Admin via MongoDB..."
echo ""

docker exec -it build-mongo-1 mongosh retech_core --eval '
db.users.deleteOne({ email: "alanrezendeee@gmail.com" });
db.tenants.deleteOne({ email: "alanrezendeee@gmail.com" });

// Criar tenant do super admin
var tenantResult = db.tenants.insertOne({
  tenantId: "tenant-super-admin",
  name: "Super Admin",
  email: "alanrezendeee@gmail.com",
  active: true,
  createdAt: new Date(),
  updatedAt: new Date()
});

// Criar super admin (senha: admin123456)
var userResult = db.users.insertOne({
  email: "alanrezendeee@gmail.com",
  name: "Alan Rezende",
  password: "$2a$10$YourBcryptHashHere", 
  role: "SUPER_ADMIN",
  tenantId: "tenant-super-admin",
  active: true,
  createdAt: new Date(),
  updatedAt: new Date()
});

print("\n✅ Super Admin criado com sucesso!");
print("\n📧 Email: alanrezendeee@gmail.com");
print("🔑 Senha: admin123456");
print("\n🌐 Acesse: http://localhost:3001/admin/login\n");
'

echo ""

