#!/bin/bash
# Monitor de Segurança - WSICRMREST
# Analisa logs em tempo real para detectar ataques
# Data: 2025-11-24

# Cores para output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Configurações
LOG_DIR="log"
ALERT_THRESHOLD_404=10
ALERT_THRESHOLD_AUTH=5
TODAY=$(date +%Y-%m-%d)
LOG_FILE="$LOG_DIR/wsicrmrest_$TODAY.log"

echo "========================================="
echo "  WSICRMREST - Monitor de Segurança"
echo "========================================="
echo ""
echo "Analisando: $LOG_FILE"
echo "Data: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Verificar se arquivo de log existe
if [ ! -f "$LOG_FILE" ]; then
    echo -e "${RED}ERRO: Arquivo de log não encontrado: $LOG_FILE${NC}"
    echo "Verifique se o servidor está rodando e gerando logs."
    exit 1
fi

# ============================================
# 1. IPs com múltiplos 404s (Scanning)
# ============================================
echo -e "${YELLOW}=== IPs SUSPEITOS (Múltiplos 404s) ===${NC}"
echo ""

grep '"status":404' "$LOG_FILE" 2>/dev/null | \
    grep -oP '"ip":"[^"]*"' | \
    sed 's/"ip":"//;s/"//' | \
    sort | uniq -c | sort -rn | \
    while read count ip; do
        if [ "$count" -gt "$ALERT_THRESHOLD_404" ]; then
            echo -e "${RED}🚨 ALERTA: $ip - $count tentativas 404${NC}"
        elif [ "$count" -gt 5 ]; then
            echo -e "${YELLOW}⚠️  ATENÇÃO: $ip - $count tentativas 404${NC}"
        else
            echo -e "${GREEN}✓ Normal: $ip - $count tentativas 404${NC}"
        fi
    done

echo ""

# ============================================
# 2. IPs Banidos pelo Fail2Ban
# ============================================
echo -e "${YELLOW}=== IPs BANIDOS (Fail2Ban) ===${NC}"
echo ""

BANNED_COUNT=$(grep -c "IP BANIDO" "$LOG_FILE" 2>/dev/null)

if [ "$BANNED_COUNT" -gt 0 ]; then
    echo -e "${RED}Total de IPs banidos hoje: $BANNED_COUNT${NC}"
    echo ""

    grep "IP BANIDO" "$LOG_FILE" | \
        grep -oP '"ip":"[^"]*"' | \
        sed 's/"ip":"//;s/"//' | \
        sort | uniq -c | sort -rn | \
        while read count ip; do
            echo -e "${RED}  🔒 $ip - banido $count vez(es)${NC}"
        done
else
    echo -e "${GREEN}✓ Nenhum IP foi banido hoje${NC}"
fi

echo ""

# ============================================
# 3. Falhas de Autenticação (401)
# ============================================
echo -e "${YELLOW}=== FALHAS DE AUTENTICAÇÃO (401) ===${NC}"
echo ""

AUTH_FAILURES=$(grep '"status":401' "$LOG_FILE" 2>/dev/null | wc -l)

if [ "$AUTH_FAILURES" -gt 0 ]; then
    echo "Total de falhas de autenticação: $AUTH_FAILURES"
    echo ""

    grep '"status":401' "$LOG_FILE" 2>/dev/null | \
        grep -oP '"ip":"[^"]*"' | \
        sed 's/"ip":"//;s/"//' | \
        sort | uniq -c | sort -rn | \
        while read count ip; do
            if [ "$count" -gt "$ALERT_THRESHOLD_AUTH" ]; then
                echo -e "${RED}🚨 ALERTA: $ip - $count falhas de auth${NC}"
            elif [ "$count" -gt 2 ]; then
                echo -e "${YELLOW}⚠️  ATENÇÃO: $ip - $count falhas de auth${NC}"
            else
                echo -e "${GREEN}✓ Normal: $ip - $count falhas de auth${NC}"
            fi
        done
else
    echo -e "${GREEN}✓ Nenhuma falha de autenticação detectada${NC}"
fi

echo ""

# ============================================
# 4. Paths mais Atacados
# ============================================
echo -e "${YELLOW}=== PATHS MAIS ATACADOS (404s) ===${NC}"
echo ""

grep '"status":404' "$LOG_FILE" 2>/dev/null | \
    grep -oP '"path":"[^"]*"' | \
    sed 's/"path":"//;s/"//' | \
    sort | uniq -c | sort -rn | head -10 | \
    while read count path; do
        echo "  $count vezes - $path"
    done

echo ""

