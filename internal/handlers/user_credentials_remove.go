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
