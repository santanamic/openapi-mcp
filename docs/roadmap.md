# Roadmap

## Visão Geral
Este documento define a evolução do OpenAPI MCP.

## Etapas Realizadas

### Etapa 1 — Fundação
- Inicialização do módulo Go.
- Configuração de dependências (`go-sdk`, `kin-openapi`, `yaml.v3`).
- Implementação da CLI básica (`init`, `data-dir`, `--version`, `--help`).
- Resolução do diretório de dados (independente do binário).
- Carregamento de `.env` com precedência do ambiente do sistema.
- Sistema de logs direcionado a `stderr`.
- Testes unitários para configuração, logs e CLI.

### Versão 0.2.0 — Conectores e OpenAPI
- Implementação do Connector Manager.
- Validação de manifestos `connector.yaml`.
- Resolução de variáveis de ambiente (`${VAR}`).
- Carregamento e validação de documentos OpenAPI locais (`kin-openapi`).
- Identificação de operações suportadas.

### Versão 0.3.0 — Ferramentas MCP
- Mapeamento de operações OpenAPI para ferramentas MCP.
- Registro de ferramentas (`ToolRegistry`).
- Integração com o SDK MCP para exposição de ferramentas.
- Implementação do transporte `stdio`.

### Versão 0.4.0 — Execução e Segurança
- Implementação do pacote `internal/security` para validação de URLs e proteção SSRF.
- Implementação do pacote `internal/auth` com suporte a Basic, Bearer e API Key.
- Implementação do pacote `internal/execution` (Policy, Request Builder, HTTP Client, Response Formatter).
- Integração da execução HTTP nas chamadas de ferramentas MCP.
- Testes unitários para segurança, autenticação e execução.

### Versão 1.0.0 — MVP Funcional
- Testes de integração end-to-end.
- Cobertura de diferentes métodos HTTP e exemplos práticos.
- Compilação multiplataforma.
- Documentação final de instalação e uso.

## Próximas Etapas (Pós-MVP)

### Versão 1.1.0 — Ferramentas CLI Avançadas (`auth`)

- **Comando `openapi-mcp auth` (OAuth2 Interativo):**
  - Subcomando CLI para fluxo OAuth2 (Authorization Code Flow) utilizando servidor HTTP local temporário.
  - Abertura automática do navegador padrão do usuário.
  - Página de callback e saída formatada no terminal destacando os tokens obtidos (`access_token` e `refresh_token`) para cópia e inserção manual segura pelo usuário.

- **Renovação Automática de Tokens & Gestão de Expiração:**
  - Suporte a `refresh_token`, `client_id`, `client_secret` e `token_url` nos manifestos de conectores.
  - Renovação transparente e automática de tokens de acesso expirados utilizando `golang.org/x/oauth2` TokenSource.
  - Em caso de expiração irrecuperável detectada por erro `401 Unauthorized`, as ferramentas MCP retornam mensagens de erro acionáveis orientando humanos e LLMs a reexecutarem o comando `openapi-mcp auth`.

### Versão 1.2.0 — Comando de Validação Pré-Execução (validate):

- **Comando `openapi-mcp validate` (Validação Pré-Execução):**
  - Subcomando CLI para verificação diagnóstica completa antes de iniciar o servidor MCP.
  - **Validação de Manifesto:** Checagem sintática e estrutural do `connector.yaml` (campos obrigatórios, caminhos e permissões).
  - **Validação OpenAPI:** Verificação do documento OpenAPI utilizando o validador `kin-openapi`.
  - **Checagem de Conectividade:** Teste de alcance HTTP (`ping`/HEAD/GET) à URL base do servidor (`server.url`).
  - **Validação de Credenciais:** Verificação da presença de variáveis de ambiente obrigatórias e teste opcional de autenticação contra a API de destino.

## Conclusão do MVP
A versão 1.0.0 do OpenAPI MCP está completa, testada e pronta para uso em produção com o transporte `stdio` e conectores OpenAPI declarativos.
