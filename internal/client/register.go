package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "google.golang.org/protobuf/proto"

	"github.com/Vidkin/gophkeeper/pkg/hash"
	"github.com/Vidkin/gophkeeper/proto"
)

// Register registers a new user with the GophKeeper server.
//
// Parameters:
//   - login: The login name of the user.
//   - password: The password for the user account.
//
// Returns an error if the registration fails, for example, if the user already exists
// or if there is a connection issue with the server. If successful, a confirmation message
// is printed to the console.
func Register(login, password string) error {
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
	req := &proto.RegisterUserRequest{
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

	_, err = client.RegisterUser(ctxTimeout, req)
	if err == nil {
		fmt.Println("User successfully registered!")
	}
	return err
}
