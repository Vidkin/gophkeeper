// Package app provides an implementation of a GophKeeper gRPC server application.
// It encapsulates all necessary components for running the server, including configuration,
// gRPC server instance, listener, and storage. The package handles signal management to
// gracefully stop the server and ensures proper resource cleanup.
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

// ServerApp represents the gRPC server application, encapsulating all necessary components
// for running the server, including configuration, gRPC server instance, listener, and storage.
//
// This struct is responsible for initializing and managing the lifecycle of the gRPC server,
// handling signals to gracefully stop the server, and managing storage connections.
type ServerApp struct {
	config     *config.ServerConfig
	gRPCServer *grpc.Server
	listener   net.Listener
	storage    *storage.PostgresStorage
}

// NewServerApp creates and returns a new instance of the ServerApp initialized with the provided
// configuration. It sets up logging, initializes storage connections (PostgreSQL and MinIO),
// configures gRPC server with interceptors for logging, hashing, and token validation,
// and prepares a TLS listener.
//
// Parameters:
//   - cfg: A pointer to the ServerConfig struct containing all necessary server configurations.
//
// Returns:
//   - A pointer to the newly created ServerApp instance.
//   - An error if any initialization step fails.
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
			interceptors.ValidateToken(cfg.JWTKey),
		),
	)
	proto.RegisterGophkeeperServer(gRPCServer, &handlers.GophkeeperServer{
		RetryCount:  cfg.RetryCount,
		Storage:     repo,
		Minio:       minioClient,
		DatabaseKey: cfg.DatabaseKey,
		JWTKey:      cfg.JWTKey,
	})
	listener, err := getTLSListener(cfg.ServerAddress.Address, cfg.CryptoKeyPublic, cfg.CryptoKeyPrivate)
	if err != nil {
		logger.Log.Fatal("failed to create TLS listener", zap.Error(err))
		return nil, err
	}
	return &ServerApp{
		config:     cfg,
		gRPCServer: gRPCServer,
		listener:   listener,
		storage:    repo,
	}, nil
}

// getTLSListener creates a TLS listener on the specified address using provided public and private keys.
//
// Parameters:
//   - addr: The network address to listen on.
//   - certFile: Path to the certificate file.
//   - keyFile: Path to the private key file.
//
// Returns:
//   - A net.Listener instance configured for TLS.
//   - An error if the listener creation fails.
func getTLSListener(addr, certFile, keyFile string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	return tls.Listen("tcp", addr, config)
}

// Run starts the gRPC server on the configured listener and handles graceful shutdown.
// It listens for OS signals to trigger a graceful stop of the server.
func (s *ServerApp) Run() {
	logger.Log.Info("running server", zap.String("address", s.config.ServerAddress.Address))
	go func() {
		if err := s.gRPCServer.Serve(s.listener); err != nil {
			logger.Log.Fatal("failed to serve", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	s.Stop()
}

// Stop gracefully shuts down the gRPC server and closes the storage connection.
func (s *ServerApp) Stop() {
	logger.Log.Info("stopping server", zap.String("address", s.config.ServerAddress.Address))
	if s.gRPCServer != nil {
		s.gRPCServer.GracefulStop()
	}
	if err := s.storage.Close(); err != nil {
		logger.Log.Error("error closing storage before exit", zap.Error(err))
	}
}
