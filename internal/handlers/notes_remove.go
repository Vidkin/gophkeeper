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

// RemoveNote deletes a note associated with the user from the storage.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.RemoveNoteRequest structure, which contains the ID of the note to be removed.
//
// Returns:
//   - A pointer to an empty proto.Empty response indicating successful removal of the note.
//   - An error if the operation fails, for example, if the note ID is not provided, if the note ID is
//     invalid, or if there is an internal error while removing the note from the storage.
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
