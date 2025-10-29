#!/bin/bash

echo "📦 Criando usuário admin local..."

# Conectar no MongoDB e criar usuário
docker exec build-mongo-1 mongosh retech_core --eval '
db.users.deleteMany({ email: "alanrezendeee@gmail.com" });
db.tenants.deleteMany({ email: "admin@theretech.com.br" });

// Criar tenant SUPER_ADMIN
db.tenants.insertOne({
  "_id": "tenant-super-admin",
  "tenantId": "tenant-super-admin",
  "name": "Super Admin",
  "email": "admin@theretech.com.br",
  "company": "Retech Core",
  "purpose": "Administração do sistema",
  "active": true,
  "createdAt": new Date(),
  "updatedAt": new Date()
});

// Criar usuário admin
// Hash bcrypt de "admin123456" (gerado com bcrypt.DefaultCost)
db.users.insertOne({
  "email": "alanrezendeee@gmail.com",
  "name": "Alan Leite",
  "password": "$2a$10$b3GJK12gSehRihTTUkEWNulhT4UgUKNfFMHDNYg6HkOSf..uoz.Ra",
  "role": "SUPER_ADMIN",
  "tenantId": "tenant-super-admin",
  "active": true,
  "createdAt": new Date(),
  "updatedAt": new Date()
});

print("✅ Usuário admin criado!");
print("📧 Email: alanrezendeee@gmail.com");
print("🔑 Senha: admin123456");
'
