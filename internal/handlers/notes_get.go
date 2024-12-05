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

// GetNotes retrieves all notes associated with the user from the storage.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - _: A pointer to the proto.GetNotesRequest structure (not used in this method).
//
// Returns:
//   - A pointer to the proto.GetNotesResponse containing the list of notes associated with the user.
//   - An error if the operation fails, for example, if there is an internal error while retrieving the
//     notes from the database.
//
// The function fetches the user's notes from the storage using the user ID extracted from the context.
// If an error occurs during the retrieval, it logs the error and returns an Internal status. If the
// operation is successful, it constructs a response containing the notes and returns it.
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
