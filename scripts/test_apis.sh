#!/bin/bash

# Script de teste das APIs do WSICRMREST
# Ajuste as variáveis conforme necessário

BASE_URL="http://localhost:8080"
CLIENT_ID="CLIENTE1234567890"
CLIENT_SECRET="a1234567890b"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "======================================"
echo "  Testando APIs WSICRMREST"
echo "======================================"
echo ""

# Teste 1: WSTest - Teste de conexão
echo -e "${YELLOW}[1] Testando GET /connect/v1/wsteste${NC}"
echo "URL: ${BASE_URL}/connect/v1/wsteste"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" "${BASE_URL}/connect/v1/wsteste")
HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

if [ "$HTTP_CODE" -eq 200 ]; then
    echo -e "${GREEN}✓ Sucesso (HTTP $HTTP_CODE)${NC}"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
    echo -e "${RED}✗ Erro (HTTP $HTTP_CODE)${NC}"
    echo "$BODY"
fi

echo ""
echo "======================================"
echo ""

# Teste 2: Gerar Token
echo -e "${YELLOW}[2] Testando GET /connect/v1/token${NC}"
echo "URL: ${BASE_URL}/connect/v1/token"
echo "Client ID: ${CLIENT_ID}"
echo "Client Secret: ${CLIENT_SECRET}"
echo ""

# Criar Basic Auth
AUTH=$(echo -n "${CLIENT_ID}:${CLIENT_SECRET}" | base64)

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  -H "Authorization: Basic ${AUTH}" \
  -H "Grant_type: client_credentials" \
  "${BASE_URL}/connect/v1/token")

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

if [ "$HTTP_CODE" -eq 200 ]; then
    echo -e "${GREEN}✓ Sucesso (HTTP $HTTP_CODE)${NC}"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"

    # Extrair token para uso futuro
    TOKEN=$(echo "$BODY" | jq -r '.access_token' 2>/dev/null)
    if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
        echo ""
        echo -e "${GREEN}Token gerado com sucesso!${NC}"
        echo "Token: ${TOKEN:0:50}..."
    fi
else
    echo -e "${RED}✗ Erro (HTTP $HTTP_CODE)${NC}"
    echo "$BODY"
fi

echo ""
echo "======================================"
echo ""
echo "Testes concluídos!"
