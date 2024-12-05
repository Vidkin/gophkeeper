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
	"time"

	"github.com/spf13/viper"
	"golang.org/x/text/unicode/norm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pb "google.golang.org/protobuf/proto"

	"github.com/Vidkin/gophkeeper/pkg/hash"
	"github.com/Vidkin/gophkeeper/proto"
)

// UploadFile uploads a file to the GophKeeper server with an optional description.
//
// Parameters:
//   - filePath: The path to the file to be uploaded.
//   - description: A description of the file being uploaded.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, file opening,
//     gRPC communication, or streaming errors.
func UploadFile(filePath, description string) error {
	tokenFile, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(tokenFile)

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
	file.FileName = norm.NFC.String(path.Base(filePath))

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

// DownloadFile downloads a file from the GophKeeper server and saves it to the specified path.
//
// Parameters:
//   - fileName: The name of the file to download from the server.
//   - filePath: The local path where the file will be saved.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, file creation,
//     gRPC communication, or streaming errors.
func DownloadFile(fileName, filePath string) error {
	tokenFile, err := os.ReadFile(path.Join(os.TempDir(), TokenFileName))
	if err != nil {
		return fmt.Errorf("error open JWT file, need to authorize: %v", err)
	}
	token := string(tokenFile)

	f, err := os.OpenFile(path.Join(filePath, fileName), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			fmt.Println("failed to close file")
		}
	}(f)

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

	writer := bufio.NewWriter(f)
	req := &proto.FileDownloadRequest{
		FileName: norm.NFC.String(fileName),
	}
	stream, err := client.Download(ctx, req)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error receive file")
			return err
		}

		_, err = writer.Write(res.Chunk)
		if err != nil {
			fmt.Println("failed to write chunk")
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("failed to flush file")
		return err
	}

	fmt.Println("Successfully download file!")
	return err
}

// RemoveFile removes a file from the GophKeeper server by its name.
//
// Parameters:
//   - fileName: The name of the file to be removed from the server.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, gRPC communication,
//     or authorization issues.
func RemoveFile(fileName string) error {
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

	req := &proto.FileRemoveRequest{FileName: fileName}

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

// GetAllFiles retrieves a list of all files stored on the GophKeeper server and displays their details.
//
// Returns:
//   - An error if any step in the process fails, including JWT file access, gRPC communication,
//     or errors in retrieving the file list.
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
		fmt.Printf("id=%d, fileName=%s, size=%d, description=%s\n", file.Id, norm.NFC.String(file.FileName), file.FileSize, file.Description)
	}
	return err
}
