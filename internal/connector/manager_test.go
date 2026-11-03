package connector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManagerLoadAll(t *testing.T) {
	tempDir := t.TempDir()
	connectorsDir := filepath.Join(tempDir, "connectors")
	os.Mkdir(connectorsDir, 0755)

	// Criar conector de teste
	wcDir := filepath.Join(connectorsDir, "wc")
	os.Mkdir(wcDir, 0755)
	
	yamlContent := `
id: woocommerce
name: WooCommerce
specification:
  path: ./openapi.yaml
authentication:
  type: basic
  username: ${WC_KEY}
`
	os.WriteFile(filepath.Join(wcDir, "connector.yaml"), []byte(yamlContent), 0600)
	
	// Definir variável de ambiente
	os.Setenv("WC_KEY", "secret_key")
	defer os.Unsetenv("WC_KEY")

	mgr := NewManager(connectorsDir)
	connectors, err := mgr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll falhou: %v", err)
	}

	if len(connectors) != 1 {
		t.Fatalf("Esperava 1 conector, obteve %d", len(connectors))
	}

	c := connectors["woocommerce"]
	if c.Authentication.Username != "secret_key" {
		t.Errorf("Variável WC_KEY não resolvida. Esperava 'secret_key', obteve '%s'", c.Authentication.Username)
	}
	
	if !*c.Tools.Enabled {
		t.Errorf("Tools.Enabled deveria ser true por padrão")
	}
}
