// Package interceptors provides gRPC interceptors for handling requests and responses.
//
// This package includes the ValidateToken function, which validates JWT tokens for
// incoming requests to secure gRPC methods.
package interceptors

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	jwtPKG "github.com/Vidkin/gophkeeper/pkg/jwt"
)

type contextKey string

const (
	GrpcRegisterUserMethod            = "/gophkeeper.Gophkeeper/RegisterUser"
	GrpcAuthorizeMethod               = "/gophkeeper.Gophkeeper/Authorize"
	GrpcEchoMethod                    = "/gophkeeper.Gophkeeper/Echo"
	UserID                 contextKey = "UserID"
)

// ValidateToken returns a gRPC unary server interceptor that validates JWT tokens
// for incoming requests, allowing access to secured methods based on the token's validity.
//
// Parameters:
//   - key: A string representing the secret key used for signing the JWT tokens.
//
// The interceptor checks if the incoming request's method is one of the public methods
// (RegisterUser, Authorize, or Echo). If it is, or if the key is empty, the interceptor
// allows the request to proceed without validation. Otherwise, it extracts the token from
// the metadata of the incoming context and attempts to parse it using the provided key.
//
// Returns:
//   - A function that implements the gRPC UnaryHandler signature, which processes the
//     request if the token is valid, or returns an error if the token is missing or invalid.
//
// If the token is successfully parsed and validated, the interceptor extracts the UserID
// from the claims and stores it in the context for further use in the request handling.
func ValidateToken(key string) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == GrpcRegisterUserMethod ||
			info.FullMethod == GrpcAuthorizeMethod ||
			info.FullMethod == GrpcEchoMethod {
			return handler(ctx, req)
		}
		if key == "" {
			return handler(ctx, req)
		}

		var tokenString string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("token")
			if len(values) > 0 {
				tokenString = values[0]
			}
		}
		if len(tokenString) == 0 {
			return nil, status.Error(codes.PermissionDenied, "missing token")
		}

		claims := &jwtPKG.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					logger.Log.Error("unexpected signing method")
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(key), nil
			})

		if err != nil {
			logger.Log.Error("error parse claims", zap.Error(err))
			return nil, status.Errorf(codes.PermissionDenied, "error parse claims")
		}

		if !token.Valid {
			logger.Log.Error("token  is not valid")
			return nil, status.Errorf(codes.PermissionDenied, "token is not valid")
		}

		ctx = context.WithValue(ctx, UserID, claims.UserID)
		return handler(ctx, req)
	}
}
