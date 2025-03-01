package client

import (
	"context"
	"fmt"
	"github.com/bjlag/go-keeper/internal/cli/form/list"
	"github.com/bjlag/go-keeper/internal/cli/form/password"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/cli"
	formLogin "github.com/bjlag/go-keeper/internal/cli/form/login"
	formRegister "github.com/bjlag/go-keeper/internal/cli/form/register"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
	"github.com/bjlag/go-keeper/internal/usecase/client/login"
	"github.com/bjlag/go-keeper/internal/usecase/client/register"
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

	rpcClient, err := rpc.NewRPCClient(a.cfg.Server.Host, a.cfg.Server.Port, a.log)
	if err != nil {
		a.log.Error("failed to create rpc client", zap.Error(err))
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		_ = rpcClient.Close()
	}()

	ucLogin := login.NewUsecase(rpcClient)
	ucRegister := register.NewUsecase(rpcClient)

	model := cli.InitModel(
		cli.WithLoginForm(formLogin.NewForm(ucLogin)),
		cli.WithRegisterForm(formRegister.NewForm(ucRegister)),
		cli.WithListFormForm(list.NewForm()),
		cli.WithShowPasswordForm(password.NewForm()),
	)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = tea.NewProgram(model, tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	if err != nil {
		a.log.Error("failed to run cli program", zap.Error(err))
	}

	return err
}
