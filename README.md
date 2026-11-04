# OpenAPI MCP

Servidor MCP genérico que transforma operações descritas em documentos OpenAPI em ferramentas MCP, executando as chamadas por HTTP.

> **Status:** Versão 1.0.0 (MVP Funcional) lançada e totalmente operacional.

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

Usar um diretório personalizado:

```bash
openapi-mcp --connector-dir "/caminho/openapi-mcp-data"
```

Iniciar o servidor com o diretório padrão:

```bash
openapi-mcp
```

O diretório padrão será `~/.openapi-mcp` no Linux/macOS e `%USERPROFILE%\.openapi-mcp` no Windows. `--connector-dir` prevalece sobre o padrão. O processo não depende do diretório de trabalho atual nem da localização do binário.

## Execução Direta e Variáveis de Ambiente (v1.0.0+)

Além do diretório de conectores padrão, o OpenAPI MCP suporta a execução direta de um conector específico (local ou via URL remota) e a injeção flexível de variáveis de ambiente.

> **Nota:** Recursos como `--env-file` e `--env` devem ser utilizados apenas em **casos específicos** (como testes locais, depuração ou quando o cliente de IA utilizado não oferece suporte nativo à injeção de variáveis de ambiente). Em uso normal com clientes compatíveis, recomenda-se configurar as variáveis no arquivo `.env` padrão do diretório de conectores (`--connector-dir`) ou diretamente na configuração do cliente MCP.

### 1. Executar um conector específico (`--connector` / `-c`)
Você pode apontar diretamente para um arquivo `connector.yaml` local ou URL remota:

```bash
# Conector local
openapi-mcp --connector "./examples/openapi-mcp/connectors/viacep/connector.yaml"

# Conector remoto (URL)
openapi-mcp --connector "https://exemplo.com/path/to/connector.yaml"
```

### 2. Carregar arquivo `.env` personalizado (`--env-file` / `-e`)
```bash
openapi-mcp -c connector.yaml --env-file "./secrets.env"
```

### 3. Injetar variáveis via CLI (`--env`)
- **Definir variáveis de ambiente:**
  ```bash
  openapi-mcp -c connector.yaml --env "WOOCOMMERCE_CONSUMER_KEY=ck_123,KEY2=val2"
  ```

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
      "args": ["--connector-dir", "/caminho/absoluto/.openapi-mcp"]
    }
  }
}
```

No Windows, use o caminho absoluto para `openapi-mcp.exe`. O formato exato do arquivo de configuração depende do cliente MCP utilizado.

## Autenticação OAuth2 Interativa (`openapi-oauth2`) & Renovação Automática (v1.1.0+)

O utilitário independente **`openapi-oauth2`** fornece um fluxo OAuth2 Authorization Code Flow interativo completo:

1. **Execução do utilitário:**
   ```bash
   openapi-oauth2 --client-id "seu-client-id" --auth-url "https://auth.exemplo.com/authorize" --token-url "https://auth.exemplo.com/token" [--client-secret "seu-secret"]
   ```
2. O utilitário inicia um servidor HTTP local de callback (`http://localhost:8080/callback`) e abre automaticamente o navegador padrão.
3. Após a autorização, exibe no terminal os tokens obtidos (`OAUTH2_ACCESS_TOKEN` e `OAUTH2_REFRESH_TOKEN`) para cópia ao arquivo `.env`.

### Renovação Automática de Tokens (`TokenSource`)
Se o conector estiver configurado com `refresh_token`, `client_id`, `client_secret` e `token_url`, o servidor MCP gerencia a expiração de forma 100% transparente:
- Utiliza o `access_token` normalmente para alta performance.
- Quando o token expira, o Core utiliza automaticamente o `refresh_token` para obter um novo `access_token` nos bastidores antes de despachar a chamada HTTP.
- Caso ocorra um erro `401 Unauthorized` irrecuperável, a ferramenta MCP retorna uma mensagem de erro acionável orientando a reexecução de `openapi-oauth2`.

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
