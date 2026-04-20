include .env
export

export PROJECT_ROOT=$(shell pwd)


env-up:
	@docker compose up -d todolist-postgres

env-down:
	@docker compose down todolist-postgres

env-cleanup:
	@read -p "Clear current DB session? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todolist-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "DB session cleared."; \
	else \
		echo "DB session preserved."; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-down:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Environment variable SEQ is not set. Example: make migrate-create seq=1"; \
		exit 1; \
	fi; \
	docker compose run --rm todolist-postgres-migrate create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: variable 'version' is required. Example: make migrate-force version=0"; \
		exit 1; \
	fi; \
	docker compose run --rm todolist-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todolist-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		force "$(version)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Environment variable action is not set. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todolist-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todolist-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

logs-cleanup:
	@read -p "Очистить все log файлы? Опасность утери логов.  [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены."; \
	else \
		echo "Очистка логов отменена."; \
	fi

todo-run:
	@export LOGGER_FOLDER=$(PROJECT_ROOT)/out/logs && \
	export POSTGRES_HOST=localhost && \
	go run ${PROJECT_ROOT}/cmd/todoApp/main.go

todo-deploy:
	@docker compose up -d --build todoapp

ps:
	@docker compose ps