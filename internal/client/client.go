package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Vidkin/gophkeeper/proto"
)

func NewGophkeeperClient(serverAddress, certPath string) (proto.GophkeeperClient, *grpc.ClientConn, error) {
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
