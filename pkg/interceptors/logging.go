// Package interceptors provides gRPC interceptors for handling requests and responses.
//
// This package includes the LoggingInterceptor function, which logs the details of
// incoming requests and their corresponding responses.
package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Vidkin/gophkeeper/internal/logger"
)

// LoggingInterceptor is a gRPC unary server interceptor that logs the details of
// incoming requests and their responses, including the method name, duration,
// and response status.
//
// Parameters:
//   - ctx: A context.Context for managing request-scoped values and cancellation.
//   - req: An interface{} representing the incoming request message.
//   - info: A pointer to grpc.UnaryServerInfo containing information about the method being called.
//   - handler: A grpc.UnaryHandler that processes the request and returns a response.
//
// Returns:
//   - An interface{} containing the response from the handler.
//   - An error if the handler returns an error.
//
// The interceptor records the start time of the request, invokes the handler to process
// the request, and calculates the duration of the request. It then logs the method name,
// duration, and response status. If an error occurs during the handling of the request,
// the error status is logged instead. This interceptor is useful for monitoring and
// debugging gRPC services by providing insights into request processing times and outcomes.
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	itf, err := handler(ctx, req)
	duration := time.Since(startTime)

	var respStatus string
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			respStatus = st.Code().String()
		} else {
			respStatus = err.Error()
		}
	} else {
		respStatus = codes.OK.String()
	}
	logger.Log.Info(
		"Request data",
		zap.String("method", info.FullMethod),
		zap.Duration("duration", duration),
	)
	logger.Log.Info(
		"Response data",
		zap.String("status", respStatus),
	)
	return itf, err
}
