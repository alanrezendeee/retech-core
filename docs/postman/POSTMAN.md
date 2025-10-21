# 📮 Guia do Postman - Retech Core API

Este guia explica como importar e usar a collection do Postman para testar a API Retech Core.

## 📥 Importando a Collection

### 1. Importar Collection

1. Abra o Postman
2. Clique em **Import** (canto superior esquerdo)
3. Arraste o arquivo `postman_collection.json` ou clique em **Upload Files**
4. Confirme a importação

### 2. Importar Environment (Opcional, mas recomendado)

1. Clique em **Import** novamente
2. Arraste o arquivo `postman_environment.json`
3. Confirme a importação
4. Selecione o environment "Retech Core - Local" no dropdown (canto superior direito)

## 🎯 Estrutura da Collection

A collection está organizada em 5 grupos:

### 0. Infra
- ✅ Health Check
- ✅ Version
- ✅ Docs (Redoc)
- ✅ OpenAPI YAML

### 1. GEO - Estados
- ✅ Listar Todos os Estados (27 UFs)
- ✅ Buscar Estados (filtro por nome/sigla)
- ✅ Buscar Estado por Sigla (PE, SP, RJ)

### 2. GEO - Municípios
- ✅ Listar Municípios de PE
- ✅ Listar Municípios de SP
- ✅ Buscar Municípios por Nome (Recife, São Paulo)
- ✅ Buscar por Nome (sem filtro UF)
- ✅ Buscar Município por Código IBGE
- ✅ Listar Municípios com Filtro por UF

### 3. Tenants
- ✅ Criar Tenant
- ✅ Listar Tenants
- ✅ Buscar Tenant por ID
- ✅ Atualizar Tenant
- ✅ Deletar Tenant

### 4. API Keys
- ✅ Criar API Key (com auto-save da key)
- ✅ Testar Endpoint Protegido
- ✅ Rotacionar API Key
- ✅ Revogar API Key

### 5. Exemplos de Casos de Uso
- ✅ Listar Estados do Nordeste
- ✅ Listar Cidades do Interior de PE
- ✅ Buscar Todas as Cidades chamadas 'Santa'
- ✅ Setup Completo (Tenant + API Key)
- ✅ Validar Dados de Endereço

## 🚀 Como Usar

### Passo 1: Verificar se a API está rodando

1. Execute a API localmente:
   ```bash
   cd /path/to/retech-core
   go run cmd/api/main.go
   ```

2. No Postman, execute a requisição **"0. Infra → Health Check"**
3. Você deve receber um status `200 OK` com:
   ```json
   {
     "status": "ok",
     "timestamp": "...",
     "database": "connected"
   }
   ```

### Passo 2: Testar Endpoints GEO

#### Listar Estados
Execute: **"1. GEO - Estados → Listar Todos os Estados"**

Resposta esperada:
```json
{
  "success": true,
  "code": "OK",
  "data": [
    {
      "id": 11,
      "sigla": "RO",
      "nome": "Rondônia",
      "regiao": {
        "id": 1,
        "sigla": "N",
        "nome": "Norte"
      }
    },
    // ... 26 outros estados
  ]
}
```

#### Buscar Municípios
Execute: **"2. GEO - Municípios → Listar Municípios de PE"**

Resposta esperada: Array com 185 municípios de Pernambuco

#### Buscar por Código IBGE
Execute: **"2. GEO - Municípios → Buscar Município por Código IBGE - Recife"**

Resposta esperada:
```json
{
  "success": true,
  "code": "OK",
  "data": {
    "id": 2611606,
    "nome": "Recife",
    "microrregiao": {
      "id": 26017,
      "nome": "Recife",
      "mesorregiao": {
        "id": 2605,
        "nome": "Metropolitana de Recife",
        "UF": {
          "id": 26,
          "sigla": "PE",
          "nome": "Pernambuco",
          "regiao": {
            "id": 2,
            "sigla": "NE",
            "nome": "Nordeste"
          }
        }
      }
    },
    "regiao-imediata": {
      "id": 260001,
      "nome": "Recife",
      "regiao-intermediaria": {
        "id": 2601,
        "nome": "Recife",
        "UF": { /* ... */ }
      }
    }
  }
}
```

### Passo 3: Criar Tenant e API Key

#### 3.1 Criar Tenant
Execute: **"3. Tenants → Criar Tenant"**

Body padrão:
```json
{
  "tenantId": "empresa-exemplo",
  "name": "Empresa Exemplo LTDA",
  "email": "contato@exemplo.com",
  "active": true
}
```

Resposta: Status `201 Created` com dados do tenant

**💡 Dica**: Guarde o `id` retornado. Ele será usado em outras requisições.

#### 3.2 Criar API Key
Execute: **"4. API Keys → Criar API Key"**

