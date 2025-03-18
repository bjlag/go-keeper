package server_test

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
)

func (s *TestSuite) login(ctx context.Context, email, password string) context.Context {
	out, err := s.client.Login(ctx, &rpc.LoginIn{
		Email:    email,
		Password: password,
	})
	s.Require().NoError(err)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set("authorization", fmt.Sprintf("%s %s", "Bearer", out.GetAccessToken()))

	return metadata.NewOutgoingContext(ctx, md)
}
