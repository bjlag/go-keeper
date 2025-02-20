package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/server"
	"github.com/bjlag/go-keeper/internal/rpc/register"
)

type App struct {
	log *zap.Logger
}

func NewApp(log *zap.Logger) *App {
	return &App{
		log: log,
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "app.Run"

	server := server.NewServer(
		server.WithLogger(a.log),
		server.WithHandler(server.RegisterMethod, register.New().Handle),
	)

	err := server.Start(ctx)
	if err != nil {
		a.log.Error("Failed to start gRPC server", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("Server shutdown gracefully")

	return nil
}
