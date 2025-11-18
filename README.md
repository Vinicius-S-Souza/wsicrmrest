# WSICRMREST - Web Service REST em Go

Web service REST desenvolvido em Go utilizando Gin e Zap, convertido de procedures WinDev.

## 🚀 Quick Start

**Novo no projeto?** Consulte o [QUICKSTART.md](QUICKSTART.md) para começar em 5 minutos!

## Características

- ✅ Conexão com Oracle via TNSNAMES
- ✅ Autenticação JWT (HS256)
- ✅ Logs rotativos por data
- ✅ Configuração via arquivo INI
- ✅ APIs REST com Gin
- ✅ Log completo de requisições no banco de dados
- ✅ Conversão completa de WinDev para Go
- ✅ Documentação Swagger/OpenAPI interativa

## Pré-requisitos

- Go 1.21 ou superior
- Oracle Client instalado
- Variável de ambiente `LD_LIBRARY_PATH` ou `DYLD_LIBRARY_PATH` configurada para o Oracle Client
- Arquivo `tnsnames.ora` configurado

## Configuração

1. Copie o arquivo de exemplo:
```bash
cp dbinit.ini.example dbinit.ini
```

2. Edite `dbinit.ini` com suas credenciais:
```ini
[database]
driver = 2
tns_name = ORCL
username = seu_usuario
password = sua_senha

[jwt]
secret_key = sua_chave_secreta_minimo_32_caracteres
issuer = WSICRMREST
timezone = -3

[organization]
codigo = 1
nome = Minha Empresa
cnpj = 12345678000199
loja_matriz = 1
cod_isga = 1001

[application]
version = 1.0.0
version_date = 2025-01-27
environment = production
port = 8080
```

## Instalação

```bash
# Baixar dependências
go mod download

# Compilar
go build -o wsicrmrest

# Executar
./wsicrmrest
```

## 📖 Documentação da API

**Documentação Swagger interativa disponível em:** `http://localhost:8080/swagger/index.html`

A documentação Swagger oferece:
- Interface interativa para testar todos os endpoints
- Exemplos de request/response
- Descrição detalhada de cada parâmetro
- Códigos de status HTTP

Consulte [docs/SWAGGER.md](docs/SWAGGER.md) para mais informações sobre como usar a documentação Swagger.

## APIs Disponíveis

### 1. Gerar Token JWT

**Endpoint:** `GET /connect/v1/token`

**Headers:**
- `Authorization`: Basic base64(client_id:client_secret)
- `Grant_type`: client_credentials

**Exemplo:**
```bash
# Encode: "CLIENTE1234567890:a1234567890b"
AUTH=$(echo -n "CLIENTE1234567890:a1234567890b" | base64)

curl -X GET http://localhost:8080/connect/v1/token \
  -H "Authorization: Basic $AUTH" \
  -H "Grant_type: client_credentials"
```

**Resposta de Sucesso (200):**
```json
{
  "code": "000",
  "access_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 1738012345,
  "datetime": 1737926345,
  "scope": "clientes",
  "modulos": 0
}
```

**Respostas de Erro:**
- `401`: Credenciais inválidas
- `403`: Erro no banco de dados
- `409`: Aplicação desabilitada ou Client_secret inválido

### 2. Teste de Conexão

**Endpoint:** `GET /connect/v1/wsteste`

**Exemplo:**
```bash
curl -X GET http://localhost:8080/connect/v1/wsteste
```

**Resposta de Sucesso (200):**
```json
{
  "code": "000",
  "organizadorCodigo": 1,
  "organizadorNome": "Minha Empresa",
  "organizadorCnpj": "12345678000199",
  "organizadorLojaMatriz": 1,
  "organizadorCodIsga": 1001,
  "versao": "1.0.0",
  "versaoData": "2025-01-27"
}
```

## Estrutura do Projeto

