package interceptor

import (
	"context"
	"pkg/log"

	"google.golang.org/grpc"
)

type LoggerInterceptor struct{}

func NewLoggerInterceptor() *LoggerInterceptor {
	return &LoggerInterceptor{}
}

func (_ *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		log.WithContextParams(ctx).Info(info.FullMethod)

		return handler(ctx, req)
	}
}

func (interceptor *LoggerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		log.WithContextParams(stream.Context()).Info(info.FullMethod)

		return handler(srv, stream)
	}
}
