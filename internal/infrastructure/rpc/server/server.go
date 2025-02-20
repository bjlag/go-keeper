package server

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/interceptor"
)

type Server struct {
	pb.UnimplementedKeeperServer

	handlers map[string]any
	log      *zap.Logger
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		handlers: make(map[string]any),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s Server) Start(ctx context.Context) error {
	const op = "rpc.Start"

	s.log.Info("Starting gRPC server")

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerServerInterceptor(s.log),
		),
	)

	pb.RegisterKeeperServer(grpcServer, s)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return grpcServer.Serve(listen)
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
