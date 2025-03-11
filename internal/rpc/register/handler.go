package register

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/infrastructure/validator"
	"github.com/bjlag/go-keeper/internal/usecase/server/user/register"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.RegisterIn) (*pb.RegisterOut, error) {
	log := logger.FromCtx(ctx)

	if !validator.ValidateEmail(in.GetEmail()) {
		return nil, status.Error(codes.InvalidArgument, "email is invalid")
	}

	if !validator.ValidatePassword(in.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "password is invalid (min. length 8 characters)")
	}

	result, err := h.usecase.Do(ctx, register.Data{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	})
	if err != nil {
		if errors.Is(err, register.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}

		log.Error("failed to register user", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.RegisterOut{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}
