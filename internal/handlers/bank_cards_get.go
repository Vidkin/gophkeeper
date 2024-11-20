package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

func (g *GophkeeperServer) GetBankCards(ctx context.Context, _ *proto.GetBankCardsRequest) (*proto.GetBankCardsResponse, error) {
	var response proto.GetBankCardsResponse

	cards, err := g.Storage.GetBankCards(ctx, ctx.Value(interceptors.UserID).(int64))
	if err != nil {
		logger.Log.Error("error get bank cards from DB", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get bank cards from DB")
	}

	protoCards := make([]*proto.BankCard, len(cards))
	for i, card := range cards {
		protoCards[i] = &proto.BankCard{
			Owner:       card.Owner,
			Number:      card.Number,
			ExpireDate:  card.ExpireDate,
			Cvv:         card.CVV,
			Description: card.Description,
			Id:          card.ID,
		}
	}
	response.Cards = protoCards
	return &response, nil
}
