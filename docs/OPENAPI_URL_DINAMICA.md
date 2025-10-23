# ✅ OpenAPI com URL Única e Dinâmica

**Data:** 23 de Outubro de 2025  
**Status:** ✅ **IMPLEMENTADO**

---

## 🎯 **OBJETIVO**

Ter **apenas UMA URL** no OpenAPI que muda automaticamente conforme o ambiente (desenvolvimento ou produção), sem duplicação de servers.

---

## 🔧 **IMPLEMENTAÇÃO**

### **1. Placeholder no OpenAPI YAML**

**Arquivo:** `internal/docs/openapi.yaml`

**Antes:**
```yaml
servers:
  - url: https://api-core.theretech.com.br
    description: Produção
  - url: http://localhost:8080
    description: Desenvolvimento Local
```

**Depois:**
```yaml
servers:
  - url: __API_BASE_URL__
    description: API Base URL
```

---

### **2. Substituição Dinâmica no Backend**

**Arquivo:** `internal/http/handlers/docs.go`

**Código:**
```go
// Substituir URL do servidor dinamicamente
apiBaseURL := os.Getenv("API_BASE_URL")
if apiBaseURL == "" {
    apiBaseURL = "https://api-core.theretech.com.br"
}

// Substituir o placeholder __API_BASE_URL__ no YAML
yamlString := string(content)
yamlString = strings.Replace(
    yamlString,
    "__API_BASE_URL__",
    apiBaseURL,
    -1, // Substituir todas as ocorrências
)
```

**Como funciona:**
1. Lê `API_BASE_URL` do ambiente
2. Se não existir, usa `https://api-core.theretech.com.br` (produção)
3. Substitui **TODAS** as ocorrências de `__API_BASE_URL__` no YAML
4. Retorna YAML com URL correta

---

### **3. Exemplos de Código Também Dinâmicos**

Todos os exemplos de `curl` no OpenAPI também usam o placeholder:

**Antes:**
```bash
curl https://api-core.theretech.com.br/geo/ufs \
  -H "X-API-Key: sua_api_key_aqui"
```

**Depois:**
```bash
curl __API_BASE_URL__/geo/ufs \
  -H "X-API-Key: sua_api_key_aqui"
```

**Resultado em localhost:**
```bash
curl http://localhost:8080/geo/ufs \
  -H "X-API-Key: sua_api_key_aqui"
```

**Resultado em produção:**
```bash
curl https://api-core.theretech.com.br/geo/ufs \
  -H "X-API-Key: sua_api_key_aqui"
```

---

## 🌍 **CONFIGURAÇÃO POR AMBIENTE**

### **Desenvolvimento Local**

**Docker Compose:** `build/docker-compose.yml`
```yaml
services:
  api:
    environment:
      - API_BASE_URL=http://localhost:8080
```

**Resultado no OpenAPI:**
```yaml
servers:
  - url: http://localhost:8080
    description: API Base URL
```

---

### **Produção (Railway)**

**Variável de ambiente:**
```
API_BASE_URL=https://api-core.theretech.com.br
```

**Resultado no OpenAPI:**
```yaml
servers:
  - url: https://api-core.theretech.com.br
    description: API Base URL
```

---

## ✅ **VANTAGENS**

1. ✅ **Uma única URL** no Redoc (não confunde o usuário)
2. ✅ **Exemplos corretos** conforme ambiente
3. ✅ **Sem duplicação** de configuração
4. ✅ **Simples de manter** (um único placeholder)
5. ✅ **Funciona em qualquer ambiente** (dev, staging, prod)

---

## 🧪 **TESTES**

### **Verificar URL Dinâmica**

**Desenvolvimento:**
```bash
curl http://localhost:8080/openapi.yaml | grep "url:"
# Resultado: - url: http://localhost:8080
```

**Produção:**
```bash
curl https://api-core.theretech.com.br/openapi.yaml | grep "url:"
# Resultado: - url: https://api-core.theretech.com.br
```

---

### **Verificar Exemplos**

```bash
curl http://localhost:8080/openapi.yaml | grep "curl "
# Resultado: curl http://localhost:8080/geo/ufs
```

---

## 📋 **PLACEHOLDERS SUBSTITUÍDOS**

Todas as ocorrências de `__API_BASE_URL__` são substituídas:

1. ✅ `servers[0].url`
2. ✅ Descrição da API (exemplo de curl)
3. ✅ Schema `apiBaseURL.example`
4. ✅ `ApiKeyAuth` description (exemplo de curl)

---

## 🚀 **DEPLOY**

### **Local**
```bash
docker compose -f build/docker-compose.yml up --build
# URL será: http://localhost:8080
```

### **Produção (Railway)**
```bash
# Configurar variável:
API_BASE_URL=https://api-core.theretech.com.br

# Deploy automático
git push origin main
```

---

## 📊 **RESULTADO FINAL**

| Ambiente | `API_BASE_URL` | URL no OpenAPI |
|----------|---------------|----------------|
| **Dev Local** | `http://localhost:8080` | `http://localhost:8080` |
| **Produção** | `https://api-core.theretech.com.br` | `https://api-core.theretech.com.br` |
| **Staging** | `https://staging.api.theretech.com.br` | `https://staging.api.theretech.com.br` |

---

## ✨ **BENEFÍCIOS PARA O DESENVOLVEDOR**

Quando um desenvolvedor acessa `/docs`:

1. 🎯 **Vê apenas UMA URL** (a correta para aquele ambiente)
2. 📋 **Exemplos de código funcionam** direto (copy/paste)
3. 🚀 **Try it out** usa a URL certa automaticamente
4. 📚 **Documentação sempre atualizada** sem intervenção manual

---

**Implementação 100% funcional e profissional! 🎯**

