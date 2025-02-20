package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
)

const (
	RegisterMethod = "Register"
)

func (s Server) Register(ctx context.Context, in *pb.RegisterIn) (*pb.RegisterOut, error) {
	handler, ok := s.handlers[RegisterMethod]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "handler for %s methos not found", RegisterMethod)
	}

	h, ok := handler.(func(context.Context, *pb.RegisterIn) (*pb.RegisterOut, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", RegisterMethod)
	}

	return h(ctx, in)
}
