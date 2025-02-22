package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
)

const (
	RegisterMethod = "Register"
	LoginMethod    = "Login"
)

func (s Server) Register(ctx context.Context, in *pb.RegisterIn) (*pb.RegisterOut, error) {
	handler, err := s.getHandler(RegisterMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.RegisterIn) (*pb.RegisterOut, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", RegisterMethod)
	}

	return h(ctx, in)
}

func (s Server) Login(ctx context.Context, in *pb.LoginIn) (*pb.LoginOut, error) {
	handler, err := s.getHandler(LoginMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.LoginIn) (*pb.LoginOut, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", LoginMethod)
	}

	return h(ctx, in)
}

func (s Server) getHandler(name string) (any, error) {
	handler, ok := s.handlers[name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "handler for %s methos not found", name)
	}

	return handler, nil
}
