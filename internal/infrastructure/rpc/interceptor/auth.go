package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bjlag/go-keeper/internal/generated/rpc"
	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
	"github.com/bjlag/go-keeper/internal/infrastructure/store/client/token"
)

const (
	// authMeta мета заголовок, в котором лежит access токен.
	authMeta = "authorization"
	// bearerAuth метод аутентификации.
	bearerAuth = "Bearer"
)

// methodSkip содержит gRPC методы, у которых не надо проверять аутентификацию.
var methodSkip = map[string]struct{}{
	rpc.Keeper_Login_FullMethodName:         {},
	rpc.Keeper_Register_FullMethodName:      {},
	rpc.Keeper_RefreshTokens_FullMethodName: {},
}

// CheckAccessTokenServerInterceptor клиентский интерцептор, который перед запросом к серверу кладет мета заголовок
// в запрос с access токеном, чтобы сервер мог проверить аутентифицирован пользователь или нет.
func CheckAccessTokenServerInterceptor(jwt *auth.JWT, log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := methodSkip[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		meta := md.Get(authMeta)
		if len(meta) == 0 {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		accessToken, found := strings.CutPrefix(meta[0], bearerAuth)
		if !found {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		userGUID, err := jwt.GetUserGUIDFromAccessToken(strings.TrimLeft(accessToken, " "))
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) {
				return nil, status.Errorf(codes.PermissionDenied, "permission denied")
			}
			log.Error("Failed to get user GUID", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "internal error")
		}

		return handler(auth.UserGUIDWithCtx(ctx, userGUID), req)
	}
}

// AuthClientInterceptor интерцептор, который на стороне сервера, перед исполнением запроса, проверяет
// по переданной информации с клиента, аутентифицирован пользователь или нет.
func AuthClientInterceptor(tokens *token.Store) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if _, ok := methodSkip[method]; ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, authMeta, fmt.Sprintf("%s %s", bearerAuth, tokens.AccessToken()))
		err := invoker(ctx, method, req, reply, cc, opts...)
		if status.Code(err) == codes.PermissionDenied {
			in := &rpc.RefreshTokensIn{
				RefreshToken: tokens.RefreshToken(),
			}
			out := &rpc.RefreshTokensOut{}
			err = cc.Invoke(ctx, rpc.Keeper_RefreshTokens_FullMethodName, in, out)
			if err != nil {
				if status.Code(err) == codes.FailedPrecondition {
					return status.Errorf(codes.PermissionDenied, err.Error())
				}
				return err
			}

			tokens.SaveTokens(out.GetAccessToken(), out.GetRefreshToken())

			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.New(nil)
			}

			md.Set(authMeta, fmt.Sprintf("%s %s", bearerAuth, tokens.AccessToken()))
			err = invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
		}

		return err
	}
}
