package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "google.golang.org/protobuf/proto"

	"github.com/Vidkin/gophkeeper/pkg/hash"
	"github.com/Vidkin/gophkeeper/proto"
)

func Auth(login, password string) error {
	client, conn, err := NewGophkeeperClient()
	if err != nil {
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			fmt.Println("failed to close grpc connection")
		}
	}(conn)

	cred := proto.Credentials{
		Login:    login,
		Password: password,
	}
	req := &proto.AuthorizeRequest{
		Credentials: &cred,
	}
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if viper.GetString("hash_key") != "" {
		data, err := pb.Marshal(req)
		if err != nil {
			return err
		}
		h := hash.GetHashSHA256(viper.GetString("hash_key"), data)
		hEnc := base64.StdEncoding.EncodeToString(h)
		md := metadata.New(map[string]string{"HashSHA256": hEnc})
		ctxTimeout = metadata.NewOutgoingContext(ctxTimeout, md)
	}

	resp, err := client.Authorize(ctxTimeout, req)

	if err != nil {
		return err
	}

	err = os.Remove(path.Join(os.TempDir(), "gophkeeperJWT.tmp"))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	f, err := os.Create(path.Join(os.TempDir(), "gophkeeperJWT.tmp"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			fmt.Println("failed to close jwt temp file")
		}
	}(f)

	if _, err = f.WriteString(resp.Token); err != nil {
		return err
	}

	fmt.Println("Successfully authorized!")

	return err
}
