# Versionamento da API WSICRMREST

**Data de criação:** 2025-11-17
**Última atualização:** 2025-11-17

## Localização da Versão

A versão da aplicação é definida através de variáveis globais em:

```
internal/config/config.go
```

## Variáveis de Versão

```go
var (
    // Version - Versão do sistema (atualizado manualmente)
    Version = "Versão 3.0.0.1 (GO)"

    // VersionDate - Data da versão (injetado automaticamente durante a compilação)
    VersionDate = "development"

    // BuildTime - Data e hora da compilação (injetado automaticamente)
    BuildTime = "unknown"
)
```

## ⚡ Novidade: Data/Hora Automática na Compilação

A partir de agora, **não é mais necessário atualizar manualmente** `VersionDate` e `BuildTime`. Esses valores são injetados automaticamente durante a compilação usando **ldflags** do Go.

## Como Atualizar a Versão

### Passo 1: Atualizar o Número da Versão (Manual)

Edite apenas o campo `Version` no arquivo `internal/config/config.go`:

```go
// Antes
Version = "Versão 3.0.0.1 (GO)"

// Depois
Version = "Versão 3.0.0.2 (GO)"
```

**IMPORTANTE:** Não modifique `VersionDate` ou `BuildTime` manualmente! Eles serão preenchidos automaticamente na compilação.

### Passo 2: Compilar com Makefile

Use o Makefile para compilar. Ele injeta automaticamente a data/hora:

```bash
# Linux
make build

# Windows 32 bits
make build-windows-32

# Windows 64 bits
make build-windows-64

# Ambas as versões Windows
make build-windows
```

**Saída da compilação:**
```
Compilando wsicrmrest para Linux...
Data/Hora da compilação: 2025-11-17T17:51:37
go build -ldflags "-X 'wsicrmrest/internal/config.VersionDate=2025-11-17T17:51:37' -X 'wsicrmrest/internal/config.BuildTime=2025-11-17T17:51:37'" -o build/wsicrmrest ./cmd/server
✓ Compilação concluída: build/wsicrmrest
```

### Passo 3: Verificar Injeção

Verifique se a data/hora foi injetada corretamente no binário:

```bash
# Linux
strings build/wsicrmrest | grep "2025-"

# Windows (PowerShell)
Select-String -Path build\wsicrmrest_win64.exe -Pattern "2025-"
```

### Passo 4: Reiniciar o Servidor

Reinicie a aplicação para usar a nova versão:

```bash
# Linux
./build/wsicrmrest

# Windows
.\build\wsicrmrest_win64.exe
```

## Formato da Versão

### Version (Manual) - gsVersao

Formato recomendado: `Versão X.Y.Z.W (GO)`

- **Versão**: Prefixo fixo
- **X**: Major - Mudanças incompatíveis na API
- **Y**: Minor - Novas funcionalidades compatíveis
- **Z**: Patch - Correções de bugs
- **W**: Build - Número incremental de build
- **(GO)**: Indicador de implementação em Go

Exemplos:
- `Versão 3.0.0.1 (GO)` - Versão 3.0, inicial, build 1
- `Versão 3.0.0.2 (GO)` - Versão 3.0, inicial, build 2
- `Versão 3.0.1.0 (GO)` - Versão 3.0, patch 1
- `Versão 3.1.0.0 (GO)` - Versão 3.1, nova funcionalidade

### VersionDate e BuildTime (Automáticos) - gsDataVersao

Formato: ISO 8601 em hora local (YYYY-MM-DDTHH:MM:SS)

- **YYYY**: Ano (4 dígitos)
- **MM**: Mês (2 dígitos)
- **DD**: Dia (2 dígitos)
- **HH**: Hora (2 dígitos, formato 24h, hora local do sistema)
- **MM**: Minuto (2 dígitos)
- **SS**: Segundo (2 dígitos)

Exemplos (injetados automaticamente):
- `2025-11-17T17:46:29` - 17 de novembro de 2025, 17h46m29s (hora local)
- `2025-11-18T14:15:42` - 18 de novembro de 2025, 14h15m42s (hora local)

## Onde a Versão Aparece

### 1. Logs de Inicialização

```json
{
  "level": "INFO",
  "message": "Iniciando WSICRMREST",
  "version": "Versão 3.0.0.1 (GO)",
  "version_date": "2025-11-17T17:51:37",
  "build_time": "2025-11-17T17:51:37"
}
```

### 2. Endpoint de Teste (GET /connect/v1/wsteste)

```json
{
  "code": "000",
  "organizadorCodigo": 1,
  "organizadorNome": "Empresa Exemplo",
  "organizadorCnpj": "12345678000190",
  "organizadorLojaMatriz": 1,
  "organizadorCodIsga": 123,
  "versao": "Versão 3.0.0.1 (GO)",
  "versaoData": "2025-11-17T17:51:37"
}
```

