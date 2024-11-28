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
	"golang.org/x/text/unicode/norm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Vidkin/gophkeeper/internal/client"
	minioStorage "github.com/Vidkin/gophkeeper/internal/storage"
	"github.com/Vidkin/gophkeeper/pkg/interceptors"
	"github.com/Vidkin/gophkeeper/proto"
)

const (
	TokenFileNameFiles = "gophkeeperJWTFiles.tmp"
	expiredToken       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzIzODczMzgsIlVzZXJJRCI6MX0.B6kBiV1YOiDZd1oxp4weHgkFtJcN5VebwWpRD70uQDw"
)

func setExpiredTokenFiles(t *testing.T) {
	err := os.Remove(path.Join(os.TempDir(), TokenFileNameFiles))
	if !os.IsNotExist(err) {
		require.NoError(t, err)
	}
	f, err := os.Create(path.Join(os.TempDir(), TokenFileNameFiles))
	require.NoError(t, err)
	defer f.Close()
	_, err = f.WriteString(expiredToken)
	require.NoError(t, err)
}

func TestFiles(t *testing.T) {
	storage, dbName := setupTestDB(t)
	defer teardownTestDB(t, storage.Conn, dbName)

	minioClient, err := minioStorage.NewMinioStorage(
		"127.0.0.1:9000",
		"minioadmin",
		"minioadmin",
		nil,
		"../../certs/public.crt",
	)
	require.NoError(t, err)

	gs := &GophkeeperServer{
		Minio:       minioClient,
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

	setExpiredTokenFiles(t)

	f, err := os.Open(path.Join(os.TempDir(), TokenFileNameFiles))
	require.NoError(t, err)
	defer f.Close()
	fStat, err := f.Stat()
	require.NoError(t, err)

	file := proto.File{}
	file.FileSize = fStat.Size()
	file.Description = "description"
	file.FileName = norm.NFC.String(path.Base(TokenFileNameFiles))

	ft, err := os.ReadFile(path.Join(os.TempDir(), TokenFileNameFiles))
	require.NoError(t, err)
	token := string(ft)
	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	t.Run("upload file: error expired token", func(t *testing.T) {
		stream, err := client.Upload(ctx)
		require.NoError(t, err)
		req := &proto.FileUploadRequest{
			FileName:    file.FileName,
			Description: file.Description,
			FileSize:    file.FileSize,
			Chunk:       ft,
		}
		require.NoError(t, err)
		err = stream.Send(req)
		require.NoError(t, err)
		_, err = stream.CloseAndRecv()
		require.ErrorContains(t, err, "error parse claims")
	})

	t.Run("test download file: error expired token", func(t *testing.T) {
		req := &proto.FileDownloadRequest{
			FileName: norm.NFC.String("badFileName"),
		}
		stream, err := client.Download(ctx, req)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.ErrorContains(t, err, "error parse claims")
	})

	resp, err := client.Authorize(context.Background(), &proto.AuthorizeRequest{Credentials: &cred})
	require.NoError(t, err)
	os.Remove(path.Join(os.TempDir(), TokenFileNameFiles))

	tf, err := os.Create(path.Join(os.TempDir(), TokenFileNameFiles))
	require.NoError(t, err)
	defer tf.Close()
	_, err = tf.WriteString(resp.Token)
	require.NoError(t, err)

	md = metadata.New(map[string]string{"token": resp.Token})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	t.Run("upload file ok", func(t *testing.T) {
		stream, err := client.Upload(ctx)
		require.NoError(t, err)
		req := &proto.FileUploadRequest{
			FileName:    file.FileName,
			Description: file.Description,
			FileSize:    file.FileSize,
			Chunk:       ft,
		}
		err = stream.Send(req)
		require.NoError(t, err)
		_, err = stream.CloseAndRecv()
		require.NoError(t, err)
	})

	t.Run("upload file 2000 bytes", func(t *testing.T) {
		stream, err := client.Upload(ctx)
		req := &proto.FileUploadRequest{
			FileName:    file.FileName + "chunks",
			Description: file.Description,
			FileSize:    2000,
			Chunk:       make([]byte, 2000),
		}
		err = stream.Send(req)
		require.NoError(t, err)
		_, err = stream.CloseAndRecv()
		require.NoError(t, err)
	})

	t.Run("upload file error upload file with wrong file size", func(t *testing.T) {
		stream, err := client.Upload(ctx)
		require.NoError(t, err)
		req := &proto.FileUploadRequest{
			FileName:    file.FileName,
			Description: file.Description,
			FileSize:    file.FileSize * 2,
			Chunk:       ft,
		}
		err = stream.Send(req)
		require.NoError(t, err)
		_, err = stream.CloseAndRecv()
		require.Error(t, err)
	})

	t.Run("test get files: ok", func(t *testing.T) {
		resp, err := client.GetFiles(ctx, &proto.GetFilesRequest{})
		require.NoError(t, err)
		assert.Equal(t, file.FileSize, resp.Files[0].FileSize)
		assert.Equal(t, file.FileName, resp.Files[0].FileName)
		assert.Equal(t, file.Description, resp.Files[0].Description)
	})

	df, err := os.OpenFile(path.Join(os.TempDir(), "testDownloadFile.tmp"), os.O_WRONLY|os.O_CREATE, 0666)
	require.NoError(t, err)
	defer df.Close()

	t.Run("test download file: file not found", func(t *testing.T) {
		req := &proto.FileDownloadRequest{
			FileName: norm.NFC.String("badFileName"),
		}
		stream, err := client.Download(ctx, req)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.ErrorContains(t, err, "error getting file info")
	})

	t.Run("test download file: file name is empty", func(t *testing.T) {
		req := &proto.FileDownloadRequest{
			FileName: "",
		}
		stream, err := client.Download(ctx, req)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.ErrorContains(t, err, "file name is required")
	})

	t.Run("test download file: ok", func(t *testing.T) {
		req := &proto.FileDownloadRequest{
			FileName: norm.NFC.String(TokenFileNameFiles),
		}
		stream, err := client.Download(ctx, req)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.NoError(t, err)
	})

	t.Run("download file 2000 bytes", func(t *testing.T) {
		req := &proto.FileDownloadRequest{
			FileName: file.FileName + "chunks",
		}
		stream, err := client.Download(ctx, req)
		require.NoError(t, err)
		_, err = stream.Recv()
		require.NoError(t, err)
	})

	t.Run("test remove file: empty file name", func(t *testing.T) {
		_, err = client.RemoveFile(ctx, &proto.FileRemoveRequest{FileName: ""})
		require.ErrorContains(t, err, "you must provide file name")
	})

	t.Run("test remove file: file not found", func(t *testing.T) {
		_, err = client.RemoveFile(ctx, &proto.FileRemoveRequest{FileName: "badFileName"})
		require.ErrorContains(t, err, "file not found")
	})

	t.Run("test remove file: ok", func(t *testing.T) {
		_, err = client.RemoveFile(ctx, &proto.FileRemoveRequest{FileName: TokenFileNameFiles})
		require.NoError(t, err)
	})
}
