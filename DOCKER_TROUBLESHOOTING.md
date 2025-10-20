# 🐳 Docker Troubleshooting - Retech Core

Guia de resolução de problemas comuns ao executar a aplicação com Docker.

## ❌ Erro: "arquivo estados.json não encontrado"

### Sintomas
```
{"level":"fatal","error":"erro ao aplicar migration 001_seed_estados: arquivo estados.json não encontrado","time":"2025-10-20T22:59:07Z","message":"migration_error"}
```

### Causa
Os arquivos de seed (`estados.json` e `municipios.json`) não estão disponíveis no container Docker.

### Solução

#### Opção 1: Garantir que os arquivos estão no diretório correto

```bash
# Verificar se os arquivos existem
ls -la seeds/
# Deve mostrar: estados.json e municipios.json

# Se não existirem, copiar:
cp ~/Downloads/estados.json seeds/
cp ~/Downloads/municipios.json seeds/
```

#### Opção 2: Rebuild do container (se necessário)

```bash
cd build

# Parar containers
docker-compose down

# Rebuild com cache limpo
docker-compose build --no-cache

# Subir novamente
docker-compose up
```

#### Opção 3: Verificar volume montado

Verifique se o volume está corretamente montado no `docker-compose.yml`:

```yaml
services:
  api:
    volumes:
      - ../seeds:/app/seeds:ro  # ← Esta linha é essencial!
```

#### Opção 4: Verificar permissões

```bash
# Dar permissão de leitura aos arquivos
chmod 644 seeds/estados.json
chmod 644 seeds/municipios.json

# Verificar
ls -la seeds/
```

---

## ❌ Erro: "connection refused" ao conectar MongoDB

### Sintomas
```
{"level":"fatal","error":"mongo_connect_error","message":"..."}
```

### Soluções

#### 1. Verificar se o MongoDB está rodando

```bash
docker-compose ps

# Deve mostrar:
# NAME                COMMAND                  SERVICE   STATUS
# build-mongo-1       "docker-entrypoint.s…"   mongo     Up
# build-api-1         "/app/retech-core"       api       Up
```

#### 2. Verificar logs do MongoDB

```bash
docker-compose logs mongo

# Se houver erro, reinicie:
docker-compose restart mongo
```

#### 3. Verificar conexão

```bash
# Entrar no container da API
docker-compose exec api sh

# (Não funcionará com distroless, então teste externamente)

# Testar conexão do host
docker exec -it build-mongo-1 mongosh --eval "db.adminCommand('ping')"
```

---

## ❌ Erro: "port is already allocated"

### Sintomas
```
Error: Bind for 0.0.0.0:8080 failed: port is already allocated
```

### Soluções

#### 1. Verificar o que está usando a porta

```bash
# macOS/Linux
lsof -i :8080

# Se encontrar processo, matar:
kill -9 <PID>
```

#### 2. Mudar a porta no docker-compose

```yaml
services:
  api:
    ports:
      - "8081:8080"  # Usar porta 8081 no host
```

---

## ❌ Seeds não carregam (collection vazia)

### Diagnóstico

```bash
# Conectar ao MongoDB
docker exec -it build-mongo-1 mongosh retech_core

# Verificar estados
db.estados.countDocuments()  # Deve retornar 27

# Verificar municípios
db.municipios.countDocuments()  # Deve retornar 5570

# Verificar migrations
db.migrations.find()

# Sair
exit
```

### Soluções

#### 1. Se as migrations não rodaram

```bash
# Ver logs da aplicação
docker-compose logs api

# Procurar por:
# - "Executando migrations e seeds..."
# - "Popular estados brasileiros"
# - "Popular municípios brasileiros"
```

#### 2. Re-executar migrations

```bash
# Conectar ao MongoDB
docker exec -it build-mongo-1 mongosh retech_core

# Limpar migrations
db.migrations.deleteMany({})
db.estados.deleteMany({})
db.municipios.deleteMany({})
exit

# Reiniciar API
docker-compose restart api

# Ver logs
docker-compose logs -f api
```

