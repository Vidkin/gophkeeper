package app

import (
	"google.golang.org/grpc"

	"github.com/Vidkin/gophkeeper/internal/config"
	"github.com/Vidkin/gophkeeper/internal/logger"
)

type ServerApp struct {
	config     *config.ServerConfig
	gRPCServer *grpc.Server
	//repository router.Repository
}

func NewServerApp(cfg *config.ServerConfig) (*ServerApp, error) {
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}

	serverApp := &ServerApp{
		config: cfg,
		//repository: repo,
	}

	return serverApp, nil
}
