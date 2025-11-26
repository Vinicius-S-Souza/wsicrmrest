//go:build windows
// +build windows

package service

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"wsicrmrest/internal/config"
	"wsicrmrest/internal/database"
	"wsicrmrest/internal/logger"
	"wsicrmrest/internal/middleware"
	"wsicrmrest/internal/routes"
	tlsloader "wsicrmrest/internal/tls"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

// WindowsService implementa a interface do Windows Service
type WindowsService struct {
	config     *config.Config
	log        *zap.SugaredLogger
	db         *database.Database
	eventLog   *eventlog.Log
	stopChan   chan struct{}
	serverDone chan error
}

// Execute é chamado pelo Service Control Manager do Windows
func (ws *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}

	// Iniciar o servidor HTTP em goroutine
	go ws.runHTTPServer()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	ws.log.Info("=======================================")
	ws.logEvent(eventlog.Info, "Serviço WSICRMREST iniciado com sucesso")
	ws.log.Info("=======================================")

	// Loop principal do serviço
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				ws.logEvent(eventlog.Info, "Recebido comando de parada do serviço")
				changes <- svc.Status{State: svc.StopPending}
				ws.stop()
				break loop
			default:
				ws.logEvent(eventlog.Error, fmt.Sprintf("Comando não esperado: %d", c.Cmd))
			}
		case err := <-ws.serverDone:
			if err != nil {
				ws.logEvent(eventlog.Error, fmt.Sprintf("Servidor HTTP falhou: %v", err))
				changes <- svc.Status{State: svc.StopPending}
				ws.stop()
				break loop
			}
		}
	}

	changes <- svc.Status{State: svc.Stopped}
	ws.log.Info("=======================================")
	ws.logEvent(eventlog.Info, "Serviço WSICRMREST parado")
	ws.log.Info("=======================================")
	return
}

// runHTTPServer inicia o servidor HTTP/HTTPS
func (ws *WindowsService) runHTTPServer() {
	defer close(ws.serverDone)

	// Configurar Gin
	if ws.config.Application.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Aplicar middlewares de segurança
	router.Use(middleware.SecurityMiddleware(ws.config))
	router.Use(middleware.Fail2BanMiddleware(ws.log)) // Proteção contra ataques de força bruta
	router.Use(middleware.RateLimitMiddleware(ws.config))
	router.Use(middleware.CORS(ws.config.CORS, ws.config.Application.Environment, ws.log))

	// Configurar rotas
	routes.SetupRoutes(router, ws.config, ws.db, ws.log)

	// Iniciar servidor com ou sem TLS
	port := ws.config.Application.Port
	if port == "" {
		port = "8080"
	}

	var err error
	if ws.config.TLS.Enabled {
		// Validar existência dos certificados
		if _, err := os.Stat(ws.config.TLS.CertFile); os.IsNotExist(err) {
			ws.log.Error("Certificado TLS não encontrado", "cert_file", ws.config.TLS.CertFile)
			ws.serverDone <- fmt.Errorf("certificado TLS não encontrado: %s", ws.config.TLS.CertFile)
			return
		}
		if _, err := os.Stat(ws.config.TLS.KeyFile); os.IsNotExist(err) {
			ws.log.Error("Chave privada TLS não encontrada", "key_file", ws.config.TLS.KeyFile)
			ws.serverDone <- fmt.Errorf("chave privada TLS não encontrada: %s", ws.config.TLS.KeyFile)
			return
		}

		tlsPort := ws.config.TLS.Port
		if tlsPort == "" {
			tlsPort = "8443"
		}

		// Carregar configuração TLS (suporta chaves criptografadas)
		tlsConfig, err := tlsloader.LoadEncryptedTLSConfig(ws.config.TLS.CertFile, ws.config.TLS.KeyFile, ws.config.TLS.KeyPassword)
		if err != nil {
			ws.log.Error("Erro ao carregar certificado TLS", "error", err)
			ws.serverDone <- fmt.Errorf("erro ao carregar certificado TLS: %w", err)
			return
		}

		// Criar servidor HTTP customizado com TLS
		server := &http.Server{
			Addr:      ":" + tlsPort,
			Handler:   router,
			TLSConfig: tlsConfig,
		}

		ws.log.Info("====================================================================")
		ws.log.Info("🔒 Servidor HTTPS/TLS iniciado",
			" - porta: ", tlsPort,
			" - cert: ", ws.config.TLS.CertFile,
			" - environment: ", ws.config.Application.Environment)
		ws.log.Info("====================================================================")

		err = server.ListenAndServeTLS("", "")
	} else {
		ws.log.Info("====================================================================")
		ws.log.Info("Servidor HTTP iniciado", " - porta: ", port, " - environment: ", ws.config.Application.Environment)
		ws.log.Warn("⚠️  TLS/HTTPS desabilitado - dados trafegam sem criptografia")
		ws.log.Info("====================================================================")

		err = router.Run(":" + port)
	}

	if err != nil {
		ws.serverDone <- err
	}
}

