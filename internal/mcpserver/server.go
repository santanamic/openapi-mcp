package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"openapi-mcp/internal/connector"
	"openapi-mcp/internal/execution"
	"openapi-mcp/internal/tools"
)

type Server struct {
	mcpServer  *mcp.Server
	registry   *tools.Registry
	connectors map[string]*connector.Manifest
	executor   *execution.Executor
}

func NewServer(version string, registry *tools.Registry, connectors map[string]*connector.Manifest) *Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "openapi-mcp",
		Version: version,
	}, nil)

	srv := &Server{
		mcpServer:  s,
		registry:   registry,
		connectors: connectors,
		executor:   execution.NewExecutor(),
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
			connectorID := toolDef.Operation.ConnectorID
			manifest, ok := s.connectors[connectorID]
			if !ok {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Conector não encontrado para a ferramenta: %s", connectorID),
						},
					},
				}, nil
			}

			args := make(map[string]any)
			if len(req.Params.Arguments) > 0 {
				_ = json.Unmarshal(req.Params.Arguments, &args)
			}

			result, err := s.executor.Execute(ctx, manifest, toolDef.Operation, args)
			if err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Erro ao executar operação: %v", err),
						},
					},
				}, nil
			}

			return result, nil
		})
	}
}

func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

func (s *Server) RunWithTransport(ctx context.Context, transport mcp.Transport) error {
	return s.mcpServer.Run(ctx, transport)
}
