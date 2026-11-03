package mcpserver

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"openapi-mcp/internal/tools"
)

func TestMCPServerTools(t *testing.T) {
	registry := tools.NewRegistry()
	td := tools.ToolDefinition{
		Name:        "test_get_item",
		Description: "Get test item",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id": map[string]any{"type": "string"},
			},
		},
		Operation: tools.OperationDefinition{
			ConnectorID: "testcon",
			Method:      "GET",
			Path:        "/items/{id}",
		},
	}
	if err := registry.Register(td); err != nil {
		t.Fatalf("falha ao registrar ferramenta: %v", err)
	}

	srv := NewServer("0.1.0-dev", registry)

	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run server in background
	go func() {
		_ = srv.mcpServer.Run(ctx, serverTransport)
	}()

	// Connect client
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("cliente falhou ao conectar: %v", err)
	}
	defer session.Close()

	// List tools
	listRes, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools falhou: %v", err)
	}

	if len(listRes.Tools) != 1 {
		t.Fatalf("esperava 1 ferramenta, obteve %d", len(listRes.Tools))
	}
	if listRes.Tools[0].Name != "test_get_item" {
		t.Errorf("nome da ferramenta incorreto: %s", listRes.Tools[0].Name)
	}

	// Call tool
	callRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "test_get_item",
		Arguments: map[string]any{"id": "123"},
	})
	if err != nil {
		t.Fatalf("CallTool falhou: %v", err)
	}

	if len(callRes.Content) == 0 {
		t.Fatalf("CallTool não retornou conteúdo")
	}

	textCont, ok := callRes.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("conteúdo esperado do tipo TextContent")
	}

	if textCont.Text == "" {
		t.Errorf("texto de resposta vazio")
	}
}
