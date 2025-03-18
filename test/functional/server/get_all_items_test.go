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

	s.Run("success limit offset", func() {
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

		getAllItemsOut1, err := s.client.GetAllItems(metadata.NewOutgoingContext(ctx, md), &rpc.GetAllItemsIn{
			Offset: 0,
			Limit:  2,
		})
		s.Require().NoError(err)
		s.Assert().Len(getAllItemsOut1.GetItems(), 2)
		s.Assert().Equal(getAllItemsOut1.GetItems()[0].GetGuid(), "127e1a2d-1943-4fb1-ba60-7dc4fc820ed4")
		s.Assert().Equal(getAllItemsOut1.GetItems()[1].GetGuid(), "60308368-7729-4d2d-a510-67926f5a159b")

		getAllItemsOut2, err := s.client.GetAllItems(metadata.NewOutgoingContext(ctx, md), &rpc.GetAllItemsIn{
			Offset: 2,
			Limit:  2,
		})
		s.Require().NoError(err)
		s.Assert().Len(getAllItemsOut2.GetItems(), 2)
		s.Assert().Equal(getAllItemsOut2.GetItems()[0].GetGuid(), "6e7fc4fa-31aa-4d75-8b6e-0479122e0147")
		s.Assert().Equal(getAllItemsOut2.GetItems()[1].GetGuid(), "b2bd09eb-2c84-4149-b2b8-29040472264a")
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
