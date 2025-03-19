package server_test

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

func (s *TestSuite) TestDeleteItem() {
	err := fixture.Load(s.db, "test/fixture/server")
	s.Require().NoError(err)

	s.Run("success", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.DeleteItem(ctx, &rpc.DeleteItemIn{
			Guid: "60308368-7729-4d2d-a510-67926f5a159b",
		})
		s.Require().NoError(err)

		items := s.getAllFromDBByUserGUID(ctx, "bf4e6232-f1ae-41da-8535-73048891b1e3")

		s.Len(items, 3)
		s.Condition(func() (success bool) {
			success = true
			for _, item := range items {
				if item.GUID == "60308368-7729-4d2d-a510-67926f5a159b" {
					success = false
				}
			}
			return
		}, "Не должно быть элемента с GUID 60308368-7729-4d2d-a510-67926f5a159b")
	})

	s.Run("permission denied", func() {
		_, err := s.client.DeleteItem(context.Background(), &rpc.DeleteItemIn{
			Guid: "60308368-7729-4d2d-a510-67926f5a159b",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.PermissionDenied, st.Code())
		s.Equal("permission denied", st.Message())
	})

	s.Run("not found", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.DeleteItem(ctx, &rpc.DeleteItemIn{
			Guid: "00000000-0000-0000-0000-000000000000",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.NotFound, st.Code())
		s.Equal("item not found", st.Message())
	})

	s.Run("not found if item belongs to other user", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.DeleteItem(ctx, &rpc.DeleteItemIn{
			Guid: "d07f9605-0b8e-42f6-a07b-3e3f839a7bee",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.NotFound, st.Code())
		s.Equal("item not found", st.Message())
	})
}
