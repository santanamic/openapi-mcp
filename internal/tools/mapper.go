package tools

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"openapi-mcp/internal/connector"
)

// OperationDefinition represents an operation independent of the provider.
type OperationDefinition struct {
	ConnectorID string
	Method      string
	Path        string
	OperationID string
	Summary     string
	Description string
	Operation   *openapi3.Operation
}

// ToolDefinition represents an MCP tool definition.
type ToolDefinition struct {
	Name        string
	Description string
	InputSchema map[string]any // JSON Schema for tool inputs
	Operation   OperationDefinition
}

// Mapper maps OpenAPI operations to ToolDefinitions for a given connector.
type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

// MapConnector maps all supported operations of a connector's OpenAPI document.
func (m *Mapper) MapConnector(manifest *connector.Manifest, doc *openapi3.T) ([]ToolDefinition, error) {
	var tools []ToolDefinition
	prefix := manifest.Tools.Prefix
	if prefix == "" {
		prefix = manifest.ID
	}

	// Allowed methods from permissions
	allowedMethods := make(map[string]bool)
	if len(manifest.Permissions.Methods) == 0 {
		// Default to GET if not specified
		allowedMethods["GET"] = true
	} else {
		for _, m := range manifest.Permissions.Methods {
			allowedMethods[strings.ToUpper(m)] = true
		}
	}

	if doc.Paths == nil {
		return tools, nil
	}

	for path, pathItem := range doc.Paths {
		if pathItem == nil {
			continue
		}

		methods := map[string]*openapi3.Operation{
			"GET":     pathItem.Get,
			"POST":    pathItem.Post,
			"PUT":     pathItem.Put,
			"DELETE":  pathItem.Delete,
			"PATCH":   pathItem.Patch,
			"OPTIONS": pathItem.Options,
			"HEAD":    pathItem.Head,
		}

		for method, op := range methods {
			if op == nil {
				continue
			}
			if !allowedMethods[method] {
				continue
			}

			toolName, err := deriveToolName(prefix, method, path, op)
			if err != nil {
				return nil, fmt.Errorf("falha ao derivar nome da ferramenta para %s %s: %w", method, path, err)
			}

			description := op.Description
			if description == "" {
				description = op.Summary
			}
			if description == "" {
				description = fmt.Sprintf("Executa %s em %s", method, path)
			}

			// Build input schema from parameters and requestBody
			inputSchema := buildInputSchema(pathItem, op)

			opDef := OperationDefinition{
				ConnectorID: manifest.ID,
				Method:      method,
				Path:        path,
				OperationID: op.OperationID,
				Summary:     op.Summary,
				Description: op.Description,
				Operation:   op,
			}

			tools = append(tools, ToolDefinition{
				Name:        toolName,
				Description: description,
				InputSchema: inputSchema,
				Operation:   opDef,
			})
		}
	}

	return tools, nil
}

func deriveToolName(prefix, method, path string, op *openapi3.Operation) (string, error) {
	var base string
	if op.OperationID != "" {
		base = toSnakeCase(op.OperationID)
	} else {
		// Derive from method and path: e.g. GET /users/{id} -> get_users_id
		cleanPath := normalizePathToSnake(path)
		base = strings.ToLower(method) + "_" + cleanPath
	}

	cleanPrefix := toSnakeCase(prefix)
	if cleanPrefix == "" {
		return base, nil
	}
	return cleanPrefix + "_" + base, nil
}

func toSnakeCase(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('_')
		}
	}
	// Clean up double underscores
	result := sb.String()
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.Trim(result, "_")
}

func normalizePathToSnake(path string) string {
	// Replace {param} or :param with param, and non-alphanumeric with underscore
	var sb strings.Builder
	inParam := false
	for _, r := range path {
		if r == '{' || r == ':' {
			inParam = true
			continue
		}
		if r == '}' {
			inParam = false
			sb.WriteRune('_')
			continue
		}
		if inParam {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				sb.WriteRune(unicode.ToLower(r))
			}
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune('_')
		}
	}
	result := sb.String()
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.Trim(result, "_")
}

func buildInputSchema(pathItem *openapi3.PathItem, op *openapi3.Operation) map[string]any {
	properties := make(map[string]any)
	var required []string

	// Combine path-level and operation-level parameters
	allParams := append([]*openapi3.ParameterRef{}, pathItem.Parameters...)
	allParams = append(allParams, op.Parameters...)

	for _, paramRef := range allParams {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}
		param := paramRef.Value
		paramSchema := map[string]any{
			"type": "string",
		}
		if param.Schema != nil && param.Schema.Value != nil {
			if param.Schema.Value.Type != "" {
				paramSchema["type"] = param.Schema.Value.Type
			}
			if param.Schema.Value.Description != "" {
				paramSchema["description"] = param.Schema.Value.Description
			}
		}
		if param.Description != "" {
			paramSchema["description"] = param.Description
		}

		properties[param.Name] = paramSchema
		if param.Required {
			required = append(required, param.Name)
		}
	}

	// Request Body if present
	if op.RequestBody != nil && op.RequestBody.Value != nil {
		rb := op.RequestBody.Value
		if content, ok := rb.Content["application/json"]; ok && content.Schema != nil {
			properties["request_body"] = map[string]any{
				"type":        "object",
				"description": "JSON Request Body",
			}
			if rb.Required {
				required = append(required, "request_body")
			}
		}
	}

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}
