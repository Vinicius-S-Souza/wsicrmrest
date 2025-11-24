# Monitor de Segurança - WSICRMREST
# Analisa logs em tempo real para detectar ataques
# Data: 2025-11-24
# Uso: .\scripts\monitor_security.ps1

# Configurações
$LogDir = "log"
$AlertThreshold404 = 10
$AlertThresholdAuth = 5
$Today = Get-Date -Format "yyyy-MM-dd"
$LogFile = "$LogDir\wsicrmrest_$Today.log"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "  WSICRMREST - Monitor de Segurança" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Analisando: $LogFile"
Write-Host "Data: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
Write-Host ""

# Verificar se arquivo de log existe
if (-not (Test-Path $LogFile)) {
    Write-Host "ERRO: Arquivo de log não encontrado: $LogFile" -ForegroundColor Red
    Write-Host "Verifique se o servidor está rodando e gerando logs."
    exit 1
}

# Ler conteúdo do log
$logContent = Get-Content $LogFile

# ============================================
# 1. IPs com múltiplos 404s (Scanning)
# ============================================
Write-Host "=== IPs SUSPEITOS (Múltiplos 404s) ===" -ForegroundColor Yellow
Write-Host ""

$ips404 = $logContent | Where-Object { $_ -match '"status":404' } |
    ForEach-Object {
        if ($_ -match '"ip":"([^"]+)"') {
            $matches[1]
        }
    } | Group-Object | Sort-Object Count -Descending

foreach ($ipGroup in $ips404) {
    $ip = $ipGroup.Name
    $count = $ipGroup.Count

    if ($count -gt $AlertThreshold404) {
        Write-Host "🚨 ALERTA: $ip - $count tentativas 404" -ForegroundColor Red
    } elseif ($count -gt 5) {
        Write-Host "⚠️  ATENÇÃO: $ip - $count tentativas 404" -ForegroundColor Yellow
    } else {
        Write-Host "✓ Normal: $ip - $count tentativas 404" -ForegroundColor Green
    }
}

Write-Host ""

# ============================================
# 2. IPs Banidos pelo Fail2Ban
# ============================================
Write-Host "=== IPs BANIDOS (Fail2Ban) ===" -ForegroundColor Yellow
Write-Host ""

$bannedLines = $logContent | Where-Object { $_ -match "IP BANIDO" }
$bannedCount = $bannedLines.Count

if ($bannedCount -gt 0) {
    Write-Host "Total de IPs banidos hoje: $bannedCount" -ForegroundColor Red
    Write-Host ""

    $bannedIPs = $bannedLines | ForEach-Object {
        if ($_ -match '"ip":"([^"]+)"') {
            $matches[1]
        }
    } | Group-Object | Sort-Object Count -Descending

    foreach ($ipGroup in $bannedIPs) {
        $ip = $ipGroup.Name
        $count = $ipGroup.Count
        Write-Host "  🔒 $ip - banido $count vez(es)" -ForegroundColor Red
    }
} else {
    Write-Host "✓ Nenhum IP foi banido hoje" -ForegroundColor Green
}

Write-Host ""

# ============================================
# 3. Falhas de Autenticação (401)
# ============================================
Write-Host "=== FALHAS DE AUTENTICAÇÃO (401) ===" -ForegroundColor Yellow
Write-Host ""

$authFailures = $logContent | Where-Object { $_ -match '"status":401' }
$authFailureCount = $authFailures.Count

if ($authFailureCount -gt 0) {
    Write-Host "Total de falhas de autenticação: $authFailureCount"
    Write-Host ""

    $authIPs = $authFailures | ForEach-Object {
        if ($_ -match '"ip":"([^"]+)"') {
            $matches[1]
        }
    } | Group-Object | Sort-Object Count -Descending

    foreach ($ipGroup in $authIPs) {
        $ip = $ipGroup.Name
        $count = $ipGroup.Count

        if ($count -gt $AlertThresholdAuth) {
            Write-Host "🚨 ALERTA: $ip - $count falhas de auth" -ForegroundColor Red
        } elseif ($count -gt 2) {
            Write-Host "⚠️  ATENÇÃO: $ip - $count falhas de auth" -ForegroundColor Yellow
        } else {
            Write-Host "✓ Normal: $ip - $count falhas de auth" -ForegroundColor Green
        }
    }
} else {
    Write-Host "✓ Nenhuma falha de autenticação detectada" -ForegroundColor Green
}

Write-Host ""

# ============================================
# 4. Paths mais Atacados
# ============================================
Write-Host "=== PATHS MAIS ATACADOS (404s) ===" -ForegroundColor Yellow
Write-Host ""

$paths404 = $logContent | Where-Object { $_ -match '"status":404' } |
    ForEach-Object {
        if ($_ -match '"path":"([^"]+)"') {
            $matches[1]
        }
    } | Group-Object | Sort-Object Count -Descending | Select-Object -First 10

