package myhttp

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config func(*Server)

type Server struct {
	srv         *http.Server
	mux         *http.ServeMux
	middlewares []Middleware
	logger      *slog.Logger
	vebose      bool
}

func defaultServer() *Server {

	return &Server{
		srv: &http.Server{
			Addr:    ":9000",
			Handler: nil,
		},
		mux:    &http.ServeMux{},
		logger: slog.Default(),
		vebose: true,
	}
}

func Quite() Config {
	return func(s *Server) {
		s.vebose = false
	}
}

func (s *Server) info(msg string) {
	if s.vebose {
		s.logger.Info(msg)
	}
}
func NewServer(c ...Config) *Server {
	srv := defaultServer()
	for _, i := range c {
		i(srv)
	}
	return srv
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
func (s *Server) StartWithGracefulShutdown(t int, ctx context.Context, fn func() error) error {
	s.info("start server with graceful shutdown fucntion on " + s.srv.Addr)
	s.srv.Handler = stack(s.middlewares)(s.mux)
	e := make(chan error, 1)
	shutdown := make(chan int, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT)
		<-sig
		s.info("starting shutting down...")
		err := fn()
		if err != nil {
			e <- err
		}
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(t))
		defer cancel()
		if err = s.srv.Shutdown(ctx); err != nil {
			e <- err
		}
		shutdown <- 1
	}()
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e <- err
		}
	}()
	for {
		select {
		case <-shutdown:
			{
				return nil
			}
		case err := <-e:
			{
				return err
			}
		}
	}

}
func (s *Server) Start() error {
	s.info("starting the server on " + s.srv.Addr)
	s.srv.Handler = stack(s.middlewares)(s.mux)
	return s.srv.ListenAndServe()
}
