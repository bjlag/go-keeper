package server

import (
	"go.uber.org/zap"
	"net"

	"github.com/bjlag/go-keeper/internal/infrastructure/auth"
)

// Option тип параметра сервера.
type Option func(*RPCServer)

// WithAddress передача адреса сервера.
func WithAddress(host string, port int) Option {
	return func(s *RPCServer) {
		s.host = host
		s.port = port
	}
}

// WithListener передача сетевого прослушивателя сервера.
func WithListener(listener net.Listener) Option {
	return func(s *RPCServer) {
		s.listener = listener
	}
}

// WithLogger передача логгера.
func WithLogger(logger *zap.Logger) Option {
	return func(s *RPCServer) {
		s.log = logger
	}
}

// WithJWT передача сервиса для работы с JWT токенами.
func WithJWT(jwt *auth.JWT) Option {
	return func(s *RPCServer) {
		s.jwt = jwt
	}
}

// WithHandler регистрация обработчика запроса по указанному имени.
func WithHandler(method string, handler any) Option {
	return func(s *RPCServer) {
		s.handlers[method] = handler
	}
}
