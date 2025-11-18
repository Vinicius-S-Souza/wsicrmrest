package database

import (
	"database/sql"
	"fmt"
	"wsicrmrest/internal/config"

	_ "github.com/godror/godror"
	"go.uber.org/zap"
)

// Database encapsula a conexão com o banco de dados
type Database struct {
	DB     *sql.DB
	Config *config.Config
	Logger *zap.SugaredLogger
}

// NewDatabase cria uma nova conexão com o banco de dados Oracle usando TNSNAMES
func NewDatabase(cfg *config.Config, logger *zap.SugaredLogger) (*Database, error) {
	// String de conexão usando TNS
	// Formato: user/password@tnsname
	connStr := fmt.Sprintf("%s/%s@%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.TNSName,
	)

	db, err := sql.Open("godror", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	// Testar conexão
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &Database{
		DB:     db,
		Config: cfg,
		Logger: logger,
	}, nil
}

// Close fecha a conexão com o banco de dados
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Query executa uma query e retorna os resultados
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.Logger.Debugw("Executando query", "query", query, "args", args)
	return d.DB.Query(query, args...)
}

// Exec executa um comando SQL
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.Logger.Debugw("Executando comando", "query", query, "args", args)
	return d.DB.Exec(query, args...)
}

// QueryRow executa uma query que retorna uma única linha
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	d.Logger.Debugw("Executando query row", "query", query, "args", args)
	return d.DB.QueryRow(query, args...)
}
