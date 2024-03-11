package common

import (
	"testing"
)

func TestGenerateServerID(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Standard HTTP URL",
			url:     "http://127.0.0.1:5000",
			wantErr: false,
		},
		{
			name:    "Standard HTTPS URL",
			url:     "https://example.com:443",
			wantErr: false,
		},
		{
			name:    "URL Without Port",
			url:     "http://example.com",
			wantErr: true, // Expect error due to missing port
		},
		{
			name:    "Invalid URL",
			url:     "htp://abc",
			wantErr: true, // Expect error due to invalid URL format
		},
		{
			name:    "URL With Uncommon Port",
			url:     "http://example.com:8080",
			wantErr: false,
		},
		{
			name:    "URL With IPv6 Address",
			url:     "http://[::1]:3000",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateServerID(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateServerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Additional checks can be added here, such as verifying the format of the returned ID,
			// if the error message is as expected, etc.
		})
	}
}
