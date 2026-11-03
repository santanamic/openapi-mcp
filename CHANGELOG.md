# Changelog

Todos os registros notáveis deste projeto serão documentados neste arquivo.

## [0.3.0-dev] - 2026-11-03

### Adicionado
- Integração com o SDK MCP (`github.com/modelcontextprotocol/go-sdk` v1.8.0).
- Implementação do pacote `internal/mcpserver` para transporte `stdio` e registro de ferramentas.
- Implementação do pacote `internal/tools` (Operation Mapper e ToolRegistry) com suporte a conversão para `snake_case`, derivação de nomes e detecção de colisões.
- Testes unitários para mapeamento de ferramentas, registro e servidor MCP via transporte in-memory.

## [0.2.0-dev] - 2026-11-03

### Adicionado
- Implementação do Connector Manager e validação de manifestos `connector.yaml`.
- Resolução de variáveis de ambiente (`${VAR}`) nos manifestos de conectores.
- Implementação do carregador e validador de documentos OpenAPI (`internal/openapi`) utilizando `kin-openapi`.
- Testes unitários para carregadores de OpenAPI e gerenciador de conectores.

## [0.1.0-dev] - 2026-11-03

### Adicionado
- Inicialização da estrutura do projeto (Go module).
- CLI inicial com suporte a comandos `init`.
- Gerenciamento de configuração (`config`) com resolução de diretório de dados e `.env`.
- Sistema de logs (`logger`) direcionado para `stderr`.
- Testes unitários para toda a fundação.
- Arquivo de Roadmap.
