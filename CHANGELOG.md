# Changelog - retech-core

## [Unreleased] - 2025-10-20

### ✨ Adicionado

#### Sistema de Migrations e Seeds
- Implementado gerenciador de migrations automático
- Seeds para popular estados brasileiros (27 UFs)
- Seeds para popular municípios brasileiros (5570 municípios)
- Migrations executam automaticamente na inicialização
- Sistema rastreia migrations executadas na collection `migrations`
- Suporte a múltiplas localizações de arquivos seed

#### Domínio - Estados e Municípios
- Modelo `Estado` com ID, sigla, nome e região
- Modelo `Municipio` com dados completos IBGE (microrregião, mesorregião, regiões imediata/intermediária)
- Modelos aninhados: `Regiao`, `Microrregiao`, `Mesorregiao`, `UFReference`, `RegiaoImediata`, `RegiaoIntermediaria`

#### Repositórios
- `EstadosRepo` - CRUD completo para estados
  - `FindAll()` - Lista todos os estados ordenados por nome
  - `FindBySigla()` - Busca por sigla (ex: PE, SP)
  - `FindByID()` - Busca por ID do IBGE
  - `InsertMany()` - Inserção em lote
  - `Count()` - Contagem de registros
  - `DeleteAll()` - Limpeza para re-seed

- `MunicipiosRepo` - CRUD completo para municípios
  - `FindAll()` - Lista todos os municípios ordenados por nome
  - `FindByUF()` - Filtra por estado
  - `FindByID()` - Busca por código IBGE
  - `Search()` - Busca por nome (case-insensitive, parcial)
  - `InsertMany()` - Inserção em lotes de 1000 para otimização
  - `Count()` - Contagem de registros
  - `DeleteAll()` - Limpeza para re-seed

#### Endpoints GEO
- `GET /geo/ufs` - Lista estados (com filtro opcional `?q=`)
- `GET /geo/ufs/:sigla` - Busca estado por sigla
- `GET /geo/municipios` - Lista municípios (com filtros `?uf=` e `?q=`)
- `GET /geo/municipios/:uf` - Lista municípios por estado
- `GET /geo/municipios/id/:id` - Busca município por código IBGE

#### Handlers
- `GeoHandler` - Gerencia todos os endpoints de GEO
- Respostas padronizadas com envelope `SuccessResponse`
- Erros padronizados seguindo RFC 7807 (`ErrorResponse`)
- Validações de entrada
- Tratamento adequado de erros 404, 400, 500

#### Índices MongoDB
- Estados: índice único por `id` e `sigla`
- Municípios: índice único por `id`
- Municípios: índice por `microrregiao.mesorregiao.UF.sigla` (busca por UF)
- Municípios: índice por `nome` (busca textual)
- Criação automática de índices na inicialização

#### Documentação
- README atualizado com endpoints implementados
- Separação clara entre "Implementado" e "Planejado"
- Seção de Migrations e Seeds
- Exemplos de uso atualizados
- Roadmap detalhado
- Estrutura do projeto documentada
- Notas de implementação
- `QUICK_START.md` - Guia rápido de configuração
- `seeds/README.md` - Documentação de seeds
- `env.example` - Exemplo de configuração

### 🔧 Modificado

#### Router
- Atualizado para receber repositórios de estados e municípios
- Adicionadas rotas GEO

#### Main
- Integrado sistema de migrations na inicialização
- Criação de índices antes de iniciar servidor
- Timeout de 10 minutos para migrations (suporta grandes volumes)
- Instanciação de repositórios de estados e municípios

### 📊 Estatísticas

- **Arquivos criados**: 12
- **Arquivos modificados**: 3
- **Linhas de código**: ~1500
- **Endpoints novos**: 5
- **Collections MongoDB**: 2 (estados, municipios)
- **Registros populados**: 5597 (27 estados + 5570 municípios)

### 🎯 Cobertura de Funcionalidades

#### Implementado ✅
- [x] Sistema de migrations automático
- [x] Seeds de estados e municípios
- [x] Endpoints de consulta de estados
- [x] Endpoints de consulta de municípios
- [x] Busca e filtros
- [x] Índices para performance
- [x] Documentação completa

#### Próximos Passos 🚧
- [ ] Cache com Redis
- [ ] Rate limiting
- [ ] Endpoints de CEP
- [ ] Validadores de CPF/CNPJ
- [ ] Prometheus metrics
- [ ] OpenTelemetry traces

### 🔍 Detalhes Técnicos

**Stack:**
- Go 1.22+
- MongoDB 6.0+
- Gin Web Framework
- Zerolog (logging estruturado)

**Padrões:**
- Repository pattern
- Envelope de resposta padronizado
- RFC 7807 para erros
- Migrations versionadas
- Índices otimizados

**Performance:**
- Inserção em lotes (1000 registros/batch)
- Índices MongoDB estratégicos
- Queries otimizadas
- Timeout configurável

### 📝 Notas de Migração

Para atualizar de versão anterior:

1. Os seeds serão executados automaticamente na primeira inicialização
2. Coloque os arquivos `estados.json` e `municipios.json` no diretório `seeds/`
3. As migrations rodam apenas uma vez
4. Para re-executar, remova os documentos da collection `migrations`

### 🙏 Créditos

Dados geográficos fornecidos pelo IBGE:
- https://servicodados.ibge.gov.br/api/v1/localidades/estados
- https://servicodados.ibge.gov.br/api/v1/localidades/municipios

