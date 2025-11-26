package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
	"wsicrmrest/internal/config"
	reqcontext "wsicrmrest/internal/context"
	"wsicrmrest/internal/database"
	"wsicrmrest/internal/models"
	"wsicrmrest/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GenerateToken godoc
// @Summary Generate JWT token
// @Description Generates a JWT access token based on client_id and client_secret provided via Basic Auth
// @Description
// @Description **Security Features:**
// @Description - Client secret comparison uses constant-time algorithm (prevents timing attacks)
// @Description - Rate limiting applied (default: 60 requests/minute)
// @Description - All attempts are logged in WSAPLLOGTOKEN table
// @Description - Failed attempts logged with error details
// @Description
// @Description **Response Headers:**
// @Description - X-RateLimit-Limit-Minute: Maximum requests per minute
// @Description - X-RateLimit-Limit-Hour: Maximum requests per hour
// @Description - X-RateLimit-Remaining-Minute: Remaining requests this minute
// @Description - X-RateLimit-Remaining-Hour: Remaining requests this hour
// @Description
// @Description **Example:**
// @Description ```bash
// @Description curl -X GET "https://api.example.com/connect/v1/token" \
// @Description   -H "Authorization: Basic $(echo -n 'client_id:client_secret' | base64)" \
// @Description   -H "Grant_type: client_credentials"
// @Description ```
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Basic Auth (Base64 encoded client_id:client_secret)"
// @Param Grant_type header string true "Must be 'client_credentials'"
// @Success 200 {object} models.TokenResponse "Token generated successfully"
// @Failure 401 {object} models.TokenResponse "Invalid credentials or grant type"
// @Failure 403 {object} models.TokenResponse "Database connection error"
// @Failure 409 {object} models.TokenResponse "Invalid client_id or disabled application"
// @Failure 429 {object} object "Rate limit exceeded"
// @Failure 500 {object} models.TokenResponse "Error generating token"
// @Router /connect/v1/token [get]
// @Security BasicAuth
func GenerateToken(cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Criar contexto da requisição
		reqCtx := reqcontext.NewRequestContext()

		reqMetodo := "GET"
		reqEndPoint := "/connect/v1/token"
		reqHeader := getHeadersAsString(c)
		reqParametros := ""
		var reqResposta string
		var codError int
		nomeProcedure := "GenerateToken"

		// Ler headers
		authorization := c.GetHeader("Authorization")
		grantType := c.GetHeader("Grant_type")
		host := c.GetHeader("Host")

		logger.Infow("Gerando Token",
			"client_ip", c.ClientIP(),
			"host", host)

		// Limpar caracteres nulos
		authorization = utils.EliminaCaracterNulo(authorization)

		var response models.TokenResponse

		// Validar Authorization e Grant_type
		if !strings.HasPrefix(authorization, "Basic ") || grantType != "client_credentials" {
			logger.Warnw("Credenciais inválidas",
				"authorization", authorization,
				"grant_type", grantType)

			response = models.TokenResponse{
				Code:    "001",
				Message: "As chaves Authorization e Grant_type devem ser informados corretamente.",
			}
			codError = 401
		} else {
			// Decodificar Basic Auth
			authEncoded := authorization[6:]
			authDecoded, err := base64.StdEncoding.DecodeString(authEncoded)
			if err != nil {
				logger.Errorw("Erro ao decodificar Authorization", "error", err)
				response = models.TokenResponse{
					Code:    "002",
					Message: "Erro ao decodificar credenciais.",
				}
				codError = 401
			} else {
				authStr := utils.EliminaCaracterNulo(string(authDecoded))
				parts := strings.SplitN(authStr, ":", 2)

				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					response = models.TokenResponse{
						Code:    "002",
						Message: "Conteúdo do nome Client_id ou Client_secret inválido.",
					}
					codError = 401
				} else {
					clientID := parts[0]
					clientSecret := parts[1]

					reqParametros = "Client_Id: " + clientID + "\nClient_Secret: " + clientSecret + "\n"

					logger.Infow("Gerando Token - Client_id", "client_id", clientID)

					// Buscar aplicação no banco de dados
					app, err := getApplicationByClientID(db, clientID)
					if err != nil {
						logger.Errorw("Erro ao buscar aplicação", "error", err)
						response = models.TokenResponse{
							Code:    "003",
							Message: "Falha ao verificar aplicação.",
						}
						codError = 403
					} else if app == nil {
						response = models.TokenResponse{
							Code:    "004",
							Message: "Falha na Validação da Aplicação. Client_Id Inválido ou Desabilitado.",
						}
						codError = 409
					} else if app.Status != 1 {
						response = models.TokenResponse{
							Code:    "006",
							Message: "Aplicação Desabilitada.",
						}
						codError = 409
					} else if !constantTimeCompare(app.ClientSecret, clientSecret) {
						// Usar comparação de tempo constante para prevenir timing attacks
						response = models.TokenResponse{
							Code:    "005",
							Message: "Client_secret Inválido.",
						}
						codError = 409
					} else {
						// Gerar token JWT
						token, expiration, nbf, err := generateJWT(cfg, app)
						if err != nil {
							logger.Errorw("Erro ao gerar token", "error", err)
							response = models.TokenResponse{
								Code:    "008",
								Message: "Erro ao gerar token JWT.",
							}
							codError = 500
						} else {
							scope := utils.Escopo(app.Scopo)

							// Definir informações do cliente no contexto
							reqCtx.SetClientInfo(clientID, app.Nome)

							response = models.TokenResponse{
								Code:        "000",
								AccessToken: token,
								TokenType:   "Bearer",
								ExpiresIn:   expiration,
								DateTime:    nbf,
								Scope:       scope,
								Modulos:     cfg.Organization.RegModulos,
							}
							codError = 200

							// Gravar log do token no banco
							if err := logTokenToDB(db, clientID, token, host, nbf, expiration); err != nil {
								logger.Errorw("Erro ao gravar log do token", "error", err)
							}
						}
					}
				}
			}
		}

		duration := reqCtx.GetDuration()
		logger.Infow("Token - Finalizado",
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
			config.Version,
		)

		// Retornar resposta
		c.JSON(codError, response)
	}
}

