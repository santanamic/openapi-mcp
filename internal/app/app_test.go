package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunVersion(t *testing.T) {
	code := Run([]string{"openapi-mcp", "--version"})
	if code != 0 {
		t.Errorf("Esperava código 0 para --version, obteve %d", code)
	}
}

func TestRunHelp(t *testing.T) {
	code := Run([]string{"openapi-mcp", "-h"})
	if code != 0 {
		t.Errorf("Esperava código 0 para -h, obteve %d", code)
	}
}

func TestRunInitAndServer(t *testing.T) {
	tempDir := filepath.Join(t.TempDir(), "mcp-app-test")

	// 1. Executar init
	codeInit := Run([]string{"openapi-mcp", "init", "--data-dir", tempDir})
	if codeInit != 0 {
		t.Fatalf("Init falhou com código %d", codeInit)
	}

	// Verificar se arquivos foram criados
	if _, err := os.Stat(filepath.Join(tempDir, ".env")); os.IsNotExist(err) {
		t.Errorf("Arquivo .env não encontrado após init")
	}

	// 2. Executar servidor com o data-dir inicializado
	codeRun := Run([]string{"openapi-mcp", "--data-dir", tempDir})
	if codeRun != 0 {
		t.Errorf("Run do servidor falhou com código %d", codeRun)
	}

	// Verificar se log foi gravado
	logFile := filepath.Join(tempDir, "logs", "server.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Arquivo server.log não encontrado no diretório de logs")
	}
}
