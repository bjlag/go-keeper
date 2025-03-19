package server_test

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

func (s *TestSuite) TestUpdateItem() {
	err := fixture.Load(s.db, "test/fixture/server")
	s.Require().NoError(err)

	s.Run("success", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.UpdateItem(ctx, &rpc.UpdateItemIn{
			Guid:          "60308368-7729-4d2d-a510-67926f5a159b",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Version:       time.Date(2025, time.March, 15, 13, 0, 0, 0, time.UTC).UnixMicro(),
		})
		s.Require().NoError(err)

		item := s.getFromDBByGUID(ctx, "60308368-7729-4d2d-a510-67926f5a159b")

		s.Equal("60308368-7729-4d2d-a510-67926f5a159b", item.GUID)
		s.Equal([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, item.EncryptedData)
		s.Equal("bf4e6232-f1ae-41da-8535-73048891b1e3", item.UserGUID)
		s.True(time.Date(2025, time.March, 15, 13, 0, 0, 0, time.UTC).Equal(item.CreatedAt.UTC()))
		s.InDelta(time.Now().UTC().Unix(), item.UpdatedAt.UTC().Unix(), 2)
	})

	s.Run("permission denied", func() {
		_, err := s.client.UpdateItem(context.Background(), &rpc.UpdateItemIn{
			Guid:          "60308368-7729-4d2d-a510-67926f5a159b",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Version:       time.Date(2025, time.March, 15, 13, 0, 0, 0, time.UTC).UnixMicro(),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.PermissionDenied, st.Code())
		s.Equal("permission denied", st.Message())
	})

	s.Run("invalid item guid", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.UpdateItem(ctx, &rpc.UpdateItemIn{
			Guid:          "invalid guid",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Version:       time.Date(2025, time.March, 15, 13, 0, 0, 0, time.UTC).UnixMicro(),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("invalid item guid", st.Message())
	})

	s.Run("encrypted data is empty", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.UpdateItem(ctx, &rpc.UpdateItemIn{
			Guid:          "60308368-7729-4d2d-a510-67926f5a159b",
			EncryptedData: []byte{},
			Version:       time.Date(2025, time.March, 15, 13, 0, 0, 0, time.UTC).UnixMicro(),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("encrypted data is empty", st.Message())
	})

	s.Run("item belongs to other user", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.UpdateItem(ctx, &rpc.UpdateItemIn{
			Guid:          "d07f9605-0b8e-42f6-a07b-3e3f839a7bee",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Version:       time.Date(2025, time.March, 15, 12, 0, 0, 0, time.UTC).UnixMicro(),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.NotFound, st.Code())
		s.Equal("item not found", st.Message())
	})

	s.Run("item is outdated", func() {
		ctx := s.login(context.Background(), "test@test.ru", "12345678")

		_, err := s.client.UpdateItem(ctx, &rpc.UpdateItemIn{
			Guid:          "60308368-7729-4d2d-a510-67926f5a159b",
			EncryptedData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Version:       time.Date(2025, time.March, 15, 12, 0, 0, 0, time.UTC).UnixMicro(),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.FailedPrecondition, st.Code())
		s.Equal("item is outdated", st.Message())
	})
}
