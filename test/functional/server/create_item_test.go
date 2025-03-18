package server_test

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

func (s *TestSuite) TestCreateItem() {
	err := fixture.Load(s.db, "test/fixture/server")
	s.Require().NoError(err)

	s.Run("success", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.CreateItem(ctx, &rpc.CreateItemIn{
			Guid:          "c904fe47-4ec8-4ff4-a642-774a7bf4e351",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			CreatedAt:     timestamppb.New(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
		})
		s.Require().NoError(err)

		item := s.getFromDBByGUID(ctx, "c904fe47-4ec8-4ff4-a642-774a7bf4e351")

		s.Equal("c904fe47-4ec8-4ff4-a642-774a7bf4e351", item.GUID)
		s.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, item.EncryptedData)
		s.Equal("bf4e6232-f1ae-41da-8535-73048891b1e3", item.UserGUID)
		s.True(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC).Equal(item.CreatedAt.UTC()))
		s.True(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC).Equal(item.UpdatedAt.UTC()))
	})

	s.Run("permission denied", func() {
		_, err := s.client.CreateItem(context.Background(), &rpc.CreateItemIn{
			Guid:          "c904fe47-4ec8-4ff4-a642-774a7bf4e351",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			CreatedAt:     timestamppb.New(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.PermissionDenied, st.Code())
		s.Equal("permission denied", st.Message())
	})

	s.Run("invalid item guid", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.CreateItem(ctx, &rpc.CreateItemIn{
			Guid:          "invalid guid",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			CreatedAt:     timestamppb.New(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("invalid item guid", st.Message())
	})

	s.Run("encrypted data is empty", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.CreateItem(ctx, &rpc.CreateItemIn{
			Guid:          "c904fe47-4ec8-4ff4-a642-774a7bf4e351",
			EncryptedData: []byte{},
			CreatedAt:     timestamppb.New(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("encrypted data is empty", st.Message())
	})
}
