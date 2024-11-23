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

func (g *GophkeeperServer) GetNotes(ctx context.Context, _ *proto.GetNotesRequest) (*proto.GetNotesResponse, error) {
	var response proto.GetNotesResponse

	notes, err := g.Storage.GetNotes(ctx, ctx.Value(interceptors.UserID).(int64))
	if err != nil {
		logger.Log.Error("error get notes from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get notes from DB")
	}

	protoNotes := make([]*proto.Note, len(notes))
	for i, note := range notes {
		protoNotes[i] = &proto.Note{
			Text:        note.Text,
			Description: note.Description,
			Id:          note.ID,
		}
	}
	response.Notes = protoNotes
	return &response, nil
}
