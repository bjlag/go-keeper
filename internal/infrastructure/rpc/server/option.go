package server

import (
	"go.uber.org/zap"

	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
)

type Option func(*RPCServer)

func WithAddress(host string, port int) Option {
	return func(s *RPCServer) {
		s.host = host
		s.port = port
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *RPCServer) {
		s.log = logger
	}
}

func WithJWT(jwt *auth.JWT) Option {
	return func(s *RPCServer) {
		s.jwt = jwt
	}
}

func WithHandler(method string, handler any) Option {
	return func(s *RPCServer) {
		s.handlers[method] = handler
	}
}
