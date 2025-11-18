package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger middleware para logging de requisições
func Logger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Processar requisição
		c.Next()

		// Calcular duração
		duration := time.Since(start)

		// Log da requisição
		logger.Infow("Request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", c.Writer.Status(),
			"ip", c.ClientIP(),
			"duration", duration,
			"user_agent", c.Request.UserAgent(),
		)
	}
}
