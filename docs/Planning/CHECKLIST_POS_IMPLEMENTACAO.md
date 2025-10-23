# ✅ CHECKLIST PÓS-IMPLEMENTAÇÃO

**Propósito:** Garantir que após implementar uma nova API, todos os componentes do sistema sejam atualizados.

---

## 📋 CHECKLIST GERAL (Para cada nova API)

### **1. Backend** ✅
- [ ] Handler implementado (`internal/http/handlers/`)
- [ ] Rotas registradas (`internal/http/router.go`)
- [ ] **SCOPE** adicionado (`internal/auth/scope_middleware.go`) 🔒
  - [ ] Adicionar scope no `validScopes` map
  - [ ] Status: `true` (ativo)
- [ ] **SCOPE** aplicado no router (`router.go`)
  - [ ] `auth.RequireScope(apikeys, "nome_api")`
- [ ] **SCOPE** no frontend (`components/apikeys/apikey-drawer.tsx`)
  - [ ] Adicionar em `availableScopes` array
- [ ] Testes unitários
- [ ] Testes de integração
- [ ] Cache implementado
- [ ] Rate limiting funcionando
- [ ] Usage logging ativo

### **2. Documentação Técnica** 📚
- [ ] **OpenAPI/Redoc** (`internal/docs/openapi.yaml`)
  - [ ] Endpoint documentado
  - [ ] Schemas definidos
  - [ ] Exemplos de request/response
  - [ ] Códigos de erro documentados
- [ ] **README** atualizado (se necessário)
- [ ] **CHANGELOG** atualizado

### **3. Frontend - Landing Page** 🎨
- [ ] **Badge de status** atualizado (`app/page.tsx`)
  - De: `<Badge className="bg-blue-600">Fase 2</Badge>`
  - Para: `<Badge className="bg-green-600">Disponível</Badge>`
- [ ] **Descrição** da API (se necessário)
- [ ] **Screenshot** (se relevante)

### **4. Frontend - Painel do Desenvolvedor** 👨‍💻
- [ ] **Documentação** (`app/painel/docs/page.tsx`)
  - [ ] Endpoint adicionado à lista
  - [ ] Exemplo de código atualizado
  - [ ] Link para Redoc (se complexo)

### **5. Frontend - Dashboard do Desenvolvedor** 📊
- [ ] **Métricas** (`app/painel/dashboard/page.tsx`)
  - [ ] Contador de requests por API
  - [ ] Gráfico de uso por endpoint
  - [ ] Top APIs mais usadas
- [ ] **Usage** (`app/painel/usage/page.tsx`)
  - [ ] Gráfico de uso por API
  - [ ] Tabela de endpoints por API
  - [ ] Filtros por API

### **6. Frontend - Admin Dashboard** 👑
- [ ] **Analytics** (`app/admin/analytics/page.tsx`)
  - [ ] Métricas globais por API
  - [ ] Top APIs do sistema
  - [ ] Uso por tenant (por API)
  - [ ] Gráficos de crescimento

### **7. Backend - Analytics & Metrics** 📈
- [ ] **Endpoint `/admin/analytics`** atualizado
  - [ ] Métricas por API
  - [ ] Agregação por endpoint
- [ ] **Endpoint `/me/usage`** atualizado
  - [ ] `byAPI` field adicionado
  - [ ] Breakdown por endpoint
- [ ] **Endpoint `/me/stats`** atualizado
  - [ ] Contador por API (se relevante)

### **8. Database** 🗄️
- [ ] **Schema** atualizado (se necessário)
- [ ] **Indexes** otimizados para queries por API
  - [ ] `api_usage_logs`: index em `endpoint` + `api_name`
  - [ ] Agregações por API eficientes
- [ ] **Migration** (se necessário)

### **9. Testes E2E** 🧪
- [ ] Teste da nova API (request/response)
- [ ] Teste de rate limiting
- [ ] Teste de cache
- [ ] Teste de fallback (se houver)

### **10. Deploy** 🚀
- [ ] Variáveis de ambiente (se necessário)
- [ ] Secrets configurados (API keys externas)
- [ ] Build bem-sucedido
- [ ] Health check OK
- [ ] Smoke test em produção

---

## 🎯 CHECKLIST ESPECÍFICO: CEP

### **Backend**
- [x] Criar handler `internal/http/handlers/cep.go`
- [x] Implementar client ViaCEP
- [x] Implementar fallback Brasil API
- [x] Cache de 7 dias
- [x] Enriquecimento com coordenadas (Nominatim - opcional)
- [x] Registrar rota `GET /cep/:codigo`
- [x] Rate limiting aplicado
- [x] Usage logging com campo `api_name: "cep"`

### **OpenAPI/Redoc**
- [x] Adicionar tag "CEP"
- [x] Documentar endpoint `GET /cep/:codigo`
- [x] Schema `CEPResponse`
- [x] Exemplos de sucesso/erro
- [x] Limites e cache documentados

