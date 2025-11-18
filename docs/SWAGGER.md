# Documentação Swagger da API WSICRMREST

## Visão Geral

A API WSICRMREST possui documentação interativa Swagger/OpenAPI que permite explorar e testar todos os endpoints disponíveis.

## Acessando a Documentação

Com o servidor em execução, acesse:

```
http://localhost:8080/swagger/index.html
```

Ou substitua `localhost:8080` pelo endereço e porta onde o servidor está rodando.

## Recursos Disponíveis

A documentação Swagger oferece:

- **Listagem de todos os endpoints**: Organizados por tags (Authentication, System, Webhooks)
- **Descrição detalhada**: Cada endpoint possui descrição completa da funcionalidade
- **Parâmetros**: Documentação de todos os parâmetros necessários
- **Exemplos de Request/Response**: Modelos de dados para requisições e respostas
- **Teste interativo**: Execute requisições diretamente pela interface
- **Códigos de resposta**: Todos os códigos HTTP possíveis documentados

## Endpoints Documentados

### Authentication
- **GET /connect/v1/token**: Geração de token JWT
  - Requer autenticação Basic Auth (client_id:client_secret)
  - Retorna token de acesso, tipo, expiração e escopos

### System
- **GET /connect/v1/wsteste**: Teste de conexão com banco de dados
  - Retorna informações da organização
  - Útil para verificar saúde do sistema

### Webhooks
- **POST /webhook/zenvia/email**: Webhook para eventos de email Zenvia
  - Processa status: SENT, DELIVERED, READ, CLICKED, REJECTED, NOT_DELIVERED
  - Atualiza status de emails no banco
  - Cria ocorrências para emails rejeitados

- **POST /webhook/zenvia/sms**: Webhook para eventos de SMS Zenvia
  - Processa status: SENT, DELIVERED, REJECTED, NOT_DELIVERED
  - Atualiza status de SMS no banco
  - Cria ocorrências para números inválidos

## Autenticação

A API utiliza dois tipos de autenticação:

### 1. Basic Auth
Usado apenas para geração de token:
- Username: `client_id`
- Password: `client_secret`
- Header: `Authorization: Basic base64(client_id:client_secret)`

### 2. Bearer Token (JWT)
Usado para endpoints protegidos (futuros):
- Header: `Authorization: Bearer {token}`

## Modelos de Dados

### TokenResponse
```json
{
  "code": "000",
  "message": "Token gerado com sucesso",
  "access_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "datetime": 1635724800,
  "scope": "clientes lojas ofertas",
  "modulos": 15
}
```

### WSTestResponse
```json
{
  "code": "000",
  "organizadorCodigo": 1,
  "organizadorNome": "Empresa Exemplo",
  "organizadorCnpj": "12345678901234",
  "organizadorLojaMatriz": 1,
  "organizadorCodIsga": 100,
  "versao": "1.0.0",
  "versaoData": "2025-01-15"
}
```

### ZenviaWebhookRequest
```json
{
  "type": "MESSAGE_STATUS",
  "message": {
    "to": "destinatario@example.com",
    "externalId": "12345"
  },
  "messageStatus": {
    "code": "DELIVERED",
    "description": "Mensagem entregue com sucesso",
    "causes": [
      {
        "reason": "INVALID_EMAIL",
        "details": "Endereço inválido"
      }
    ]
  }
}
```

### ZenviaWebhookResponse
```json
{
  "success": true,
  "message": "Webhook processado com sucesso"
}
```

## Testando pela Interface

### Exemplo: Gerar Token

1. Acesse `/swagger/index.html`
2. Localize o endpoint `GET /connect/v1/token` na seção **Authentication**
3. Clique em "Try it out"
4. Preencha os headers:
   - `Authorization`: `Basic {base64(client_id:client_secret)}`
   - `Grant_type`: `client_credentials`
5. Clique em "Execute"
6. Veja a resposta com o token gerado

### Exemplo: Testar Conexão

1. Localize o endpoint `GET /connect/v1/wsteste` na seção **System**
2. Clique em "Try it out"
3. Clique em "Execute"
4. Veja informações da organização na resposta

### Exemplo: Simular Webhook

1. Localize o endpoint `POST /webhook/zenvia/email` na seção **Webhooks**
2. Clique em "Try it out"
3. Edite o JSON de exemplo no campo "body"
4. Clique em "Execute"
5. Veja o resultado do processamento

## Regenerando a Documentação

Quando adicionar ou modificar endpoints:

```bash
# Instalar swag CLI (uma vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Regenerar documentação
~/go/bin/swag init -g cmd/server/main.go --output docs/swagger
```

Ou usando o caminho completo:
```bash
$HOME/go/bin/swag init -g cmd/server/main.go --output docs/swagger
```

## Anotações Swagger

As anotações seguem o padrão `swaggo`:

```go
// NomeDoHandler godoc
// @Summary Resumo curto
// @Description Descrição detalhada do endpoint
// @Tags Nome da Tag
// @Accept json
// @Produce json
// @Param nome_param tipo tipo obrigatorio "descrição"
// @Success 200 {object} models.ResponseType "descrição"
// @Failure 400 {object} models.ErrorType "descrição"
// @Router /caminho/endpoint [metodo]
// @Security TipoDeAutenticacao
func NomeDoHandler() {}
```

## Estrutura de Arquivos

```
docs/swagger/
├── docs.go          # Documentação Go gerada
├── swagger.json     # Especificação OpenAPI 3.0 em JSON
└── swagger.yaml     # Especificação OpenAPI 3.0 em YAML
```

## Integração com Ferramentas

### Postman
1. Acesse `http://localhost:8080/swagger/doc.json`
2. Importe o JSON no Postman

### Insomnia
1. Acesse `http://localhost:8080/swagger/doc.json`
2. Importe o JSON no Insomnia

### OpenAPI Generator
```bash
# Gerar client em várias linguagens
npx @openapitools/openapi-generator-cli generate \
  -i http://localhost:8080/swagger/doc.json \
  -g python \
  -o ./client
```

## Referências

- Swagger/OpenAPI: https://swagger.io/specification/
- swaggo/swag: https://github.com/swaggo/swag
- gin-swagger: https://github.com/swaggo/gin-swagger

## Notas

- A documentação Swagger está disponível apenas quando o servidor está em execução
- Headers de autenticação sensíveis não são logados na documentação
- Todos os exemplos de request/response são baseados nos modelos reais da API
- A interface Swagger UI é interativa e pode ser usada para testes em ambiente de desenvolvimento
