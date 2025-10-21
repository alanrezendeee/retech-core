# 📝 Changelog - 21 de Outubro de 2025

## ✅ Ajustes Implementados

### 1. 🎨 Melhorado Select de Tenant na Drawer de API Key

**Problema:** Select básico sem informações visuais
**Solução:**
- Adicionado avatar colorido com iniciais do tenant
- Exibição do `tenantId` abaixo do nome
- Altura aumentada para 56px (h-14)
- Hover effect e cursor pointer

**Arquivos modificados:**
- `retech-core-admin/components/apikeys/apikey-drawer.tsx`

**Visual:**
```
┌─────────────────────────────────────┐
│ [AB] Super Admin                    │
│      tenant-20251021145821          │
├─────────────────────────────────────┤
│ [GC] Global Corp                    │
│      tenant-global-001              │
└─────────────────────────────────────┘
```

---

### 2. 📊 Dashboard - Dados Reais dos Cards

**Problema:** Cards exibindo valores zerados (mock data)
**Solução:**
- Integrado com endpoint `/admin/stats`
- Loading state enquanto busca dados
- Tratamento de erros com toast
- Dados atualizados em tempo real

**Dados exibidos:**
- ✅ Total de Tenants (totalTenants)
- ✅ API Keys Ativas (activeAPIKeys / totalAPIKeys)
- ✅ Requests Hoje (requestsToday)
- ✅ Requests Mês (requestsMonth)

**Arquivos modificados:**
- `retech-core-admin/app/admin/dashboard/page.tsx`
- `retech-core-admin/lib/api/admin.ts` (já existia)

**Exemplo de resposta da API:**
```json
{
  "totalTenants": 2,
  "activeTenants": 2,
  "totalAPIKeys": 1,
  "activeAPIKeys": 1,
  "requestsToday": 0,
  "requestsMonth": 0,
  "systemUptime": "83ns",
  "timestamp": "2025-10-21T19:26:23Z"
}
```

---

### 3. 🔔 Dashboard - Atividade Recente Funcional

**Problema:** Seção vazia sem dados
**Solução:**
- Geração de atividades baseadas em dados reais
- Exibição dinâmica com timestamp
- Estado vazio com mensagem informativa
- Hover effects e animações

**Atividades geradas:**
1. "Sistema iniciado" - quando há tenants/keys ativos
2. "Requisições da API" - quando há requests hoje

**Visual:**
```
┌────────────────────────────────────────┐
│  🔵 Sistema iniciado            19:26  │
│     2 tenant(s) e 1 API key(s) ativos  │
├────────────────────────────────────────┤
│  🔵 Requisições da API          17:26  │
│     124 requisições processadas hoje   │
└────────────────────────────────────────┘
```

---

### 4. 🛠️ Corrigido Rate Limit no TenantDrawer

**Problema 1:** Erro "uncontrolled → controlled input"
**Causa:** Campos `requestsPerDay`/`requestsPerMinute` vinham `undefined` do backend
**Solução:**
```typescript
if (tenant.rateLimit && tenant.rateLimit.requestsPerDay) {
  setCustomRateLimit(true);
  setRateLimit({
    requestsPerDay: tenant.rateLimit.requestsPerDay || 1000,
    requestsPerMinute: tenant.rateLimit.requestsPerMinute || 60,
  });
}
```

**Problema 2:** Campos bloqueados para edição
**Solução:** Aplicada mesma correção do `/admin/settings`:
- `onChange` permite valor vazio temporariamente
- `onBlur` aplica valor padrão se vazio

**Arquivos modificados:**
- `retech-core-admin/components/tenants/tenant-drawer.tsx`

---

### 5. 🎯 Correção: Inputs de Número Editáveis

**Problema:** Impossível apagar todos os dígitos e digitar novo valor
**Causa:** Fallback `|| 1000` aplicado no `onChange`
**Solução:**
```typescript
// Antes (❌):
onChange={(e) => handleInputChange('field', parseInt(e.target.value) || 1000)}

// Depois (✅):
onChange={(e) => {
  const value = e.target.value === '' ? '' : parseInt(e.target.value);
  handleInputChange('field', value);
}}
onBlur={(e) => {
  if (e.target.value === '' || parseInt(e.target.value) < 1) {
    handleInputChange('field', 1000); // Fallback só ao sair
  }
}}
```

