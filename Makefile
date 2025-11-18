.PHONY: help build run clean test deps install

# Variáveis
BINARY_NAME=wsicrmrest
MAIN_PATH=./cmd/server
BUILD_DIR=build

# Variáveis de versão
VERSION_PACKAGE=wsicrmrest/internal/config
BUILD_TIME=$(shell date '+%Y-%m-%dT%H:%M:%S')
VERSION_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS=-ldflags "-X '$(VERSION_PACKAGE).VersionDate=$(VERSION_DATE)' -X '$(VERSION_PACKAGE).BuildTime=$(BUILD_TIME)'"

help: ## Mostra esta ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Baixa as dependências
	@echo "Baixando dependências..."
	go mod download
	go mod tidy
	@echo "✓ Dependências instaladas"

build: ## Compila o projeto para Linux
	@echo "Compilando $(BINARY_NAME) para Linux..."
	@echo "Data/Hora da compilação: $(BUILD_TIME)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Compilação concluída: $(BUILD_DIR)/$(BINARY_NAME)"

build-windows: build-windows-32 build-windows-64 ## Compila para Windows 32 e 64 bits

build-windows-32: ## Compila o projeto para Windows 32 bits
	@echo "Compilando $(BINARY_NAME) para Windows 32 bits..."
	@echo "Data/Hora da compilação: $(BUILD_TIME)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_win32.exe $(MAIN_PATH)
	@echo "✓ Compilação concluída: $(BUILD_DIR)/$(BINARY_NAME)_win32.exe"

build-windows-64: ## Compila o projeto para Windows 64 bits
	@echo "Compilando $(BINARY_NAME) para Windows 64 bits..."
	@echo "Data/Hora da compilação: $(BUILD_TIME)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_win64.exe $(MAIN_PATH)
	@echo "✓ Compilação concluída: $(BUILD_DIR)/$(BINARY_NAME)_win64.exe"

run: ## Executa o servidor
	@if [ ! -f dbinit.ini ]; then \
		echo "❌ Erro: arquivo dbinit.ini não encontrado!"; \
		echo "Copie dbinit.ini.example para dbinit.ini e configure."; \
		exit 1; \
	fi
	@echo "Iniciando $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

dev: build run ## Compila e executa em modo desenvolvimento

clean: ## Remove arquivos de build
	@echo "Limpando arquivos de build..."
	rm -rf $(BUILD_DIR)
	rm -rf log/*.log
	@echo "✓ Limpeza concluída"

test: ## Executa testes
	@echo "Executando testes..."
	go test -v ./...

test-api: ## Testa as APIs (requer servidor rodando)
	@if [ ! -f scripts/test_apis.sh ]; then \
		echo "❌ Erro: scripts/test_apis.sh não encontrado!"; \
		exit 1; \
	fi
	@echo "Testando APIs..."
	@chmod +x scripts/test_apis.sh
	@./scripts/test_apis.sh

install: deps build ## Instala dependências e compila

fmt: ## Formata o código
	@echo "Formatando código..."
	go fmt ./...
	@echo "✓ Código formatado"

vet: ## Executa go vet
	@echo "Executando go vet..."
	go vet ./...
	@echo "✓ Verificação concluída"

check: fmt vet ## Formata e verifica o código

docker-build: ## Cria imagem Docker (TODO)
	@echo "TODO: Implementar build Docker"

setup: ## Configuração inicial do projeto
	@echo "Configurando projeto..."
	@if [ ! -f dbinit.ini ]; then \
		cp dbinit.ini.example dbinit.ini; \
		echo "✓ Arquivo dbinit.ini criado. Configure suas credenciais!"; \
	else \
		echo "ℹ dbinit.ini já existe"; \
	fi
	@mkdir -p log
	@mkdir -p docs
	@mkdir -p scripts
	@mkdir -p $(BUILD_DIR)
	@echo "✓ Diretórios criados"
	@$(MAKE) deps
	@echo ""
	@echo "✓ Setup concluído!"
	@echo "Próximos passos:"
	@echo "  1. Configure dbinit.ini com suas credenciais"
	@echo "  2. Execute: make build (Linux) ou make build-windows (Windows)"
	@echo "  3. Execute: make run"

.DEFAULT_GOAL := help
