package storage

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMinioClient is a mock of the minio.Client
type MockMinioClient struct {
	mock.Mock
}

func (m *MockMinioClient) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	args := m.Called(ctx, bucketName, opts)
	return args.Error(0)
}

func (m *MockMinioClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	args := m.Called(ctx, bucketName)
	return args.Bool(0), args.Error(1)
}

func (m *MockMinioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Get(0).(*minio.Object), args.Error(1)
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *MockMinioClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Error(0)
}

// TestNewMinioStorage tests the NewMinioStorage function
func TestNewMinioStorage(t *testing.T) {
	tests := []struct {
		name                string
		endpoint            string
		accessKeyID         string
		secretAccessKey     string
		mockMakeBucketErr   error
		mockBucketExists    bool
		mockBucketExistsErr error
		expectErr           bool
	}{
		{
			name:              "Successful Bucket Creation",
			endpoint:          "localhost:9000",
			accessKeyID:       "testAccessKey",
			secretAccessKey:   "testSecretKey",
			mockMakeBucketErr: nil,
			mockBucketExists:  false,
			expectErr:         false,
		},
		{
			name:                "Bucket Already Exists",
			endpoint:            "localhost:9000",
			accessKeyID:         "testAccessKey",
			secretAccessKey:     "testSecretKey",
			mockMakeBucketErr:   errors.New("bucket already exists"),
			mockBucketExists:    true,
			mockBucketExistsErr: nil,
			expectErr:           false,
		},
		{
			name:                "Error Creating Bucket",
			endpoint:            "localhost:9000",
			accessKeyID:         "testAccessKey",
			secretAccessKey:     "testSecretKey",
			mockMakeBucketErr:   errors.New("failed to create bucket"),
			mockBucketExists:    false,
			mockBucketExistsErr: nil,
			expectErr:           true,
		},
		{
			name:                "Error Checking Bucket Existence",
			endpoint:            "localhost:9000",
			accessKeyID:         "testAccessKey",
			secretAccessKey:     "testSecretKey",
			mockMakeBucketErr:   errors.New("failed to create bucket"),
			mockBucketExists:    false,
			mockBucketExistsErr: errors.New("failed to check bucket existence"),
			expectErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMinioClient)

			// Mock the MakeBucket method
			mockClient.On("MakeBucket", mock.Anything, MinioBucketName, mock.Anything).Return(tt.mockMakeBucketErr)

			// Mock the BucketExists method
			if tt.mockBucketExists || tt.name == "Error Creating Bucket" || tt.name == "Error Checking Bucket Existence" {
				mockClient.On("BucketExists", mock.Anything, MinioBucketName).Return(tt.mockBucketExists, tt.mockBucketExistsErr)
			}

			// Call NewMinioStorage with the mock client
			minioClient, err := NewMinioStorage(tt.endpoint, tt.accessKeyID, tt.secretAccessKey, mockClient)

			if (err != nil) != tt.expectErr {
				t.Errorf("NewMinioStorage() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr {
				assert.NotNil(t, minioClient)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
