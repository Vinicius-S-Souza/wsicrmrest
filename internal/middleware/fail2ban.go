// Data de criação: 26/11/2025 10:55
// Última atualização: 26/11/2025 10:55
// Versão: 3.0.0.6
// Middleware de Fail2Ban para proteção contra ataques de força bruta
// Implementação compatível com WSICRMMDB
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Fail2BanConfig define configurações do Fail2Ban
type Fail2BanConfig struct {
	MaxAttempts     int           // Tentativas máximas antes de bloquear
	BanDuration     time.Duration // Duração do banimento
	WindowDuration  time.Duration // Janela de tempo para contar tentativas
	CleanupInterval time.Duration // Intervalo de limpeza de registros antigos
	WhitelistIPs    []string      // IPs que nunca serão bloqueados
}

// IPAttempt rastreia tentativas de um IP
type IPAttempt struct {
	attempts      []time.Time
	banned        bool
	banExpiry     time.Time
	totalAttempts int // Total de tentativas desde o início
	firstAttempt  time.Time
	lastAttempt   time.Time
}

// Fail2Ban gerencia bloqueios por IP
type Fail2Ban struct {
	config    Fail2BanConfig
	ips       map[string]*IPAttempt
	whitelist map[string]bool
	mu        sync.RWMutex
	logger    *zap.SugaredLogger
}

var (
	fail2ban     *Fail2Ban
	fail2banOnce sync.Once
)

// NewFail2Ban cria uma nova instância de Fail2Ban
func NewFail2Ban(config Fail2BanConfig, logger *zap.SugaredLogger) *Fail2Ban {
	fail2banOnce.Do(func() {
		fail2ban = &Fail2Ban{
			config:    config,
			ips:       make(map[string]*IPAttempt),
			whitelist: make(map[string]bool),
			logger:    logger,
		}

		// Adiciona IPs da whitelist
		for _, ip := range config.WhitelistIPs {
			fail2ban.whitelist[ip] = true
		}

		// Inicia limpeza periódica
		go fail2ban.cleanupRoutine()

		logger.Infow("Fail2Ban inicializado",
			"max_attempts", config.MaxAttempts,
			"ban_duration", config.BanDuration,
			"window_duration", config.WindowDuration,
			"whitelist_count", len(config.WhitelistIPs),
		)
	})

	return fail2ban
}

// GetFail2Ban retorna a instância singleton
func GetFail2Ban() *Fail2Ban {
	return fail2ban
}

// cleanupRoutine remove registros antigos periodicamente
func (f *Fail2Ban) cleanupRoutine() {
	ticker := time.NewTicker(f.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		f.cleanup()
	}
}

// cleanup remove IPs que não estão mais banidos e não têm atividade recente
func (f *Fail2Ban) cleanup() {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	cleanedCount := 0

	for ip, attempt := range f.ips {
		// Remove se não está banido e última tentativa foi há mais de 1 hora
		if !attempt.banned && now.Sub(attempt.lastAttempt) > time.Hour {
			delete(f.ips, ip)
			cleanedCount++
		}
	}

	if cleanedCount > 0 {
		f.logger.Debugw("Limpeza Fail2Ban concluída",
			"cleaned_ips", cleanedCount,
			"tracked_ips", len(f.ips),
		)
	}
}

// isWhitelisted verifica se um IP está na whitelist
func (f *Fail2Ban) isWhitelisted(ip string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.whitelist[ip]
}

