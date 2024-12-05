package handlers

import (
	"context"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/proto"
)

// RemoveFile removes a file associated with the user from both the storage and MinIO.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.FileRemoveRequest structure, which contains the name of the file to be removed.
//
// Returns:
//   - A pointer to an empty proto.Empty response indicating successful removal of the file.
//   - An error if the operation fails, for example, if the file name is not provided, if the file is not found,
//     or if there is an internal error while removing the file from MinIO or the database.
//
// The function first checks if the file name is provided in the request. If not, it logs an error and returns
// an InvalidArgument status. It then attempts to retrieve the file from the storage. If the file is not found,
// it logs the error and returns a NotFound status. If the file is found, it proceeds to remove the file from
// MinIO. If an error occurs during the removal from MinIO, it logs the error and returns an Internal status.
// Finally, it attempts to remove the file from the database, logging any errors that occur and returning an
// Internal status if the operation fails. If all operations are successful, it returns an empty response.
func (g *GophkeeperServer) RemoveFile(ctx context.Context, in *proto.FileRemoveRequest) (*emptypb.Empty, error) {
	if in.FileName == "" {
		logger.Log.Error("you must provide file name")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide file name")
	}

	file, err := g.Storage.GetFile(ctx, in.FileName)
	if err != nil {
		logger.Log.Error("file not found", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "file not found")
	}

	err = g.Minio.RemoveObject(ctx, storage.MinioBucketName, file.FileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		logger.Log.Error("error remove file from minio", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove file from minio")
	}

	if err = g.Storage.RemoveFile(ctx, in.FileName); err != nil {
		logger.Log.Error("error remove file from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove file from DB")
	}
	return &emptypb.Empty{}, nil
}
