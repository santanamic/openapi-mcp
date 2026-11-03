package connector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Manager handles connector discovery and loading.
type Manager struct {
	connectorsDir string
}

func NewManager(connectorsDir string) *Manager {
	return &Manager{connectorsDir: connectorsDir}
}

// LoadAll scans the directory and loads all valid connectors.
func (m *Manager) LoadAll() (map[string]*Manifest, error) {
	connectors := make(map[string]*Manifest)
	
	entries, err := os.ReadDir(m.connectorsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return connectors, nil // Or return error if connectors dir required
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		connectorDir := filepath.Join(m.connectorsDir, entry.Name())
		manifestPath := filepath.Join(connectorDir, "connector.yaml")
		
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue // Skip directories without manifest
		}
		
		manifest, err := m.loadManifest(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar manifesto %s: %w", manifestPath, err)
		}
		
		// Apply defaults and validation
		if err := m.validateAndNormalize(manifest); err != nil {
			return nil, fmt.Errorf("manifesto inválido %s: %w", manifestPath, err)
		}
		
		connectors[manifest.ID] = manifest
	}
	
	return connectors, nil
}

func (m *Manager) loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	// Expand ${VAR} in the raw data
	expanded := os.Expand(string(data), func(key string) string {
		return os.Getenv(key) // Resolver usando ambiente do SO
	})
	
	var manifest Manifest
	if err := yaml.Unmarshal([]byte(expanded), &manifest); err != nil {
		return nil, err
	}
	
	return &manifest, nil
}

func (m *Manager) validateAndNormalize(manifest *Manifest) error {
	if manifest.ID == "" {
		return fmt.Errorf("id é obrigatório")
	}
	if manifest.Name == "" {
		return fmt.Errorf("name é obrigatório")
	}
	if manifest.Specification.Path == "" && manifest.Specification.URL == "" {
		return fmt.Errorf("especificação (path ou url) é obrigatória")
	}
	
	// Apply defaults
	if manifest.Tools.Enabled == nil {
		enabled := true
		manifest.Tools.Enabled = &enabled
	}
	if manifest.Tools.Prefix == "" {
		manifest.Tools.Prefix = manifest.ID
	}
	
	// Normalize methods
	for i, method := range manifest.Permissions.Methods {
		manifest.Permissions.Methods[i] = strings.ToUpper(method)
	}
	
	return nil
}
