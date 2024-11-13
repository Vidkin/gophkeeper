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
