package main

import (
	"fmt"
	"os"
	"wsicrmrest/internal/config"
	"wsicrmrest/internal/database"
	"wsicrmrest/internal/logger"
	"wsicrmrest/internal/middleware"
	"wsicrmrest/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title WSICRMREST API
// @version 1.26.4.27
// @description REST API service for CRM integration, converted from WinDev. Provides JWT-based authentication and webhook endpoints for Zenvia email and SMS events.
// @description
// @description **Security Features:**
// @description - SQL Injection Protection (bind variables)
// @description - HTTPS/TLS Support (configurable)
// @description - Rate Limiting (60 req/min, 1000 req/hour by default)
// @description - CORS with origin validation
// @description - Request size limits (1MB max body)
// @description - Request timeout (30s default)
// @description - Security headers (HSTS, X-Frame-Options, etc.)
// @description
// @description **Important Notes:**
// @description - All API endpoints require JWT Bearer token (except /token and /wsteste)
// @description - Webhooks from Zenvia are public (no authentication required)
// @description - Rate limit headers are included in all responses
// @description - HTTPS is strongly recommended for production
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.basic BasicAuth
// @type basic
// @description Basic authentication with client_id as username and client_secret as password (for token generation only).

// @x-extension-openapi {"example": "value"}
// @tag.name Authentication
// @tag.description Endpoints for JWT token generation and management
// @tag.name Testing
// @tag.description Endpoints for testing database connectivity and API health
// @tag.name Webhooks
// @tag.description Endpoints for receiving Zenvia webhook events (email and SMS status updates)

func main() {
	// Carregar configurações do dbinit.ini
	cfg, err := config.LoadConfig("dbinit.ini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao carregar configurações: %v\n", err)
		os.Exit(1)
	}

	// Inicializar logger
	log, err := logger.NewLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Iniciando WSICRMREST",
		"version", config.Version,
		"version_date", config.VersionDate,
		"build_time", config.BuildTime)

	// Inicializar conexão com banco de dados
	db, err := database.NewDatabase(cfg, log)
	if err != nil {
		log.Error("Erro ao conectar ao banco de dados", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	log.Info("Conexão com banco de dados estabelecida com sucesso")

	// Carregar dados do organizador (obrigatório)
	if err := db.LeOrganizador(cfg); err != nil {
		log.Error("Erro ao carregar dados do organizador", "error", err)
		log.Error("Verifique se a tabela ORGANIZADOR existe e possui ao menos um registro com OrgCodigo > 0")
		os.Exit(1)
	}

	log.Info("Dados do organizador carregados com sucesso",
		"codigo", cfg.Organization.Codigo,
		"nome", cfg.Organization.Nome)

	// Log configurações CORS
	if len(cfg.CORS.AllowedOrigins) == 0 {
		log.Info("CORS configurado para permitir TODAS as origens (*) - Modo Desenvolvimento")
	} else {
		log.Info("CORS configurado com origens restritas",
			"allowed_origins", cfg.CORS.AllowedOrigins)
	}
	log.Debug("Configurações CORS",
		"methods", cfg.CORS.AllowedMethods,
		"headers", cfg.CORS.AllowedHeaders,
		"credentials", cfg.CORS.AllowCredentials,
		"max_age", cfg.CORS.MaxAge)

	// Configurar Gin
	if cfg.Application.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Aplicar middlewares de segurança
	router.Use(middleware.SecurityMiddleware(cfg))
	router.Use(middleware.RateLimitMiddleware(cfg))

	// Aplicar middleware CORS com validação de produção
	router.Use(middleware.CORS(cfg.CORS, cfg.Application.Environment, log))

	// Log configurações de segurança
	if cfg.Security.RateLimitEnabled {
		log.Info("Rate limiting habilitado",
			"per_min", cfg.Security.RateLimitPerMin,
			"per_hour", cfg.Security.RateLimitPerHour)
	}
	log.Info("Segurança configurada",
		"max_body_size", cfg.Security.MaxBodySize,
		"request_timeout", cfg.Application.RequestTimeout)

	// Configurar rotas
	routes.SetupRoutes(router, cfg, db, log)

	// Iniciar servidor com ou sem TLS
	port := cfg.Application.Port
	if port == "" {
		port = "8080"
	}

	if cfg.TLS.Enabled {
		// Validar existência dos certificados
		if _, err := os.Stat(cfg.TLS.CertFile); os.IsNotExist(err) {
			log.Error("Certificado TLS não encontrado", "cert_file", cfg.TLS.CertFile)
			os.Exit(1)
		}
		if _, err := os.Stat(cfg.TLS.KeyFile); os.IsNotExist(err) {
			log.Error("Chave privada TLS não encontrada", "key_file", cfg.TLS.KeyFile)
			os.Exit(1)
		}

		tlsPort := cfg.TLS.Port
		if tlsPort == "" {
			tlsPort = "8443"
		}

		log.Info("🔒 Servidor HTTPS/TLS iniciado",
			"port", tlsPort,
			"cert", cfg.TLS.CertFile,
			"environment", cfg.Application.Environment)

		if err := router.RunTLS(":"+tlsPort, cfg.TLS.CertFile, cfg.TLS.KeyFile); err != nil {
			log.Error("Erro ao iniciar servidor HTTPS", "error", err)
			os.Exit(1)
		}
	} else {
		log.Info("Servidor HTTP iniciado", "port", port, "environment", cfg.Application.Environment)
		log.Warn("⚠️  TLS/HTTPS desabilitado - dados trafegam sem criptografia")

		if err := router.Run(":" + port); err != nil {
			log.Error("Erro ao iniciar servidor", "error", err)
			os.Exit(1)
		}
	}
}
