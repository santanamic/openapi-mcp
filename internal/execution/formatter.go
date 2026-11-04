package execution

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const MaxResponseBodySize = 1024 * 1024 // 1 MB limit

// FormatResponse formats an HTTP response into an MCP CallToolResult.
func FormatResponse(resp *http.Response) (*mcp.CallToolResult, error) {
	defer resp.Body.Close()

	limitedReader := io.LimitReader(resp.Body, MaxResponseBodySize+1)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler corpo da resposta HTTP: %w", err)
	}

	truncated := false
	if int64(len(bodyBytes)) > MaxResponseBodySize {
		truncated = true
		bodyBytes = bodyBytes[:MaxResponseBodySize]
	}

	bodyStr := string(bodyBytes)

	// Try to pretty-print if JSON
	var jsonVal any
	if json.Unmarshal(bodyBytes, &jsonVal) == nil {
		pretty, err := json.MarshalIndent(jsonVal, "", "  ")
		if err == nil {
			bodyStr = string(pretty)
		}
	}

	resultText := fmt.Sprintf("Status HTTP: %d %s\n", resp.StatusCode, resp.Status)
	if resp.StatusCode == http.StatusUnauthorized {
		resultText += "[AÇÃO REQUERIDA: 401 Unauthorized] Credenciais de autenticação expiradas ou inválidas. Por favor, reexecute o comando 'openapi-mcp auth' para atualizar as credenciais no seu ambiente/arquivo .env.\n"
	}
	if truncated {
		resultText += "[Aviso: Resposta truncada por exceder o limite de 1MB]\n"
	}
	resultText += "Corpo:\n" + bodyStr

	isError := resp.StatusCode < 200 || resp.StatusCode >= 300

	return &mcp.CallToolResult{
		IsError: isError,
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: resultText,
			},
		},
	}, nil
}
