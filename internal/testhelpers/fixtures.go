// Data de criação: 26/11/2025 14:35
// Versão: 3.0.0.6
// Fixtures com payloads de teste para benchmarks de webhook
package testhelpers

// Payloads JSON de exemplo do Zenvia para testes e benchmarks
const (
	// ZenviaEmailPayloadSent - Email enviado (status 121)
	ZenviaEmailPayloadSent = `{
		"id": "msg-123456789",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:30:00Z",
		"subscriptionId": "sub-webhook-email",
		"channel": "email",
		"direction": "OUT",
		"messageId": "msg-email-12345",
		"messageStatus": {
			"timestamp": "2025-11-26T14:30:00Z",
			"code": "SENT",
			"description": "Message sent to carrier"
		},
		"contents": [
			{
				"type": "text",
				"text": "Email de teste"
			}
		]
	}`

	// ZenviaEmailPayloadDelivered - Email entregue (status 122)
	ZenviaEmailPayloadDelivered = `{
		"id": "msg-123456790",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:31:00Z",
		"subscriptionId": "sub-webhook-email",
		"channel": "email",
		"direction": "OUT",
		"messageId": "msg-email-12346",
		"messageStatus": {
			"timestamp": "2025-11-26T14:31:00Z",
			"code": "DELIVERED",
			"description": "Message delivered to recipient"
		},
		"contents": [
			{
				"type": "text",
				"text": "Email de teste entregue"
			}
		]
	}`

	// ZenviaEmailPayloadRead - Email lido (status 123)
	ZenviaEmailPayloadRead = `{
		"id": "msg-123456791",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:32:00Z",
		"subscriptionId": "sub-webhook-email",
		"channel": "email",
		"direction": "OUT",
		"messageId": "msg-email-12347",
		"messageStatus": {
			"timestamp": "2025-11-26T14:32:00Z",
			"code": "READ",
			"description": "Message read by recipient"
		},
		"contents": [
			{
				"type": "text",
				"text": "Email de teste lido"
			}
		]
	}`

	// ZenviaEmailPayloadRejected - Email rejeitado (status 124)
	ZenviaEmailPayloadRejected = `{
		"id": "msg-123456792",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:33:00Z",
		"subscriptionId": "sub-webhook-email",
		"channel": "email",
		"direction": "OUT",
		"messageId": "msg-email-12348",
		"messageStatus": {
			"timestamp": "2025-11-26T14:33:00Z",
			"code": "NOT_SENT",
			"description": "Message rejected by carrier"
		},
		"contents": [
			{
				"type": "text",
				"text": "Email de teste rejeitado"
			}
		]
	}`

	// ZenviaEmailPayloadBounce - Email devolvido (status 125)
	ZenviaEmailPayloadBounce = `{
		"id": "msg-123456793",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:34:00Z",
		"subscriptionId": "sub-webhook-email",
		"channel": "email",
		"direction": "OUT",
		"messageId": "msg-email-12349",
		"messageStatus": {
			"timestamp": "2025-11-26T14:34:00Z",
			"code": "BOUNCED",
			"description": "Message bounced - invalid email"
		},
		"contents": [
			{
				"type": "text",
				"text": "Email de teste devolvido"
			}
		]
	}`

	// ZenviaSMSPayloadSent - SMS enviado (status 121)
	ZenviaSMSPayloadSent = `{
		"id": "msg-sms-123456789",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:30:00Z",
		"subscriptionId": "sub-webhook-sms",
		"channel": "sms",
		"direction": "OUT",
		"messageId": "msg-sms-54321",
		"messageStatus": {
			"timestamp": "2025-11-26T14:30:00Z",
			"code": "SENT",
			"description": "Message sent to carrier"
		},
		"contents": [
			{
				"type": "text",
				"text": "SMS de teste"
			}
		]
	}`

	// ZenviaSMSPayloadDelivered - SMS entregue (status 122)
	ZenviaSMSPayloadDelivered = `{
		"id": "msg-sms-123456790",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:31:00Z",
		"subscriptionId": "sub-webhook-sms",
		"channel": "sms",
		"direction": "OUT",
		"messageId": "msg-sms-54322",
		"messageStatus": {
			"timestamp": "2025-11-26T14:31:00Z",
			"code": "DELIVERED",
			"description": "Message delivered to recipient"
		},
		"contents": [
			{
				"type": "text",
				"text": "SMS de teste entregue"
			}
		]
	}`

	// ZenviaSMSPayloadRead - SMS lido (status 123)
	ZenviaSMSPayloadRead = `{
		"id": "msg-sms-123456791",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:32:00Z",
		"subscriptionId": "sub-webhook-sms",
		"channel": "sms",
		"direction": "OUT",
		"messageId": "msg-sms-54323",
		"messageStatus": {
			"timestamp": "2025-11-26T14:32:00Z",
			"code": "READ",
			"description": "Message read by recipient"
		},
		"contents": [
			{
				"type": "text",
				"text": "SMS de teste lido"
			}
		]
	}`

	// ZenviaSMSPayloadRejected - SMS rejeitado (status 124)
	ZenviaSMSPayloadRejected = `{
		"id": "msg-sms-123456792",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:33:00Z",
		"subscriptionId": "sub-webhook-sms",
		"channel": "sms",
		"direction": "OUT",
		"messageId": "msg-sms-54324",
		"messageStatus": {
			"timestamp": "2025-11-26T14:33:00Z",
			"code": "NOT_SENT",
			"description": "Message rejected by carrier"
		},
		"contents": [
			{
				"type": "text",
				"text": "SMS de teste rejeitado"
			}
		]
	}`

	// ZenviaSMSPayloadBounce - SMS devolvido (status 125)
	ZenviaSMSPayloadBounce = `{
		"id": "msg-sms-123456793",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:34:00Z",
		"subscriptionId": "sub-webhook-sms",
		"channel": "sms",
		"direction": "OUT",
		"messageId": "msg-sms-54325",
		"messageStatus": {
			"timestamp": "2025-11-26T14:34:00Z",
			"code": "BOUNCED",
			"description": "Message bounced - invalid number"
		},
		"contents": [
			{
				"type": "text",
				"text": "SMS de teste devolvido"
			}
		]
	}`

	// InvalidJSON - JSON inválido para testes de erro
	InvalidJSON = `{"id": "invalid", "type": "MESSAGE_STATUS", "timestamp": }`

	// EmptyJSON - JSON vazio
	EmptyJSON = `{}`

	// MissingMessageID - Payload sem messageId (campo obrigatório)
	MissingMessageID = `{
		"id": "msg-no-messageid",
		"type": "MESSAGE_STATUS",
		"timestamp": "2025-11-26T14:30:00Z",
		"subscriptionId": "sub-webhook-test",
		"channel": "email"
	}`
)

