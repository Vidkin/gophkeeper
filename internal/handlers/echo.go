package handlers

import (
	"context"

	"github.com/Vidkin/gophkeeper/proto"
)

func (g *GophkeeperServer) Echo(_ context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {
	var response proto.EchoResponse
	response.Message = in.Message
	return &response, nil
}
