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

func (g *GophkeeperServer) RemoveBankCard(ctx context.Context, in *proto.RemoveBankCardRequest) (*emptypb.Empty, error) {
	if in.Id == "" {
		logger.Log.Error("you should provide card id")
		return nil, status.Errorf(codes.InvalidArgument, "you should provide card id")
	}

	cardID, err := strconv.ParseInt(in.Id, 10, 64)
	if err != nil {
		logger.Log.Error("invalid card id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid card id")
	}

	if err = g.Storage.RemoveBankCard(ctx, cardID); err != nil {
		logger.Log.Error("error remove bank card", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove bank card")
	}
	return &emptypb.Empty{}, nil
}
