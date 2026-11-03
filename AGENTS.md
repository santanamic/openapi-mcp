# PROMPT MESTRE — IMPLEMENTAÇÃO DO OPENAPI MCP

Atue como engenheiro de software sênior especializado em Go, arquitetura de sistemas, APIs REST, OpenAPI e Model Context Protocol (MCP).

Sua tarefa é implementar o projeto **OpenAPI MCP**, um servidor MCP genérico capaz de transformar operações descritas em especificações OpenAPI em ferramentas MCP executáveis.

## 1. Contexto e documentação

Antes de escrever qualquer código, examine integralmente a documentação existente no repositório.

A documentação técnica está organizada em:

- `README.md`
- `docs/overview.md`
- `docs/architecture.md`
- `docs/components.md`
- `docs/connectors.md`
- `docs/configuration.md`
- `docs/execution-flow.md`
- `docs/security.md`
- `docs/mvp.md`
- `docs/roadmap.md`

Também existem exemplos em `examples/`.

Esses documentos representam as decisões arquiteturais e os contratos definidos para o projeto.

**A documentação deverá ser tratada como fonte primária de referência para a implementação.**

Não substitua decisões já estabelecidas por alternativas arquiteturais sem identificar uma incompatibilidade técnica concreta.

Caso encontre inconsistências entre documentos, identifique-as e proponha uma solução antes de realizar alterações estruturais.

## 2. Objetivo da implementação

Construir um servidor MCP funcional que permita:

1. Carregar configurações gerais de um diretório de dados independente do executável.
2. Descobrir e carregar conectores declarativos em YAML.
3. Resolver referências a variáveis de ambiente.
4. Carregar especificações OpenAPI locais.
5. Interpretar operações OpenAPI 3.0.x.
6. Gerar automaticamente ferramentas MCP.
7. Registrar essas ferramentas no servidor MCP.
8. Receber chamadas dos clientes MCP.
9. Converter argumentos em requisições HTTP.
10. Aplicar autenticação e políticas de segurança.
11. Executar as requisições.
12. Converter respostas HTTP em resultados MCP.

O sistema deverá funcionar sem implementações específicas para APIs como WooCommerce, Stripe ou Shopee.

Toda integração deverá ser realizada por meio de conectores declarativos e especificações OpenAPI.

## 3. Tecnologias obrigatórias

Utilize:

- Go como linguagem principal.
- `github.com/modelcontextprotocol/go-sdk` para o protocolo MCP.
- `github.com/getkin/kin-openapi` para interpretação OpenAPI.
- `gopkg.in/yaml.v3` para leitura YAML.
- `net/http` para execução HTTP.
- Bibliotecas padrão de Go para JSON, arquivos, caminhos e configuração.

Antes de definir as dependências, confira as versões estáveis e seus requisitos de compatibilidade com a versão de Go escolhida.

Não utilize bibliotecas adicionais sem necessidade técnica justificada.

Utilize `map[string]any`, `json.RawMessage` e estruturas dinâmicas quando necessário para representar operações e schemas desconhecidos em tempo de compilação.

Não crie estruturas específicas para endpoints ou APIs individuais.

## 4. Arquitetura

Implemente a arquitetura documentada, respeitando a separação de responsabilidades.

Estrutura inicial de referência:

```text
openapi-mcp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── connector/
│   ├── openapi/
│   ├── tools/
│   ├── mcpserver/
│   ├── execution/
│   ├── auth/
│   └── security/
├── docs/
├── examples/
├── go.mod
└── README.md
```

Essa estrutura é uma referência, não uma obrigação de criar pacotes vazios ou fragmentar excessivamente a implementação.

Cada pacote deverá possuir uma responsabilidade real.

Evite abstrações especulativas, interfaces desnecessárias e camadas que apenas repassem chamadas.

Utilize interfaces quando houver uma necessidade concreta de desacoplamento, substituição de implementação ou testabilidade.

## 5. Configuração e distribuição

O projeto deverá ser distribuído por binários independentes.

O usuário final não deverá precisar instalar Go ou dependências de desenvolvimento.

Diretório de dados padrão:

- Windows: `%USERPROFILE%\.openapi-mcp`
- Linux/macOS: `~/.openapi-mcp`

Permitir personalização por:

```bash
openapi-mcp --data-dir '/caminho/personalizado'
```

Implementar também:

```bash
openapi-mcp init
```

O comando `init` deverá criar a estrutura inicial de dados e os arquivos de configuração necessários.

O executável não deverá depender do diretório de trabalho atual ou da localização do próprio binário.

O diretório de dados deverá permanecer intacto durante atualizações do executável.

## 6. Configuração geral

Utilize um arquivo `.env` localizado na raiz do diretório de dados.

Exemplo:

```ini
APP_ENV=production
APP_LOG_LEVEL=info
MCP_TRANSPORT=stdio
CONNECTORS_DIR=./connectors
LOGS_DIR=./logs
HTTP_TIMEOUT=30
```

As variáveis de ambiente do sistema operacional deverão prevalecer sobre os valores definidos no `.env`.

