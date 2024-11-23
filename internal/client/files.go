package client

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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

	"github.com/Vidkin/gophkeeper/pkg/hash"
	"github.com/Vidkin/gophkeeper/proto"
)

func UploadFile(filePath, description string) error {
	tokenFile, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(tokenFile)
	fmt.Println(token)

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			fmt.Println("failed to close file")
		}
	}(f)
	fStat, err := f.Stat()
	if err != nil {
		return err
	}

	file := proto.File{}
	file.FileSize = fStat.Size()
	file.Description = description
	file.FileName = path.Base(filePath)

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

	ctx := context.Background()
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	reader := bufio.NewReader(f)
	buffer := make([]byte, 1024*1024)

	stream, err := client.Upload(ctx)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		req := &proto.FileUploadRequest{
			FileName:    file.FileName,
			Description: file.Description,
			FileSize:    file.FileSize,
			Chunk:       buffer[:n],
		}
		err = stream.Send(req)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	fmt.Println("Successfully upload file!")
	return err
}

func RemoveFile(fileID int64) error {
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

	req := &proto.FileRemoveRequest{Id: strconv.FormatInt(fileID, 10)}

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

	_, err = client.RemoveFile(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("File has been successfully removed")
	return err
}

func GetAllFiles() error {
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

	req := &proto.GetFilesRequest{}

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

	resp, err := client.GetFiles(ctxTimeout, req)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.PermissionDenied {
				return errors.New("need to re-authorize, call auth command")
			}
		}
		return err
	}

	fmt.Println("Files:")
	for _, file := range resp.Files {
		fmt.Printf("id=%d, fileName=%s, size=%d, description=%s\n", file.Id, file.FileName, file.FileSize, file.Description)
	}
	return err
}
