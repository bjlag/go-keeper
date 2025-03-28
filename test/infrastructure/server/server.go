package server

import (
	"context"
	"net"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc/test/bufconn"

	"github.com/bjlag/go-keeper/internal/app/server"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
)

const bufSize = 1024 * 1024

func Start(ctx context.Context, db *sqlx.DB, jwt *auth.JWT, log *zap.Logger) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(bufSize)

	go func() {
		if err := server.NewApp(db, jwt, listener, log).Run(ctx); err != nil {
			panic(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
