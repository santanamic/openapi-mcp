package app

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"openapi-mcp/internal/config"
	"openapi-mcp/internal/connector"
	"openapi-mcp/internal/logger"
	"openapi-mcp/internal/mcpserver"
	"openapi-mcp/internal/openapi"
	"openapi-mcp/internal/tools"
)

const Version = "1.0.0"

type envVarFlag []string

func (e *envVarFlag) String() string {
	return strings.Join(*e, ", ")
}

func (e *envVarFlag) Set(value string) error {
	*e = append(*e, value)
	return nil
}

// Run é o ponto de entrada principal para a execução da CLI e servidor.
// Retorna o código de saída do processo (0 para sucesso, != 0 para erro).
func Run(args []string) int {
	fs := flag.NewFlagSet("openapi-mcp", flag.ContinueOnError)
	var connectorDir string
	var connectorPath string
	var envFilePath string
	var logFilePath string
	var envFlags envVarFlag
	var showVersion bool
	var showHelp bool

	fs.StringVar(&connectorDir, "connector-dir", "", "Caminho personalizado para o diretório de conectores")
	fs.StringVar(&connectorDir, "d", "", "Caminho personalizado para o diretório de conectores (abreviação)")
	fs.StringVar(&connectorPath, "connector", "", "Caminho local ou URL para um arquivo connector.yaml específico")
	fs.StringVar(&connectorPath, "c", "", "Caminho local ou URL para um arquivo connector.yaml específico (abreviação)")
	fs.StringVar(&envFilePath, "env-file", "", "Caminho para um arquivo .env personalizado")
	fs.StringVar(&envFilePath, "e", "", "Caminho para um arquivo .env personalizado (abreviação)")
	fs.StringVar(&logFilePath, "log-file", "", "Caminho personalizado para o arquivo de log")
	fs.StringVar(&logFilePath, "l", "", "Caminho personalizado para o arquivo de log (abreviação)")
	fs.Var(&envFlags, "env", "Define ou sobrescreve variáveis de ambiente (formato CHAVE=VALOR ou CHAVE1=VALOR1,CHAVE2=VALOR2)")
	fs.BoolVar(&showVersion, "version", false, "Exibe a versão do programa")
	fs.BoolVar(&showVersion, "version-short", false, "Exibe a versão do programa")
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

	// 1. Carregar arquivo .env personalizado se especificado via --env-file / -e
	if envFilePath != "" {
		loadRootEnv(envFilePath)
	}

	// 2. Injetar variáveis passadas via CLI (--env "KEY1=VAL1,KEY2=VAL2")
	for _, envArg := range envFlags {
		pairs := strings.Split(envArg, ",")
		for _, pair := range pairs {
			pair = strings.TrimSpace(pair)
			if pair == "" {
				continue
			}
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				_ = os.Setenv(key, val)
			}
		}
	}

	// 4. Carregar .env padrão da raiz / working directory se existir
	loadRootEnv(".env")

	registry := tools.NewRegistry()
	mapper := tools.NewMapper()
	connectors := make(map[string]*connector.Manifest)

	if connectorPath != "" {
		var manifestPath string
		var cleanupTemp func()

		// Verificar se o conector é uma URL remota ou arquivo local
		if strings.HasPrefix(connectorPath, "http://") || strings.HasPrefix(connectorPath, "https://") {
			tmpPath, err := fetchRemoteToTemp(connectorPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao baixar conector remoto de %s: %v\n", connectorPath, err)
				return 1
			}
			manifestPath = tmpPath
			cleanupTemp = func() { _ = os.Remove(manifestPath) }
			defer cleanupTemp()
		} else {
			absPath, err := filepath.Abs(connectorPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Caminho de conector inválido: %v\n", err)
				return 1
			}
			manifestPath = absPath
			// Carregar .env do diretório do conector se existir
			loadRootEnv(filepath.Join(filepath.Dir(manifestPath), ".env"))
		}

		manifest, err := connector.LoadSingleManifest(manifestPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao carregar manifesto do conector: %v\n", err)
			return 1
		}

		connectors[manifest.ID] = manifest

		absSpecPath, err := resolveSpecPath(connectorPath, manifest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao resolver especificação OpenAPI: %v\n", err)
			return 1
		}

		doc, err := openapi.LoadDocument(absSpecPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao carregar documento OpenAPI %s: %v\n", absSpecPath, err)
			return 1
		}

		toolDefs, err := mapper.MapConnector(manifest, doc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao mapear ferramentas do conector %s: %v\n", manifest.ID, err)
			return 1
		}

		for _, td := range toolDefs {
			if err := registry.Register(td); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao registrar ferramenta %s: %v\n", td.Name, err)
				return 1
			}
		}

	} else {
		// Modo Padrão (Diretório de Conectores)
		resolvedConnectorDir, err := config.ResolveConnectorDir(connectorDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao resolver diretório de conectores: %v\n", err)
			return 1
		}

		// Carregar .env do diretório de conectores
		loadRootEnv(filepath.Join(resolvedConnectorDir, ".env"))

		cfg, err := config.Load(resolvedConnectorDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao carregar configurações de %s: %v\n", resolvedConnectorDir, err)
			return 1
		}

		targetLogTarget := cfg.AbsLogsDir()
		if logFilePath != "" {
			if filepath.IsAbs(logFilePath) {
				targetLogTarget = logFilePath
			} else {
				abs, err := filepath.Abs(logFilePath)
				if err == nil {
					targetLogTarget = abs
				} else {
					targetLogTarget = logFilePath
				}
			}
		}

		log, err := logger.InitLogger(cfg.AppLogLevel, targetLogTarget)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
			return 1
		}
		defer log.Close()

		log.Infof("Iniciando OpenAPI MCP v%s...", Version)
		log.Infof("Diretório de conectores: %s", cfg.ConnectorDir)
		log.Infof("Caminho de conectores: %s", cfg.AbsConnectorsDir())

		connectorManager := connector.NewManager(cfg.AbsConnectorsDir())
		loadedConns, err := connectorManager.LoadAll()
		if err != nil {
			log.Errorf("Erro ao carregar conectores: %v", err)
			fmt.Fprintf(os.Stderr, "Erro ao carregar conectores: %v\n", err)
			return 1
		}
		connectors = loadedConns

		for id, manifest := range connectors {
			if manifest.Tools.Enabled != nil && !*manifest.Tools.Enabled {
				continue
			}

			specPath := manifest.Specification.Path
			if specPath == "" {
				continue
			}

			absSpecPath := specPath
			if !filepath.IsAbs(specPath) {
				absSpecPath = filepath.Join(cfg.AbsConnectorsDir(), id, specPath)
			}

			doc, err := openapi.LoadDocument(absSpecPath)
			if err != nil {
				log.Errorf("Erro ao carregar documento OpenAPI para conector %s: %v", id, err)
				fmt.Fprintf(os.Stderr, "Erro ao carregar OpenAPI do conector %s: %v\n", id, err)
				return 1
			}

			toolDefs, err := mapper.MapConnector(manifest, doc)
			if err != nil {
				log.Errorf("Erro ao mapear ferramentas do conector %s: %v", id, err)
				fmt.Fprintf(os.Stderr, "Erro ao mapear ferramentas do conector %s: %v\n", id, err)
				return 1
			}

			for _, td := range toolDefs {
				if err := registry.Register(td); err != nil {
					log.Errorf("Erro ao registrar ferramenta %s: %v", td.Name, err)
					fmt.Fprintf(os.Stderr, "Erro ao registrar ferramenta %s: %v\n", td.Name, err)
					return 1
				}
			}
		}
	}

	// Iniciar Servidor MCP
	srv := mcpserver.NewServer(Version, registry, connectors)
	ctx := context.Background()
	if err := srv.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Erro na execução do servidor MCP: %v\n", err)
		return 1
	}

	return 0
}

