package tools

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"openapi-mcp/internal/connector"
)

func TestMapperAndRegistry(t *testing.T) {
	manifest := &connector.Manifest{
		ID:   "testapi",
		Name: "Test API",
		Tools: connector.Tools{
			Prefix: "test",
		},
		Permissions: connector.Permissions{
			Methods: []string{"GET", "POST"},
		},
	}

	doc := &openapi3.T{
		Info: &openapi3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: openapi3.Paths{
			"/users": &openapi3.PathItem{
				Get: &openapi3.Operation{
					OperationID: "getUsers",
					Summary:     "Get all users",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: &openapi3.Parameter{
								Name:     "status",
								In:       "query",
								Required: false,
							},
						},
					},
				},
				Post: &openapi3.Operation{
					OperationID: "createUser",
					Summary:     "Create user",
				},
			},
		},
	}

	mapper := NewMapper()
	toolDefs, err := mapper.MapConnector(manifest, doc)
	if err != nil {
		t.Fatalf("MapConnector falhou: %v", err)
	}

	if len(toolDefs) != 2 {
		t.Fatalf("Esperava 2 ferramentas, obteve %d", len(toolDefs))
	}

	registry := NewRegistry()
	for _, td := range toolDefs {
		if err := registry.Register(td); err != nil {
			t.Errorf("Falha ao registrar ferramenta %s: %v", td.Name, err)
		}
	}

	// Test collision
	err = registry.Register(toolDefs[0])
	if err == nil {
		t.Errorf("Deveria ter detectado colisão de ferramentas com mesmo nome")
	}

	// Test retrieval
	got, ok := registry.Get("test_get_users")
	if !ok {
		t.Fatalf("Ferramenta test_get_users não encontrada no registry")
	}
	if got.Operation.OperationID != "getUsers" {
		t.Errorf("OperationID incorreto: %s", got.Operation.OperationID)
	}
}
