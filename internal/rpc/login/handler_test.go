package login_test

import (
	"context"
	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/db/pg"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
	server2 "github.com/bjlag/go-keeper/internal/infrastructure/rpc/server"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/server/user"
	rpcLogin "github.com/bjlag/go-keeper/internal/rpc/login"
	"github.com/bjlag/go-keeper/internal/usecase/server/user/login"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
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

	db, err := pg.New(pg.GetDSN("localhost", "5444", "master", "postgres", "secret")).Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(ctx, db)))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := rpc.NewKeeperClient(conn)

	//

	t.Run("success", func(t *testing.T) {
		out, err := client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.com",
			Password: "test",
		})

		if err != nil {
		}

		out.GetAccessToken()
	})
}
