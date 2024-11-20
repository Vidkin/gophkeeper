package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "google.golang.org/protobuf/proto"

	"github.com/Vidkin/gophkeeper/pkg/aes"
	"github.com/Vidkin/gophkeeper/pkg/hash"
	"github.com/Vidkin/gophkeeper/proto"
)

func AddCard(card *proto.BankCard) error {
	f, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(f)

	card.Cvv, err = aes.Encrypt(viper.GetString("secret_key"), card.Cvv)
	if err != nil {
		return err
	}
	card.Owner, err = aes.Encrypt(viper.GetString("secret_key"), card.Owner)
	if err != nil {
		return err
	}
	card.Number, err = aes.Encrypt(viper.GetString("secret_key"), card.Number)
	if err != nil {
		return err
	}
	card.Description, err = aes.Encrypt(viper.GetString("secret_key"), card.Description)
	if err != nil {
		return err
	}
	card.ExpireDate, err = aes.Encrypt(viper.GetString("secret_key"), card.ExpireDate)
	if err != nil {
		return err
	}

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

	req := &proto.AddBankCardRequest{
		Card: card,
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	md := metadata.New(map[string]string{"token": token})
	ctxTimeout = metadata.NewOutgoingContext(ctxTimeout, md)

	if viper.GetString("hash_key") != "" {
		data, err := pb.Marshal(req)
		if err != nil {
			return err
		}
		h := hash.GetHashSHA256(viper.GetString("hash_key"), data)
		hEnc := base64.StdEncoding.EncodeToString(h)
		md.Append("HashSHA256", hEnc)
		ctxTimeout = metadata.NewOutgoingContext(ctxTimeout, md)
	}

	_, err = client.AddBankCard(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Successfully add a new bank card!")
	return err
}
