package execution

import (
	"fmt"
	"strings"

	"openapi-mcp/internal/connector"
)

// CheckPolicy validates if the operation method is allowed by connector permissions.
func CheckPolicy(manifest *connector.Manifest, method string) error {
	method = strings.ToUpper(method)
	allowedMethods := make(map[string]bool)

	if len(manifest.Permissions.Methods) == 0 {
		// Default to GET only
		allowedMethods["GET"] = true
	} else {
		for _, m := range manifest.Permissions.Methods {
			allowedMethods[strings.ToUpper(m)] = true
		}
	}

	if !allowedMethods[method] {
		return fmt.Errorf("método HTTP '%s' não autorizado pelas permissões do conector '%s'", method, manifest.ID)
	}

	return nil
}
