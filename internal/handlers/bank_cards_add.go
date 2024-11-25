package handlers

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/model"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

// AddBankCard adds a new bank card for the user.
//
// Parameters:
//   - ctx: The context for the gRPC call, which may contain user identification information.
//   - in: A pointer to the proto.AddBankCardRequest structure containing the bank card details.
//
// Returns:
//   - An empty response (emptypb.Empty) if the operation is successful.
//   - An error if the operation fails, for example, if required fields are missing or if
//     there is an internal error while adding the card to the storage.
//
// The function validates the input fields and logs an error if any required field is missing.
// It then creates a model.BankCard instance and attempts to add it to the storage.
func (g *GophkeeperServer) AddBankCard(ctx context.Context, in *proto.AddBankCardRequest) (*emptypb.Empty, error) {
	if in.Card.Cvv == "" || in.Card.ExpireDate == "" || in.Card.Number == "" || in.Card.Owner == "" {
		logger.Log.Error("you must provide: CVV, expire date, card number, card owner")
		return nil, status.Errorf(codes.InvalidArgument, "you must provide: CVV, expire date, card number, card owner")
	}

	card := &model.BankCard{
		UserID:      ctx.Value(interceptors.UserID).(int64),
		CVV:         in.Card.Cvv,
		Owner:       in.Card.Owner,
		Number:      in.Card.Number,
		ExpireDate:  in.Card.ExpireDate,
		Description: in.Card.Description,
	}

	if err := g.Storage.AddCard(ctx, card); err != nil {
		logger.Log.Error("error add bank card", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error add bank card")
	}
	return &emptypb.Empty{}, nil
}
