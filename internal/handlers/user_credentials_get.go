package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

// GetUserCredentials retrieves all user credentials associated with the user from the storage.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - _: A pointer to the proto.GetUserCredentialsRequest structure (not used in this method).
//
// Returns:
//   - A pointer to the proto.GetUserCredentialsResponse containing the list of user credentials.
//   - An error if the operation fails, for example, if there is an internal error while retrieving the
//     credentials from the database.
//
// The function fetches the user's credentials from the storage using the user ID extracted from the context.
// If an error occurs during the retrieval, it logs the error and returns an Internal status. If the
// operation is successful, it constructs a response containing the credentials and returns it.
func (g *GophkeeperServer) GetUserCredentials(ctx context.Context, _ *proto.GetUserCredentialsRequest) (*proto.GetUserCredentialsResponse, error) {
	var response proto.GetUserCredentialsResponse

	creds, err := g.Storage.GetUserCredentials(ctx, ctx.Value(interceptors.UserID).(int64))
	if err != nil {
		logger.Log.Error("error get user credentials from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get user credentials from DB")
	}

	protoCreds := make([]*proto.Credentials, len(creds))
	for i, cred := range creds {
		protoCreds[i] = &proto.Credentials{
			Login:       cred.Login,
			Password:    cred.Password,
			Description: cred.Description,
			Id:          cred.ID,
		}
	}
	response.Credentials = protoCreds
	return &response, nil
}
