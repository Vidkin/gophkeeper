package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/pkg/aes"
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
		protoCards[i] = &proto.BankCard{}

		owner, err := aes.Decrypt(g.DatabaseKey, card.Owner)
		if err != nil {
			logger.Log.Error("error decrypt data", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "error decrypt data")
		}
		protoCards[i].Owner = owner

		number, err := aes.Decrypt(g.DatabaseKey, card.Number)
		if err != nil {
			logger.Log.Error("error decrypt data", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "error decrypt data")
		}
		protoCards[i].Number = number

		expireDate, err := aes.Decrypt(g.DatabaseKey, card.ExpireDate)
		if err != nil {
			logger.Log.Error("error decrypt data", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "error decrypt data")
		}
		protoCards[i].ExpireDate = expireDate

		cvv, err := aes.Decrypt(g.DatabaseKey, card.CVV)
		if err != nil {
			logger.Log.Error("error decrypt data", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "error decrypt data")
		}
		protoCards[i].Cvv = cvv
		protoCards[i].Id = card.ID
	}
	response.Cards = protoCards
	return &response, nil
}