Body:
```json
{
  "tenantId": "empresa-exemplo",
  "name": "Chave de Produção"
}
```

**🎉 Automático**: A collection tem um script que salva automaticamente a API key retornada na variável de ambiente `api_key`!

#### 3.3 Testar Autenticação
Execute: **"4. API Keys → Testar Endpoint Protegido"**

A requisição usará automaticamente a API key salva no header `x-api-key`.

## 🔧 Variáveis de Environment

As seguintes variáveis são utilizadas:

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `base_url` | URL base da API | `http://localhost:8080` |
| `tenant_id` | ID do tenant criado | Salve manualmente após criar tenant |
| `api_key` | API Key para autenticação | Salva automaticamente ao criar key |
| `key_id` | ID da API key | Salve manualmente se necessário |

### Como Editar Variáveis

1. Clique no ícone de **olho** (👁️) no canto superior direito
2. Clique em **Edit** ao lado do environment "Retech Core - Local"
3. Edite os valores
4. Clique em **Save**

## 📝 Exemplos de Testes

### Exemplo 1: Buscar Estados do Nordeste

```
GET {{base_url}}/geo/ufs?q=NE
```

Retorna todos os estados da região Nordeste.

### Exemplo 2: Buscar Cidades com "Santa" no nome

```
GET {{base_url}}/geo/municipios?q=santa
```

Retorna até 100 municípios que contenham "santa" no nome.

### Exemplo 3: Validar Município por Código IBGE

```
GET {{base_url}}/geo/municipios/id/3550308
```

Valida que o código 3550308 corresponde à cidade de São Paulo/SP.

### Exemplo 4: Filtrar Municípios por Estado

```
GET {{base_url}}/geo/municipios?uf=RS
```

Lista todos os municípios do Rio Grande do Sul.

## 🎯 Fluxo de Teste Completo

1. **Health Check** - Verificar se API está online
2. **Listar Estados** - Ver todos os 27 estados
3. **Buscar Estado PE** - Pegar detalhes de Pernambuco
4. **Listar Municípios PE** - Ver todos os 185 municípios
5. **Buscar Recife** - Buscar por nome
6. **Buscar Recife por Código** - Validar código IBGE
7. **Criar Tenant** - Setup inicial
8. **Criar API Key** - Gerar chave de autenticação
9. **Testar Endpoint Protegido** - Validar autenticação

## 🐛 Troubleshooting

### Erro: "Connection refused"
- Verifique se a API está rodando
- Confirme a porta correta (padrão: 8080)
- Teste: `curl http://localhost:8080/health`

### Erro: 404 "Estado não encontrado"
- Verifique se as migrations foram executadas
- Confirme que os seeds foram carregados
- Verifique os logs da aplicação

### Erro: 401 "Unauthorized"
- Verifique se a API key está configurada corretamente
- Confirme que o header `x-api-key` está sendo enviado
- Crie uma nova API key se necessário

### API Key não salva automaticamente
- Verifique se o environment está selecionado
- Vá em **Console** (View → Show Postman Console) para ver os logs
- Salve manualmente: copie a key da resposta e cole na variável `api_key`

## 📊 Códigos IBGE Úteis

Para facilitar seus testes:

| Município | UF | Código IBGE |
|-----------|----|-----------:|
| Recife | PE | 2611606 |
| São Paulo | SP | 3550308 |
| Rio de Janeiro | RJ | 3304557 |
| Belo Horizonte | MG | 3106200 |
| Salvador | BA | 2927408 |
| Brasília | DF | 5300108 |
| Fortaleza | CE | 2304400 |
| Manaus | AM | 1302603 |
| Curitiba | PR | 4106902 |
| Porto Alegre | RS | 4314902 |

## 💡 Dicas Pro

1. **Use o Runner**: Execute toda a collection de uma vez
   - Click direito na collection → **Run collection**

2. **Testes Automatizados**: As requisições de API Key salvam automaticamente as keys
   - Veja os scripts em **Tests** nas requisições

3. **Organize por Casos de Uso**: A pasta "5. Exemplos de Casos de Uso" tem cenários reais

4. **Documente**: Use a descrição das requisições para notas pessoais

5. **Compartilhe**: Exporte e compartilhe a collection com o time

## 📚 Recursos Adicionais

- [Documentação completa](README.md)
- [Quick Start Guide](QUICK_START.md)
- [Changelog](CHANGELOG.md)
- [OpenAPI Docs](http://localhost:8080/docs) (quando API estiver rodando)

## 🤝 Contribuindo

Encontrou um problema ou tem sugestões para a collection?
- Abra uma issue
- Adicione novas requisições úteis
- Melhore a documentação

---

**Última atualização**: 2025-10-20
**Versão da API**: 0.3.0
**Versão da Collection**: 1.0.0

