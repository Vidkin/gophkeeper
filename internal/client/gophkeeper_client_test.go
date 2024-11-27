package client

import (
	"context"
	"database/sql"
	"io"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Vidkin/gophkeeper/internal/storage"
)

const (
	expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzIzODczMzgsIlVzZXJJRCI6MX0.B6kBiV1YOiDZd1oxp4weHgkFtJcN5VebwWpRD70uQDw"
)

func setExpiredToken(t *testing.T) {
	err := os.Remove(path.Join(os.TempDir(), TokenFileName))
	if !os.IsNotExist(err) {
		require.NoError(t, err)
	}
	f, err := os.Create(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	defer f.Close()
	_, err = f.WriteString(expiredToken)
	require.NoError(t, err)
}

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

const letters = "abcdefghijklmnopqrstuvwxyz"

func randomDBName(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func setupTestDB(t *testing.T) (*storage.PostgresStorage, string) {
	connStr := "user=postgres password=postgres dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	tempDBName := randomDBName(10)
	_, err = db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1", tempDBName)
	require.NoError(t, err)
	db.Exec("DROP DATABASE " + tempDBName)
	_, err = db.Exec("CREATE DATABASE " + tempDBName)
	require.NoError(t, err)

	tempConnStr := "user=postgres password=postgres dbname=" + tempDBName + " host=127.0.0.1 port=5432 sslmode=disable"
	st, err := storage.NewPostgresStorage(tempConnStr)
	require.NoError(t, err)

	return st, tempDBName
}

func teardownTestDB(t *testing.T, db *sql.DB, dbName string) {
	db.Close()

	connStr := "user=postgres password=postgres dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	mainDB, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	_, err = mainDB.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1", dbName)
	require.NoError(t, err)

	_, err = mainDB.Exec("DROP DATABASE " + dbName)
	require.NoError(t, err)

	mainDB.Close()
}

//mockClient := new(MockMinioClient)
//minioClient, err := mStorage.NewMinioStorage("endpoint", "accessKeyID", "secretAccessKey", mockClient)
