package execution

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"openapi-mcp/internal/auth"
	"openapi-mcp/internal/connector"
	"openapi-mcp/internal/security"
	"openapi-mcp/internal/tools"
)

// Executor coordinates policy enforcement, request building, authentication, security validation, and HTTP execution.
type Executor struct {
	httpClient *http.Client
}

func NewExecutor() *Executor {
	return &Executor{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			// Disable automatic redirects or handle them securely if needed
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return fmt.Errorf("muitos redirecionamentos")
				}
				// Validate redirect destination
				return security.ValidateURL(req.URL.String())
			},
		},
	}
}

// Execute handles a tool call execution.
func (e *Executor) Execute(ctx context.Context, manifest *connector.Manifest, opDef tools.OperationDefinition, args map[string]any) (*mcp.CallToolResult, error) {
	// 1. Check Policy
	if err := CheckPolicy(manifest, opDef.Method); err != nil {
		return nil, err
	}

	// 2. Build Request
	baseURL := manifest.Server.URL
	if baseURL == "" {
		return nil, fmt.Errorf("url do servidor não definida no manifesto do conector '%s'", manifest.ID)
	}

	req, err := BuildRequest(opDef, baseURL, args)
	if err != nil {
		return nil, fmt.Errorf("falha ao construir requisição: %w", err)
	}

	// 3. Apply Authentication
	if err := auth.ApplyAuth(req, manifest.Authentication); err != nil {
		return nil, fmt.Errorf("falha ao aplicar autenticação: %w", err)
	}

	// 4. Validate Target URL & SSRF Protection
	targetURL := req.URL.String()
	if err := security.ValidateURL(targetURL); err != nil {
		return nil, fmt.Errorf("validação de segurança de destino falhou: %w", err)
	}

	// 5. Execute HTTP Request
	req = req.WithContext(ctx)
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na execução da requisição HTTP: %w", err)
	}

	// 6. Format Response
	return FormatResponse(resp)
}
