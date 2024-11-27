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

func TestFiles(t *testing.T) {
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

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test get all files: missing hash", func(t *testing.T) {
		err = GetAllFiles()
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test get all files: missed token file", func(t *testing.T) {
		err = GetAllFiles()
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test get all files: expired token", func(t *testing.T) {
		err = GetAllFiles()
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test get all files: ok", func(t *testing.T) {
		err = GetAllFiles()
		require.NoError(t, err)
	})

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test remove file: missing hash", func(t *testing.T) {
		err = RemoveFile("file")
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test remove file: missed token file", func(t *testing.T) {
		err = RemoveFile("file")
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test remove note: expired token", func(t *testing.T) {
		err = RemoveFile("file")
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test remove unknown file error", func(t *testing.T) {
		err = RemoveFile("fileUnknown")
		require.ErrorContains(t, err, "file not found")
	})

	//t.Run("test remove file: ok", func(t *testing.T) {
	//	err = RemoveFile("file")
	//	require.NoError(t, err)
	//})
}
