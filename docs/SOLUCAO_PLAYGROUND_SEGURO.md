# 🔒 Solução: Playground Seguro com SEO

## 🎯 Objetivo

Manter o playground **funcional para SEO** (sem forçar login) MAS **protegido contra abuso**.

---

## ✅ Solução Escolhida: API Key Demo Compartilhada

### **Como Funciona:**

1. **Criar API Key especial** `rtc_demo_playground`
2. **Configurar limites agressivos:**
   - 100 requests/dia **total** (compartilhado por todos)
   - 10 requests/minuto **total**
   - Rate limit **por IP** também: 20 requests/dia/IP
3. **Hardcoded no frontend** (apenas na página `/playground`)
4. **Scopes**: `cep`, `cnpj`, `geo` (read-only)

---

## 🏗️ Implementação

### **1. Backend: Criar API Key Demo (Manual)**

No admin (`/admin/apikeys`):
- **Name**: `Demo Playground (Public)`
- **Key**: `rtc_demo_playground_2024` (customizada)
- **Tenant**: Criar tenant especial "Retech Demo"
- **Scopes**: `cep`, `cnpj`, `geo`
- **Daily Limit**: 100 (compartilhado)
- **Rate Limit**: 10 req/min

### **2. Backend: Middleware de Rate Limit por IP**

Adicionar middleware que:
- Rastreia IP + API Key
- Limita 20 requests/dia por IP
- Retorna erro com CTA: "Crie conta grátis para mais requests"

### **3. Frontend: Usar API Key Demo**

```typescript
// app/playground/page.tsx
const DEMO_API_KEY = process.env.NEXT_PUBLIC_DEMO_API_KEY || 'rtc_demo_playground_2024';

const res = await fetch(url, {
  headers: {
    'X-API-Key': DEMO_API_KEY
  }
});
```

### **4. Mensagem de Erro Amigável**

Quando atingir limite:
```json
{
  "error": "rate_limit_exceeded",
  "message": "Limite de demonstração atingido",
  "cta": "Crie uma conta grátis e ganhe 1.000 requests/dia",
  "link": "https://core.theretech.com.br/painel/register"
}
```

---

## 🎯 Benefícios

| Aspecto | Solução |
|---------|---------|
| **SEO** | ✅ Playground funciona sem login |
| **Segurança** | ✅ Rate limit duplo (global + por IP) |
| **UX** | ✅ Funciona imediatamente, sem fricção |
| **Conversão** | ✅ CTA forte quando atingir limite |
| **Custo** | ✅ Controlado (max 100 req/dia) |

---

## 🚀 Alternativas Consideradas

### **❌ Opção 1: Rotas 100% Públicas (sem auth)**
- Problema: Abuso ilimitado
- Custo: Potencialmente infinito

### **❌ Opção 2: Remover Playground Público**
- Problema: Perde SEO
- Problema: Perde conversão (fricção no cadastro)

### **✅ Opção 3: API Key Demo (ESCOLHIDA)**
- Benefícios: Balanceia SEO, segurança e conversão
- Custo: Previsível e controlado

---

## 📊 Monitoramento

### **Métricas a Acompanhar:**
1. **Uso da API Key Demo**: quantas requests/dia
2. **IPs únicos**: quantos usuários testam
3. **Taxa de conversão**: % que cria conta após usar playground
4. **Abuso**: IPs que atingem limite rapidamente

### **Alertas:**
- Se uso > 80 requests/dia → aumentar limite OU reduzir
- Se muitos IPs batem limite → ajustar rate limit por IP

---

## 🔧 Implementação Técnica

### **Mudanças Necessárias:**

1. **Backend:**
   - ✅ Remover rotas `/public/*` (já temos)
   - ✅ Criar API Key demo manualmente
   - ⏳ Adicionar middleware de rate limit por IP
   - ⏳ Mensagem de erro com CTA

2. **Frontend:**
   - ⏳ Adicionar `NEXT_PUBLIC_DEMO_API_KEY`
   - ⏳ Usar API key nas requests do playground
   - ⏳ Exibir modal de erro com CTA quando atingir limite

3. **Infraestrutura:**
   - ⏳ Adicionar variável de ambiente no Railway (frontend)
   - ⏳ Criar tenant e API key demo no banco

---

## ✅ Status

- [x] Análise de segurança
- [x] Definição da solução
- [ ] Implementação backend (middleware)
- [ ] Implementação frontend (usar API key)
- [ ] Testes de rate limiting
- [ ] Deploy e monitoramento

---

**Conclusão**: Esta solução **preserva o SEO** (playground funciona sem login), **protege contra abuso** (rate limit duplo), e **aumenta conversão** (CTA quando atingir limite). 🎯

