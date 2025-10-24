# 🧪 Teste Redis Pós-Deploy

## ✅ Checklist de Validação

### **1. Verificar Logs do Railway**

Após o deploy, verifique os logs do serviço `retech-core`:

```bash
# Procure por estas mensagens:
✅ "⚡ Redis conectado com sucesso!"
✅ "Server running on port 8080"

# OU (se Redis falhar):
⚠️ "⚠️ Redis não disponível, usando apenas MongoDB"
```

---

### **2. Testar API CEP (Verificar Cache)**

#### **Teste 1: Primeira Request (Origin)**
```bash
curl -X GET "https://api-core.theretech.com.br/cep/01310100" \
  -H "x-api-key: SUA_API_KEY" \
  -w "\n⏱️ Tempo: %{time_total}s\n"
```

**Resultado Esperado:**
```json
{
  "cep": "01310100",
  "logradouro": "Avenida Paulista",
  "bairro": "Bela Vista",
  "localidade": "São Paulo",
  "uf": "SP",
  "source": "viacep",  // ← Veio da API externa
  "cachedAt": "2025-10-24T..."
}
⏱️ Tempo: 0.150s  // ~150ms (primeira vez)
```

#### **Teste 2: Segunda Request (Redis Cache)**
```bash
# MESMA REQUEST (executar imediatamente após)
curl -X GET "https://api-core.theretech.com.br/cep/01310100" \
  -H "x-api-key: SUA_API_KEY" \
  -w "\n⏱️ Tempo: %{time_total}s\n"
```

**Resultado Esperado:**
```json
{
  "cep": "01310100",
  "logradouro": "Avenida Paulista",
  "bairro": "Bela Vista",
  "localidade": "São Paulo",
  "uf": "SP",
  "source": "redis-cache",  // ← 🎯 Veio do Redis!
  "cachedAt": "2025-10-24T..."
}
⏱️ Tempo: 0.003s  // ⚡ ~3ms! (97% mais rápido!)
```

---

### **3. Testar API CNPJ (Verificar Cache)**

#### **Teste 1: Primeira Request**
```bash
curl -X GET "https://api-core.theretech.com.br/cnpj/00000000000191" \
  -H "x-api-key: SUA_API_KEY" \
  -w "\n⏱️ Tempo: %{time_total}s\n"
```

**Resultado Esperado:**
```json
{
  "cnpj": "00000000000191",
  "razao_social": "BANCO DO BRASIL S.A.",
  "nome_fantasia": "BANCO DO BRASIL",
  "source": "brasilapi",  // ← Origem
  "cachedAt": "..."
}
⏱️ Tempo: 0.250s  // ~250ms (primeira vez)
```

#### **Teste 2: Segunda Request (Redis)**
```bash
curl -X GET "https://api-core.theretech.com.br/cnpj/00000000000191" \
  -H "x-api-key: SUA_API_KEY" \
  -w "\n⏱️ Tempo: %{time_total}s\n"
```

**Resultado Esperado:**
```json
{
  "source": "redis-cache",  // ← 🎯 Redis!
  "cachedAt": "..."
}
⏱️ Tempo: 0.002s  // ⚡ ~2ms!
```

---

### **4. Testar API Geografia (Verificar Cache)**

#### **Teste 1: Lista de UFs**
```bash
curl -X GET "https://api-core.theretech.com.br/geo/ufs" \
  -H "x-api-key: SUA_API_KEY" \
  -w "\n⏱️ Tempo: %{time_total}s\n"
```

**1ª Request**: ~30ms (MongoDB)  
**2ª Request**: **<3ms** (Redis) ⚡

#### **Teste 2: UF Específica**
```bash
curl -X GET "https://api-core.theretech.com.br/geo/ufs/SP" \
  -H "x-api-key: SUA_API_KEY" \
  -w "\n⏱️ Tempo: %{time_total}s\n"
```

**1ª Request**: ~20ms (MongoDB)  
**2ª Request**: **<2ms** (Redis) ⚡

---

## 🎯 Métricas de Sucesso

### **Redis Funcionando Corretamente:**
| Métrica | Valor Esperado |
|---------|---------------|
| **Latência 1ª request** | 100-300ms (origin) |
| **Latência 2ª request** | **<5ms** (Redis) ⚡ |
| **Campo `source`** | `"redis-cache"` na 2ª request |
| **Logs Railway** | `"⚡ Redis conectado com sucesso!"` |

