package app

import (
	"google.golang.org/grpc"

	"github.com/Vidkin/gophkeeper/internal/config"
	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/storage"
)

type ServerApp struct {
	config     *config.ServerConfig
	gRPCServer *grpc.Server
	storage    *storage.PostgresStorage
}

func NewServerApp(cfg *config.ServerConfig) (*ServerApp, error) {
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}

	repo, err := storage.NewPostgresStorage(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	serverApp := &ServerApp{
		config:     cfg,
		storage:    repo,
		gRPCServer: grpc.NewServer(),
	}

	return serverApp, nil
}
