package interceptor

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth/jwt"
)

const (
	AuthMeta   = "authorization"
	bearerAuth = "Bearer"
)

var methodSkip = map[string]struct{}{
	rpc.Keeper_Login_FullMethodName:         {},
	rpc.Keeper_Register_FullMethodName:      {},
	rpc.Keeper_RefreshTokens_FullMethodName: {},
}

func CheckAccessTokenInterceptor(jwtGenerator *jwt.Generator, log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := methodSkip[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		meta := md.Get(AuthMeta)
		if len(meta) == 0 {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		token, found := strings.CutPrefix(meta[0], bearerAuth)
		if !found {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		userGUID, err := jwtGenerator.GetUserGUIDFromAccessToken(strings.TrimLeft(token, " "))
		if err != nil {
			if errors.Is(err, jwt.ErrInvalidToken) {
				return nil, status.Errorf(codes.PermissionDenied, "permission denied")
			}
			log.Error("Failed to get user GUID", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "internal error")
		}

		return handler(auth.UserGUIDWithCtx(ctx, userGUID), req)
	}
}
