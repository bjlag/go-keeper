package get_by_guid

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/usecase/server/item/get_by_guid"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

// Handle метод получение элемента по его GUID аутентифицированного пользователя.
func (h *Handler) Handle(ctx context.Context, in *pb.GetByGuidIn) (*pb.GetByGuidOut, error) {
	log := logger.FromCtx(ctx)

	userGUID := auth.UserGUIDFromCtx(ctx)
	if userGUID == uuid.Nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	itemGUID, err := uuid.Parse(in.GetGuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.usecase.Do(ctx, get_by_guid.Data{
		GUID:     itemGUID,
		UserGUID: userGUID,
	})
	if err != nil {
		if !errors.Is(err, get_by_guid.ErrNoData) {
			log.Error("Failed to get user item by guid", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}

		return nil, status.Error(codes.NotFound, "item not found")
	}

	return &pb.GetByGuidOut{
		Guid:          result.GUID.String(),
		EncryptedData: result.EncryptedData,
		CreatedAt:     timestamppb.New(result.CreatedAt),
		UpdatedAt:     timestamppb.New(result.UpdatedAt),
	}, nil
}
