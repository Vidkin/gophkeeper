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

func (g *GophkeeperServer) AddUserCredentials(ctx context.Context, in *proto.AddUserCredentialsRequest) (*emptypb.Empty, error) {
	if in.Credentials.Login == "" || in.Credentials.Password == "" {
		logger.Log.Error("you should provide: login and password")
		return nil, status.Errorf(codes.InvalidArgument, "you should provide: login and password")
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
