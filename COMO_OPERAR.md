# 🎯 COMO OPERAR - Guia Direto

**Última atualização**: 2025-10-21

---

## ⚡ SETUP INICIAL (Uma vez só)

### 1. Criar Super Admin (VOCÊ)

```bash
# Executar (já está no diretório retech-core):
./create-admin.sh
```

**Quando pedir, digite:**
- Email: `admin@theretech.com.br` (ou o que quiser)
- Nome: `Alan Leite` (seu nome)
- Senha: `sua_senha_forte` (min 8 caracteres)

**Resultado:**
```
✅ Super Admin criado com sucesso!
```

---

## 🚀 EXECUTAR DIARIAMENTE

### Backend (Docker)

```bash
cd /Users/alanleitederezende/source/theretech/projetos-producao/retech-core
docker-compose -f build/docker-compose.yml up
```

**Status esperado:**
```
✅ listening on :8080 (env=development)
```

### Frontend

```bash
cd /Users/alanleitederezende/source/theretech/projetos-producao/retech-admin
yarn dev
```

**Status esperado:**
```
✅ Ready on http://localhost:3000
```

---

## 🌐 ACESSAR A PLATAFORMA

### URLs Locais

| URL | O que é | Quem acessa |
|-----|---------|-------------|
| http://localhost:3000 | Landing page | Público |
| http://localhost:3000/admin/login | Admin login | **VOCÊ** (super admin) |
| http://localhost:3000/painel/login | Dev login | Desenvolvedores |
| http://localhost:3000/painel/register | Criar conta | Novos devs |
| http://localhost:8080/docs | API docs | Todos |

### Suas Credenciais Admin

- Email: `admin@theretech.com.br` (o que você definiu)
- Senha: (a que você criou)
- Acesso: `/admin/login`

---

## 👨‍💼 COMO ADMIN (VOCÊ)

### 1. Fazer Login

- Acesse: http://localhost:3000/admin/login
- Email: `admin@theretech.com.br`
- Senha: (sua senha)

### 2. Dashboard Admin

Você verá:
- Total de tenants
- Total de API Keys
- Requests hoje/mês
- Status do sistema

### 3. Gerenciar (em desenvolvimento)

**Já funciona:**
- ✅ Ver dashboard
- ✅ Ver seu perfil

**Em breve:**
- 🚧 Listar tenants
- 🚧 Criar/editar tenants
- 🚧 Gerenciar API Keys
- 🚧 Ver analytics

---

## 👨‍💻 COMO DESENVOLVEDOR (Tenants)

### 1. Registrar

- Acesse: http://localhost:3000/painel/register
- Preencha dados da empresa
- Preencha seus dados pessoais
- Clique "Criar Conta Gratuita"

### 2. Login

- Acesse: http://localhost:3000/painel/login
- Use email e senha criados

### 3. Dashboard

Desenvolvedor vê:
- Uso atual (requests)
- Limite diário (1000)
- Plano (Free)
- Quick start guide

### 4. Criar API Key (em breve)

Por enquanto, você (admin) precisa criar via MongoDB ou esperar implementarmos a UI.

---

## 🔑 CRIAR API KEY MANUAL (Temporário)

Enquanto a UI não está pronta:

```bash
# Conectar ao MongoDB
docker exec -it build-mongo-1 mongosh retech_core

# Criar API Key (ajustar tenantId)
db.api_keys.insertOne({
  keyId: "rtc_" + Date.now(),
  keyHash: "hash_temporario",
  scopes: ["geo:read"],
  ownerId: "tenant-id-aqui",
  expiresAt: new Date(Date.now() + 365*24*60*60*1000),
  revoked: false,
  createdAt: new Date()
})

exit
```

---

## 🧪 TESTAR A API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Estados (precisa API Key)

```bash
curl http://localhost:8080/geo/ufs \
  -H "x-api-key: rtc_sua_key"
```

**Resposta esperada:**
```json
{
  "success": true,
  "code": "OK",
  "data": [ ... 27 estados ... ]
}
```

### 3. Ver Headers de Rate Limit

```bash
curl -i http://localhost:8080/geo/ufs \
  -H "x-api-key: rtc_sua_key"
```

Procure por:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1729641600
```

---

## 📊 VERIFICAR BANCO DE DADOS

### Ver collections criadas

```bash
docker exec -it build-mongo-1 mongosh retech_core

# Ver collections
show collections

# Deve ter:
# - estados (27)
# - municipios (5570)
# - migrations (2)
# - tenants (quando criar)
# - users (depois do create-admin)
# - api_keys (quando criar)
# - rate_limits (quando usar API)
# - api_usage_logs (quando usar API)
```

### Verificar dados

```bash
# Estados
db.estados.countDocuments()  # 27

# Municípios
db.municipios.countDocuments()  # 5570

# Users (depois de criar admin)
db.users.find()

# Sair
exit
```

---

## 🔄 FLUXO COMPLETO DE TESTE

### Cenário 1: Você como Admin

1. ✅ Executar `./create-admin.sh` (uma vez)
2. ✅ Acessar `/admin/login`
3. ✅ Ver dashboard admin
4. 🚧 Gerenciar tenants (implementar depois)

### Cenário 2: Desenvolvedor se Registrando

1. ✅ Acessa `/painel/register`
2. ✅ Preenche formulário
3. ✅ Conta criada automaticamente
4. ✅ Redirect para `/painel/dashboard`
5. ✅ Vê uso e limites
6. 🚧 Cria API Key (implementar depois)
7. 🚧 Testa API com a key

---

## ⚠️ PROBLEMAS COMUNS

### "Connection refused" no frontend

**Causa**: Backend não está rodando

**Solução:**
```bash
cd retech-core
docker-compose -f build/docker-compose.yml up
```

### "arquivo estados.json não encontrado"

**Causa**: Seeds não estão no lugar

**Solução:**
```bash
ls seeds/
# Deve ter: estados.json e municipios.json
```

### Collections não aparecem

**Causa**: Collections só são criadas no primeiro uso

**Solução**: Normal! Execute `./create-admin.sh` e elas aparecerão.

### Erro ao criar admin

**Causa**: MongoDB não está rodando

**Solução:**
```bash
docker ps | grep mongo
# Deve mostrar container rodando
```

---

## 📝 RESUMO RÁPIDO

```bash
# 1. Subir backend (Docker)
cd retech-core
docker-compose -f build/docker-compose.yml up

# 2. Criar admin (uma vez)
./create-admin.sh

# 3. Subir frontend (outro terminal)
cd retech-admin
yarn dev

# 4. Acessar
# http://localhost:3000/admin/login
```

---

## 🎯 PRÓXIMOS PASSOS

Após criar o admin:

1. ✅ Login no admin
2. ✅ Ver dashboard
3. ✅ Testar registro de tenant
4. 🔜 Implementar páginas de gerenciamento
5. 🔜 Deploy no Railway

---

**Alguma dúvida? Consulte [ROADMAP.md](ROADMAP.md)** 📚

