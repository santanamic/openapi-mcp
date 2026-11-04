package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"openapi-mcp/internal/auth"
)

func main() {
	fs := flag.NewFlagSet("openapi-oauth2", flag.ContinueOnError)
	var clientID string
	var clientSecret string
	var authURL string
	var tokenURL string
	var scopesStr string
	var port int
	var showHelp bool

	fs.StringVar(&clientID, "client-id", "", "ID do Cliente OAuth2")
	fs.StringVar(&clientID, "i", "", "ID do Cliente OAuth2 (abreviação)")
	fs.StringVar(&clientSecret, "client-secret", "", "Segredo do Cliente OAuth2")
	fs.StringVar(&clientSecret, "s", "", "Segredo do Cliente OAuth2 (abreviação)")
	fs.StringVar(&authURL, "auth-url", "", "URL de Autorização OAuth2")
	fs.StringVar(&tokenURL, "token-url", "", "URL de Token OAuth2")
	fs.StringVar(&scopesStr, "scopes", "", "Escopos OAuth2 separados por vírgula")
	fs.IntVar(&port, "port", 8080, "Porta local para o servidor de callback HTTP")
	fs.BoolVar(&showHelp, "help", false, "Exibe esta mensagem de ajuda")
	fs.BoolVar(&showHelp, "h", false, "Exibe esta mensagem de ajuda (abreviação)")

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	if clientID == "" || authURL == "" || tokenURL == "" {
		fmt.Fprintln(os.Stderr, "Erro: Os parâmetros --client-id, --auth-url e --token-url são obrigatórios para o fluxo OAuth2.")
		fmt.Fprintln(os.Stderr, "Uso: openapi-oauth2 --client-id <id> --auth-url <url> --token-url <url> [--client-secret <secret>] [--scopes <scopes>] [--port <port>]")
		os.Exit(1)
	}

	var scopes []string
	if scopesStr != "" {
		for _, s := range strings.Split(scopesStr, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				scopes = append(scopes, s)
			}
		}
	}

	token, err := auth.RunInteractiveOAuthFlow(auth.OAuthFlowConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		Scopes:       scopes,
		Port:         port,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro durante fluxo de autenticação OAuth2: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAutenticação OAuth2 concluída com sucesso!")
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("OAUTH2_ACCESS_TOKEN=%s\n", token.AccessToken)
	if token.RefreshToken != "" {
		fmt.Printf("OAUTH2_REFRESH_TOKEN=%s\n", token.RefreshToken)
	}
	if !token.Expiry.IsZero() {
		fmt.Printf("# Expira em: %s\n", token.Expiry.Format(time.RFC3339))
	}
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Copie os valores acima e insira-os no seu arquivo .env ou variáveis de ambiente.")
}

func printHelp() {
	help := `OpenAPI OAuth2 - Utilitário de Autenticação OAuth2 Interativo

Uso:
  openapi-oauth2 [opções]

Opções:
  -i, --client-id     ID do Cliente OAuth2 (obrigatório).
  -s, --client-secret Segredo do Cliente OAuth2.
      --auth-url      URL de Autorização OAuth2 (obrigatório).
      --token-url     URL de Token OAuth2 (obrigatório).
      --scopes        Escopos OAuth2 separados por vírgula.
      --port          Porta local para o servidor de callback HTTP (padrão: 8080).
  -h, --help          Exibe esta mensagem de ajuda.
`
	fmt.Print(help)
}
