package server_test

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

func (s *TestSuite) TestGetAllItems() {

	err := fixture.Load(s.db, "test/fixture/server")
	s.Require().NoError(err)

	s.Run("success", func() {
		ctx := context.Background()
		loginOut, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "12345678",
		})
		s.Require().NoError(err)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		md.Set("authorization", fmt.Sprintf("%s %s", "Bearer", loginOut.GetAccessToken()))

		getAllItemsOut, err := s.client.GetAllItems(metadata.NewOutgoingContext(ctx, md), &rpc.GetAllItemsIn{})
		s.Require().NoError(err)
		s.Assert().Len(getAllItemsOut.GetItems(), 4)
	})

	s.Run("permission denied", func() {
		ctx := context.Background()
		out, err := s.client.GetAllItems(ctx, &rpc.GetAllItemsIn{})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.PermissionDenied, st.Code())
		s.Equal("permission denied", st.Message())

		s.Empty(out.GetItems())
	})
}
