package myhttp

import (
	"crypto/tls"
	"log/slog"
	"time"
)

type Config func(*Server)

func Quite() Config {
	return func(s *Server) {
		s.vebose = false
	}
}
func WithAddress(add string) Config {
	return func(s *Server) {
		s.srv.Addr = add
	}
}
func WithCustomLogger(logger *slog.Logger) Config {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithIdleTimeOut(time time.Duration) Config {
	return func(s *Server) {
		s.srv.IdleTimeout = time
	}
}

func WithReadTimeout(time time.Duration) Config {
	return func(s *Server) {
		s.srv.ReadTimeout = time
	}
}
func WithWriteTimeout(time time.Duration) Config {
	return func(s *Server) {
		s.srv.WriteTimeout = time
	}
}
func WithTLS(tlsConfig *tls.Config) Config {
	return func(s *Server) {
		s.srv.TLSConfig = tlsConfig
	}
}
