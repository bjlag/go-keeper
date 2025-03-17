package login_test

import (
	"context"
	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/infrastructure/migrator"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/user"
	rpcLogin "github.com/bjlag/go-keeper/internal/rpc/login"
	"github.com/bjlag/go-keeper/internal/usecase/server/user/login"
	"github.com/bjlag/go-keeper/test/config"
	container2 "github.com/bjlag/go-keeper/test/infrastructure/container"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"

	server2 "github.com/bjlag/go-keeper/internal/infrastructure/rpc/server"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

func dialer(ctx context.Context, db *sqlx.DB) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1204 * 1024)

	jwt := auth.NewJWT("secret", 15*time.Minute, 15*time.Minute)

	userStore := user.NewStore(db)
	ucLogin := login.NewUsecase(userStore, jwt)

	server := server2.NewRPCServer(
		server2.WithListener(listener),
		server2.WithLogger(logger.Get("test")),

		server2.WithHandler(server2.LoginMethod, rpcLogin.New(ucLogin).Handle),
	)

	go func() {
		if err := server.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestHandler_Handle(t *testing.T) {
	const pathToConfig = "./config/server_test.yaml"

	ctx := context.Background()

	var cfg config.Config
	if err := cleanenv.ReadConfig(pathToConfig, &cfg); err != nil {
		panic(err)
	}

	pgContainer, err := container2.NewPostgreSQLContainer(ctx, container2.PostgreSQLConfig{
		Database: cfg.Container.PG.DBName,
		Username: cfg.Container.PG.DBUser,
		Password: cfg.Container.PG.DBPassword,
		ImageTag: cfg.Container.PG.Tag,
	})
	if err != nil {
		log.Fatal(err)
	}

	db, err := pg.New(pg.GetDSN(pgContainer.Host, pgContainer.Port, pgContainer.Database, pgContainer.Username, pgContainer.Password)).Connect()
	require.NoError(t, err)
	defer db.Close()

	m, err := migrator.Get(db, migrator.TypePG, pgContainer.Database, cfg.Migration.SourcePath, cfg.Migration.Table)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	err = fixture.Load(db, "test/fixture")
	require.NoError(t, err)

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(ctx, db)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := rpc.NewKeeperClient(conn)

	// todo выключение логгера
	// todo конфиги для тестов
	// todo тестовое приложение
	// todo своя база
	// todo фикстуры базы
	// todo test suit, tear down

	t.Run("success", func(t *testing.T) {
		out, err := client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "12345678",
		})

		if err != nil {
		}

		out.GetAccessToken()
	})
}
