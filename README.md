# retech-core — Core APIs for The Retech

Centraliza **serviços utilitários** e **integrações públicas** para acelerar diversos projetos (web, mobile, backoffice). Foco em **estabilidade**, **observabilidade** e **padronização** de respostas.

---

## ✨ Objetivos

* Expor endpoints **estáveis** para dados públicos/derivados (UFs, municípios, RAs do DF, CEP, bancos, feriados etc.).
* **Padronizar** contratos (envelope de resposta, erros RFC 7807) entre projetos.
* **Desacoplar** front-ends das fontes externas (IBGE, BrasilAPI, etc.), com **cache** e **rate limiting**.
* **Observabilidade** pronta (Prometheus + OpenTelemetry) e **resiliência** (retries, circuit breaker).

---

## 📐 Arquitetura (resumo)

```
Cliente → retech-core (API Gateway utilitário)
            ├─ /geo/*   → (cache + normalização) → IBGE/BrasilAPI/fonte DF
            ├─ /cep/*   → (cache + fallback) → ViaCEP/BrasilAPI
            ├─ /utils/* → serviços puros (CPF/CNPJ, phone, currency, slug)
            ├─ /biz/*   → bancos, feriados, etc. (cache)
            └─ /health,/metrics,/docs
```

* **Cache**: Redis (TTL configurável).
* **Circuit breaker**: impedir cascata de falhas de upstream.
* **Rate limiting**: token bucket por IP/chave.
* **Auth**: API Key opcional por rota; JWT passthrough quando aplicável.

---

## 📦 Estrutura de pastas (Go + Gin)

```
retech-core/
  cmd/api/main.go
  internal/
    http/ (rotas, handlers)
    services/ (casos de uso)
    clients/ (ibge, brasilapi, viacep, gdf)
    cache/ (redis)
    config/
    observability/ (otel, prometheus)
    middleware/ (auth, rate-limit, cors, idempotency)
    domain/ (DTOs/validators)
  pkg/
  build/
    Dockerfile
    docker-compose.yml
```

---

## 🔐 Segurança

* **API Key** por projeto/cliente (header `x-api-key`).
* **JWT passthrough** (quando consumir serviços internos que exigem identidade).
* **CORS** configurável por ambiente.
* **Auditoria** básica (request-id, caller, rota, status, latência).

---

## 📊 Observabilidade

* **/metrics** (Prometheus)
* **OpenTelemetry** traces (HTTP client/server)
* **Log estruturado** (JSON; correlação por `X-Request-ID`)

---

## 📑 Convenções de resposta

### Sucesso

```json
{
  "success": true,
  "code": "OK",
  "data": { "..." },
  "meta": { "cache": {"hit": true, "ttl": 300} }
}
```

### Erro (RFC 7807)

```json
{
  "type": "https://retech-core/errors/upstream-timeout",
  "title": "Upstream Timeout",
  "status": 504,
  "detail": "IBGE não respondeu em 5s",
  "instance": "/geo/ufs?q=pern",
  "traceId": "01H..."
}
```

---

## 🚀 Endpoints propostos (v1)

### 0) Infra

* `GET /health` → { status: "ok", uptime }
* `GET /version` → { name, version, commit, buildDate }
* `GET /docs` → Swagger UI / Redoc
* `GET /metrics` → Prometheus

### 1) GEO (Estados, Municípios, RA-DF)

* `GET /geo/ufs`
  **Query**: `q` (opcional, busca parcial client-side: o serviço pode retornar todos e filtrar)
  **Resposta** (normalizada):

  ```json
  [{ "id": 26, "sigla": "PE", "nome": "Pernambuco", "regiao": {"sigla": "NE", "nome": "Nordeste"} }]
  ```

  **Fonte**: IBGE `/localidades/estados` (com cache).

* `GET /geo/municipios/{uf}`
  **Resposta**: array de municípios (id, nome, microrregião, mesorregião opcional).
  **Fonte**: IBGE `/localidades/estados/{UF}/municipios` (cache).

* `GET /geo/df/ras`
  **Resposta**: 33 Regiões Administrativas do DF (id simbólico, nome, sigla RA, bbox opcional).
  **Fonte**: dataset próprio (seed baseado em GDF) versionado no repo; opção de cliente GDF.

### 2) CEP

* `GET /cep/{cep}`
  **Resposta** (normalizada): logradouro, bairro, cidade, UF, IBGE, lat/lng (quando disponível).
  **Fonte** primária: BrasilAPI. **Fallback**: ViaCEP. **Cache** agressivo (7–30d).

