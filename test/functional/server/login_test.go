package server_test

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/test/infrastructure/fixture"
	_ "github.com/bjlag/go-keeper/test/infrastructure/init"
)

func (s *TestSuite) TestLogin() {
	err := fixture.Load(s.db, "test/fixture/server")
	s.Require().NoError(err)

	s.Run("success", func() {
		ctx := context.Background()
		out, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "12345678",
		})

		s.NoError(err)
		s.NotEmpty(out.GetAccessToken())
		s.NotEmpty(out.GetRefreshToken())
	})

	s.Run("wrong credentials", func() {
		ctx := context.Background()
		out, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "1111111",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.Unauthenticated, st.Code())
		s.Equal("credentials incorrect", st.Message())

		s.Empty(out.GetAccessToken())
		s.Empty(out.GetRefreshToken())
	})

	s.Run("empty body", func() {
		ctx := context.Background()
		out, err := s.client.Login(ctx, &rpc.LoginIn{})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("email is invalid", st.Message())

		s.Empty(out.GetAccessToken())
		s.Empty(out.GetRefreshToken())
	})

	s.Run("wrong email", func() {
		ctx := context.Background()
		out, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test",
			Password: "12345678",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal(st.Message(), "email is invalid")

		s.Empty(out.GetAccessToken())
		s.Empty(out.GetRefreshToken())
	})

	s.Run("wrong password", func() {
		ctx := context.Background()
		out, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("password is empty", st.Message())

		s.Empty(out.GetAccessToken())
		s.Empty(out.GetRefreshToken())
	})
}
