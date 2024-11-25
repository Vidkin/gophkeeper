package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/model"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

// AddUserCredentials stores user credentials in the database.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.AddUserCredentialsRequest structure, which contains the user credentials
//     to be added.
//
// Returns:
//   - A pointer to an empty proto.Empty response indicating successful addition of the user credentials.
//   - An error if the operation fails, for example, if the login or password is not provided, or if there
//     is an internal error while adding the credentials to the storage.
//
// The function first checks if both the login and password are provided in the request. If either is missing,
// it logs an error and returns an InvalidArgument status. It then creates a new Credentials model instance,
// populating it with the user ID (extracted from the context), login, password, and description from the
// request. If an error occurs while adding the credentials to the storage, it logs the error and returns
// an Internal status. If the operation is successful, it returns an empty response.
func (g *GophkeeperServer) AddUserCredentials(ctx context.Context, in *proto.AddUserCredentialsRequest) (*emptypb.Empty, error) {
	if in.Credentials.Login == "" || in.Credentials.Password == "" {
		logger.Log.Error("you must provide: login and password")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide: login and password")
	}

	cred := &model.Credentials{
		UserID:      ctx.Value(interceptors.UserID).(int64),
		Login:       in.Credentials.Login,
		Password:    in.Credentials.Password,
		Description: in.Credentials.Description,
	}

	if err := g.Storage.AddUserCredentials(ctx, cred); err != nil {
		logger.Log.Error("error add user credentials", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error add user credentials")
	}
	return &emptypb.Empty{}, nil
}