### **Redis NÃO Funcionando (Fallback):**
| Métrica | Valor |
|---------|-------|
| **Latência 1ª request** | 100-300ms (origin) |
| **Latência 2ª request** | ~10-50ms (MongoDB) |
| **Campo `source`** | `"mongodb-cache"` |
| **Logs Railway** | `"⚠️ Redis não disponível"` |

---

## 🔧 Troubleshooting

### **Problema 1: `source` sempre `"viacep"` ou `"brasilapi"`**
**Causa**: Cache não está funcionando  
**Solução**:
1. Verificar se `REDIS_URL` está configurada no Railway
2. Verificar logs: `"⚡ Redis conectado com sucesso!"`
3. Testar conexão Redis manualmente (ver abaixo)

### **Problema 2: `source` sempre `"mongodb-cache"`**
**Causa**: Redis não está conectado, mas MongoDB cache funciona  
**Solução**:
1. Verificar se `REDIS_URL` está correta
2. Verificar se serviço Redis está rodando (status verde no Railway)
3. Testar conexão: `redis-cli -u $REDIS_URL ping` (deve retornar `PONG`)

### **Problema 3: Latência ainda alta (>50ms)**
**Causa**: Possível problema de rede Railway  
**Solução**:
1. Verificar região do Redis (deve estar na mesma do `retech-core`)
2. Verificar métricas Railway (CPU, memória)
3. Aguardar alguns minutos (cache warming)

---

## 🧪 Script de Teste Completo

Copie e cole este script para testar tudo de uma vez:

```bash
#!/bin/bash

API_KEY="SUA_API_KEY"
BASE_URL="https://api-core.theretech.com.br"

echo "🧪 TESTE 1: CEP (Origem)"
curl -X GET "$BASE_URL/cep/01310100" -H "x-api-key: $API_KEY" -w "\n⏱️ %{time_total}s\n\n"
sleep 1

echo "🧪 TESTE 2: CEP (Redis Cache)"
curl -X GET "$BASE_URL/cep/01310100" -H "x-api-key: $API_KEY" -w "\n⏱️ %{time_total}s\n\n"
sleep 1

echo "🧪 TESTE 3: CNPJ (Origem)"
curl -X GET "$BASE_URL/cnpj/00000000000191" -H "x-api-key: $API_KEY" -w "\n⏱️ %{time_total}s\n\n"
sleep 1

echo "🧪 TESTE 4: CNPJ (Redis Cache)"
curl -X GET "$BASE_URL/cnpj/00000000000191" -H "x-api-key: $API_KEY" -w "\n⏱️ %{time_total}s\n\n"
sleep 1

echo "🧪 TESTE 5: Geografia (Origem)"
curl -X GET "$BASE_URL/geo/ufs/SP" -H "x-api-key: $API_KEY" -w "\n⏱️ %{time_total}s\n\n"
sleep 1

echo "🧪 TESTE 6: Geografia (Redis Cache)"
curl -X GET "$BASE_URL/geo/ufs/SP" -H "x-api-key: $API_KEY" -w "\n⏱️ %{time_total}s\n\n"

echo "✅ Testes concluídos!"
echo "📊 Verifique se as 2ªs requests têm 'source: redis-cache' e latência <5ms"
```

---

## 📊 Monitoramento Contínuo

### **Railway Logs**
```bash
# Verificar se Redis está conectado
railway logs | grep -i redis

# Contar cache hits
railway logs | grep "redis-cache" | wc -l

# Contar cache misses
railway logs | grep "viacep\|brasilapi" | wc -l
```

### **Métricas Esperadas (após 24h)**
- **Cache Hit Rate**: >85%
- **Latência P50**: <5ms
- **Latência P99**: <20ms
- **Redis Memory**: <100MB (para 10K CEPs/CNPJs cacheados)

---

## ✅ Conclusão

Se os testes mostrarem:
- ✅ `"source": "redis-cache"` na 2ª request
- ✅ Latência <5ms na 2ª request
- ✅ Logs `"⚡ Redis conectado com sucesso!"`

**🎉 REDIS ESTÁ FUNCIONANDO PERFEITAMENTE! 🚀**

---

**Última Atualização**: 2025-10-24  
**Próximo Passo**: Monitorar cache hit rate e ajustar TTLs se necessário 📊

