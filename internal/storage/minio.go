// Package storage provides functionality for interacting with storage services.
//
// This package includes the NewMinioStorage function, which initializes a MinIO client for object storage.
package storage

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Vidkin/gophkeeper/internal/logger"
)

// MinioBucketName contains default MinIO bucket name for gophkeeper app
const MinioBucketName = "gophkeeper"

// MinioClientInterface is the interface for MinIO client and mock-client
type MinioClientInterface interface {
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
}

// NewMinioStorage initializes a new MinIO client and creates a bucket if it does not already exist.
//
// Parameters:
//   - endpoint: A string representing the MinIO server endpoint (e.g., "localhost:9000").
//   - accessKeyID: A string representing the access key ID for MinIO authentication.
//   - secretAccessKey: A string representing the secret access key for MinIO authentication.
//
// Returns:
//   - A pointer to a minio.Client instance for interacting with the MinIO storage.
//   - An error if the client could not be created or if there was an issue creating the bucket.
//
// The function creates a new MinIO client with the provided credentials and a secure connection.
// It then attempts to create a bucket with the name defined by MinioBucketName. If the bucket
// already exists, it logs an informational message. If the bucket creation fails for any other
// reason, it returns an error. If successful, it returns the MinIO client instance.
func NewMinioStorage(endpoint, accessKeyID, secretAccessKey string, minioClient MinioClientInterface) (MinioClientInterface, error) {
	ctx := context.Background()
	if minioClient == nil {
		var err error
		minioClient, err = minio.New(endpoint, &minio.Options{
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
	}

	err := minioClient.MakeBucket(ctx, MinioBucketName, minio.MakeBucketOptions{})
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
