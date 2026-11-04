package security

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// ValidateURL validates a target URL against SSRF and security rules:
// - Must be http or https.
// - Must not contain embedded credentials (userinfo).
// - Host must resolve and not be loopback, private, link-local, multicast, or unspecified (unless allowed).
func ValidateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("url inválida: %w", err)
	}

	// 1. Check scheme
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("esquema não permitido: '%s' (apenas http e https são permitidos)", parsed.Scheme)
	}

	// 2. Check embedded credentials
	if parsed.User != nil {
		return errors.New("credenciais embutidas na url não são permitidas")
	}

	// 3. Check host
	host := parsed.Hostname()
	if host == "" {
		return errors.New("url sem host definido")
	}

	// Check if host is literal IP or domain name
	ip := net.ParseIP(host)
	if ip != nil {
		if err := checkIP(ip); err != nil {
			return fmt.Errorf("endereço IP bloqueado por política de segurança: %w", err)
		}
	} else {
		// Resolve domain name and check all resolved IPs (DNS rebinding protection)
		ips, err := net.LookupIP(host)
		if err != nil {
			return fmt.Errorf("falha ao resolver host '%s': %w", host, err)
		}
		if len(ips) == 0 {
			return fmt.Errorf("nenhum endereço IP encontrado para o host '%s'", host)
		}
		for _, resolvedIP := range ips {
			if err := checkIP(resolvedIP); err != nil {
				return fmt.Errorf("ip resolvido bloqueado para o host '%s': %w", host, err)
			}
		}
	}

	return nil
}

func checkIP(ip net.IP) error {
	if ip.IsLoopback() {
		if os.Getenv("OPENAPI_MCP_ALLOW_LOCAL") == "true" || os.Getenv("APP_ENV") == "test" {
			return nil
		}
		return errors.New("endereço de loopback não permitido")
	}
	if ip.IsPrivate() {
		return errors.New("endereço IP privado não permitido")
	}
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return errors.New("endereço IP link-local não permitido")
	}
	if ip.IsUnspecified() {
		return errors.New("endereço IP não especificado não permitido")
	}
	if ip.IsMulticast() {
		return errors.New("endereço IP multicast não permitido")
	}
	return nil
}
