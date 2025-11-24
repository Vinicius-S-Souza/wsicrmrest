# Análise de Segurança - WSICRMREST

**Data:** 2025-11-24
**Versão:** 3.0.0.2

---

## 🔍 Logs de Ataques Detectados

### Exemplo 1: Scanning de Login
```json
{
  "timestamp": "2025-11-24T14:26:22.617-0300",
  "method": "GET",
  "path": "/login",
  "ip": "45.38.44.221",
  "user_agent": "Go-http-client/1.1",
  "status": 404
}
```

**Análise:**
- **Tipo:** Bot automatizado procurando páginas de login
- **Origem:** 45.38.44.221 (scanner automatizado)
- **User-Agent:** `Go-http-client/1.1` (script em Go)
- **Objetivo:** Encontrar interface de login para tentar credenciais padrão
- **Resultado:** ✅ 404 - Rota não existe na API

### Exemplo 2: Exploit OpenWrt/Router
```json
{
  "timestamp": "2025-11-24T14:46:31.204-0300",
  "method": "GET",
  "path": "/cgi-bin/luci/;stok=/locale",
  "ip": "193.142.147.209",
  "user_agent": "",
  "status": 404
}
```

**Análise:**
- **Tipo:** Tentativa de exploit CVE conhecido em roteadores
- **Origem:** 193.142.147.209 (Europa - scanner distribuído)
- **Path:** `/cgi-bin/luci/;stok=/locale`
  - **LuCI:** Interface web de roteadores OpenWrt
  - **`;stok=`:** Tentativa de injeção de comando
- **Objetivo:** Explorar vulnerabilidade em roteadores OpenWrt
- **Resultado:** ✅ 404 - API não é roteador

---

## 📊 Avaliação de Risco

### Risco Atual: 🟡 MÉDIO

| Categoria | Status | Nível |
|-----------|--------|-------|
| **Exposição Pública** | ⚠️ Ativo | Alto |
| **Scanning Automatizado** | ⚠️ Frequente | Médio |
| **Ataques Direcionados** | ✅ Nenhum detectado | Baixo |
| **Vulnerabilidades Exploradas** | ✅ Nenhuma | Baixo |

### Por que risco médio?

**Pontos Positivos (✅):**
1. Rate limiting implementado (60 req/min, 1000 req/hora)
2. Security headers configurados
3. Request timeout (30s)
4. Body size limit (1MB)
5. Todas tentativas de exploit retornaram 404
6. Autenticação JWT implementada

**Pontos de Atenção (⚠️):**
1. **TLS/HTTPS desabilitado** (tráfego sem criptografia)
2. **CORS permite todas origens** (modo desenvolvimento)
3. Servidor exposto publicamente na internet
4. Scanning constante de bots
5. Sem firewall/whitelist de IPs
6. Sem sistema de detecção de intrusão (IDS)

---

## 🛡️ Recomendações de Segurança

### 🔴 CRÍTICO - Implementar Imediatamente

#### 1. Habilitar HTTPS/TLS

**Problema:** Dados trafegam sem criptografia, incluindo tokens JWT.

**Solução:**
```ini
# dbinit.ini
[tls]
enabled = true
cert_file = certs/server.crt
key_file = certs/server.key
port = 8443
```

