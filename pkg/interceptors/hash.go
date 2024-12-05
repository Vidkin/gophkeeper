// Package interceptors provides gRPC interceptors for handling requests and responses.
//
// This package includes the HashInterceptor function, which verifies the integrity of
// incoming requests by comparing a provided hash with a computed hash of the request data.
package interceptors

import (
	"bytes"
	"context"
	"encoding/base64"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/Vidkin/gophkeeper/internal/logger"
	"github.com/Vidkin/gophkeeper/pkg/hash"
)

// HashInterceptor returns a gRPC unary server interceptor that verifies the SHA-256 hash
// of incoming requests against a provided key.
//
// Parameters:
//   - key: A string representing the key used for hash computation. If the key is empty,
//     the interceptor will skip hash verification.
//
// The interceptor extracts the "HashSHA256" metadata from the incoming context and decodes
// it from a base64 string. It then marshals the request into a byte slice and computes
// the SHA-256 hash using the provided key. If the computed hash does not match the
// provided hash, an error is returned, indicating that the hashes do not match.
//
// Returns:
//   - A function that implements the gRPC UnaryHandler signature, which processes the
//     request if the hash verification is successful, or returns an error if it fails.
func HashInterceptor(key string) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if key == "" {
			return handler(ctx, req)
		}

		var hEnc string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("HashSHA256")
			if len(values) > 0 {
				hEnc = values[0]
			}
		}
		if len(hEnc) == 0 {
			return nil, status.Error(codes.InvalidArgument, "missing hash")
		}

		hashA, err := base64.StdEncoding.DecodeString(hEnc)
		if err != nil {
			logger.Log.Error("error decode hash from base64 string", zap.Error(err))
			return nil, status.Error(codes.Internal, "missing hash")
		}

		var data []byte
		if msg, ok := req.(proto.Message); ok {
			data, err = proto.Marshal(msg)
			if err != nil {
				logger.Log.Error("failed to marshal request: %v", zap.Error(err))
				return nil, status.Errorf(codes.Internal, "failed to marshal request")
			}
		} else {
			return nil, status.Errorf(codes.Internal, "failed to get proto.Message")
		}

		hashB := hash.GetHashSHA256(key, data)
		if !bytes.Equal(hashA, hashB) {
			logger.Log.Error("hashes don't match")
			return nil, status.Errorf(codes.InvalidArgument, "hashes don't match")
		}

		return handler(ctx, req)
	}
}
