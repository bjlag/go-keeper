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
	container2 "github.com/bjlag/go-keeper/test/util/container"
	"github.com/go-testfixtures/testfixtures/v3"
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
	_ "github.com/bjlag/go-keeper/test/util/init"
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
	ctx := context.Background()

	// todo из конфига берем данные для базы

	const (
		migrationsSourcePath = "./migrations/server"
		migrationsTable      = "migrations"
	)

	pgContainer, err := container2.NewPostgreSQLContainer(ctx, container2.PostgreSQLConfig{
		Database: "master_test",
		Username: "postgres",
		Password: "secret",
		ImageTag: "16.4-alpine3.20",
	})
	if err != nil {
		log.Fatal(err)
	}

	db, err := pg.New(pg.GetDSN(pgContainer.Host, pgContainer.Port, pgContainer.Database, pgContainer.Username, pgContainer.Password)).Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	m, err := migrator.Get(db, migrator.TypePG, pgContainer.Database, migrationsSourcePath, migrationsTable)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("test/fixture"),
	)

	require.NoError(t, err)
	require.NoError(t, fixtures.Load())

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
