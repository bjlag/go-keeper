package update_item

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/usecase/server/item/update"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.UpdateItemIn) (*pb.UpdateItemOut, error) {
	log := logger.FromCtx(ctx)

	userGUID := auth.UserGUIDFromCtx(ctx)
	if userGUID == uuid.Nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if len(in.GetEncryptedData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "encrypted data is empty")
	}

	itemGUID, err := uuid.Parse(in.GetGuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item guid")
	}

	newVersion, err := h.usecase.Do(ctx, update.In{
		UserGUID:      userGUID,
		ItemGUID:      itemGUID,
		EncryptedData: in.GetEncryptedData(),
		Version:       in.GetVersion(),
	})
	if err != nil {
		if errors.Is(err, update.ErrConflict) {
			return nil, status.Error(codes.FailedPrecondition, "data outdated")
		}

		if errors.Is(err, update.ErrNotFoundUpdatedData) {
			return nil, status.Error(codes.NotFound, "item not found")
		}

		log.Error("Failed to update item", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.UpdateItemOut{
		NewVersion: newVersion,
	}, nil
}
