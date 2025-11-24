package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// IPTracker rastreia tentativas falhas e IPs banidos
type IPTracker struct {
	mu             sync.RWMutex
	failedAttempts map[string][]time.Time // IP -> lista de timestamps de falhas
	bannedIPs      map[string]time.Time   // IP -> timestamp de quando o ban expira
	maxAttempts    int                    // Número máximo de tentativas antes do ban
	banDuration    time.Duration          // Duração do banimento
	timeWindow     time.Duration          // Janela de tempo para contar tentativas
	logger         *zap.SugaredLogger
}

// NewIPTracker cria um novo rastreador de IPs
func NewIPTracker(maxAttempts int, banDuration, timeWindow time.Duration, logger *zap.SugaredLogger) *IPTracker {
	tracker := &IPTracker{
		failedAttempts: make(map[string][]time.Time),
		bannedIPs:      make(map[string]time.Time),
		maxAttempts:    maxAttempts,
		banDuration:    banDuration,
		timeWindow:     timeWindow,
		logger:         logger,
	}

	// Iniciar limpeza periódica de dados antigos
	go tracker.cleanupLoop()

	return tracker
}

// IsBanned verifica se um IP está banido
func (t *IPTracker) IsBanned(ip string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if banTime, exists := t.bannedIPs[ip]; exists {
		if time.Now().Before(banTime) {
			return true
		}
		// Ban expirou, remover
		delete(t.bannedIPs, ip)
	}
	return false
}

// RecordFailure registra uma tentativa falha e retorna true se o IP foi banido
func (t *IPTracker) RecordFailure(ip string, path string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	attempts := t.failedAttempts[ip]

	// Filtrar apenas tentativas recentes (dentro da janela de tempo)
	var recentAttempts []time.Time
	for _, attempt := range attempts {
		if now.Sub(attempt) < t.timeWindow {
			recentAttempts = append(recentAttempts, attempt)
		}
	}

	// Adicionar nova tentativa
	recentAttempts = append(recentAttempts, now)
	t.failedAttempts[ip] = recentAttempts

	// Verificar se excedeu o limite
	if len(recentAttempts) >= t.maxAttempts {
		t.bannedIPs[ip] = now.Add(t.banDuration)
		delete(t.failedAttempts, ip) // Limpar contadores após ban

		t.logger.Warn("IP banido por múltiplas tentativas suspeitas",
			"ip", ip,
			"attempts", len(recentAttempts),
			"path", path,
			"ban_until", t.bannedIPs[ip].Format(time.RFC3339))

		return true // IP foi banido
	}

	return false
}

// GetBannedIPs retorna lista de IPs atualmente banidos
func (t *IPTracker) GetBannedIPs() map[string]time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[string]time.Time)
	now := time.Now()

	for ip, banTime := range t.bannedIPs {
		if now.Before(banTime) {
			result[ip] = banTime
		}
	}

	return result
}

// UnbanIP remove um IP da lista de banidos (uso manual)
func (t *IPTracker) UnbanIP(ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.bannedIPs[ip]; exists {
		delete(t.bannedIPs, ip)
		t.logger.Info("IP desbanido manualmente", "ip", ip)
		return true
	}
	return false
}

// cleanupLoop limpa periodicamente dados expirados
func (t *IPTracker) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		t.cleanup()
	}
}

// cleanup remove dados expirados
func (t *IPTracker) cleanup() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()

	// Remover bans expirados
	for ip, banTime := range t.bannedIPs {
		if now.After(banTime) {
			delete(t.bannedIPs, ip)
			t.logger.Info("Ban expirado removido", "ip", ip)
		}
	}

	// Remover tentativas antigas
	for ip, attempts := range t.failedAttempts {
		var recent []time.Time
		for _, attempt := range attempts {
			if now.Sub(attempt) < t.timeWindow {
				recent = append(recent, attempt)
			}
		}

		if len(recent) == 0 {
			delete(t.failedAttempts, ip)
		} else {
			t.failedAttempts[ip] = recent
		}
	}

	t.logger.Debug("Limpeza de Fail2Ban concluída",
		"banned_ips", len(t.bannedIPs),
		"tracked_ips", len(t.failedAttempts))
}

// Fail2BanMiddleware cria middleware de proteção contra ataques
func Fail2BanMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	// Configuração padrão:
	// - 10 tentativas 404 em 5 minutos = ban de 1 hora
	// - 5 tentativas 401 (auth falha) em 5 minutos = ban de 2 horas
	tracker404 := NewIPTracker(
		10,            // max 10 404s
		1*time.Hour,   // ban por 1 hora
		5*time.Minute, // janela de 5 minutos
		logger,
	)

	tracker401 := NewIPTracker(
		5,              // max 5 falhas de auth
		2*time.Hour,    // ban por 2 horas
		5*time.Minute,  // janela de 5 minutos
		logger,
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path

		// Verificar se IP está banido (por 404s)
		if tracker404.IsBanned(ip) {
			logger.Warn("IP banido (404s) tentou acessar",
				"ip", ip,
				"path", path,
				"user_agent", c.Request.UserAgent())

			c.AbortWithStatusJSON(403, gin.H{
				"error":   "Access forbidden",
				"message": "Too many invalid requests. Please try again later.",
			})
			return
		}

		// Verificar se IP está banido (por auth falhas)
		if tracker401.IsBanned(ip) {
			logger.Warn("IP banido (auth) tentou acessar",
				"ip", ip,
				"path", path,
				"user_agent", c.Request.UserAgent())

			c.AbortWithStatusJSON(403, gin.H{
				"error":   "Access forbidden",
				"message": "Too many authentication failures. Please try again later.",
			})
			return
		}

		c.Next()

		// Após processar requisição, verificar status
		status := c.Writer.Status()

		// Registrar 404s (possível scanning)
		if status == 404 {
			if tracker404.RecordFailure(ip, path) {
				// IP foi banido
				logger.Error("IP BANIDO por múltiplos 404s",
					"ip", ip,
					"path", path,
					"user_agent", c.Request.UserAgent())
			}
		}

		// Registrar falhas de autenticação (401)
		if status == 401 {
			if tracker401.RecordFailure(ip, path) {
				// IP foi banido
				logger.Error("IP BANIDO por múltiplas falhas de autenticação",
					"ip", ip,
					"path", path,
					"user_agent", c.Request.UserAgent())
			}
		}
	}
}
