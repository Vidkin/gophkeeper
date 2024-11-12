package storage

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Vidkin/gophkeeper/internal/logger"
)

const MinioBucketName = "gophkeeper"

func NewMinioStorage(endpoint, accessKeyID, secretAccessKey string) (*minio.Client, error) {
	ctx := context.Background()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	err = minioClient.MakeBucket(ctx, MinioBucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, MinioBucketName)
		if errBucketExists == nil && exists {
			logger.Log.Info("bucket already created")
		} else {
			return nil, err
		}
	} else {
		logger.Log.Info("successfully created bucket")
	}

	return minioClient, nil
}
