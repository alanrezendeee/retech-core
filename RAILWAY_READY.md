# ✅ PRONTO PARA RAILWAY - Checklist Final

## 🎯 Status: TUDO PREPARADO! 

Sua aplicação está 100% pronta para deploy no Railway com:
- ✅ Seeds embedados na imagem Docker
- ✅ Migrations automáticas
- ✅ MongoDB suportado
- ✅ Health checks configurados
- ✅ Dockerfile otimizado para produção
- ✅ Scripts de verificação

---

## 📦 O que foi criado/configurado

### Arquivos de Deploy

| Arquivo | Descrição | Status |
|---------|-----------|--------|
| `Dockerfile.railway` | Dockerfile produção com seeds embedados | ✅ Pronto |
| `railway.json` | Configuração Railway (JSON) | ✅ Pronto |
| `railway.toml` | Configuração Railway (TOML) | ✅ Pronto |
| `railway.env.example` | Exemplo de variáveis de ambiente | ✅ Pronto |
| `railway-check.sh` | Script de verificação pré-deploy | ✅ Pronto |
| `RAILWAY_DEPLOY.md` | Guia completo de deploy | ✅ Pronto |

### Arquivos Atualizados

| Arquivo | O que foi mudado |
|---------|------------------|
| `seeds/.gitignore` | Ajustado para permitir commit dos JSONs |
| `internal/bootstrap/migrations.go` | Busca de seeds melhorada |
| `README.md` | Seção de Railway Deploy adicionada |

### Seeds

| Arquivo | Tamanho | Status no Git |
|---------|---------|---------------|
| `seeds/estados.json` | 5.1 KB | ⚠️ Precisa add |
| `seeds/municipios.json` | 6.0 MB | ⚠️ Precisa add |

---

## 🚀 Como fazer deploy AGORA

### Opção A: Deploy Rápido (3 minutos)

```bash
# 1. Adicionar seeds ao git (necessário para Railway)
git add seeds/estados.json seeds/municipios.json

# 2. Commit
git commit -m "chore: preparar para Railway - adicionar seeds"

# 3. Push para GitHub
git push origin main

# 4. Ir para Railway
# - Acesse https://railway.app
# - Criar projeto → Deploy MongoDB
# - Adicionar serviço → GitHub repo → selecionar retech-core
# - Configurar variáveis (ver abaixo)
```

### Opção B: Verificar antes (recomendado)

```bash
# 1. Executar verificação
./railway-check.sh

# 2. Seguir instruções do script

# 3. Commit e push
git add .
git commit -m "chore: preparar para Railway"
git push origin main
```

---

## ⚙️ Variáveis de Ambiente no Railway

Configure estas variáveis no **Railway Dashboard → Variables**:

```bash
# Mínimo necessário:
PORT=8080
ENV=production
MONGO_URI=${{MongoDB.MONGO_URL}}
MONGO_DB=retech_core
CORS_ENABLE=true
```

**💡 Dica**: Use `${{MongoDB.MONGO_URL}}` para referenciar o serviço MongoDB do mesmo projeto Railway.

---

## 📋 Checklist Pré-Deploy

Marque conforme for completando:

### Git e Seeds
- [ ] Seeds existem localmente (`ls seeds/`)
- [ ] Seeds adicionados ao git (`git add seeds/*.json`)
- [ ] Mudanças commitadas
- [ ] Push para GitHub (`git push origin main`)

### Railway Setup
- [ ] Conta Railway criada
- [ ] Projeto criado
- [ ] MongoDB provisionado no Railway
- [ ] Repositório GitHub conectado
- [ ] Variáveis de ambiente configuradas

### Verificação
- [ ] Railway detectou `Dockerfile.railway`
- [ ] Build iniciou sem erros
- [ ] Deploy completou com sucesso
- [ ] Health check passou (`/health`)
- [ ] Seeds carregaram (ver logs)

### Testes
- [ ] `/health` retorna 200 OK
- [ ] `/version` retorna versão
- [ ] `/geo/ufs` retorna 27 estados
- [ ] `/geo/municipios/PE` retorna municípios

---

## 🔍 Verificação Pós-Deploy

### 1. Verificar logs

```bash
# Via Railway CLI
railway logs --service api

# Procurar por:
# ✅ "27 estados inseridos com sucesso"
# ✅ "5570 municípios inseridos com sucesso"
# ✅ "listening on :8080"
```

### 2. Testar API

