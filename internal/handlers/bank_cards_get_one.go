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
