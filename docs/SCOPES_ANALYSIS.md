# 🔒 Análise Completa do Sistema de Scopes

## 📊 **Situação Atual**

### **APIs Disponíveis (Todas READ-ONLY):**

| API | Métodos | Endpoints | Scopes Necessários |
|-----|---------|-----------|-------------------|
| **GEO** | `GET` | `/geo/ufs`, `/geo/ufs/:sigla`, `/geo/municipios`, etc | `geo` ou `all` |
| **CEP** | `GET` | `/cep/:codigo` | `cep` ou `all` |
| **CNPJ** | `GET` | `/cnpj/:numero` | `cnpj` ou `all` |

### **✅ Todas as APIs são READ-ONLY (apenas GET)**

**Não há operações de escrita (POST, PUT, DELETE) nas APIs públicas.**

---

## 🤔 **Precisamos de `:read` e `:write`?**

### **RESPOSTA: NÃO! (Por enquanto)**

**Motivos:**

1. ✅ **Todas as APIs são consulta apenas (GET)**
   - GEO: Consulta estados e municípios
   - CEP: Consulta endereços
   - CNPJ: Consulta dados de empresas
   - Nenhuma permite criar/editar/deletar

2. ✅ **Futuro previsível também será READ-ONLY**
   - CPF: Consulta (não faz sentido criar CPF via API)
   - FIPE: Consulta (dados públicos)
   - Moedas: Consulta (cotações)
   - Bancos: Consulta (lista de bancos)

3. ✅ **Simplicidade > Complexidade desnecessária**
   - Menos scopes = Mais fácil de gerenciar
   - UX melhor para desenvolvedores
   - Menos chances de erro de configuração

---

## 🎯 **Formato Padronizado ATUAL (RECOMENDADO)**

### **Scopes sem sufixo:**

| Scope | Descrição | Endpoints |
|-------|-----------|-----------|
| `geo` | Acesso a dados geográficos | `/geo/*` |
| `cep` | Acesso a consulta de CEP | `/cep/:codigo` |
| `cnpj` | Acesso a dados de empresas | `/cnpj/:numero` |
| `all` | Acesso total a todas as APIs | `/*` |

### **Futuros (planejados):**

| Scope | Descrição | Status |
|-------|-----------|--------|
| `cpf` | Validação e consulta CPF | 🔜 Fase 3 |
| `fipe` | Tabela FIPE (veículos) | 🔜 Fase 3 |
| `moedas` | Cotações de moedas | 🔜 Fase 4 |
| `bancos` | Lista de bancos brasileiros | 🔜 Fase 4 |

---

## 🚨 **Quando usar `:read` e `:write`?**

### **Cenário hipotético (FUTURO):**

Se no futuro tivermos APIs que **criam/modificam dados**, exemplo:

```
POST   /documentos          → Criar documento
PUT    /documentos/:id      → Atualizar documento
DELETE /documentos/:id      → Deletar documento
GET    /documentos          → Listar documentos
```

**Aí sim usaríamos:**
- `documentos:read` → Apenas GET
- `documentos:write` → POST, PUT, DELETE
- `documentos` → Ambos (shorthand para read+write)

---

## ✅ **Decisão Final: MANTER FORMATO ATUAL**

### **Padrão:**
```
geo, cep, cnpj, all
```

### **NÃO usar:**
```
geo:read, geo:write ❌
```

### **Motivos:**

1. ✅ **Todas as APIs são read-only**
2. ✅ **Mais simples e intuitivo**
3. ✅ **Menos verboso**
4. ✅ **Backend já padronizado (hasScope remove sufixo automaticamente)**
5. ✅ **Frontend já cria no formato correto**
6. ✅ **Produção já está padronizado**

---

## 🛠️ **Implementação Atual (CORRETA)**