Os caminhos relativos de `CONNECTORS_DIR` e `LOGS_DIR` deverão ser resolvidos em relação ao diretório de dados.

Não utilizar o diretório atual como referência implícita.

## 7. Conectores

Implementar o carregamento automático dos diretórios de conectores.

Cada conector deverá possuir um `connector.yaml`.

Exemplo:

```yaml
id: woocommerce
name: WooCommerce
description: WooCommerce REST API
version: 1.0.0

specification:
  path: ./openapi.yaml

server:
  url: https://loja.com.br/wp-json/wc/v3

authentication:
  type: basic
  username: ${WOOCOMMERCE_CONSUMER_KEY}
  password: ${WOOCOMMERCE_CONSUMER_SECRET}

permissions:
  methods:
    - GET
    - POST
    - PUT
    - PATCH
    - DELETE

tools:
  enabled: true
  prefix: woocommerce
```

Respeitar integralmente o contrato definido em `docs/connectors.md`.

Implementar:

- Descoberta dos conectores.
- Leitura dos manifestos.
- Validação dos campos.
- Validação dos identificadores.
- Aplicação de valores padrão.
- Resolução de variáveis de ambiente.
- Validação de referências ausentes.
- Normalização das configurações.

Os valores secretos deverão permanecer apenas em memória durante a execução.

Não registrar credenciais em logs ou mensagens de erro.

## 8. Interpretação OpenAPI

Utilizar exclusivamente `kin-openapi` para carregar e interpretar as especificações.

Priorizar OpenAPI 3.0.x.

Implementar:

- Identificação das operações.
- Paths e métodos HTTP.
- `operationId`.
- Descrições e resumos.
- Path parameters.
- Query parameters.
- Header parameters.
- Request bodies.
- Responses.
- Schemas.
- Referências internas `$ref`.
- Servidores definidos na especificação.

Quando o conector informar `server.url`, ele deverá prevalecer sobre os servidores da especificação.

Caso contrário, utilizar o primeiro servidor aplicável, respeitando os valores padrão de suas variáveis.

Não implementar um parser OpenAPI próprio.

Operações que dependam de recursos ainda não suportados deverão ser diagnosticadas e ignoradas, sem gerar ferramentas que executem requisições incorretas.

## 9. Geração de ferramentas MCP

Converter dinamicamente as operações OpenAPI em ferramentas MCP.

Regras para nomes:

1. Utilizar `operationId` quando disponível.
2. Na ausência de `operationId`, derivar o nome a partir do método e path.
3. Normalizar o nome para `snake_case`.
4. Aplicar o prefixo do conector.
5. Garantir unicidade global.

Exemplos:

```text
woocommerce_get_products
woocommerce_get_product
woocommerce_create_product
woocommerce_update_product
woocommerce_delete_product
```

Cada ferramenta deverá possuir:

- Nome.
- Descrição.
- JSON Schema de entrada.
- Referência à operação OpenAPI original.

Os schemas de entrada deverão ser derivados dos parâmetros e request bodies da operação.

Não criar ferramentas manualmente para APIs específicas.

O `ToolRegistry` deverá manter o relacionamento entre os nomes das ferramentas e suas operações.

O servidor MCP deverá utilizar o SDK Go para registrar e expor as ferramentas.

## 10. Execução HTTP

Implementar um fluxo genérico de execução.

O executor deverá:

1. Localizar a operação associada à ferramenta.
2. Validar os argumentos recebidos.
3. Verificar as permissões.
4. Construir a URL.
5. Substituir parâmetros de path.
6. Adicionar query parameters.
7. Aplicar headers.
8. Preparar o request body JSON.
9. Aplicar autenticação.
10. Validar o destino.
11. Executar a requisição HTTP.
12. Interpretar a resposta.
13. Retornar o resultado MCP.

Utilizar `net/http`.

O timeout deverá ser configurável globalmente.

Implementar tratamento de:

- Erros de transporte.
- Respostas HTTP de erro.
- Respostas JSON.
- Respostas textuais.
- Respostas vazias.
- Respostas excessivamente grandes.

Não registrar corpos de requisição ou respostas por padrão.

## 11. Autenticação

Implementar inicialmente:

- `none`
- `basic`
- `api_key`
- `bearer`

A autenticação deverá ser genérica e configurada pelo manifesto.

Não implementar autenticação específica para fornecedores.

OAuth2 e HMAC ficam fora do MVP, salvo necessidade técnica comprovada.

## 12. Segurança

A implementação deverá respeitar `docs/security.md`.

Prioridades:

- Prevenção contra SSRF.
- Validação de URLs.
- Bloqueio de IPs locais e privados por padrão.
- Validação de DNS e redirecionamentos.
- Proteção de credenciais.
- Timeout HTTP.
- Limites de tamanho de respostas.
- Validação de argumentos.
- Permissões por método HTTP.

Na ausência de `permissions.methods`, permitir somente `GET`.

Métodos de escrita deverão exigir autorização explícita.

Não encaminhar credenciais para outros hosts por redirecionamento.

