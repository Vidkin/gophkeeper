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

// GetFiles retrieves all files associated with the user.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - _: A pointer to the proto.GetFilesRequest structure (not used in this method).
//
// Returns:
//   - A pointer to the proto.GetFilesResponse containing the list of files.
//   - An error if the operation fails, for example, if there is an internal error while
//     retrieving the files from the storage.
//
// The function fetches the user's files from the storage and constructs a response
// containing the file details. If an error occurs during the retrieval, it logs the error
// and returns an appropriate gRPC status code.
func (g *GophkeeperServer) GetFiles(ctx context.Context, _ *proto.GetFilesRequest) (*proto.GetFilesResponse, error) {
	var response proto.GetFilesResponse

	files, err := g.Storage.GetFiles(ctx, ctx.Value(interceptors.UserID).(int64))
	if err != nil {
		logger.Log.Error("error get files from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get files from DB")
	}

	protoFiles := make([]*proto.File, len(files))
	for i, file := range files {
		protoFiles[i] = &proto.File{}
		protoFiles[i].Id = file.ID
		protoFiles[i].FileName = file.FileName
		protoFiles[i].FileSize = file.FileSize
		protoFiles[i].Description = file.Description
		protoFiles[i].CreatedAt = file.CreatedAt
	}
	response.Files = protoFiles
	return &response, nil
}
