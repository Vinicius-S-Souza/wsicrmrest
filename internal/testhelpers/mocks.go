// Data de criação: 26/11/2025 14:30
// Versão: 3.0.0.6
// Mocks para testes e benchmarks - Simula database e logger sem conexões reais
package testhelpers

import (
	"time"

	"go.uber.org/zap"
)

// EmailData estrutura simplificada para testes (evita dependência circular)
type EmailData struct {
	EmailNumero int
	CliCodigo   int
	LogsApiId   int
}

// SMSData estrutura simplificada para testes (evita dependência circular)
type SMSData struct {
	SmsNumero int
	CliCodigo int
	LogsApiId int
}

// MockDatabase simula o database.Database para testes
// NÃO conecta ao Oracle real - apenas retorna dados fake
type MockDatabase struct {
	// Dados simulados de retorno
	EmailData *EmailData
	SMSData   *SMSData
	Error     error

	// Contadores para verificar chamadas
	GetEmailCalls     int
	GetSMSCalls       int
	InsereLogsCalls   int
	InsereOcorrCalls  int
}

// NewMockDatabase cria um mock database com dados padrão
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		EmailData: &EmailData{
			EmailNumero: 12345,
			CliCodigo:   67890,
			LogsApiId:   111,
		},
		SMSData: &SMSData{
			SmsNumero: 54321,
			CliCodigo: 98765,
			LogsApiId: 222,
		},
		Error: nil,
	}
}

// GetEmailByAPIMessageID simula busca de email por MessageID
func (m *MockDatabase) GetEmailByAPIMessageID(messageID string) (*EmailData, error) {
	m.GetEmailCalls++
	if m.Error != nil {
		return nil, m.Error
	}
	return m.EmailData, nil
}

// GetSMSByMessageID simula busca de SMS por MessageID
func (m *MockDatabase) GetSMSByMessageID(messageID string) (*SMSData, error) {
	m.GetSMSCalls++
	if m.Error != nil {
		return nil, m.Error
	}
	return m.SMSData, nil
}

// InsereLogsAPI simula inserção de logs de email
func (m *MockDatabase) InsereLogsAPI(logsApiId, emailNumero, statusEnvio int, statusCode, statusDescription, externalID string) error {
	m.InsereLogsCalls++
	return m.Error
}

// InsereLogsAPISMS simula inserção de logs de SMS
func (m *MockDatabase) InsereLogsAPISMS(logsApiId, smsNumero, statusEnvio int, statusCode, statusDescription, externalID string) error {
	m.InsereLogsCalls++
	return m.Error
}

// InsereOcorrenciaEmailInconsistente simula criação de ocorrência de email
func (m *MockDatabase) InsereOcorrenciaEmailInconsistente(cliCodigo int, email, statusCode, statusDescription string) error {
	m.InsereOcorrCalls++
	return m.Error
}

// InsereOcorrenciaSmsInconsistente simula criação de ocorrência de SMS
func (m *MockDatabase) InsereOcorrenciaSmsInconsistente(cliCodigo int, celular, statusCode, statusDescription string) error {
	m.InsereOcorrCalls++
	return m.Error
}

// GravaLogDB simula gravação de log de requisição (usado por todos os handlers)
func (m *MockDatabase) GravaLogDB(
	uuid, method, endpoint, headers, params string,
	statusCode int,
	response, procedure, clientID, appName string,
	startTime time.Time,
	clientIP string,
	gravaLogDB, detalheLogAPI bool,
	detalheLog, version string,
) {
	// Não faz nada - apenas simula
	// Na implementação real, grava no banco
}

// Exec simula execução de query SQL genérica
func (m *MockDatabase) Exec(query string, args ...interface{}) (interface{}, error) {
	return nil, m.Error
}

// Query simula execução de query SQL com retorno
func (m *MockDatabase) Query(query string, args ...interface{}) (interface{}, error) {
	return nil, m.Error
}

// QueryRow simula execução de query SQL com retorno de uma linha
func (m *MockDatabase) QueryRow(query string, args ...interface{}) interface{} {
	return nil
}

// MockLogger simula o zap.SugaredLogger para testes
// NÃO grava logs reais - apenas conta chamadas
type MockLogger struct {
	InfoCalls  int
	WarnCalls  int
	ErrorCalls int
	DebugCalls int
}

// NewMockLogger cria um logger mock
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Infow simula log de info
func (m *MockLogger) Infow(msg string, keysAndValues ...interface{}) {
	m.InfoCalls++
}

// Info simula log de info
func (m *MockLogger) Info(args ...interface{}) {
	m.InfoCalls++
}

// Warnw simula log de warning
func (m *MockLogger) Warnw(msg string, keysAndValues ...interface{}) {
	m.WarnCalls++
}

// Warn simula log de warning
func (m *MockLogger) Warn(args ...interface{}) {
	m.WarnCalls++
}

// Errorw simula log de erro
func (m *MockLogger) Errorw(msg string, keysAndValues ...interface{}) {
	m.ErrorCalls++
}

// Error simula log de erro
func (m *MockLogger) Error(args ...interface{}) {
	m.ErrorCalls++
}

// Debugw simula log de debug
func (m *MockLogger) Debugw(msg string, keysAndValues ...interface{}) {
	m.DebugCalls++
}

// Debug simula log de debug
func (m *MockLogger) Debug(args ...interface{}) {
	m.DebugCalls++
}

// With simula criação de logger com campos adicionais
func (m *MockLogger) With(args ...interface{}) *zap.SugaredLogger {
	// Retorna nil pois não é usado nos benchmarks
	return nil
}

// Sync simula sincronização de logs
func (m *MockLogger) Sync() error {
	return nil
}
