package interceptor

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
)

func LoggerServerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		hLog := log
		hLog = hLog.
			With(zap.String("method", info.FullMethod)).
			With(zap.Any("request", req))

		resp, err := handler(logger.WithCtx(ctx, hLog), req)

		hLog.Info("Got RPC request",
			zap.Error(err),
			zap.String("code", status.Code(err).String()),
		)

		return resp, err
	}
}

func LoggerClientInterceptor(log *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)

		hLog := log
		hLog.Info("Send RPC request",
			zap.String("method", method),
			zap.Any("request", req),
			zap.String("code", status.Code(err).String()),
		)

		return err
	}
}
