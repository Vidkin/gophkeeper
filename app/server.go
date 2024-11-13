package app

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Vidkin/gophkeeper/internal/config"
	"github.com/Vidkin/gophkeeper/internal/handlers"
	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

type ServerApp struct {
	config     *config.ServerConfig
	gRPCServer *grpc.Server
	listener   net.Listener
	storage    *storage.PostgresStorage
}

func NewServerApp(cfg *config.ServerConfig) (*ServerApp, error) {
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}

	repo, err := storage.NewPostgresStorage(cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Error("error init postgres storage", zap.Error(err))
		return nil, err
	}

	minioClient, err := storage.NewMinioStorage(cfg.MinioEndpoint, cfg.MinioAccessKeyID, cfg.MinioSecretAccessKey)
	if err != nil {
		logger.Log.Error("error init minio storage", zap.Error(err))
		return nil, err
	}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor,
			interceptors.HashInterceptor(cfg.Key),
			interceptors.ValidateToken(cfg.JWTKey)))
	proto.RegisterGophkeeperServer(gRPCServer, &handlers.GophkeeperServer{
		RetryCount:  cfg.RetryCount,
		Storage:     repo,
		Minio:       minioClient,
		DatabaseKey: cfg.DatabaseKey,
		JWTKey:      cfg.JWTKey,
	})

	listener, err := getTLSListener(cfg.ServerAddress.Address, cfg.CryptoKeyPublic, cfg.CryptoKeyPrivate)
	if err != nil {
		logger.Log.Error("error init listener", zap.Error(err))
		return nil, err
	}

	serverApp := &ServerApp{
		config:     cfg,
		storage:    repo,
		listener:   listener,
		gRPCServer: gRPCServer,
	}

	return serverApp, nil
}

func getTLSListener(address, publicKeyPath, privateKeyPath string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(publicKeyPath, privateKeyPath)
	if err != nil {
		logger.Log.Error("failed to load server certificate", zap.Error(err))
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	listener, err := tls.Listen("tcp", address, cfg)
	if err != nil {
		logger.Log.Error("failed to listen", zap.Error(err))
		return nil, err
	}

	return listener, nil
}

func (a *ServerApp) Run() {
	logger.Log.Info("running server", zap.String("address", a.config.ServerAddress.Address))

	go func() {
		err := a.gRPCServer.Serve(a.listener)
		if err != nil {
			logger.Log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	defer a.Stop()
}

func (a *ServerApp) Stop() {
	logger.Log.Info("stop server", zap.String("address", a.config.ServerAddress.Address))

	if a.gRPCServer != nil {
		a.gRPCServer.GracefulStop()
	}

	err := a.storage.Close()
	if err != nil {
		logger.Log.Info("error close repository before exit", zap.Error(err))
	}
}
