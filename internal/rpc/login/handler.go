package login

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/infrastructure/validator"
	"github.com/bjlag/go-keeper/internal/usecase/server/user/login"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

// Handle метод для аутентификации пользователя.
func (h *Handler) Handle(ctx context.Context, in *pb.LoginIn) (*pb.LoginOut, error) {
	log := logger.FromCtx(ctx)

	if !validator.ValidateEmail(in.GetEmail()) {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	if len(in.GetPassword()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	result, err := h.usecase.Do(ctx, login.Data{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	})
	if err != nil {
		switch {
		case errors.Is(err, login.ErrUserNotFound):
			return nil, status.Error(codes.Unauthenticated, "credentials incorrect")
		case errors.Is(err, login.ErrPasswordIncorrect):
			return nil, status.Error(codes.Unauthenticated, "credentials incorrect")
		default:
			log.Error("Failed to login user", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &pb.LoginOut{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}
