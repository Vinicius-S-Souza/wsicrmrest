# Documentação Swagger - WSICRMREST API

**Data de atualização:** 2025-11-17
**Versão da API:** 1.26.4.27

## 📚 Visão Geral

Esta documentação Swagger descreve todos os endpoints da API WSICRMREST, incluindo:

- Autenticação JWT
- Webhooks da Zenvia (Email e SMS)
- Endpoints de teste

## 🔗 Acessar Documentação

### Localmente

Após iniciar o servidor, acesse:

```
http://localhost:8080/swagger/index.html
```

Ou com HTTPS habilitado:

```
https://localhost:8443/swagger/index.html
```

## 📖 Arquivos Gerados

- **`swagger.yaml`** - Especificação OpenAPI em formato YAML
- **`swagger.json`** - Especificação OpenAPI em formato JSON
- **`docs.go`** - Código Go para servir a documentação

## 🔄 Regenerar Documentação

Sempre que modificar comentários Swagger nos handlers, execute:

```bash
~/go/bin/swag init -g cmd/server/main.go -o docs/swagger
```

Ou use o Makefile (se disponível):

```bash
make swagger
```

## 📋 Endpoints Documentados

### 🔐 Authentication

#### `GET /connect/v1/token`
- **Descrição:** Gera token JWT com Basic Auth
- **Autenticação:** Basic (client_id:client_secret)
- **Rate Limiting:** 60 req/min, 1000 req/hour
- **Security:** Comparação constant-time de secrets

**Novidades na documentação:**
- ✅ Documentação de headers de rate limiting
- ✅ Exemplo de uso com curl
- ✅ Nota sobre proteção contra timing attacks
- ✅ Status code 429 documentado

---

### 🧪 Testing

#### `GET /connect/v1/wsteste`
- **Descrição:** Testa conexão com banco e retorna dados do organizador
- **Autenticação:** Nenhuma
- **Rate Limiting:** Aplicado

---

### 📨 Webhooks

#### `POST /webhook/zenvia/email`
- **Descrição:** Recebe eventos de status de email da Zenvia
- **Autenticação:** Nenhuma (endpoint público)
- **Rate Limiting:** Aplicado

**Novidades na documentação:**
- ✅ **Comportamento para IDs não identificados** documentado:
  - Retorna HTTP 200 (evita retry da Zenvia)
  - Loga warning com messageId
  - Armazena requisição para auditoria
  - Mensagem: "Mensagem não encontrada no banco de dados"
- ✅ Eventos suportados documentados:
  - `sent` → 121 - Agendado
  - `delivered` → 122 - Entregue
  - `read/clicked` → 123 - Aberto
  - `rejected/not_delivered` → 124 - Não Entregue
- ✅ Nota sobre segurança (endpoint público)

#### `POST /webhook/zenvia/sms`
- **Descrição:** Recebe eventos de status de SMS da Zenvia
- **Autenticação:** Nenhuma (endpoint público)
- **Rate Limiting:** Aplicado

**Mesmas melhorias do webhook de email**

---

## 🔒 Segurança Documentada

A documentação Swagger agora inclui:

### Na Descrição Geral da API:

```yaml
description: |
  REST API service for CRM integration, converted from WinDev.

  **Security Features:**
  - SQL Injection Protection (bind variables)
  - HTTPS/TLS Support (configurable)
  - Rate Limiting (60 req/min, 1000 req/hour by default)
  - CORS with origin validation
  - Request size limits (1MB max body)
  - Request timeout (30s default)
  - Security headers (HSTS, X-Frame-Options, etc.)

  **Important Notes:**
  - All API endpoints require JWT Bearer token (except /token and /wsteste)
  - Webhooks from Zenvia are public (no authentication required)
  - Rate limit headers are included in all responses
  - HTTPS is strongly recommended for production
```

### Nos Endpoints:

- **Rate Limiting Headers** documentados
- **Timing Attack Protection** mencionado
- **Comportamento para erros** detalhado
- **Exemplos de uso** incluídos

### Schemes Suportados:

```yaml
schemes:
  - http
  - https
```

---

## 📊 Modelos de Dados

### TokenResponse
```json
{
  "code": "000",
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "clientes lojas ofertas",
  "datetime": 1700000000,
  "modulos": 1,
  "message": "Token gerado com sucesso"
}
```

