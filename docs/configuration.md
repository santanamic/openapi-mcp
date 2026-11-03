# Configuração do servidor

## Diretório de dados

Padrão:

- Windows: `%USERPROFILE%\\.openapi-mcp`
- Linux/macOS: `~/.openapi-mcp`

Personalização:

```bash
openapi-mcp --data-dir "/caminho/personalizado"
```

A opção `--data-dir` tem precedência sobre o diretório padrão. O caminho deve ser normalizado e convertido em absoluto durante a inicialização. O sistema nunca deve procurar `.env` ou conectores no diretório atual como fallback implícito.

## Estrutura

```text
.openapi-mcp/
├── .env
├── connectors/
│   └── <connector-id>/
│       ├── connector.yaml
│       └── openapi.yaml
└── logs/
```

`CONNECTORS_DIR` e `LOGS_DIR` podem ser caminhos absolutos ou relativos ao diretório de dados.

## Arquivo `.env`

Exemplo:

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

### Variáveis gerais

| Variável | Padrão | Descrição |
|---|---|---|
| `APP_ENV` | `production` | Ambiente/identificação operacional |
| `APP_LOG_LEVEL` | `info` | `debug`, `info`, `warn` ou `error` |
| `MCP_TRANSPORT` | `stdio` | Transporte inicial; no MVP somente `stdio` |
| `CONNECTORS_DIR` | `./connectors` | Diretório de conectores, relativo ao data-dir |
| `LOGS_DIR` | `./logs` | Diretório de logs, relativo ao data-dir |
| `HTTP_TIMEOUT` | `30` | Timeout total por chamada HTTP, em segundos |

## Precedência

1. Variáveis presentes no ambiente do processo.
2. Valores do arquivo `.env` no diretório de dados.
3. Defaults documentados.

O `.env` não sobrescreve uma variável já definida no ambiente do sistema. O carregamento do arquivo deve ser explícito a partir do data-dir.

## Expansão de variáveis nos conectores

O Connector Manager resolve placeholders `${VAR}` nos campos de autenticação do manifesto.

Regras:

- O nome deve corresponder a uma variável de ambiente válida.
- A expansão ocorre após carregar `.env` e aplicar a precedência do ambiente do processo.
- Variável inexistente ou vazia em campo obrigatório gera erro de configuração.
- Não suportar sintaxe de fallback, interpolação aninhada ou execução de comandos no MVP.
- Os valores resolvidos permanecem somente em memória.
- Não registrar valores expandidos em logs, mensagens de erro, dumps de configuração ou diagnósticos.

O `.env` pode conter credenciais para conveniência local, mas deve ter permissões restritas e ser excluído do controle de versão. Em ambientes gerenciados, as mesmas variáveis podem ser injetadas pelo sistema operacional.

## Inicialização da estrutura

`openapi-mcp init` cria o diretório de dados, `.env` de exemplo, diretórios `connectors/` e `logs/`. Não deve substituir arquivos existentes. Se o diretório já existir, completar somente itens ausentes e informar o que foi preservado.

## Logs

Logs são direcionados para `LOGS_DIR`; mensagens operacionais essenciais também podem ir para stderr. Em `stdio`, stdout é reservado exclusivamente para mensagens do protocolo MCP. Não emitir banners, logs ou progresso em stdout.
