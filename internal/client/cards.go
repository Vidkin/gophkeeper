package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
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

// AddCard adds a new bank card to the GophKeeper server after encrypting its sensitive information.
//
// Parameters:
//   - card: A pointer to the proto.BankCard struct containing the card details to be added.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, encryption, or gRPC communication.
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

// GetAllCards retrieves all bank cards from the GophKeeper server and decrypts their information for display.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, gRPC communication, or decryption.
func GetAllCards() error {
	f, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(f)

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

	req := &proto.GetBankCardsRequest{}

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

	resp, err := client.GetBankCards(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Bank cards:")
	for _, card := range resp.Cards {
		card.Owner, err = aes.Decrypt(viper.GetString("secret_key"), card.Owner)
		if err != nil {
			return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
		}
		card.Description, err = aes.Decrypt(viper.GetString("secret_key"), card.Description)
		if err != nil {
			return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
		}
		card.Number, err = aes.Decrypt(viper.GetString("secret_key"), card.Number)
		if err != nil {
			return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
		}
		card.Cvv, err = aes.Decrypt(viper.GetString("secret_key"), card.Cvv)
		if err != nil {
			return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
		}
		card.ExpireDate, err = aes.Decrypt(viper.GetString("secret_key"), card.ExpireDate)
		if err != nil {
			return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
		}
		fmt.Printf("id=%d, number=%s, owner=%s, cvv=%s, expire=%s, description=%s\n", card.Id, card.Number, card.Cvv, card.Owner, card.Description, card.ExpireDate)
	}
	return err
}

// GetCard retrieves a specific bank card by its ID from the GophKeeper server and decrypts its information for display.
//
// Parameters:
//   - cardID: The ID of the bank card to retrieve.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, gRPC communication, or decryption.
func GetCard(cardID int64) error {
	f, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(f)

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

	req := &proto.GetBankCardRequest{Id: strconv.FormatInt(cardID, 10)}

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

	resp, err := client.GetBankCard(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Bank card:")
	resp.Card.Owner, err = aes.Decrypt(viper.GetString("secret_key"), resp.Card.Owner)
	if err != nil {
		return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
	}
	resp.Card.Description, err = aes.Decrypt(viper.GetString("secret_key"), resp.Card.Description)
	if err != nil {
		return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
	}
	resp.Card.Number, err = aes.Decrypt(viper.GetString("secret_key"), resp.Card.Number)
	if err != nil {
		return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
	}
	resp.Card.Cvv, err = aes.Decrypt(viper.GetString("secret_key"), resp.Card.Cvv)
	if err != nil {
		return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
	}
	resp.Card.ExpireDate, err = aes.Decrypt(viper.GetString("secret_key"), resp.Card.ExpireDate)
	if err != nil {
		return fmt.Errorf("failed to decrypt card info, check secret key, original error: %v", err)
	}
	fmt.Printf(
		"id=%d, number=%s, owner=%s, cvv=%s, expire=%s, description=%s\n",
		resp.Card.Id, resp.Card.Number, resp.Card.Cvv, resp.Card.Owner, resp.Card.Description, resp.Card.ExpireDate)
	return err
}

// RemoveCard removes a bank card from the GophKeeper server by its ID.
//
// Parameters:
//   - cardID: The ID of the bank card to remove.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, gRPC communication, or authorization issues.
func RemoveCard(cardID int64) error {
	f, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(f)

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

	req := &proto.RemoveBankCardRequest{Id: strconv.FormatInt(cardID, 10)}

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

	_, err = client.RemoveBankCard(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Bank card has been successfully removed")
	return err
}