O transporte `stdio` deverá manter stdout exclusivamente para mensagens MCP. Logs e diagnósticos deverão utilizar stderr ou arquivos.

## 13. Estratégia de implementação

Implemente incrementalmente, mas mantenha o foco em entregar um MVP funcional.

### Etapa 1 — Fundação

- Inicializar o módulo Go.
- Configurar dependências.
- Implementar CLI básica.
- Implementar resolução do diretório de dados.
- Implementar comando `init`.
- Implementar carregamento do `.env`.
- Configurar logging.

Critério de conclusão: o executável inicia corretamente e prepara o ambiente.

### Etapa 2 — Conectores e OpenAPI

- Implementar Connector Manager.
- Validar manifestos.
- Resolver variáveis de ambiente.
- Carregar documentos OpenAPI locais.
- Validar especificações.
- Identificar operações suportadas.

Critério de conclusão: o sistema consegue carregar o exemplo WooCommerce e identificar suas operações.

### Etapa 3 — Ferramentas MCP

- Implementar Operation Mapper.
- Gerar schemas de entrada.
- Implementar Tool Registry.
- Integrar o SDK MCP.
- Registrar ferramentas dinamicamente.
- Implementar transporte stdio.

Critério de conclusão: um cliente MCP consegue listar as ferramentas geradas pelo conector.

### Etapa 4 — Execução HTTP

- Implementar políticas.
- Implementar Request Builder.
- Implementar autenticação.
- Implementar validação de destino.
- Implementar execução HTTP.
- Implementar Response Formatter.

Critério de conclusão: um cliente MCP consegue executar operações reais da API utilizando as ferramentas geradas.

### Etapa 5 — Testes e distribuição

- Criar testes unitários.
- Criar testes de integração.
- Testar diferentes métodos HTTP.
- Testar autenticação.
- Testar erros HTTP.
- Testar permissões.
- Testar carregamento de múltiplos conectores.
- Testar colisões de nomes.
- Compilar para Windows, Linux e macOS.

Critério de conclusão: MVP funcional, testado e compilável para as plataformas previstas.

## 14. Testes

Priorizar testes automatizados.

Utilizar `httptest` para simular APIs REST.

Não depender de credenciais ou serviços externos para os testes automatizados.

Cobrir especialmente:

- Resolução de configurações.
- Validação dos manifestos.
- Geração de ferramentas.
- Construção de URLs.
- Serialização de parâmetros.
- Autenticação.
- Políticas HTTP.
- Segurança dos destinos.
- Formatação das respostas.

Os testes deverão ser independentes e reproduzíveis.

## 15. Qualidade e manutenção

Seguir as convenções idiomáticas de Go.

Requisitos:

- Código simples e legível.
- Tratamento explícito de erros.
- Context propagation nas operações HTTP.
- Concorrência segura.
- Evitar estado global mutável.
- Evitar duplicação.
- Não criar abstrações sem necessidade.
- Não ignorar erros.
- Não utilizar dependências desnecessárias.

Executar `gofmt`, `go vet` e `go test ./...` durante o desenvolvimento.

## 16. Documentação

Manter a documentação atualizada conforme a implementação evoluir.

Atualizar o README quando houver mudanças em comandos, configuração, instalação ou execução.

Documentar decisões técnicas relevantes.

Não duplicar conteúdo entre documentos.

Criar exemplos funcionais que correspondam aos contratos efetivamente implementados.

Não documentar funcionalidades como disponíveis antes de estarem implementadas.

## 17. Regras de trabalho

- Não começar reescrevendo toda a documentação.
- Não alterar decisões arquiteturais sem justificativa técnica.
- Não antecipar funcionalidades do roadmap.
- Não criar implementações específicas para APIs.
- Não implementar upload, download avançado, OAuth2 ou paginação automática no MVP.
- Não criar interface gráfica.
- Não introduzir hot reload.
- Não criar sistema de plugins executáveis.
- Não fragmentar o projeto em pacotes desnecessários.
- Não substituir bibliotecas previstas sem motivo concreto.

Se encontrar uma dificuldade técnica, analise primeiro como resolvê-la dentro da arquitetura estabelecida.

Se uma decisão importante precisar ser revista, explique o problema, apresente a alternativa e aguarde aprovação antes de modificar a arquitetura.

## 18. Primeira tarefa

Comece examinando o repositório e toda a documentação existente.

Em seguida:

1. Apresente um diagnóstico breve da estrutura atual.
2. Identifique eventuais inconsistências técnicas.
3. Confirme as dependências e suas compatibilidades.
4. Apresente o plano de implementação da primeira etapa.
5. Inicie a implementação da fundação do projeto.

Não se limite a apresentar sugestões ou pseudocódigo. Trabalhe diretamente nos arquivos do projeto.

Ao concluir a primeira etapa, apresente:

- Arquivos criados e modificados.
- Decisões técnicas tomadas.
- Funcionalidades implementadas.
- Testes executados e resultados.
- Pendências para a próxima etapa.

**Objetivo final:** construir um servidor MCP genérico, confiável e extensível, mantendo o Core independente das APIs integradas e os conectores completamente declarativos.