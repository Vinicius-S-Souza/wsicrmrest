# Exemplos de Uso - WSICRMREST

Este documento demonstra como usar as principais funcionalidades do sistema.

## 📝 Índice

- [Acessar Variáveis Globais](#acessar-variáveis-globais)
- [Criar um Novo Handler](#criar-um-novo-handler)
- [Gerar Token JWT](#gerar-token-jwt)
- [Validar Token JWT](#validar-token-jwt)
- [Gravar Log no Banco](#gravar-log-no-banco)
- [Usar Sistema de Escopos](#usar-sistema-de-escopos)

---

## Acessar Variáveis Globais

### Exemplo Completo

```go
package handlers

import (
    "wsicrmrest/internal/config"
    "github.com/gin-gonic/gin"
)

func MyHandler(cfg *config.Config, ...) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Variáveis JWT
        secretKey := cfg.JWT.SecretKey      // gsKey = "CloudI0812IcrMmDB"
        issuer := cfg.JWT.Issuer            // gsIss = "WSCloudICrmIntellsys"
        keyDelivery := cfg.JWT.KeyDelivery  // gsKeyDelivery = "Ped2505IcrM"
        timezone := cfg.JWT.Timezone        // gnFusoHorario = 0

        // Variáveis de Versão
        version := cfg.Application.Version  // gsVersao = "Ver 1.26.4.27"
        versionDate := cfg.Application.VersionDate // gsDataVersao

        // Variáveis de Organização
        orgCodigo := cfg.Organization.Codigo     // gnOrgCodigo
        orgNome := cfg.Organization.Nome         // gsOrgNome
        orgCNPJ := cfg.Organization.CNPJ         // gnOrgCnpj
        regModulos := cfg.Organization.RegModulos // gnRegModulos

        // Usar as variáveis...
        logger.Infow("Processando requisição",
            "version", version,
            "organization", orgNome,
            "timezone", timezone)
    }
}
```

---

## Criar um Novo Handler

### Passo 1: Criar o Handler

Arquivo: `internal/handlers/meu_handler.go`

```go
package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "wsicrmrest/internal/config"
    reqcontext "wsicrmrest/internal/context"
    "wsicrmrest/internal/database"
    "wsicrmrest/internal/models"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// MeuHandler handler para GET /api/v1/minha-rota
func MeuHandler(cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Criar contexto da requisição
        reqCtx := reqcontext.NewRequestContext()

        reqMetodo := "GET"
        reqEndPoint := "/api/v1/minha-rota"
        reqHeader := getHeadersAsString(c)
        reqParametros := ""
        var reqResposta string
        var codError int
        nomeProcedure := "MeuHandler"

        logger.Infow("MeuHandler iniciado",
            "client_ip", c.ClientIP(),
            "version", cfg.Application.Version)

        var response models.MeuResponse

        // Sua lógica aqui...
        response = models.MeuResponse{
            Code:    "000",
            Message: "Sucesso",
            Data:    "Seus dados aqui",
        }
        codError = 200

        duration := reqCtx.GetDuration()
        logger.Infow("MeuHandler - Finalizado",
            "duration", duration,
            "cod_error", codError)

        // Converter resposta para JSON
        respJSON, _ := json.Marshal(response)
        reqResposta = string(respJSON)

        // Gravar log no banco de dados
        go db.GravaLogDB(
            reqCtx.UUID,
            reqMetodo,
            reqEndPoint,
            reqHeader,
            reqParametros,
            codError,
            reqResposta,
            nomeProcedure,
            reqCtx.ClientID,
            reqCtx.NomeAplicacao,
            reqCtx.StartTime,
            c.ClientIP(),
            cfg.Application.WSGravaLogDB,
            cfg.Application.WSDetalheLogAPI,
            reqCtx.DetalheLogAPI,
            cfg.Application.Version,
        )

        c.JSON(codError, response)
    }
}
```

### Passo 2: Criar o Model

Arquivo: `internal/models/models.go`

```go
// MeuResponse representa a resposta da minha API
type MeuResponse struct {
    Code    string `json:"code"`
    Message string `json:"message,omitempty"`
    Data    string `json:"data,omitempty"`
}
```

### Passo 3: Adicionar a Rota

Arquivo: `internal/routes/routes.go`

```go
func SetupRoutes(router *gin.Engine, cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) {
    // ... rotas existentes ...

    // Nova rota
    apiGroup := router.Group("/api/v1")
    {
        apiGroup.GET("/minha-rota", handlers.MeuHandler(cfg, db, logger))
    }
}
```

---

## Gerar Token JWT

### Exemplo com Variáveis Globais

```go
package mypackage

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "time"

    "wsicrmrest/internal/config"
    "wsicrmrest/internal/utils"
)

func GerarMeuToken(cfg *config.Config, userID string) (string, error) {
    // Header
    header := map[string]string{
        "typ": "JWT",
        "alg": "HS256",
    }
    headerJSON, _ := json.Marshal(header)
    headerB64 := base64.StdEncoding.EncodeToString(headerJSON)
    headerB64 = utils.StringChange(headerB64, "\n", "")
    headerB64 = utils.StringChange(headerB64, "\r", "")

    // Payload
    now := time.Now()
    nbf := utils.CalcTimeStampUnix(now, cfg.JWT.Timezone)
    exp := nbf + 86400 // 24 horas

    payload := map[string]interface{}{
        "iss":     cfg.JWT.Issuer,      // gsIss
        "nbf":     nbf,
        "exp":     exp,
        "user_id": userID,
    }
    payloadJSON, _ := json.Marshal(payload)
    payloadB64 := base64.StdEncoding.EncodeToString(payloadJSON)
    payloadB64 = utils.StringChange(payloadB64, "\n", "")
    payloadB64 = utils.StringChange(payloadB64, "=", "")

    // Token sem assinatura
    tokenUnsigned := headerB64 + "." + payloadB64

    // Assinatura com gsKey
    h := hmac.New(sha256.New, []byte(cfg.JWT.SecretKey)) // gsKey
    h.Write([]byte(tokenUnsigned))
    signature := h.Sum(nil)
    signatureB64 := base64.StdEncoding.EncodeToString(signature)
    signatureB64 = utils.StringChange(signatureB64, "+", "-")
    signatureB64 = utils.StringChange(signatureB64, "=", "")

    // Token completo
    token := headerB64 + "." + payloadB64 + "." + signatureB64

    return token, nil
}
```

---

## Validar Token JWT

### Middleware de Validação

```go
package middleware

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "wsicrmrest/internal/config"

    "github.com/gin-gonic/gin"
)

// ValidateJWT middleware para validar token JWT
func ValidateJWT(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ler token do header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": "401",
                "message": "Token não fornecido",
            })
            c.Abort()
            return
        }

        token := authHeader[7:] // Remove "Bearer "

        // Separar partes do token
        parts := strings.Split(token, ".")
        if len(parts) != 3 {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": "401",
                "message": "Token inválido",
            })
            c.Abort()
            return
        }

        headerB64 := parts[0]
        payloadB64 := parts[1]
        signatureB64 := parts[2]

        // Verificar assinatura usando gsKey
        tokenUnsigned := headerB64 + "." + payloadB64
        h := hmac.New(sha256.New, []byte(cfg.JWT.SecretKey)) // gsKey
        h.Write([]byte(tokenUnsigned))
        expectedSignature := h.Sum(nil)
        expectedSignatureB64 := base64.StdEncoding.EncodeToString(expectedSignature)
        expectedSignatureB64 = strings.ReplaceAll(expectedSignatureB64, "+", "-")
        expectedSignatureB64 = strings.ReplaceAll(expectedSignatureB64, "=", "")

        if signatureB64 != expectedSignatureB64 {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": "401",
                "message": "Assinatura inválida",
            })
            c.Abort()
            return
        }

        // Decodificar payload
        payloadB64 += strings.Repeat("=", (4-len(payloadB64)%4)%4)
        payloadJSON, err := base64.StdEncoding.DecodeString(payloadB64)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": "401",
                "message": "Erro ao decodificar payload",
            })
            c.Abort()
            return
        }

        var payload map[string]interface{}
        if err := json.Unmarshal(payloadJSON, &payload); err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code": "401",
                "message": "Payload inválido",
            })
            c.Abort()
            return
        }

        // Verificar expiração
        if exp, ok := payload["exp"].(float64); ok {
            if time.Now().Unix() > int64(exp) {
                c.JSON(http.StatusUnauthorized, gin.H{
                    "code": "401",
                    "message": "Token expirado",
                })
                c.Abort()
                return
            }
        }

        // Adicionar payload ao contexto
        c.Set("jwt_payload", payload)
        c.Next()
    }
}
```

### Usar o Middleware

```go
// Em routes.go
protectedGroup := router.Group("/api/v1/protected")
protectedGroup.Use(middleware.ValidateJWT(cfg))
{
    protectedGroup.GET("/dados", handlers.DadosProtegidos(cfg, db, logger))
}
```

---

## Gravar Log no Banco

### Uso Simples

```go
// Gravar log de forma assíncrona (recomendado)
go db.GravaLogDB(
    uuid,           // UUID da requisição
    "GET",          // Método HTTP
    "/api/v1/test", // Endpoint
    headers,        // Headers HTTP
    params,         // Parâmetros
    200,            // Código HTTP
    response,       // Resposta JSON
    "MeuHandler",   // Nome do handler
    clientID,       // Client ID
    appName,        // Nome da aplicação
    startTime,      // Tempo de início
    clientIP,       // IP do cliente
    true,           // Gravar no banco?
    false,          // Gravar detalhes?
    "",             // Detalhes adicionais
    cfg.Application.Version, // Versão
)
```

---

## Usar Sistema de Escopos

### Verificar Escopos

```go
package handlers

import "wsicrmrest/internal/utils"

func VerificarEscopos(scopeCode int64, requiredScopes ...string) bool {
    // Obter escopos da aplicação
    appScopes := utils.Escopo(scopeCode)

    // Verificar se contém todos os escopos necessários
    for _, required := range requiredScopes {
        if !strings.Contains(appScopes, required) {
            return false
        }
    }
    return true
}

// Exemplo de uso
func MeuHandler(cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Buscar aplicação
        app, _ := getApplicationByClientID(db, clientID)

        // Verificar se tem escopo "clientes" e "produtos"
        if !VerificarEscopos(app.Scopo, "clientes", "produtos") {
            c.JSON(403, gin.H{
                "code": "403",
                "message": "Sem permissão para acessar este recurso",
            })
            return
        }

        // Continuar processamento...
    }
}
```

### Códigos de Escopo (Bitwise)

```go
// Exemplos de códigos de escopo
const (
    ScopeClientes   = 1    // Bit 1
    ScopeLojas      = 2    // Bit 2
    ScopeOfertas    = 4    // Bit 3
    ScopeProdutos   = 8    // Bit 4
    ScopePontos     = 16   // Bit 5
    ScopePrivate    = 32   // Bit 6
    ScopeConvenio   = 64   // Bit 7
    ScopeGiftcard   = 128  // Bit 8
    ScopeCobranca   = 256  // Bit 9
    ScopeBasico     = 512  // Bit 10
    ScopeSistema    = 1024 // Bit 11
    ScopeTerceiros  = 2048 // Bit 12
    ScopeTotem      = 4096 // Bit 13
)

// Combinar escopos
scopeClientesProdutos := ScopeClientes | ScopeProdutos // 1 + 8 = 9
utils.Escopo(9) // Retorna: "clientes produtos"

// Todos os escopos
scopeTodos := ScopeClientes | ScopeLojas | ScopeOfertas | ScopeProdutos |
              ScopePontos | ScopePrivate | ScopeConvenio | ScopeGiftcard |
              ScopeCobranca | ScopeBasico | ScopeSistema | ScopeTerceiros |
              ScopeTotem // = 8191
utils.Escopo(8191) // Retorna todos os escopos
```

---

## Trabalhar com Fuso Horário

### Converter para Timestamp Unix

```go
package handlers

import (
    "time"
    "wsicrmrest/internal/utils"
    "wsicrmrest/internal/config"
)

func ConverterDataHora(cfg *config.Config) {
    // Data/hora local
    dataHora := time.Now()

    // Converter para timestamp Unix usando fuso horário configurado
    timestamp := utils.CalcTimeStampUnix(dataHora, cfg.JWT.Timezone)

    // timestamp agora está em UTC
    fmt.Printf("Timestamp Unix: %d\n", timestamp)
}
```

### Formatar para Oracle

```go
package handlers

import (
    "time"
    "wsicrmrest/internal/utils"
)

func FormatarParaOracle() {
    dataHora := time.Now()

    // Com milissegundos
    sqlWithMs := utils.FormatDateTimeOracle(dataHora, true)
    // Retorna: TO_TIMESTAMP('01/27/2025 15:30:45.50', 'MM/DD/YYYY HH24:MI:SS.FF')

    // Sem milissegundos
    sqlNoMs := utils.FormatDateTimeOracle(dataHora, false)
    // Retorna: TO_TIMESTAMP('01/27/2025 15:30:45', 'MM/DD/YYYY HH24:MI:SS')
}
```

---

## 📚 Mais Exemplos

- **Handlers:** Veja `internal/handlers/token.go` e `internal/handlers/wstest.go`
- **Funções Auxiliares:** Veja `internal/utils/helpers.go`
- **Database:** Veja `internal/database/log.go`
- **Configuração:** Veja `internal/config/config.go`

---

**Última atualização:** 2025-01-27