### 3. Logs de Requisições (Tabela WSREQUISICOES)

Cada requisição registrada no banco inclui o campo `WSVERSAO` com a versão atual.

### 4. Documentação Swagger

A versão aparece no topo da documentação Swagger:
```
http://localhost:8080/swagger/index.html
```

## 🌍 Fuso Horário

A data/hora de compilação usa a **hora local do sistema** onde a compilação é executada:

```makefile
# Hora local (padrão)
BUILD_TIME=$(shell date '+%Y-%m-%dT%H:%M:%S')
VERSION_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')
```

Para usar UTC, adicione o flag `-u` no Makefile:

```makefile
# UTC (opcional)
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%S')
VERSION_DATE=$(shell date -u '+%Y-%m-%dT%H:%M:%S')
```

## Por Que Não Usar dbinit.ini?

A versão é definida como variável global (e não no arquivo `dbinit.ini`) pelos seguintes motivos:

1. **Imutabilidade**: A versão é parte integral do código, não uma configuração de ambiente
2. **Rastreabilidade**: A versão fica versionada junto com o código no Git
3. **Build Único**: Cada build possui uma versão específica, independente do ambiente
4. **Compatibilidade WinDev**: Mantém compatibilidade com as variáveis globais `gsVersao` e `gsDataVersao` do sistema original

## Automatização (Opcional)

### Usando Build Tags

Você pode automatizar a versão durante o build usando ldflags:

```bash
# Build com versão injetada
VERSION="Ver 1.27.0.1"
DATE=$(date +"%Y-%m-%dT%H:%M:%S")

go build -ldflags="-X 'wsicrmrest/internal/config.Version=$VERSION' -X 'wsicrmrest/internal/config.VersionDate=$DATE'" -o bin/wsicrmrest ./cmd/server
```

### Usando Makefile

Adicione ao `Makefile`:

```makefile
# Obter versão do Git
GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "Ver 1.0.0.0")
BUILD_DATE := $(shell date +"%Y-%m-%dT%H:%M:%S")

build:
	go build -ldflags="-X 'wsicrmrest/internal/config.Version=$(GIT_VERSION)' -X 'wsicrmrest/internal/config.VersionDate=$(BUILD_DATE)'" -o bin/wsicrmrest ./cmd/server
```

### Usando CI/CD

Em pipelines de CI/CD, você pode injetar a versão automaticamente:

```yaml
# GitHub Actions exemplo
- name: Build
  run: |
    VERSION="Ver ${{ github.ref_name }}"
    DATE=$(date +"%Y-%m-%dT%H:%M:%S")
    go build -ldflags="-X 'wsicrmrest/internal/config.Version=$VERSION' -X 'wsicrmrest/internal/config.VersionDate=$DATE'" -o bin/wsicrmrest ./cmd/server
```

## Verificando a Versão Atual

### Via Logs

Ao iniciar o servidor, a versão é exibida nos logs:
```bash
./bin/wsicrmrest
```

### Via API

Consulte o endpoint de teste:
```bash
curl http://localhost:8080/connect/v1/wsteste | jq '.versao, .versaoData'
```

### Via Código

Importe e use as variáveis:
```go
import "wsicrmrest/internal/config"

fmt.Println("Version:", config.Version)
fmt.Println("Date:", config.VersionDate)
```

## Boas Práticas

1. **Sempre atualize ambas**: Quando mudar `Version`, atualize também `VersionDate`
2. **Use datas reais**: A data deve refletir quando a versão foi criada
3. **Documente mudanças**: Mantenha um CHANGELOG.md com as alterações de cada versão
4. **Tag no Git**: Crie tags Git correspondentes às versões principais
5. **Incremente corretamente**:
   - Major: Mudanças incompatíveis
   - Minor: Novas funcionalidades
   - Patch: Correções de bugs
   - Build: Builds incrementais

## Exemplo de Workflow

```bash
# 1. Fazer alterações no código
git add .
git commit -m "feat: adicionar novo endpoint"

# 2. Atualizar versão
# Editar internal/config/config.go
# Version = "Ver 1.27.0.1"
# VersionDate = "2025-11-01T10:00:00"

# 3. Comitar mudança de versão
git add internal/config/config.go
git commit -m "chore: bump version to 1.27.0.1"

# 4. Criar tag Git
git tag -a v1.27.0.1 -m "Release version 1.27.0.1"

# 5. Push com tags
git push origin main --tags

# 6. Recompilar
make build

# 7. Deploy
./bin/wsicrmrest
```

## Referências

- [Semantic Versioning](https://semver.org/)
- [ISO 8601 Date Format](https://en.wikipedia.org/wiki/ISO_8601)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
