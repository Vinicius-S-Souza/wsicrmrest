package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// Variáveis globais da aplicação
// Estas variáveis podem ser injetadas durante a compilação usando -ldflags
var (
	// Version - Versão do sistema (atualizado manualmente)
	Version = "Versão 3.0.0.2 (GO)"
	// VersionDate - Data da versão (injetado automaticamente durante a compilação)
	VersionDate = "2025-11-24"
	// BuildTime - Data e hora da compilação (injetado automaticamente)
	BuildTime = "unknown"
)

// Config representa as configurações da aplicação
type Config struct {
	Database     DatabaseConfig
	JWT          JWTConfig
	Organization OrganizationConfig
	Application  ApplicationConfig
	CORS         CORSConfig
	TLS          TLSConfig
	Security     SecurityConfig
}

// DatabaseConfig representa as configurações do banco de dados Oracle
type DatabaseConfig struct {
	TNSName  string
	Username string
	Password string
}

// JWTConfig representa as configurações de JWT
type JWTConfig struct {
	SecretKey   string // gsKey - Chave HMAC para JWT
	Issuer      string // gsIss - Issuer do JWT
	KeyDelivery string // gsKeyDelivery - Chave adicional para delivery
	Timezone    int    // gnFusoHorario - Fuso horário em horas (ex: -3 para Brasília, 0 para UTC)
}

// OrganizationConfig representa os dados da organização
type OrganizationConfig struct {
	Codigo                 int
	Nome                   string
	CNPJ                   string
	LojaMatriz             int
	CodISGA                int
	RegModulos             int // Valor padrão ou carregado do banco
	FormaLimite            int
	CalcDispFuturoCartao   int
	CalcDispFuturoConvenio int
	DiaVectoGrupo1         int
	DiaVectoGrupo2         int
	DiaVectoGrupo3         int
	DiaVectoGrupo4         int
	DiaVectoGrupo5         int
	DiaVectoGrupo6         int
	DiaCorteGrupo1         int
	DiaCorteGrupo2         int
	DiaCorteGrupo3         int
	DiaCorteGrupo4         int
	DiaCorteGrupo5         int
	DiaCorteGrupo6         int
}

// ApplicationConfig representa as configurações da aplicação
type ApplicationConfig struct {
	Environment     string
	Port            string
	LogDir          string // Diretório de logs
	WSGravaLogDB    bool   // Se deve gravar log no banco
	WSDetalheLogAPI bool   // Se deve gravar detalhes do log
	RequestTimeout  int    // Timeout de requisição em segundos (padrão: 30)
}

// CORSConfig representa as configurações de CORS
type CORSConfig struct {
	AllowedOrigins   []string // Lista de origens permitidas (vazio = permite todas)
	AllowedMethods   string   // Métodos HTTP permitidos
	AllowedHeaders   string   // Headers permitidos
	AllowCredentials bool     // Permite credenciais (cookies, auth headers)
	MaxAge           string   // Tempo de cache do preflight em segundos
}

// TLSConfig representa as configurações de HTTPS/TLS
type TLSConfig struct {
	Enabled  bool   // Habilitar HTTPS
	CertFile string // Caminho do certificado TLS
	KeyFile  string // Caminho da chave privada TLS
	Port     string // Porta HTTPS (padrão: 8443)
}

// SecurityConfig representa as configurações de segurança
type SecurityConfig struct {
	MaxBodySize      int64 // Tamanho máximo do body em bytes (padrão: 1MB)
	RateLimitPerMin  int   // Limite de requests por minuto (0 = desabilitado)
	RateLimitPerHour int   // Limite de requests por hora (0 = desabilitado)
	RateLimitEnabled bool  // Habilitar rate limiting
}

