package routes

import (
	"wsicrmrest/internal/config"
	"wsicrmrest/internal/database"
	"wsicrmrest/internal/handlers"
	"wsicrmrest/internal/middleware"

	_ "wsicrmrest/docs/swagger" // Swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(router *gin.Engine, cfg *config.Config, db *database.Database, logger *zap.SugaredLogger) {
	// Middleware de logging
	router.Use(middleware.Logger(logger))

	// Grupo de rotas /connect
	connectGroup := router.Group("/connect/v1")
	{
		// GET /connect/v1/token - Gerar token JWT
		connectGroup.GET("/token", handlers.GenerateToken(cfg, db, logger))

		// GET /connect/v1/wsteste - Teste de conexão
		connectGroup.GET("/wsteste", handlers.WSTest(cfg, db, logger))
	}

	// Grupo de rotas /webhook
	webhookGroup := router.Group("/webhook")
	{
		// POST /webhook/zenvia/email - Webhook Zenvia para eventos de email
		webhookGroup.POST("/zenvia/email", handlers.ZenviaEmailWebhook(cfg, db, logger))

		// POST /webhook/zenvia/sms - Webhook Zenvia para eventos de SMS
		webhookGroup.POST("/zenvia/sms", handlers.ZenviaSMSWebhook(cfg, db, logger))
	}

	// Grupo de rotas /connect/v1/fail2ban - Administração do Fail2Ban
	fail2banGroup := router.Group("/connect/v1/fail2ban")
	{
		// GET /connect/v1/fail2ban/status - Lista IPs banidos
		fail2banGroup.GET("/status", handlers.Fail2BanGetStatus(logger))

		// POST /connect/v1/fail2ban/unban - Desbanir IP manualmente
		fail2banGroup.POST("/unban", handlers.Fail2BanUnbanIP(logger))

		// GET /connect/v1/fail2ban/ip/:ip - Estatísticas de IP específico
		fail2banGroup.GET("/ip/:ip", handlers.Fail2BanGetIPStats(logger))
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
