# Seeds - Dados para Popular o Banco

Este diretório deve conter os arquivos JSON com dados para popular o banco de dados.

## Arquivos Necessários

- `estados.json` - Lista de estados brasileiros (27 estados)
- `municipios.json` - Lista de municípios brasileiros (5570 municípios)

## Como Usar

1. Coloque os arquivos `estados.json` e `municipios.json` neste diretório
2. Execute a aplicação normalmente
3. O sistema detectará automaticamente os arquivos e executará as migrations/seeds

## Localização dos Arquivos

O sistema busca os arquivos nas seguintes localizações (em ordem):

1. Diretório `seeds/` (este diretório)
2. `~/Downloads/` (conveniente para desenvolvimento)
3. Diretório `data/`
4. Diretório raiz do projeto

## Formato dos Arquivos

### estados.json

```json
[
  {
    "id": 26,
    "sigla": "PE",
    "nome": "Pernambuco",
    "regiao": {
      "id": 2,
      "sigla": "NE",
      "nome": "Nordeste"
    }
  }
]
```

### municipios.json

```json
[
  {
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
    }
  }
]
```

## Migrations

O sistema mantém um registro das migrations executadas na collection `migrations`. 
As seeds só serão executadas uma vez. Para re-executar:

1. Remova o documento da migration da collection `migrations`
2. Ou limpe as collections `estados` e `municipios`
3. Reinicie a aplicação

## Fonte dos Dados

Os dados são baseados nas APIs públicas do IBGE:
- https://servicodados.ibge.gov.br/api/v1/localidades/estados
- https://servicodados.ibge.gov.br/api/v1/localidades/municipios

