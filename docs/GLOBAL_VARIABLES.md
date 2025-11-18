# Variáveis Globais - WSICRMREST

Este documento detalha todas as variáveis globais do sistema convertidas do WinDev para Go.

## 📋 Mapeamento WinDev → Go

### Variáveis JWT e Segurança

| WinDev | Go Config | Tipo | Valor Padrão | Descrição |
|--------|-----------|------|--------------|-----------|
| `gsKey` | `config.JWT.SecretKey` | string | `CloudI0812IcrMmDB` | Chave secreta para HMAC-SHA256 do JWT |
| `gsIss` | `config.JWT.Issuer` | string | `WSCloudICrmIntellsys` | Issuer do token JWT |
| `gsKeyDelivery` | `config.JWT.KeyDelivery` | string | `Ped2505IcrM` | Chave adicional para delivery |
| `gnFusoHorario` | `config.JWT.Timezone` | int | `0` | Fuso horário (0=UTC, -3=Brasília) |

### Variáveis de Versão

| WinDev | Go Config | Tipo | Valor Padrão | Descrição |
|--------|-----------|------|--------------|-----------|
| `gsVersao` | `config.Application.Version` | string | `Ver 1.26.4.27` | Versão do sistema |
| `gsDataVersao` | `config.Application.VersionDate` | string | `2025-10-16T11:55:00` | Data da versão |

### Variáveis de Organização

| WinDev | Go Config | Tipo | Descrição |
|--------|-----------|------|-----------|
| `gnOrgCodigo` | `config.Organization.Codigo` | int | Código da organização |
| `gsOrgNome` | `config.Organization.Nome` | string | Nome da organização |
| `gnOrgCnpj` | `config.Organization.CNPJ` | string | CNPJ da organização |
| `gnOrgLojMatriz` | `config.Organization.LojaMatriz` | int | Código da loja matriz |
| `gnOrgCodISGA` | `config.Organization.CodISGA` | int | Código ISGA |
| `gnRegModulos` | `config.Organization.RegModulos` | int | Registro de módulos (padrão: 1) |
| `gnOrgFormaLimite` | `config.Organization.FormaLimite` | int | Forma de cálculo do limite |
| `gnOrgCalcDispFuturoCartao` | `config.Organization.CalcDispFuturoCartao` | int | Cálculo disponível futuro cartão |
| `gnOrgCalcDispFuturoConvenio` | `config.Organization.CalcDispFuturoConvenio` | int | Cálculo disponível futuro convênio |
| `gnOrgDiaVectoGrupo1-6` | `config.Organization.DiaVectoGrupo1-6` | int | Dia de vencimento por grupo (1 a 6) |
| `gnOrgDiaCorteGrupo1-6` | `config.Organization.DiaCorteGrupo1-6` | int | Dia de corte por grupo (1 a 6) |

**Importante:** Os dados da organização são carregados automaticamente da tabela `ORGANIZADOR` na inicialização do sistema.

### Variáveis de Contexto de Requisição

| WinDev | Go Context | Tipo | Descrição |
|--------|------------|------|-----------|
| `gsClient_Id` | `RequestContext.ClientID` | string | Client ID da requisição atual |
| `gsNomeAplicacao` | `RequestContext.NomeAplicacao` | string | Nome da aplicação |
| `gsWSLogReqUUID` | `RequestContext.UUID` | string | UUID único da requisição |
| `gdtInicio` | `RequestContext.StartTime` | time.Time | Data/hora de início da requisição |
| `gsDetalheLogApi` | `RequestContext.DetalheLogAPI` | string | Detalhes do log |

---

## 🗄️ Tabela ORGANIZADOR

### Carregamento Automático

Os dados da organização são carregados automaticamente da tabela `ORGANIZADOR` durante a inicialização do sistema através da função `pgLeOrganizador()` (convertida para `database.LeOrganizador()`).

**Processo:**
1. Sistema conecta ao banco de dados
2. Executa query: `SELECT ... FROM ORGANIZADOR WHERE OrgCodigo > 0`
3. Carrega primeiro registro encontrado
4. Atualiza objeto `config.Organization` com dados do banco
5. Dados do banco sobrescrevem valores do `dbinit.ini`

