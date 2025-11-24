# Monitoramento de Segurança - WSICRMREST

**Data de criação:** 2025-11-24
**Última atualização:** 2025-11-24

Este documento descreve as ferramentas de monitoramento de segurança implementadas no WSICRMREST.

---

## 📋 Índice

- [Fail2Ban Middleware](#fail2ban-middleware)
- [Scripts de Monitoramento](#scripts-de-monitoramento)
- [Uso dos Scripts](#uso-dos-scripts)
- [Interpretação de Resultados](#interpretação-de-resultados)
- [Ações Recomendadas](#ações-recomendadas)

---

## 🛡️ Fail2Ban Middleware

### O que é?

Middleware implementado em Go que detecta e bloqueia automaticamente IPs com comportamento suspeito.

### Como Funciona?

**Duas categorias de proteção:**

1. **Proteção contra Scanning (404s)**
   - **Limite:** 10 requisições 404 em 5 minutos
   - **Ação:** Ban de 1 hora
   - **Objetivo:** Bloquear bots que fazem scanning de vulnerabilidades

2. **Proteção contra Brute Force (401s)**
   - **Limite:** 5 falhas de autenticação em 5 minutos
   - **Ação:** Ban de 2 horas
   - **Objetivo:** Bloquear tentativas de força bruta em autenticação

### Implementação

**Arquivo:** `internal/middleware/fail2ban.go`

**Recursos:**
- ✅ Rastreamento em memória (sem banco de dados)
- ✅ Limpeza automática de dados antigos a cada 5 minutos
- ✅ Thread-safe (usa sync.RWMutex)
- ✅ Logs detalhados de IPs banidos
- ✅ Resposta 403 com mensagem clara ao usuário

**Integração:**
```go
// cmd/server/main.go e internal/service/windows_service.go
router.Use(middleware.Fail2BanMiddleware(log))
```

### Exemplo de Log

**IP sendo banido:**
```json
{
  "level": "WARN",
  "timestamp": "2025-11-24T15:30:22-0300",
  "message": "IP banido por múltiplas tentativas suspeitas",
  "ip": "45.38.44.221",
  "attempts": 10,
  "path": "/login",
  "ban_until": "2025-11-24T16:30:22-0300"
}
```

**IP banido tentando acessar:**
```json
{
  "level": "WARN",
  "timestamp": "2025-11-24T15:31:00-0300",
  "message": "IP banido (404s) tentou acessar",
  "ip": "45.38.44.221",
  "path": "/admin",
  "user_agent": "Go-http-client/1.1"
}
```

---

## 📊 Scripts de Monitoramento

### 1. Monitor de Segurança (Análise Completa)

**Linux/WSL:** `scripts/monitor_security.sh`
**Windows:** `scripts/monitor_security.ps1`

**O que faz:**
- Analisa logs do dia atual
- Detecta IPs suspeitos com múltiplos 404s
- Lista IPs banidos pelo Fail2Ban
- Identifica falhas de autenticação
- Mostra paths mais atacados
- Identifica User-Agents suspeitos
- Calcula estatísticas gerais
- Fornece recomendações de ação

**Quando usar:**
- Diariamente para revisar segurança
- Após detectar atividade suspeita
- Antes de tomar decisões sobre firewall

### 2. Monitor em Tempo Real

**Linux/WSL:** `scripts/watch_security.sh`

**O que faz:**
- Monitora logs em tempo real usando `tail -f`
- Alerta imediatamente sobre:
  - 🚨 IPs banidos
  - 🔒 Tentativas bloqueadas (403)
  - ⚠️  Requisições 404 (scanning)
  - 🔑 Falhas de autenticação (401)
  - ❌ Erros do servidor (5xx)
  - ✓ Requisições normais (sample de 10%)

**Quando usar:**
- Durante ataques ativos
- Para monitoramento ao vivo
- Debugging de problemas de segurança

---

## 🚀 Uso dos Scripts

### Linux/WSL

#### Análise Completa
```bash
cd /home/vinicius/projetos/wsicrmrest
./scripts/monitor_security.sh
```

**Saída exemplo:**
```
=========================================
  WSICRMREST - Monitor de Segurança
=========================================

Analisando: log/wsicrmrest_2025-11-24.log
Data: 2025-11-24 15:45:30

=== IPs SUSPEITOS (Múltiplos 404s) ===

🚨 ALERTA: 45.38.44.221 - 15 tentativas 404
⚠️  ATENÇÃO: 193.142.147.209 - 8 tentativas 404
✓ Normal: 192.168.1.100 - 2 tentativas 404

=== IPs BANIDOS (Fail2Ban) ===

Total de IPs banidos hoje: 2

  🔒 45.38.44.221 - banido 1 vez(es)
  🔒 193.142.147.209 - banido 1 vez(es)

...
```

#### Monitoramento em Tempo Real
```bash
./scripts/watch_security.sh
```

**Saída exemplo:**
```
=========================================
  WSICRMREST - Monitor em Tempo Real
=========================================

Monitorando: log/wsicrmrest_2025-11-24.log
Pressione Ctrl+C para sair

⚠️  [15:46:12] 404: 45.38.44.221 -> /admin
⚠️  [15:46:15] 404: 45.38.44.221 -> /login
⚠️  [15:46:18] 404: 45.38.44.221 -> /wp-admin
🚨 [15:46:21] IP BANIDO: 45.38.44.221 (por múltiplos 404s)
🔒 [15:46:25] BLOQUEADO: 45.38.44.221 tentou acessar /api
✓ [15:46:30] OK: /wsteste
```

### Windows

#### Análise Completa (PowerShell)
```powershell
cd C:\CRM\WSICRMREST
.\scripts\monitor_security.ps1
```

#### Análise Agendada (Task Scheduler)

Criar tarefa que executa diariamente:

```powershell
# Criar script de tarefa
$action = New-ScheduledTaskAction -Execute "PowerShell.exe" `
    -Argument "-ExecutionPolicy Bypass -File C:\CRM\WSICRMREST\scripts\monitor_security.ps1" `
    -WorkingDirectory "C:\CRM\WSICRMREST"

$trigger = New-ScheduledTaskTrigger -Daily -At "23:00"

$principal = New-ScheduledTaskPrincipal -UserId "SYSTEM" -LogonType ServiceAccount -RunLevel Highest

Register-ScheduledTask -TaskName "WSICRMREST Security Monitor" `
    -Action $action -Trigger $trigger -Principal $principal `
    -Description "Análise diária de segurança do WSICRMREST"
```

---

## 📖 Interpretação de Resultados

### IPs Suspeitos

| Tentativas 404 | Status | Ação |
|----------------|--------|------|
| 1-5 | ✓ Normal | Nenhuma |
| 6-10 | ⚠️ Atenção | Monitorar |
| >10 | 🚨 Alerta | Verificar/Banir |

### Taxa de 404

| Taxa | Status | Ação |
|------|--------|------|
| <10% | ✓ Normal | Nenhuma |
| 10-20% | ⚠️ Elevada | Investigar |
| >20% | 🚨 Muito Alta | Ação imediata |

### Falhas de Autenticação

| Tentativas 401 | Status | Ação |
|----------------|--------|------|
| 1-2 | ✓ Normal | Nenhuma (erro de usuário) |
| 3-5 | ⚠️ Atenção | Verificar se é legítimo |
| >5 | 🚨 Alerta | Possível ataque |

### User-Agents Suspeitos

**Indicadores de bot/scanner:**
- User-Agent vazio: `""`
- Genéricos: `curl`, `wget`, `python-requests`, `Go-http-client`
- Desatualizados: versões antigas de navegadores
- Malformados: strings estranhas ou muito longas

---

## 🔧 Ações Recomendadas

### Quando detectar IP suspeito

#### 1. Verificar se é Fail2Ban está funcionando
```bash
# Linux
grep "IP BANIDO" log/wsicrmrest_$(date +%Y-%m-%d).log

# Windows
Select-String -Path "log\wsicrmrest_$(Get-Date -Format 'yyyy-MM-dd').log" -Pattern "IP BANIDO"
```

Se Fail2Ban já baniu, **não precisa ação manual**.

#### 2. Banir IP manualmente no firewall (se necessário)

**Linux (UFW):**
```bash
sudo ufw deny from 45.38.44.221
sudo ufw status numbered
```

**Linux (iptables):**
```bash
sudo iptables -A INPUT -s 45.38.44.221 -j DROP
sudo iptables-save > /etc/iptables/rules.v4
```

**Windows (Firewall):**
```powershell
New-NetFirewallRule -DisplayName "Block 45.38.44.221" `
    -Direction Inbound -RemoteAddress 45.38.44.221 -Action Block
```

#### 3. Verificar padrão de ataque

```bash
# Ver todos os paths que o IP tentou acessar
grep "45.38.44.221" log/wsicrmrest_$(date +%Y-%m-%d).log | \
    grep -oP '"path":"[^"]*"' | sort | uniq -c
```

**Padrões comuns:**
- `/admin`, `/wp-admin`, `/phpmyadmin` → Scanning de CMS
- `/cgi-bin/*`, `/api/*` → Procurando vulnerabilidades
- Paths aleatórios longos → SQL injection tentativas

#### 4. Relatar abuso (opcional)

Se ataque for persistente, considere reportar ao provedor do IP:

```bash
# Descobrir provedor
whois 45.38.44.221 | grep -i abuse
```

Enviar email ao abuse contact com:
- IP atacante
- Timestamps do ataque
- Logs relevantes
- Tipo de ataque detectado

### Quando taxa de 404 está muito alta

1. **Verificar se é ataque distribuído:**
   ```bash
   ./scripts/monitor_security.sh | grep "404"
   ```

2. **Habilitar HTTPS/TLS** (se ainda não estiver)
   - Reduz visibilidade em scanners automáticos

3. **Considerar CloudFlare** ou similar
   - Proteção DDoS
   - Bot protection
   - Rate limiting adicional

4. **Revisar endpoints públicos**
   - Remover endpoints desnecessários
   - Adicionar autenticação em endpoints sensíveis

### Quando detectar falhas de auth constantes

1. **Verificar se são usuários legítimos:**
   - Perguntar time se estão com problemas
   - Verificar se credenciais expiraram

2. **Se for ataque:**
   - Fail2Ban já deve ter banido
   - Considerar aumentar penalidade (ban mais longo)
   - Notificar time de segurança

3. **Medidas preventivas:**
   - Implementar MFA (Multi-Factor Authentication)
   - Política de senha forte
   - Alertar usuários sobre senhas comprometidas

---

## 📈 Métricas para Dashboard

Se integrar com sistema de monitoramento (Grafana, Prometheus, etc.):

**Métricas importantes:**
- Total de requisições por minuto
- Taxa de 404 em tempo real
- IPs únicos por hora
- IPs banidos acumulados
- Falhas de autenticação por hora
- Latência média de requests
- Erros 5xx (problemas do servidor)

**Exemplo de query:**
```bash
# Requisições por minuto (últimos 10 minutos)
grep '"message":"Request"' log/wsicrmrest_$(date +%Y-%m-%d).log | \
    tail -1000 | \
    grep -oP '"timestamp":"[^"]*"' | \
    cut -c17-21 | \
    uniq -c
```

---

## 🔔 Alertas Automáticos

### Via Email (Linux com sendmail)

Criar `scripts/alert_security.sh`:
```bash
#!/bin/bash
LOG_FILE="log/wsicrmrest_$(date +%Y-%m-%d).log"
BANNED_COUNT=$(grep -c "IP BANIDO" "$LOG_FILE")

if [ "$BANNED_COUNT" -gt 5 ]; then
    echo "ALERTA: $BANNED_COUNT IPs foram banidos hoje!" | \
        mail -s "WSICRMREST Security Alert" admin@example.com
fi
```

**Agendar no crontab:**
```cron
0 * * * * /path/to/scripts/alert_security.sh
```

### Via Webhook (Slack, Discord, etc.)

```bash
#!/bin/bash
WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
BANNED_IPS=$(grep "IP BANIDO" log/wsicrmrest_$(date +%Y-%m-%d).log | \
    grep -oP '"ip":"[^"]*"' | sed 's/"ip":"//;s/"//' | sort -u)

if [ -n "$BANNED_IPS" ]; then
    MESSAGE="IPs banidos hoje:\n$BANNED_IPS"
    curl -X POST -H 'Content-type: application/json' \
        --data "{\"text\":\"$MESSAGE\"}" \
        "$WEBHOOK_URL"
fi
```

---

## 📚 Referências

- [OWASP Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)
- [Fail2Ban Official Documentation](https://www.fail2ban.org/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)

---

**Última atualização:** 2025-11-24
**Próxima revisão:** 2025-12-24
