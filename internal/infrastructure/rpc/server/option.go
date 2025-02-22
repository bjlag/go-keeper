package server

import (
	"go.uber.org/zap"
)

type Option func(*Server)

func WithAddress(host string, port int) Option {
	return func(s *Server) {
		s.host = host
		s.port = port
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.log = logger
	}
}

func WithHandler(method string, handler any) Option {
	return func(s *Server) {
		s.handlers[method] = handler
	}
}
