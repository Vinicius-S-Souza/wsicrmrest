package database

import (
	"fmt"
	"strings"
	"time"
)

// EmailData representa os dados do email recuperados do banco
type EmailData struct {
	EmailNumero  int
	CliCodigo    int
	LogsApiId    int
	EmailAddress string
}

// SMSData representa os dados do SMS recuperados do banco
type SMSData struct {
	SMSNumero     int
	CliCodigo     int
	LogsApiId     int
	NumeroCelular string
}

// GetEmailByAPIMessageID busca dados do email pelo ID da mensagem da API
// Equivalente ao SELECT com SQLExec/SQLFetch do WinDev
func (d *Database) GetEmailByAPIMessageID(apiMessageID string) (*EmailData, error) {
	query := `SELECT e.emsgcodigo, e.clicodigo, l.logsapiid
	          FROM emailmensagem e
	          INNER JOIN logsapi l ON e.emsgcodigo = l.emsgcodigo
	          WHERE e.Emsgapimsgid = :1
	          AND l.logsapitipmensagem = 1`

	d.Logger.Infow("Executando query para buscar Email",
		"query", query,
		"emsgapimsgid_buscado", apiMessageID,
		"emsgapimsgid_length", len(apiMessageID),
		"emsgapimsgid_is_empty", apiMessageID == "")

	var emailData EmailData
	err := d.DB.QueryRow(query, apiMessageID).Scan(
		&emailData.EmailNumero,
		&emailData.CliCodigo,
		&emailData.LogsApiId,
	)

	if err != nil {
		d.Logger.Warnw("Query não retornou resultados",
			"error", err,
			"emsgapimsgid_buscado", apiMessageID)
		return nil, err
	}

	d.Logger.Infow("Email encontrado no banco",
		"emsgcodigo", emailData.EmailNumero,
		"clicodigo", emailData.CliCodigo,
		"logsapiid", emailData.LogsApiId)

	return &emailData, nil
}

// InsereLogsAPI insere registro na tabela logsApi e logsApiHistorico para Email
// Equivalente a pgInsereLogsAPI do WinDev
func (d *Database) InsereLogsAPI(
	emailNumero int,
	logsApiTipId int,
	logsApiStatus int,
	logsApiEnvio string,
	logsApiRetorno string,
	logsApiDtaCadastro string,
	logsApiHisDescricao string,
	tipoMensagem int,
	enderecoEmail string,
	logsApiId int,
	logsApiTag string,
) error {
	// Atualizar tabela logsApi usando bind variables (proteção SQL injection)
	queryUpdateLogsApi := `UPDATE logsapi
		SET logsapistatus = :1
		WHERE logsapiid = :2`

	if _, err := d.Exec(queryUpdateLogsApi, logsApiStatus, logsApiId); err != nil {
		d.Logger.Errorw("Falha ao atualizar logsApi",
			"error", err,
			"logsApiId", logsApiId)
		return err
	}

	d.Logger.Infow("LogsApi atualizado com sucesso",
		"logsApiId", logsApiId,
		"status", logsApiStatus)

	// Inserir em logsApiHistorico
	if err := d.InsertLogsApiHistorico(logsApiId, logsApiHisDescricao, logsApiStatus, logsApiTag); err != nil {
		return err
	}

	// Atualizar status da mensagem na tabela EmailMensagem
	if err := d.SetMsgStatus(logsApiStatus, logsApiHisDescricao, emailNumero); err != nil {
		return err
	}

	d.Logger.Infow("Logs API atualizados com sucesso",
		"emailNumero", emailNumero,
		"logsApiId", logsApiId,
		"status", logsApiStatus,
		"tag", logsApiTag)

	return nil
}