**Se a tabela não existir ou estiver vazia:**
- Sistema emite erro no log
- **Sistema NÃO inicia** (erro fatal)
- É obrigatório ter a tabela ORGANIZADOR com pelo menos um registro

**⚠️ IMPORTANTE:** Diferente de outras configurações, os dados da organização **NÃO** vêm do `dbinit.ini`. Eles são **exclusivamente** carregados da tabela `ORGANIZADOR` no banco de dados.

### Estrutura da Tabela

```sql
CREATE TABLE ORGANIZADOR (
    ORGCODIGO                  NUMBER PRIMARY KEY,
    ORGNOME                    VARCHAR2(200) NOT NULL,
    ORGCNPJ                    VARCHAR2(20),
    ORGCODLOJAMATRIZ           NUMBER,
    ORGFORMALIMITE             NUMBER DEFAULT 0,
    ORGCALCDISPFUTUROCARTAO    NUMBER DEFAULT 0,
    ORGCALCDISPFUTUROCONVENIO  NUMBER DEFAULT 0,
    ORGCODISGA                 NUMBER,
    ORGDIAFATGRUPO1            NUMBER DEFAULT 0,
    ORGDIAFATGRUPO2            NUMBER DEFAULT 0,
    ORGDIAFATGRUPO3            NUMBER DEFAULT 0,
    ORGDIAFATGRUPO4            NUMBER DEFAULT 0,
    ORGDIAFATGRUPO5            NUMBER DEFAULT 0,
    ORGDIAFATGRUPO6            NUMBER DEFAULT 0,
    ORGDIACORGRUPO1            NUMBER DEFAULT 0,
    ORGDIACORGRUPO2            NUMBER DEFAULT 0,
    ORGDIACORGRUPO3            NUMBER DEFAULT 0,
    ORGDIACORGRUPO4            NUMBER DEFAULT 0,
    ORGDIACORGRUPO5            NUMBER DEFAULT 0,
    ORGDIACORGRUPO6            NUMBER DEFAULT 0
);
```

### Exemplo de Dados

```sql
INSERT INTO ORGANIZADOR (
    ORGCODIGO,
    ORGNOME,
    ORGCNPJ,
    ORGCODLOJAMATRIZ,
    ORGCODISGA,
    ORGDIAFATGRUPO1,
    ORGDIAFATGRUPO2,
    ORGDIACORGRUPO1,
    ORGDIACORGRUPO2
) VALUES (
    1,                      -- Código
    'Minha Empresa',        -- Nome
    '12345678000199',       -- CNPJ
    1,                      -- Loja Matriz
    1001,                   -- Código ISGA
    10,                     -- Vencimento Grupo 1 (dia 10)
    25,                     -- Vencimento Grupo 2 (dia 25)
    5,                      -- Corte Grupo 1 (dia 5)
    20                      -- Corte Grupo 2 (dia 20)
);
```

### Campos Detalhados

**Grupos de Faturamento:**
- `ORGDIAFATGRUPO1-6`: Dia do mês para vencimento de cada grupo (1-31)
- `ORGDIACORGRUPO1-6`: Dia do mês para corte de cada grupo (1-31)

**Exemplo de Uso:**
- Grupo 1: Clientes com vencimento no dia 10, corte no dia 5
- Grupo 2: Clientes com vencimento no dia 25, corte no dia 20
- E assim por diante até o Grupo 6

---

## 🔧 Configuração via dbinit.ini

### ~~Seção [jwt]~~ (NÃO UTILIZADA)

**⚠️ A seção `[jwt]` foi REMOVIDA do `dbinit.ini`**

As credenciais JWT são **variáveis globais hardcoded** no código, conforme o WinDev original:

```go
// Valores fixos (não configuráveis)
SecretKey:   "CloudI0812IcrMmDB"      // gsKey
Issuer:      "WSCloudICrmIntellsys"   // gsIss
KeyDelivery: "Ped2505IcrM"            // gsKeyDelivery
Timezone:    0                         // gnFusoHorario
```

