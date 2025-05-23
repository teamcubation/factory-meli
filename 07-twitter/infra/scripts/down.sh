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

# Pergunta se o usuário deseja remover TUDO (containers, volumes, imagens, redes)
echo -e "${CYAN}Deseja remover também TODOS os volumes, imagens e redes do projeto? (y/n)${NC}"
read -r resposta

if [[ "$resposta" =~ ^[Yy]$ ]]; then
    echo -e "${CYAN}🧹 Removendo todos os containers, volumes, imagens e redes do projeto...${NC}"
    if docker compose down -v --rmi all --remove-orphans; then
        echo -e "${GREEN}✅ Todos os containers, volumes, imagens e redes do projeto foram removidos com sucesso!${NC}"
    else
        echo -e "${RED}❌ Falha ao remover containers, volumes, imagens ou redes.${NC}"
        exit 1
    fi
    exit 0
fi

# Obtendo os IDs dos containers em execução
containers=$(docker ps -q)

# Verificando se há containers em execução
if [ -z "$containers" ]; then
    echo -e "${CYAN}ℹ️  Nenhum container em execução no momento.${NC}"
    exit 1
fi

# Verifica se o usuário passou argumentos (nomes de containers)
if [ "$#" -gt 0 ]; then
    echo -e "${CYAN}🛑 Parando os containers especificados...${NC}"
    
    for container in "$@"; do
        # Parando cada container específico
        if docker stop "$container"; then
            echo -e "${GREEN}✅ Container '$container' parado com sucesso!${NC}"
        else
            echo -e "${RED}❌ Falha ao parar o container '$container'.${NC}"
        fi
    done
else
    echo -e "${CYAN}🛑 Parando todos os containers...${NC}"
    
    # Parando todos os containers
    if docker compose down -v; then
        echo -e "${GREEN}✅ Todos os containers foram parados e removidos com sucesso!${NC}"
    else
        echo -e "${RED}❌ Falha ao parar e remover os containers.${NC}"
        exit 1
    fi
fi