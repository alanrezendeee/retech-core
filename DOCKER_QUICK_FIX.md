# ⚡ Fix Rápido - Seeds não encontrados no Docker

## 🎯 Solução Rápida (1 minuto)

Se você está vendo o erro:
```
arquivo estados.json não encontrado
```

Execute estes comandos:

```bash
# 1. Parar containers
cd build
docker-compose down

# 2. Voltar para raiz e verificar seeds
cd ..
ls seeds/

# 3. Se arquivos não existirem, copiar:
cp ~/Downloads/estados.json seeds/
cp ~/Downloads/municipios.json seeds/

# 4. Rebuild e iniciar
cd build
docker-compose up --build
```

## ✅ Solução Completa (2 minutos)

Use o script automático:

```bash
# Da raiz do projeto
./docker-setup.sh
```

O script faz tudo automaticamente! ✨

## 📋 Checklist Manual

Se preferir fazer manualmente:

- [ ] **Arquivos existem?**
  ```bash
  ls -la seeds/
  # Deve mostrar: estados.json e municipios.json
  ```

- [ ] **Permissões corretas?**
  ```bash
  chmod 644 seeds/*.json
  ```

- [ ] **Volume montado?** Verifique em `build/docker-compose.yml`:
  ```yaml
  volumes:
    - ../seeds:/app/seeds:ro
  ```

- [ ] **Rebuild da imagem:**
  ```bash
  cd build
  docker-compose build --no-cache
  docker-compose up
  ```

## 🔍 Verificar se Funcionou

Após iniciar, você deve ver nos logs:

```
✅ Logs esperados:
{"level":"info","message":"Executando migrations e seeds..."}
{"level":"info","message":"[migration] Aplicando 001_seed_estados..."}
{"level":"info","message":"[seed] Carregando estados de: seeds/estados.json"}
{"level":"info","message":"[seed] 27 estados inseridos com sucesso"}
{"level":"info","message":"[migration] Aplicando 002_seed_municipios..."}
{"level":"info","message":"[seed] 5570 municípios inseridos com sucesso"}
{"level":"info","message":"listening on :8080 (env=development)"}
```

Teste a API:
```bash
curl http://localhost:8080/health
curl http://localhost:8080/geo/ufs | jq
```

## 🆘 Ainda com Problemas?

Veja o guia completo: [DOCKER_TROUBLESHOOTING.md](DOCKER_TROUBLESHOOTING.md)

## 📝 O que foi corrigido

As seguintes mudanças foram feitas para resolver o problema:

1. ✅ `build/docker-compose.yml` - Adicionado volume dos seeds
2. ✅ `build/Dockerfile` - WORKDIR mudado de `/` para `/app`
3. ✅ `internal/bootstrap/migrations.go` - Busca melhorada de arquivos
4. ✅ `docker-setup.sh` - Script automático criado

---

**Data**: 2025-10-20
**Versão**: 0.3.0

