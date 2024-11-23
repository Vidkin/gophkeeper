package handlers

import (
	"fmt"
	"io"

	"github.com/golang-jwt/jwt/v4"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/storage"
	jwtPKG "github.com/Vidkin/gophkeeper/pkg/jwt"
	"github.com/Vidkin/gophkeeper/proto"
)

func (g *GophkeeperServer) Upload(stream proto.Gophkeeper_UploadServer) error {
	var fileName, description string
	var fileSize int64
	var tokenString string

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		values := md.Get("token")
		if len(values) > 0 {
			tokenString = values[0]
		}
	}
	if len(tokenString) == 0 {
		return status.Error(codes.PermissionDenied, "missing token")
	}

	claims := &jwtPKG.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Log.Error("unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(g.JWTKey), nil
		})

	if err != nil || !token.Valid {
		logger.Log.Error("error parse claims", zap.Error(err))
		return status.Errorf(codes.PermissionDenied, "error parse claims")
	}

	pr, pw := io.Pipe()
	defer func(pw *io.PipeWriter) {
		err = pw.Close()
		if err != nil {
			logger.Log.Error("error closing pipe writer", zap.Error(err))
		}
	}(pw)
	defer func(pr *io.PipeWriter) {
		err = pr.Close()
		if err != nil {
			logger.Log.Error("error closing pipe reader", zap.Error(err))
		}
	}(pw)

	req, err := stream.Recv()
	if err == io.EOF {
		logger.Log.Error("empty file")
		return status.Errorf(codes.FailedPrecondition, "empty file")
	}
	if err != nil {
		logger.Log.Error("error receive file", zap.Error(err))
		return status.Errorf(codes.Internal, "error receive file")
	}

	fileName = req.FileName
	description = req.Description
	fileSize = req.FileSize
	if fileName == "" || fileSize == 0 {
		return status.Errorf(codes.InvalidArgument, "filename, file-size are required")
	}
	chunk := req.GetChunk()
	resultChan := make(chan error, 1)
	go func() {
		for {
			if chunk != nil {
				if _, err = pw.Write(chunk); err != nil {
					logger.Log.Error("error writing chunk to pipe", zap.Error(err))
					resultChan <- err
					return
				}
			}

			req, err = stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Log.Error("error receive file", zap.Error(err))
				resultChan <- status.Errorf(codes.Internal, "error receive file")
			}

			chunk = req.GetChunk()
			if _, err = pw.Write(chunk); err != nil {
				logger.Log.Error("error writing chunk to pipe", zap.Error(err))
				resultChan <- err
				return
			}
			chunk = nil
		}
		resultChan <- nil
	}()

	_, err = g.Minio.PutObject(stream.Context(), storage.MinioBucketName, fileName, pr, fileSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		logger.Log.Error("failed to upload file to MinIO", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to upload file to MinIO")
	}

	err = <-resultChan
	if err != nil {
		if errRm := g.Minio.RemoveObject(stream.Context(), storage.MinioBucketName, fileName, minio.RemoveObjectOptions{ForceDelete: true}); errRm != nil {
			logger.Log.Error("failed to remove file from MinIO", zap.Error(err))
		}
		return status.Errorf(codes.Internal, "failed to upload file to MinIO")
	}

	err = g.Storage.AddFile(stream.Context(), storage.MinioBucketName, fileName, description, claims.UserID, fileSize)
	if err != nil {
		if errRm := g.Minio.RemoveObject(stream.Context(), storage.MinioBucketName, fileName, minio.RemoveObjectOptions{ForceDelete: true}); errRm != nil {
			logger.Log.Error("failed to remove file from MinIO", zap.Error(err))
		}
		logger.Log.Error("failed to save file info to database", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to upload file to MinIO")
	}

	return stream.SendAndClose(&proto.FileUploadResponse{
		FileName: fileName,
		FileSize: fileSize,
	})
}
