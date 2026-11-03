package app

import (
	"flag"
	"fmt"
	"os"

	"openapi-mcp/internal/config"
	"openapi-mcp/internal/logger"
)

const Version = "0.1.0-dev"

// Run é o ponto de entrada principal para a execução da CLI e servidor.
// Retorna o código de saída do processo (0 para sucesso, != 0 para erro).
func Run(args []string) int {
	if len(args) > 1 && args[1] == "init" {
		return runInit(args[2:])
	}

	fs := flag.NewFlagSet("openapi-mcp", flag.ContinueOnError)
	var dataDir string
	var showVersion bool
	var showHelp bool

	fs.StringVar(&dataDir, "data-dir", "", "Caminho personalizado para o diretório de dados")
	fs.StringVar(&dataDir, "d", "", "Caminho personalizado para o diretório de dados (abreviação)")
	fs.BoolVar(&showVersion, "version", false, "Exibe a versão do programa")
	fs.BoolVar(&showVersion, "v", false, "Exibe a versão do programa (abreviação)")
	fs.BoolVar(&showHelp, "help", false, "Exibe esta mensagem de ajuda")
	fs.BoolVar(&showHelp, "h", false, "Exibe esta mensagem de ajuda (abreviação)")

	if err := fs.Parse(args[1:]); err != nil {
		return 1
	}

	if showHelp {
		printHelp()
		return 0
	}

	if showVersion {
		fmt.Printf("openapi-mcp versão %s\n", Version)
		return 0
	}

	// Execução principal do servidor
	resolvedDataDir, err := config.ResolveDataDir(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao resolver diretório de dados: %v\n", err)
		return 1
	}

	cfg, err := config.Load(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao carregar configurações de %s: %v\n", resolvedDataDir, err)
		return 1
	}

	log, err := logger.InitLogger(cfg.AppLogLevel, cfg.AbsLogsDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
		return 1
	}
	defer log.Close()

	log.Infof("Iniciando OpenAPI MCP v%s...", Version)
	log.Infof("Diretório de dados: %s", cfg.DataDir)
	log.Infof("Diretório de conectores: %s", cfg.AbsConnectorsDir())
	log.Infof("Diretório de logs: %s", cfg.AbsLogsDir())
	log.Infof("Ambiente: %s | Transport: %s", cfg.AppEnv, cfg.McpTransport)

	// Etapa 1 concluída: o ambiente está inicializado e validado.
	log.Infof("Fundação inicializada com sucesso.")

	return 0
}

func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	var dataDir string

	fs.StringVar(&dataDir, "data-dir", "", "Caminho personalizado para o diretório de dados")
	fs.StringVar(&dataDir, "d", "", "Caminho personalizado para o diretório de dados (abreviação)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	resolvedDataDir, err := config.ResolveDataDir(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao resolver diretório de dados: %v\n", err)
		return 1
	}

	msgs, err := config.InitDataDir(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar estrutura de dados: %v\n", err)
		return 1
	}

	fmt.Println("Inicialização concluída com sucesso:")
	for _, msg := range msgs {
		fmt.Printf(" - %s\n", msg)
	}

	return 0
}

func printHelp() {
	help := `OpenAPI MCP - Servidor MCP para especificações OpenAPI

Uso:
  openapi-mcp [opções]
  openapi-mcp init [opções]

Comandos:
  init            Cria a estrutura inicial de diretórios e o arquivo .env no diretório de dados.

Opções:
  -d, --data-dir  Caminho personalizado para o diretório de dados.
                  Padrão: ~/.openapi-mcp (Linux/macOS) ou (USERPROFILE)\.openapi-mcp (Windows).
  -v, --version   Exibe a versão do programa.
  -h, --help      Exibe esta mensagem de ajuda.
`
	fmt.Print(help)
}
