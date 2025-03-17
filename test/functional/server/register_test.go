package server_test

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

func (s *TestSuite) TestRegister() {
	ctx := context.Background()

	err := fixture.Load(s.db, "test/fixture/server")
	s.Require().NoError(err)

	s.Run("success", func() {
		out, err := s.client.Register(ctx, &rpc.RegisterIn{
			Email:    "new@test.ru",
			Password: "12345678",
		})

		s.NoError(err)
		s.NotEmpty(out.GetAccessToken())
		s.NotEmpty(out.GetRefreshToken())
	})

	s.Run("already exists", func() {
		out, err := s.client.Register(ctx, &rpc.RegisterIn{
			Email:    "test@test.ru",
			Password: "12345678",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.AlreadyExists, st.Code())
		s.Equal("user with this email already exists", st.Message())

		s.Empty(out.GetAccessToken())
		s.Empty(out.GetRefreshToken())
	})

	s.Run("invalid email", func() {
		out, err := s.client.Register(ctx, &rpc.RegisterIn{
			Email:    "test@test",
			Password: "12345678",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("email is invalid", st.Message())

		s.Empty(out.GetAccessToken())
		s.Empty(out.GetRefreshToken())
	})

	s.Run("invalid password", func() {
		out, err := s.client.Register(ctx, &rpc.RegisterIn{
			Email:    "test@test.ru",
			Password: "1234567",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("password is invalid (min. length 8 characters)", st.Message())

		s.Empty(out.GetAccessToken())
		s.Empty(out.GetRefreshToken())
	})
}
