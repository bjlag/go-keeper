// Package client настраивает и запускает клиент.
package client

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
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
	"github.com/bjlag/go-keeper/internal/infrastructure/migrator"
	rpc "github.com/bjlag/go-keeper/internal/infrastructure/rpc/client"
	backups "github.com/bjlag/go-keeper/internal/infrastructure/store/client/backup"
	sItem "github.com/bjlag/go-keeper/internal/infrastructure/store/client/item"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/option"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
	"github.com/bjlag/go-keeper/internal/usecase/client/backup"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/create"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/edit"
	"github.com/bjlag/go-keeper/internal/usecase/client/item/remove"
	itemSync "github.com/bjlag/go-keeper/internal/usecase/client/item/sync"
	"github.com/bjlag/go-keeper/internal/usecase/client/login"
	mkey "github.com/bjlag/go-keeper/internal/usecase/client/master_key"
	"github.com/bjlag/go-keeper/internal/usecase/client/register"
	"github.com/bjlag/go-keeper/internal/usecase/client/restore"
	"github.com/bjlag/go-keeper/internal/usecase/client/sync"
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

	tokens := token.NewStore()

	client, err := rpc.NewRPCClient(a.cfg.Server.Host, a.cfg.Server.Port, tokens, a.log)
	if err != nil {
		a.log.Error("Failed to create rpc client", zap.Error(err))
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		_ = client.Close()
	}()

	email, pass, err := a.login(client, tokens)
	if err != nil {
		a.log.Error("Failed to login", zap.Error(err))
		return fmt.Errorf("%s:%w", op, err)
	}
	if email == "" {
		return nil
	}

	db, err := a.initDB(email)
	if err != nil {
		a.log.Error("Failed to init db", zap.Error(err))
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		_ = db.Close()
	}()

	salter := master_key.NewSaltGenerator(a.cfg.MasterKey.SaltLength)
	keymaker := master_key.NewKeyGenerator(a.cfg.MasterKey.IterCount, a.cfg.MasterKey.Length)
	cipher := new(crypt.Cipher)

	storeItem := sItem.NewStore(db)
	storeOption := option.NewStore(db)
	storeBackup := backups.NewStore(db)

	ucMasterKey := mkey.NewUsecase(tokens, storeOption, salter, keymaker)
	err = ucMasterKey.Do(ctx, mkey.Data{Password: pass})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	ucSync := sync.NewUsecase(client, storeItem, tokens, cipher)
	ucItemSync := itemSync.NewUsecase(client, storeItem, tokens, cipher)
	ucCreateItem := create.NewUsecase(client, storeItem, tokens, cipher)
	ucSaveItem := edit.NewUsecase(client, storeItem, tokens, cipher)
	ucRemoveItem := remove.NewUsecase(client, storeItem)
	ucBackup := backup.NewUsecase(storeItem, tokens, storeBackup, cipher)
	ucRestore := restore.NewUsecase(storeItem, tokens, storeBackup, cipher)

	fetchItem := item.NewFetcher(storeItem)

	frmSync := syncItem.InitModel(ucItemSync)
	frmPasswordItem := password.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)
	frmTextItem := text.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)
	frmBankCardItem := bank_card.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)
	frmFileItem := file.InitModel(ucCreateItem, ucSaveItem, ucRemoveItem, frmSync)

	err = ucRestore.Do(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	m := master.InitModel(
		master.WithCreatForm(formCreate.InitModel(frmPasswordItem, frmTextItem, frmBankCardItem, frmFileItem)),
		master.WithListForm(list.InitModel(ucSync, fetchItem, frmPasswordItem, frmTextItem, frmBankCardItem, frmFileItem)),
	)
	defer func() {
		err = ucBackup.Do(ctx)
		if err != nil {
			a.log.Error("Failed to backup", zap.Error(err))
		}
	}()

	_, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	if err != nil {
		a.log.Error("failed to run cli program", zap.Error(err))
	}

	return err
}

func (a *App) login(client *rpc.RPCClient, tokens *token.Store) (email string, pass string, err error) {
	ucLogin := login.NewUsecase(client, tokens)
	ucRegister := register.NewUsecase(client, tokens)

	frmRegister := formRegister.InitModel(ucRegister)
	frmLogin := formLogin.InitModel(ucLogin, frmRegister)

	mLogin, err := tea.NewProgram(frmLogin, tea.WithAltScreen()).Run()
	if err != nil {
		a.log.Error("Failed to run cli program for login", zap.Error(err))
		return
	}

	ml, ok := mLogin.(*formLogin.Model)
	if !ok {
		panic("failed to run model cli program")
	}

	email = ml.UserEmail()
	pass = ml.UserPass()

	return
}

func (a *App) initDB(email string) (db *sqlx.DB, err error) {
	emailHash := md5.Sum([]byte(email)) //nolint:gosec
	dbName := fmt.Sprintf("%s_%x.db", a.cfg.Database.Prefix, emailHash)
	pathToDB := path.Join(a.cfg.Database.Dir, dbName)

	db, err = sqlite.New(pathToDB).Connect()
	if err != nil {
		return
	}

	m, err := migrator.Get(db, migrator.TypeSqlite, "", a.cfg.Migration.SourcePath, a.cfg.Migration.Table)
	if err != nil {
		return
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return
		}
		err = nil
	}

	return
}
