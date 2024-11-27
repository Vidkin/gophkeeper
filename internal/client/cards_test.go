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

func TestCards(t *testing.T) {
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

	card := proto.BankCard{
		ExpireDate:  "12.03.2024",
		Number:      "12312313",
		Cvv:         "123",
		Owner:       "owner",
		Description: "description",
	}
	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test add card: invalid key size", func(t *testing.T) {
		err = AddCard(&card)
		require.ErrorContains(t, err, "invalid key size")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	t.Run("test add card: missing hash", func(t *testing.T) {
		err = AddCard(&card)
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test add card: missed token file", func(t *testing.T) {
		err = AddCard(&card)
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test add card: expired token", func(t *testing.T) {
		err = AddCard(&card)
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test add card: ok", func(t *testing.T) {
		err = AddCard(&card)
		require.NoError(t, err)
	})

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test get all cards: missing hash", func(t *testing.T) {
		err = GetAllCards()
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test get all cards: missed token file", func(t *testing.T) {
		err = GetAllCards()
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test get all cards: expired token", func(t *testing.T) {
		err = GetAllCards()
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test get all cards: ok", func(t *testing.T) {
		err = GetAllCards()
		require.NoError(t, err)
	})

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test get card: missing hash", func(t *testing.T) {
		err = GetCard(1)
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test get card: missed token file", func(t *testing.T) {
		err = GetCard(1)
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test get card: expired token", func(t *testing.T) {
		err = GetCard(1)
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test get card: ok", func(t *testing.T) {
		err = GetCard(1)
		require.NoError(t, err)
	})

	viper.Set("secret_key", "")
	viper.Set("hash_key", "")
	t.Run("test remove card: missing hash", func(t *testing.T) {
		err = RemoveCard(1)
		require.ErrorContains(t, err, "missing hash")
	})

	err = os.Remove(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	t.Run("test remove card: missed token file", func(t *testing.T) {
		err = RemoveCard(1)
		require.ErrorContains(t, err, "no such file or directory")
	})

	viper.Set("secret_key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x")
	viper.Set("hash_key", "defaultHashKey")
	setExpiredToken(t)
	t.Run("test remove card: expired token", func(t *testing.T) {
		err = RemoveCard(1)
		require.ErrorContains(t, err, "need to re-authorize")
	})

	err = Auth("test_login", "test_pass")
	require.NoError(t, err)
	t.Run("test remove unknown card: ok", func(t *testing.T) {
		err = RemoveCard(765)
		require.NoError(t, err)
	})

	t.Run("test remove card: ok", func(t *testing.T) {
		err = RemoveCard(1)
		require.NoError(t, err)
	})
}
