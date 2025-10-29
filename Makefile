.PHONY: help admin up down logs build

help: ## Mostrar ajuda
	@echo "🛠️  Retech Core - Comandos Disponíveis:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

admin: ## Criar Super Admin (alanrezendeee@gmail.com / admin123456)
	@./scripts/create-local-admin.sh

up: ## Subir containers (API + MongoDB)
	@docker-compose -f build/docker-compose.yml up -d
	@echo "✅ Containers iniciados!"
	@echo "📡 API: http://localhost:8080"
	@echo "🗄️  MongoDB: localhost:27017"

down: ## Parar containers
	@docker-compose -f build/docker-compose.yml down
	@echo "✅ Containers parados!"

logs: ## Ver logs da API
	@docker-compose -f build/docker-compose.yml logs -f api

logs-mongo: ## Ver logs do MongoDB
	@docker-compose -f build/docker-compose.yml logs -f mongo

build: ## Rebuild da API
	@docker-compose -f build/docker-compose.yml up --build -d
	@echo "✅ API rebuilded!"

restart: ## Reiniciar containers
	@docker-compose -f build/docker-compose.yml restart
	@echo "✅ Containers reiniciados!"

ps: ## Ver status dos containers
	@docker-compose -f build/docker-compose.yml ps

shell-api: ## Shell no container da API
	@docker exec -it build-api-1 sh

shell-mongo: ## Shell no MongoDB
	@docker exec -it build-mongo-1 mongosh retech_core

clean: ## Limpar volumes e containers
	@docker-compose -f build/docker-compose.yml down -v
	@echo "✅ Containers e volumes removidos!"

setup: up admin ## Setup completo (up + criar admin)
	@echo ""
	@echo "✅ Setup completo!"
	@echo "🌐 Acesse: http://localhost:3001/admin/login"
	@echo "📧 Email: alanrezendeee@gmail.com"
	@echo "🔑 Senha: admin123456"