// LoadConfig carrega as configurações do arquivo dbinit.ini
func LoadConfig(filename string) (*Config, error) {
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo %s: %w", filename, err)
	}

	config := &Config{
		Database: DatabaseConfig{
			TNSName:  cfg.Section("database").Key("tns_name").String(),
			Username: cfg.Section("database").Key("username").String(),
			Password: cfg.Section("database").Key("password").String(),
		},
		JWT: JWTConfig{
			// Valores fixos das variáveis globais WinDev (não configuráveis)
			SecretKey:   "CloudI0812IcrMmDB",    // gsKey
			Issuer:      "WSCloudICrmIntellsys", // gsIss
			KeyDelivery: "Ped2505IcrM",          // gsKeyDelivery
			Timezone:    0,                      // gnFusoHorario (0 = UTC)
		},
		Organization: OrganizationConfig{
			// Todos os dados serão carregados da tabela ORGANIZADOR via pgLeOrganizador()
			RegModulos: 1, // Valor inicial conforme especificado
		},
		Application: ApplicationConfig{
			Environment:     cfg.Section("application").Key("environment").MustString("development"),
			Port:            cfg.Section("application").Key("port").MustString("8080"),
			LogDir:          cfg.Section("application").Key("log_dir").MustString("log"),
			WSGravaLogDB:    cfg.Section("application").Key("ws_grava_log_db").MustBool(true),
			WSDetalheLogAPI: cfg.Section("application").Key("ws_detalhe_log_api").MustBool(false),
			RequestTimeout:  cfg.Section("application").Key("request_timeout").MustInt(30),
		},
		CORS:     loadCORSConfig(cfg),
		TLS:      loadTLSConfig(cfg),
		Security: loadSecurityConfig(cfg),
	}

	// Validações básicas
	if config.Database.TNSName == "" {
		return nil, fmt.Errorf("tns_name não configurado em [database]")
	}
	if config.Database.Username == "" {
		return nil, fmt.Errorf("username não configurado em [database]")
	}
	if config.Database.Password == "" {
		return nil, fmt.Errorf("password não configurado em [database]")
	}
	if config.JWT.SecretKey == "" {
		return nil, fmt.Errorf("secret_key não configurado em [jwt]")
	}

	return config, nil
}

// loadCORSConfig carrega as configurações de CORS do arquivo ini
func loadCORSConfig(cfg *ini.File) CORSConfig {
	corsSection := cfg.Section("CORS")

	// Origens permitidas (separadas por vírgula)
	originsStr := corsSection.Key("AllowedOrigins").String()
	var origins []string
	if originsStr != "" {
		for _, origin := range splitAndTrim(originsStr, ",") {
			origins = append(origins, origin)
		}
	} else {
		origins = []string{} // Vazio = permite todas (*)
	}

	// Métodos permitidos
	methods := corsSection.Key("AllowedMethods").MustString("GET,POST,PUT,PATCH,DELETE,OPTIONS")

	// Headers permitidos
	headers := corsSection.Key("AllowedHeaders").MustString("Origin,Content-Type,Content-Length,Accept-Encoding,Authorization,Grant_type,X-CSRF-Token")

	// Allow Credentials
	allowCredentials := corsSection.Key("AllowCredentials").MustBool(true)

	// Max Age
	maxAge := corsSection.Key("MaxAge").MustString("43200") // 12 horas

	return CORSConfig{
		AllowedOrigins:   origins,
		AllowedMethods:   methods,
		AllowedHeaders:   headers,
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	}
}

// splitAndTrim divide uma string e remove espaços em branco
func splitAndTrim(s string, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

// splitString divide uma string pelo separador
func splitString(s string, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

// trimSpace remove espaços em branco do início e fim de uma string
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

// loadTLSConfig carrega as configurações de TLS do arquivo ini
func loadTLSConfig(cfg *ini.File) TLSConfig {
	tlsSection := cfg.Section("tls")

	return TLSConfig{
		Enabled:  tlsSection.Key("enabled").MustBool(false),
		CertFile: tlsSection.Key("cert_file").MustString("certs/server.crt"),
		KeyFile:  tlsSection.Key("key_file").MustString("certs/server.key"),
		Port:     tlsSection.Key("port").MustString("8443"),
	}
}

// loadSecurityConfig carrega as configurações de segurança do arquivo ini
func loadSecurityConfig(cfg *ini.File) SecurityConfig {
	secSection := cfg.Section("security")

	return SecurityConfig{
		MaxBodySize:      secSection.Key("max_body_size").MustInt64(1048576), // 1MB padrão
		RateLimitPerMin:  secSection.Key("rate_limit_per_min").MustInt(60),
		RateLimitPerHour: secSection.Key("rate_limit_per_hour").MustInt(1000),
		RateLimitEnabled: secSection.Key("rate_limit_enabled").MustBool(true),
	}
}
