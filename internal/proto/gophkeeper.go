package proto

import (
	"context"

	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/proto"
)

type GophkeeperServer struct {
	proto.UnimplementedGophkeeperServer
	Storage    *storage.PostgresStorage // Repository for storing data
	RetryCount int                      // Number of retry attempts for database operations
}

func (g *GophkeeperServer) RegisterUser(ctx context.Context, in *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	var response proto.RegisterUserResponse
	response.Id = "test"
	return &response, nil
}

func (g *GophkeeperServer) Echo(ctx context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {
	var response proto.EchoResponse
	response.Message = in.Message
	return &response, nil
}
