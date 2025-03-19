// Package client содержит логику для работы с сервером через RPC протокол.
package client

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/interceptor"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
)

type RPCClient struct {
	conn   *grpc.ClientConn
	client rpc.KeeperClient
}

// NewRPCClient создает gRPC подключение к серверу расположенному на хосте serverHost и порте serverPort.
// tokensStore используется для интерцептора отвечающего за авторизацию клиента.
// log используется для интерцептора по логированию запросов.
func NewRPCClient(serverHost string, serverPort int, tokensStore *token.Store, log *zap.Logger) (*RPCClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", serverHost, serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.LoggerClientInterceptor(log),
			interceptor.AuthClientInterceptor(tokensStore),
		),
	)
	if err != nil {
		return nil, err
	}

	return &RPCClient{
		conn:   conn,
		client: rpc.NewKeeperClient(conn),
	}, nil
}

// Close закрывает подключение.
func (c RPCClient) Close() error {
	return c.conn.Close()
}
