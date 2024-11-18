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

func (g *GophkeeperServer) GetUserCredential(ctx context.Context, in *proto.GetUserCredentialRequest) (*proto.GetUserCredentialResponse, error) {
	credID, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing credentials id")
	}
	var response proto.GetUserCredentialResponse

	cred, err := g.Storage.GetUserCredential(ctx, int64(credID))
	if err != nil {
		logger.Log.Error("error get bank card from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get bank card from DB")
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
