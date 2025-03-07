package server

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/server"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/item"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/user"
	rpcDeleteItem "github.com/bjlag/go-keeper/internal/rpc/delete_item"
	rpcGetAllItems "github.com/bjlag/go-keeper/internal/rpc/get_all_items"
	rpcLogin "github.com/bjlag/go-keeper/internal/rpc/login"
	rpcRefreshTokens "github.com/bjlag/go-keeper/internal/rpc/refresh_tokens"
	rpcRegister "github.com/bjlag/go-keeper/internal/rpc/register"
	rpcUpdateItem "github.com/bjlag/go-keeper/internal/rpc/update_item"
	"github.com/bjlag/go-keeper/internal/usecase/server/data/delete"
	"github.com/bjlag/go-keeper/internal/usecase/server/data/get_all"
	"github.com/bjlag/go-keeper/internal/usecase/server/data/update"
	"github.com/bjlag/go-keeper/internal/usecase/server/user/login"
	rt "github.com/bjlag/go-keeper/internal/usecase/server/user/refresh_tokens"
	"github.com/bjlag/go-keeper/internal/usecase/server/user/register"
)

type App struct {
	cfg Config
	log *zap.Logger
}

func NewApp(cfg Config, log *zap.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "app.Run"

	dbConf := a.cfg.Database
	db, err := pg.New(pg.GetDSN(dbConf.Host, dbConf.Port, dbConf.Name, dbConf.User, dbConf.Password)).Connect()
	if err != nil {
		a.log.Error("Failed to get db connection", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = db.Close()
	}()

	userStore := user.NewStore(db)
	dataStore := item.NewStore(db)
	jwt := auth.NewJWT(a.cfg.Auth.SecretKey, a.cfg.Auth.AccessTokenExp, a.cfg.Auth.RefreshTokenExp)

	ucRegister := register.NewUsecase(userStore, jwt)
	ucLogin := login.NewUsecase(userStore, jwt)
	ucRefreshTokens := rt.NewUsecase(userStore, jwt)
	ucGetAllData := get_all.NewUsecase(dataStore)
	ucUpdateItem := update.NewUsecase(dataStore)
	ucDeleteItem := delete.NewUsecase(dataStore)

	s := server.NewRPCServer(
		server.WithAddress(a.cfg.Address.Host, a.cfg.Address.Port),
		server.WithJWT(jwt),
		server.WithLogger(a.log),

		server.WithHandler(server.RegisterMethod, rpcRegister.New(ucRegister).Handle),
		server.WithHandler(server.LoginMethod, rpcLogin.New(ucLogin).Handle),
		server.WithHandler(server.RefreshTokensMethod, rpcRefreshTokens.New(ucRefreshTokens).Handle),
		server.WithHandler(server.GetAllItemsMethod, rpcGetAllItems.New(ucGetAllData).Handle),
		server.WithHandler(server.UpdateItemMethod, rpcUpdateItem.New(ucUpdateItem).Handle),
		server.WithHandler(server.DeleteItemMethod, rpcDeleteItem.New(ucDeleteItem).Handle),
	)

	err = s.Start(ctx)
	if err != nil {
		a.log.Error("Failed to start gRPC server", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("Server shutdown gracefully")

	return nil
}
