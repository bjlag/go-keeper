package server

import (
	"github.com/bjlag/go-keeper/internal/infrastructure/auth/jwt"
	"go.uber.org/zap"
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

func WithJWTGenerator(jwt *jwt.Generator) Option {
	return func(s *RPCServer) {
		s.jwt = jwt
	}
}

func WithHandler(method string, handler any) Option {
	return func(s *RPCServer) {
		s.handlers[method] = handler
	}
}
