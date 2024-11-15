package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/model"
	"github.com/Vidkin/gophkeeper/pkg/aes"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

func (g *GophkeeperServer) AddBankCard(ctx context.Context, in *proto.AddBankCardRequest) (*emptypb.Empty, error) {
	if in.Card.Cvv == "" || in.Card.ExpireDate == "" || in.Card.Number == "" || in.Card.Owner == "" {
		logger.Log.Error("you should provide: CVV, expire date, card number, card owner")
		return nil, status.Errorf(codes.InvalidArgument, "you should provide: CVV, expire date, card number, card owner")
	}

	cvv, err := aes.Encrypt(g.DatabaseKey, in.Card.Cvv)
	if err != nil {
		logger.Log.Error("error encrypt data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt data")
	}

	owner, err := aes.Encrypt(g.DatabaseKey, in.Card.Owner)
	if err != nil {
		logger.Log.Error("error encrypt data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt data")
	}

	number, err := aes.Encrypt(g.DatabaseKey, in.Card.Number)
	if err != nil {
		logger.Log.Error("error encrypt data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt data")
	}

	description, err := aes.Encrypt(g.DatabaseKey, in.Card.Description)
	if err != nil {
		logger.Log.Error("error encrypt data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt data")
	}

	expireDate, err := aes.Encrypt(g.DatabaseKey, in.Card.ExpireDate)
	if err != nil {
		logger.Log.Error("error encrypt data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt data")
	}

	card := &model.BankCard{
		UserID:      ctx.Value(interceptors.UserID).(int64),
		CVV:         cvv,
		Owner:       owner,
		Number:      number,
		ExpireDate:  expireDate,
		Description: description,
	}

	if err := g.Storage.AddCard(ctx, card); err != nil {
		logger.Log.Error("error add bank card", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error add bank card")
	}
	return &emptypb.Empty{}, nil
}
