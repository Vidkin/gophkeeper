// Package aes provides functionality for AES encryption and decryption.
//
// This package includes the Encrypt function, which encrypts a plaintext string
// using the AES algorithm in GCM mode.
package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Encrypt encrypts a plaintext string using the provided AES key.
//
// Parameters:
//   - key: A string representing the AES key used for encryption. The key must be either
//     16, 24, or 32 bytes long, corresponding to AES-128, AES-192, or AES-256.
//   - src: A string representing the plaintext data to be encrypted.
//
// Returns:
//   - A base64-encoded string containing the encrypted data, which includes the nonce
//     used during encryption.
//   - An error if the encryption process fails, including issues with the key length,
//     GCM initialization, or random nonce generation.
//
// The function creates a new AES cipher block using the provided key and initializes a GCM
// (Galois/Counter Mode) cipher for encryption. It generates a random nonce of the appropriate
// size and uses it to encrypt the plaintext. The resulting ciphertext, which includes the nonce,
// is then base64-encoded and returned. If any step fails, an appropriate error is returned.
func Encrypt(key, src string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(src), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
