package execution

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"openapi-mcp/internal/connector"
	"openapi-mcp/internal/tools"
)

func TestCheckPolicy(t *testing.T) {
	manifest := &connector.Manifest{
		ID: "test",
		Permissions: connector.Permissions{
			Methods: []string{"GET", "POST"},
		},
	}

	if err := CheckPolicy(manifest, "GET"); err != nil {
		t.Errorf("expected GET to be allowed: %v", err)
	}
	if err := CheckPolicy(manifest, "POST"); err != nil {
		t.Errorf("expected POST to be allowed: %v", err)
	}
	if err := CheckPolicy(manifest, "DELETE"); err == nil {
		t.Errorf("expected DELETE to be denied")
	}
}

func TestBuildRequest(t *testing.T) {
	op := &openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "id",
					In:       "path",
					Required: true,
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:     "filter",
					In:       "query",
					Required: false,
				},
			},
		},
	}

	opDef := tools.OperationDefinition{
		ConnectorID: "test",
		Method:      "GET",
		Path:        "/users/{id}",
		Operation:   op,
	}

	args := map[string]any{
		"id":     "123",
		"filter": "active",
	}

	req, err := BuildRequest(opDef, "https://api.example.com/v1", args)
	if err != nil {
		t.Fatalf("unexpected error building request: %v", err)
	}

	expectedURL := "https://api.example.com/v1/users/123?filter=active"
	if req.URL.String() != expectedURL {
		t.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL.String())
	}
	if req.Method != "GET" {
		t.Errorf("expected method GET, got '%s'", req.Method)
	}
}
