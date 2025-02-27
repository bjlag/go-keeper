package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

type RPCClient struct {
	conn   *grpc.ClientConn
	client rpc.KeeperClient
}

func NewRPCClient(serverHost string, serverPort int) (*RPCClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", serverHost, serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithChainUnaryInterceptor(
		//	interceptor.LoggerClientInterceptor(log),
		//	interceptor.RealIPClientInterceptor,
		//	interceptor.SignatureClientInterceptor(sign),
		//),
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
