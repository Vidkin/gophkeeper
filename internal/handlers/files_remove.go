package handlers

import (
	"context"
	"strconv"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/proto"
)

func (g *GophkeeperServer) RemoveFile(ctx context.Context, in *proto.FileRemoveRequest) (*emptypb.Empty, error) {
	if in.Id == "" {
		logger.Log.Error("you must provide file id")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide file id")
	}

	fileID, err := strconv.ParseInt(in.Id, 10, 64)
	if err != nil {
		logger.Log.Error("invalid file id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid file id")
	}

	file, err := g.Storage.GetFile(ctx, fileID)
	if err != nil {
		logger.Log.Error("file not found", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "file not found")
	}

	err = g.Minio.RemoveObject(ctx, storage.MinioBucketName, file.FileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		logger.Log.Error("error remove file from minio", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove file from minio")
	}

	if err = g.Storage.RemoveFile(ctx, fileID); err != nil {
		logger.Log.Error("error remove file from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove file from DB")
	}
	return &emptypb.Empty{}, nil
}
