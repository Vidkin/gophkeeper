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

func (g *GophkeeperServer) AddNote(ctx context.Context, in *proto.AddNoteRequest) (*emptypb.Empty, error) {
	if in.Note.Text == "" {
		logger.Log.Error("you must provide note text")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide note text")
	}

	note := &model.Note{
		UserID:      ctx.Value(interceptors.UserID).(int64),
		Text:        in.Note.Text,
		Description: in.Note.Description,
	}

	if err := g.Storage.AddNote(ctx, note); err != nil {
		logger.Log.Error("error add note", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error add note")
	}
	return &emptypb.Empty{}, nil
}
