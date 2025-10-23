# Configuração da URL Base da API

**Data**: 2025-10-23  
**Status**: ✅ **IMPLEMENTADO**

---

## 🎯 **O QUE FOI FEITO:**

A URL base da API agora é **100% dinâmica** e configurável via variável de ambiente!

Isso permite que a documentação (`/painel/docs`) e outros endpoints sempre usem a URL correta sem hardcoding.

---

## ⚙️ **CONFIGURAÇÃO:**

### **1. Desenvolvimento (Local)**

No arquivo `.env` (na raiz do projeto):

```env
# API Base URL (usado na documentação)
API_BASE_URL=http://localhost:8080
```

**OU** diretamente no `docker-compose.yml` (já configurado):

```yaml
api:
  environment:
    - API_BASE_URL=http://localhost:8080
```

---

### **2. Produção (Railway)**

#### **Passos:**

1. Acesse o projeto no Railway
2. Vá em **retech-core** (backend) → **Variables**
3. Adicione a variável:

```
API_BASE_URL=https://api-core.theretech.com.br
```

4. Clique em **Add**
5. O Railway fará **redeploy automático**

---

## 🔄 **COMO FUNCIONA:**

### **Backend (`/me/config`)**

```go
// internal/http/handlers/tenant.go

func (h *TenantHandler) GetMyConfig(c *gin.Context) {
    // Busca da variável de ambiente
    apiBaseURL := os.Getenv("API_BASE_URL")
    if apiBaseURL == "" {
        // Fallback para produção
        apiBaseURL = "https://api-core.theretech.com.br"
    }

    c.JSON(http.StatusOK, gin.H{
        "apiBaseURL": apiBaseURL,  // ← Retorna para o frontend
        // ...
    })
}
```

### **Frontend (`/painel/docs`)**

```tsx
// app/painel/docs/page.tsx

const config = await getMyConfig();  // Chama backend

// Usa a URL dinâmica
const curlExample = `curl ${config.apiBaseURL}/geo/ufs \\
  -H "X-API-Key: sua_api_key_aqui"`;
```

---

## 📊 **BENEFÍCIOS:**

✅ **Sem hardcoding**: Frontend nunca tem URL hardcoded  
✅ **Ambiente-específico**: Development usa `localhost`, Production usa `api-core.theretech.com.br`  
✅ **Documentação sempre correta**: Exemplos de código sempre com a URL certa  
✅ **Fácil de mudar**: Basta alterar a env var, não precisa mudar código  

---

## 🧪 **TESTAR:**

### **Development:**

1. Garanta que `docker-compose.yml` tem:
   ```yaml
   environment:
     - API_BASE_URL=http://localhost:8080
   ```

2. Faça rebuild:
   ```bash
   docker compose -f build/docker-compose.yml up --build -d
   ```

3. Acesse: `http://localhost:3000/painel/docs`

4. Verifique que os exemplos usam `http://localhost:8080`

### **Production:**

1. Adicione env var no Railway:
   ```
   API_BASE_URL=https://api-core.theretech.com.br
   ```

2. Acesse: `https://core.theretech.com.br/painel/docs`

3. Verifique que os exemplos usam `https://api-core.theretech.com.br`

---

## 📋 **CHECKLIST:**

- [x] ✅ Backend lê `API_BASE_URL` da env
- [x] ✅ Fallback configurado
- [x] ✅ Endpoint `/me/config` retorna URL
- [x] ✅ Frontend usa URL dinâmica
- [x] ✅ `docker-compose.yml` configurado
- [ ] ⏳ Railway env var configurada
- [ ] ⏳ Testado em produção

---

## 🚀 **PRÓXIMOS PASSOS:**

1. ✅ Fazer commit das mudanças
2. ✅ Push para o repositório
3. ⏳ Adicionar env var no Railway
4. ⏳ Verificar em produção

---

**Documentação sempre atualizada e correta! 🎯✨**

