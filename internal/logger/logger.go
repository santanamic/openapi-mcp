package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func parseLevel(lvl string) Level {
	switch strings.ToLower(strings.TrimSpace(lvl)) {
	case "debug":
		return LevelDebug
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

type Logger struct {
	mu     sync.Mutex
	level  Level
	logger *log.Logger
	file   *os.File
}

var globalLogger *Logger
var once sync.Once

// InitLogger inicializa o logger global direcionando saída exclusivamente para stderr e/ou arquivo no logsDir.
// STDOUT NUNCA é utilizado para logs, garantindo integridade do transporte stdio do MCP.
func InitLogger(logLevel, logsDir string) (*Logger, error) {
	level := parseLevel(logLevel)

	writers := []io.Writer{os.Stderr}

	var logFile *os.File
	if logsDir != "" {
		if err := os.MkdirAll(logsDir, 0755); err == nil {
			filePath := filepath.Join(logsDir, "server.log")
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if err == nil {
				logFile = f
				writers = append(writers, f)
			}
		}
	}

	multi := io.MultiWriter(writers...)
	l := log.New(multi, "", log.LstdFlags)

	instance := &Logger{
		level:  level,
		logger: l,
		file:   logFile,
	}

	once.Do(func() {
		globalLogger = instance
	})

	return instance, nil
}

// Get retorna o logger global configurado. Se não inicializado, retorna um default para stderr.
func Get() *Logger {
	if globalLogger == nil {
		InitLogger("info", "")
	}
	return globalLogger
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		return err
	}
	return nil
}

func (l *Logger) Debugf(format string, v ...any) {
	if l.level <= LevelDebug {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Infof(format string, v ...any) {
	if l.level <= LevelInfo {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Warnf(format string, v ...any) {
	if l.level <= LevelWarn {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	if l.level <= LevelError {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}
