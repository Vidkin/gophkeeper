package handlers

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
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

func TestUserCredentials(t *testing.T) {
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

	_, err = client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &cred})
	require.NoError(t, err)

	f, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	token := string(f)
	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	cred.Login = ""
	t.Run("test add credentials: error login is empty", func(t *testing.T) {
		_, err = client.AddUserCredentials(ctx, &proto.AddUserCredentialsRequest{Credentials: &cred})
		require.ErrorContains(t, err, "you must provide: login and password")
	})

	cred.Login = "login"
	t.Run("test add credentials: ok", func(t *testing.T) {
		_, err = client.AddUserCredentials(ctx, &proto.AddUserCredentialsRequest{Credentials: &cred})
		require.NoError(t, err)
	})

	t.Run("test get credentials: ok", func(t *testing.T) {
		resp, err := client.GetUserCredentials(ctx, &proto.GetUserCredentialsRequest{})
		require.NoError(t, err)
		assert.Equal(t, cred.Login, resp.Credentials[0].Login)
		assert.Equal(t, cred.Password, resp.Credentials[0].Password)
		assert.Equal(t, cred.Description, resp.Credentials[0].Description)
	})

	t.Run("test get credential: invalid id", func(t *testing.T) {
		_, err = client.GetUserCredential(ctx, &proto.GetUserCredentialRequest{Id: "badId"})
		require.ErrorContains(t, err, "missing credentials id")
	})

	t.Run("test get credential: unknown id", func(t *testing.T) {
		_, err = client.GetUserCredential(ctx, &proto.GetUserCredentialRequest{Id: "435"})
		require.ErrorContains(t, err, "error get credentials from DB")
	})

	t.Run("test get credential: ok", func(t *testing.T) {
		resp, err := client.GetUserCredential(ctx, &proto.GetUserCredentialRequest{Id: "1"})
		require.NoError(t, err)
		assert.Equal(t, cred.Login, resp.Credentials.Login)
		assert.Equal(t, cred.Password, resp.Credentials.Password)
		assert.Equal(t, cred.Description, resp.Credentials.Description)
	})

	t.Run("test remove credentials: empty id", func(t *testing.T) {
		_, err = client.RemoveUserCredentials(ctx, &proto.RemoveUserCredentialsRequest{Id: ""})
		require.ErrorContains(t, err, "you must provide credentials id")
	})

	t.Run("test remove credentials: bad id", func(t *testing.T) {
		_, err = client.RemoveUserCredentials(ctx, &proto.RemoveUserCredentialsRequest{Id: "badID"})
		require.ErrorContains(t, err, "invalid credentials id")
	})

	t.Run("test remove credentials: ok", func(t *testing.T) {
		_, err = client.RemoveUserCredentials(ctx, &proto.RemoveUserCredentialsRequest{Id: "1"})
		require.NoError(t, err)
	})
}
