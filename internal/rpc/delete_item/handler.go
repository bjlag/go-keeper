package delete_item

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/usecase/server/item/delete"
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(ctx context.Context, in *pb.DeleteItemIn) (*emptypb.Empty, error) {
	log := logger.FromCtx(ctx)

	userGUID := auth.UserGUIDFromCtx(ctx)
	if userGUID == uuid.Nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	itemGUID, err := uuid.Parse(in.GetGuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item guid")
	}

	err = h.usecase.Do(ctx, delete.In{
		UserGUID: userGUID,
		ItemGUID: itemGUID,
	})
	if err != nil {
		if errors.Is(err, delete.ErrNotFoundUpdatedData) {
			return nil, status.Error(codes.NotFound, "item not found")
		}

		log.Error("Failed to update item", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}
