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
