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
		protoFiles[i].ContentType = file.ContentType
		protoFiles[i].Description = file.Description
		protoFiles[i].CreatedAt = file.CreatedAt
	}
	response.Files = protoFiles
	return &response, nil
}
