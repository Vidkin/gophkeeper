package handlers

import (
	"github.com/minio/minio-go/v7"

	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/proto"
)

type GophkeeperServer struct {
	proto.UnimplementedGophkeeperServer
	Storage     *storage.PostgresStorage // Repository for storing data
	Minio       *minio.Client            // Client to minio storage
	DatabaseKey string                   // Hash key
	JWTKey      string                   // JWT secret key
	RetryCount  int                      // Number of retry attempts for database operations
}