// SetMsgStatus atualiza o status da mensagem na tabela EmailMensagem
// Equivalente a pgSetMsgStatus do WinDev
func (d *Database) SetMsgStatus(logsApiStatus int, logsApiHisDescricao string, emailNumero int) error {
	query := `UPDATE EmailMensagem
		SET EMsgStsEnvio = :1
		WHERE EMsgCodigo = :2`

	if _, err := d.Exec(query, logsApiStatus, emailNumero); err != nil {
		d.Logger.Errorw("Falha ao atualizar status Email",
			"error", err,
			"emailNumero", emailNumero)
		return err
	}

	d.Logger.Infow("Status Email atualizado com sucesso",
		"emailNumero", emailNumero,
		"status", logsApiStatus,
		"descricao", logsApiHisDescricao)

	return nil
}

// InsereOcorrenciaEmailInconsistente insere ocorrência de email inconsistente
// Equivalente a pgInsereOcorrenciaEmailInconsistente do WinDev
func (d *Database) InsereOcorrenciaEmailInconsistente(
	enderecoEmail string,
	cliCodigo int,
	ocorrenciaTipo int,
	provedor string,
) error {
	// Buscar dados do cliente e emails de ClientesExtensao
	queryCliente := `SELECT c.clinome, c.clicpfcnpj, ce.CliExtEmail2, ce.CliExtEmail3
		FROM clientes c
		INNER JOIN Clientesextensao ce ON c.clicodigo = ce.clicodigo
		WHERE c.clicodigo = :1`

	var cliNome, cliCpfCnpj, cliExtEmail2, cliExtEmail3 string
	err := d.DB.QueryRow(queryCliente, cliCodigo).Scan(&cliNome, &cliCpfCnpj, &cliExtEmail2, &cliExtEmail3)
	if err != nil {
		d.Logger.Errorw("Erro ao buscar dados do cliente", "error", err, "cliCodigo", cliCodigo)
		return err
	}

	d.Logger.Infow("Dados do cliente recuperados",
		"cliCodigo", cliCodigo,
		"cliNome", cliNome)

	// Determinar qual campo deve ser limpo
	var campoLimpar string
	if cliExtEmail2 == enderecoEmail {
		campoLimpar = "CliExtEmail2"
	} else if cliExtEmail3 == enderecoEmail {
		campoLimpar = "CliExtEmail3"
	} else {
		d.Logger.Warnw("Email do cliente não corresponde ao email inconsistente",
			"cliExtEmail2", cliExtEmail2,
			"cliExtEmail3", cliExtEmail3,
			"emailInconsistente", enderecoEmail)
		return nil
	}

	// Obter próximo ID usando MAX+1
	var maxOcoCod int
	queryMax := `SELECT NVL(MAX(OcoCod), 0) + 1 FROM Ocorrencia`
	err = d.DB.QueryRow(queryMax).Scan(&maxOcoCod)
	if err != nil {
		d.Logger.Errorw("Erro ao obter MAX(OcoCod)", "error", err)
		return err
	}

	d.Logger.Infow("Próximo OcoCod calculado",
		"ococod", maxOcoCod)

	// Inserir ocorrência usando bind variables (proteção SQL injection)
	now := time.Now()
	queryInsert := `INSERT INTO Ocorrencia(
		OcoCod,
		EntCod,
		CliCod,
		CodPltCod,
		TocNum,
		OcoTip,
		OcoCliNon,
		OcoDsc,
		OcoUsrSol,
		OcoSolDta,
		UsrAlt,
		DatCad
	) VALUES(:1, :2, :3, :4, :5, 2, :6, :7, :8, :9, :10, :11)`

	if _, err := d.Exec(queryInsert,
		maxOcoCod,
		d.Config.Organization.Codigo,
		cliCodigo,
		cliCpfCnpj,
		ocorrenciaTipo,
		cliNome,
		"Email inválido. Não foi possível o envio de mensagem para esse email, favor preencher o email corretamente.",
		"WebHookSendGrid",
		now,
		"WebHookSendGrid",
		now,
	); err != nil {
		d.Logger.Errorw("Falha ao inserir ocorrência de email inconsistente",
			"error", err,
			"cliCodigo", cliCodigo,
			"email", enderecoEmail)
		return err
	}

	d.Logger.Infow("Ocorrência de email inconsistente inserida com sucesso",
		"cliCodigo", cliCodigo,
		"email", enderecoEmail,
		"tipo", ocorrenciaTipo,
		"provedor", provedor,
		"campo", campoLimpar)

	// Limpar email inconsistente
	if err := d.LimpaEmailInconsistente(campoLimpar, cliCodigo); err != nil {
		d.Logger.Errorw("Erro ao limpar email inconsistente", "error", err)
		// Não retorna erro pois a ocorrência já foi inserida
	}

	return nil
}

