COMPOSE := docker compose
QUEUE   ?= messages
ADMIN   := $(COMPOSE) exec -T rabbitmq rabbitmqadmin --username app --password app --non-interactive

# Toolchain Go roda em container: nao precisa de Go instalado na maquina.
# --network host: o app alcanca o RabbitMQ em localhost:5672, igual ao .env
# -u: arquivos gerados saem com o seu usuario, nao root
# GOPATH/GOCACHE em .cache/: cache persiste entre runs, sem volume nomeado
DOCKER_GO := docker run --rm -t --network host \
		-v "$(PWD)":/app -w /app \
		-u $(shell id -u):$(shell id -g) \
		-e GOPATH=/app/.cache/gopath -e GOCACHE=/app/.cache/gobuild
GO        := $(DOCKER_GO) golang:1.26 go
GO_LAMBDA := $(DOCKER_GO) -e GOOS=linux -e GOARCH=arm64 -e CGO_ENABLED=0 golang:1.26 go

.DEFAULT_GOAL := help

help: ## Lista os comandos
	@grep -hE '^[a-z][a-z-]*:.*##' $(MAKEFILE_LIST) | sed -E 's/:[^#]*## /\t/' | expand -t22

# --- ambiente ---------------------------------------------------------------

up: ## Sobe o RabbitMQ e declara a fila
	$(COMPOSE) up -d --wait
	$(ADMIN) declare queue --name $(QUEUE) --durable true

down: ## Derruba o RabbitMQ (apaga os dados)
	$(COMPOSE) down

logs: ## Acompanha os logs do RabbitMQ
	$(COMPOSE) logs -f rabbitmq

publish: ## Publica payload.json na fila
	$(ADMIN) publish message --routing-key $(QUEUE) --payload-file - < payload.json

peek: ## Mostra quantas mensagens estao na fila
	$(COMPOSE) exec -T rabbitmq rabbitmqctl list_queues name messages

# --- go (tudo em container) -------------------------------------------------

run: ## Roda o dispatcher local (consumer RabbitMQ)
	$(GO) run ./cmd

test: ## Testes unitarios (sem rede, sem e-mail real)
	$(GO) test -short ./...

test-integration: ## Testes de integracao (exige 'make up' + envia e-mail de verdade)
	$(GO) test ./...

vet: ## Analise estatica
	$(GO) vet ./...

tidy: ## Arruma o go.mod
	$(GO) mod tidy

deps: ## Baixa as dependencias
	$(GO) mod download

build: ## Compila o binario local
	$(GO) build -o bin/dispatcher ./cmd

build-lambda: ## Compila o bootstrap do Lambda (arm64/Graviton) - Fase 2
	$(GO_LAMBDA) build -tags lambda.norpc -o bin/bootstrap ./cmd/lambda

shell: ## Abre um shell no container Go
	$(DOCKER_GO) -i golang:1.26 bash

clean: ## Limpa binarios e cache de build
	rm -rf bin .cache

.PHONY: help up down logs publish peek run test test-integration vet tidy deps build build-lambda shell clean
