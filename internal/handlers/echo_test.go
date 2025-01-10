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

	"github.com/Vidkin/gophkeeper/internal/client"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

func TestEcho(t *testing.T) {
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

	t.Run("test echo: ok", func(t *testing.T) {
		resp, err := client.Echo(context.Background(), &proto.EchoRequest{Message: "test"})
		require.NoError(t, err)
		assert.Equal(t, "test", resp.Message)
	})
}