**Passos:**
1. Gerar certificado SSL (Let's Encrypt recomendado)
2. Configurar certificado no `dbinit.ini`
3. Redirecionar porta 8080 (HTTP) para 8443 (HTTPS)

**Comandos para gerar certificado auto-assinado (desenvolvimento):**
```bash
mkdir -p certs
openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes -subj "/CN=localhost"
```

**Para produção, use Let's Encrypt:**
```bash
certbot certonly --standalone -d seu-dominio.com
# Certificado: /etc/letsencrypt/live/seu-dominio.com/fullchain.pem
# Chave: /etc/letsencrypt/live/seu-dominio.com/privkey.pem
```

#### 2. Restringir CORS

**Problema:** Qualquer site pode fazer requisições à sua API.

**Solução:**
```ini
# dbinit.ini
[CORS]
AllowedOrigins=https://seu-dominio.com,https://app.seu-dominio.com
AllowCredentials=true
```

**Nunca use `*` em produção com `AllowCredentials=true`!**

---

### 🟡 IMPORTANTE - Implementar em Breve

#### 3. Implementar Fail2Ban ou Similar

**Problema:** IPs maliciosos podem tentar ataques repetidos.

**Solução:** Banir automaticamente IPs com comportamento suspeito.

**Implementação no código Go:**

Criar `internal/middleware/fail2ban.go`:
```go
package middleware

import (
    "sync"
    "time"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type IPTracker struct {
    mu             sync.RWMutex
    failedAttempts map[string][]time.Time
    bannedIPs      map[string]time.Time
    maxAttempts    int
    banDuration    time.Duration
    timeWindow     time.Duration
}

func NewIPTracker(maxAttempts int, banDuration, timeWindow time.Duration) *IPTracker {
    return &IPTracker{
        failedAttempts: make(map[string][]time.Time),
        bannedIPs:      make(map[string]time.Time),
        maxAttempts:    maxAttempts,
        banDuration:    banDuration,
        timeWindow:     timeWindow,
    }
}

func (t *IPTracker) IsBanned(ip string) bool {
    t.mu.RLock()
    defer t.mu.RUnlock()

    if banTime, exists := t.bannedIPs[ip]; exists {
        if time.Now().Before(banTime) {
            return true
        }
        delete(t.bannedIPs, ip)
    }
    return false
}

func (t *IPTracker) RecordFailure(ip string) bool {
    t.mu.Lock()
    defer t.mu.Unlock()

    now := time.Now()
    attempts := t.failedAttempts[ip]

    // Remover tentativas antigas
    var recent []time.Time
    for _, attempt := range attempts {
        if now.Sub(attempt) < t.timeWindow {
            recent = append(recent, attempt)
        }
    }

    recent = append(recent, now)
    t.failedAttempts[ip] = recent

    if len(recent) >= t.maxAttempts {
        t.bannedIPs[ip] = now.Add(t.banDuration)
        delete(t.failedAttempts, ip)
        return true // IP foi banido
    }

    return false
}

func Fail2BanMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
    tracker := NewIPTracker(
        10,              // Máximo de 404s
        1*time.Hour,     // Ban por 1 hora
        5*time.Minute,   // Janela de 5 minutos
    )

    return func(c *gin.Context) {
        ip := c.ClientIP()

        // Verificar se IP está banido
        if tracker.IsBanned(ip) {
            logger.Warn("IP banido tentou acessar",
                "ip", ip,
                "path", c.Request.URL.Path)
            c.AbortWithStatusJSON(403, gin.H{
                "error": "Access forbidden",
            })
            return
        }

        c.Next()

        // Registrar falhas 404
        if c.Writer.Status() == 404 {
            if tracker.RecordFailure(ip) {
                logger.Warn("IP banido por múltiplas tentativas 404",
                    "ip", ip,
                    "path", c.Request.URL.Path)
            }
        }
    }
}
```

**Adicionar no `main.go`:**
```go
router.Use(middleware.Fail2BanMiddleware(log))
```

#### 4. Configurar Firewall

**No servidor (Linux):**

```bash
# Permitir apenas portas necessárias
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 8443/tcp  # HTTPS
sudo ufw allow 22/tcp    # SSH (remova se não usar)
sudo ufw enable

# Se usar CloudFlare, permitir apenas IPs deles
# https://www.cloudflare.com/ips/
```

**No Windows Server:**
```powershell
# Criar regra de firewall
New-NetFirewallRule -DisplayName "WSICRMREST HTTPS" -Direction Inbound -Protocol TCP -LocalPort 8443 -Action Allow

# Bloquear porta HTTP (forçar HTTPS)
New-NetFirewallRule -DisplayName "Block HTTP" -Direction Inbound -Protocol TCP -LocalPort 8080 -Action Block
```

#### 5. Whitelist de IPs (Opcional)

Se sua API é usada apenas por servidores conhecidos:

```ini
# dbinit.ini
[security]
allowed_ips = 192.168.1.0/24,10.0.0.0/8,seu-ip-externo
```

**Implementar middleware:**
```go
func IPWhitelistMiddleware(allowedIPs []string, logger *zap.SugaredLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()

        allowed := false
        for _, allowedIP := range allowedIPs {
            if matchCIDR(ip, allowedIP) {
                allowed = true
                break
            }
        }

        if !allowed {
            logger.Warn("IP não autorizado bloqueado", "ip", ip)
            c.AbortWithStatusJSON(403, gin.H{"error": "Access forbidden"})
            return
        }

        c.Next()
    }
}
```

---

### 🟢 RECOMENDADO - Boas Práticas

#### 6. Monitoramento de Segurança

**Criar alertas para:**
- Múltiplos 404s do mesmo IP (possível scanning)
- Tentativas de SQL injection em parâmetros
- User-Agents suspeitos
- Requests muito grandes (DoS)
- Picos de tráfego anormais

**Script de monitoramento de logs:**

Criar `scripts/monitor_security.sh`:
```bash
#!/bin/bash
# Monitor de segurança - Analisa logs em tempo real

LOG_FILE="log/wsicrmrest_$(date +%Y-%m-%d).log"
ALERT_THRESHOLD=10

echo "Monitorando segurança em $LOG_FILE..."

# IPs com mais de 10 requests 404
echo ""
echo "=== IPs suspeitos (>$ALERT_THRESHOLD 404s) ==="
grep '"status":404' "$LOG_FILE" | \
    grep -oP '"ip":"[^"]*"' | \
    sort | uniq -c | sort -rn | \
    awk -v threshold=$ALERT_THRESHOLD '$1 > threshold {print $0}'

# Paths mais atacados
echo ""
echo "=== Paths mais atacados (404s) ==="
grep '"status":404' "$LOG_FILE" | \
    grep -oP '"path":"[^"]*"' | \
    sort | uniq -c | sort -rn | head -10

# User-Agents suspeitos
echo ""
echo "=== User-Agents suspeitos ==="
grep '"status":404' "$LOG_FILE" | \
    grep -oP '"user_agent":"[^"]*"' | \
    sort | uniq -c | sort -rn | head -10
```

**Executar:**
```bash
chmod +x scripts/monitor_security.sh
./scripts/monitor_security.sh
```

#### 7. Headers de Segurança Adicionais

Verificar se estão implementados em `internal/middleware/security.go`:

```go
// Headers de segurança obrigatórios
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "DENY")
c.Header("X-XSS-Protection", "1; mode=block")
c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
c.Header("Content-Security-Policy", "default-src 'self'")
c.Header("Referrer-Policy", "no-referrer")
c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
```

#### 8. Rotação de Logs

**Problema:** Logs podem crescer indefinidamente.

**Solução:** Configurar logrotate (Linux) ou Task Scheduler (Windows).

**Linux - `/etc/logrotate.d/wsicrmrest`:**
```
/home/vinicius/projetos/wsicrmrest/log/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    missingok
    create 0640 vinicius vinicius
}
```

**Windows - Script PowerShell:**
```powershell
# Criar script em scripts/rotate_logs.ps1
$LogPath = "C:\CRM\WSICRMREST\log"
$DaysToKeep = 30

Get-ChildItem $LogPath -Filter *.log |
    Where-Object { $_.LastWriteTime -lt (Get-Date).AddDays(-$DaysToKeep) } |
    Remove-Item -Force

# Agendar no Task Scheduler para rodar diariamente
```

---

## 📈 Métricas de Segurança Recomendadas

### Monitorar Diariamente:

1. **Taxa de 404s por IP**
   - Normal: < 5 por IP por hora
   - Suspeito: > 10 por IP por hora
   - Crítico: > 50 por IP por hora

2. **Tentativas de autenticação falhas**
   - Normal: < 3 por IP por dia
   - Suspeito: > 10 por IP por dia
   - Crítico: > 50 por IP por dia

3. **Picos de tráfego**
   - Normal: Variação < 200% da média
   - Suspeito: Variação > 300% da média
   - Crítico: Variação > 1000% da média (possível DDoS)

4. **User-Agents vazios ou suspeitos**
   - User-Agent vazio geralmente é scanner
   - User-Agents genéricos (`curl`, `wget`, `python-requests`) sem contexto

---

## ✅ Checklist de Segurança

### Configuração Atual (baseado em dbinit.ini.example):

- [x] Rate limiting habilitado (60/min, 1000/hora)
- [x] Request timeout configurado (30s)
- [x] Body size limit (1MB)
- [x] Security headers
- [x] Autenticação JWT
- [x] Logging de requisições
- [ ] **TLS/HTTPS habilitado** ⚠️ CRÍTICO
- [ ] **CORS restrito** ⚠️ CRÍTICO
- [ ] Firewall configurado
- [ ] IP whitelist (se aplicável)
- [ ] Fail2Ban ou similar
- [ ] Monitoramento de alertas
- [ ] Rotação de logs
- [ ] Backup de dados

---

## 🚨 Ações Imediatas Recomendadas

### Prioridade 1 (Hoje):
1. **Habilitar HTTPS/TLS**
2. **Restringir CORS** para domínios conhecidos

### Prioridade 2 (Esta semana):
3. Configurar firewall
4. Implementar monitoramento de IPs suspeitos
5. Criar script de análise de logs

### Prioridade 3 (Este mês):
6. Implementar Fail2Ban automático
7. Configurar rotação de logs
8. Criar documentação de resposta a incidentes

---

## 📞 Resposta a Incidentes

Se detectar ataque ativo:

1. **Identificar:** Verificar logs e identificar IPs atacantes
2. **Bloquear:** Adicionar IPs ao firewall
   ```bash
   sudo ufw deny from 45.38.44.221
   ```
3. **Documentar:** Registrar tipo de ataque, horário, IPs
4. **Notificar:** Se for ataque sério, considerar notificar hosting provider
5. **Revisar:** Verificar se houve comprometimento de dados

---

## 📚 Referências

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [CWE Top 25](https://cwe.mitre.org/top25/archive/2023/2023_top25_list.html)
- [CVE Database](https://cve.mitre.org/)

---

**Última atualização:** 2025-11-24
**Próxima revisão:** 2025-12-24
