package proto

import (
	"context"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/model"
	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/pkg/aes"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/pkg/jwt"
	"github.com/Vidkin/gophkeeper/proto"
)

type GophkeeperServer struct {
	proto.UnimplementedGophkeeperServer
	Storage     *storage.PostgresStorage // Repository for storing data
	Minio       *minio.Client            // Client to minio storage
	DatabaseKey string                   // Hash key
	JWTKey      string                   // JWT secret key
	RetryCount  int                      // Number of retry attempts for database operations
}

func (g *GophkeeperServer) RegisterUser(ctx context.Context, in *proto.RegisterUserRequest) (*emptypb.Empty, error) {
	if in.User.Login == "" {
		logger.Log.Error("invalid user login")
		return nil, status.Errorf(codes.InvalidArgument, "invalid user login")
	}
	if in.User.Password == "" {
		logger.Log.Error("invalid user password")
		return nil, status.Errorf(codes.InvalidArgument, "invalid user password")
	}

	_, err := g.Storage.GetUser(ctx, in.User.Login)
	if err == nil {
		logger.Log.Error("user already exists")
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	encPwd, err := aes.Encrypt(g.DatabaseKey, in.User.Password)
	if err != nil {
		logger.Log.Error("error encrypt password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt password")
	}

	if err := g.Storage.AddUser(ctx, in.User.Login, encPwd); err != nil {
		logger.Log.Error("error create user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error create user")
	}
	return &emptypb.Empty{}, nil
}

func (g *GophkeeperServer) Authorize(ctx context.Context, in *proto.AuthorizeRequest) (*proto.AuthorizeResponse, error) {
	var response proto.AuthorizeResponse
	if in.User.Login == "" || in.User.Password == "" {
		logger.Log.Error("invalid user login or password")
		return nil, status.Errorf(codes.PermissionDenied, "invalid user login or password")
	}

	u, err := g.Storage.GetUser(ctx, in.User.Login)
	if err != nil {
		logger.Log.Error("error get user from db", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error get user from db")
	}

	decPwd, err := aes.Decrypt(g.DatabaseKey, u.Password)
	if err != nil {
		logger.Log.Error("error encrypt password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt password")
	}

	if in.User.Password != decPwd {
		logger.Log.Error("invalid user login or password", zap.Error(err))
		return nil, status.Errorf(codes.PermissionDenied, "invalid user login or password")
	}

	token, err := jwt.BuildJWTString(g.JWTKey, u.ID)
	if err != nil {
		logger.Log.Error("error build jwt string", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error build jwt string")
	}

	response.Token = token
	return &response, nil
}

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

	expireDate, err := aes.Encrypt(g.DatabaseKey, in.Card.ExpireDate)
	if err != nil {
		logger.Log.Error("error encrypt data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error encrypt data")
	}

	card := &model.BankCard{
		UserID:     ctx.Value(interceptors.UserID).(int64),
		CVV:        cvv,
		Owner:      owner,
		Number:     number,
		ExpireDate: expireDate,
	}

	if err := g.Storage.AddCard(ctx, card); err != nil {
		logger.Log.Error("error add bank card", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "error add bank card")
	}
	return &emptypb.Empty{}, nil
}

func (g *GophkeeperServer) GetBankCards(ctx context.Context, in *proto.GetBankCardsRequest) (*proto.GetBankCardsResponse, error) {
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

func (g *GophkeeperServer) Echo(_ context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {
	var response proto.EchoResponse
	response.Message = in.Message
	return &response, nil
}
