package client

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/cli/model/item"
	"github.com/bjlag/go-keeper/internal/cli/model/list"
	formLogin "github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/model/master"
	formRegister "github.com/bjlag/go-keeper/internal/cli/model/register"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/sqlite"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
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

	storeTokens := token.NewStore()

	rpcClient, err := rpc.NewRPCClient(a.cfg.Server.Host, a.cfg.Server.Port, storeTokens, a.log)
	if err != nil {
		a.log.Error("failed to create rpc client", zap.Error(err))
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		_ = rpcClient.Close()
	}()

	db, err := sqlite.New("./client.db").Connect()
	if err != nil {
		a.log.Error("failed to open db", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = db

	ucLogin := login.NewUsecase(rpcClient)
	ucRegister := register.NewUsecase(rpcClient)

	m := master.InitModel(
		master.WithStoreTokens(storeTokens),

		master.WithLoginForm(formLogin.InitModel(ucLogin)),
		master.WithRegisterForm(formRegister.InitModel(ucRegister)),
		master.WithListFormForm(list.InitModel(rpcClient)),
		master.WithShowPasswordForm(item.InitModel()),
	)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	if err != nil {
		a.log.Error("failed to run cli program", zap.Error(err))
	}

	return err
}