foreach ($pathGroup in $paths404) {
    $path = $pathGroup.Name
    $count = $pathGroup.Count
    Write-Host "  $count vezes - $path"
}

Write-Host ""

# ============================================
# 5. User-Agents Suspeitos
# ============================================
Write-Host "=== USER-AGENTS SUSPEITOS ===" -ForegroundColor Yellow
Write-Host ""

# User-Agents vazios
$emptyUA = $logContent | Where-Object { $_ -match '"status":404' -and $_ -match '"user_agent":""' }
$emptyUACount = $emptyUA.Count

if ($emptyUACount -gt 0) {
    Write-Host "⚠️  $emptyUACount requisições com User-Agent vazio" -ForegroundColor Red
}

# Listar User-Agents únicos em 404s
Write-Host ""
Write-Host "User-Agents em requisições 404:"

$userAgents = $logContent | Where-Object { $_ -match '"status":404' } |
    ForEach-Object {
        if ($_ -match '"user_agent":"([^"]*)"') {
            $ua = $matches[1]
            if ([string]::IsNullOrEmpty($ua)) {
                "(vazio)"
            } else {
                $ua
            }
        }
    } | Group-Object | Sort-Object Count -Descending | Select-Object -First 10

foreach ($uaGroup in $userAgents) {
    $ua = $uaGroup.Name
    $count = $uaGroup.Count
    Write-Host "  $count vezes - $ua"
}

Write-Host ""

# ============================================
# 6. Estatísticas Gerais
# ============================================
Write-Host "=== ESTATÍSTICAS GERAIS ===" -ForegroundColor Yellow
Write-Host ""

$totalRequests = ($logContent | Where-Object { $_ -match '"message":"Request"' }).Count
$total404 = ($logContent | Where-Object { $_ -match '"status":404' }).Count
$total401 = ($logContent | Where-Object { $_ -match '"status":401' }).Count
$total403 = ($logContent | Where-Object { $_ -match '"status":403' }).Count
$total500 = ($logContent | Where-Object { $_ -match '"status":5' }).Count

$uniqueIPs = $logContent | ForEach-Object {
    if ($_ -match '"ip":"([^"]+)"') {
        $matches[1]
    }
} | Select-Object -Unique

$uniqueIPCount = $uniqueIPs.Count

Write-Host "Total de requisições: $totalRequests"
Write-Host "IPs únicos: $uniqueIPCount"
Write-Host ""
Write-Host "Status HTTP:"
Write-Host "  404 (Not Found): $total404"
Write-Host "  401 (Unauthorized): $total401"
Write-Host "  403 (Forbidden/Banidos): $total403"
Write-Host "  5xx (Erros do servidor): $total500"

# Calcular taxa de 404
if ($totalRequests -gt 0) {
    $rate404 = [math]::Round(($total404 / $totalRequests) * 100, 2)
    Write-Host ""
    Write-Host "Taxa de 404: $rate404%"

    if ($rate404 -gt 20) {
        Write-Host "⚠️  ALERTA: Taxa de 404 muito alta! Possível scanning." -ForegroundColor Red
    } elseif ($rate404 -gt 10) {
        Write-Host "⚠️  Atenção: Taxa de 404 elevada." -ForegroundColor Yellow
    } else {
        Write-Host "✓ Taxa de 404 normal." -ForegroundColor Green
    }
}

Write-Host ""

# ============================================
# 7. Top 10 IPs Mais Ativos
# ============================================
Write-Host "=== TOP 10 IPs MAIS ATIVOS ===" -ForegroundColor Yellow
Write-Host ""

$topIPs = $logContent | ForEach-Object {
    if ($_ -match '"ip":"([^"]+)"') {
        $matches[1]
    }
} | Group-Object | Sort-Object Count -Descending | Select-Object -First 10

foreach ($ipGroup in $topIPs) {
    $ip = $ipGroup.Name
    $count = $ipGroup.Count
    Write-Host "  $count requisições - $ip"
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "  Análise concluída" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Sugestão de ações
if ($total403 -gt 0) {
    Write-Host "✓ Fail2Ban está funcionando - $total403 requisições bloqueadas" -ForegroundColor Green
}

if ($bannedCount -gt 5) {
    Write-Host "⚠️  Considere adicionar IPs banidos ao firewall permanentemente" -ForegroundColor Yellow
}

if ($rate404 -gt 20) {
    Write-Host "🚨 AÇÃO RECOMENDADA: Revisar configurações de segurança" -ForegroundColor Red
    Write-Host "   - Verificar se há scanning ativo"
    Write-Host "   - Considerar whitelist de IPs se API é interna"
    Write-Host "   - Habilitar HTTPS/TLS se ainda não estiver"
}

Write-Host ""
