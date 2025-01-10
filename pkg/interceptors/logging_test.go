package interceptors

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Invoke(ctx context.Context, req interface{}) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestLoggingInterceptor_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Method"}

	handler := &MockHandler{}
	handler.On("Invoke", ctx, req).Return("test response", nil)

	resp, err := LoggingInterceptor(ctx, req, info, handler.Invoke)

	assert.NoError(t, err)
	assert.Equal(t, "test response", resp)
	handler.AssertExpectations(t)
}

func TestLoggingInterceptor_Error(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Method"}

	handler := &MockHandler{}
	handler.On("Invoke", ctx, req).Return(nil, errors.New("test error"))

	resp, err := LoggingInterceptor(ctx, req, info, handler.Invoke)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "test error", err.Error())
	handler.AssertExpectations(t)
}

func TestLoggingInterceptor_StatusError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Method"}

	handler := &MockHandler{}
	handler.On("Invoke", ctx, req).Return(nil, status.Error(codes.Internal, "internal error"))

	resp, err := LoggingInterceptor(ctx, req, info, handler.Invoke)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal.String(), status.Code(err).String())
	handler.AssertExpectations(t)
}
