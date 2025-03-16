-include .env
export

.DEFAULT_GOAL := help

.PHONY: help
help: ## Available commands
	@clear
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[0;33m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""


##@ Services

.PHONY: run-api
run-api: ## Run API server
	go run cmd/api/main.go

.PHONY: run-worker
run-worker: ## Run background worker
	go run cmd/worker/main.go

.PHONY: run-scheduler
run-scheduler: ## Run task scheduler
	go run cmd/scheduler/main.go

.PHONY: run-all
run-all: ## Run all services
	make run-api & make run-worker & make run-scheduler


##@ Development

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: lint
lint: ## Run linter (golangci-lint)
	golangci-lint run ./...

.PHONY: format
format: ## Format code
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -l -w .

.PHONY: seed
seed: ## Initialize database and RabbitMQ
	go run cmd/seed/main.go

.PHONY: setup-local
setup-local: ## Setup database and RabbitMQ locally
	./scripts/setup.sh

.PHONY: setup-and-run
setup-and-run: setup-local ## Initialize database and run all services
	make run-all

.PHONY: setup-and-run-api
setup-and-run-api: setup-local ## Initialize database and run API server
	make run-api


##@ Aliases

.PHONY: r
r: ## Run all services
	@make run-all

.PHONY: t
t: ## Run tests
	@make test

.PHONY: l
l: ## Run linter (golangci-lint)
	@make lint

.PHONY: f
f: ## Format code
	@make format

.PHONY: s
s: ## Initialize database and RabbitMQ
	@make seed

.PHONY: sl
sl: ## Setup database and RabbitMQ locally
	@make setup-local

.PHONY: sr
sr: ## Initialize database and run all services
	@make setup-and-run

.PHONY: sa
sa: ## Initialize database and run API server
	@make setup-and-run-api

