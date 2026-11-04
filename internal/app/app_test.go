package app

import (
	"os"
	"path/filepath"
	"testing"

	"openapi-mcp/internal/config"
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

func TestDataDirInitialization(t *testing.T) {
	tempDir := filepath.Join(t.TempDir(), "mcp-app-test")

	msgs, err := config.InitDataDir(tempDir)
	if err != nil {
		t.Fatalf("InitDataDir falhou: %v", err)
	}
	if len(msgs) == 0 {
		t.Errorf("Esperava mensagens de inicialização")
	}

	// Verificar se arquivos foram criados
	if _, err := os.Stat(filepath.Join(tempDir, ".env")); os.IsNotExist(err) {
		t.Errorf("Arquivo .env não encontrado após inicialização")
	}
}