### **Backend (`scope_middleware.go`):**
```go
func hasScope(scopes []string, scope string) bool {
	// Extrair apenas o nome do scope (ex: "geo:read" → "geo")
	// ✅ Suporta ambos os formatos (retrocompatível)
	scopeName := scope
	if idx := strings.Index(scope, ":"); idx != -1 {
		scopeName = scope[:idx]
	}
	
	for _, s := range scopes {
		sName := s
		if idx := strings.Index(s, ":"); idx != -1 {
			sName = s[:idx]
		}
		
		// geo == geo ✅
		// geo:read == geo ✅
		// geo == geo:read ✅
		if sName == scopeName || s == scope || s == "all" {
			return true
		}
	}
	return false
}
```

### **Frontend (`apikey-drawer.tsx`):**
```typescript
const availableScopes = [
  { value: 'geo', label: '🗺️ GEO - Dados geográficos' },
  { value: 'cep', label: '📮 CEP - Consulta de endereços' },
  { value: 'cnpj', label: '🏢 CNPJ - Dados de empresas' },
  { value: 'all', label: '⭐ ALL - Acesso total' },
];
```

### **Rotas protegidas (`router.go`):**
```go
geoGroup.Use(
	auth.AuthAPIKey(apikeys),          // Valida API Key
	auth.RequireScope(apikeys, "geo"), // ✅ Requer 'geo' ou 'all'
	rateLimiter.Middleware(),
	usageLogger.Middleware(),
)

cepGroup.Use(
	auth.AuthAPIKey(apikeys),
	auth.RequireScope(apikeys, "cep"), // ✅ Requer 'cep' ou 'all'
	rateLimiter.Middleware(),
	usageLogger.Middleware(),
)

cnpjGroup.Use(
	auth.AuthAPIKey(apikeys),
	auth.RequireScope(apikeys, "cnpj"), // ✅ Requer 'cnpj' ou 'all'
	rateLimiter.Middleware(),
	usageLogger.Middleware(),
)
```

---

## 📈 **Evolução Futura**

### **Quando adicionar novos scopes:**

1. ✅ Adicionar em `scope_middleware.go` → `validScopes` map
2. ✅ Adicionar em `apikey-drawer.tsx` → `availableScopes` array
3. ✅ Criar handler para a nova API
4. ✅ Registrar rota em `router.go` com `RequireScope`
5. ✅ Atualizar Redoc (`openapi.yaml`)
6. ✅ Atualizar landing page
7. ✅ Atualizar documentação do desenvolvedor

### **Se precisar de `:write` no futuro:**

```go
// Exemplo hipotético
documentsGroup.GET("/", auth.RequireScope(apikeys, "documents:read"))
documentsGroup.POST("/", auth.RequireScope(apikeys, "documents:write"))
documentsGroup.PUT("/:id", auth.RequireScope(apikeys, "documents:write"))
documentsGroup.DELETE("/:id", auth.RequireScope(apikeys, "documents:write"))
```

---

## 🎯 **Conclusão**

### **✅ FORMATO ATUAL ESTÁ PERFEITO!**

**Não precisa mudar nada:**
- Backend: ✅ Padronizado
- Frontend: ✅ Padronizado
- Produção: ✅ Padronizado
- Documentação: ✅ Clara

**Formato:**
```
geo, cep, cnpj, all
```

**Compatibilidade:**
- ✅ Aceita `geo:read` (legacy, se existir)
- ✅ Cria apenas `geo` (formato novo)
- ✅ Ambos funcionam perfeitamente

---

## 📝 **Checklist de Validação**

- [x] Todas as APIs são GET-only
- [x] Backend suporta ambos os formatos (retrocompatível)
- [x] Frontend cria no formato padronizado
- [x] Produção já está no formato correto
- [x] Documentação alinhada
- [x] Proteção de rotas funcionando
- [x] Rate limiting por tenant (não por scope)
- [x] Scopes validados na criação de API keys

**🎉 TUDO PADRONIZADO E FUNCIONANDO!**

