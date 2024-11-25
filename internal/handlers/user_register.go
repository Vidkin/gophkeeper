package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/pkg/aes"
	"github.com/Vidkin/gophkeeper/proto"
)

// RegisterUser registers a new user in the system.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.RegisterUserRequest structure, which contains the user's credentials
//     (login and password) for registration.
//
// Returns:
//   - A pointer to an empty proto.Empty response indicating successful registration of the user.
//   - An error if the operation fails, for example, if the user login or password is invalid, if the
//     user already exists, or if there is an internal error while encrypting the password or adding the
//     user to the storage.
//
// The function first checks if the user login and password are provided in the request. If either is
// missing, it logs an error and returns an InvalidArgument status. It then checks if the user already
// exists in the storage. If the user exists, it logs an error and returns an AlreadyExists status. If
// the user does not exist, it encrypts the password using AES encryption. If the encryption fails, it
// logs the error and returns an Internal status. Finally, if the user is successfully added to the
// storage, it returns an empty response.
func (g *GophkeeperServer) RegisterUser(ctx context.Context, in *proto.RegisterUserRequest) (*emptypb.Empty, error) {
	if in.Credentials.Login == "" {
		logger.Log.Error("invalid user login")
		return nil, status.Errorf(codes.InvalidArgument, "invalid user login")
	}
	if in.Credentials.Password == "" {
		logger.Log.Error("invalid user password")
		return nil, status.Errorf(codes.InvalidArgument, "invalid user password")
	}

	_, err := g.Storage.GetUser(ctx, in.Credentials.Login)
	if err == nil {
		logger.Log.Error("user already exists")
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	encPwd, err := aes.Encrypt(g.DatabaseKey, in.Credentials.Password)
	if err != nil {
		logger.Log.Error("error encrypt password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt password")
	}

	if err := g.Storage.AddUser(ctx, in.Credentials.Login, encPwd); err != nil {
		logger.Log.Error("error create user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error create user")
	}
	return &emptypb.Empty{}, nil
}
