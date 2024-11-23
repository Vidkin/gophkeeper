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

func (g *GophkeeperServer) GetNote(ctx context.Context, in *proto.GetNoteRequest) (*proto.GetNoteResponse, error) {
	noteID, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing note id")
	}
	var response proto.GetNoteResponse

	note, err := g.Storage.GetNote(ctx, int64(noteID))
	if err != nil {
		logger.Log.Error("error get note from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get note from DB")
	}

	protoNote := &proto.Note{
		Text:        note.Text,
		Description: note.Description,
		Id:          note.ID,
	}
	response.Note = protoNote
	return &response, nil
}
