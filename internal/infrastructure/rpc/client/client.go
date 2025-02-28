package client

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/rpc/interceptor"
)

type RPCClient struct {
	conn   *grpc.ClientConn
	client rpc.KeeperClient
}

func NewRPCClient(serverHost string, serverPort int, log *zap.Logger) (*RPCClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", serverHost, serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.LoggerClientInterceptor(log),
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

func (c RPCClient) Close() error {
	return c.conn.Close()
}
