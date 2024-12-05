package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Vidkin/gophkeeper/internal/logger"
)

func TestBuildJWTString_ValidInput(t *testing.T) {
	secretKey := "my_secret_key"
	userID := int64(123)

	tokenString, err := BuildJWTString(secretKey, userID)

	require.NoError(t, err)

	assert.NotEmpty(t, tokenString)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Log.Error("unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})

	require.NoError(t, err)
	assert.NotNil(t, token)

	assert.Equal(t, userID, claims.UserID)

	assert.False(t, claims.ExpiresAt.Time.Before(time.Now()))
}