### **Landing Page**
- [x] Badge: `Fase 2` → `Disponível`
- [x] Cor: `border-blue-300` → `border-green-300`

### **Painel Docs**
- [x] Adicionar seção "CEP"
- [x] Exemplo de código curl
- [x] Response exemplo
- [x] Rate limit info

### **Analytics (Backend)**
- [x] Campo `api_name` no log de uso
- [x] Agregação por API implementada
- [x] Endpoint `/admin/analytics` com `byAPI`

### **Dashboard Dev**
- [x] Gráfico "Uso por API"
- [x] Filtro por API
- [x] Tabela com breakdown

### **Dashboard Admin**
- [x] Métricas por API
- [x] Top APIs
- [x] Crescimento por API

---

## 📊 ESTRUTURA DE DADOS: Logging por API

### **API Usage Log (MongoDB)**
```json
{
  "apiKey": "rtc_abc123.xyz789",
  "tenantId": "tenant-123",
  "endpoint": "/cep/01310-100",
  "api_name": "cep",           // ✅ NOVO CAMPO
  "method": "GET",
  "statusCode": 200,
  "responseTime": 45,
  "date": "2025-10-23",
  "timestamp": "2025-10-23T14:30:00Z",
  "ip": "192.168.1.1",
  "userAgent": "axios/1.0"
}
```

### **Agregação por API**
```javascript
// Admin Analytics - Uso por API
db.api_usage_logs.aggregate([
  { $match: { tenantId: "tenant-123" } },
  { $group: {
      _id: "$api_name",
      count: { $sum: 1 },
      avgResponseTime: { $avg: "$responseTime" }
    }
  },
  { $sort: { count: -1 } }
])

// Resultado:
[
  { _id: "geografia", count: 1500, avgResponseTime: 50 },
  { _id: "cep", count: 850, avgResponseTime: 45 },
  { _id: "cnpj", count: 320, avgResponseTime: 120 }
]
```

---

## 🎨 FRONTEND: Componentes a Criar/Atualizar

### **1. Gráfico "Uso por API" (Developer Dashboard)**
```tsx
// app/painel/dashboard/page.tsx
<Card>
  <CardHeader>
    <CardTitle>Uso por API</CardTitle>
  </CardHeader>
  <CardContent>
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie data={usageByAPI} dataKey="count" nameKey="api" />
      </PieChart>
    </ResponsiveContainer>
  </CardContent>
</Card>
```

### **2. Tabela de Endpoints por API (Developer Usage)**
```tsx
// app/painel/usage/page.tsx
<Tabs defaultValue="all">
  <TabsList>
    <TabsTrigger value="all">Todas</TabsTrigger>
    <TabsTrigger value="cep">CEP</TabsTrigger>
    <TabsTrigger value="geografia">Geografia</TabsTrigger>
  </TabsList>
  <TabsContent value="cep">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Endpoint</TableHead>
          <TableHead>Requests</TableHead>
          <TableHead>Tempo Médio</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {endpoints.filter(e => e.api === 'cep').map(...)}
      </TableBody>
    </Table>
  </TabsContent>
</Tabs>
```

### **3. Admin Analytics - Top APIs**
```tsx
// app/admin/analytics/page.tsx
<Card>
  <CardHeader>
    <CardTitle>Top 5 APIs Mais Usadas</CardTitle>
  </CardHeader>
  <CardContent>
    <div className="space-y-4">
      {topAPIs.map(api => (
        <div key={api.name}>
          <div className="flex justify-between mb-1">
            <span>{api.name}</span>
            <span>{api.count.toLocaleString()} req</span>
          </div>
          <Progress value={api.percentage} />
        </div>
      ))}
    </div>
  </CardContent>
</Card>
```

---

## 🚀 ORDEM DE IMPLEMENTAÇÃO SUGERIDA

1. **Backend Handler** (API funcionando)
2. **Usage Logging** com `api_name` (tracking)
3. **OpenAPI/Redoc** (documentação técnica)
4. **Landing Page** badge (comunicação)
5. **Backend Analytics** endpoints (dados para frontend)
6. **Frontend Dashboards** (visualização)
7. **Painel Docs** (guia do desenvolvedor)
8. **Testes** (qualidade)
9. **Deploy** (produção)

---

## 📝 TEMPLATE DE COMMIT

```bash
feat(cep): implementa API de consulta de CEP

Backend:
- Handler com ViaCEP + fallback Brasil API
- Cache 7 dias
- Rate limiting + usage logging
- Campo api_name="cep" nos logs

Frontend:
- Badge "Disponível" na landing page
- Documentação no painel/docs
- Gráfico "Uso por API" no dashboard
- Analytics admin com breakdown por API

Docs:
- OpenAPI atualizado
- Exemplos de código
- Rate limits documentados

Closes #XX
```

---

**✅ Use este checklist para cada nova API implementada!**

