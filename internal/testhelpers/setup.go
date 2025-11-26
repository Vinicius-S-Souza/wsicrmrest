// Data de criação: 26/11/2025 14:38
// Versão: 3.0.0.6
// Funções helper para setup de testes e benchmarks
package testhelpers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"wsicrmrest/internal/config"

	"github.com/gin-gonic/gin"
)

// SetupTestConfig cria uma configuração de teste padrão
func SetupTestConfig() *config.Config {
	return &config.Config{
		Application: config.ApplicationConfig{
			Environment:     "test",
			Port:            "8080",
			LogDir:          "log",
			RequestTimeout:  30,
			WSGravaLogDB:    false, // Desabilitar log de DB em testes
			WSDetalheLogAPI: false,
		},
		Database: config.DatabaseConfig{
			TNSName:  "TEST_DB",
			Username: "test_user",
			Password: "test_pass",
		},
		Security: config.SecurityConfig{
			MaxBodySize:      1048576,
			RateLimitPerMin:  60,
			RateLimitPerHour: 1000,
			RateLimitEnabled: false, // Desabilitar rate limit em testes
		},
		CORS: config.CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   "GET,POST,PUT,DELETE,OPTIONS",
			AllowedHeaders:   "Origin,Content-Type,Authorization",
			AllowCredentials: true,
			MaxAge:           "43200",
		},
		Organization: config.OrganizationConfig{
			Codigo: 1,
			Nome:   "Organização Teste",
		},
		JWT: config.JWTConfig{
			SecretKey:   "CloudI0812IcrMmDB",
			Issuer:      "WSCloudICrmIntellsys",
			KeyDelivery: "Ped2505IcrM",
			Timezone:    0,
		},
		TLS: config.TLSConfig{
			Enabled: false,
		},
		Fail2Ban: config.Fail2BanConfig{
			Enabled: false, // Desabilitar Fail2Ban em testes
		},
	}
}

// SetupGinTestMode configura o Gin para modo de teste
func SetupGinTestMode() {
	gin.SetMode(gin.TestMode)
}

// CreateTestRequest cria uma requisição HTTP de teste com payload JSON
func CreateTestRequest(method, url, payload string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateTestRecorder cria um ResponseRecorder para capturar resposta HTTP
func CreateTestRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// SetupTestRouter cria um router Gin básico para testes
func SetupTestRouter() *gin.Engine {
	SetupGinTestMode()
	router := gin.New()
	router.Use(gin.Recovery())
	return router
}

// CreateWebhookEmailRequest cria requisição para webhook de email
func CreateWebhookEmailRequest(payload string) *http.Request {
	return CreateTestRequest("POST", "/webhook/zenvia/email", payload)
}

// CreateWebhookSMSRequest cria requisição para webhook de SMS
func CreateWebhookSMSRequest(payload string) *http.Request {
	return CreateTestRequest("POST", "/webhook/zenvia/sms", payload)
}

// SetupWebhookTest prepara ambiente completo para teste de webhook
// Retorna: router, config, mockDB, mockLogger, recorder
func SetupWebhookTest() (*gin.Engine, *config.Config, *MockDatabase, *MockLogger, *httptest.ResponseRecorder) {
	router := SetupTestRouter()
	cfg := SetupTestConfig()
	mockDB := NewMockDatabase()
	mockLogger := NewMockLogger()
	recorder := CreateTestRecorder()

	return router, cfg, mockDB, mockLogger, recorder
}

// ResetMockCounters reseta contadores de chamadas do mock database e logger
func ResetMockCounters(mockDB *MockDatabase, mockLogger *MockLogger) {
	// Reset database counters
	mockDB.GetEmailCalls = 0
	mockDB.GetSMSCalls = 0
	mockDB.InsereLogsCalls = 0
	mockDB.InsereOcorrCalls = 0

	// Reset logger counters
	mockLogger.InfoCalls = 0
	mockLogger.WarnCalls = 0
	mockLogger.ErrorCalls = 0
	mockLogger.DebugCalls = 0
}

// BenchmarkOptions opções para configuração de benchmarks
type BenchmarkOptions struct {
	PayloadType string // "email" ou "sms"
	StatusCode  int    // 121-125
	WithError   bool   // Simular erro de database
	InvalidJSON bool   // Usar JSON inválido
}

// GetPayloadForBenchmark retorna payload apropriado baseado nas opções
func GetPayloadForBenchmark(opts BenchmarkOptions) string {
	if opts.InvalidJSON {
		return InvalidJSON
	}

	if opts.PayloadType == "email" {
		return GetEmailPayloadByStatus(opts.StatusCode)
	}

	return GetSMSPayloadByStatus(opts.StatusCode)
}

// SetupBenchmark prepara ambiente para benchmark com opções customizadas
func SetupBenchmark(opts BenchmarkOptions) (*gin.Engine, *config.Config, *MockDatabase, *MockLogger) {
	SetupGinTestMode()

	router := gin.New()
	router.Use(gin.Recovery())

	cfg := SetupTestConfig()
	mockDB := NewMockDatabase()
	mockLogger := NewMockLogger()

	// Configurar erro no mock database se solicitado
	if opts.WithError {
		mockDB.Error = &MockError{Message: "Simulated database error"}
	}

	return router, cfg, mockDB, mockLogger
}

// MockError implementa interface error para simular erros
type MockError struct {
	Message string
}

func (e *MockError) Error() string {
	return e.Message
}

// AssertStatusCode verifica se o status code da resposta é o esperado
func AssertStatusCode(recorder *httptest.ResponseRecorder, expected int) bool {
	return recorder.Code == expected
}

// GetResponseBody retorna o body da resposta como string
func GetResponseBody(recorder *httptest.ResponseRecorder) string {
	return recorder.Body.String()
}

// CreateBenchmarkEmailPayloads cria conjunto de payloads para benchmark de email
func CreateBenchmarkEmailPayloads() []struct {
	Name    string
	Payload string
} {
	return []struct {
		Name    string
		Payload string
	}{
		{"EmailSent", ZenviaEmailPayloadSent},
		{"EmailDelivered", ZenviaEmailPayloadDelivered},
		{"EmailRead", ZenviaEmailPayloadRead},
		{"EmailRejected", ZenviaEmailPayloadRejected},
		{"EmailBounce", ZenviaEmailPayloadBounce},
	}
}

// CreateBenchmarkSMSPayloads cria conjunto de payloads para benchmark de SMS
func CreateBenchmarkSMSPayloads() []struct {
	Name    string
	Payload string
} {
	return []struct {
		Name    string
		Payload string
	}{
		{"SMSSent", ZenviaSMSPayloadSent},
		{"SMSDelivered", ZenviaSMSPayloadDelivered},
		{"SMSRead", ZenviaSMSPayloadRead},
		{"SMSRejected", ZenviaSMSPayloadRejected},
		{"SMSBounce", ZenviaSMSPayloadBounce},
	}
}
