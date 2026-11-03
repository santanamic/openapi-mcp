# OpenAPI MCP

Servidor MCP genérico que transforma operações descritas em documentos OpenAPI em ferramentas MCP, executando as chamadas por HTTP.

> **Status:** planejamento e documentação. Este repositório ainda não contém implementação.

## Objetivo

Permitir integrar APIs REST compatíveis com OpenAPI por meio de conectores declarativos, sem implementar código específico para cada API. O Core interpreta os documentos, gera ferramentas e executa as operações autorizadas.

## Requisitos de execução

- Binário compilado para o sistema operacional e arquitetura do usuário.
- Diretório de dados independente do local do executável.
- Arquivo `.env` para configurações gerais e valores de ambiente referenciados pelos conectores.
- Conectores YAML e documentos OpenAPI.
- Transporte MCP inicial: `stdio`.

O usuário final não precisa instalar Go ou ferramentas de desenvolvimento.

## Instalação e inicialização

Distribuições futuras fornecerão um binário por plataforma:

- Windows: `openapi-mcp.exe`
- Linux: `openapi-mcp`
- macOS: `openapi-mcp`

Inicializar a estrutura de dados:

```bash
openapi-mcp init
```

Usar um diretório personalizado:

```bash
openapi-mcp init --data-dir "/caminho/openapi-mcp-data"
openapi-mcp --data-dir "/caminho/openapi-mcp-data"
```

Iniciar o servidor com o diretório padrão:

```bash
openapi-mcp
```

O diretório padrão será `~/.openapi-mcp` no Linux/macOS e `%USERPROFILE%\.openapi-mcp` no Windows. `--data-dir` prevalece sobre o padrão. O processo não depende do diretório de trabalho atual nem da localização do binário.

## Estrutura de dados

```text
.openapi-mcp/
├── .env
├── connectors/
│   ├── woocommerce/
│   │   ├── connector.yaml
│   │   └── openapi.yaml
│   └── outro-servico/
│       └── connector.yaml
└── logs/
```

Cada diretório de conector contém um `connector.yaml`. A especificação pode ser local ou remota. Caminhos locais relativos são resolvidos em relação ao diretório do conector.

## Configuração inicial

Exemplo de `.env`:

```ini
APP_ENV=production
APP_LOG_LEVEL=info
MCP_TRANSPORT=stdio
CONNECTORS_DIR=./connectors
LOGS_DIR=./logs
HTTP_TIMEOUT=30
WOOCOMMERCE_CONSUMER_KEY=ck_example
WOOCOMMERCE_CONSUMER_SECRET=cs_example
```

Variáveis definidas no ambiente do sistema têm precedência sobre o `.env`. Não versione valores secretos. O manifesto referencia os nomes das variáveis, não seus valores.

## Conectar a um cliente MCP

Exemplo conceitual de configuração para um cliente que aceita servidores MCP locais por `stdio`:

```json
{
  "mcpServers": {
    "openapi-mcp": {
      "command": "/caminho/absoluto/openapi-mcp",
      "args": ["--data-dir", "/caminho/absoluto/.openapi-mcp"]
    }
  }
}
```

No Windows, use o caminho absoluto para `openapi-mcp.exe`. O formato exato do arquivo de configuração depende do cliente MCP utilizado.

## Documentação

- [Visão geral](docs/overview.md)
- [Arquitetura](docs/architecture.md)
- [Componentes](docs/components.md)
- [Contrato dos conectores](docs/connectors.md)
- [Configuração](docs/configuration.md)
- [Fluxo de execução](docs/execution-flow.md)
- [Segurança](docs/security.md)
- [MVP](docs/mvp.md)
- [Roadmap](docs/roadmap.md)
- [Exemplo WooCommerce](examples/woocommerce/connector.yaml)
- [Exemplo de cliente MCP](examples/client/mcp-config.json)

## Tecnologias de referência

| Finalidade | Dependência |
|---|---|
| Protocolo MCP | `github.com/modelcontextprotocol/go-sdk` |
| OpenAPI | `github.com/getkin/kin-openapi` |
| YAML | `gopkg.in/yaml.v3` |
| HTTP | `net/http` |
| JSON | `encoding/json` |

As versões serão fixadas na implementação. O planejamento de dependências está em [architecture.md](docs/architecture.md).

## Escopo atual

O MVP prioriza OpenAPI 3.0.x, especificações locais, ferramentas MCP dinâmicas, path/query/header parameters, corpos JSON, HTTP, autenticação básica/API key/Bearer e transporte `stdio`. Upload/download avançado, OAuth2, paginação e transformações específicas ficam fora do MVP.
