// Package server отвечает за работу gRPC сервера.
package server

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/interceptor"
)

type RPCServer struct {
	pb.UnimplementedKeeperServer

	listener net.Listener
	handlers map[string]any
	jwt      *auth.JWT
	log      *zap.Logger
}

func NewRPCServer(opts ...Option) *RPCServer {
	s := &RPCServer{
		handlers: make(map[string]any),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s RPCServer) Start(ctx context.Context) error {
	const op = "server.rpc.Start"

	s.log.Info("Starting gRPC server",
		zap.String("addr", s.listener.Addr().String()),
	)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerServerInterceptor(s.log),
			interceptor.CheckAccessTokenServerInterceptor(s.jwt, s.log),
		),
	)

	pb.RegisterKeeperServer(grpcServer, s)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return grpcServer.Serve(s.listener)
	})
	g.Go(func() error {
		<-gCtx.Done()

		s.log.Info("Shutting down gRPC server")
		grpcServer.GracefulStop()

		return nil
	})
	if err := g.Wait(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
