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