# ============================================
# 5. User-Agents Suspeitos
# ============================================
echo -e "${YELLOW}=== USER-AGENTS SUSPEITOS ===${NC}"
echo ""

# User-Agents vazios ou suspeitos
SUSPICIOUS_UA=$(grep '"status":404' "$LOG_FILE" 2>/dev/null | \
    grep -E '"user_agent":""' | wc -l)

if [ "$SUSPICIOUS_UA" -gt 0 ]; then
    echo -e "${RED}⚠️  $SUSPICIOUS_UA requisições com User-Agent vazio${NC}"
fi

# Listar User-Agents únicos em 404s
echo ""
echo "User-Agents em requisições 404:"
grep '"status":404' "$LOG_FILE" 2>/dev/null | \
    grep -oP '"user_agent":"[^"]*"' | \
    sed 's/"user_agent":"//;s/"//' | \
    sort | uniq -c | sort -rn | head -10 | \
    while read count ua; do
        if [ -z "$ua" ]; then
            ua="(vazio)"
        fi
        echo "  $count vezes - $ua"
    done

echo ""

# ============================================
# 6. Estatísticas Gerais
# ============================================
echo -e "${YELLOW}=== ESTATÍSTICAS GERAIS ===${NC}"
echo ""

TOTAL_REQUESTS=$(grep -c '"message":"Request"' "$LOG_FILE" 2>/dev/null)
TOTAL_404=$(grep -c '"status":404' "$LOG_FILE" 2>/dev/null)
TOTAL_401=$(grep -c '"status":401' "$LOG_FILE" 2>/dev/null)
TOTAL_403=$(grep -c '"status":403' "$LOG_FILE" 2>/dev/null)
TOTAL_500=$(grep -c '"status":5' "$LOG_FILE" 2>/dev/null)
UNIQUE_IPS=$(grep '"ip":"' "$LOG_FILE" 2>/dev/null | \
    grep -oP '"ip":"[^"]*"' | sed 's/"ip":"//;s/"//' | sort -u | wc -l)

echo "Total de requisições: $TOTAL_REQUESTS"
echo "IPs únicos: $UNIQUE_IPS"
echo ""
echo "Status HTTP:"
echo "  404 (Not Found): $TOTAL_404"
echo "  401 (Unauthorized): $TOTAL_401"
echo "  403 (Forbidden/Banidos): $TOTAL_403"
echo "  5xx (Erros do servidor): $TOTAL_500"

# Calcular taxa de 404
if [ "$TOTAL_REQUESTS" -gt 0 ]; then
    RATE_404=$(awk "BEGIN {printf \"%.2f\", ($TOTAL_404/$TOTAL_REQUESTS)*100}")
    echo ""
    echo "Taxa de 404: $RATE_404%"

    if (( $(echo "$RATE_404 > 20" | bc -l) )); then
        echo -e "${RED}⚠️  ALERTA: Taxa de 404 muito alta! Possível scanning.${NC}"
    elif (( $(echo "$RATE_404 > 10" | bc -l) )); then
        echo -e "${YELLOW}⚠️  Atenção: Taxa de 404 elevada.${NC}"
    else
        echo -e "${GREEN}✓ Taxa de 404 normal.${NC}"
    fi
fi

echo ""

# ============================================
# 7. Países de Origem (Top 10 IPs)
# ============================================
echo -e "${YELLOW}=== TOP 10 IPs MAIS ATIVOS ===${NC}"
echo ""

grep '"ip":"' "$LOG_FILE" 2>/dev/null | \
    grep -oP '"ip":"[^"]*"' | \
    sed 's/"ip":"//;s/"//' | \
    sort | uniq -c | sort -rn | head -10 | \
    while read count ip; do
        echo "  $count requisições - $ip"
    done

echo ""
echo "========================================="
echo "  Análise concluída"
echo "========================================="
echo ""

# Sugestão de ações
if [ "$TOTAL_403" -gt 0 ]; then
    echo -e "${GREEN}✓ Fail2Ban está funcionando - $TOTAL_403 requisições bloqueadas${NC}"
fi

if [ "$BANNED_COUNT" -gt 5 ]; then
    echo -e "${YELLOW}⚠️  Considere adicionar IPs banidos ao firewall permanentemente${NC}"
fi

if (( $(echo "$RATE_404 > 20" | bc -l) )); then
    echo -e "${RED}🚨 AÇÃO RECOMENDADA: Revisar configurações de segurança${NC}"
    echo "   - Verificar se há scanning ativo"
    echo "   - Considerar whitelist de IPs se API é interna"
    echo "   - Habilitar HTTPS/TLS se ainda não estiver"
fi

echo ""
