package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Vidkin/gophkeeper/proto"
)

// TokenFileName is the name of the temporary file used to store the JWT token.
const TokenFileName = "gophkeeperJWT.tmp"

// NewGophkeeperClient creates a new gRPC client for the GophKeeper server with secure TLS credentials.
//
// It reads the server address and the public key certificate path from the configuration,
// establishes a TLS connection using the provided CA certificate, and returns a new GophkeeperClient instance
// along with the gRPC connection.
//
// Returns:
//   - A pointer to the GophkeeperClient interface for making gRPC calls.
//   - A pointer to the grpc.ClientConn for managing the connection.
//   - An error if any step in the process fails, including reading the certificate,
//     appending it to the certificate pool, or establishing the gRPC connection.
func NewGophkeeperClient() (proto.GophkeeperClient, *grpc.ClientConn, error) {
	serverAddress := viper.GetString("address")
	certPath := viper.GetString("crypto_key_public_path")

	certPool := x509.NewCertPool()
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, nil, errors.New("error add CA cert into the pool")
	}

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		NextProtos: []string{"h2"},
	}
	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, err
	}

	return proto.NewGophkeeperClient(conn), conn, nil
}