### 3) Documentos & utilidades

* `POST /utils/cpf/validate`
  **Body**: `{ "cpf": "00000000191" }` → `{ "valid": true }` (offline, sem upstream)
* `POST /utils/cnpj/validate`
  **Body**: `{ "cnpj": "11222333000181" }` → `{ "valid": true }`
* `POST /utils/phone/format`
  **Body**: `{ "phone": "+55 (48) 99999-6679" }` → `{ "e164": "+554899996679", "national": "(48) 99996-6679" }`
* `POST /utils/slugify`
  **Body**: `{ "text": "Sudoeste / Octogonal" }` → `{ "slug": "sudoeste-octogonal" }`

### 4) Negócio (dados abertos úteis)

* `GET /biz/bancos`
  Lista bancos ativos (código, nome). **Fonte**: BrasilAPI / Bacen (cache 24h).
* `GET /biz/feriados/{ano}`
  Feriados nacionais (e opcional `?uf=SC`). **Fonte**: BrasilAPI / Nager.Date (cache 24h).

---

## 🧰 Exemplo de uso (cURL)

```bash
# Estados
curl -s 'https://api.retech-core.com/v1/geo/ufs' | jq

# Municípios por UF
curl -s 'https://api.retech-core.com/v1/geo/municipios/PE' | jq

# Regiões Administrativas do DF
curl -s 'https://api.retech-core.com/v1/geo/df/ras' | jq

# CEP
curl -s 'https://api.retech-core.com/v1/cep/88160116' | jq

# Validar CNPJ
curl -s -X POST 'https://api.retech-core.com/v1/utils/cnpj/validate' \
  -H 'Content-Type: application/json' \
  -d '{"cnpj":"11222333000181"}' | jq
```

---

## 🧪 Contratos (OpenAPI sketch)

```yaml
openapi: 3.1.0
info:
  title: retech-core API
  version: 1.0.0
servers:
  - url: https://api.retech-core.com/v1
paths:
  /geo/ufs:
    get:
      summary: Lista UFs
      parameters:
        - in: query
          name: q
          schema: { type: string }
      responses:
        '200':
          description: OK
  /geo/municipios/{uf}:
    get:
      parameters:
        - in: path
          name: uf
          required: true
          schema: { type: string, minLength: 2, maxLength: 2 }
      responses:
        '200': { description: OK }
```

---

## ⚙️ Configuração (ENV)

```
PORT=8080
ENV=production
API_KEY_REQUIRED=true
REDIS_URL=redis://localhost:6379
CACHE_TTL_SECONDS=3600
UPSTREAM_TIMEOUT_MS=5000
RATE_LIMIT_RPS=10
CORS_ORIGINS=https://*.theretech.com,https://*.brbit.com
```

---

## 🐳 Deploy (Docker)

**Dockerfile (Go)**

```dockerfile
FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/retech-core ./cmd/api

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/retech-core /retech-core
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/retech-core"]
```

**docker-compose.yml (dev)**

```yaml
version: '3.9'
services:
  api:
    build: .
    ports: ["8080:8080"]
    env_file: .env
    depends_on: [redis]
  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
```

---

## 🛡️ Boas práticas

* **Timeouts** e **retries** exponenciais para upstreams.
* **Cache key** determinística e invalidação por rota.
* **Schema validation** (ex.: go-playground/validator) em inputs.
* **E2E contract tests** (Dredd/Pact) para não quebrar clientes.
* **Idempotência** em POSTs sensíveis (chave idempotency em header).

---

## 🗺️ Roadmap (sugestão)

* [ ] `/geo/municipios/{uf}/search?q=` (filtro local com índice em memória)
* [ ] `/biz/cnaes` e `/biz/naturesa-juridica` (catálogos públicos)
* [ ] `/utils/coordinates/geocode` (Nominatim, com rate-limit estrito)
* [ ] **Bulk** endpoints (`/geo/municipios/bulk`)
* [ ] **Webhooks** para cache-invalidation de listas estáticas

---

## 📄 Licença

Definir conforme política interna (MIT/Proprietária).

---

### Notas finais

* Para **RAs do DF**, manter **seed estático** versionado (fonte: GDF) e opcionalmente cliente para atualização periódica.
* Para **UF q=**, como o IBGE não filtra de forma confiável, retornar todos e aplicar **filtro local** (normalizado p/ acentuação).