// stop para o serviço graciosamente
func (ws *WindowsService) stop() {
	ws.log.Info("====================================================================")
	ws.log.Info("Parando serviço WSICRMREST...")
	ws.log.Info("====================================================================")

	// Fechar conexão com banco de dados
	if ws.db != nil {
		ws.db.Close()
		ws.log.Info("Conexão com banco de dados fechada")
	}

	// Sincronizar logs
	if ws.log != nil {
		ws.log.Sync()
	}

	close(ws.stopChan)
}

// logEvent registra eventos no Event Log do Windows
func (ws *WindowsService) logEvent(etype uint16, msg string) {
	if ws.eventLog != nil {
		ws.eventLog.Info(1, msg)
	}
	if ws.log != nil {
		switch etype {
		case eventlog.Error:
			ws.log.Error(msg)
		case eventlog.Warning:
			ws.log.Warn(msg)
		default:
			ws.log.Info(msg)
		}
	}
}

// RunAsWindowsService executa a aplicação como Windows Service
func RunAsWindowsService() error {
	serviceName := "WSICRMREST"

	// Verificar se está rodando como serviço
	isService, err := svc.IsWindowsService()
	if err != nil {
		return fmt.Errorf("falha ao verificar se é serviço Windows: %v", err)
	}

	if !isService {
		return fmt.Errorf("não está sendo executado como serviço Windows")
	}

	// Mudar para o diretório do executável
	// Serviços Windows iniciam em C:\Windows\System32 por padrão
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		os.Chdir(exeDir)
	}

	// Configurar Event Log do Windows
	elog, err := eventlog.Open(serviceName)
	if err != nil {
		// Se não conseguir abrir, continua sem event log
		elog = nil
	}
	defer func() {
		if elog != nil {
			elog.Close()
		}
	}()

	// Carregar configurações do dbinit.ini
	cfg, err := config.LoadConfig("dbinit.ini")
	if err != nil {
		if elog != nil {
			elog.Error(1, fmt.Sprintf("Erro ao carregar configurações: %v", err))
		}
		return fmt.Errorf("erro ao carregar configurações: %v", err)
	}

	// Inicializar logger
	log, err := logger.NewLogger()
	if err != nil {
		if elog != nil {
			elog.Error(1, fmt.Sprintf("Erro ao inicializar logger: %v", err))
		}
		return fmt.Errorf("erro ao inicializar logger: %v", err)
	}
	defer log.Sync()

	log.Info("Iniciando WSICRMREST como Windows Service",
		"version", config.Version,
		"version_date", config.VersionDate,
		"build_time", config.BuildTime)

	// Inicializar conexão com banco de dados
	db, err := database.NewDatabase(cfg, log)
	if err != nil {
		log.Error("Erro ao conectar ao banco de dados", "error", err)
		if elog != nil {
			elog.Error(1, fmt.Sprintf("Erro ao conectar ao banco de dados: %v", err))
		}
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}

	log.Info("Conexão com banco de dados estabelecida com sucesso")

	// Carregar dados do organizador (obrigatório)
	if err := db.LeOrganizador(cfg); err != nil {
		log.Error("Erro ao carregar dados do organizador", "error", err)
		log.Error("Verifique se a tabela ORGANIZADOR existe e possui ao menos um registro com OrgCodigo > 0")
		if elog != nil {
			elog.Error(1, fmt.Sprintf("Erro ao carregar organizador: %v", err))
		}
		return fmt.Errorf("erro ao carregar organizador: %v", err)
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

	// Criar instância do serviço
	ws := &WindowsService{
		config:     cfg,
		log:        log,
		db:         db,
		eventLog:   elog,
		stopChan:   make(chan struct{}),
		serverDone: make(chan error, 1),
	}

	// Executar como serviço Windows
	if err := svc.Run(serviceName, ws); err != nil {
		log.Error("Erro ao executar serviço Windows", "error", err)
		if elog != nil {
			elog.Error(1, fmt.Sprintf("Erro ao executar serviço: %v", err))
		}
		return err
	}

	return nil
}

// InstallEventLog instala o Event Log para o serviço (deve ser chamado durante instalação)
func InstallEventLog(serviceName string) error {
	err := eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		return fmt.Errorf("falha ao instalar event log: %v", err)
	}
	return nil
}

// UninstallEventLog remove o Event Log do serviço (deve ser chamado durante desinstalação)
func UninstallEventLog(serviceName string) error {
	err := eventlog.Remove(serviceName)
	if err != nil {
		return fmt.Errorf("falha ao remover event log: %v", err)
	}
	return nil
}
