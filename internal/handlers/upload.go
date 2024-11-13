package handlers

import (
	"fmt"
	"io"
	"os"

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
	var fileName, contentType string
	var fileSize uint64

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

	if err != nil {
		logger.Log.Error("error parse claims", zap.Error(err))
		return status.Errorf(codes.PermissionDenied, "error parse claims")
	}

	if !token.Valid {
		logger.Log.Error("token is not valid")
		return status.Errorf(codes.PermissionDenied, "token is not valid")
	}

	f, err := os.CreateTemp("", "*")
	if err != nil {
		logger.Log.Error("error creating temp file", zap.Error(err))
		return status.Errorf(codes.Internal, "error creating temp file")
	}
	defer func(name string) {
		err = f.Close()
		if err != nil {
			logger.Log.Error("error close temp file", zap.Error(err))
		}
		err = os.Remove(name)
		if err != nil {
			logger.Log.Error("error removing temp file", zap.Error(err))
		}
	}(f.Name())

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error("error receive file", zap.Error(err))
			return status.Errorf(codes.Internal, "error receive file")
		}

		if fileName == "" {
			fileName = req.FileName
			if fileName == "" {
				logger.Log.Error("need to provide filaname", zap.Error(err))
				return status.Errorf(codes.InvalidArgument, "need to provide filaname")
			}
		}

		if contentType == "" {
			contentType = req.ContentType
			if contentType == "" {
				logger.Log.Error("need to provide content-type", zap.Error(err))
				return status.Errorf(codes.InvalidArgument, "need to provide content-type")
			}
		}

		chunk := req.GetChunk()
		fileSize += uint64(len(chunk))

		if _, err = f.Write(chunk); err != nil {
			logger.Log.Error("error writing chunk to temp file", zap.Error(err))
			return status.Errorf(codes.Internal, "error writing chunk to temp file")
		}
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		logger.Log.Error("error seek temp file to start", zap.Error(err))
		return status.Errorf(codes.Internal, "error seek temp file to start")
	}

	_, err = g.Minio.PutObject(stream.Context(), storage.MinioBucketName, fileName, f, int64(fileSize), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		logger.Log.Error("failed to upload file to MinIO", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to upload file to MinIO")
	}

	err = g.Storage.AddFile(stream.Context(), storage.MinioBucketName, fileName, contentType, claims.UserID, fileSize)
	if err != nil {
		if errRm := g.Minio.RemoveObject(stream.Context(), storage.MinioBucketName, fileName, minio.RemoveObjectOptions{ForceDelete: true}); errRm != nil {
			logger.Log.Error("failed to remove file from MinIO", zap.Error(err))
		}
		logger.Log.Error("failed to sava file info to database", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to upload file to MinIO")
	}

	return stream.SendAndClose(&proto.FileUploadResponse{
		FileName: fileName,
		Size:     fileSize,
	})
}