// LimpaEmailInconsistente limpa o campo de email inconsistente na tabela ClientesExtensao
// Equivalente a pgLimpaEmailInconsistente do WinDev
func (d *Database) LimpaEmailInconsistente(campoDeveLimpar string, cliCodigo int) error {
	// Validar nome da coluna para prevenir SQL injection
	// Apenas CliExtEmail2 e CliExtEmail3 são permitidos
	if campoDeveLimpar != "CliExtEmail2" && campoDeveLimpar != "CliExtEmail3" {
		return fmt.Errorf("campo inválido: %s", campoDeveLimpar)
	}

	// Usar nome de coluna validado diretamente (não é parametrizável em Oracle)
	query := fmt.Sprintf(`UPDATE Clientesextensao
		SET %s = NULL
		WHERE clicodigo = :1`, campoDeveLimpar)

	if _, err := d.Exec(query, cliCodigo); err != nil {
		d.Logger.Errorw("Falha ao limpar email inconsistente",
			"error", err,
			"cliCodigo", cliCodigo,
			"campo", campoDeveLimpar)
		return err
	}

	d.Logger.Infow("Email inconsistente limpo com sucesso",
		"cliCodigo", cliCodigo,
		"campo", campoDeveLimpar)

	return nil
}

// GetSMSByMessageID busca dados do SMS pelo ID da mensagem (smscodigo)
// Equivalente ao SELECT com SQLExec/SQLFetch do WinDev para SMS
func (d *Database) GetSMSByMessageID(smscodigo string) (*SMSData, error) {
	query := `SELECT s.smscodigo, s.clicodigo, l.logsapiid
	          FROM smsmensagem s
	          INNER JOIN logsapi l ON s.smscodigo = l.emsgcodigo
	          WHERE s.SMSAPIID = :1
	          AND l.logsapitipmensagem = 2`

	d.Logger.Infow("Executando query para buscar SMS",
		"query", query,
		"smsapiid_buscado", smscodigo,
		"smsapiid_length", len(smscodigo),
		"smsapiid_is_empty", smscodigo == "")

	var smsData SMSData
	err := d.DB.QueryRow(query, smscodigo).Scan(
		&smsData.SMSNumero,
		&smsData.CliCodigo,
		&smsData.LogsApiId,
	)

	if err != nil {
		d.Logger.Warnw("Query não retornou resultados",
			"error", err,
			"smsapiid_buscado", smscodigo)
		return nil, err
	}

	d.Logger.Infow("SMS encontrado no banco",
		"smscodigo", smsData.SMSNumero,
		"clicodigo", smsData.CliCodigo,
		"logsapiid", smsData.LogsApiId)

	return &smsData, nil
}

// InsereLogsAPISMS insere/atualiza registro na tabela logsApi e logsApiHistorico para SMS
// Equivalente a pgInsereLogsAPISms do WinDev
func (d *Database) InsereLogsAPISMS(
	smscodigo int,
	logsApiTipId int,
	logsApiStatus int,
	logsApiEnvio string,
	logsApiRetorno string,
	logsApiDtaCadastro string,
	logsApiHisDescricao string,
	tipoMensagem int,
	numeroCelular string,
	logsApiId int,
	logsApiTag string,
) error {
	// Atualizar tabela logsApi usando bind variables (proteção SQL injection)
	queryUpdateLogsApi := `UPDATE logsapi
		SET logsapistatus = :1,
		    logsapiretorno = :2
		WHERE logsapiid = :3`

	if _, err := d.Exec(queryUpdateLogsApi,
		logsApiStatus,
		logsApiRetorno,
		logsApiId,
	); err != nil {
		d.Logger.Errorw("Falha ao atualizar logsApi",
			"error", err,
			"logsApiId", logsApiId)
		return err
	}

	d.Logger.Infow("LogsApi atualizado com sucesso",
		"logsApiId", logsApiId,
		"status", logsApiStatus)

	// Inserir em logsApiHistorico
	if err := d.InsertLogsApiHistorico(logsApiId, logsApiHisDescricao, logsApiStatus, logsApiTag); err != nil {
		return err
	}

	// Atualizar status da mensagem na tabela smsmensagem
	if err := d.SetMsgStatusSMS(logsApiStatus, logsApiHisDescricao, smscodigo); err != nil {
		return err
	}

	d.Logger.Infow("Logs API SMS atualizados com sucesso",
		"smscodigo", smscodigo,
		"logsApiId", logsApiId,
		"status", logsApiStatus,
		"tag", logsApiTag)

	return nil
}

