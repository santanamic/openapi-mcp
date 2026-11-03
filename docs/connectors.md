# Contrato dos conectores

## Estrutura de diretórios

Cada conector tem seu próprio diretório dentro de `CONNECTORS_DIR`:

```text
connectors/
├── woocommerce/
│   ├── connector.yaml
│   └── openapi.yaml
└── stripe/
    └── connector.yaml
```

O arquivo obrigatório é `connector.yaml`. O documento OpenAPI pode ser local ou remoto. Não há código executável no conector.

## Manifesto de referência

```yaml
id: woocommerce
name: WooCommerce
description: WooCommerce REST API
version: 1.0.0

specification:
  path: ./openapi.yaml

server:
  url: https://loja.example/wp-json/wc/v3

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

## Campos

### Identificação

- `id` (obrigatório): identificador único entre conectores. Regex recomendada: `^[a-z0-9][a-z0-9-]*$`. É usado internamente e no prefixo padrão.
- `name` (obrigatório): nome de exibição.
- `description` (opcional): descrição do conector.
- `version` (opcional): versão da configuração/conector, não necessariamente versão da API. Recomenda-se SemVer.

### `specification` (obrigatório)

Deve conter exatamente uma propriedade:

- `path`: caminho para arquivo OpenAPI local.
- `url`: URL HTTP(S) de documento OpenAPI remoto.

`path` relativo é resolvido em relação ao diretório que contém o `connector.yaml`, nunca em relação ao working directory. Caminho absoluto é aceito. A URL é usada somente para obter a especificação; não é o endpoint base das operações.

### `server` (opcional)

- `url`: URL base efetiva para as operações. Quando definida, substitui os `servers` do OpenAPI.
- Quando omitida, selecionar o primeiro `servers[].url` aplicável no documento, respeitando defaults de variáveis de servidor.
- A URL deve ser absoluta, HTTP ou HTTPS, sem fragmento. Para produção, HTTPS é recomendado.

### `authentication` (opcional)

- `type`: `none`, `basic`, `api_key` ou `bearer` no MVP.
- `username` e `password`: Basic.
- `value`: segredo para Bearer ou API key.
- `name`: nome do header ou query parameter para API key.
- `in`: `header` ou `query` para API key.

Exemplos:

```yaml
authentication:
  type: bearer
  value: ${SERVICE_TOKEN}
```

```yaml
authentication:
  type: api_key
  name: X-API-Key
  in: header
  value: ${SERVICE_API_KEY}
```

```yaml
authentication:
  type: api_key
  name: api_key
  in: query
  value: ${SERVICE_API_KEY}
```

Se `authentication` for omitido, o padrão é `type: none`. A autenticação configurada no conector é aplicada genericamente às operações desse conector. A especificação OpenAPI continua sendo a fonte dos contratos e pode documentar `securitySchemes`; no MVP, o Core não implementa negociação OAuth2 nem combinações complexas de múltiplos esquemas. Se uma operação exige um esquema não atendido pela configuração do conector, ela deve ser diagnosticada e não executada sem autenticação adequada.

### `permissions` (opcional)

- `methods`: lista de métodos HTTP autorizados, normalizados para maiúsculas.
- Se omitido, o padrão é somente `GET`.
- Valores válidos no MVP: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS`.
- Método fora da lista é bloqueado antes da criação/envio da requisição.

### `tools` (opcional)

- `enabled`: booleano; padrão `true`. Se `false`, o conector é carregado/validado, mas não publica ferramentas.
- `prefix`: prefixo dos nomes MCP; padrão igual a `id`.

## Regras de validação

- `id` deve ser válido e único.
- `specification` deve conter exatamente `path` ou `url`.
- `server.url`, quando definido, deve ser URL absoluta HTTP(S).
- `authentication.type` deve ser conhecido; campos obrigatórios do tipo devem existir.
- `permissions.methods` não pode conter valores desconhecidos.
- `tools.enabled` deve ser booleano e `tools.prefix` deve resultar em nomes válidos.
- Variáveis de ambiente obrigatórias devem estar resolvidas antes de iniciar o servidor.
- IDs duplicados, manifestos inválidos e configurações de autenticação incompletas devem gerar diagnóstico com o caminho do arquivo.

## Descoberta e falhas

O Connector Manager percorre diretórios imediatos de `CONNECTORS_DIR` e procura `connector.yaml` em cada diretório. Diretórios sem manifesto são ignorados com diagnóstico informativo. Um manifesto inválido desabilita aquele conector; não deve impedir necessariamente os demais de carregar. Falha global de configuração, como diretório inexistente ou inválido, impede inicialização.

## Evolução

OAuth2, autenticação HMAC, parâmetros adicionais de autenticação e políticas por operação podem ser adicionados ao contrato em versões futuras, mantendo compatibilidade e sem permitir scripts no YAML.
