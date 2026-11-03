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
- Implementação do fluxo de execução HTTP.
- Autenticação genérica (Basic, API Key, Bearer).
- Políticas de segurança (validação de destinos, SSRF, limites).
- Formatação de respostas.

### Versão 1.0.0 — MVP Funcional
- Testes de integração.
- Cobertura de diferentes métodos HTTP.
- Compilação multiplataforma.
- Documentação final de instalação e uso.
