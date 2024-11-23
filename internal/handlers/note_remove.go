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

func (g *GophkeeperServer) RemoveNote(ctx context.Context, in *proto.RemoveNoteRequest) (*emptypb.Empty, error) {
	if in.Id == "" {
		logger.Log.Error("you must provide note id")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide note id")
	}

	noteID, err := strconv.ParseInt(in.Id, 10, 64)
	if err != nil {
		logger.Log.Error("invalid note id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid note id")
	}

	if err = g.Storage.RemoveNote(ctx, noteID); err != nil {
		logger.Log.Error("error remove note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove note")
	}
	return &emptypb.Empty{}, nil
}
