.PHONY: help up down restart rebuild build logs ps clean test migrate

help: ## Показать список команд
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

up: ## Собрать и запустить все сервисы
	docker compose up -d --build

down: ## Остановить и удалить контейнеры
	docker compose down

restart: ## Перезапустить сервисы
	docker compose restart

rebuild: ## Пересобрать образы и перезапустить контейнеры
	docker compose up -d --build --force-recreate

build: ## Собрать образы
	docker compose build

logs: ## Логи всех сервисов (follow)
	docker compose logs -f

ps: ## Список контейнеров проекта
	docker compose ps

clean: ## Остановить всё и удалить тома (данные БД удалятся)
	docker compose down -v

test: ## Прогнать go-тесты в контейнере app
	docker compose run --rm app go test ./...

migrate: ## Применить миграции goose (migrations/)
	docker compose run --rm app sh -c 'goose -dir migrations postgres "$$DATABASE_URL" up'