// InsertLogsApiHistorico insere registro no histórico da API
// Equivalente a pgInsertLogsApiHistorico do WinDev
func (d *Database) InsertLogsApiHistorico(logsApiId int, logsApiHisDescricao string, logsApiStatus int, plataforma string) error {
	// Obter próximo sequencial
	querySeq := `SELECT NVL(MAX(LogsApiHisSequencial), 0) + 1 FROM LogsApiHistorico WHERE LogsApiId = :1`

	var sequencial int
	err := d.DB.QueryRow(querySeq, logsApiId).Scan(&sequencial)
	if err != nil {
		d.Logger.Errorw("Erro ao obter sequencial", "error", err, "logsApiId", logsApiId)
		return err
	}

	// Inserir histórico usando bind variables (proteção SQL injection)
	queryInsert := `INSERT INTO logsApiHistorico(
		logsApiId,
		logsApiHisSequencial,
		logsApiHisTag,
		logsApiHisDescricao,
		logsapihisstatus,
		logsapihisdata
	) VALUES(:1, :2, :3, :4, :5, :6)`

	if _, err := d.Exec(queryInsert,
		logsApiId,
		sequencial,
		plataforma,
		logsApiHisDescricao,
		logsApiStatus,
		time.Now(),
	); err != nil {
		d.Logger.Errorw("Falha ao inserir em logsApiHistorico",
			"error", err,
			"logsApiId", logsApiId)
		return err
	}

	d.Logger.Infow("Histórico inserido com sucesso",
		"logsApiId", logsApiId,
		"sequencial", sequencial,
		"plataforma", plataforma)

	return nil
}

// SetMsgStatusSMS atualiza o status da mensagem na tabela smsmensagem
// Equivalente a pgSetMsgStatusSMS do WinDev
func (d *Database) SetMsgStatusSMS(logsApiStatus int, logsApiHisDescricao string, smscodigo int) error {
	query := `UPDATE smsmensagem
		SET smsstsenvio = :1
		WHERE smscodigo = :2`

	if _, err := d.Exec(query, logsApiStatus, smscodigo); err != nil {
		d.Logger.Errorw("Falha ao atualizar status SMS",
			"error", err,
			"smscodigo", smscodigo)
		return err
	}

	d.Logger.Infow("Status SMS atualizado com sucesso",
		"smscodigo", smscodigo,
		"status", logsApiStatus,
		"descricao", logsApiHisDescricao)

	return nil
}

