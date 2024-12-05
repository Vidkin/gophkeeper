package handlers

import (
	"context"

	"github.com/Vidkin/gophkeeper/proto"
)

// Echo returns the received message back to the caller.
//
// Parameters:
//   - ctx: The context for the gRPC call (not used in this method).
//   - in: A pointer to the proto.EchoRequest structure containing the message to echo.
//
// Returns:
//   - A pointer to the proto.EchoResponse containing the echoed message.
//   - An error if the operation fails (this method does not generate errors).
//
// The function simply takes the input message from the request and returns it in the response.
func (g *GophkeeperServer) Echo(_ context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {
	var response proto.EchoResponse
	response.Message = in.Message
	return &response, nil
}
