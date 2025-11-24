#!/bin/bash
# Monitor em Tempo Real - WSICRMREST
# Monitora logs continuamente e alerta sobre eventos de segurança
# Data: 2025-11-24
# Uso: ./scripts/watch_security.sh

# Cores
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

LOG_DIR="log"
TODAY=$(date +%Y-%m-%d)
LOG_FILE="$LOG_DIR/wsicrmrest_$TODAY.log"

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}  WSICRMREST - Monitor em Tempo Real${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""
echo "Monitorando: $LOG_FILE"
echo "Pressione Ctrl+C para sair"
echo ""

# Verificar se arquivo existe
if [ ! -f "$LOG_FILE" ]; then
    echo -e "${RED}Aguardando criação do arquivo de log...${NC}"
    while [ ! -f "$LOG_FILE" ]; do
        sleep 1
    done
    echo -e "${GREEN}✓ Arquivo de log criado${NC}"
fi

# Monitorar log em tempo real
tail -f "$LOG_FILE" | while read line; do
    # Detectar IPs banidos
    if echo "$line" | grep -q "IP BANIDO"; then
        ip=$(echo "$line" | grep -oP '"ip":"[^"]*"' | sed 's/"ip":"//;s/"//')
        reason=$(echo "$line" | grep -oP 'por [^"]*' || echo "razão desconhecida")
        echo -e "${RED}🚨 [$(date '+%H:%M:%S')] IP BANIDO: $ip ($reason)${NC}"

    # Detectar tentativas banidas (403)
    elif echo "$line" | grep -q '"status":403'; then
        ip=$(echo "$line" | grep -oP '"ip":"[^"]*"' | sed 's/"ip":"//;s/"//')
        path=$(echo "$line" | grep -oP '"path":"[^"]*"' | sed 's/"path":"//;s/"//')
        echo -e "${RED}🔒 [$(date '+%H:%M:%S')] BLOQUEADO: $ip tentou acessar $path${NC}"

    # Detectar múltiplos 404s do mesmo IP
    elif echo "$line" | grep -q '"status":404'; then
        ip=$(echo "$line" | grep -oP '"ip":"[^"]*"' | sed 's/"ip":"//;s/"//')
        path=$(echo "$line" | grep -oP '"path":"[^"]*"' | sed 's/"path":"//;s/"//')
        echo -e "${YELLOW}⚠️  [$(date '+%H:%M:%S')] 404: $ip -> $path${NC}"

    # Detectar falhas de autenticação
    elif echo "$line" | grep -q '"status":401'; then
        ip=$(echo "$line" | grep -oP '"ip":"[^"]*"' | sed 's/"ip":"//;s/"//')
        path=$(echo "$line" | grep -oP '"path":"[^"]*"' | sed 's/"path":"//;s/"//')
        echo -e "${YELLOW}🔑 [$(date '+%H:%M:%S')] AUTH FALHA: $ip -> $path${NC}"

    # Detectar erros 500
    elif echo "$line" | grep -q '"status":5'; then
        path=$(echo "$line" | grep -oP '"path":"[^"]*"' | sed 's/"path":"//;s/"//')
        status=$(echo "$line" | grep -oP '"status":[0-9]+' | sed 's/"status"://')
        echo -e "${RED}❌ [$(date '+%H:%M:%S')] ERRO $status: $path${NC}"

    # Requisições normais (apenas mostrar de vez em quando)
    elif echo "$line" | grep -q '"status":200'; then
        # Mostrar 1 a cada 10 requisições 200
        if [ $((RANDOM % 10)) -eq 0 ]; then
            path=$(echo "$line" | grep -oP '"path":"[^"]*"' | sed 's/"path":"//;s/"//')
            echo -e "${GREEN}✓ [$(date '+%H:%M:%S')] OK: $path${NC}"
        fi
    fi
done
