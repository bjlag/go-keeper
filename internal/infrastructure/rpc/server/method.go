package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
)

const (
	RegisterMethod      = "Register"
	LoginMethod         = "Login"
	RefreshTokensMethod = "RefreshTokens"
	GetAllItemsMethod   = "GetAllItems"
)

func (s RPCServer) Register(ctx context.Context, in *pb.RegisterIn) (*pb.RegisterOut, error) {
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

func (s RPCServer) Login(ctx context.Context, in *pb.LoginIn) (*pb.LoginOut, error) {
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

func (s RPCServer) RefreshTokens(ctx context.Context, in *pb.RefreshTokensIn) (*pb.RefreshTokensOut, error) {
	handler, err := s.getHandler(RefreshTokensMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.RefreshTokensIn) (*pb.RefreshTokensOut, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", RefreshTokensMethod)
	}

	return h(ctx, in)
}

func (s RPCServer) GetAllItems(ctx context.Context, in *pb.GetAllItemsIn) (*pb.GetAllItemsOut, error) {
	handler, err := s.getHandler(GetAllItemsMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.GetAllItemsIn) (*pb.GetAllItemsOut, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", GetAllItemsMethod)
	}

	return h(ctx, in)
}

func (s RPCServer) getHandler(name string) (any, error) {
	handler, ok := s.handlers[name]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "handler for %s methos not found", name)
	}

	return handler, nil
}
