package auth

import (
	"net/http"
	"testing"

	"openapi-mcp/internal/connector"
)

func TestApplyAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com/api", nil)

	// Test Basic
	basicAuth := connector.Authentication{
		Type:     "basic",
		Username: "user",
		Password: "password",
	}
	if err := ApplyAuth(req, basicAuth); err != nil {
		t.Fatalf("unexpected error for basic auth: %v", err)
	}
	if got := req.Header.Get("Authorization"); got == "" {
		t.Errorf("expected Authorization header for basic auth, got empty")
	}

	// Test Bearer
	req2, _ := http.NewRequest("GET", "https://example.com/api", nil)
	bearerAuth := connector.Authentication{
		Type:  "bearer",
		Value: "secret-token",
	}
	if err := ApplyAuth(req2, bearerAuth); err != nil {
		t.Fatalf("unexpected error for bearer auth: %v", err)
	}
	if got := req2.Header.Get("Authorization"); got != "Bearer secret-token" {
		t.Errorf("expected 'Bearer secret-token', got '%s'", got)
	}

	// Test OAuth2
	reqOAuth, _ := http.NewRequest("GET", "https://example.com/api", nil)
	oauthAuth := connector.Authentication{
		Type:  "oauth2",
		Value: "oauth-access-token-123",
	}
	if err := ApplyAuth(reqOAuth, oauthAuth); err != nil {
		t.Fatalf("unexpected error for oauth2 auth: %v", err)
	}
	if got := reqOAuth.Header.Get("Authorization"); got != "Bearer oauth-access-token-123" {
		t.Errorf("expected 'Bearer oauth-access-token-123', got '%s'", got)
	}

	// Test API Key Header
	req3, _ := http.NewRequest("GET", "https://example.com/api", nil)
	apiKeyAuth := connector.Authentication{
		Type:  "api_key",
		Name:  "X-Custom-Key",
		Value: "12345",
		In:    "header",
	}
	if err := ApplyAuth(req3, apiKeyAuth); err != nil {
		t.Fatalf("unexpected error for api_key header: %v", err)
	}
	if got := req3.Header.Get("X-Custom-Key"); got != "12345" {
		t.Errorf("expected '12345', got '%s'", got)
	}

	// Test API Key Query
	req4, _ := http.NewRequest("GET", "https://example.com/api", nil)
	apiKeyQueryAuth := connector.Authentication{
		Type:  "api_key",
		Name:  "api_key",
		Value: "abcde",
		In:    "query",
	}
	if err := ApplyAuth(req4, apiKeyQueryAuth); err != nil {
		t.Fatalf("unexpected error for api_key query: %v", err)
	}
	if got := req4.URL.Query().Get("api_key"); got != "abcde" {
		t.Errorf("expected query param 'abcde', got '%s'", got)
	}
}
