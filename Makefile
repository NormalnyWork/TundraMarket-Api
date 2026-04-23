SHELL := /bin/bash

-include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

.DEFAULT_GOAL := help

MIGRATIONS_DIR ?= migrations
DB_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

.PHONY: help
help: ## Показать список команд
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | \
    sort | \
    awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'

## go install goose
## go install github.com/pressly/goose/v3/cmd/goose@latest
# ==============================================================================
# Migrations

.PHONY: migrate-up
migrate-up: ## Применить все миграции
	goose -dir ./$(MIGRATIONS_DIR) postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down: ## Откатить последнюю миграцию
	goose -dir ./$(MIGRATIONS_DIR) postgres "$(DB_URL)" down

.PHONY: generate-migration
generate-migration: ## Создать миграцию (make generate-migration name=...)
	@test -n "$(name)" || (echo "ERROR: name is required. Usage: make generate-migration name=add_users" && exit 1)
	goose -dir ./$(MIGRATIONS_DIR) create $(name) sql