```
wsicrmrest/
├── cmd/
│   └── server/
│       └── main.go       # Ponto de entrada
├── dbinit.ini             # Configurações (não versionado)
├── go.mod                 # Dependências Go
├── internal/
│   ├── config/           # Gerenciamento de configurações
│   ├── database/         # Conexão com banco de dados
│   ├── handlers/         # Handlers das APIs
│   ├── logger/           # Sistema de logs
│   ├── middleware/       # Middlewares do Gin
│   ├── models/           # Estruturas de dados
│   ├── routes/           # Definição de rotas
│   └── utils/            # Funções auxiliares
└── log/                  # Diretório de logs (criado automaticamente)
```

## Logs

Os logs são gravados automaticamente na pasta `log/` com o nome:
- `wsicrmrest_YYYY-MM-DD.log`

O arquivo de log é alterado automaticamente a cada dia.

## Tabelas do Banco de Dados

### WSAPLICACOES
Tabela de aplicações registradas:
- `WSAPLCLIENTID` - Client ID
- `WSAPLCLIENTSECRET` - Client Secret
- `WSAPLIJWTEXPIRACAO` - Tempo de expiração do JWT em segundos
- `WSAPLSCOPO` - Código de escopo (bitwise)
- `WSAPLSTATUS` - Status (1=ativo, 0=inativo)
- `WSAPLNOME` - Nome da aplicação

### WSAPLLOGTOKEN
Log de tokens gerados:
- `WSLTKNUMERO` - Número sequencial
- `WSLTKDATA` - Data/hora de geração
- `WSLTKEXPIRACAO` - Data/hora de expiração
- `WSAPLCLIENTID` - Client ID
- `WSAPLTOKEN` - Token gerado
- `WSAPLHOST` - Host da requisição

## Desenvolvimento

### Adicionar novas APIs

1. Crie um novo handler em `internal/handlers/`
2. Adicione a rota em `internal/routes/routes.go`
3. Crie os modelos necessários em `internal/models/`

### Variáveis Globais

O sistema utiliza variáveis configuráveis via `dbinit.ini`. Consulte [docs/GLOBAL_VARIABLES.md](docs/GLOBAL_VARIABLES.md) para documentação completa.

**Principais variáveis:**
- `gsKey` → `JWT.SecretKey` = `CloudI0812IcrMmDB` - Chave HMAC do JWT
- `gsIss` → `JWT.Issuer` = `WSCloudICrmIntellsys` - Issuer do JWT
- `gsKeyDelivery` → `JWT.KeyDelivery` = `Ped2505IcrM` - Chave delivery
- `gnFusoHorario` → `JWT.Timezone` = `0` - Fuso horário (0=UTC, -3=Brasília)
- `gsVersao` → `Application.Version` = `Ver 1.26.4.27` - Versão do sistema
- `gsDataVersao` → `Application.VersionDate` = `2025-10-16T11:55:00` - Data da versão

## Estrutura das Tabelas

Consulte o arquivo [docs/DATABASE_SCHEMA.md](docs/DATABASE_SCHEMA.md) para o script completo de criação das tabelas Oracle necessárias:
- `WSAPLICACOES` - Aplicações registradas
- `WSAPLLOGTOKEN` - Log de tokens gerados
- `WSREQUISICOES` - Log de requisições

## Scripts Úteis

### Testar APIs
```bash
# Editar as variáveis CLIENT_ID e CLIENT_SECRET no script
./scripts/test_apis.sh
```

### Executar o serviço
```bash
./wsicrmrest
```

## Estrutura de Diretórios

```
wsicrmrest/
├── cmd/
│   └── server/
│       └── main.go              # Ponto de entrada
├── go.mod                       # Dependências
├── dbinit.ini                   # Configurações (criar a partir do .example)
├── wsicrmrest                   # Binário compilado
├── internal/
│   ├── config/                  # Gerenciamento de configurações
│   ├── context/                 # Contexto de requisições
│   ├── database/                # Conexão e operações de BD
│   ├── handlers/                # Handlers das APIs
│   ├── logger/                  # Sistema de logs em arquivo
│   ├── middleware/              # Middlewares Gin
│   ├── models/                  # Modelos de dados
│   ├── routes/                  # Definição de rotas
│   └── utils/                   # Funções auxiliares
├── docs/                        # Documentação
│   └── DATABASE_SCHEMA.md       # Estrutura das tabelas
├── scripts/                     # Scripts úteis
│   └── test_apis.sh            # Teste das APIs
└── log/                         # Logs (criado automaticamente)
    └── wsicrmrest_YYYY-MM-DD.log
```

