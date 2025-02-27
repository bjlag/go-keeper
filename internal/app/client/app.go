package client

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/cli"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
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

	rpcClient, err := rpc.NewRPCClient(a.cfg.Server.Host, a.cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		_ = rpcClient.Close()
	}()

	_, err = tea.NewProgram(cli.InitModel(rpcClient), tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	return err
}
