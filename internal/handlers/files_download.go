package handlers

import (
	"fmt"
	"io"

	"github.com/golang-jwt/jwt/v4"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"golang.org/x/text/unicode/norm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	jwtPKG "github.com/Vidkin/gophkeeper/pkg/jwt"
	"github.com/Vidkin/gophkeeper/proto"
)

func (g *GophkeeperServer) Download(in *proto.FileDownloadRequest, srv proto.Gophkeeper_DownloadServer) error {
	var tokenString string
	if md, ok := metadata.FromIncomingContext(srv.Context()); ok {
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

	fileName := norm.NFC.String(in.FileName)
	if fileName == "" {
		logger.Log.Error("file name is required")
		return status.Error(codes.InvalidArgument, "file name is required")
	}

	fileInfo, err := g.Storage.GetFile(srv.Context(), fileName)
	if err != nil {
		logger.Log.Error("error getting file info", zap.Error(err))
		return status.Error(codes.Internal, "error getting file info")
	}

	object, err := g.Minio.GetObject(srv.Context(), fileInfo.BucketName, fileInfo.FileName, minio.GetObjectOptions{})
	if err != nil {
		logger.Log.Error("error getting object from MinIO", zap.Error(err))
		return status.Error(codes.Internal, "error getting object from MinIO")
	}

	var totalSize int64
	chunk := make([]byte, 1024)
	for {
		clear(chunk)
		n, err := object.Read(chunk)
		if err == io.EOF {
			if n != 0 {
				totalSize += int64(n)

				resp := proto.FileDownloadResponse{
					Chunk:       chunk[:n],
					FileSize:    int64(n),
					Filename:    fileInfo.FileName,
					Description: fileInfo.Description,
				}

				if err = srv.Send(&resp); err != nil {
					logger.Log.Error("error send chunk", zap.Error(err))
					return err
				}
			}
			break
		}
		if err != nil {
			logger.Log.Error("error reading object", zap.Error(err))
			return status.Error(codes.Internal, "error reading object")
		}
		totalSize += int64(n)

		resp := proto.FileDownloadResponse{
			Chunk:       chunk[:n],
			FileSize:    int64(n),
			Filename:    fileInfo.FileName,
			Description: fileInfo.Description,
		}

		if err = srv.Send(&resp); err != nil {
			logger.Log.Error("error send chunk", zap.Error(err))
			return err
		}
	}
	return nil
}
