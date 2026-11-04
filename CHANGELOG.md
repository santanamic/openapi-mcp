# Changelog

Todos os registros notáveis deste projeto serão documentados neste arquivo.

## [1.2.0] - 2026-11-03

### Adicionado
- **Comando `openapi-mcp validate`**: Subcomando CLI para verificação diagnóstica pré-execução de manifestos (`connector.yaml`), validação de documentos OpenAPI via `kin-openapi`, teste de conectividade HTTP à URL base do servidor e verificação de variáveis de ambiente/credenciais.

## [1.1.0] - 2026-11-03

### Adicionado
- **Comando `openapi-mcp auth`**: Subcomando CLI para fluxo OAuth2 interativo (Authorization Code Flow) utilizando servidor HTTP local temporário, abertura automática do navegador e formatação de tokens (`access_token` e `refresh_token`).
- **Renovação Automática de Tokens OAuth2**: Suporte a `refresh_token` e renovação transparente e automática de tokens de acesso expirados utilizando `golang.org/x/oauth2` TokenSource.
- **Tratamento de Expiração de Tokens (401 Unauthorized)**: Respostas HTTP `401 Unauthorized` agora retornam mensagens de erro acionáveis orientando humanos e LLMs a reexecutarem `openapi-mcp auth`.
- Suporte a autenticação `oauth2` no pacote `internal/auth`.

## [1.0.0] - 2026-11-03

### Adicionado
- Testes de integração End-to-End (E2E) cobrindo o fluxo completo do servidor MCP à execução HTTP.
- Scripts de compilação multiplataforma (`scripts/build.sh` e `scripts/build.ps1`).
- Documentação final de instalação, inicialização e uso.
- Lançamento oficial do MVP funcional (v1.0.0).

## [0.4.0-dev] - 2026-11-03

### Adicionado
- Implementação do pacote `internal/security` para validação de URLs e proteção contra SSRF (bloqueio de loopback, IPs privados, esquemas não HTTP e credenciais embutidas).
- Implementação do pacote `internal/auth` com suporte a autenticação genérica (Basic, Bearer e API Key via header ou query).
- Implementação do pacote `internal/execution` incluindo política de permissões de métodos, construtor de requisições, executor HTTP e formatador de respostas com limite de tamanho (1MB) e tratamento de erros.
- Integração da execução de ferramentas no `internal/mcpserver` ligando chamadas MCP a requisições HTTP reais.
- Testes unitários para segurança, autenticação e execução.

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
