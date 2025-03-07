package create_item

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/usecase/server/item/create"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.CreateItemIn) (*pb.CreateItemOut, error) {
	log := logger.FromCtx(ctx)

	userGUID := auth.UserGUIDFromCtx(ctx)
	if userGUID == uuid.Nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if len(in.GetEncryptedData()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "encrypted data is empty")
	}

	itemGUID := uuid.New()

	err := h.usecase.Do(ctx, create.In{
		ItemGUID:      itemGUID,
		UserGUID:      userGUID,
		EncryptedData: in.GetEncryptedData(),
	})
	if err != nil {
		log.Error("Failed to update item", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.CreateItemOut{
		Guid: itemGUID.String(),
	}, nil
}
