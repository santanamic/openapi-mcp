package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"openapi-mcp/internal/connector"
	"openapi-mcp/internal/mcpserver"
	"openapi-mcp/internal/openapi"
	"openapi-mcp/internal/tools"
)

func TestE2EIntegrationFlow(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	// 1. Create a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/items" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id": "1", "name": "Item One"}]`))
			return
		}
		if r.Method == "POST" && r.URL.Path == "/items" {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": "created",
				"data":   body,
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	// 2. Create temporary directory for connector & openapi
	tmpDir, err := os.MkdirTemp("", "openapi-mcp-e2e-*")
	if err != nil {
		t.Fatalf("falha ao criar dir temporário: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	connectorDir := filepath.Join(tmpDir, "testapi")
	if err := os.MkdirAll(connectorDir, 0755); err != nil {
		t.Fatalf("falha ao criar dir do conector: %v", err)
	}

	connectorYAML := `
id: testapi
name: Test API
version: "1.0.0"
server:
  url: ` + ts.URL + `
specification:
  path: openapi.yaml
permissions:
  methods: ["GET", "POST"]
tools:
  prefix: testapi
`
	if err := os.WriteFile(filepath.Join(connectorDir, "connector.yaml"), []byte(connectorYAML), 0644); err != nil {
		t.Fatalf("falha ao escrever connector.yaml: %v", err)
	}

	openapiYAML := `
openapi: 3.0.3
info:
  title: Test API
  version: "1.0.0"
paths:
  /items:
    get:
      operationId: listItems
      summary: List items
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
    post:
      operationId: createItem
      summary: Create item
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
      responses:
        '201':
          description: Created
          content:
            application/json:
              schema:
                type: object
`
	if err := os.WriteFile(filepath.Join(connectorDir, "openapi.yaml"), []byte(openapiYAML), 0644); err != nil {
		t.Fatalf("falha ao escrever openapi.yaml: %v", err)
	}

	// 3. Load connector and OpenAPI
	manager := connector.NewManager(tmpDir)
	connectors, err := manager.LoadAll()
	if err != nil {
		t.Fatalf("falha ao carregar conectores: %v", err)
	}

	manifest, ok := connectors["testapi"]
	if !ok {
		t.Fatalf("conector testapi não encontrado")
	}

	doc, err := openapi.LoadDocument(filepath.Join(connectorDir, manifest.Specification.Path))
	if err != nil {
		t.Fatalf("falha ao carregar openapi: %v", err)
	}

	// 4. Map tools and build registry
	mapper := tools.NewMapper()
	toolDefs, err := mapper.MapConnector(manifest, doc)
	if err != nil {
		t.Fatalf("falha ao mapear ferramentas: %v", err)
	}

	registry := tools.NewRegistry()
	for _, td := range toolDefs {
		if err := registry.Register(td); err != nil {
			t.Fatalf("falha ao registrar ferramenta %s: %v", td.Name, err)
		}
	}

	// 5. Start MCP Server with in-memory transport
	srv := mcpserver.NewServer("1.0.0", registry, connectors)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = srv.RunWithTransport(ctx, serverTransport)
	}()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-e2e-client", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("cliente falhou ao conectar: %v", err)
	}
	defer session.Close()

	// 6. Test listing tools
	listRes, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools falhou: %v", err)
	}

	if len(listRes.Tools) != 2 {
		t.Fatalf("esperava 2 ferramentas, obteve %d", len(listRes.Tools))
	}

	// 7. Test calling GET tool
	callGetRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "testapi_list_items",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("CallTool GET falhou: %v", err)
	}
	if callGetRes.IsError {
		for _, c := range callGetRes.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				t.Logf("GET Error Content: %s", tc.Text)
			}
		}
		t.Fatalf("CallTool GET retornou erro")
	}

	// 8. Test calling POST tool
	callPostRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "testapi_create_item",
		Arguments: map[string]any{
			"request_body": map[string]any{
				"name": "New Item",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool POST falhou: %v", err)
	}
	if callPostRes.IsError {
		for _, c := range callPostRes.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				t.Logf("POST Error Content: %s", tc.Text)
			}
		}
		t.Fatalf("CallTool POST retornou erro")
	}
}
