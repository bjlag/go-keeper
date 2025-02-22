package interceptor

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggerServerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		log.Info("Got RPC request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Any("response", resp),
			zap.Error(err),
			zap.String("code", status.Code(err).String()),
		)

		return resp, err
	}
}
