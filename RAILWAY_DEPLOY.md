# 🚂 Deploy no Railway - Retech Core API

Guia completo para deploy da aplicação no Railway com MongoDB e seeds automáticos.

## 📋 Pré-requisitos

- ✅ Conta no [Railway](https://railway.app)
- ✅ Repositório no GitHub com o código
- ✅ Arquivos `estados.json` e `municipios.json` no diretório `seeds/`
- ✅ Railway CLI (opcional, mas recomendado)

## 🚀 Deploy Passo a Passo

### 1. Preparar o Repositório GitHub

#### 1.1. Commit dos arquivos necessários

```bash
# Verificar que os seeds estão no repo
git status seeds/

# IMPORTANTE: Adicionar seeds ao git (são necessários na imagem Docker)
# O .gitignore em seeds/ deve ser ajustado para Railway
```

#### 1.2. Criar um .gitignore específico para seeds no Railway

Crie ou edite `seeds/.gitignore`:

```gitignore
# Manter arquivos JSON para Railway
!*.json
```

#### 1.3. Commit e push

```bash
git add .
git commit -m "chore: preparar para deploy no Railway"
git push origin main
```

### 2. Configurar MongoDB no Railway

#### 2.1. Criar novo projeto no Railway

1. Acesse [railway.app](https://railway.app)
2. Clique em **"New Project"**
3. Escolha **"Deploy MongoDB"**
4. Aguarde o MongoDB ser provisionado

#### 2.2. Anotar credenciais do MongoDB

Na aba **"Variables"** do MongoDB, você verá:
- `MONGOHOST`
- `MONGOPASSWORD`
- `MONGOPORT`
- `MONGOUSER`
- `MONGO_URL` ← **Use esta variável completa**

Exemplo:
```
mongodb://mongo:senha@mongodb.railway.internal:27017
```

### 3. Deploy da API

#### 3.1. Adicionar serviço API ao projeto

1. No mesmo projeto, clique em **"New Service"**
2. Escolha **"GitHub Repo"**
3. Conecte sua conta GitHub (se ainda não conectou)
4. Selecione o repositório **retech-core**
5. Clique em **"Deploy"**

#### 3.2. Configurar variáveis de ambiente

Na aba **"Variables"** do serviço da API, adicione:

```bash
# Servidor
PORT=8080
ENV=production

# MongoDB (usar a URL do serviço MongoDB do Railway)
MONGO_URI=${{MongoDB.MONGO_URL}}
MONGO_DB=retech_core

# CORS
CORS_ENABLE=true

# IMPORTANTE: Railway detecta PORT automaticamente, mas definimos 8080 por padrão
```

**💡 Dica**: Use a referência `${{MongoDB.MONGO_URL}}` para pegar automaticamente a URL do MongoDB.

#### 3.3. Configurar Dockerfile

O Railway detectará automaticamente o `Dockerfile.railway` através do `railway.json`.

Verifique que o serviço está usando:
- **Builder**: Dockerfile
- **Dockerfile Path**: `Dockerfile.railway`

### 4. Verificar Deploy

#### 4.1. Acompanhar logs

Na aba **"Deployments"**, clique no deployment ativo e veja os logs:

```
✅ Logs esperados:
[build] Step 1/16 : FROM golang:1.22-alpine AS build
[build] ...
[build] Successfully built abc123def456
[deploy] {"level":"info","message":"Executando migrations e seeds..."}
[deploy] {"level":"info","message":"[migration] Aplicando 001_seed_estados..."}
[deploy] {"level":"info","message":"[seed] Carregando estados de: seeds/estados.json"}
[deploy] {"level":"info","message":"[seed] 27 estados inseridos com sucesso"}
[deploy] {"level":"info","message":"[migration] Aplicando 002_seed_municipios..."}
[deploy] {"level":"info","message":"[seed] 5570 municípios inseridos com sucesso"}
[deploy] {"level":"info","message":"listening on :8080 (env=production)"}
```

#### 4.2. Obter URL pública

1. Na aba **"Settings"** do serviço API
2. Seção **"Networking"**
3. Clique em **"Generate Domain"**
4. Anote a URL (ex: `https://retech-core-production.up.railway.app`)

#### 4.3. Testar a API

```bash
# Substituir pela sua URL do Railway
export API_URL="https://retech-core-production.up.railway.app"

# Health check
curl $API_URL/health

# Versão
curl $API_URL/version

# Listar estados
curl $API_URL/geo/ufs | jq

# Buscar estado
curl $API_URL/geo/ufs/PE | jq

# Municípios
curl $API_URL/geo/municipios/PE | jq
```

## 🔐 Segurança e Produção

### Variáveis de Ambiente Recomendadas

```bash
# Obrigatórias
PORT=8080
ENV=production
MONGO_URI=${{MongoDB.MONGO_URL}}
MONGO_DB=retech_core

# Segurança (adicionar conforme implementação futura)
# API_KEY_REQUIRED=true
# CORS_ORIGINS=https://seu-frontend.com,https://app.seu-frontend.com

# Observabilidade (futuro)
# OTEL_ENABLED=true
# SENTRY_DSN=...
```

### CORS para Produção

Se você tem um frontend específico, configure o CORS:

```bash
CORS_ORIGINS=https://seu-app.com,https://app.seu-app.com
```

## 📊 Monitoramento

### Métricas no Railway

O Railway fornece automaticamente:
- ✅ CPU usage
- ✅ Memory usage
- ✅ Network traffic
- ✅ Deployment history

### Logs

```bash
# Via Railway CLI (recomendado)
railway logs

# Filtrar por serviço
railway logs --service api

# Tail em tempo real
railway logs --follow
```

### Health Checks

O Railway verificará automaticamente:
- **Path**: `/health`
- **Timeout**: 300s (5 minutos - tempo suficiente para seeds)
- **Restart Policy**: ON_FAILURE (até 10 tentativas)

## 🔄 Atualizações e Re-deploys

### Deploy Automático (CI/CD)

O Railway faz deploy automático em cada push para `main`:

```bash
git add .
git commit -m "feat: adicionar novo endpoint"
git push origin main

# Railway detecta e faz deploy automaticamente
```

### Deploy Manual

Via Railway CLI:

```bash
# Login
railway login

# Linkar ao projeto
railway link

# Deploy
railway up
```

### Rollback

1. Vá para **"Deployments"**
2. Encontre o deployment anterior estável
3. Clique nos 3 pontos → **"Redeploy"**

## 🗄️ Banco de Dados

### Backup do MongoDB (Railway)

O Railway não faz backup automático do MongoDB Community (free tier).

**Recomendações:**

#### Opção 1: Upgrade para MongoDB Pro (Railway)
- Backups automáticos diários
- Point-in-time recovery
- Alta disponibilidade

#### Opção 2: MongoDB Atlas
1. Criar cluster no [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Obter connection string
3. Atualizar `MONGO_URI` no Railway

```bash
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/retech_core?retryWrites=true&w=majority
```

#### Opção 3: Backup Manual via Script

```bash
# Conectar via Railway CLI
railway connect MongoDB

# Ou via mongodump
mongodump --uri="$MONGO_URI" --out=backup-$(date +%Y%m%d)
```

### Acesso ao MongoDB

```bash
# Via Railway CLI
railway connect MongoDB

# Verificar dados
use retech_core
db.estados.countDocuments()    // 27
db.municipios.countDocuments() // 5570
```

## 🔧 Troubleshooting

### Erro: "arquivo estados.json não encontrado"

**Causa**: Seeds não estão commitados no Git

**Solução**:
```bash
# Verificar
ls seeds/

# Adicionar ao git
git add seeds/*.json
git commit -m "chore: adicionar seeds para Railway"
git push origin main
```

### Erro: Build falha no Dockerfile

**Causa**: Dockerfile não encontrado

**Solução**:
1. Verificar que `Dockerfile.railway` existe na raiz
2. Verificar que `railway.json` aponta para ele
3. Re-deploy

### Seeds não carregam

**Diagnóstico**:
```bash
railway logs --service api | grep seed

# Deve mostrar:
# "27 estados inseridos"
# "5570 municípios inseridos"
```

**Solução**: Re-run migrations
```bash
# Conectar ao MongoDB
railway connect MongoDB

# Limpar
use retech_core
db.migrations.deleteMany({})
db.estados.deleteMany({})
db.municipios.deleteMany({})
exit

# Restart API
railway restart --service api
```

### Timeout no Health Check

**Causa**: Seeds demoram para carregar

**Solução**: O timeout está configurado para 300s (5min) no `railway.json`. Se ainda assim timeout:

1. Verificar logs de erro
2. Considerar índices pré-criados
3. Otimizar inserção em lotes maiores

## 💰 Custos

### Free Tier (Railway)

- **Horas gratuitas**: $5 USD/mês de crédito
- **MongoDB Community**: ~$3-5/mês
- **API**: ~$2-3/mês
- **Total estimado**: Dentro do free tier inicialmente

### Produção (Recomendado)

- **Hobby Plan**: $5-10/mês
- **MongoDB Atlas M0** (free): 512 MB
- **MongoDB Atlas M10**: $9/mês (recomendado para produção)

## 📚 Estrutura de Arquivos para Railway

```
retech-core/
├── Dockerfile.railway        ← Dockerfile de produção
├── railway.json              ← Configuração Railway (JSON)
├── railway.toml              ← Configuração Railway (TOML)
├── seeds/
│   ├── estados.json          ← DEVE estar no git
│   └── municipios.json       ← DEVE estar no git
├── internal/
└── cmd/api/main.go
```

## ✅ Checklist de Deploy

Antes de fazer deploy:

- [ ] Seeds commitados no git (`seeds/*.json`)
- [ ] `Dockerfile.railway` existe
- [ ] `railway.json` ou `railway.toml` existe
- [ ] MongoDB provisionado no Railway
- [ ] Variáveis de ambiente configuradas
- [ ] Domínio gerado
- [ ] Health check funcionando
- [ ] Estados carregados (27)
- [ ] Municípios carregados (5570)
- [ ] Testes de API passando

## 🔗 Links Úteis

- [Railway Dashboard](https://railway.app/dashboard)
- [Railway Docs](https://docs.railway.app)
- [Railway CLI](https://docs.railway.app/develop/cli)
- [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)

## 🆘 Suporte

Se encontrar problemas:
1. Verificar logs: `railway logs --service api`
2. Verificar health: `curl https://sua-url.railway.app/health`
3. Consultar [DOCKER_TROUBLESHOOTING.md](DOCKER_TROUBLESHOOTING.md)

---

**Última atualização**: 2025-10-20
**Versão**: 0.3.0
**Status**: ✅ Pronto para Produção

