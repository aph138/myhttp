package myhttp

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config func(*http.Server)

type Server struct {
	srv         *http.Server
	mux         *http.ServeMux
	middlewares []Middleware
}

func defaultServer() *http.Server {
	return &http.Server{
		Addr:    ":9000",
		Handler: nil,
	}
}

func NewServer(c ...Config) *Server {
	srv := defaultServer()
	for _, i := range c {
		i(srv)
	}
	return &Server{
		srv: srv,
		mux: http.NewServeMux(),
	}
}
func WithAddress(add string) Config {
	return func(s *http.Server) {
		s.Addr = add
	}
}

func (s *Server) StartWithGracefulShutdown(t int, ctx context.Context, fn func() error) error {
	s.srv.Handler = stack(s.middlewares)(s.mux)
	e := make(chan error, 1)
	shutdown := make(chan int, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT)
		<-sig
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
	s.srv.Handler = stack(s.middlewares)(s.mux)
	return s.srv.ListenAndServe()
}
