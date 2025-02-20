package server

import (
	"go.uber.org/zap"
)

type Option func(*Server)

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
