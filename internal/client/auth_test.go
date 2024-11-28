package client

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/Vidkin/gophkeeper/app"
	"github.com/Vidkin/gophkeeper/internal/handlers"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

func TestAuth(t *testing.T) {
	storage, dbName := setupTestDB(t)
	defer teardownTestDB(t, storage.Conn, dbName)

	gs := &handlers.GophkeeperServer{
		Storage:     storage,
		JWTKey:      "JWTKey",
		DatabaseKey: "strongDBKey2Ks5nM2J5JaI59PPEhL1x",
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor,
			interceptors.HashInterceptor("defaultHashKey"),
			interceptors.ValidateToken("JWTKey")))
	proto.RegisterGophkeeperServer(s, gs)

	listen, err := app.GetTLSListener(
		"127.0.0.1:8080",
		"../../certs/public.crt",
		"../../certs/private.key")
	require.NoError(t, err)
	go func() {
		err = s.Serve(listen)
		require.NoError(t, err)
	}()
	defer s.Stop()

	viper.Set("address", "127.0.0.1:8080")
	viper.Set("crypto_key_public_path", "../../certs/public.crt")
	viper.Set("hash_key", "defaultHashKey")

	err = Register("test_login", "test_pass")
	require.NoError(t, err)

	t.Run("test bad login", func(t *testing.T) {
		err = Auth("test_loginBad", "test_pass")
		require.ErrorContains(t, err, "invalid user login or password")
	})

	t.Run("test bad password", func(t *testing.T) {
		err = Auth("test_login", "test_passBad")
		require.ErrorContains(t, err, "invalid user login or password")
	})

	t.Run("test auth ok", func(t *testing.T) {
		err = Auth("test_login", "test_pass")
		require.NoError(t, err)
	})

	tokenFilePath := path.Join(os.TempDir(), TokenFileName)
	defer os.Remove(tokenFilePath)
	t.Run("check token file exists", func(t *testing.T) {
		content, err := os.ReadFile(tokenFilePath)
		require.NoError(t, err)
		assert.NotEmpty(t, string(content))
	})
}
