package server_test

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/test/functional/server"
)

//nolint:unparam
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

func (s *TestSuite) getFromDBByGUID(ctx context.Context, guid string) server.Item {
	query := `
		SELECT guid, user_guid, encrypted_data, created_at, updated_at
		FROM items
		WHERE guid = $1
	`

	var row server.Item
	err := s.db.GetContext(ctx, &row, query, guid)
	s.Require().NoError(err)

	return row
}

func (s *TestSuite) getAllFromDBByUserGUID(ctx context.Context, guid string) []server.Item {
	query := `
		SELECT guid, user_guid, encrypted_data, created_at, updated_at
		FROM items
		WHERE user_guid = $1
	`

	var rows []server.Item
	err := s.db.SelectContext(ctx, &rows, query, guid)
	s.Require().NoError(err)

	return rows
}
