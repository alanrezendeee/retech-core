# 🛠️ Scripts Utilitários

## 📋 Visão Geral

Scripts auxiliares para facilitar o desenvolvimento e administração do Retech Core.

---

## 🚀 Scripts Disponíveis

### **1. `quick-admin.sh` - Criar Super Admin (RECOMENDADO)**

**Uso mais simples e rápido para criar um super admin local.**

```bash
cd /path/to/retech-core
./scripts/quick-admin.sh
```

**O que faz:**
- ✅ Remove usuário anterior (se existir)
- ✅ Cria novo tenant via API
- ✅ Cria novo usuário via API
- ✅ Altera role para `SUPER_ADMIN`
- ✅ Pronto para usar!

**Credenciais criadas:**
```
📧 Email: alanrezendeee@gmail.com
🔑 Senha: admin123456
🌐 URL: http://localhost:3001/admin/login
```

---

### **2. `create-admin-api.sh` - Criar Admin via API**

Similar ao `quick-admin.sh`, mas com mais detalhes no output.

```bash
cd /path/to/retech-core
./scripts/create-admin-api.sh
```

---

### **3. `create-admin.sh` - Criar Admin via CLI (Go)**

**Requer que o binário Go tenha o comando `create-admin` implementado.**

```bash
cd /path/to/retech-core
./scripts/create-admin.sh
```

**⚠️ Nota:** Este script só funciona se você tiver implementado o comando CLI no Go.

---

## 🧪 Pré-requisitos

### **Para todos os scripts:**

1. **Docker containers rodando:**
   ```bash
   cd /path/to/retech-core
   docker-compose -f build/docker-compose.yml up -d
   ```

2. **Verificar se estão ativos:**
   ```bash
   docker-compose -f build/docker-compose.yml ps
   ```

   Você deve ver:
   ```
   build-api-1     Up
   build-mongo-1   Up
   ```

---

## 📝 Exemplos de Uso

### **Cenário 1: Primeira vez (setup inicial)**

```bash
# 1. Subir os containers
cd /Users/alanleitederezende/source/theretech/projetos-producao/retech-core
docker-compose -f build/docker-compose.yml up -d

# 2. Aguardar ~10 segundos para tudo iniciar

# 3. Criar super admin
./scripts/quick-admin.sh

# 4. Acessar
open http://localhost:3001/admin/login
```

---

### **Cenário 2: Resetar senha do admin**

```bash
# Executar o script novamente (ele remove e recria)
./scripts/quick-admin.sh

# Senha resetada para: admin123456
```

---

### **Cenário 3: Após rebuild do backend**

```bash
# 1. Rebuild
docker-compose -f build/docker-compose.yml up --build -d

# 2. Recriar admin
./scripts/quick-admin.sh
```

---

## 🔧 Customização

### **Alterar email/senha:**

Edite o arquivo `scripts/quick-admin.sh`:

```bash
# Linha ~25-35 (no curl)
"userEmail": "SEU_EMAIL@exemplo.com",
"userPassword": "SUA_SENHA_AQUI"
```

**OU** crie um novo script personalizado:

```bash
cp scripts/quick-admin.sh scripts/my-admin.sh
# Edite my-admin.sh com seus dados
chmod +x scripts/my-admin.sh
./scripts/my-admin.sh
```

---

## ⚠️ Troubleshooting

### **Erro: "Container MongoDB não está rodando"**

**Solução:**
```bash
cd /path/to/retech-core
docker-compose -f build/docker-compose.yml up -d
```

---

### **Erro: "connection refused" ou "Cannot connect"**

**Causa:** API não está respondendo ainda

**Solução:**
```bash
# Aguarde alguns segundos e tente novamente
sleep 10
./scripts/quick-admin.sh
```

---

### **Erro: "Email já cadastrado" (409)**

**Solução:**
```bash
# O script já remove automaticamente, mas se persistir:
docker exec build-mongo-1 mongosh retech_core --eval '
db.users.deleteOne({ email: "alanrezendeee@gmail.com" });
db.tenants.deleteOne({ email: "alanrezendeee@gmail.com" });
'

# Tente novamente
./scripts/quick-admin.sh
```

---

### **Verificar se o admin foi criado:**

```bash
docker exec build-mongo-1 mongosh retech_core --eval '
db.users.findOne({ email: "alanrezendeee@gmail.com" }, { password: 0 })
'
```

**Resultado esperado:**
```json
{
  "_id": ObjectId("..."),
  "email": "alanrezendeee@gmail.com",
  "name": "Alan Rezende",
  "role": "SUPER_ADMIN",  // ✅ Importante!
  "active": true,
  ...
}
```

---

## 🎯 Fluxo Completo

```
┌─────────────────────────────────────────┐
│  1. Docker containers UP                │
│     docker-compose up -d                │
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│  2. Executar script                     │
│     ./scripts/quick-admin.sh            │
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│  3. Script faz:                         │
│     a) Remove usuário anterior          │
│     b) Cria tenant via POST /register   │
│     c) Cria user via POST /register     │
│     d) Altera role para SUPER_ADMIN     │
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│  4. Acessar login                       │
│     http://localhost:3001/admin/login   │
└─────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│  5. Logar com:                          │
│     Email: alanrezendeee@gmail.com      │
│     Senha: admin123456                  │
└─────────────────────────────────────────┘
```

---

## 📚 Scripts Relacionados

| Script | Função | Quando usar |
|--------|--------|-------------|
| `quick-admin.sh` | Criar super admin | ✅ Use sempre |
| `create-admin-api.sh` | Criar via API (verboso) | Debug |
| `create-admin-simple.sh` | Criar via MongoDB direto | Emergência |
| `create-admin.sh` | Criar via CLI Go | Se implementado |

---

## 🎨 Output Esperado

```
🔧 Criando Super Admin...

🗑️  Removendo usuário anterior (se existir)...
📝 Criando novo usuário...
✅ Usuário criado!
🔄 Alterando para SUPER_ADMIN...

✅ SUPER ADMIN CRIADO COM SUCESSO!

┌─────────────────────────────────────────┐
│  📧 Email: alanrezendeee@gmail.com      │
│  🔑 Senha: admin123456                  │
│  🌐 URL: http://localhost:3001/admin/login │
└─────────────────────────────────────────┘
```

---

## ✅ Checklist Rápido

Antes de executar o script:

- [ ] Docker instalado
- [ ] Containers rodando (`docker-compose up -d`)
- [ ] API respondendo em `http://localhost:8080`
- [ ] MongoDB acessível em `localhost:27017`
- [ ] Script tem permissão de execução (`chmod +x`)

---

**Pronto! Use `./scripts/quick-admin.sh` sempre que precisar! 🚀**

