package auth

import (
	"context"

	"github.com/google/uuid"
)

type ctxKeyUserGUID int

const UserGUIDKey ctxKeyUserGUID = 0

func UserGUIDWithCtx(ctx context.Context, guid uuid.UUID) context.Context {
	if ctxGUID, ok := ctx.Value(UserGUIDKey).(string); ok {
		if ctxGUID == guid.String() {
			return ctx
		}
	}

	return context.WithValue(ctx, UserGUIDKey, guid.String())
}

func UserGUIDFromCtx(ctx context.Context) uuid.UUID {
	if s, ok := ctx.Value(UserGUIDKey).(string); ok {
		if guid, err := uuid.Parse(s); err == nil {
			return guid
		}
	}

	return uuid.Nil
}
