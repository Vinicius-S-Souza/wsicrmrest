package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DailyRotatingWriter é um writer que rotaciona o arquivo de log diariamente
type DailyRotatingWriter struct {
	logDir      string
	currentDate string
	file        *os.File
	mu          sync.Mutex
}

// NewDailyRotatingWriter cria um novo writer com rotação diária
func NewDailyRotatingWriter(logDir string) (*DailyRotatingWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de logs: %w", err)
	}

	writer := &DailyRotatingWriter{
		logDir: logDir,
	}

	// Abrir arquivo inicial
	if err := writer.rotate(); err != nil {
		return nil, err
	}

	return writer, nil
}

// Write implementa io.Writer
func (w *DailyRotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Verificar se a data mudou
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != w.currentDate {
		// Rotacionar arquivo
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	return w.file.Write(p)
}

// Sync implementa zapcore.WriteSyncer
func (w *DailyRotatingWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// rotate fecha o arquivo atual e abre um novo para a data de hoje
func (w *DailyRotatingWriter) rotate() error {
	// Fechar arquivo anterior
	if w.file != nil {
		w.file.Close()
	}

	// Atualizar data atual
	w.currentDate = time.Now().Format("2006-01-02")

	// Criar nome do novo arquivo
	logFileName := fmt.Sprintf("wsicrmrest_%s.log", w.currentDate)
	logFilePath := filepath.Join(w.logDir, logFileName)

	// Abrir novo arquivo
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %w", err)
	}

	w.file = file
	return nil
}

// Close fecha o arquivo de log
func (w *DailyRotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// NewLogger cria um novo logger configurado para gravar em arquivo
// O arquivo será gravado na pasta "log" com o nome "wsicrmrest_YYYY-MM-DD.log"
// O arquivo é rotacionado automaticamente à meia-noite
func NewLogger() (*zap.SugaredLogger, error) {
	// Criar writer com rotação diária
	logDir := "log"
	fileWriter, err := NewDailyRotatingWriter(logDir)
	if err != nil {
		return nil, err
	}

	// Configurar encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Core que escreve no arquivo (com rotação diária)
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(fileWriter),
		zapcore.InfoLevel,
	)

	// Core que escreve no console (para debug)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	// Combinar os dois cores
	core := zapcore.NewTee(fileCore, consoleCore)

	// Criar logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger.Sugar(), nil
}

// GetLogFileName retorna o nome do arquivo de log para a data especificada
func GetLogFileName(date time.Time) string {
	return fmt.Sprintf("wsicrmrest_%s.log", date.Format("2006-01-02"))
}
