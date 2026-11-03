package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	tempDir := t.TempDir()

	l, err := InitLogger("debug", tempDir)
	if err != nil {
		t.Fatalf("InitLogger falhou: %v", err)
	}
	defer l.Close()

	l.Infof("Mensagem de informação para teste")
	l.Debugf("Mensagem de debug para teste")

	logFile := filepath.Join(tempDir, "server.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Falha ao ler arquivo de log: %v", err)
	}

	strContent := string(content)
	if !strings.Contains(strContent, "[INFO] Mensagem de informação para teste") {
		t.Errorf("Arquivo de log não contém a mensagem [INFO]")
	}
	if !strings.Contains(strContent, "[DEBUG] Mensagem de debug para teste") {
		t.Errorf("Arquivo de log não contém a mensagem [DEBUG]")
	}
}
