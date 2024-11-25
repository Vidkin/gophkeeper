// Package hash provides functionality for generating cryptographic hashes.
//
// This package includes the GetHashSHA256 function, which computes a SHA-256 hash
// of the provided data combined with a key.
package hash

import "crypto/sha256"

// GetHashSHA256 computes the SHA-256 hash of the given data combined with a key.
//
// Parameters:
//   - key: A string representing the key to be included in the hash computation.
//   - data: A byte slice containing the data to be hashed.
//
// Returns:
//   - A byte slice containing the resulting SHA-256 hash.
//
// The function creates a new SHA-256 hash instance and writes the provided data
// and key into it. It then computes the final hash and returns it as a byte slice.
// This function can be used for generating secure hashes for data integrity
// verification or authentication purposes.
func GetHashSHA256(key string, data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	h.Write([]byte(key))
	return h.Sum(nil)
}
