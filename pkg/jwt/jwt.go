// Package jwt provides functionality for creating and managing JSON Web Tokens (JWT).
//
// This package includes the Claims struct and the BuildJWTString function for generating
// signed JWTs with user-specific claims.
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the custom claims for the JWT, including registered claims
// and a UserID field.
type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

// TokenExpireTime defines the duration for which the JWT is valid.
const TokenExpireTime = time.Hour * 1

// BuildJWTString generates a signed JWT string for a given user ID using the provided secret key.
//
// Parameters:
//   - secretKey: A string representing the secret key used for signing the JWT.
//   - userID: An int64 representing the user ID to be included in the JWT claims.
//
// Returns:
//   - A string containing the signed JWT.
//   - An error if the token could not be created or signed.
//
// The function creates a new JWT with claims that include an expiration time set to
// one hour from the current time and the specified user ID. It then signs the token
// using the provided secret key and returns the resulting token string. If any error
// occurs during the signing process, it returns an error.
func BuildJWTString(secretKey string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireTime)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
