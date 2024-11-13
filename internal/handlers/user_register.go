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
