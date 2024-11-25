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

// AddNote adds a new note to the GophKeeper server.
//
// Parameters:
//   - note: A pointer to the proto.Note structure containing the text and description of the note.
//
// Returns an error if the operation fails, for example, if re-authorization is required.
func AddNote(note *proto.Note) error {
	f, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(f)

	note.Text, err = aes.Encrypt(viper.GetString("secret_key"), note.Text)
	if err != nil {
		return err
	}
	note.Description, err = aes.Encrypt(viper.GetString("secret_key"), note.Description)
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

	req := &proto.AddNoteRequest{
		Note: note,
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

	_, err = client.AddNote(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Successfully add a new user note!")
	return err
}

// GetAllNotes retrieves all user notes from the GophKeeper server.
//
// Returns an error if the operation fails, for example, if re-authorization is required.
func GetAllNotes() error {
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

	req := &proto.GetNotesRequest{}

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

	resp, err := client.GetNotes(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("User notes:")
	for _, note := range resp.Notes {
		note.Text, err = aes.Decrypt(viper.GetString("secret_key"), note.Text)
		if err != nil {
			return fmt.Errorf("failed to decrypt note info, check secret key, original error: %v", err)
		}
		note.Description, err = aes.Decrypt(viper.GetString("secret_key"), note.Description)
		if err != nil {
			return fmt.Errorf("failed to decrypt note info, check secret key, original error: %v", err)
		}
		fmt.Printf("id=%d, text=%s, description=%s\n", note.Id, note.Text, note.Description)
	}
	return err
}

// GetNote retrieves a note by its ID from the GophKeeper server.
//
// Parameters:
//   - noteID: The ID of the note to retrieve.
//
// Returns an error if the operation fails, for example, if re-authorization is required.
func GetNote(noteID int64) error {
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

	req := &proto.GetNoteRequest{Id: strconv.FormatInt(noteID, 10)}

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

	resp, err := client.GetNote(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Note:")
	resp.Note.Text, err = aes.Decrypt(viper.GetString("secret_key"), resp.Note.Text)
	if err != nil {
		return fmt.Errorf("failed to decrypt note, check secret key, original error: %v", err)
	}

	resp.Note.Description, err = aes.Decrypt(viper.GetString("secret_key"), resp.Note.Description)
	if err != nil {
		return fmt.Errorf("failed to decrypt note, check secret key, original error: %v", err)
	}

	fmt.Printf(
		"id=%d, text=%s, description=%s\n",
		resp.Note.Id, resp.Note.Text, resp.Note.Description)
	return err
}

// RemoveNote removes a note by its ID from the GophKeeper server.
//
// Parameters:
//   - noteID: The ID of the note to remove.
//
// Returns an error if the operation fails, for example, if re-authorization is required.
func RemoveNote(noteID int64) error {
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

	req := &proto.RemoveNoteRequest{Id: strconv.FormatInt(noteID, 10)}

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

	_, err = client.RemoveNote(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Note has been successfully removed")
	return err
}