### ZenviaWebhookRequest
```json
{
  "type": "message_status",
  "message": {
    "messageId": "MSG-123",
    "to": "cliente@example.com",
    "externalId": "12345"
  },
  "messageStatus": {
    "code": "delivered",
    "description": "Message delivered",
    "causes": []
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

**Caso ID não encontrado:**
```json
{
  "success": true,
  "message": "Mensagem não encontrada no banco de dados"
}
```

---

## 🔍 Como Usar

### 1. Testar no Swagger UI

1. Acesse `http://localhost:8080/swagger/index.html`
2. Clique em "Authorize"
3. Para Basic Auth:
   - Username: `seu_client_id`
   - Password: `seu_client_secret`
4. Teste os endpoints

### 2. Importar no Postman

```bash
# Importar arquivo
File > Import > docs/swagger/swagger.json
```

### 3. Usar com Ferramentas CLI

```bash
# curl
curl -X GET "http://localhost:8080/connect/v1/token" \
  -H "Authorization: Basic $(echo -n 'client_id:client_secret' | base64)" \
  -H "Grant_type: client_credentials"

# httpie
http GET http://localhost:8080/connect/v1/token \
  Authorization:"Basic $(echo -n 'client_id:client_secret' | base64)" \
  Grant_type:client_credentials
```

---

## 🎨 Personalização

### Tags

As tags organizam os endpoints:

- **Authentication** - Geração de tokens
- **Testing** - Endpoints de teste
- **Webhooks** - Webhooks da Zenvia

### Security Definitions

```yaml
securityDefinitions:
  BasicAuth:
    type: basic
    description: Basic authentication with client_id as username and client_secret as password

  BearerAuth:
    type: apiKey
    name: Authorization
    in: header
    description: Type "Bearer" followed by a space and JWT token
```

---

## 📝 Changelog Swagger

### 2025-11-17 - v1.26.4.27

**Adicionado:**
- ✅ Documentação completa de segurança na descrição da API
- ✅ Headers de rate limiting documentados
- ✅ Comportamento para IDs não identificados nos webhooks
- ✅ Eventos suportados pelos webhooks
- ✅ Exemplos de uso com curl
- ✅ Status code 429 (Rate Limit Exceeded)
- ✅ Schemes HTTP e HTTPS
- ✅ Tags descritivas

**Atualizado:**
- ✅ Versão da API para 1.26.4.27
- ✅ Descrições de todos os endpoints
- ✅ Mensagens de erro mais detalhadas
- ✅ Security definitions

---

## 🚀 Melhores Práticas

### Para Desenvolvedores

1. **Sempre atualize os comentários Swagger** quando modificar handlers
2. **Regenere a documentação** após mudanças: `swag init`
3. **Teste os endpoints** usando Swagger UI
4. **Valide o JSON/YAML** antes de commitar

### Para Integradores

1. **Use HTTPS em produção** (não HTTP)
2. **Respeite os rate limits** (veja headers nas respostas)
3. **Implemente retry com backoff** para erros 429
4. **Valide os modelos** antes de enviar requests

---

## 📚 Recursos

- **OpenAPI Specification:** https://swagger.io/specification/
- **Swaggo:** https://github.com/swaggo/swag
- **Swagger UI:** https://swagger.io/tools/swagger-ui/

---

## ❓ FAQ

### Como adicionar um novo endpoint?

1. Adicione comentários Swagger no handler:
```go
// MyHandler godoc
// @Summary Breve descrição
// @Description Descrição detalhada
// @Tags NomeDaTag
// @Accept json
// @Produce json
// @Param id path int true "ID do recurso"
// @Success 200 {object} models.Response
// @Router /my-endpoint/{id} [get]
func MyHandler(c *gin.Context) {
    // ...
}
```

2. Regenere o Swagger:
```bash
swag init -g cmd/server/main.go -o docs/swagger
```

### Como documentar headers customizados?

```go
// @Param X-Custom-Header header string false "Descrição do header"
```

### Como documentar rate limiting?

Já está documentado automaticamente! Os headers são:
- `X-RateLimit-Limit-Minute`
- `X-RateLimit-Limit-Hour`
- `X-RateLimit-Remaining-Minute`
- `X-RateLimit-Remaining-Hour`

---

**Documentação mantida por:** Claude Code
**Última atualização:** 2025-11-17
