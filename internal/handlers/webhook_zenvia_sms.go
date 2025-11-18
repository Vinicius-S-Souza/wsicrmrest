package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
	"wsicrmrest/internal/config"
	reqcontext "wsicrmrest/internal/context"
	"wsicrmrest/internal/database"
	"wsicrmrest/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZenviaSMSWebhook godoc
// @Summary Zenvia SMS Webhook
// @Description Processes SMS status events received from Zenvia webhook. Updates SMS status, logs API history, and creates occurrences for bounced messages.
// @Description
// @Description **Behavior for Unknown Message IDs:**
// @Description - Returns HTTP 200 with success=true (prevents Zenvia retries)
// @Description - Logs warning with messageId and error details
// @Description - Stores webhook request in WSREQUISICOES table for audit
// @Description - Message: "Mensagem não encontrada no banco de dados"
// @Description
// @Description **Supported Events:**
// @Description - sent (121 - Agendado)
// @Description - delivered (122 - Entregue)
// @Description - read (123 - Aberto)
// @Description - rejected/not_delivered (124 - Não Entregue)
// @Description
// @Description **Security:** No authentication required (public webhook endpoint)
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param body body models.ZenviaWebhookRequest true "Zenvia webhook payload"
// @Success 200 {object} models.ZenviaWebhookResponse "Webhook processed successfully (includes cases where message ID is not found)"
// @Failure 400 {object} models.ZenviaWebhookResponse "Invalid request body or JSON"
// @Router /webhook/zenvia/sms [post]
func ZenviaSMSWebhook(cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Criar contexto da requisição
		reqCtx := reqcontext.NewRequestContext()

		reqMetodo := "POST"
		reqEndPoint := "/webhook/zenvia/sms"
		reqHeader := getHeadersAsString(c)
		var reqParametros string
		var reqResposta string
		var codError int
		nomeProcedure := "ZenviaSMSWebhook"

		// Log inicial
		logger.Infow("====================================================================================================")
		logger.Infow("Iniciando tratamento de mensagem recebida pelo WebHookZenvia")
		logger.Infow("====================================================================================================")

		// Ler o body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Errorw("Erro ao ler body da requisição", "error", err)
			response := models.ZenviaWebhookResponse{
				Success: false,
				Message: "Erro ao ler body da requisição",
			}
			codError = 400

			respJSON, _ := json.Marshal(response)
			reqResposta = string(respJSON)

			// Log no banco
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

			c.JSON(codError, response)
			return
		}

		reqParametros = string(bodyBytes)

		// Log do payload recebido
		logger.Infow("Payload recebido do webhook Zenvia SMS", "payload", string(bodyBytes))

		// Parse JSON
		var webhookRequest models.ZenviaWebhookRequest
		if err := json.Unmarshal(bodyBytes, &webhookRequest); err != nil {
			logger.Errorw("Erro ao fazer parse do JSON", "error", err, "body", string(bodyBytes))
			response := models.ZenviaWebhookResponse{
				Success: false,
				Message: "JSON inválido",
			}
			codError = 400

			respJSON, _ := json.Marshal(response)
			reqResposta = string(respJSON)

			// Log no banco
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

			c.JSON(codError, response)
			return
		}

		// Validar tipo de callback
		tipoCallback := strings.ToLower(webhookRequest.Type)
		if tipoCallback != "message_status" {
			logger.Warnw("Tipo de mensagem inválida",
				"tipo", tipoCallback,
				"esperado", "message_status")

			logger.Infow("Encerrando tratamento de mensagem recebida pelo WebHookZenvia: tipo de mensagem inválida",
				"tipo", tipoCallback)
			logger.Infow("====================================================================================================")

			response := models.ZenviaWebhookResponse{
				Success: true,
				Message: fmt.Sprintf("Tipo de mensagem não processado: %s", tipoCallback),
			}
			codError = 200

			respJSON, _ := json.Marshal(response)
			reqResposta = string(respJSON)

			// Log no banco
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

			c.JSON(codError, response)
			return
		}

		// Extrair dados da mensagem
		numeroCelular := webhookRequest.Message.To
		sEvento := strings.ToLower(webhookRequest.MessageStatus.Code)
		sDescricaoEvento := webhookRequest.MessageStatus.Description
		if len(sDescricaoEvento) > 1000 {
			sDescricaoEvento = sDescricaoEvento[:1000]
		}
		sGMessageID := webhookRequest.Message.MessageId
		sExternalID := webhookRequest.Message.ExternalID

		logger.Infow("Processando ocorrência",
			"from", numeroCelular,
			"evento", sEvento,
			"messageId", sGMessageID,
			"externalId", sExternalID,
			"description", sDescricaoEvento)

		/*
			// Validar se o ID é numérico (smscodigo)
			matched, _ := regexp.MatchString(`^\d+$`, sGMessageID)
			if !matched {
				logger.Warnw("O ID recebido do Zenvia não corresponde ao tipo de ID esperado (smscodigo)",
					"id_recebido", sGMessageID)

				response := models.ZenviaWebhookResponse{
					Success: true,
					Message: "ID não corresponde ao formato esperado",
				}
				codError = 200

				respJSON, _ := json.Marshal(response)
				reqResposta = string(respJSON)

				// Log no banco
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

				c.JSON(codError, response)
				return
			}
		*/

		// Determinar status e tag baseado no evento
		var nLogsApiStatus int
		var sLogsApiHisDescricao string
		var sLogsApiTag string

		switch sEvento {
		case "sent":
			nLogsApiStatus = 121
			sLogsApiHisDescricao = fmt.Sprintf("121 - Agendado: [%s]", sDescricaoEvento)
			sLogsApiTag = "AgendadoProvedor"
		case "delivered":
			nLogsApiStatus = 122
			sLogsApiHisDescricao = fmt.Sprintf("122 - Entregue: [%s]", sDescricaoEvento)
			sLogsApiTag = "Entregue"
		case "rejected", "not_delivered":
			sDetalhes := ""
			if len(webhookRequest.MessageStatus.Causes) > 0 {
				cause := webhookRequest.MessageStatus.Causes[0]
				sDetalhes = fmt.Sprintf(" %s (%s)", cause.Reason, cause.Details)
			}
			nLogsApiStatus = 124
			sLogsApiHisDescricao = fmt.Sprintf("124 - Não Entregue: [%s][%s]", sDescricaoEvento, sDetalhes)
			sLogsApiTag = "NãoEntregue"
		default:
			logger.Warnw("Status fora do escopo de tratamento",
				"evento", sEvento,
				"descricao", sDescricaoEvento)

			logger.Infow("Encerrando tratamento de mensagem recebida pelo WebHookZenvia: Status fora do escopo")
			logger.Infow("====================================================================================================")

			response := models.ZenviaWebhookResponse{
				Success: true,
				Message: fmt.Sprintf("Status não processado: %s", sEvento),
			}
			codError = 200

			respJSON, _ := json.Marshal(response)
			reqResposta = string(respJSON)

			// Log no banco
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

			c.JSON(codError, response)
			return
		}

		// Buscar dados do SMS no banco
		logger.Infow("Buscando mensagem no banco de dados",
			"messageId", sGMessageID,
			"externalId", sExternalID,
			"celular", numeroCelular)

		smsData, err := db.GetSMSByMessageID(sGMessageID)
		if err != nil {
			logger.Warnw("O ID da mensagem NÃO foi encontrado na tabela smsmensagem",
				"messageId", sGMessageID,
				"externalId", sExternalID,
				"celular", numeroCelular,
				"evento", sEvento,
				"error", err)

			response := models.ZenviaWebhookResponse{
				Success: true,
				Message: "Mensagem não encontrada no banco de dados",
			}
			codError = 200

			respJSON, _ := json.Marshal(response)
			reqResposta = string(respJSON)

			// Log no banco
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

			c.JSON(codError, response)
			return
		}

		logger.Infow("ID da mensagem encontrado na tabela smsmensagem",
			"messageId", sGMessageID,
			"smscodigo", smsData.SMSNumero,
			"clicodigo", smsData.CliCodigo,
			"logsapiid", smsData.LogsApiId)

		// Preparar dados para inserção no log
		sLogsApiEnvio := "{}"
		sLogsApiRetorno := string(bodyBytes)
		sLogsApiRetorno = strings.ReplaceAll(sLogsApiRetorno, "'", "`")
		sLogsApiDtaCadastro := time.Now().Format("20060102150405")
		sLogsApiHisDescricao = strings.ReplaceAll(sLogsApiHisDescricao, "'", "")

		// Inserir nas tabelas logsApi e logsApiHistorico
		if err := db.InsereLogsAPISMS(
			smsData.SMSNumero,
			0, // logsApiTipId
			nLogsApiStatus,
			sLogsApiEnvio,
			sLogsApiRetorno,
			sLogsApiDtaCadastro,
			sLogsApiHisDescricao,
			1, // ocorrenciaNumero
			numeroCelular,
			smsData.LogsApiId,
			sLogsApiTag,
		); err != nil {
			logger.Errorw("Erro ao inserir logs API SMS", "error", err)
		}

		// Inserir ocorrência se o SMS for inconsistente - STATUS 125 (bounce)
		// Nota: O código WinDev usa 125, mas o switch case só gera 121, 122, 124
		// Mantendo a lógica original para 125
		if nLogsApiStatus == 125 {
			if err := db.InsereOcorrenciaSmsInconsistente(
				numeroCelular,
				smsData.CliCodigo,
				721, // Tipo de ocorrência para SMS inconsistente
				"Zenvia",
			); err != nil {
				logger.Errorw("Erro ao inserir ocorrência de SMS inconsistente", "error", err)
			}
		}

		logger.Infow("Processamento concluído com sucesso",
			"from", numeroCelular,
			"status", nLogsApiStatus,
			"tag", sLogsApiTag)

		logger.Infow("Finalizando tratamento de mensagem recebida pelo WebHookZenvia")
		logger.Infow("====================================================================================================")

		response := models.ZenviaWebhookResponse{
			Success: true,
			Message: "Webhook processado com sucesso",
		}
		codError = 200

		respJSON, _ := json.Marshal(response)
		reqResposta = string(respJSON)

		// Log no banco
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

		c.JSON(codError, response)
	}
}
