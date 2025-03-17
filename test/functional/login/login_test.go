package login_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/infrastructure/migrator"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/server"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/user"
	rpcLogin "github.com/bjlag/go-keeper/internal/rpc/login"
	"github.com/bjlag/go-keeper/internal/usecase/server/user/login"
	"github.com/bjlag/go-keeper/test/config"
	"github.com/bjlag/go-keeper/test/infrastructure/container"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

type TestSuite struct {
	suite.Suite
	pgContainer *container.PostgreSQLContainer
	db          *sqlx.DB
	conn        *grpc.ClientConn
	client      rpc.KeeperClient
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	const pathToConfig = "./config/server_test.yaml"

	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	// load config
	var cfg config.Config
	err := cleanenv.ReadConfig(pathToConfig, &cfg)
	s.Require().NoError(err)

	// crate db container
	pgContainer, err := container.NewPostgreSQLContainer(ctx, container.PostgreSQLConfig{
		Database: cfg.Container.PG.DBName,
		Username: cfg.Container.PG.DBUser,
		Password: cfg.Container.PG.DBPassword,
		ImageTag: cfg.Container.PG.Tag,
	})
	s.Require().NoError(err)

	s.pgContainer = pgContainer

	// get db connection
	db, err := pg.New(pg.GetDSN(pgContainer.Host, pgContainer.Port, pgContainer.Database, pgContainer.Username, pgContainer.Password)).Connect()
	s.Require().NoError(err)

	s.db = db

	// apply migrations
	m, err := migrator.Get(db, migrator.TypePG, pgContainer.Database, cfg.Migration.SourcePath, cfg.Migration.Table)
	s.Require().NoError(err)

	err = m.Up()
	s.Require().NoError(err)

	// create grpc client
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(ctx, db)))
	s.Require().NoError(err)

	s.client = rpc.NewKeeperClient(conn)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.pgContainer.Terminate(ctx))

	err := s.conn.Close()
	s.Require().NoError(err)

	err = s.db.Close()
	s.Require().NoError(err)
}

func dialer(ctx context.Context, db *sqlx.DB) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1204 * 1024)

	jwt := auth.NewJWT("secret", 15*time.Minute, 15*time.Minute)

	userStore := user.NewStore(db)
	ucLogin := login.NewUsecase(userStore, jwt)

	server := server.NewRPCServer(
		server.WithListener(listener),
		server.WithLogger(logger.Get("test")),

		server.WithHandler(server.LoginMethod, rpcLogin.New(ucLogin).Handle),
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

func (s *TestSuite) TestHandler_Handle() {
	ctx := context.Background()

	err := fixture.Load(s.db, "test/fixture")
	s.Require().NoError(err)

	s.Run("success", func() {
		out, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "12345678",
		})

		if err != nil {
		}

		out.GetAccessToken()
	})
}