## Configurações Avançadas

### dbinit.ini Completo

```ini
[database]
driver = 2
tns_name = ORCL
username = wsuser
password = wspass

[jwt]
secret_key = sua_chave_secreta_muito_segura_aqui_min_32_chars
issuer = WSICRMREST
timezone = -3

[organization]
codigo = 1
nome = Minha Empresa
cnpj = 12345678000199
loja_matriz = 1
cod_isga = 1001

[application]
version = 1.0.0
version_date = 2025-01-27
environment = production
port = 8080
log_dir = log
ws_grava_log_db = true
ws_detalhe_log_api = false
```

**Configurações de Log:**
- `ws_grava_log_db`: Habilita gravação de logs no banco de dados (padrão: true)
- `ws_detalhe_log_api`: Habilita gravação de detalhes adicionais (padrão: false)

## Códigos de Resposta

### Sucesso
- `200 OK` - Requisição bem-sucedida

### Erros de Autenticação
- `401 Unauthorized` - Credenciais inválidas ou ausentes
- `403 Forbidden` - Erro no acesso ao banco de dados
- `409 Conflict` - Aplicação desabilitada ou credenciais incorretas

### Códigos de Erro Customizados

**Token:**
- `000` - Sucesso
- `001` - Headers Authorization ou Grant_type incorretos
- `002` - Client_id ou Client_secret inválido
- `003` - Falha ao verificar aplicação no banco
- `004` - Client_id inválido ou desabilitado
- `005` - Client_secret inválido
- `006` - Aplicação desabilitada
- `007` - Falha na abertura do banco de dados
- `008` - Erro ao gerar token JWT

**WSTest:**
- `000` - Sucesso
- `005` - Falha na abertura do banco de dados

## Conversão WinDev → Go

Este projeto é uma conversão de procedures WinDev para Go. As principais conversões incluem:

| WinDev | Go |
|--------|-----|
| `pgGerarToken()` | `handlers.GenerateToken()` |
| `pgWSRestTeste()` | `handlers.WSTest()` |
| `pgGravaLogDB()` | `database.GravaLogDB()` |
| `pgScopo()` | `utils.Escopo()` |
| `fgEliminaCaracterNulo()` | `utils.EliminaCaracterNulo()` |
| `pgStringChange()` | `utils.StringChange()` |
| `FcDateTime()` | `utils.FormatDateTimeOracle()` |
| `pgCalcTimeStampUnix()` | `utils.CalcTimeStampUnix()` |

## Troubleshooting

### Erro: "cannot connect to database"
- Verifique se o Oracle Client está instalado
- Confirme se `LD_LIBRARY_PATH` está configurado
- Verifique se o `tnsnames.ora` contém a entrada correta
- Teste a conexão: `sqlplus username/password@tns_name`

### Erro: "table or view does not exist"
- Execute os scripts de criação das tabelas (ver `docs/DATABASE_SCHEMA.md`)
- Verifique as permissões do usuário no banco de dados

### Logs não são gravados
- Verifique permissões do diretório `log/`
- Confirme `ws_grava_log_db = true` no `dbinit.ini`
- Verifique se as tabelas de log existem no banco

## TODO

- [ ] Adicionar validação de JWT nas rotas protegidas
- [ ] Implementar testes unitários
- [ ] Adicionar suporte a SQL Server (driver = 1)
- [ ] Implementar rotação automática de logs
- [ ] Adicionar métricas e health check

## Autor

Convertido de WinDev para Go

## Licença

[Definir licença]