**Aplicado em:**
- `/admin/settings` - Requests por Dia/Minuto
- `TenantDrawer` - Rate Limit personalizado

---

### 6. 🔄 Normalização PascalCase ↔ camelCase

**Problema:** Backend (Go) usa `RequestsPerDay`, Frontend (TS) usa `requestsPerDay`
**Solução:**
```typescript
// Ao carregar (PascalCase → camelCase):
const normalized = {
  defaultRateLimit: {
    requestsPerDay: data.defaultRateLimit?.RequestsPerDay || 1000,
    requestsPerMinute: data.defaultRateLimit?.RequestsPerMinute || 60,
  },
};

// Ao salvar (camelCase → PascalCase):
const payload = {
  defaultRateLimit: {
    RequestsPerDay: settings.defaultRateLimit.requestsPerDay,
    RequestsPerMinute: settings.defaultRateLimit.requestsPerMinute,
  },
};
```

**Arquivos modificados:**
- `retech-core-admin/app/admin/settings/page.tsx`

---

## 📋 Resumo das Correções

| #  | Issue                                    | Status | Impacto        |
|----|------------------------------------------|--------|----------------|
| 1️⃣  | Select de Tenant pouco informativo       | ✅ OK  | UX melhorado   |
| 2️⃣  | Dashboard com dados mock                 | ✅ OK  | Dados reais    |
| 3️⃣  | Atividade Recente vazia                  | ✅ OK  | Funcional      |
| 4️⃣  | Erro uncontrolled input (Rate Limit)     | ✅ OK  | Sem erros      |
| 5️⃣  | Inputs bloqueados para edição            | ✅ OK  | Editável       |
| 6️⃣  | Incompatibilidade PascalCase/camelCase   | ✅ OK  | Normalizado    |

---

## 🧪 Como Testar

### Dashboard:
1. Acessar `http://localhost:3001/admin/dashboard`
2. Verificar se cards mostram dados reais (2 tenants, 1 API key, etc)
3. Verificar seção "Atividade Recente" com dados

### API Key Drawer:
1. Acessar `/admin/apikeys`
2. Clicar "Nova API Key"
3. Abrir select de Tenant
4. Verificar visual melhorado com avatares

### Rate Limit Personalizado:
1. Acessar `/admin/tenants`
2. Editar um tenant
3. Ativar "Limite Personalizado"
4. Digitar valores (ex: 5000, 100)
5. Salvar e verificar se persiste

### Inputs de Número:
1. Acessar `/admin/settings`
2. Campo "Requests por Dia"
3. Apagar todos os dígitos (Ctrl+A → Del)
4. Digitar novo valor (ex: 3000)
5. ✅ Deve funcionar sem preencher automaticamente

---

## 📁 Arquivos Modificados

```
retech-core-admin/
├── app/admin/
│   ├── dashboard/page.tsx          ✅ Dados reais + Atividade
│   └── settings/page.tsx           ✅ Inputs editáveis + Normalização
├── components/
│   ├── apikeys/apikey-drawer.tsx   ✅ Select melhorado
│   └── tenants/tenant-drawer.tsx   ✅ Rate Limit corrigido
└── lib/api/
    └── admin.ts                     (já existia)
```

---

## 🎉 Resultado Final

✅ Dashboard 100% funcional com dados reais
✅ Atividade recente exibindo informações
✅ Rate Limit personalizado salvando corretamente
✅ Inputs de número totalmente editáveis
✅ Select de Tenant com visual profissional
✅ Sem erros de uncontrolled input
✅ Compatibilidade PascalCase/camelCase

---

**Próximos Passos:**
1. Implementar endpoint `/admin/activity` para logs reais de atividade
2. Adicionar cálculo de crescimento (%) nos cards
3. Melhorar analytics com gráficos
4. Implementar filtros de data no dashboard

---

**Data:** 21 de Outubro de 2025
**Versão:** v1.2.0
**Status:** ✅ Todas as correções aplicadas e testadas

