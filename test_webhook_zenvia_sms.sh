#!/bin/bash

# Script para testar o webhook Zenvia SMS
# Endpoint: POST /webhook/zenvia/sms

BASE_URL="http://localhost:8080"

echo "======================================"
echo "Testando Webhook Zenvia SMS"
echo "======================================"
echo ""

# Teste 1: Evento SENT (Agendado)
echo "Teste 1: Evento SENT (Agendado)"
curl -X POST "$BASE_URL/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999998888",
      "externalId": "54321"
    },
    "messageStatus": {
      "code": "SENT",
      "description": "Mensagem enviada ao provedor de SMS"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 2: Evento DELIVERED (Entregue)
echo "Teste 2: Evento DELIVERED (Entregue)"
curl -X POST "$BASE_URL/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999998888",
      "externalId": "54321"
    },
    "messageStatus": {
      "code": "DELIVERED",
      "description": "SMS entregue com sucesso"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 3: Evento REJECTED (Rejeitado)
echo "Teste 3: Evento REJECTED (Rejeitado)"
curl -X POST "$BASE_URL/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999997777",
      "externalId": "54322"
    },
    "messageStatus": {
      "code": "REJECTED",
      "description": "SMS rejeitado pela operadora",
      "causes": [
        {
          "reason": "INVALID_NUMBER",
          "details": "Número de telefone inválido ou inexistente"
        }
      ]
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 4: Evento NOT_DELIVERED (Não Entregue)
echo "Teste 4: Evento NOT_DELIVERED (Não Entregue)"
curl -X POST "$BASE_URL/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999996666",
      "externalId": "54323"
    },
    "messageStatus": {
      "code": "NOT_DELIVERED",
      "description": "SMS não entregue",
      "causes": [
        {
          "reason": "CARRIER_BLOCKED",
          "details": "Número bloqueado pela operadora"
        }
      ]
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 5: Tipo de callback inválido
echo "Teste 5: Tipo de callback inválido"
curl -X POST "$BASE_URL/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "OTHER_TYPE",
    "message": {
      "to": "+5511999998888",
      "externalId": "54321"
    },
    "messageStatus": {
      "code": "SENT",
      "description": "Mensagem enviada"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 6: ID não numérico
echo "Teste 6: ID não numérico (formato inválido)"
curl -X POST "$BASE_URL/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999998888",
      "externalId": "abc123xyz"
    },
    "messageStatus": {
      "code": "SENT",
      "description": "Mensagem enviada"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 7: Status não reconhecido
echo "Teste 7: Status não reconhecido (fora do escopo)"
curl -X POST "$BASE_URL/webhook/zenvia/sms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "+5511999998888",
      "externalId": "54321"
    },
    "messageStatus": {
      "code": "UNKNOWN_STATUS",
      "description": "Status desconhecido"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

echo "Testes concluídos!"
