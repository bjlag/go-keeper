package get_all_items

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
	"github.com/bjlag/go-keeper/internal/usecase/server/item/get_all"
)

const (
	limitDefault = 40
	limitMax     = 100
)

type Handler struct {
	usecase usecase
}

func New(usecase usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

// Handle метод получения всех элементов аутентифицированного пользователя.
func (h *Handler) Handle(ctx context.Context, in *pb.GetAllItemsIn) (*pb.GetAllItemsOut, error) {
	log := logger.FromCtx(ctx)

	userGUID := auth.UserGUIDFromCtx(ctx)
	if userGUID == uuid.Nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if in.GetLimit() > limitMax {
		in.Limit = limitMax
	}

	if in.GetLimit() <= 0 {
		in.Limit = limitDefault
	}

	if in.GetOffset() <= 0 {
		in.Offset = 0
	}

	result, err := h.usecase.Do(ctx, get_all.Data{
		UserGUID: userGUID,
		Limit:    in.GetLimit(),
		Offset:   in.GetOffset(),
	})
	if err != nil {
		if !errors.Is(err, get_all.ErrNoData) {
			log.Error("Failed to get all item", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	itemsOut := make([]*pb.Item, 0, len(result.Items))
	for _, item := range result.Items {
		itemsOut = append(itemsOut, &pb.Item{
			Guid:          item.GUID.String(),
			EncryptedData: item.EncryptedData,
			CreatedAt:     timestamppb.New(item.CreatedAt),
			UpdatedAt:     timestamppb.New(item.UpdatedAt),
		})
	}

	return &pb.GetAllItemsOut{
		Items: itemsOut,
	}, nil
}
