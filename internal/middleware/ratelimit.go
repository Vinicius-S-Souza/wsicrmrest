package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"wsicrmrest/internal/config"

	"github.com/gin-gonic/gin"
)

// RateLimiter gerencia rate limiting em memória
type RateLimiter struct {
	requests map[string]*clientLimit
	mu       sync.RWMutex
	perMin   int
	perHour  int
}

type clientLimit struct {
	minuteCount int
	hourCount   int
	minuteReset time.Time
	hourReset   time.Time
}

// NewRateLimiter cria um novo rate limiter
func NewRateLimiter(perMin, perHour int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*clientLimit),
		perMin:   perMin,
		perHour:  perHour,
	}

	// Limpar entradas antigas a cada minuto
	go rl.cleanup()

	return rl
}

// cleanup remove entradas antigas periodicamente
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, limit := range rl.requests {
			if now.After(limit.hourReset) {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow verifica se a requisição deve ser permitida
func (rl *RateLimiter) Allow(ip string) (bool, int, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	limit, exists := rl.requests[ip]

	if !exists {
		rl.requests[ip] = &clientLimit{
			minuteCount: 1,
			hourCount:   1,
			minuteReset: now.Add(1 * time.Minute),
			hourReset:   now.Add(1 * time.Hour),
		}
		return true, rl.perMin - 1, rl.perHour - 1
	}

	// Reset contador por minuto se passou 1 minuto
	if now.After(limit.minuteReset) {
		limit.minuteCount = 0
		limit.minuteReset = now.Add(1 * time.Minute)
	}

	// Reset contador por hora se passou 1 hora
	if now.After(limit.hourReset) {
		limit.hourCount = 0
		limit.hourReset = now.Add(1 * time.Minute)
	}

	// Verificar limites
	if rl.perMin > 0 && limit.minuteCount >= rl.perMin {
		return false, 0, rl.perHour - limit.hourCount
	}

	if rl.perHour > 0 && limit.hourCount >= rl.perHour {
		return false, rl.perMin - limit.minuteCount, 0
	}

	// Incrementar contadores
	limit.minuteCount++
	limit.hourCount++

	return true, rl.perMin - limit.minuteCount, rl.perHour - limit.hourCount
}

// RateLimitMiddleware cria middleware de rate limiting
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	if !cfg.Security.RateLimitEnabled {
		// Rate limiting desabilitado, não aplicar limites
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := NewRateLimiter(cfg.Security.RateLimitPerMin, cfg.Security.RateLimitPerHour)

	return func(c *gin.Context) {
		// Obter IP do cliente
		ip := c.ClientIP()

		// Verificar rate limit
		allowed, remainingMin, remainingHour := limiter.Allow(ip)

		// Adicionar headers de rate limit
		c.Header("X-RateLimit-Limit-Minute", formatInt(cfg.Security.RateLimitPerMin))
		c.Header("X-RateLimit-Limit-Hour", formatInt(cfg.Security.RateLimitPerHour))
		c.Header("X-RateLimit-Remaining-Minute", formatInt(remainingMin))
		c.Header("X-RateLimit-Remaining-Hour", formatInt(remainingHour))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    "429",
				"message": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

// formatInt converte int para string
func formatInt(n int) string {
	if n < 0 {
		return "0"
	}
	return fmt.Sprintf("%d", n)
}
