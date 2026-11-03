package tools

import (
	"fmt"
)

// Registry maintains the index of MCP tools and ensures unique names.
type Registry struct {
	tools map[string]ToolDefinition
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]ToolDefinition),
	}
}

// Register adds a tool definition to the registry. Returns an error if a collision occurs.
func (r *Registry) Register(tool ToolDefinition) error {
	if existing, ok := r.tools[tool.Name]; ok {
		return fmt.Errorf("colisão de nome de ferramenta '%s': já registrada pelo conector '%s' (%s %s)",
			tool.Name, existing.Operation.ConnectorID, existing.Operation.Method, existing.Operation.Path)
	}
	r.tools[tool.Name] = tool
	return nil
}

// Get retrieves a tool definition by name.
func (r *Registry) Get(name string) (ToolDefinition, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// ListAll returns all registered tool definitions.
func (r *Registry) ListAll() []ToolDefinition {
	list := make([]ToolDefinition, 0, len(r.tools))
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}
