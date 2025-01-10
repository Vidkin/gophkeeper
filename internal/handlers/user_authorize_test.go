package handlers

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/Vidkin/gophkeeper/internal/client"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

func TestAuthorize(t *testing.T) {
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

	t.Run("test authorize: empty login", func(t *testing.T) {
		_, err = client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &proto.Credentials{Login: "", Password: "password"}})
		require.ErrorContains(t, err, "invalid user login or password")
	})

	t.Run("test authorize: empty password", func(t *testing.T) {
		_, err = client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &proto.Credentials{Login: "login", Password: ""}})
		require.ErrorContains(t, err, "invalid user login or password")
	})

	t.Run("test authorize: bad login", func(t *testing.T) {
		_, err = client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &proto.Credentials{Login: "badLogin", Password: "password"}})
		require.ErrorContains(t, err, "invalid user login or password")
	})

	gs.DatabaseKey = ""
	t.Run("test authorize: empty database key", func(t *testing.T) {
		_, err = client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &cred})
		require.ErrorContains(t, err, "invalid user login or password")
	})

	gs.DatabaseKey = "strongDBKey2Ks5nM2J5JaI59PPEhL1x"
	t.Run("test authorize: invalid password", func(t *testing.T) {
		_, err = client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &proto.Credentials{Login: "login", Password: "badPassword"}})
		require.ErrorContains(t, err, "invalid user login or password")
	})

	t.Run("test authorize: ok", func(t *testing.T) {
		_, err = client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &cred})
		require.NoError(t, err)
	})

}
