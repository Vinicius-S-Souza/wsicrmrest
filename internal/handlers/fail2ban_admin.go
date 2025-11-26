// Data de criação: 26/11/2025 11:00
// Última atualização: 26/11/2025 11:00
// Versão: 3.0.0.6
// Handlers administrativos para gerenciar Fail2Ban
// Implementação compatível com WSICRMMDB
package handlers

import (
	"net/http"
	"wsicrmrest/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Fail2BanStatusResponse resposta de status do Fail2Ban
type Fail2BanStatusResponse struct {
	BannedIPs   []string               `json:"banned_ips"`
	BannedCount int                    `json:"banned_count"`
	IPStats     map[string]interface{} `json:"ip_stats,omitempty"`
}

// Fail2BanUnbanRequest requisição para desbanir IP
type Fail2BanUnbanRequest struct {
	IP string `json:"ip" binding:"required"`
}

// Fail2BanGetStatus godoc
// @Summary Lista IPs banidos pelo Fail2Ban
// @Description Retorna lista de IPs atualmente banidos e estatísticas detalhadas
// @Tags Fail2Ban
// @Accept json
// @Produce json
// @Param ip query string false "IP específico para consultar estatísticas"
// @Success 200 {object} Fail2BanStatusResponse
// @Failure 500 {object} map[string]interface{}
// @Router /connect/v1/fail2ban/status [get]
func Fail2BanGetStatus(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		f2b := middleware.GetFail2Ban()

		if f2b == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    "503",
				"message": "Fail2Ban não está inicializado",
			})
			return
		}

		bannedIPs := f2b.GetBannedIPs()

		response := Fail2BanStatusResponse{
			BannedIPs:   bannedIPs,
			BannedCount: len(bannedIPs),
		}

		// Se forneceu IP específico, retorna estatísticas
		queryIP := c.Query("ip")
		if queryIP != "" {
			response.IPStats = f2b.GetIPStats(queryIP)
		}

		logger.Infow("Consulta de status Fail2Ban",
			"admin_ip", c.ClientIP(),
			"banned_count", len(bannedIPs),
			"query_ip", queryIP,
		)

		c.JSON(http.StatusOK, response)
	}
}

// Fail2BanUnbanIP godoc
// @Summary Remove ban de um IP
// @Description Remove manualmente o banimento de um IP específico
// @Tags Fail2Ban
// @Accept json
// @Produce json
// @Param request body Fail2BanUnbanRequest true "IP a ser desbanido"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /connect/v1/fail2ban/unban [post]
func Fail2BanUnbanIP(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Fail2BanUnbanRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "Parâmetros inválidos",
				"error":   err.Error(),
			})
			return
		}

		f2b := middleware.GetFail2Ban()

		if f2b == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    "503",
				"message": "Fail2Ban não está inicializado",
			})
			return
		}

		if f2b.UnbanIP(req.IP) {
			logger.Infow("IP desbanido manualmente",
				"ip", req.IP,
				"admin_ip", c.ClientIP(),
			)

			c.JSON(http.StatusOK, gin.H{
				"code":    "200",
				"message": "IP desbanido com sucesso",
				"ip":      req.IP,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "404",
				"message": "IP não está banido ou não encontrado",
				"ip":      req.IP,
			})
		}
	}
}

// Fail2BanGetIPStats godoc
// @Summary Estatísticas de um IP específico
// @Description Retorna estatísticas detalhadas de tentativas e banimentos de um IP
// @Tags Fail2Ban
// @Accept json
// @Produce json
// @Param ip path string true "Endereço IP"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /connect/v1/fail2ban/ip/{ip} [get]
func Fail2BanGetIPStats(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Param("ip")
		if ip == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "400",
				"message": "IP não fornecido",
			})
			return
		}

		f2b := middleware.GetFail2Ban()

		if f2b == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    "503",
				"message": "Fail2Ban não está inicializado",
			})
			return
		}

		stats := f2b.GetIPStats(ip)

		if tracked, ok := stats["tracked"].(bool); ok && !tracked {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "404",
				"message": "IP não está sendo rastreado",
				"ip":      ip,
			})
			return
		}

		logger.Infow("Consulta de estatísticas de IP",
			"ip", ip,
			"admin_ip", c.ClientIP(),
		)

		c.JSON(http.StatusOK, stats)
	}
}
