# 🚀 Quick Start - retech-core

Guia rápido para configurar e executar o projeto.

## Pré-requisitos

- Go 1.22+
- MongoDB 6.0+
- Docker e Docker Compose (opcional)

## Setup Rápido com Docker

### 1. Clone o repositório

```bash
cd /path/to/retech-core
```

### 2. Prepare os arquivos de seed

Coloque os arquivos `estados.json` e `municipios.json` no diretório `seeds/`:

```bash
# Se os arquivos estão em Downloads
cp ~/Downloads/estados.json seeds/
cp ~/Downloads/municipios.json seeds/
```

### 3. Inicie com Docker Compose

```bash
cd build
docker-compose up -d
```

Isso irá:
- Iniciar MongoDB na porta 27017
- Buildar e iniciar a API na porta 8080
- Executar automaticamente as migrations/seeds

### 4. Verifique o status

```bash
# Health check
curl http://localhost:8080/health

# Versão
curl http://localhost:8080/version

# Estados
curl http://localhost:8080/geo/ufs | jq
```

## Setup Local (sem Docker)

### 1. Instale MongoDB

```bash
# macOS (Homebrew)
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community

# Linux (Ubuntu)
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt-get update
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
```

### 2. Configure as variáveis de ambiente

```bash
# Crie um arquivo .env na raiz do projeto
cat > .env << EOF
PORT=8080
ENV=development
MONGO_URI=mongodb://localhost:27017
MONGO_DB=retech_core
CORS_ENABLE=true
EOF
```

### 3. Prepare os seeds

```bash
# Coloque os arquivos JSON no diretório seeds
cp ~/Downloads/estados.json seeds/
cp ~/Downloads/municipios.json seeds/
```

### 4. Instale as dependências

```bash
go mod download
```

### 5. Execute a aplicação

```bash
go run cmd/api/main.go
```

Na primeira execução, você verá:

```
[config] ENV=development PORT=8080 MONGO_URI=mongodb://localhost:27017 DB=retech_core
{"level":"info","time":"2025-10-20T...","message":"Executando migrations e seeds..."}
{"level":"info","time":"2025-10-20T...","message":"[migration] Aplicando 001_seed_estados: Popular estados brasileiros"}
{"level":"info","time":"2025-10-20T...","message":"[seed] Carregando estados de: seeds/estados.json"}
{"level":"info","time":"2025-10-20T...","message":"[seed] 27 estados inseridos com sucesso"}
{"level":"info","time":"2025-10-20T...","message":"[migration] 001_seed_estados aplicada com sucesso em 123ms"}
{"level":"info","time":"2025-10-20T...","message":"[migration] Aplicando 002_seed_municipios: Popular municípios brasileiros"}
{"level":"info","time":"2025-10-20T...","message":"[seed] Carregando municípios de: seeds/municipios.json"}
{"level":"info","time":"2025-10-20T...","message":"[seed] Inserindo 5570 municípios (isso pode demorar)..."}
{"level":"info","time":"2025-10-20T...","message":"[seed] 5570 municípios inseridos com sucesso"}
{"level":"info","time":"2025-10-20T...","message":"[migration] 002_seed_municipios aplicada com sucesso em 2.3s"}
{"level":"info","time":"2025-10-20T...","message":"Migrations concluídas com sucesso"}
{"level":"info","time":"2025-10-20T...","message":"Criando índices..."}
{"level":"info","time":"2025-10-20T...","message":"Índices criados com sucesso"}
{"level":"info","time":"2025-10-20T...","message":"listening on :8080 (env=development)"}
```

## Testando os Endpoints

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Resposta esperada:
```json
{
  "status": "ok",
  "timestamp": "2025-10-20T...",
  "database": "connected"
}
```

### 2. Listar Estados

```bash
curl http://localhost:8080/geo/ufs | jq
```

### 3. Buscar Estado por Sigla

```bash
curl http://localhost:8080/geo/ufs/PE | jq
```

### 4. Listar Municípios por Estado

```bash
curl http://localhost:8080/geo/municipios/PE | jq
```

### 5. Buscar Municípios por Nome

```bash
curl http://localhost:8080/geo/municipios?uf=PE&q=recife | jq
```

### 6. Buscar Município por Código IBGE

```bash
curl http://localhost:8080/geo/municipios/id/2611606 | jq
```

## Criando Tenants e API Keys

### 1. Criar Tenant

```bash
curl -X POST http://localhost:8080/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "empresa-exemplo",
    "name": "Empresa Exemplo LTDA",
    "email": "contato@exemplo.com",
    "active": true
  }' | jq
```

### 2. Criar API Key

```bash
curl -X POST http://localhost:8080/apikeys \
  -H "Content-Type: application/json" \
  -d '{
    "tenantId": "empresa-exemplo",
    "name": "Chave de Produção"
  }' | jq
```

Resposta (guarde o `key` retornado):
```json
{
  "id": "...",
  "key": "rtc_1234567890abcdef...",
  "tenantId": "empresa-exemplo",
  "name": "Chave de Produção",
  "active": true,
  "createdAt": "2025-10-20T..."
}
```

### 3. Usar API Key

```bash
curl http://localhost:8080/protected/apikey \
  -H "x-api-key: rtc_1234567890abcdef..."
```

## Troubleshooting

### Erro: "arquivo estados.json não encontrado"

Verifique se os arquivos estão no diretório correto:

```bash
ls -la seeds/
# Deve listar: estados.json, municipios.json
```

Se os arquivos não estiverem lá, copie-os:

```bash
cp ~/Downloads/estados.json seeds/
cp ~/Downloads/municipios.json seeds/
```

### Erro: "mongo_connect_error"

Verifique se o MongoDB está rodando:

```bash
# macOS
brew services list | grep mongodb

# Linux
sudo systemctl status mongod

# Docker
docker ps | grep mongo
```

### Seeds não estão sendo executados novamente

As migrations rodam apenas uma vez. Para re-executar:

```bash
# Conectar ao MongoDB
mongo retech_core

# Remover registros de migration
db.migrations.deleteMany({})

# Limpar dados
db.estados.deleteMany({})
db.municipios.deleteMany({})

# Sair
exit

# Reiniciar a aplicação
go run cmd/api/main.go
```

### Verificar dados no MongoDB

```bash
# Conectar ao MongoDB
mongo retech_core

# Verificar estados
db.estados.countDocuments()  # Deve retornar 27

# Verificar municípios
db.municipios.countDocuments()  # Deve retornar 5570

# Ver um exemplo de estado
db.estados.findOne({sigla: "PE"})

# Ver um exemplo de município
db.municipios.findOne({nome: "Recife"})

# Sair
exit
```

## Próximos Passos

1. Explore a [documentação completa](README.md)
2. Veja a [especificação OpenAPI](http://localhost:8080/docs)
3. Configure o [ambiente de produção](#deploy-produção)
4. Implemente [autenticação via API Key](README.md#-segurança)

## Deploy Produção

Para deploy em produção:

1. Configure as variáveis de ambiente adequadas
2. Use MongoDB Atlas ou cluster gerenciado
3. Configure CORS apropriadamente
4. Habilite logs estruturados
5. Configure backup do MongoDB

Veja mais detalhes na [documentação principal](README.md#-deploy-docker).

