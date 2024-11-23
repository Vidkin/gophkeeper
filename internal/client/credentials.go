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

func AddCredentials(credentials *proto.Credentials) error {
	f, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(f)

	credentials.Login, err = aes.Encrypt(viper.GetString("secret_key"), credentials.Login)
	if err != nil {
		return err
	}
	credentials.Password, err = aes.Encrypt(viper.GetString("secret_key"), credentials.Password)
	if err != nil {
		return err
	}
	credentials.Description, err = aes.Encrypt(viper.GetString("secret_key"), credentials.Description)
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

	req := &proto.AddUserCredentialsRequest{
		Credentials: credentials,
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

	_, err = client.AddUserCredentials(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Successfully add a new user credentials!")
	return err
}

func GetAllCredentials() error {
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

	req := &proto.GetUserCredentialsRequest{}

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

	resp, err := client.GetUserCredentials(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("User credentials:")
	for _, cred := range resp.Credentials {
		cred.Login, err = aes.Decrypt(viper.GetString("secret_key"), cred.Login)
		if err != nil {
			return fmt.Errorf("failed to decrypt credentials info, check secret key, original error: %v", err)
		}
		cred.Password, err = aes.Decrypt(viper.GetString("secret_key"), cred.Password)
		if err != nil {
			return fmt.Errorf("failed to decrypt credentials info, check secret key, original error: %v", err)
		}
		cred.Description, err = aes.Decrypt(viper.GetString("secret_key"), cred.Description)
		if err != nil {
			return fmt.Errorf("failed to decrypt credentials info, check secret key, original error: %v", err)
		}
		fmt.Printf("id=%d, login=%s, password=%s, description=%s\n", cred.Id, cred.Login, cred.Password, cred.Description)
	}
	return err
}

func GetCredentials(credID int64) error {
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

	req := &proto.GetUserCredentialRequest{Id: strconv.FormatInt(credID, 10)}

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

	resp, err := client.GetUserCredential(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Credentials:")
	resp.Credentials.Login, err = aes.Decrypt(viper.GetString("secret_key"), resp.Credentials.Login)
	if err != nil {
		return fmt.Errorf("failed to decrypt credentials info, check secret key, original error: %v", err)
	}

	resp.Credentials.Password, err = aes.Decrypt(viper.GetString("secret_key"), resp.Credentials.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt credentials info, check secret key, original error: %v", err)
	}

	resp.Credentials.Description, err = aes.Decrypt(viper.GetString("secret_key"), resp.Credentials.Description)
	if err != nil {
		return fmt.Errorf("failed to decrypt credentials info, check secret key, original error: %v", err)
	}

	fmt.Printf(
		"id=%d, login=%s, password=%s, description=%s\n",
		resp.Credentials.Id, resp.Credentials.Login, resp.Credentials.Password, resp.Credentials.Description)
	return err
}

func RemoveCredentials(credID int64) error {
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

	req := &proto.RemoveUserCredentialsRequest{Id: strconv.FormatInt(credID, 10)}

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

	_, err = client.RemoveUserCredentials(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Credentials has been successfully removed")
	return err
}