```bash
# Substituir pela sua URL Railway
export API_URL="https://retech-core-production.up.railway.app"

# Health
curl $API_URL/health

# Estados
curl $API_URL/geo/ufs | jq

# Município específico
curl $API_URL/geo/municipios/id/2611606 | jq
```

### 3. Verificar MongoDB

```bash
# Via Railway CLI
railway connect MongoDB

# No MongoDB shell:
use retech_core
db.estados.countDocuments()    // Deve retornar: 27
db.municipios.countDocuments() // Deve retornar: 5570
exit
```

---

## 🎛️ Configurações Railway Recomendadas

### Networking
- ✅ Gerar domínio público
- ✅ Configurar domínio customizado (opcional)

### Deployment
- ✅ Auto-deploy habilitado (main branch)
- ✅ Health check configurado (`/health`)

### Resources
- **Starter**: Suficiente para começar
- **Hobby** ($5/mês): Recomendado para produção
- **Pro**: Para alta demanda

### MongoDB
- **Community** (free): Para desenvolvimento
- **MongoDB Atlas M0** (free): Melhor opção free
- **MongoDB Atlas M10** ($9/mês): Recomendado para produção

---

## 📊 O que acontece no primeiro deploy

1. **Build** (~2-3 minutos)
   - Compila código Go
   - Copia seeds para imagem
   - Otimiza imagem

2. **Deploy** (~1 minuto)
   - Inicia container
   - Conecta ao MongoDB
   - **Executa migrations** ⏱️ (~10-30 segundos)
   - **Carrega 27 estados** ⏱️ (~1 segundo)
   - **Carrega 5570 municípios** ⏱️ (~5-15 segundos)
   - Cria índices
   - Inicia servidor

3. **Health Check**
   - Railway verifica `/health`
   - Timeout: 300s (5 minutos)
   - Status: ✅ Healthy

**Tempo total estimado**: 3-5 minutos

---

## 🐛 Troubleshooting Rápido

### Erro: "arquivo estados.json não encontrado"

```bash
# Seeds não commitados
git add seeds/*.json
git commit -m "fix: adicionar seeds"
git push origin main
```

### Build falha

```bash
# Testar localmente
docker build -f Dockerfile.railway -t test .

# Ver logs detalhados
railway logs --service api
```

### Seeds não carregam

```bash
# Re-executar migrations
railway connect MongoDB
use retech_core
db.migrations.deleteMany({})
exit

railway restart --service api
```

### Timeout no Health Check

- Aguardar mais (seeds podem demorar)
- Verificar logs: `railway logs`
- Timeout configurado para 300s (5min)

---

## 📚 Documentação Completa

| Documento | Quando usar |
|-----------|-------------|
| **[RAILWAY_DEPLOY.md](RAILWAY_DEPLOY.md)** | Guia passo a passo completo |
| **[DOCKER_TROUBLESHOOTING.md](DOCKER_TROUBLESHOOTING.md)** | Problemas com Docker |
| **[README.md](README.md)** | Visão geral da API |

---

## 💰 Custos Estimados

### Free Tier (Starter)
- **Railway**: $5 USD crédito/mês
- **MongoDB Community**: ~$3/mês
- **API**: ~$2/mês
- ✅ **Total**: Dentro do free tier

### Produção (Recomendado)
- **Railway Hobby**: $5/mês fixo
- **MongoDB Atlas M10**: $9/mês
- ✅ **Total**: ~$14/mês

### Alta Demanda
- **Railway Pro**: A partir de $20/mês
- **MongoDB Atlas M20**: A partir de $18/mês
- ✅ **Total**: A partir de ~$38/mês

---

## 🎉 Próximos Passos

Depois do deploy:

1. ✅ Configurar domínio customizado
2. ✅ Configurar CORS para seu frontend
3. ✅ Monitorar uso e custos
4. ✅ Configurar backup do MongoDB
5. ✅ Importar collection Postman
6. ✅ Testar todos os endpoints

---

## 🆘 Suporte

Problemas?

1. ✅ Execute `./railway-check.sh`
2. ✅ Consulte [RAILWAY_DEPLOY.md](RAILWAY_DEPLOY.md)
3. ✅ Verifique logs: `railway logs`
4. ✅ Consulte [Railway Docs](https://docs.railway.app)

---

**Status**: ✅ PRONTO PARA PRODUÇÃO

**Versão**: 0.3.0

**Última verificação**: 2025-10-20

**Deploy estimado**: 3-5 minutos

---

## 🚀 COMECE AGORA!

```bash
./railway-check.sh
```

Boa sorte com o deploy! 🎉

