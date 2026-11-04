package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
)

// OAuthFlowConfig holds configuration for interactive OAuth2 Authorization Code Flow.
type OAuthFlowConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	Scopes       []string
	Port         int
}

// RunInteractiveOAuthFlow executes the interactive OAuth2 Authorization Code Flow.
func RunInteractiveOAuthFlow(cfg OAuthFlowConfig) (*oauth2.Token, error) {
	port := cfg.Port
	if port <= 0 {
		port = 8080
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		// Fallback to random available port if specified port is busy
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, fmt.Errorf("falha ao iniciar servidor local para callback OAuth2: %w", err)
		}
		port = listener.Addr().(*net.TCPAddr).Port
	}
	defer listener.Close()

	redirectURL := fmt.Sprintf("http://localhost:%d/callback", port)

	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.AuthURL,
			TokenURL: cfg.TokenURL,
		},
		RedirectURL: redirectURL,
		Scopes:      cfg.Scopes,
	}

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return nil, fmt.Errorf("falha ao gerar estado aleatório: %w", err)
	}
	state := hex.EncodeToString(stateBytes)

	tokenChan := make(chan *oauth2.Token, 1)
	errChan := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		queryState := r.URL.Query().Get("state")
		if queryState != state {
			http.Error(w, "Estado OAuth inválido", http.StatusBadRequest)
			errChan <- fmt.Errorf("estado OAuth inválido retornado pelo provedor")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Código de autorização ausente", http.StatusBadRequest)
			errChan <- fmt.Errorf("código de autorização ausente na resposta")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		token, err := oauthCfg.Exchange(ctx, code)
		if err != nil {
			http.Error(w, "Falha na troca do token: "+err.Error(), http.StatusInternalServerError)
			errChan <- fmt.Errorf("falha ao trocar código por token: %w", err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Autenticação Concluída</title></head>
<body style="font-family: sans-serif; text-align: center; padding-top: 50px;">
  <h2 style="color: #2e7d32;">Autenticação realizada com sucesso!</h2>
  <p>Você já pode fechar esta aba e retornar ao terminal do <strong>openapi-mcp</strong>.</p>
</body>
</html>`))

		tokenChan <- token
	})

	srv := &http.Server{Handler: mux}
	go func() {
		_ = srv.Serve(listener)
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	authURL := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf(" Abrindo navegador para autenticação OAuth2:\n %s\n\n", authURL)

	if err := openBrowser(authURL); err != nil {
		fmt.Printf(" Não foi possível abrir o navegador automaticamente. Por favor, acesse o link acima manualmente.\n")
	}

	select {
	case token := <-tokenChan:
		return token, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("tempo limite excedido (5 minutos) aguardando callback de autenticação")
	}
}

func openBrowser(rawURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start()
}