// ZenviaStatusCodes - Mapeamento de códigos Zenvia para códigos internos
var ZenviaStatusCodes = map[string]int{
	"SENT":      121, // Enviado
	"DELIVERED": 122, // Entregue
	"READ":      123, // Lido
	"NOT_SENT":  124, // Rejeitado
	"BOUNCED":   125, // Devolvido
}

// GetEmailPayloadByStatus retorna payload de email pelo código de status interno
func GetEmailPayloadByStatus(statusCode int) string {
	switch statusCode {
	case 121:
		return ZenviaEmailPayloadSent
	case 122:
		return ZenviaEmailPayloadDelivered
	case 123:
		return ZenviaEmailPayloadRead
	case 124:
		return ZenviaEmailPayloadRejected
	case 125:
		return ZenviaEmailPayloadBounce
	default:
		return ZenviaEmailPayloadSent
	}
}

// GetSMSPayloadByStatus retorna payload de SMS pelo código de status interno
func GetSMSPayloadByStatus(statusCode int) string {
	switch statusCode {
	case 121:
		return ZenviaSMSPayloadSent
	case 122:
		return ZenviaSMSPayloadDelivered
	case 123:
		return ZenviaSMSPayloadRead
	case 124:
		return ZenviaSMSPayloadRejected
	case 125:
		return ZenviaSMSPayloadBounce
	default:
		return ZenviaSMSPayloadSent
	}
}

// GetAllEmailPayloads retorna todos os payloads de email para testes completos
func GetAllEmailPayloads() []string {
	return []string{
		ZenviaEmailPayloadSent,
		ZenviaEmailPayloadDelivered,
		ZenviaEmailPayloadRead,
		ZenviaEmailPayloadRejected,
		ZenviaEmailPayloadBounce,
	}
}

// GetAllSMSPayloads retorna todos os payloads de SMS para testes completos
func GetAllSMSPayloads() []string {
	return []string{
		ZenviaSMSPayloadSent,
		ZenviaSMSPayloadDelivered,
		ZenviaSMSPayloadRead,
		ZenviaSMSPayloadRejected,
		ZenviaSMSPayloadBounce,
	}
}

// GetInvalidPayloads retorna payloads inválidos para testes de erro
func GetInvalidPayloads() []string {
	return []string{
		InvalidJSON,
		EmptyJSON,
		MissingMessageID,
	}
}
