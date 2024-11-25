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

// GetNote retrieves a specific note associated with the user from the storage.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.GetNoteRequest structure, which contains the ID of the note to be retrieved.
//
// Returns:
//   - A pointer to the proto.GetNoteResponse containing the requested note.
//   - An error if the operation fails, for example, if the note ID is invalid or if there is an internal
//     error while retrieving the note from the database.
//
// The function first attempts to convert the note ID from a string to an integer. If the conversion fails,
// it returns an InvalidArgument error. It then retrieves the note from the storage using the note ID.
// If an error occurs during the retrieval, it logs the error and returns an Internal status. If the
// operation is successful, it constructs a response containing the note and returns it.
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
