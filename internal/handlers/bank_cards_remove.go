package handlers

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/proto"
)

// RemoveBankCard removes a bank card associated with the user by its ID.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.RemoveBankCardRequest structure containing the ID of the bank card to remove.
//
// Returns:
//   - An empty response (emptypb.Empty) if the operation is successful.
//   - An error if the operation fails, for example, if the provided ID is missing or invalid,
//     or if there is an internal error while removing the card from the storage.
//
// The function validates the input ID, converts it to an integer, and attempts to remove the
// corresponding bank card from the storage. If an error occurs during the removal, it logs the
// error and returns an appropriate gRPC status code.
func (g *GophkeeperServer) RemoveBankCard(ctx context.Context, in *proto.RemoveBankCardRequest) (*emptypb.Empty, error) {
	if in.Id == "" {
		logger.Log.Error("you must provide card id")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide card id")
	}

	cardID, err := strconv.ParseInt(in.Id, 10, 64)
	if err != nil {
		logger.Log.Error("invalid card id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid card id")
	}

	if err = g.Storage.RemoveBankCard(ctx, cardID); err != nil {
		logger.Log.Error("error remove bank card", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error remove bank card")
	}
	return &emptypb.Empty{}, nil
}
