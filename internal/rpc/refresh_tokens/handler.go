package refresh_tokens

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/usecase/user/refresh_tokens"
)

const lenRefreshToken = 24

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.RefreshTokensIn) (*pb.RefreshTokensOut, error) {
	log := logger.FromCtx(ctx)

	if len(in.GetRefreshToken()) < lenRefreshToken {
		return nil, status.Error(codes.InvalidArgument, "refresh token too short")
	}

	result, err := h.usecase.Do(ctx, refresh_tokens.Data{RefreshToken: in.GetRefreshToken()})
	if err != nil {
		if errors.Is(err, refresh_tokens.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.FailedPrecondition, "invalid refresh token")
		}

		log.Error("Failed to refresh tokens", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.RefreshTokensOut{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}
