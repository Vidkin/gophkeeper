package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Vidkin/gophkeeper/proto"
)

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
