# Visão geral

## Propósito

O OpenAPI MCP é um servidor MCP genérico que converte operações descritas em documentos OpenAPI em ferramentas MCP. Um conector informa qual especificação utilizar, qual servidor HTTP chamar, como autenticar e quais métodos autorizar. O Core não contém código específico de WooCommerce, Stripe, Shopee ou qualquer outro fornecedor.

## Modelo conceitual

- **Core:** coordena carregamento, interpretação, registro e execução.
- **Conector:** configuração declarativa que associa uma API a uma especificação, servidor, autenticação, permissões e opções de exposição MCP.
- **Especificação OpenAPI:** contrato das operações, parâmetros, corpos, schemas e respostas.
- **Cliente MCP:** descobre e invoca as ferramentas publicadas pelo servidor.

## Ciclo de vida

Na inicialização, o Core carrega a configuração geral, descobre os conectores, valida os manifestos, carrega as especificações, converte operações elegíveis em ferramentas e inicia o transporte MCP. Durante a execução, cada chamada é mapeada para uma operação, validada por políticas, transformada em HTTP e formatada como resultado MCP.

## Princípios

1. Conectores são dados, não código executável.
2. A descrição das operações vem do OpenAPI.
3. O Core não conhece endpoints de fornecedores.
4. Configuração e dados persistentes não ficam junto ao binário.
5. Comportamento não suportado deve ser explícito, nunca aproximado silenciosamente.
6. Segurança e permissões são aplicadas antes da execução HTTP.
7. O MVP deve ser pequeno e extensível, sem abstrações especulativas.

## Fora do propósito

O projeto não é um cliente API com interface gráfica, não substitui um gateway de API, não implementa lógica de negócio dos fornecedores e não executa scripts arbitrários fornecidos por conectores.
