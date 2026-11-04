package connector

// Specification defines the OpenAPI document source.
type Specification struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// Server defines the base URL for the API.
type Server struct {
	URL string `yaml:"url"`
}

// Authentication defines the auth mechanism.
type Authentication struct {
	Type         string `yaml:"type"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Value        string `yaml:"value"`
	Name         string `yaml:"name"`
	In           string `yaml:"in"`
	RefreshToken string `yaml:"refresh_token"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	TokenURL     string `yaml:"token_url"`
}

// Permissions defines allowed HTTP methods.
type Permissions struct {
	Methods []string `yaml:"methods"`
}

// Tools defines MCP tool exposure settings.
type Tools struct {
	Enabled *bool  `yaml:"enabled"` // Pointer to distinguish between false and missing
	Prefix  string `yaml:"prefix"`
}

// Manifest represents the structure of connector.yaml.
type Manifest struct {
	ID             string         `yaml:"id"`
	Name           string         `yaml:"name"`
	Description    string         `yaml:"description"`
	Version        string         `yaml:"version"`
	Specification  Specification  `yaml:"specification"`
	Server         Server         `yaml:"server"`
	Authentication Authentication `yaml:"authentication"`
	Permissions    Permissions    `yaml:"permissions"`
	Tools          Tools          `yaml:"tools"`
}
