package middleware

import (
	"wsicrmrest/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CORS configura políticas de Cross-Origin Resource Sharing
func CORS(corsConfig config.CORSConfig, environment string, logger *zap.SugaredLogger) gin.HandlerFunc {
	// Avisar se CORS aberto em produção
	if len(corsConfig.AllowedOrigins) == 0 && environment == "production" {
		logger.Warn("⚠️  CORS configurado para permitir TODAS as origens (*) em PRODUÇÃO! Isso é um risco de segurança.")
		logger.Warn("Configure AllowedOrigins no dbinit.ini para restringir as origens permitidas.")
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Se allowedOrigins está vazio, permite todos (desenvolvimento)
		if len(corsConfig.AllowedOrigins) == 0 {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			// Verifica se origem está na lista permitida
			allowed := false
			for _, allowedOrigin := range corsConfig.AllowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			} else if origin != "" {
				// Origem não permitida - logar e bloquear
				logger.Warnw("CORS: Origem não permitida bloqueada",
					"origin", origin,
					"ip", c.ClientIP())
				// Não definir header CORS - navegador bloqueará
				return
			}
		}

		// Define métodos permitidos
		c.Header("Access-Control-Allow-Methods", corsConfig.AllowedMethods)

		// Define headers permitidos
		c.Header("Access-Control-Allow-Headers", corsConfig.AllowedHeaders)

		// Headers expostos
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		// Permite credenciais
		if corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Define tempo de cache do preflight
		c.Header("Access-Control-Max-Age", corsConfig.MaxAge)

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
