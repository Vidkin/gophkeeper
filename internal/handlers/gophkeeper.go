package handlers

import (
	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/proto"
)

// GophkeeperServer struct implements the methods defined in the proto.GophkeeperServer interface and serves as the
// main entry point for handling file-related operations.
type GophkeeperServer struct {
	proto.UnimplementedGophkeeperServer
	Storage     *storage.PostgresStorage     // Repository for storing data
	Minio       storage.MinioClientInterface // Client to minio storage
	DatabaseKey string                       // Hash key
	JWTKey      string                       // JWT secret key
	RetryCount  int                          // Number of retry attempts for database operations
}
