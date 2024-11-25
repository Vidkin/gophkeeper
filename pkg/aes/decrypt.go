// Package aes provides functionality for AES encryption and decryption.
//
// This package includes the Decrypt function, which decrypts a base64-encoded string
// using the AES algorithm in GCM mode.
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// Decrypt decrypts a base64-encoded string using the provided AES key.
//
// Parameters:
//   - key: A string representing the AES key used for decryption. The key must be either
//     16, 24, or 32 bytes long, corresponding to AES-128, AES-192, or AES-256.
//   - src: A base64-encoded string representing the encrypted data, which includes the nonce
//     used during encryption.
//
// Returns:
//   - A string containing the decrypted plaintext.
//   - An error if the decryption process fails, including issues with the key length,
//     base64 decoding, or GCM decryption.
//
// The function first decodes the base64-encoded input string. It then creates a new AES cipher
// block using the provided key. A GCM (Galois/Counter Mode) cipher is created for decryption.
// The function extracts the nonce from the beginning of the decoded data and separates the
// ciphertext. Finally, it attempts to decrypt the ciphertext using the nonce and returns the
// resulting plaintext. If any step fails, an appropriate error is returned.
func Decrypt(key, src string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("can't extract nonce")
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	res, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
