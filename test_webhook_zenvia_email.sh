#!/bin/bash

# Script para testar o webhook Zenvia Email
# Endpoint: POST /webhook/zenvia/email

BASE_URL="http://localhost:8080"

echo "======================================"
echo "Testando Webhook Zenvia Email"
echo "======================================"
echo ""

# Teste 1: Evento SENT (Agendado)
echo "Teste 1: Evento SENT (Agendado)"
curl -X POST "$BASE_URL/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "cliente@example.com",
      "externalId": "12345"
    },
    "messageStatus": {
      "code": "SENT",
      "description": "Mensagem enviada ao provedor de email"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 2: Evento DELIVERED (Entregue)
echo "Teste 2: Evento DELIVERED (Entregue)"
curl -X POST "$BASE_URL/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "cliente@example.com",
      "externalId": "12345"
    },
    "messageStatus": {
      "code": "DELIVERED",
      "description": "Email entregue com sucesso"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 3: Evento READ (Aberto)
echo "Teste 3: Evento READ (Aberto)"
curl -X POST "$BASE_URL/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "cliente@example.com",
      "externalId": "12345"
    },
    "messageStatus": {
      "code": "READ",
      "description": "Email foi aberto pelo destinatário"
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 4: Evento REJECTED (Rejeitado)
echo "Teste 4: Evento REJECTED (Rejeitado)"
curl -X POST "$BASE_URL/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "invalido@example.com",
      "externalId": "12345"
    },
    "messageStatus": {
      "code": "REJECTED",
      "description": "Email rejeitado pelo servidor",
      "causes": [
        {
          "reason": "INVALID_EMAIL",
          "details": "Endereço de email inválido ou inexistente"
        }
      ]
    }
  }' | jq '.'
echo ""
echo "======================================"
echo ""

# Teste 5: Tipo de callback inválido
echo "Teste 5: Tipo de callback inválido"
curl -X POST "$BASE_URL/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "OTHER_TYPE",
    "message": {
      "to": "cliente@example.com",
      "externalId": "12345"
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
curl -X POST "$BASE_URL/webhook/zenvia/email" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token_exemplo" \
  -d '{
    "type": "MESSAGE_STATUS",
    "message": {
      "to": "cliente@example.com",
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

echo "Testes concluídos!"
