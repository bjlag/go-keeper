package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
)

// Зарегистрированные методы.
const (
	RegisterMethod      = "Register"
	LoginMethod         = "Login"
	RefreshTokensMethod = "RefreshTokens"
	GetByGUIDMethod     = "GetByGUID"
	GetAllItemsMethod   = "GetAllItems"
	CreateItemMethod    = "CreateItem"
	UpdateItemMethod    = "UpdateItem"
	DeleteItemMethod    = "DeleteItem"
)

// Register регистрация пользователя.
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

// Login аутентификация пользователя.
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

// RefreshTokens обновление токенов.
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

// GetByGuid получение элемента по его GUID.
func (s RPCServer) GetByGuid(ctx context.Context, in *pb.GetByGuidIn) (*pb.GetByGuidOut, error) { //nolint:revive
	handler, err := s.getHandler(GetByGUIDMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.GetByGuidIn) (*pb.GetByGuidOut, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", GetByGUIDMethod)
	}

	return h(ctx, in)
}

// GetAllItems получение всех элементов.
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

// CreateItem создание элемента.
func (s RPCServer) CreateItem(ctx context.Context, in *pb.CreateItemIn) (*emptypb.Empty, error) {
	handler, err := s.getHandler(CreateItemMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.CreateItemIn) (*emptypb.Empty, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", CreateItemMethod)
	}

	return h(ctx, in)
}

// UpdateItem обновление элемента.
func (s RPCServer) UpdateItem(ctx context.Context, in *pb.UpdateItemIn) (*pb.UpdateItemOut, error) {
	handler, err := s.getHandler(UpdateItemMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.UpdateItemIn) (*pb.UpdateItemOut, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", UpdateItemMethod)
	}

	return h(ctx, in)
}

// DeleteItem удаление элемента.
func (s RPCServer) DeleteItem(ctx context.Context, in *pb.DeleteItemIn) (*emptypb.Empty, error) {
	handler, err := s.getHandler(DeleteItemMethod)
	if err != nil {
		return nil, err
	}

	h, ok := handler.(func(context.Context, *pb.DeleteItemIn) (*emptypb.Empty, error))
	if !ok {
		return nil, status.Errorf(codes.Internal, "handler for %s method not found", DeleteItemMethod)
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
