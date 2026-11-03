package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"openapi-mcp/internal/tools"
)

type Server struct {
	mcpServer *mcp.Server
	registry  *tools.Registry
}

func NewServer(version string, registry *tools.Registry) *Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "openapi-mcp",
		Version: version,
	}, nil)

	srv := &Server{
		mcpServer: s,
		registry:  registry,
	}

	srv.registerTools()
	return srv
}

func (s *Server) registerTools() {
	for _, td := range s.registry.ListAll() {
		toolDef := td // capture
		mcpTool := &mcp.Tool{
			Name:        toolDef.Name,
			Description: toolDef.Description,
			InputSchema: toolDef.InputSchema,
		}

		s.mcpServer.AddTool(mcpTool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			argsJSON, _ := json.Marshal(req.Params.Arguments)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Executando operação %s %s (Conector: %s) com argumentos: %s",
							toolDef.Operation.Method, toolDef.Operation.Path, toolDef.Operation.ConnectorID, string(argsJSON)),
					},
				},
			}, nil
		})
	}
}

func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}
