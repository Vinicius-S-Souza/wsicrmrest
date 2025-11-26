// Data de criação: 26/11/2025 14:55
// Versão: 3.0.0.6
// Benchmarks para webhooks Zenvia (email e SMS)
// NÃO conecta ao banco de dados real - apenas testa performance de parsing JSON
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"wsicrmrest/internal/models"

	"github.com/gin-gonic/gin"
)

// Setup do Gin para testes
func init() {
	gin.SetMode(gin.TestMode)
}

// BenchmarkZenviaEmailWebhook_JSONParsing - Benchmark de parsing JSON para email
func BenchmarkZenviaEmailWebhook_JSONParsing(b *testing.B) {
	payload := `{
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
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal([]byte(payload), &webhook)
	}
}

// BenchmarkZenviaSMSWebhook_JSONParsing - Benchmark de parsing JSON para SMS
func BenchmarkZenviaSMSWebhook_JSONParsing(b *testing.B) {
	payload := `{
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
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal([]byte(payload), &webhook)
	}
}

// BenchmarkZenviaEmailWebhook_HTTPRequest - Benchmark de requisição HTTP completa
func BenchmarkZenviaEmailWebhook_HTTPRequest(b *testing.B) {
	payload := []byte(`{
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
		}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/webhook/zenvia/email", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		// Simular parse do body
		var webhook models.ZenviaWebhookRequest
		_ = json.NewDecoder(req.Body).Decode(&webhook)

		// Simular resposta
		w.WriteHeader(200)
		w.Write([]byte(`{"success":true}`))
	}
}

// BenchmarkZenviaSMSWebhook_HTTPRequest - Benchmark de requisição HTTP completa
func BenchmarkZenviaSMSWebhook_HTTPRequest(b *testing.B) {
	payload := []byte(`{
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
		}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/webhook/zenvia/sms", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		// Simular parse do body
		var webhook models.ZenviaWebhookRequest
		_ = json.NewDecoder(req.Body).Decode(&webhook)

		// Simular resposta
		w.WriteHeader(200)
		w.Write([]byte(`{"success":true}`))
	}
}

// BenchmarkZenviaEmailWebhook_MultipleStatuses - Benchmark com diferentes status
func BenchmarkZenviaEmailWebhook_MultipleStatuses(b *testing.B) {
	statuses := []string{"SENT", "DELIVERED", "READ", "NOT_SENT", "BOUNCED"}
	payloads := make([][]byte, len(statuses))

	for i, status := range statuses {
		payloads[i] = []byte(`{
			"id": "msg-` + status + `",
			"type": "MESSAGE_STATUS",
			"messageId": "msg-email-12345",
			"messageStatus": {
				"timestamp": "2025-11-26T14:30:00Z",
				"code": "` + status + `",
				"description": "Status ` + status + `"
			}
		}`)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		payload := payloads[i%len(payloads)]
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal(payload, &webhook)
	}
}

// BenchmarkZenviaSMSWebhook_MultipleStatuses - Benchmark com diferentes status
func BenchmarkZenviaSMSWebhook_MultipleStatuses(b *testing.B) {
	statuses := []string{"SENT", "DELIVERED", "READ", "NOT_SENT", "BOUNCED"}
	payloads := make([][]byte, len(statuses))

	for i, status := range statuses {
		payloads[i] = []byte(`{
			"id": "msg-` + status + `",
			"type": "MESSAGE_STATUS",
			"messageId": "msg-sms-54321",
			"messageStatus": {
				"timestamp": "2025-11-26T14:30:00Z",
				"code": "` + status + `",
				"description": "Status ` + status + `"
			}
		}`)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		payload := payloads[i%len(payloads)]
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal(payload, &webhook)
	}
}

// BenchmarkZenviaEmailWebhook_InvalidJSON - Benchmark com JSON inválido
func BenchmarkZenviaEmailWebhook_InvalidJSON(b *testing.B) {
	invalidPayload := []byte(`{"id": "invalid", "type": "MESSAGE_STATUS", "timestamp": }`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal(invalidPayload, &webhook) // Vai falhar, mas queremos medir o tempo
	}
}

// BenchmarkZenviaSMSWebhook_InvalidJSON - Benchmark com JSON inválido
func BenchmarkZenviaSMSWebhook_InvalidJSON(b *testing.B) {
	invalidPayload := []byte(`{"id": "invalid", "type": "MESSAGE_STATUS", "timestamp": }`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal(invalidPayload, &webhook) // Vai falhar, mas queremos medir o tempo
	}
}

// BenchmarkZenviaEmailWebhook_LargePayload - Benchmark com payload grande
func BenchmarkZenviaEmailWebhook_LargePayload(b *testing.B) {
	// Criar payload grande com conteúdo repetido
	largeContent := bytes.Repeat([]byte("Large content data "), 100)
	payload := []byte(`{
		"id": "msg-large",
		"type": "MESSAGE_STATUS",
		"messageId": "msg-email-12345",
		"messageStatus": {
			"timestamp": "2025-11-26T14:30:00Z",
			"code": "SENT",
			"description": "` + string(largeContent) + `"
		}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal(payload, &webhook)
	}
}

// BenchmarkZenviaSMSWebhook_LargePayload - Benchmark com payload grande
func BenchmarkZenviaSMSWebhook_LargePayload(b *testing.B) {
	// Criar payload grande com conteúdo repetido
	largeContent := bytes.Repeat([]byte("Large content data "), 100)
	payload := []byte(`{
		"id": "msg-large",
		"type": "MESSAGE_STATUS",
		"messageId": "msg-sms-54321",
		"messageStatus": {
			"timestamp": "2025-11-26T14:30:00Z",
			"code": "SENT",
			"description": "` + string(largeContent) + `"
		}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var webhook models.ZenviaWebhookRequest
		_ = json.Unmarshal(payload, &webhook)
	}
}

// BenchmarkZenviaEmailWebhook_Parallel - Benchmark paralelo
func BenchmarkZenviaEmailWebhook_Parallel(b *testing.B) {
	payload := []byte(`{
		"id": "msg-parallel",
		"type": "MESSAGE_STATUS",
		"messageId": "msg-email-12345",
		"messageStatus": {
			"timestamp": "2025-11-26T14:30:00Z",
			"code": "SENT",
			"description": "Parallel test"
		}
	}`)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var webhook models.ZenviaWebhookRequest
			_ = json.Unmarshal(payload, &webhook)
		}
	})
}

// BenchmarkZenviaSMSWebhook_Parallel - Benchmark paralelo
func BenchmarkZenviaSMSWebhook_Parallel(b *testing.B) {
	payload := []byte(`{
		"id": "msg-parallel",
		"type": "MESSAGE_STATUS",
		"messageId": "msg-sms-54321",
		"messageStatus": {
			"timestamp": "2025-11-26T14:30:00Z",
			"code": "SENT",
			"description": "Parallel test"
		}
	}`)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var webhook models.ZenviaWebhookRequest
			_ = json.Unmarshal(payload, &webhook)
		}
	})
}
