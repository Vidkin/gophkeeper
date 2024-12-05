package aes

import (
	"testing"
)

const (
	validAES128Encrypted = "ZigOmJv7WbqOHn4bqE6zwkDemHae2AgpJlMsc4ypUoIC688IIw=="
	validAES192Encrypted = "DLTbVl0pSqoA0p7j4YxeaJ6YttkYrm58kZ/1aMtbml2UPxD6tg=="
	validAES256Encrypted = "klaukCkDbN4VoPJ8dD1x+C92TyHBeO0h8krNeMNC9ZDaQY4w2Q=="
)

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		src       string
		expected  string
		expectErr bool
	}{
		{
			name:      "Valid AES-128 Decryption",
			key:       "examplekey123456", // 16 bytes
			src:       validAES128Encrypted,
			expected:  "plaintext",
			expectErr: false,
		},
		{
			name:      "Valid AES-192 Decryption",
			key:       "examplekey12345678901234", // 24 bytes
			src:       validAES192Encrypted,
			expected:  "plaintext",
			expectErr: false,
		},
		{
			name:      "Valid AES-256 Decryption",
			key:       "examplekey1234567890123456789011", // 32 bytes
			src:       validAES256Encrypted,
			expected:  "plaintext",
			expectErr: false,
		},
		{
			name:      "Invalid Key Length",
			key:       "shortkey",
			src:       validAES128Encrypted,
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid Base64 Input",
			key:       "examplekey1234",
			src:       "invalidbase64",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Nonce Extraction Error",
			key:       "examplekey1234",
			src:       "Y2F0Y2hlc3Q=",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "GCM Decryption Error",
			key:       "examplekey1234",
			src:       "Y2F0Y2hlc3Q=",
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decrypt(tt.key, tt.src)
			if (err != nil) != tt.expectErr {
				t.Errorf("Decrypt() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if result != tt.expected {
				t.Errorf("Decrypt() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
