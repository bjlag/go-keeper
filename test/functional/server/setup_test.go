package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	"github.com/bjlag/go-keeper/internal/infrastructure/migrator"
	"github.com/bjlag/go-keeper/test/infrastructure/config"
	"github.com/bjlag/go-keeper/test/infrastructure/container"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
	"github.com/bjlag/go-keeper/test/infrastructure/server"
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
	jwt := auth.NewJWT(cfg.Auth.SecretKey, cfg.Auth.AccessTokenExp, cfg.Auth.RefreshTokenExp)
	log := logger.Get(cfg.Env)

	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(server.Start(context.Background(), db, jwt, log)),
	)
	s.Require().NoError(err)

	s.conn = conn
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
