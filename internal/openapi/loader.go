package openapi

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// LoadDocument carrega e valida um documento OpenAPI 3.0.x a partir de um caminho absoluto.
func LoadDocument(absPath string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = false // Reforçar restrição para MVP: apenas locais

	doc, err := loader.LoadFromFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar arquivo OpenAPI: %w", err)
	}

	// Validar o documento
	if err := doc.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("falha na validação do documento OpenAPI: %w", err)
	}

	return doc, nil
}
