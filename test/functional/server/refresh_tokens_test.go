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

func (s *TestSuite) TestRefreshTokens() {
	err := fixture.Load(s.db, "test/fixture/server")
	s.Require().NoError(err)

	s.Run("success", func() {
		ctx := context.Background()
		loginOut, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "12345678",
		})
		s.Require().NoError(err)

		refreshTokensOut, err := s.client.RefreshTokens(ctx, &rpc.RefreshTokensIn{
			RefreshToken: loginOut.GetRefreshToken(),
		})
		s.Require().NoError(err)
		s.NotEmpty(refreshTokensOut.GetAccessToken())
		s.NotEmpty(refreshTokensOut.GetRefreshToken())

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		md.Set("authorization", fmt.Sprintf("%s %s", "Bearer", refreshTokensOut.GetAccessToken()))

		getAllItemsOut, err := s.client.GetAllItems(metadata.NewOutgoingContext(ctx, md), &rpc.GetAllItemsIn{})
		s.Require().NoError(err)
		s.Len(getAllItemsOut.GetItems(), 4)
	})

	s.Run("sent access token", func() {
		ctx := context.Background()
		loginOut, err := s.client.Login(ctx, &rpc.LoginIn{
			Email:    "test@test.ru",
			Password: "12345678",
		})
		s.Require().NoError(err)

		refreshTokensOut, err := s.client.RefreshTokens(ctx, &rpc.RefreshTokensIn{
			RefreshToken: loginOut.GetAccessToken(),
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.FailedPrecondition, st.Code())
		s.Equal("invalid refresh token", st.Message())
		s.Nil(refreshTokensOut)
	})

	s.Run("refresh token too short", func() {
		ctx := context.Background()
		refreshTokensOut, err := s.client.RefreshTokens(ctx, &rpc.RefreshTokensIn{
			RefreshToken: "refresh token too short",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.InvalidArgument, st.Code())
		s.Equal("refresh token too short", st.Message())
		s.Nil(refreshTokensOut)
	})

	s.Run("refresh token is wrong", func() {
		ctx := context.Background()
		refreshTokensOut, err := s.client.RefreshTokens(ctx, &rpc.RefreshTokensIn{
			RefreshToken: "wrong_token_eyJhbGciOiJIUzI1NiIsInR5.eyJpc3MiOiJiZjRlNjIzMi1mMWFlLTQxZGEtODUzNS03MzA0ODg5MWIxZTMiLCJzdWIiOiJyZWZyZXNoX3Rva2VuIiwiZXhwIjoxNzQyMzE2MDQ1LCJpYXQiOjE3NDIzMDg4NDV9.YO1nbCZ-1r4u4BWBvMSPOAda5m9R_IvoTJNqsziU-ik",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.FailedPrecondition, st.Code())
		s.Equal("invalid refresh token", st.Message())
		s.Nil(refreshTokensOut)
	})

	s.Run("refresh token is expired", func() {
		ctx := context.Background()
		refreshTokensOut, err := s.client.RefreshTokens(ctx, &rpc.RefreshTokensIn{
			RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJiZjRlNjIzMi1mMWFlLTQxZGEtODUzNS03MzA0ODg5MWIxZTMiLCJzdWIiOiJyZWZyZXNoX3Rva2VuIiwiZXhwIjoxNzQyMzA5MTQwLCJpYXQiOjE3NDIzMDkxMzB9.zd-cbzbhry50DmC94xjxeBxXpX6HkHurgLpVorA9lg0",
		})

		st, ok := status.FromError(err)
		s.True(ok)
		s.Equal(codes.FailedPrecondition, st.Code())
		s.Equal("invalid refresh token", st.Message())
		s.Nil(refreshTokensOut)
	})
}
