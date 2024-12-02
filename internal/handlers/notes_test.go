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

func TestNotes(t *testing.T) {
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

	note := &proto.Note{
		Text:        "",
		Description: "description",
	}

	md := metadata.New(map[string]string{"token": resp.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	t.Run("test add note: error text is empty", func(t *testing.T) {
		_, err = client.AddNote(ctx, &proto.AddNoteRequest{Note: note})
		require.ErrorContains(t, err, "you must provide note text")
	})

	note.Text = "text"
	t.Run("test add note: ok", func(t *testing.T) {
		_, err = client.AddNote(ctx, &proto.AddNoteRequest{Note: note})
		require.NoError(t, err)
	})

	t.Run("test get notes: ok", func(t *testing.T) {
		resp, err := client.GetNotes(ctx, &proto.GetNotesRequest{})
		require.NoError(t, err)
		assert.Equal(t, note.Text, resp.Notes[0].Text)
		assert.Equal(t, note.Description, resp.Notes[0].Description)
	})

	t.Run("test get note: invalid id", func(t *testing.T) {
		_, err = client.GetNote(ctx, &proto.GetNoteRequest{Id: "badId"})
		require.ErrorContains(t, err, "missing note id")
	})

	t.Run("test get note: unknown id", func(t *testing.T) {
		_, err = client.GetNote(ctx, &proto.GetNoteRequest{Id: "435"})
		require.ErrorContains(t, err, "error get note from DB")
	})

	t.Run("test get note: ok", func(t *testing.T) {
		resp, err := client.GetNote(ctx, &proto.GetNoteRequest{Id: "1"})
		require.NoError(t, err)
		assert.Equal(t, note.Text, resp.Note.Text)
		assert.Equal(t, note.Description, resp.Note.Description)
	})

	t.Run("test remove note: empty id", func(t *testing.T) {
		_, err = client.RemoveNote(ctx, &proto.RemoveNoteRequest{Id: ""})
		require.ErrorContains(t, err, "you must provide note id")
	})

	t.Run("test remove note: bad id", func(t *testing.T) {
		_, err = client.RemoveNote(ctx, &proto.RemoveNoteRequest{Id: "badID"})
		require.ErrorContains(t, err, "invalid note id")
	})

	t.Run("test remove note: ok", func(t *testing.T) {
		_, err = client.RemoveNote(ctx, &proto.RemoveNoteRequest{Id: "1"})
		require.NoError(t, err)
	})
}
