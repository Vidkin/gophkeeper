package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pkgJwt "github.com/Vidkin/gophkeeper/pkg/jwt"
)

func TestValidateToken_ValidToken(t *testing.T) {
	secretKey := "my_secret_key"
	userID := int64(123)

	tokenString, err := pkgJwt.BuildJWTString(secretKey, userID)
	require.NoError(t, err)

	md := metadata.Pairs("token", tokenString)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	interceptor := ValidateToken(secretKey)

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/gophkeeper.Gophkeeper/GetFiles"}, handler)

	require.NoError(t, err)
	assert.Nil(t, resp)
}

func TestValidateToken_MissingToken(t *testing.T) {
	ctx := context.Background()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	interceptor := ValidateToken("my_secret_key")

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/gophkeeper.Gophkeeper/GetFiles"}, handler)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, status.Code(err), codes.PermissionDenied)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	md := metadata.Pairs("token", invalidToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	interceptor := ValidateToken("my_secret_key")

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/gophkeeper.Gophkeeper/GetFiles"}, handler)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, status.Code(err), codes.PermissionDenied)
}

func TestValidateToken_EmptyKey(t *testing.T) {
	secretKey := "my_secret_key"
	userID := int64(123)
	tokenString, err := pkgJwt.BuildJWTString(secretKey, userID)
	require.NoError(t, err)

	md := metadata.Pairs("token", tokenString)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	interceptor := ValidateToken("")

	resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: GrpcEchoMethod}, handler)

	require.NoError(t, err)
	assert.Nil(t, resp)
}
