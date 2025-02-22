package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/infrastructure/auth/jwt"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/server"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/user"
	rpcRegister "github.com/bjlag/go-keeper/internal/rpc/register"
	"github.com/bjlag/go-keeper/internal/usecase/user/register"
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
	tokeGenerator := jwt.NewGenerator(a.cfg.Auth.SecretKey, a.cfg.Auth.AccessTokenExp, a.cfg.Auth.RefreshTokenExp)

	ucRegister := register.NewUsecase(userStore, tokeGenerator)

	s := server.NewServer(
		server.WithAddress(a.cfg.Address.Host, a.cfg.Address.Port),
		server.WithLogger(a.log),
		server.WithHandler(server.RegisterMethod, rpcRegister.New(ucRegister).Handle),
	)

	err = s.Start(ctx)
	if err != nil {
		a.log.Error("Failed to start gRPC server", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("Server shutdown gracefully")

	return nil
}