func fetchRemoteToTemp(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("falha ao baixar URL %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status HTTP não OK ao baixar %s: %s", rawURL, resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "openapi-mcp-remote-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func resolveSpecPath(connectorPath string, manifest *connector.Manifest) (string, error) {
	if manifest.Specification.URL != "" {
		if strings.HasPrefix(manifest.Specification.URL, "http://") || strings.HasPrefix(manifest.Specification.URL, "https://") {
			return fetchRemoteToTemp(manifest.Specification.URL)
		}
		return manifest.Specification.URL, nil
	}

	specPath := manifest.Specification.Path
	if specPath == "" {
		return "", fmt.Errorf("especificação (path ou url) não definida")
	}

	if strings.HasPrefix(connectorPath, "http://") || strings.HasPrefix(connectorPath, "https://") {
		base, err := url.Parse(connectorPath)
		if err != nil {
			return "", err
		}
		ref, err := url.Parse(specPath)
		if err != nil {
			return "", err
		}
		specURL := base.ResolveReference(ref).String()
		return fetchRemoteToTemp(specURL)
	}

	absConnectorPath, err := filepath.Abs(connectorPath)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(specPath) {
		return specPath, nil
	}
	return filepath.Join(filepath.Dir(absConnectorPath), specPath), nil
}

func loadRootEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

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
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
	}
}



func printHelp() {
	help := `OpenAPI MCP - Servidor MCP para especificações OpenAPI

Uso:
  openapi-mcp [opções]

Opções:
  -d, --connector-dir  Configura o diretório de conectores base.
  -c, --connector      Caminho local ou URL para um arquivo connector.yaml específico.
  -e, --env-file       Caminho para um arquivo .env personalizado.
  -l, --log-file       Caminho personalizado para o arquivo de log.
      --env            Define ou sobrescreve variáveis de ambiente (formato CHAVE=VALOR ou CHAVE1=VALOR1,CHAVE2=VALOR2).
      --version        Exibe a versão do programa.
  -h, --help           Exibe esta mensagem de ajuda.
`
	fmt.Print(help)
}
