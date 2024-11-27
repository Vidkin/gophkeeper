package client

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/Vidkin/gophkeeper/app"
	"github.com/Vidkin/gophkeeper/internal/handlers"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

func TestCredentials(t *testing.T) {
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
		"/Users/skim/GolandProjects/gophkeeper/certs/public.crt",
		"/Users/skim/GolandProjects/gophkeeper/certs/private.key")
	require.NoError(t, err)
	go func() {
		err = s.Serve(listen)
		require.NoError(t, err)
	}()
	defer s.Stop()

	viper.Set("address", "127.0.0.1:8080")
	viper.Set("crypto_key_public_path", "/Users/skim/GolandProjects/gophkeeper/certs/public.crt")
	viper.Set("hash_key", "defaultHashKey")

	err = Register("test_login", "test_pass")
	require.NoError(t, err)

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)

	cred := proto.Credentials{
		Login:       "login",
		Password:    "password",
		Description: "description",
	}
	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test add credentials: invalid key size", func(t *testing.T) {
		err = AddCredentials(&cred)
		require.ErrorContains(t, err, "invalid key size")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	t.Run("test add credentials: missing hash", func(t *testing.T) {
		err = AddCredentials(&cred)
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test add credentials: missed token file", func(t *testing.T) {
		err = AddCredentials(&cred)
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test add credentials: expired token", func(t *testing.T) {
		err = AddCredentials(&cred)
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test add credentials: ok", func(t *testing.T) {
		err = AddCredentials(&cred)
		require.NoError(t, err)
	})

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test get all credentials: missing hash", func(t *testing.T) {
		err = GetAllCredentials()
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test get all credentials: missed token file", func(t *testing.T) {
		err = GetAllCredentials()
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test get all credentials: expired token", func(t *testing.T) {
		err = GetAllCredentials()
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test get all credentials: ok", func(t *testing.T) {
		err = GetAllCredentials()
		require.NoError(t, err)
	})

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test get credential: missing hash", func(t *testing.T) {
		err = GetCredentials(1)
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test get credential: missed token file", func(t *testing.T) {
		err = GetCredentials(1)
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test get credential: expired token", func(t *testing.T) {
		err = GetCredentials(1)
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test get credential: ok", func(t *testing.T) {
		err = GetCredentials(1)
		require.NoError(t, err)
	})

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test remove credentials: missing hash", func(t *testing.T) {
		err = RemoveCredentials(1)
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test remove credentials: missed token file", func(t *testing.T) {
		err = RemoveCredentials(1)
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test remove credentials: expired token", func(t *testing.T) {
		err = RemoveCredentials(1)
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test remove unknown credentials: ok", func(t *testing.T) {
		err = RemoveCredentials(765)
		require.NoError(t, err)
	})

	t.Run("test remove credentials: ok", func(t *testing.T) {
		err = RemoveCredentials(1)
		require.NoError(t, err)
	})
}
