#!/bin/bash

# Script para normalizar scopes de API keys
# Converte "geo:read" → "geo", mantém "cep", "cnpj", "all"

echo "🔄 Normalizando scopes de API keys..."

# Conectar ao MongoDB em produção
# ATENÇÃO: Execute este comando manualmente em produção:
echo ""
echo "Execute no MongoDB de PRODUÇÃO:"
echo ""
cat << 'EOF'
db.apikeys.updateMany(
  { "scopes": { $in: ["geo:read", "geo:write"] } },
  [
    {
      $set: {
        scopes: {
          $map: {
            input: "$scopes",
            as: "scope",
            in: {
              $cond: {
                if: { $in: ["$$scope", ["geo:read", "geo:write"]] },
                then: "geo",
                else: "$$scope"
              }
            }
          }
        }
      }
    }
  ]
);

db.apikeys.updateMany(
  {},
  [
    {
      $set: {
        scopes: { $setUnion: ["$scopes", []] }
      }
    }
  ]
);
EOF

echo ""
echo "✅ Após executar, todos os scopes estarão no formato padronizado:"
echo "   - geo:read → geo"
echo "   - geo:write → geo"
echo "   - cep → cep"
echo "   - cnpj → cnpj"
echo "   - all → all"

