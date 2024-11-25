package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/pkg/aes"
	"github.com/Vidkin/gophkeeper/pkg/jwt"
	"github.com/Vidkin/gophkeeper/proto"
)

// Authorize authenticates a user based on provided credentials and generates a JWT token.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.AuthorizeRequest structure, which contains the user's login credentials.
//
// Returns:
//   - A pointer to the proto.AuthorizeResponse containing the generated JWT token if authentication is successful.
//   - An error if the operation fails, for example, if the login or password is invalid, if there is an
//     error retrieving the user from the database, or if there is an internal error during token generation.
//
// The function first checks if the login and password are provided in the request. If either is missing,
// it logs an error and returns a PermissionDenied status. It then attempts to retrieve the user from the
// storage using the provided login. If the user is not found or an error occurs during retrieval, it logs
// the error and returns a PermissionDenied status. The function then decrypts the stored password and
// compares it with the provided password. If they do not match, it logs an error and returns a
// PermissionDenied status. If authentication is successful, it generates a JWT token for the user and
// returns it in the response.
func (g *GophkeeperServer) Authorize(ctx context.Context, in *proto.AuthorizeRequest) (*proto.AuthorizeResponse, error) {
	var response proto.AuthorizeResponse
	if in.Credentials.Login == "" || in.Credentials.Password == "" {
		logger.Log.Error("invalid user login or password")
		return nil, status.Errorf(codes.PermissionDenied, "invalid user login or password")
	}

	u, err := g.Storage.GetUser(ctx, in.Credentials.Login)
	if err != nil {
		logger.Log.Error("error get user from db", zap.Error(err))
		return nil, status.Errorf(codes.PermissionDenied, "invalid user login or password")
	}

	decPwd, err := aes.Decrypt(g.DatabaseKey, u.Password)
	if err != nil {
		logger.Log.Error("error encrypt password", zap.Error(err))
		return nil, status.Errorf(codes.PermissionDenied, "invalid user login or password")
	}

	if in.Credentials.Password != decPwd {
		logger.Log.Error("invalid user login or password", zap.Error(err))
		return nil, status.Errorf(codes.PermissionDenied, "invalid user login or password")
	}

	token, err := jwt.BuildJWTString(g.JWTKey, u.ID)
	if err != nil {
		logger.Log.Error("error build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error build jwt string")
	}

	response.Token = token
	return &response, nil
}
