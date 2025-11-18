package context

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// RequestContext armazena o contexto de uma requisição
type RequestContext struct {
	UUID            string
	ClientID        string
	NomeAplicacao   string
	StartTime       time.Time
	DetalheLogAPI   string
	WSGravaLogDB    bool
	WSDetalheLogAPI bool
	mu              sync.Mutex
}

// NewRequestContext cria um novo contexto de requisição
func NewRequestContext() *RequestContext {
	return &RequestContext{
		UUID:            uuid.New().String(),
		StartTime:       time.Now(),
		WSGravaLogDB:    true,  // Habilita gravação de log no BD por padrão
		WSDetalheLogAPI: false, // Desabilita detalhes por padrão
		DetalheLogAPI:   "",
	}
}

// SetClientInfo define informações do cliente
func (rc *RequestContext) SetClientInfo(clientID, nomeAplicacao string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.ClientID = clientID
	rc.NomeAplicacao = nomeAplicacao
}

// AddLogDetail adiciona detalhes ao log
func (rc *RequestContext) AddLogDetail(detail string) {
	if !rc.WSDetalheLogAPI {
		return
	}
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.DetalheLogAPI != "" {
		rc.DetalheLogAPI += "\n"
	}
	rc.DetalheLogAPI += time.Now().Format("02/01/2006 15:04:05.00") + " -> " + detail
}

// GetDuration retorna a duração desde o início da requisição
func (rc *RequestContext) GetDuration() time.Duration {
	return time.Since(rc.StartTime)
}
