package middleware

import (
	"fmt"
	"net/http"
	"time"
	"wsicrmrest/internal/config"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware cria middleware de segurança com limites de request
func SecurityMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Limitar tamanho do body
		if cfg.Security.MaxBodySize > 0 {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, cfg.Security.MaxBodySize)
		}

		// Timeout de requisição
		if cfg.Application.RequestTimeout > 0 {
			timeout := time.Duration(cfg.Application.RequestTimeout) * time.Second
			c.Request = c.Request.WithContext(c.Request.Context())

			// Criar canal para sinalizar conclusão
			done := make(chan struct{})
			defer close(done)

			// Timer de timeout
			timer := time.AfterFunc(timeout, func() {
				select {
				case <-done:
					return
				default:
					c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
						"code":    "408",
						"message": "Request timeout",
					})
				}
			})
			defer timer.Stop()
		}

		// Headers de segurança
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "no-referrer")

		// HSTS apenas se HTTPS estiver habilitado
		if cfg.TLS.Enabled {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

// MaxBodySizeMiddleware limita o tamanho do body da requisição
func MaxBodySizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()

		// Verificar se o limite foi excedido
		if err := c.Request.Context().Err(); err != nil {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    "413",
				"message": fmt.Sprintf("Request body too large. Max size: %d bytes", maxSize),
			})
			return
		}
	}
}
