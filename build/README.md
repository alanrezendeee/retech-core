# 🐳 Docker - Retech Core

Arquivos Docker para executar a aplicação em containers.

## 📦 Arquivos

- `Dockerfile` - Imagem Docker multi-stage para produção
- `docker-compose.yml` - Orquestração de serviços (API + MongoDB)

## 🚀 Uso Rápido

### Primeira vez / Com alterações no código:

```bash
docker-compose up --build
```

### Execuções seguintes:

```bash
docker-compose up
```

### Em background:

```bash
docker-compose up -d
```

## 📋 Pré-requisitos

⚠️ **IMPORTANTE**: Os arquivos de seed devem existir antes de iniciar!

```bash
# Da raiz do projeto:
ls ../seeds/
# Deve listar:
# - estados.json
# - municipios.json
```

Se não existirem, copie-os:

```bash
cp ~/Downloads/estados.json ../seeds/
cp ~/Downloads/municipios.json ../seeds/
```

## 🔧 Serviços

### API (retech-core)
- **Porta**: 8080
- **Imagem**: `theretech/retech-core:dev`
- **Volume**: `../seeds:/app/seeds:ro` (somente leitura)
- **Dependência**: MongoDB

### MongoDB
- **Porta**: 27017
- **Imagem**: `mongo:7`
- **Volume**: `mongo-data:/data/db` (persistente)
- **Database**: `retech_core`

## 📊 Comandos Úteis

```bash
# Ver logs
docker-compose logs -f

# Ver status
docker-compose ps

# Parar
docker-compose down

# Parar e remover volumes (perde dados!)
docker-compose down -v

# Rebuild
docker-compose build --no-cache

# Restart de um serviço
docker-compose restart api
```

## 🐛 Troubleshooting

### Erro: "arquivo estados.json não encontrado"

**Solução**: Veja [../DOCKER_QUICK_FIX.md](../DOCKER_QUICK_FIX.md)

### Erro: "port is already allocated"

```bash
# Verificar o que usa a porta
lsof -i :8080

# Ou mudar a porta no docker-compose.yml:
ports:
  - "8081:8080"
```

### Seeds não carregam

```bash
# Verificar no MongoDB
docker exec -it build-mongo-1 mongosh retech_core --eval "db.estados.countDocuments()"

# Deve retornar: 27
```

### Container reinicia constantemente

```bash
# Ver logs
docker-compose logs api

# Rebuild completo
docker-compose down -v
docker-compose up --build
```

## 📚 Mais Informações

- [Guia Completo de Troubleshooting](../DOCKER_TROUBLESHOOTING.md)
- [Script de Setup Automático](../docker-setup.sh)
- [README Principal](../README.md)

## 🏗️ Estrutura das Imagens

### Build Stage
- Base: `golang:1.22`
- Compila o binário estático
- Output: `/out/retech-core`

### Runtime Stage
- Base: `gcr.io/distroless/base-debian12`
- Workdir: `/app`
- Binário: `/app/retech-core`
- User: `nonroot:nonroot` (segurança)

## 🔒 Segurança

- ✅ Imagem distroless (sem shell, minimal attack surface)
- ✅ User não-root
- ✅ Seeds montados como read-only
- ✅ Binário estático (sem dependências)

## 🌍 Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto (opcional):

```env
PORT=8080
ENV=development
MONGO_URI=mongodb://mongo:27017
MONGO_DB=retech_core
CORS_ENABLE=true
```

O docker-compose.yml lerá automaticamente se existir.

---

**Última atualização**: 2025-10-20
**Versão**: 0.3.0

