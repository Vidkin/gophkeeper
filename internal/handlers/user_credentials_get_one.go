package handlers

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/proto"
)

// GetUserCredential retrieves a specific user credential associated with the user from the storage.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.GetUserCredentialRequest structure, which contains the ID of the credential
//     to be retrieved.
//
// Returns:
//   - A pointer to the proto.GetUserCredentialResponse containing the requested user credential.
//   - An error if the operation fails, for example, if the credential ID is invalid or if there is an
//     internal error while retrieving the credential from the database.
//
// The function first attempts to convert the credential ID from a string to an integer. If the conversion
// fails, it returns an InvalidArgument error. It then retrieves the credential from the storage using the
// credential ID. If an error occurs during the retrieval, it logs the error and returns an Internal status.
// If the operation is successful, it constructs a response containing the credential and returns it.
func (g *GophkeeperServer) GetUserCredential(ctx context.Context, in *proto.GetUserCredentialRequest) (*proto.GetUserCredentialResponse, error) {
	credID, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing credentials id")
	}
	var response proto.GetUserCredentialResponse

	cred, err := g.Storage.GetUserCredential(ctx, int64(credID))
	if err != nil {
		logger.Log.Error("error get credentials from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get credentials from DB")
	}

	protoCreds := &proto.Credentials{
		Login:       cred.Login,
		Password:    cred.Password,
		Description: cred.Description,
		Id:          cred.ID,
	}
	response.Credentials = protoCreds
	return &response, nil
}
