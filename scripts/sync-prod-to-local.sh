#!/bin/bash
# Script para copiar MongoDB de produção para local

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔄 Sincronizando MongoDB de Produção → Local${NC}\n"

# 1. Validar variáveis
if [ -z "$MONGO_URI_PROD" ]; then
  echo -e "${RED}❌ MONGO_URI_PROD não definida!${NC}"
  echo "Use: export MONGO_URI_PROD='mongodb+srv://...'"
  exit 1
fi

# 2. Fazer backup
echo -e "${YELLOW}📦 Fazendo dump de produção...${NC}"
rm -rf /tmp/mongo-backup-prod
mongodump --uri="$MONGO_URI_PROD" --out=/tmp/mongo-backup-prod

if [ $? -ne 0 ]; then
  echo -e "${RED}❌ Erro ao fazer dump!${NC}"
  exit 1
fi

# 3. Restaurar localmente
echo -e "${YELLOW}📥 Restaurando no MongoDB local...${NC}"
mongorestore --uri="mongodb://localhost:27017" \
  --db=retech_core \
  --drop \
  /tmp/mongo-backup-prod/retech_core

if [ $? -ne 0 ]; then
  echo -e "${RED}❌ Erro ao restaurar!${NC}"
  exit 1
fi

# 4. Limpar dados sensíveis
echo -e "${YELLOW}🧹 Limpando dados sensíveis...${NC}"
mongosh retech_core --quiet --eval "
  // Remover API keys de produção (manter só playground)
  db.api_keys.deleteMany({ ownerId: { \$ne: 'playground-public' } });
  
  // Limpar refresh tokens
  db.users.updateMany({}, { \$unset: { refreshToken: 1 } });
  
  print('✅ Limpeza concluída!');
"

# 5. Limpar backup
echo -e "${YELLOW}🗑️  Removendo arquivos temporários...${NC}"
rm -rf /tmp/mongo-backup-prod

echo -e "\n${GREEN}✅ Sincronização completa!${NC}"
echo -e "MongoDB local está atualizado com dados de produção (sem dados sensíveis)\n"

