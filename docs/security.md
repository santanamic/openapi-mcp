# Segurança

O servidor transforma chamadas de ferramentas em requisições HTTP autenticadas. Uma ferramenta publicada pode executar operações com efeitos reais; por isso, segurança deve ser aplicada antes da execução.

## Credenciais

- Manifestos armazenam somente referências `${VAR}`.
- Valores vêm do ambiente do processo ou `.env`.
- Não incluir segredos em logs, nomes de ferramentas, descrições, mensagens de erro ou dumps.
- Não devolver headers de autenticação na resposta MCP.
- Não registrar corpos de requisição por padrão.
- Restringir permissões do `.env` no sistema operacional.

## Permissões HTTP

O padrão sem `permissions.methods` é `GET` apenas. Métodos de escrita precisam ser explicitamente autorizados no manifesto. A autorização é verificada no servidor, não confiada ao cliente MCP.

A autorização por método é uma proteção mínima: não substitui controles de acesso na API remota, nem limita quais registros um token pode alterar.

## SSRF e validação de destinos

A validação se aplica a:

1. `specification.url` ao baixar OpenAPI remoto.
2. `server.url` e URL final de cada operação.
3. Destinos após resolução DNS e em cada redirecionamento.

Regras mínimas:

- Aceitar somente `http` e `https`.
- Rejeitar URL com credenciais embutidas (`user:pass@host`).
- Rejeitar loopback, link-local, multicast, unspecified e endereços privados por padrão.
- Revalidar IPs resolvidos e impedir que uma resolução DNS posterior altere o destino para IP bloqueado.
- Desabilitar redirecionamentos ou validar cada destino antes de segui-los; não permitir downgrade HTTPS→HTTP sem política explícita.
- Limitar tamanho e tempo de download da especificação.
- Não aceitar protocolos como `file`, `ftp`, `gopher` ou esquemas não HTTP.

Para APIs internas em rede privada, prever uma autorização explícita no nível de configuração do conector/global, documentada e restrita. Essa exceção não deve ser habilitada por padrão. O desenho exato da allowlist será fechado na implementação, evitando uma opção ampla que desative todas as proteções sem necessidade.

## TLS e HTTP

Usar validação TLS padrão do Go. Não oferecer opção de ignorar certificados no MVP. Aplicar timeout total, limites de tamanho de resposta e limites razoáveis de cabeçalhos. Não reutilizar cookies entre conectores.

## Dados não confiáveis

Descrições, schemas, exemplos e respostas OpenAPI são dados externos. Não interpretar descrições como instruções executáveis. Não executar código fornecido por conectores. Validar schemas e impor limites de profundidade/tamanho ao processar JSON.

## Redirecionamentos e autenticação

Não encaminhar credenciais para outro host por redirecionamento. A política recomendada no MVP é não seguir redirecionamentos cross-host; se redirecionamentos forem seguidos, validar o destino e remover credenciais sensíveis quando a origem mudar.

## Logs

Registrar conector, nome da ferramenta, método, host, duração, status HTTP e categoria de erro. Evitar query strings quando possam conter tokens; mascarar valores sensíveis de headers. Nunca registrar Authorization, API keys ou passwords.

## Transporte stdio

- stdout é canal exclusivo do protocolo MCP.
- Logs e diagnósticos vão para stderr/arquivos.
- Não aceitar argumentos de processo ou configuração que alterem silenciosamente o protocolo.

## Limitações

O servidor não substitui autenticação/autorização do provedor da API. Um cliente MCP com acesso ao servidor pode invocar todas as ferramentas expostas; por isso, conectores e tokens devem seguir o princípio do menor privilégio.
