package aes

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		src       string
		expectErr bool
	}{
		{
			name:      "Valid AES-128 Encryption",
			key:       "examplekey123456", // 16 bytes
			src:       "plaintext",
			expectErr: false,
		},
		{
			name:      "Valid AES-192 Encryption",
			key:       "examplekey12345678901234", // 24 bytes
			src:       "plaintext",
			expectErr: false,
		},
		{
			name:      "Valid AES-256 Encryption",
			key:       "examplekey1234567890123456789011", // 32 bytes
			src:       "plaintext",
			expectErr: false,
		},
		{
			name:      "Invalid Key Length (Too Short)",
			key:       "shortkey", // Invalid key length
			src:       "plaintext",
			expectErr: true,
		},
		{
			name:      "Invalid Key Length (Too Long)",
			key:       "thiskeyiswaytoolong1234567890", // Invalid key length
			src:       "plaintext",
			expectErr: true,
		},
		{
			name:      "Empty Plaintext",
			key:       "examplekey123456", // 16 bytes
			src:       "",
			expectErr: false, // Should succeed, even with empty plaintext
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encrypt(tt.key, tt.src)
			fmt.Println(result)
			if (err != nil) != tt.expectErr {
				t.Errorf("Encrypt() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr {
				// Check if the result is a valid base64 string
				if _, err := base64.StdEncoding.DecodeString(result); err != nil {
					t.Errorf("Encrypt() returned invalid base64 string: %v", result)
				}
			}
		})
	}
}
