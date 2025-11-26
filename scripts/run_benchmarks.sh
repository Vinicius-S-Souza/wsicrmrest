#!/bin/bash
# Data de criação: 26/11/2025 14:45
# Versão: 3.0.0.6
# Script para executar benchmarks dos webhooks Zenvia

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Diretório do projeto (assumindo que o script está em scripts/)
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR" || exit 1

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}WSICRMREST - Benchmarks de Webhook${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${RED}Erro: Go não está instalado${NC}"
    exit 1
fi

echo -e "${GREEN}Go version:${NC} $(go version)"
echo ""

# Criar diretório de resultados se não existir
RESULTS_DIR="$PROJECT_DIR/benchmark_results"
mkdir -p "$RESULTS_DIR"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="$RESULTS_DIR/benchmark_$TIMESTAMP.txt"
RESULTS_MEM="$RESULTS_DIR/benchmark_mem_$TIMESTAMP.txt"

# Função para executar benchmark específico
run_benchmark() {
    local benchmark_name=$1
    local description=$2

    echo -e "${YELLOW}Executando: $description${NC}"
    go test -bench="$benchmark_name" -benchmem -benchtime=3s \
        ./internal/handlers/ >> "$RESULTS_FILE" 2>&1

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Concluído${NC}"
    else
        echo -e "${RED}✗ Falhou${NC}"
    fi
    echo ""
}

echo -e "${BLUE}Resultados serão salvos em:${NC}"
echo -e "  ${RESULTS_FILE}"
echo -e "  ${RESULTS_MEM}"
echo ""

# Cabeçalho do arquivo de resultados
{
    echo "=========================================="
    echo "WSICRMREST - Benchmark Results"
    echo "Data: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "Go Version: $(go version)"
    echo "=========================================="
    echo ""
} > "$RESULTS_FILE"

# Menu de opções
echo -e "${BLUE}Selecione o tipo de benchmark:${NC}"
echo "1) Email - Todos os status (Sent, Delivered, Read, Rejected, Bounce)"
echo "2) SMS - Todos os status (Sent, Delivered, Read, Rejected, Bounce)"
echo "3) Email e SMS - Status individuais"
echo "4) Benchmarks paralelos (Email e SMS)"
echo "5) Benchmarks com payloads grandes"
echo "6) Benchmarks com JSON inválido"
echo "7) Todos os benchmarks (completo)"
echo "8) Benchmark rápido (apenas Sent)"
echo ""
read -p "Opção [1-8]: " option

case $option in
    1)
        echo -e "${BLUE}Executando benchmarks de Email...${NC}"
        echo ""
        run_benchmark "BenchmarkZenviaEmailWebhook_Sent" "Email - Status SENT (121)"
        run_benchmark "BenchmarkZenviaEmailWebhook_Delivered" "Email - Status DELIVERED (122)"
        run_benchmark "BenchmarkZenviaEmailWebhook_Read" "Email - Status READ (123)"
        run_benchmark "BenchmarkZenviaEmailWebhook_Rejected" "Email - Status REJECTED (124)"
        run_benchmark "BenchmarkZenviaEmailWebhook_Bounce" "Email - Status BOUNCE (125)"
        run_benchmark "BenchmarkZenviaEmailWebhook_AllStatuses" "Email - Todos os status alternados"
        ;;

    2)
        echo -e "${BLUE}Executando benchmarks de SMS...${NC}"
        echo ""
        run_benchmark "BenchmarkZenviaSMSWebhook_Sent" "SMS - Status SENT (121)"
        run_benchmark "BenchmarkZenviaSMSWebhook_Delivered" "SMS - Status DELIVERED (122)"
        run_benchmark "BenchmarkZenviaSMSWebhook_Read" "SMS - Status READ (123)"
        run_benchmark "BenchmarkZenviaSMSWebhook_Rejected" "SMS - Status REJECTED (124)"
        run_benchmark "BenchmarkZenviaSMSWebhook_Bounce" "SMS - Status BOUNCE (125)"
        run_benchmark "BenchmarkZenviaSMSWebhook_AllStatuses" "SMS - Todos os status alternados"
        ;;

    3)
        echo -e "${BLUE}Executando benchmarks individuais...${NC}"
        echo ""
        run_benchmark "BenchmarkZenviaEmailWebhook_Sent" "Email - SENT"
        run_benchmark "BenchmarkZenviaEmailWebhook_Delivered" "Email - DELIVERED"
        run_benchmark "BenchmarkZenviaEmailWebhook_Read" "Email - READ"
        run_benchmark "BenchmarkZenviaEmailWebhook_Rejected" "Email - REJECTED"
        run_benchmark "BenchmarkZenviaEmailWebhook_Bounce" "Email - BOUNCE"
        run_benchmark "BenchmarkZenviaSMSWebhook_Sent" "SMS - SENT"
        run_benchmark "BenchmarkZenviaSMSWebhook_Delivered" "SMS - DELIVERED"
        run_benchmark "BenchmarkZenviaSMSWebhook_Read" "SMS - READ"
        run_benchmark "BenchmarkZenviaSMSWebhook_Rejected" "SMS - REJECTED"
        run_benchmark "BenchmarkZenviaSMSWebhook_Bounce" "SMS - BOUNCE"
        ;;

    4)
        echo -e "${BLUE}Executando benchmarks paralelos...${NC}"
        echo ""
        run_benchmark "BenchmarkZenviaEmailWebhook_Parallel" "Email - Benchmark Paralelo"
        run_benchmark "BenchmarkZenviaSMSWebhook_Parallel" "SMS - Benchmark Paralelo"
        ;;

    5)
        echo -e "${BLUE}Executando benchmarks com payloads grandes...${NC}"
        echo ""
        run_benchmark "BenchmarkZenviaEmailWebhook_LargePayload" "Email - Payload Grande"
        run_benchmark "BenchmarkZenviaSMSWebhook_LargePayload" "SMS - Payload Grande"
        ;;

    6)
        echo -e "${BLUE}Executando benchmarks com JSON inválido...${NC}"
        echo ""
        run_benchmark "BenchmarkZenviaEmailWebhook_InvalidJSON" "Email - JSON Inválido"
        run_benchmark "BenchmarkZenviaSMSWebhook_InvalidJSON" "SMS - JSON Inválido"
        ;;

    7)
        echo -e "${BLUE}Executando TODOS os benchmarks (pode demorar)...${NC}"
        echo ""

        # Email benchmarks
        echo -e "${YELLOW}=== BENCHMARKS DE EMAIL ===${NC}"
        run_benchmark "BenchmarkZenviaEmailWebhook_Sent" "Email - SENT"
        run_benchmark "BenchmarkZenviaEmailWebhook_Delivered" "Email - DELIVERED"
        run_benchmark "BenchmarkZenviaEmailWebhook_Read" "Email - READ"
        run_benchmark "BenchmarkZenviaEmailWebhook_Rejected" "Email - REJECTED"
        run_benchmark "BenchmarkZenviaEmailWebhook_Bounce" "Email - BOUNCE"
        run_benchmark "BenchmarkZenviaEmailWebhook_AllStatuses" "Email - Todos os status"
        run_benchmark "BenchmarkZenviaEmailWebhook_InvalidJSON" "Email - JSON Inválido"
        run_benchmark "BenchmarkZenviaEmailWebhook_Parallel" "Email - Paralelo"
        run_benchmark "BenchmarkZenviaEmailWebhook_LargePayload" "Email - Payload Grande"

        # SMS benchmarks
        echo -e "${YELLOW}=== BENCHMARKS DE SMS ===${NC}"
        run_benchmark "BenchmarkZenviaSMSWebhook_Sent" "SMS - SENT"
        run_benchmark "BenchmarkZenviaSMSWebhook_Delivered" "SMS - DELIVERED"
        run_benchmark "BenchmarkZenviaSMSWebhook_Read" "SMS - READ"
        run_benchmark "BenchmarkZenviaSMSWebhook_Rejected" "SMS - REJECTED"
        run_benchmark "BenchmarkZenviaSMSWebhook_Bounce" "SMS - BOUNCE"
        run_benchmark "BenchmarkZenviaSMSWebhook_AllStatuses" "SMS - Todos os status"
        run_benchmark "BenchmarkZenviaSMSWebhook_InvalidJSON" "SMS - JSON Inválido"
        run_benchmark "BenchmarkZenviaSMSWebhook_Parallel" "SMS - Paralelo"
        run_benchmark "BenchmarkZenviaSMSWebhook_LargePayload" "SMS - Payload Grande"
        ;;

    8)
        echo -e "${BLUE}Executando benchmark rápido...${NC}"
        echo ""
        run_benchmark "BenchmarkZenviaEmailWebhook_Sent" "Email - SENT (Quick)"
        run_benchmark "BenchmarkZenviaSMSWebhook_Sent" "SMS - SENT (Quick)"
        ;;

    *)
        echo -e "${RED}Opção inválida!${NC}"
        exit 1
        ;;
esac

# Gerar relatório de memória separado
echo -e "${BLUE}Gerando relatório de uso de memória...${NC}"
go test -bench=. -benchmem ./internal/handlers/ 2>&1 | grep -E "Benchmark|alloc" > "$RESULTS_MEM"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Benchmarks concluídos!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Resultados salvos em:${NC}"
echo -e "  ${RESULTS_FILE}"
echo -e "  ${RESULTS_MEM}"
echo ""

# Mostrar resumo dos últimos resultados
if [ -f "$RESULTS_FILE" ]; then
    echo -e "${BLUE}Resumo dos resultados:${NC}"
    echo ""
    tail -n 50 "$RESULTS_FILE"
    echo ""
fi

# Sugestão de próximos passos
echo -e "${YELLOW}Próximos passos:${NC}"
echo "1. Analise os resultados em: $RESULTS_FILE"
echo "2. Compare com benchmarks anteriores em: $RESULTS_DIR"
echo "3. Consulte docs/BENCHMARKS.md para interpretação dos resultados"
echo ""
