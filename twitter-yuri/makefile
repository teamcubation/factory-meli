# Variáveis
DOCKER_COMPOSE = docker compose

# PHONY targets
.PHONY: up down logs clean help

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

clean:
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f
help:
	@echo "Comandos disponíveis:"
	@echo "  make up     	 - Inicia os containers"
	@echo "  make down   	 - Para os containers
	@echo "  make logs       - Mostra logs dos containers"
	@echo "  make clean      - Remove containers, volumes e cache"
	@echo "  make help       - Mostra esta mensagem de ajuda"