package database

import (
	"database/sql"
	"wsicrmrest/internal/config"
)

// LeOrganizador carrega os dados da organização do banco de dados
// Equivalente a pgLeOrganizador() do WinDev
func (d *Database) LeOrganizador(cfg *config.Config) error {
	query := `SELECT OrgCodigo, OrgNome, orgcnpj, orgcodlojamatriz, OrgFormaLimite,
	                 OrgCalcDispFuturoCartao, OrgCalcDispFuturoConvenio, ORGCODISGA,
	                 ORGDIAFATGRUPO1, ORGDIAFATGRUPO2, ORGDIAFATGRUPO3, ORGDIAFATGRUPO4, ORGDIAFATGRUPO5, ORGDIAFATGRUPO6,
	                 ORGDIACORGRUPO1, ORGDIACORGRUPO2, ORGDIACORGRUPO3, ORGDIACORGRUPO4, ORGDIACORGRUPO5, ORGDIACORGRUPO6
	          FROM ORGANIZADOR
	          WHERE OrgCodigo > 0`

	var (
		codigo                 int
		nome                   string
		cnpj                   string
		lojaMatriz             sql.NullInt64
		formaLimite            sql.NullInt64
		calcDispFuturoCartao   sql.NullInt64
		calcDispFuturoConvenio sql.NullInt64
		codISGA                sql.NullInt64
		diaVectoGrupo1         sql.NullInt64
		diaVectoGrupo2         sql.NullInt64
		diaVectoGrupo3         sql.NullInt64
		diaVectoGrupo4         sql.NullInt64
		diaVectoGrupo5         sql.NullInt64
		diaVectoGrupo6         sql.NullInt64
		diaCorteGrupo1         sql.NullInt64
		diaCorteGrupo2         sql.NullInt64
		diaCorteGrupo3         sql.NullInt64
		diaCorteGrupo4         sql.NullInt64
		diaCorteGrupo5         sql.NullInt64
		diaCorteGrupo6         sql.NullInt64
	)

	err := d.DB.QueryRow(query).Scan(
		&codigo,
		&nome,
		&cnpj,
		&lojaMatriz,
		&formaLimite,
		&calcDispFuturoCartao,
		&calcDispFuturoConvenio,
		&codISGA,
		&diaVectoGrupo1,
		&diaVectoGrupo2,
		&diaVectoGrupo3,
		&diaVectoGrupo4,
		&diaVectoGrupo5,
		&diaVectoGrupo6,
		&diaCorteGrupo1,
		&diaCorteGrupo2,
		&diaCorteGrupo3,
		&diaCorteGrupo4,
		&diaCorteGrupo5,
		&diaCorteGrupo6,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			d.Logger.Error("Organizador Não Cadastrado")
			return err
		}
		d.Logger.Errorw("Falha na Leitura do Organizador", "error", err)
		return err
	}

	// Atualizar configuração com dados do banco
	// Converter sql.NullInt64 para int (usando 0 se NULL)
	cfg.Organization.Codigo = codigo
	cfg.Organization.Nome = nome
	cfg.Organization.CNPJ = cnpj
	cfg.Organization.LojaMatriz = int(lojaMatriz.Int64)
	cfg.Organization.FormaLimite = int(formaLimite.Int64)
	cfg.Organization.CalcDispFuturoCartao = int(calcDispFuturoCartao.Int64)
	cfg.Organization.CalcDispFuturoConvenio = int(calcDispFuturoConvenio.Int64)
	cfg.Organization.CodISGA = int(codISGA.Int64)
	cfg.Organization.DiaVectoGrupo1 = int(diaVectoGrupo1.Int64)
	cfg.Organization.DiaVectoGrupo2 = int(diaVectoGrupo2.Int64)
	cfg.Organization.DiaVectoGrupo3 = int(diaVectoGrupo3.Int64)
	cfg.Organization.DiaVectoGrupo4 = int(diaVectoGrupo4.Int64)
	cfg.Organization.DiaVectoGrupo5 = int(diaVectoGrupo5.Int64)
	cfg.Organization.DiaVectoGrupo6 = int(diaVectoGrupo6.Int64)
	cfg.Organization.DiaCorteGrupo1 = int(diaCorteGrupo1.Int64)
	cfg.Organization.DiaCorteGrupo2 = int(diaCorteGrupo2.Int64)
	cfg.Organization.DiaCorteGrupo3 = int(diaCorteGrupo3.Int64)
	cfg.Organization.DiaCorteGrupo4 = int(diaCorteGrupo4.Int64)
	cfg.Organization.DiaCorteGrupo5 = int(diaCorteGrupo5.Int64)
	cfg.Organization.DiaCorteGrupo6 = int(diaCorteGrupo6.Int64)

	d.Logger.Infow("Organizador carregado com sucesso",
		"codigo", codigo,
		"nome", nome,
		"cnpj", cnpj)

	return nil
}