// IsBanned verifica se um IP está banido
func (f *Fail2Ban) IsBanned(ip string) bool {
	if f.isWhitelisted(ip) {
		return false
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	attempt, exists := f.ips[ip]
	if !exists {
		return false
	}

	// Verifica se ainda está banido
	if attempt.banned {
		if time.Now().Before(attempt.banExpiry) {
			return true
		}
		// Ban expirou, desbloqueia
		attempt.banned = false
		attempt.attempts = []time.Time{}
		f.logger.Infow("IP desbanido automaticamente",
			"ip", ip,
			"total_attempts", attempt.totalAttempts,
		)
	}

	return false
}

// RecordAttempt registra uma tentativa de acesso (falha de autenticação, etc)
func (f *Fail2Ban) RecordAttempt(ip string, path string, statusCode int) {
	if f.isWhitelisted(ip) {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()

	// Obtém ou cria registro do IP
	attempt, exists := f.ips[ip]
	if !exists {
		attempt = &IPAttempt{
			attempts:     []time.Time{},
			firstAttempt: now,
		}
		f.ips[ip] = attempt
	}

	// Atualiza contadores
	attempt.totalAttempts++
	attempt.lastAttempt = now

	// Remove tentativas fora da janela de tempo
	cutoff := now.Add(-f.config.WindowDuration)
	validAttempts := []time.Time{}
	for _, t := range attempt.attempts {
		if t.After(cutoff) {
			validAttempts = append(validAttempts, t)
		}
	}
	attempt.attempts = validAttempts

	// Adiciona tentativa atual
	attempt.attempts = append(attempt.attempts, now)

	// Verifica se deve banir
	if len(attempt.attempts) >= f.config.MaxAttempts {
		attempt.banned = true
		attempt.banExpiry = now.Add(f.config.BanDuration)

		f.logger.Warnw("IP banido por múltiplas tentativas suspeitas",
			"ip", ip,
			"path", path,
			"attempts_in_window", len(attempt.attempts),
			"total_attempts", attempt.totalAttempts,
			"ban_duration", f.config.BanDuration,
			"ban_expiry", attempt.banExpiry,
			"status_code", statusCode,
		)
	}
}

// GetIPStats retorna estatísticas de um IP
func (f *Fail2Ban) GetIPStats(ip string) map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()

	attempt, exists := f.ips[ip]
	if !exists {
		return map[string]interface{}{
			"tracked": false,
		}
	}

	now := time.Now()
	remaining := time.Duration(0)
	if attempt.banned && now.Before(attempt.banExpiry) {
		remaining = attempt.banExpiry.Sub(now)
	}

	return map[string]interface{}{
		"tracked":            true,
		"banned":             attempt.banned,
		"total_attempts":     attempt.totalAttempts,
		"recent_attempts":    len(attempt.attempts),
		"first_attempt":      attempt.firstAttempt,
		"last_attempt":       attempt.lastAttempt,
		"ban_expiry":         attempt.banExpiry,
		"ban_time_remaining": remaining.String(),
	}
}

// Fail2BanMiddleware retorna middleware Gin para Fail2Ban
func Fail2BanMiddleware(config Fail2BanConfig, logger *zap.SugaredLogger) gin.HandlerFunc {
	f2b := NewFail2Ban(config, logger)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Verifica se IP está banido
		if f2b.IsBanned(ip) {
			stats := f2b.GetIPStats(ip)

			f2b.logger.Warnw("Acesso bloqueado - IP banido",
				"ip", ip,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"user_agent", c.Request.UserAgent(),
			)

			c.JSON(http.StatusForbidden, gin.H{
				"code":    "403",
				"message": "Acesso bloqueado devido a múltiplas tentativas suspeitas. Tente novamente mais tarde.",
				"error":   "Forbidden - IP Banned",
				"details": gin.H{
					"ban_time_remaining": stats["ban_time_remaining"],
					"total_attempts":     stats["total_attempts"],
				},
			})
			c.Abort()
			return
		}

		c.Next()

		// Após processar a requisição, verifica se houve falha de autenticação
		statusCode := c.Writer.Status()

		// Registra tentativa suspeita em endpoints sensíveis
		path := c.Request.URL.Path
		isSensitiveEndpoint := path == "/connect/v1/token" ||
			path == "/connect/v1/wsteste" ||
			statusCode == http.StatusUnauthorized ||
			statusCode == http.StatusForbidden

		if isSensitiveEndpoint && (statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden) {
			f2b.RecordAttempt(ip, path, statusCode)
		}
	}
}

// SimpleFail2BanMiddleware retorna middleware com configurações padrão
func SimpleFail2BanMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return Fail2BanMiddleware(Fail2BanConfig{
		MaxAttempts:     5,
		BanDuration:     30 * time.Minute,
		WindowDuration:  10 * time.Minute,
		CleanupInterval: 5 * time.Minute,
		WhitelistIPs:    []string{"127.0.0.1", "::1"},
	}, logger)
}

// StrictFail2BanMiddleware retorna middleware mais restritivo
func StrictFail2BanMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return Fail2BanMiddleware(Fail2BanConfig{
		MaxAttempts:     3,               // Apenas 3 tentativas
		BanDuration:     1 * time.Hour,   // Ban por 1 hora
		WindowDuration:  5 * time.Minute, // Janela de 5 minutos
		CleanupInterval: 5 * time.Minute,
		WhitelistIPs:    []string{"127.0.0.1", "::1"},
	}, logger)
}

// UnbanIP remove manualmente o ban de um IP (para administradores)
func (f *Fail2Ban) UnbanIP(ip string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	attempt, exists := f.ips[ip]
	if !exists || !attempt.banned {
		return false
	}

	attempt.banned = false
	attempt.attempts = []time.Time{}

	f.logger.Infow("IP desbanido manualmente",
		"ip", ip,
	)

	return true
}

// GetBannedIPs retorna lista de IPs banidos
func (f *Fail2Ban) GetBannedIPs() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	now := time.Now()
	banned := []string{}

	for ip, attempt := range f.ips {
		if attempt.banned && now.Before(attempt.banExpiry) {
			banned = append(banned, ip)
		}
	}

	return banned
}
