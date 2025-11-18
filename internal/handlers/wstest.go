package handlers

import (
	"encoding/json"
	"net/http"
	"wsicrmrest/internal/config"
	reqcontext "wsicrmrest/internal/context"
	"wsicrmrest/internal/database"
	"wsicrmrest/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WSTest godoc
// @Summary Test database connection
// @Description Tests the database connection and returns organization information
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} models.WSTestResponse "Connection successful"
// @Failure 403 {object} models.WSTestResponse "Database connection failed"
// @Router /connect/v1/wsteste [get]
func WSTest(cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Criar contexto da requisição
		reqCtx := reqcontext.NewRequestContext()

		reqMetodo := "GET"
		reqEndPoint := "/connect/v1/wsteste"
		reqHeader := getHeadersAsString(c)
		reqParametros := ""
		var reqResposta string
		var codError int
		nomeProcedure := "WSTest"

		clientIP := c.ClientIP()
		host := c.GetHeader("Host")

		logger.Infow("WSTESTE iniciado",
			"client_ip", clientIP,
			"host", host)

		var response models.WSTestResponse

		// Testar conexão com banco de dados
		if err := db.DB.Ping(); err != nil {
			logger.Errorw("Falha no acesso ao Banco de Dados", "error", err)

			response = models.WSTestResponse{
				Code:    "005",
				Message: "Falha na Abertura do Banco de Dados.",
				Erro:    err.Error(),
			}
			codError = 403
		} else {
			logger.Info("===> WSTESTE - Sucesso !!!")

			response = models.WSTestResponse{
				Code:                  "000",
				OrganizadorCodigo:     cfg.Organization.Codigo,
				OrganizadorNome:       cfg.Organization.Nome,
				OrganizadorCNPJ:       cfg.Organization.CNPJ,
				OrganizadorLojaMatriz: cfg.Organization.LojaMatriz,
				OrganizadorCodISGA:    cfg.Organization.CodISGA,
				Versao:                config.Version,
				VersaoData:            config.VersionDate,
			}
			codError = 200
		}

		duration := reqCtx.GetDuration()
		logger.Infow("WSTest - Finalizado",
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

		// Retornar resposta baseada no código de erro
		if codError != http.StatusOK {
			c.JSON(codError, response)
		} else {
			c.JSON(http.StatusOK, response)
		}
	}
}
