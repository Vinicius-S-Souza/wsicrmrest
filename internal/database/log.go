package database

import (
	"strings"
	"time"
)

// GravaLogDB grava o log da requisição no banco de dados
// Equivalente a pgGravaLogDB do WinDev
func (d *Database) GravaLogDB(
	reqUUID string,
	reqMetodo string,
	reqEndPoint string,
	reqHeader string,
	reqParametro string,
	reqCodRetorno int,
	reqResposta string,
	reqProcedure string,
	clientID string,
	nomeAplicacao string,
	startTime time.Time,
	ipCliente string,
	wsGravaLogDB bool,
	wsDetalheLogAPI bool,
	detalheLogAPI string,
	versao string,
) bool {

	if !wsGravaLogDB {
		return true
	}

	// Testar conexão
	if err := d.DB.Ping(); err != nil {
		d.Logger.Errorw("Falha na Abertura do Banco de Dados na Gravação do Log no Banco", "error", err)
		return false
	}

	// Remover header Authorization do log por segurança
	reqHeader = removeAuthorizationHeader(reqHeader)

	// Data/hora da resposta
	dtDataResposta := time.Now()
	duracao := dtDataResposta.Sub(startTime).Milliseconds()

	// Construir query com bind variables (proteção contra SQL injection)
	var query string
	var args []interface{}

	if wsDetalheLogAPI {
		query = `INSERT INTO WSREQUISICOES(WSREQUUID, WSREQDTARECEBE, WSREQIPCLIENTE, WSREQENDPOINT, WSREQMETODO,
				WSREQHEADER, WSREQPARAMETROS, WSREQDTARESPOSTA, WSREQCODRESPOSTA, WSREQRESPOSTA, WSREQDURACAO, WSREQPROCEDURE,
				WSAPLNOME, WSAPLCLIENTID, WSVERSAO, WSLOGDETALHE)
				VALUES(:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16)`

		// Preparar valor NULL ou string para detalhe
		var detalheSQL interface{}
		if detalheLogAPI != "" {
			detalheSQL = detalheLogAPI
		} else {
			detalheSQL = nil
		}

		args = []interface{}{
			reqUUID,
			startTime,
			ipCliente,
			strings.ToUpper(reqEndPoint),
			strings.ToUpper(reqMetodo),
			reqHeader,
			reqParametro,
			dtDataResposta,
			reqCodRetorno,
			reqResposta,
			duracao,
			reqProcedure,
			nomeAplicacaoOrNil(nomeAplicacao),
			clientIDOrNil(clientID),
			strings.ToUpper(versao),
			detalheSQL,
		}
	} else {
		query = `INSERT INTO WSREQUISICOES(WSREQUUID, WSREQDTARECEBE, WSREQIPCLIENTE, WSREQENDPOINT, WSREQMETODO,
				WSREQHEADER, WSREQPARAMETROS, WSREQDTARESPOSTA, WSREQCODRESPOSTA, WSREQRESPOSTA, WSREQDURACAO, WSREQPROCEDURE,
				WSAPLNOME, WSAPLCLIENTID, WSVERSAO)
				VALUES(:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15)`

		args = []interface{}{
			reqUUID,
			startTime,
			ipCliente,
			strings.ToUpper(reqEndPoint),
			strings.ToUpper(reqMetodo),
			reqHeader,
			reqParametro,
			dtDataResposta,
			reqCodRetorno,
			reqResposta,
			duracao,
			reqProcedure,
			nomeAplicacaoOrNil(nomeAplicacao),
			clientIDOrNil(clientID),
			strings.ToUpper(versao),
		}
	}

	// Executar query com bind variables
	if _, err := d.Exec(query, args...); err != nil {
		d.Logger.Errorw("Falha na Inclusão do Log de Requisição",
			"error", err,
			"uuid", reqUUID)
		return false
	}

	return true
}

// removeAuthorizationHeader remove o header Authorization por segurança
func removeAuthorizationHeader(header string) string {
	lines := strings.Split(header, "\n")
	var filtered []string

	for _, line := range lines {
		upperLine := strings.ToUpper(line)
		if !strings.HasPrefix(upperLine, "AUTHORIZATION:") &&
			!strings.HasPrefix(upperLine, "X-API-KEY:") &&
			!strings.HasPrefix(upperLine, "X-AUTH-TOKEN:") &&
			!strings.HasPrefix(upperLine, "COOKIE:") {
			filtered = append(filtered, line)
		}
	}

	return strings.Join(filtered, "\n")
}

// clientIDOrNil retorna o clientID ou nil se vazio
func clientIDOrNil(clientID string) interface{} {
	if clientID == "" {
		return nil
	}
	return clientID
}

// nomeAplicacaoOrNil retorna o nome da aplicação em maiúsculas ou nil se vazio
func nomeAplicacaoOrNil(nome string) interface{} {
	if nome == "" {
		return nil
	}
	return strings.ToUpper(nome)
}
