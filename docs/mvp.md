# Escopo do MVP

## Objetivo verificável

Dado um conector válido e uma especificação OpenAPI 3.0.x local, iniciar o binário e disponibilizar as operações HTTP elegíveis como ferramentas MCP. Uma chamada MCP deve produzir uma requisição HTTP correta e devolver o resultado de forma compreensível.

## Incluído

1. Binário independente por sistema operacional/arquitetura suportados.
2. CLI mínima: inicialização do servidor, `init` e `--data-dir`.
3. Diretório de dados independente do binário e do working directory.
4. `.env` com precedência do ambiente do processo.
5. Descoberta e validação de conectores YAML.
6. Especificações OpenAPI locais.
7. OpenAPI 3.0.x.
8. `$ref` internos, com proteção contra referências inválidas/ciclos.
9. Paths, métodos e `operationId`.
10. Parâmetros path, query e header nos estilos implementados.
11. Request body JSON.
12. Geração de schemas MCP e ferramentas dinâmicas.
13. Registro e detecção de colisões.
14. HTTP com timeout e limites.
15. Autenticação Basic, API key e Bearer.
16. Permissões por método, padrão GET.
17. Respostas JSON e texto; status HTTP preservado.
18. Transporte MCP stdio.
19. Logs sem credenciais.
20. Validação de destino/mitigações SSRF.

## Fora do MVP

- OAuth2 completo, refresh tokens e fluxo interativo.
- HMAC e mecanismos específicos não padronizados.
- Paginação automática.
- Cache de respostas.
- Transformações por endpoint.
- Configuração manual de cada operação.
- Interface gráfica.
- Gerenciamento avançado de segredos.
- Sincronização remota de conectores.
- Atualização automática de especificações.
- Suporte integral a todas as versões OpenAPI.
- Plugins com código executável.
- Upload/download binário avançado.
- Hot reload de conectores.
- Outros transportes MCP além de stdio.

## Critérios de aceite

- `init` cria a estrutura sem sobrescrever arquivos existentes.
- O servidor encontra o data-dir mesmo iniciado de outro working directory.
- Variáveis do processo prevalecem sobre `.env`.
- Um conector válido gera ferramentas com nomes estáveis.
- Path, query, headers e JSON body são enviados corretamente.
- Métodos não permitidos são bloqueados antes do HTTP.
- Segredos ausentes impedem ativação do conector e são diagnosticados sem revelar valores.
- Hosts bloqueados e redirecionamentos inseguros são rejeitados.
- Respostas JSON/texto e erros HTTP são devolvidos com status preservado.
- stdout contém somente protocolo MCP.
- Testes cobrem erros de parsing, colisão, parâmetros, autenticação, timeout e segurança.

## Estratégia de entrega

Implementar primeiro um caminho vertical mínimo: carregar um conector, carregar OpenAPI, gerar uma ferramenta GET, registrá-la no MCP e chamar um servidor HTTP de teste. Depois ampliar parâmetros, corpo JSON, autenticação, políticas e respostas. Isso valida cedo a integração entre SDK MCP e `kin-openapi` sem criar toda a arquitetura antes de exercitá-la.
