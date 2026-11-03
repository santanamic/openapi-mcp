# Componentes

Este documento define responsabilidades e fronteiras. Os nomes descrevem responsabilidades lógicas, não obrigatoriamente pacotes ou tipos Go independentes.

| Componente lógico | Responsabilidade | Não deve fazer |
|---|---|---|
| App/Core | Compor dependências, iniciar e encerrar o processo | Regras de APIs específicas |
| Configuration Manager | Resolver diretório de dados, `.env`, defaults e ambiente | Ler manifestos de conectores |
| Connector Manager | Descobrir, desserializar, validar e normalizar conectores | Interpretar paths OpenAPI |
| Specification Loader | Obter documento local/remoto com limites e validação de destino | Mapear ferramentas |
| OpenAPI Parser | Usar `kin-openapi` para carregar/validar e expor modelo | Executar HTTP |
| Operation Mapper | Converter operações suportadas em ToolDefinition | Registrar no SDK ou fazer requests |
| Tool Registry | Indexar ferramentas e operações, garantir nomes únicos | Implementar protocolo MCP |
| MCP Server | Adaptar definições/chamadas ao SDK MCP e transporte stdio | Construir URL HTTP |
| Operation Policy | Autorizar método e aplicar regras do conector | Aplicar credenciais |
| Request Builder | Resolver parâmetros, URL, headers e corpo | Enviar requisição |
| Authentication | Aplicar `none`, Basic, API key ou Bearer | Descobrir credenciais por conta própria |
| Host Validator | Validar esquema, host, IP, DNS e redirecionamentos | Decidir permissões de negócio |
| HTTP Client | Executar a requisição com timeout e limites | Conhecer OpenAPI ou MCP |
| Response Formatter | Converter status, headers e corpo em resultado MCP | Alterar semântica da resposta da API |

## Contratos internos principais

### ConnectorConfig normalizado

Contém `id`, metadados, localização da especificação, URL base efetiva, configuração de autenticação resolvida, métodos permitidos e opções MCP. Valores secretos existem apenas em memória e não devem ser serializados em logs.

### OperationDefinition

Representa uma operação independente de fornecedor: conector, método, path, `operationId`, resumo/descrição, parâmetros, request body, respostas e JSON Schema de entrada. Deve manter referência ao objeto OpenAPI original ou a metadados suficientes para construir a requisição sem reanalisar o YAML a cada chamada.

### ToolDefinition

Contém nome MCP, descrição, schema JSON de entrada e referência estável à `OperationDefinition`. O nome é único em todo o servidor.

### ToolCall

Contém nome de ferramenta e argumentos JSON recebidos do SDK. Os argumentos devem ser validados contra o schema derivado antes da execução. Valores dinâmicos podem ser representados por `map[string]any` e `json.RawMessage`.

### ExecutionResult

Contém status HTTP, headers permitidos, corpo interpretado ou texto, indicação de erro HTTP e metadados de truncamento. Não deve incluir segredos nem dados internos desnecessários.

## Sequência de responsabilidades

1. Connector Manager entrega `ConnectorConfig` validado.
2. Specification Loader entrega bytes/documento OpenAPI obtido.
3. OpenAPI Parser entrega documento validado.
4. Operation Mapper entrega `ToolDefinition` e `OperationDefinition`.
5. Tool Registry indexa as definições.
6. MCP Server registra as ferramentas e encaminha chamadas.
7. Operation Executor resolve a operação e produz `ExecutionResult`.
8. MCP Server converte o resultado para o formato de retorno do SDK.
