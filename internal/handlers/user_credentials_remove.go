package handlers

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/proto"
)

// RemoveUserCredentials deletes a specific user credential associated with the user from the storage.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.RemoveUserCredentialsRequest structure, which contains the ID of the
//     credential to be removed.
//
// Returns:
//   - A pointer to an empty proto.Empty response indicating successful removal of the user credential.
//   - An error if the operation fails, for example, if the credential ID is not provided, if the credential
//     ID is invalid, or if there is an internal error while removing the credential from the storage.
//
// The function first checks if the credential ID is provided in the request. If not, it logs an error and
// returns an InvalidArgument status. It then attempts to parse the credential ID from a string to an int64.
// If the parsing fails, it logs the error and returns an InvalidArgument status. If the credential ID is
// valid, it proceeds to remove the credential from the storage. If an error occurs during the removal, it
// logs the error and returns an Internal status. If the operation is successful, it returns an empty response.
func (g *GophkeeperServer) RemoveUserCredentials(ctx context.Context, in *proto.RemoveUserCredentialsRequest) (*emptypb.Empty, error) {
	if in.Id == "" {
		logger.Log.Error("you must provide credentials id")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide credentials id")
	}

	credID, err := strconv.ParseInt(in.Id, 10, 64)
	if err != nil {
		logger.Log.Error("invalid credentials id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials id")
	}

	if err = g.Storage.RemoveUserCredential(ctx, credID); err != nil {
		logger.Log.Error("error remove user credentials", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove user credentials")
	}
	return &emptypb.Empty{}, nil
}
