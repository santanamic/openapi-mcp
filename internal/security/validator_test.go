package security

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Invalid scheme FTP",
			url:     "ftp://example.com/file",
			wantErr: true,
		},
		{
			name:    "Embedded credentials rejected",
			url:     "https://user:pass@example.com/",
			wantErr: true,
		},
		{
			name:    "Loopback IP rejected",
			url:     "http://127.0.0.1:8080/path",
			wantErr: true,
		},
		{
			name:    "Private IP rejected",
			url:     "http://192.168.1.1/path",
			wantErr: true,
		},
		{
			name:    "Localhost rejected",
			url:     "http://localhost/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
