package proto

import (
	"context"
	"encoding/base64"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/pkg/hash"
	"github.com/Vidkin/gophkeeper/pkg/jwt"
	"github.com/Vidkin/gophkeeper/proto"
)

type GophkeeperServer struct {
	proto.UnimplementedGophkeeperServer
	Storage     *storage.PostgresStorage // Repository for storing data
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

	pHash := hash.GetHashSHA256(g.DatabaseKey, []byte(in.User.Password))
	pHashEncoded := base64.StdEncoding.EncodeToString(pHash)

	if err := g.Storage.AddUser(ctx, in.User.Login, pHashEncoded); err != nil {
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

	pHash := hash.GetHashSHA256(g.DatabaseKey, []byte(in.User.Password))
	pHashEncoded := base64.StdEncoding.EncodeToString(pHash)

	if pHashEncoded != u.Password {
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

func (g *GophkeeperServer) Echo(_ context.Context, in *proto.EchoRequest) (*proto.EchoResponse, error) {
	var response proto.EchoResponse
	response.Message = in.Message
	return &response, nil
}
