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

func (g *GophkeeperServer) Authorize(ctx context.Context, in *proto.AuthorizeRequest) (*proto.AuthorizeResponse, error) {
	var response proto.AuthorizeResponse
	if in.Credentials.Login == "" || in.Credentials.Password == "" {
		logger.Log.Error("invalid user login or password")
		return nil, status.Errorf(codes.PermissionDenied, "invalid user login or password")
	}

	u, err := g.Storage.GetUser(ctx, in.Credentials.Login)
	if err != nil {
		logger.Log.Error("error get user from db", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get user from db")
	}

	decPwd, err := aes.Decrypt(g.DatabaseKey, u.Password)
	if err != nil {
		logger.Log.Error("error encrypt password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt password")
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
