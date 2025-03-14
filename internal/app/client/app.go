package client

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"

	formCreate "github.com/bjlag/go-keeper/internal/cli/model/create"
	"github.com/bjlag/go-keeper/internal/cli/model/item/bank_card"
	"github.com/bjlag/go-keeper/internal/cli/model/item/file"
	"github.com/bjlag/go-keeper/internal/cli/model/item/password"
	syncItem "github.com/bjlag/go-keeper/internal/cli/model/item/sync"
	"github.com/bjlag/go-keeper/internal/cli/model/item/text"
	"github.com/bjlag/go-keeper/internal/cli/model/list"
	formLogin "github.com/bjlag/go-keeper/internal/cli/model/login"
	"github.com/bjlag/go-keeper/internal/cli/model/master"
	formRegister "github.com/bjlag/go-keeper/internal/cli/model/register"
	"github.com/bjlag/go-keeper/internal/fetcher/item"
	crypt "github.com/bjlag/go-keeper/internal/infrastructure/crypt/cipher"
	"github.com/bjlag/go-keeper/internal/infrastructure/crypt/master_key"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/sqlite"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
	sItem "github.com/bjlag/go-keeper/internal/infrastructure/store/client/item"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/option"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/create"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/edit"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/remove"
	itemSync "github.com/bjlag/go-keeper/internal/usecase/client/item/sync"
	"github.com/bjlag/go-keeper/internal/usecase/client/login"
	mkey "github.com/bjlag/go-keeper/internal/usecase/client/master_key"
	"github.com/bjlag/go-keeper/internal/usecase/client/register"
	"github.com/bjlag/go-keeper/internal/usecase/client/sync"
)

const (
	saltLength         = 16
	masterKeyIterCount = 100_000
	masterKeyLength    = 256 / 8
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

	// TODO базу создавать и подключаться после успешного логин
	// TODO название файла базы должно быть уникальным под каждую учетку под которой авторизовались
	db, err := sqlite.New("./client.db").Connect()
	if err != nil {
		a.log.Error("failed to open db", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	salter := master_key.NewSaltGenerator(saltLength)
	keymaker := master_key.NewKeyGenerator(masterKeyIterCount, masterKeyLength)
	cipher := new(crypt.Cipher)

	storeItem := sItem.NewStore(db)
	storeOption := option.NewStore(db)

	ucLogin := login.NewUsecase(rpcClient, storeTokens)
	ucRegister := register.NewUsecase(rpcClient, storeTokens)
	ucMasterKey := mkey.NewUsecase(storeTokens, storeOption, salter, keymaker)
	ucSync := sync.NewUsecase(rpcClient, storeItem, storeTokens, cipher)
	ucItemSync := itemSync.NewUsecase(rpcClient, storeItem, storeTokens, cipher)
	ucCreateItem := create.NewUsecase(rpcClient, storeItem, storeTokens, cipher)
	ucSaveItem := edit.NewUsecase(rpcClient, storeItem, storeTokens, cipher)
	ucRemoveItem := remove.NewUsecase(rpcClient, storeItem)

	fetchItem := item.NewFetcher(storeItem)

	frmRegister := formRegister.InitModel(ucRegister, ucMasterKey)
	frmSync := syncItem.InitModel(ucItemSync)
	frmPasswordItem := password.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)
	frmTextItem := text.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)
	frmBankCardItem := bank_card.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)
	frmFileItem := file.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)

	m := master.InitModel(
		master.WithLoginForm(formLogin.InitModel(ucLogin, ucMasterKey, frmRegister)),
		master.WithCreatForm(formCreate.InitModel(frmPasswordItem, frmTextItem, frmBankCardItem, frmFileItem)),
		master.WithListForm(list.InitModel(ucSync, fetchItem, frmPasswordItem, frmTextItem, frmBankCardItem, frmFileItem)),
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
