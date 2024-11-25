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

// AddNote creates a new note associated with the user and stores it in the database.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.AddNoteRequest structure, which contains the note to be added.
//
// Returns:
//   - A pointer to an empty proto.Empty response indicating successful addition of the note.
//   - An error if the operation fails, for example, if the note text is not provided or if there is an
//     internal error while adding the note to the storage.
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