**Não configure credenciais JWT no `dbinit.ini` - elas são fixas no código!**

Para alterar esses valores, você deve modificar o arquivo `internal/config/config.go`.

### Seção [application]

```ini
[application]
# gsVersao - Versão do sistema
version = Ver 1.26.4.27

# gsDataVersao - Data da versão
version_date = 2025-10-16T11:55:00

# Outras configurações
environment = production
port = 8080
log_dir = log
ws_grava_log_db = true
ws_detalhe_log_api = false
```

### ~~Seção [organization]~~ (NÃO UTILIZADA)

**⚠️ A seção `[organization]` foi REMOVIDA do `dbinit.ini`**

Os dados da organização são carregados **exclusivamente** da tabela `ORGANIZADOR` no banco de dados através da função `pgLeOrganizador()`.

**Não configure dados de organização no `dbinit.ini` - eles serão ignorados!**

---

## 💻 Uso no Código

### Acessando Configurações

```go
// Em um handler
func MyHandler(cfg *config.Config, ...) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Acessar variáveis JWT
        secretKey := cfg.JWT.SecretKey      // gsKey
        issuer := cfg.JWT.Issuer            // gsIss
        timezone := cfg.JWT.Timezone        // gnFusoHorario

        // Acessar variáveis de versão
        version := cfg.Application.Version  // gsVersao
        versionDate := cfg.Application.VersionDate // gsDataVersao

        // Acessar variáveis de organização
        orgCodigo := cfg.Organization.Codigo // gnOrgCodigo
        orgNome := cfg.Organization.Nome     // gsOrgNome

        // ... resto do código
    }
}
```

### Usando Contexto de Requisição

```go
// Criar contexto no início da requisição
reqCtx := reqcontext.NewRequestContext()

// Definir Client ID e Nome da Aplicação
reqCtx.SetClientInfo(clientID, nomeAplicacao)

// Adicionar detalhes ao log (se habilitado)
reqCtx.AddLogDetail("Processando requisição...")

// Obter duração
duration := reqCtx.GetDuration()

// Acessar UUID
uuid := reqCtx.UUID

// Acessar tempo de início
startTime := reqCtx.StartTime
```

---

## 🔐 Segurança das Variáveis

### Variáveis Sensíveis

As seguintes variáveis contêm informações sensíveis e **NÃO** devem ser commitadas no Git:

- ✅ **dbinit.ini** (está no .gitignore)
- ❌ **dbinit.ini.example** (pode ser commitado - contém valores de exemplo)

### Chaves de Segurança

#### gsKey (SecretKey)
- **Uso:** Assinatura HMAC-SHA256 dos tokens JWT
- **Requisito:** Mínimo 16 caracteres (recomendado 32+)
- **Padrão:** `CloudI0812IcrMmDB`
- **Produção:** ⚠️ **TROCAR** por uma chave forte e única

#### gsKeyDelivery (KeyDelivery)
- **Uso:** Chave adicional para módulo de delivery
- **Padrão:** `Ped2505IcrM`
- **Produção:** Trocar se necessário

### Boas Práticas

1. **Em Desenvolvimento:**
   ```ini
   secret_key = CloudI0812IcrMmDB
   ```

2. **Em Produção:**
   ```ini
   # Gerar chave forte:
   # openssl rand -base64 32
   secret_key = sua_chave_gerada_aleatoriamente_muito_segura_aqui
   ```

3. **Rotação de Chaves:**
   - Manter chave antiga temporariamente
   - Gerar novos tokens com chave nova
   - Validar tokens com ambas as chaves durante transição
   - Remover chave antiga após expiração de todos os tokens

---

## 🌍 Fuso Horário (gnFusoHorario)

### Valores Comuns

