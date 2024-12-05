package handlers

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Vidkin/gophkeeper/internal/client"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

func TestBankCards(t *testing.T) {
	storage, dbName := setupTestDB(t)
	defer teardownTestDB(t, storage.Conn, dbName)

	gs := &GophkeeperServer{
		Storage:     storage,
		JWTKey:      "JWTKey",
		DatabaseKey: "",
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.ValidateToken("JWTKey")))
	proto.RegisterGophkeeperServer(s, gs)

	listen, err := GetTLSListener(
		"0.0.0.0:0",
		"../../certs/public.crt",
		"../../certs/private.key")
	require.NoError(t, err)
	go func() {
		err = s.Serve(listen)
		require.NoError(t, err)
	}()
	defer s.Stop()

	addr := listen.Addr().(*net.TCPAddr)
	viper.Set("address", fmt.Sprintf("127.0.0.1:%d", addr.Port))
	viper.Set("crypto_key_public_path", "../../certs/public.crt")
	client, conn, err := client.NewGophkeeperClient()
	require.NoError(t, err)
	defer conn.Close()

	cred := proto.Credentials{
		Login:    "login",
		Password: "password",
	}

	gs.DatabaseKey = "strongDBKey2Ks5nM2J5JaI59PPEhL1x"
	_, err = client.RegisterUser(context.Background(), &proto.RegisterUserRequest{Credentials: &cred})
	require.NoError(t, err)

	resp, err := client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &cred})
	require.NoError(t, err)

	card := &proto.BankCard{
		Number:      "1234567",
		ExpireDate:  "12.07.2024",
		Cvv:         "",
		Owner:       "owner",
		Description: "description",
	}

	md := metadata.New(map[string]string{"token": resp.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	t.Run("test add bank card: error number field is empty", func(t *testing.T) {
		_, err = client.AddBankCard(ctx, &proto.AddBankCardRequest{Card: card})
		require.ErrorContains(t, err, "you must provide: CVV, expire date, card number, card owner")
	})

	card.Cvv = "123"
	t.Run("test add bank card: ok", func(t *testing.T) {
		_, err = client.AddBankCard(ctx, &proto.AddBankCardRequest{Card: card})
		require.NoError(t, err)
	})

	t.Run("test get bank cards: ok", func(t *testing.T) {
		resp, err := client.GetBankCards(ctx, &proto.GetBankCardsRequest{})
		require.NoError(t, err)
		assert.Equal(t, card.Number, resp.Cards[0].Number)
		assert.Equal(t, card.ExpireDate, resp.Cards[0].ExpireDate)
		assert.Equal(t, card.Owner, resp.Cards[0].Owner)
		assert.Equal(t, card.Description, resp.Cards[0].Description)
		assert.Equal(t, card.Cvv, resp.Cards[0].Cvv)
	})

	t.Run("test get bank card: bad id", func(t *testing.T) {
		_, err = client.GetBankCard(ctx, &proto.GetBankCardRequest{Id: ""})
		require.ErrorContains(t, err, "missing bank card id")
	})

	t.Run("test get bank card: unknown card id", func(t *testing.T) {
		_, err = client.GetBankCard(ctx, &proto.GetBankCardRequest{Id: "435"})
		require.ErrorContains(t, err, "error get bank card from DB")
	})

	t.Run("test get bank card: ok", func(t *testing.T) {
		resp, err := client.GetBankCard(ctx, &proto.GetBankCardRequest{Id: "1"})
		require.NoError(t, err)
		assert.Equal(t, card.Number, resp.Card.Number)
		assert.Equal(t, card.ExpireDate, resp.Card.ExpireDate)
		assert.Equal(t, card.Owner, resp.Card.Owner)
		assert.Equal(t, card.Description, resp.Card.Description)
		assert.Equal(t, card.Cvv, resp.Card.Cvv)
	})

	t.Run("test remove bank card: empty id", func(t *testing.T) {
		_, err = client.RemoveBankCard(ctx, &proto.RemoveBankCardRequest{Id: ""})
		require.ErrorContains(t, err, "you must provide card id")
	})

	t.Run("test remove bank card: bad id", func(t *testing.T) {
		_, err = client.RemoveBankCard(ctx, &proto.RemoveBankCardRequest{Id: "badID"})
		require.ErrorContains(t, err, "invalid card id")
	})

	t.Run("test remove bank card: ok", func(t *testing.T) {
		_, err = client.RemoveBankCard(ctx, &proto.RemoveBankCardRequest{Id: "1"})
		require.NoError(t, err)
	})
}