---

## ❌ Container reinicia constantemente

### Sintomas
```bash
docker-compose ps
# NAME        STATUS
# build-api-1   Restarting
```

### Diagnóstico

```bash
# Ver logs completos
docker-compose logs api

# Ver últimas 50 linhas
docker-compose logs --tail 50 api
```

### Causas comuns

1. **Erro fatal na inicialização** (seeds não encontrados)
2. **Porta já em uso**
3. **MongoDB não está pronto**
4. **Configuração incorreta**

### Solução

```bash
# Parar tudo
docker-compose down

# Verificar arquivos de seed
ls -la seeds/

# Limpar volumes (CUIDADO: apaga dados)
docker-compose down -v

# Subir novamente
docker-compose up --build
```

---

## 🔍 Comandos Úteis de Diagnóstico

### Logs

```bash
# Ver logs em tempo real
docker-compose logs -f

# Logs apenas da API
docker-compose logs -f api

# Logs apenas do MongoDB
docker-compose logs -f mongo

# Últimas 100 linhas
docker-compose logs --tail 100 api
```

### Status dos Containers

```bash
# Ver status
docker-compose ps

# Ver processos
docker-compose top

# Ver uso de recursos
docker stats
```

### Acessar Container

```bash
# Entrar no container do MongoDB
docker exec -it build-mongo-1 mongosh retech_core

# Listar containers
docker ps

# Inspecionar container
docker inspect build-api-1
```

### Verificar Volumes

```bash
# Listar volumes
docker volume ls

# Inspecionar volume do MongoDB
docker volume inspect build_mongo-data

# Remover volumes não usados (CUIDADO!)
docker volume prune
```

---

## 🧹 Limpeza e Reset

### Reset Completo (perde dados)

```bash
# Parar containers
docker-compose down

# Remover volumes (apaga banco)
docker-compose down -v

# Limpar imagens
docker-compose down --rmi all

# Rebuild do zero
docker-compose build --no-cache

# Subir novamente
docker-compose up
```

### Reset apenas do banco

```bash
# Parar containers
docker-compose down

# Remover apenas volume do MongoDB
docker volume rm build_mongo-data

# Subir novamente (vai recriar volume)
docker-compose up
```

---

## ✅ Checklist de Verificação

Antes de reportar um problema, verifique:

- [ ] Arquivos `estados.json` e `municipios.json` existem em `seeds/`
- [ ] Docker e Docker Compose estão atualizados
- [ ] Porta 8080 está livre
- [ ] Porta 27017 está livre
- [ ] Você tem permissão para criar volumes
- [ ] `.env` está configurado (opcional, mas recomendado)
- [ ] Você está executando `docker-compose up` de dentro do diretório `build/`

---

## 📝 Setup Ideal

### Estrutura de diretórios esperada:

```
retech-core/
├── seeds/
│   ├── estados.json       ← DEVE existir
│   └── municipios.json    ← DEVE existir
├── build/
│   ├── docker-compose.yml
│   └── Dockerfile
└── cmd/api/main.go
```

### Comando correto para subir:

```bash
# A partir da raiz do projeto
cd build
docker-compose up

# Ou com rebuild
cd build
docker-compose up --build

# Ou em background
cd build
docker-compose up -d
```

---

## 🆘 Ainda com problemas?

1. **Verifique os logs completos**: `docker-compose logs api > api.log`
2. **Teste local sem Docker**: `go run cmd/api/main.go`
3. **Verifique versões**:
   ```bash
   docker --version    # Docker 20.10+
   docker-compose --version  # Docker Compose 2.0+
   go version          # Go 1.22+
   ```

4. **Issues conhecidos**: Verifique o GitHub/docs do projeto

---

**Última atualização**: 2025-10-20
**Versão da API**: 0.3.0

