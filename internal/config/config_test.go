package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConnectorDir(t *testing.T) {
	tempDir := t.TempDir()

	resolved, err := ResolveConnectorDir(tempDir)
	if err != nil {
		t.Fatalf("ResolveConnectorDir falhou: %v", err)
	}

	absExpected, _ := filepath.Abs(tempDir)
	if resolved != filepath.Clean(absExpected) {
		t.Errorf("Esperava %s, obteve %s", absExpected, resolved)
	}
}

func TestInitDataDir(t *testing.T) {
	tempDir := filepath.Join(t.TempDir(), "openapi-mcp-test")

	msgs, err := InitDataDir(tempDir)
	if err != nil {
		t.Fatalf("InitDataDir falhou: %v", err)
	}

	if len(msgs) == 0 {
		t.Errorf("Esperava mensagens de inicialização, recebeu lista vazia")
	}

	// Verificar se diretórios foram criados
	if _, err := os.Stat(filepath.Join(tempDir, "connectors")); os.IsNotExist(err) {
		t.Errorf("Diretório connectors não foi criado")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "logs")); os.IsNotExist(err) {
		t.Errorf("Diretório logs não foi criado")
	}
	dotEnvPath := filepath.Join(tempDir, ".env")
	if _, err := os.Stat(dotEnvPath); os.IsNotExist(err) {
		t.Errorf("Arquivo .env não foi criado")
	}

	// Testar idempotência (executar novamente)
	msgs2, err := InitDataDir(tempDir)
	if err != nil {
		t.Fatalf("Segunda execução de InitDataDir falhou: %v", err)
	}
	if len(msgs2) == 0 {
		t.Errorf("Esperava mensagens de preservação na segunda execução")
	}
}

func TestLoadPrecedence(t *testing.T) {
	tempDir := t.TempDir()

	// Escrever um .env de teste
	dotEnvContent := `
APP_ENV=staging
APP_LOG_LEVEL=debug
HTTP_TIMEOUT=45
TEST_VAR=from_dotenv
`
	err := os.WriteFile(filepath.Join(tempDir, ".env"), []byte(dotEnvContent), 0600)
	if err != nil {
		t.Fatalf("Falha ao criar .env de teste: %v", err)
	}

	// Definir variável no ambiente do SO
	t.Setenv("APP_ENV", "production_override")
	t.Setenv("SYSTEM_VAR", "from_system")

	cfg, err := Load(tempDir)
	if err != nil {
		t.Fatalf("Load falhou: %v", err)
	}

	// SO deve ter prevalecido sobre .env
	if cfg.AppEnv != "production_override" {
		t.Errorf("Precedência SO falhou. Esperava 'production_override', obteve '%s'", cfg.AppEnv)
	}

	// .env deve ser lido se o SO não sobrescrever
	if cfg.AppLogLevel != "debug" {
		t.Errorf("Leitura do .env falhou. Esperava 'debug', obteve '%s'", cfg.AppLogLevel)
	}

	if cfg.HttpTimeout != 45 {
		t.Errorf("HTTP_TIMEOUT falhou. Esperava 45, obteve %d", cfg.HttpTimeout)
	}

	if cfg.GetEnv("TEST_VAR") != "from_dotenv" {
		t.Errorf("GetEnv TEST_VAR falhou. Esperava 'from_dotenv', obteve '%s'", cfg.GetEnv("TEST_VAR"))
	}

	if cfg.GetEnv("SYSTEM_VAR") != "from_system" {
		t.Errorf("GetEnv SYSTEM_VAR falhou. Esperava 'from_system', obteve '%s'", cfg.GetEnv("SYSTEM_VAR"))
	}
}

func TestAbsPaths(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := Load(tempDir)
	if err != nil {
		t.Fatalf("Load falhou: %v", err)
	}

	expectedConnectors := filepath.Clean(filepath.Join(tempDir, "connectors"))
	if cfg.AbsConnectorsDir() != expectedConnectors {
		t.Errorf("AbsConnectorsDir incorreto. Esperava %s, obteve %s", expectedConnectors, cfg.AbsConnectorsDir())
	}

	expectedLogs := filepath.Clean(filepath.Join(tempDir, "logs"))
	if cfg.AbsLogsDir() != expectedLogs {
		t.Errorf("AbsLogsDir incorreto. Esperava %s, obteve %s", expectedLogs, cfg.AbsLogsDir())
	}
}