| Timezone | Valor | Região |
|----------|-------|--------|
| UTC | `0` | Universal |
| Brasília | `-3` | Brasil |
| New York (EST) | `-5` | EUA (Leste) |
| Los Angeles (PST) | `-8` | EUA (Oeste) |
| London (GMT) | `0` | Reino Unido |
| Paris (CET) | `+1` | Europa Central |
| Tokyo (JST) | `+9` | Japão |

### Como Funciona

O fuso horário afeta:
1. **Cálculo de Timestamps Unix:** Ajusta para UTC antes de converter
2. **Tokens JWT:** Campos `nbf` (not before) e `exp` (expiration)
3. **Logs:** Horários registrados no banco de dados

**Exemplo:**
```go
// Se timezone = -3 (Brasília)
// Data local: 2025-01-27 15:00:00 (Brasília)
// Será convertido para: 2025-01-27 18:00:00 (UTC)
timestamp := utils.CalcTimeStampUnix(dateTime, -3)
```

---

## 📊 Versão do Sistema

### Formato da Versão

**Pattern:** `Ver MAJOR.MINOR.PATCH.BUILD`

**Exemplo:** `Ver 1.26.4.27`
- **MAJOR:** 1 - Versão principal
- **MINOR:** 26 - Funcionalidades adicionadas
- **PATCH:** 4 - Correções de bugs
- **BUILD:** 27 - Número do build

### Data da Versão

**Pattern:** `YYYY-MM-DDTHH:MM:SS` (ISO 8601)

**Exemplo:** `2025-10-16T11:55:00`
- Data: 16 de Outubro de 2025
- Hora: 11:55:00

### Onde é Usado

1. **API /connect/v1/wsteste:**
   ```json
   {
     "versao": "Ver 1.26.4.27",
     "versaoData": "2025-10-16T11:55:00"
   }
   ```

2. **Logs no Banco de Dados:**
   - Coluna `WSVERSAO` na tabela `WSREQUISICOES`

3. **Logs de Arquivo:**
   - Registrado ao iniciar o servidor

---

## 🔄 Migração de Código WinDev

### Antes (WinDev)

```windev
// Variáveis globais
gsKey is string = "CloudI0812IcrMmDB"
gsIss is string = "WSCloudICrmIntellsys"
gnFusoHorario is int = 0
gsVersao is string = "Ver 1.26.4.27"

// Uso
sToken = HashString(HA_HMAC_SHA_256, sData, gsKey)
nTimestamp = pgCalcTimeStampUnix(dData, tHora, gnFusoHorario)
```

### Depois (Go)

```go
// Configuração carregada de dbinit.ini
cfg, _ := config.LoadConfig("dbinit.ini")

// Uso
h := hmac.New(sha256.New, []byte(cfg.JWT.SecretKey))
h.Write([]byte(data))
token := h.Sum(nil)

timestamp := utils.CalcTimeStampUnix(dateTime, cfg.JWT.Timezone)
```

---

## ✅ Validação

### Verificar Configurações Carregadas

```bash
# Visualizar configurações ao iniciar
./wsicrmrest

# Saída esperada nos logs:
# INFO Iniciando WSICRMREST version=Ver 1.26.4.27 version_date=2025-10-16T11:55:00
```

### Teste de Configuração

Criar arquivo `test_config.go`:

```go
package main

import (
    "fmt"
    "wsicrmrest/internal/config"
)

func main() {
    cfg, err := config.LoadConfig("dbinit.ini")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Secret Key: %s\n", cfg.JWT.SecretKey)
    fmt.Printf("Issuer: %s\n", cfg.JWT.Issuer)
    fmt.Printf("Timezone: %d\n", cfg.JWT.Timezone)
    fmt.Printf("Version: %s\n", cfg.Application.Version)
    fmt.Printf("Version Date: %s\n", cfg.Application.VersionDate)
}
```

Executar:
```bash
go run test_config.go
```

---

## 📚 Referências

- **Código WinDev Original:** Variáveis globais declaradas no início do projeto
- **Implementação Go:** `internal/config/config.go`
- **Configuração:** `dbinit.ini.example`
- **Uso:** Handlers em `internal/handlers/`

---

**Última atualização:** 2025-01-27
