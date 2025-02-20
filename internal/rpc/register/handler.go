package register

import (
	"context"

	pb "github.com/bjlag/go-keeper/internal/generated/rpc"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, in *pb.RegisterIn) (*pb.RegisterOut, error) {
	// todo проверить есть ли пользователь с указанным емейлом
	// todo проверить политику паролей
	// todo захешировать пароль
	// todo добавить пользователя в базу
	// todo выпустить токены
	// todo вернуть токены

	return &pb.RegisterOut{
		AccessToken:  "xxxx",
		RefreshToken: "yyyy",
	}, nil
}
