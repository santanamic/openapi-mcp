package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config guarda as configurações gerais da aplicação.
type Config struct {
	ConnectorDir  string
	AppEnv        string
	AppLogLevel   string
	McpTransport  string
	ConnectorsDir string
	LogsDir       string
	HttpTimeout   int
	envMap        map[string]string
}

// DefaultConnectorDir retorna o caminho padrão do diretório de conectores (~/.openapi-mcp no Unix, %USERPROFILE%\.openapi-mcp no Windows).
func DefaultConnectorDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("falha ao obter diretório home do usuário: %w", err)
	}
	return filepath.Join(home, ".openapi-mcp"), nil
}

// ResolveConnectorDir normaliza e converte em caminho absoluto o diretório de conectores fornecido ou padrão.
func ResolveConnectorDir(flagConnectorDir string) (string, error) {
	var rawDir string
	if flagConnectorDir != "" {
		rawDir = flagConnectorDir
	} else {
		var err error
		rawDir, err = DefaultConnectorDir()
		if err != nil {
			return "", err
		}
	}
	absDir, err := filepath.Abs(rawDir)
	if err != nil {
		return "", fmt.Errorf("falha ao resolver caminho absoluto de %s: %w", rawDir, err)
	}
	return filepath.Clean(absDir), nil
}

// Load carrega a configuração a partir do diretório de conectores especificado.
// Respeita a precedência:
// 1. Variáveis no ambiente do processo (os.Getenv)
// 2. Arquivo .env no diretório de conectores
// 3. Valores padrão
func Load(connectorDir string) (*Config, error) {
	absConnectorDir, err := filepath.Abs(connectorDir)
	if err != nil {
		return nil, fmt.Errorf("diretório de conectores inválido: %w", err)
	}

	envFileMap := make(map[string]string)
	dotEnvPath := filepath.Join(absConnectorDir, ".env")

	if _, err := os.Stat(dotEnvPath); err == nil {
		parsed, err := parseDotEnv(dotEnvPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler arquivo .env (%s): %w", dotEnvPath, err)
		}
		envFileMap = parsed
	}

	cfg := &Config{
		ConnectorDir: absConnectorDir,
		envMap:       make(map[string]string),
	}

	// Carrega todas as chaves do .env no mapa interno
	for k, v := range envFileMap {
		cfg.envMap[k] = v
	}

	// Sobrescreve/complementa com o ambiente do SO
	for _, envStr := range os.Environ() {
		parts := strings.SplitN(envStr, "=", 2)
		if len(parts) == 2 {
			cfg.envMap[parts[0]] = parts[1]
		}
	}

	cfg.AppEnv = cfg.getEnvOrDefault("APP_ENV", "production")
	cfg.AppLogLevel = cfg.getEnvOrDefault("APP_LOG_LEVEL", "info")
	cfg.McpTransport = cfg.getEnvOrDefault("MCP_TRANSPORT", "stdio")
	cfg.ConnectorsDir = cfg.getEnvOrDefault("CONNECTORS_DIR", "./connectors")
	cfg.LogsDir = cfg.getEnvOrDefault("LOGS_DIR", "./logs")

	timeoutStr := cfg.getEnvOrDefault("HTTP_TIMEOUT", "30")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout <= 0 {
		timeout = 30
	}
	cfg.HttpTimeout = timeout

	return cfg, nil
}

// AbsConnectorsDir resolve o caminho de CONNECTORS_DIR em relação ao ConnectorDir se for relativo.
func (c *Config) AbsConnectorsDir() string {
	if filepath.IsAbs(c.ConnectorsDir) {
		return filepath.Clean(c.ConnectorsDir)
	}
	return filepath.Clean(filepath.Join(c.ConnectorDir, c.ConnectorsDir))
}

// AbsLogsDir resolve o caminho de LOGS_DIR em relação ao ConnectorDir se for relativo.
func (c *Config) AbsLogsDir() string {
	if filepath.IsAbs(c.LogsDir) {
		return filepath.Clean(c.LogsDir)
	}
	return filepath.Clean(filepath.Join(c.ConnectorDir, c.LogsDir))
}

// GetEnv obtém o valor de uma variável considerando SO > .env.
func (c *Config) GetEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return c.envMap[key]
}

func (c *Config) getEnvOrDefault(key, defaultValue string) string {
	if val := c.GetEnv(key); val != "" {
		return val
	}
	return defaultValue
}

// parseDotEnv lê um arquivo no formato KEY=VALUE simples.
func parseDotEnv(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			result[key] = val
		}
	}
	return result, scanner.Err()
}

const defaultDotEnvTemplate = `# Configuração Geral do OpenAPI MCP
APP_ENV=production
APP_LOG_LEVEL=info
MCP_TRANSPORT=stdio
CONNECTORS_DIR=./connectors
LOGS_DIR=./logs
HTTP_TIMEOUT=30
`

// InitDataDir cria a estrutura inicial de diretórios e arquivos se não existirem.
func InitDataDir(dataDir string) ([]string, error) {
	absDataDir, err := filepath.Abs(dataDir)
	if err != nil {
		return nil, fmt.Errorf("caminho de diretório inválido: %w", err)
	}

	var messages []string

	// 1. Criar dataDir
	if err := os.MkdirAll(absDataDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de dados %s: %w", absDataDir, err)
	}
	messages = append(messages, fmt.Sprintf("Diretório de dados: %s", absDataDir))

	// 2. Criar connectors dir
	connectorsDir := filepath.Join(absDataDir, "connectors")
	if _, err := os.Stat(connectorsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(connectorsDir, 0755); err != nil {
			return nil, fmt.Errorf("falha ao criar diretório connectors: %w", err)
		}
		messages = append(messages, fmt.Sprintf("Criado diretório: %s", connectorsDir))
	} else {
		messages = append(messages, fmt.Sprintf("Preservado diretório existente: %s", connectorsDir))
	}

	// 3. Criar logs dir
	logsDir := filepath.Join(absDataDir, "logs")
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			return nil, fmt.Errorf("falha ao criar diretório logs: %w", err)
		}
		messages = append(messages, fmt.Sprintf("Criado diretório: %s", logsDir))
	} else {
		messages = append(messages, fmt.Sprintf("Preservado diretório existente: %s", logsDir))
	}

	// 4. Criar .env se não existir
	dotEnvPath := filepath.Join(absDataDir, ".env")
	if _, err := os.Stat(dotEnvPath); os.IsNotExist(err) {
		if err := os.WriteFile(dotEnvPath, []byte(defaultDotEnvTemplate), 0600); err != nil {
			return nil, fmt.Errorf("falha ao criar arquivo .env: %w", err)
		}
		messages = append(messages, fmt.Sprintf("Criado arquivo de exemplo: %s", dotEnvPath))
	} else {
		messages = append(messages, fmt.Sprintf("Preservado arquivo existente: %s", dotEnvPath))
	}

	return messages, nil
}
