package handlers

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/proto"
)

// GetBankCard retrieves a specific bank card by its ID.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.GetBankCardRequest structure containing the ID of the bank card.
//
// Returns:
//   - A pointer to the proto.GetBankCardResponse containing the details of the requested bank card.
//   - An error if the operation fails, for example, if the provided ID is invalid or if there is
//     an internal error while retrieving the card from the storage.
//
// The function converts the card ID from a string to an integer and fetches the corresponding
// bank card from the storage. If an error occurs during the retrieval, it logs the error
// and returns an appropriate gRPC status code.
func (g *GophkeeperServer) GetBankCard(ctx context.Context, in *proto.GetBankCardRequest) (*proto.GetBankCardResponse, error) {
	cardID, err := strconv.Atoi(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing bank card id")
	}
	var response proto.GetBankCardResponse

	card, err := g.Storage.GetBankCard(ctx, int64(cardID))
	if err != nil {
		logger.Log.Error("error get bank card from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get bank card from DB")
	}

	protoCard := &proto.BankCard{
		Owner:       card.Owner,
		Number:      card.Number,
		ExpireDate:  card.ExpireDate,
		Cvv:         card.CVV,
		Description: card.Description,
		Id:          card.ID,
	}
	response.Card = protoCard
	return &response, nil
}
