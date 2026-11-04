package auth

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"openapi-mcp/internal/connector"
)

// ApplyAuth applies authentication credentials to an HTTP request.
func ApplyAuth(req *http.Request, auth connector.Authentication) error {
	authType := strings.ToLower(auth.Type)
	if authType == "" || authType == "none" {
		return nil
	}

	switch authType {
	case "basic":
		req.SetBasicAuth(auth.Username, auth.Password)
	case "bearer":
		if auth.Value == "" {
			return fmt.Errorf("token bearer vazio para autenticação bearer")
		}
		req.Header.Set("Authorization", "Bearer "+auth.Value)
	case "oauth2":
		tokenStr := auth.Value
		if auth.RefreshToken != "" && auth.TokenURL != "" && auth.ClientID != "" {
			conf := &oauth2.Config{
				ClientID:     auth.ClientID,
				ClientSecret: auth.ClientSecret,
				Endpoint: oauth2.Endpoint{
					TokenURL: auth.TokenURL,
				},
			}
			existingToken := &oauth2.Token{
				AccessToken:  auth.Value,
				RefreshToken: auth.RefreshToken,
			}
			ts := conf.TokenSource(req.Context(), existingToken)
			newToken, err := ts.Token()
			if err != nil {
				return fmt.Errorf("falha ao renovar token OAuth2 com refresh token: %w. Por favor, reexecute o comando 'openapi-mcp auth'", err)
			}
			if newToken != nil {
				tokenStr = newToken.AccessToken
			}
		}

		if tokenStr == "" {
			return fmt.Errorf("token de acesso vazio para autenticação oauth2. Por favor, execute 'openapi-mcp auth'")
		}
		req.Header.Set("Authorization", "Bearer "+tokenStr)
	case "api_key", "apikey":
		if auth.Value == "" {
			return fmt.Errorf("valor de chave de API vazio para autenticação api_key")
		}
		name := auth.Name
		if name == "" {
			name = "X-API-Key"
		}
		in := strings.ToLower(auth.In)
		if in == "" || in == "header" {
			req.Header.Set(name, auth.Value)
		} else if in == "query" {
			q := req.URL.Query()
			q.Set(name, auth.Value)
			req.URL.RawQuery = q.Encode()
		} else {
			return fmt.Errorf("localização de api_key desconhecida: '%s'", auth.In)
		}
	default:
		return fmt.Errorf("tipo de autenticação não suportado: '%s'", auth.Type)
	}

	return nil
}
