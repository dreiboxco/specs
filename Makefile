.PHONY: build build-release test clean version help

# Variáveis
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")
BINARY_NAME := specs
BUILD_DIR := bin
MAIN_PATH := ./cmd/specs

help: ## Exibe ajuda
	@echo "Targets disponíveis:"
	@echo "  build          - Build local para desenvolvimento"
	@echo "  build-release  - Build para release (incrementa versão e cria tag)"
	@echo "  test           - Executa testes"
	@echo "  clean          - Remove arquivos de build"
	@echo "  version        - Exibe versão atual"

build: ## Build local para desenvolvimento
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build concluído: $(BUILD_DIR)/$(BINARY_NAME)"

build-release: ## Build para release com incremento automático
	@bash scripts/build-release.sh

test: ## Executa testes
	@echo "Executando testes..."
	@go test -v ./...

clean: ## Remove arquivos de build
	@echo "Limpando arquivos de build..."
	@rm -rf $(BUILD_DIR)
	@go clean

version: ## Exibe versão atual
	@echo $(VERSION)