// InsereOcorrenciaSmsInconsistente insere ocorrência de SMS inconsistente
// Equivalente a pgInsereOcorrenciaSmsInconsistente do WinDev
func (d *Database) InsereOcorrenciaSmsInconsistente(
	numeroCelular string,
	cliCodigo int,
	ocorrenciaTipo int,
	provedor string,
) error {
	// Buscar dados do cliente
	queryCliente := `SELECT c.clinome, c.clicpfcnpj, c.clicelular, c.clidddcelular
		FROM clientes c
		WHERE c.clicodigo = :1`

	var cliNome, cliCpfCnpj, cliCelular, clidddcelular string
	err := d.DB.QueryRow(queryCliente, cliCodigo).Scan(&cliNome, &cliCpfCnpj, &cliCelular, &clidddcelular)
	if err != nil {
		d.Logger.Errorw("Erro ao buscar dados do cliente", "error", err, "cliCodigo", cliCodigo)
		return err
	}

	d.Logger.Infow("Dados do cliente recuperados",
		"cliCodigo", cliCodigo,
		"cliNome", cliNome)

	// Formata Número do Celular do Cliente com DDD
	auxCliCelular := clidddcelular + cliCelular

	// Remover prefixo +55 ou 55 para comparação
	auxNumeroCelular := numeroCelular

	if strings.HasPrefix(auxNumeroCelular, "+55") {
		auxNumeroCelular = strings.TrimPrefix(auxNumeroCelular, "+55")
	} else if strings.HasPrefix(auxNumeroCelular, "55") {
		auxNumeroCelular = strings.TrimPrefix(auxNumeroCelular, "55")
	}

	// Verificar se o celular do cliente corresponde ao número inconsistente
	// if cliCelular != numeroCelular {
	if auxCliCelular != auxNumeroCelular {
		d.Logger.Warnw("Celular do cliente não corresponde ao número inconsistente",
			"cliCelular", cliCelular,
			"numeroCelular", numeroCelular)
		return nil
	}

	// Obter próximo ID usando MAX+1
	var maxOcoCod int
	queryMax := `SELECT NVL(MAX(OcoCod), 0) + 1 FROM Ocorrencia`
	err = d.DB.QueryRow(queryMax).Scan(&maxOcoCod)
	if err != nil {
		d.Logger.Errorw("Erro ao obter MAX(OcoCod)", "error", err)
		return err
	}

	d.Logger.Infow("Próximo OcoCod calculado",
		"ococod", maxOcoCod)

	// Inserir ocorrência usando bind variables (proteção SQL injection)
	now := time.Now()
	queryInsert := `INSERT INTO Ocorrencia(
		OcoCod,
		EntCod,
		CliCod,
		CodPltCod,
		TocNum,
		OcoTip,
		OcoCliNon,
		OcoDsc,
		OcoUsrSol,
		OcoSolDta,
		UsrAlt,
		DatCad
	) VALUES(:1, :2, :3, :4, :5, 2, :6, :7, :8, :9, :10, :11)`

	if _, err := d.Exec(queryInsert,
		maxOcoCod,
		d.Config.Organization.Codigo,
		cliCodigo,
		cliCpfCnpj,
		ocorrenciaTipo,
		cliNome,
		"Celular inválido. Não foi possível o envio de mensagem para esse celular, favor preencher o celular corretamente.",
		provedor,
		now,
		provedor,
		now,
	); err != nil {
		d.Logger.Errorw("Falha ao inserir ocorrência de SMS inconsistente",
			"error", err,
			"cliCodigo", cliCodigo,
			"celular", numeroCelular)
		return err
	}

	d.Logger.Infow("Ocorrência de SMS inconsistente inserida com sucesso",
		"cliCodigo", cliCodigo,
		"celular", numeroCelular,
		"tipo", ocorrenciaTipo,
		"provedor", provedor)

	// Limpar celular inconsistente
	if err := d.LimpaCelularInconsistente(cliCodigo); err != nil {
		d.Logger.Errorw("Erro ao limpar celular inconsistente", "error", err)
		// Não retorna erro pois a ocorrência já foi inserida
	}

	return nil
}

// LimpaCelularInconsistente limpa o campo de celular inconsistente
// Equivalente a pgLimpaCelularInconsistente do WinDev
func (d *Database) LimpaCelularInconsistente(cliCodigo int) error {
	query := `UPDATE clientes
		SET clicelular = NULL, clidddcelular = NULL
		WHERE clicodigo = :1`

	if _, err := d.Exec(query, cliCodigo); err != nil {
		d.Logger.Errorw("Falha ao limpar celular inconsistente",
			"error", err,
			"cliCodigo", cliCodigo)
		return err
	}

	d.Logger.Infow("Celular inconsistente limpo com sucesso",
		"cliCodigo", cliCodigo)

	return nil
}