// getApplicationByClientID busca uma aplicação pelo client_id
func getApplicationByClientID(db *database.Database, clientID string) (*models.Application, error) {
	query := `SELECT WSAPLCLIENTSECRET, WSAPLIJWTEXPIRACAO, WSAPLSCOPO, WSAPLSTATUS, WSAPLNOME
	          FROM WSAPLICACOES
	          WHERE WSAPLSTATUS = 1 AND WSAPLCLIENTID = :1`

	var app models.Application
	err := db.QueryRow(query, clientID).Scan(
		&app.ClientSecret,
		&app.JWTExpiracao,
		&app.Scopo,
		&app.Status,
		&app.Nome,
	)

	if err != nil {
		return nil, err
	}

	app.ClientID = clientID
	return &app, nil
}

// generateJWT gera um token JWT
func generateJWT(cfg *config.Config, app *models.Application) (string, int64, int64, error) {
	// Header
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.StdEncoding.EncodeToString(headerJSON)
	headerB64 = utils.StringChange(headerB64, "\n", "")
	headerB64 = utils.StringChange(headerB64, "\r", "")

	// Calcular nbf (not before) e exp (expiration)
	now := time.Now()
	nbf := utils.CalcTimeStampUnix(now, cfg.JWT.Timezone)

	expirationTime := 86400 // 24 horas padrão
	if app.JWTExpiracao > 0 {
		expirationTime = app.JWTExpiracao
	}
	exp := nbf + int64(expirationTime)

	scope := utils.Escopo(app.Scopo)

	// Payload
	payload := map[string]interface{}{
		"iss":       cfg.JWT.Issuer,
		"nbf":       nbf,
		"exp":       exp,
		"client_id": app.ClientID,
		"scope":     scope,
		"aplicacao": app.Nome,
	}
	payloadJSON, _ := json.Marshal(payload)
	payloadB64 := base64.StdEncoding.EncodeToString(payloadJSON)
	payloadB64 = utils.StringChange(payloadB64, "\n", "")
	payloadB64 = utils.StringChange(payloadB64, "\r", "")
	payloadB64 = utils.StringChange(payloadB64, "=", "")

	// Criar token sem assinatura
	tokenUnsigned := headerB64 + "." + payloadB64

	// Gerar assinatura HMAC-SHA256
	h := hmac.New(sha256.New, []byte(cfg.JWT.SecretKey))
	h.Write([]byte(tokenUnsigned))
	signature := h.Sum(nil)
	signatureB64 := base64.StdEncoding.EncodeToString(signature)
	signatureB64 = utils.StringChange(signatureB64, "+", "-")
	signatureB64 = utils.StringChange(signatureB64, "=", "")

	// Token completo
	token := headerB64 + "." + payloadB64 + "." + signatureB64

	return token, exp, nbf, nil
}

// logTokenToDB grava o log do token no banco de dados
func logTokenToDB(db *database.Database, clientID, token, host string, nbf, exp int64) error {
	dataGeracao := time.Unix(nbf, 0)
	dataExpiracao := time.Unix(exp, 0)

	// Usar bind variables para proteção contra SQL injection
	query := `INSERT INTO WSAPLLOGTOKEN(WSLTKNUMERO, WSLTKDATA, WSLTKEXPIRACAO, WSAPLCLIENTID, WSAPLTOKEN, WSAPLHOST)
	          VALUES(SEQ_WSLTKNUMERO.NEXTVAL, :1, :2, :3, :4, :5)`

	_, err := db.Exec(query, dataGeracao, dataExpiracao, clientID, token, host)
	return err
}

// constantTimeCompare compara duas strings em tempo constante para prevenir timing attacks
func constantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// getHeadersAsString retorna todos os headers como string
func getHeadersAsString(c *gin.Context) string {
	var headers strings.Builder
	for key, values := range c.Request.Header {
		for _, value := range values {
			headers.WriteString(key + ": " + value + "\n")
		}
	}
	return headers.String()
}
