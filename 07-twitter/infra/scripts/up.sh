#!/bin/bash

# Cores
GREEN='\033[1;32m'
RED='\033[1;31m'
CYAN='\033[1;36m'
NC='\033[0m' # Sem cor

echo -e "${CYAN}🔍 Verificando dependências...${NC}"

# Verificando se o docker compose está instalado
if ! command -v docker compose &> /dev/null
then
    echo -e "${RED}❌ docker compose não está instalado ou não está no PATH.${NC}"
    exit 1
fi

echo -e "${CYAN}🚀 Subindo os containers com docker compose...${NC}"

# Executando os containers do docker-compose.yml
if docker compose up -d; then
    echo -e "${GREEN}✅ Containers iniciados com sucesso!${NC}"
else
    echo -e "${RED}❌ Falha ao subir os containers.${NC}"
    exit 1
fi