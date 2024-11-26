package interceptors

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/Vidkin/gophkeeper/pkg/hash"
	gkProto "github.com/Vidkin/gophkeeper/proto"
)

type MockHandlerHash struct {
	mock.Mock
}

func (m *MockHandlerHash) Invoke(ctx context.Context, req interface{}) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestHashInterceptor(t *testing.T) {
	key := "test-key"
	interceptor := HashInterceptor(key)

	t.Run("valid hash", func(t *testing.T) {
		req := &gkProto.EchoRequest{Message: "test"}
		data, _ := proto.Marshal(req)
		expectedHash := hash.GetHashSHA256(key, data)
		encodedHash := base64.StdEncoding.EncodeToString(expectedHash)

		md := metadata.Pairs("HashSHA256", encodedHash)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		mockHandler := new(MockHandlerHash)
		mockHandler.On("Invoke", ctx, req).Return(req, nil)

		resp, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, mockHandler.Invoke)

		assert.NoError(t, err)
		assert.Equal(t, req, resp)
		mockHandler.AssertExpectations(t)
	})

	t.Run("missing hash", func(t *testing.T) {
		req := &gkProto.EchoRequest{Message: "test"}
		ctx := context.Background()

		resp, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return req, nil
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, status.Code(err), codes.InvalidArgument)
	})

	t.Run("hashes don't match", func(t *testing.T) {
		req := &gkProto.EchoRequest{Message: "test"}
		invalidHash := base64.StdEncoding.EncodeToString([]byte("invalid-hash"))
		md := metadata.Pairs("HashSHA256", invalidHash)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return req, nil
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, status.Code(err), codes.InvalidArgument)
	})

	t.Run("empty key", func(t *testing.T) {
		interceptor := HashInterceptor("")
		req := &gkProto.EchoRequest{Message: "test"}
		ctx := context.Background()

		resp, err := interceptor(ctx, req, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (interface{}, error) {
			return req, nil
		})

		assert.NoError(t, err)
		assert.Equal(t, req, resp)
	})
}
