package openapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDocument(t *testing.T) {
	tempDir := t.TempDir()
	
	// Valid OpenAPI
	validYaml := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths: {}
`
	validPath := filepath.Join(tempDir, "valid.yaml")
	os.WriteFile(validPath, []byte(validYaml), 0600)

	doc, err := LoadDocument(validPath)
	if err != nil {
		t.Fatalf("LoadDocument falhou para YAML válido: %v", err)
	}
	if doc.Info.Title != "Test API" {
		t.Errorf("Título incorreto: %s", doc.Info.Title)
	}

	// Invalid OpenAPI
	invalidYaml := `invalid: yaml`
	invalidPath := filepath.Join(tempDir, "invalid.yaml")
	os.WriteFile(invalidPath, []byte(invalidYaml), 0600)

	_, err = LoadDocument(invalidPath)
	if err == nil {
		t.Errorf("LoadDocument deveria falhar para YAML inválido")
	}
}
