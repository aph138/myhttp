package myhttp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	srv         *http.Server
	Mux         *Mux
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
		Mux:    NewMux(),
		logger: slog.Default(),
		vebose: true,
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

func (s *Server) AddSubRouter(path string, m *Mux) {
	//TODO:check path
	if path != "/" {
		s.Mux.mux.Handle(fmt.Sprintf("%s/", path), http.StripPrefix(path, stack(m.middlewares)(m.mux)))
	}
}
func (s *Server) Handle(method string, path string, handler Handler) {
	s.Mux.Handle(method, path, handler)
}
func (s *Server) Get(path string, handler Handler) {
	s.Mux.Get(path, handler)
}
func (s *Server) Post(path string, handler Handler) {
	s.Mux.Post(path, handler)
}
func (s *Server) ServeFolder(path string, file http.FileSystem) {
	s.Mux.ServeFolder(path, file)
}
func (s *Server) ServeFile(path string, file string) {
	s.Mux.ServeFile(path, file)
}

// you can't add add middleware after starting the server
func (s *Server) Use(m Middleware) {
	s.middlewares = append(s.middlewares, m)
	// s.srv.Handler = m(s.srv.Handler)
}

func (s *Server) StartWithGracefulShutdown(t int, ctx context.Context, fn func() error) error {
	s.info("start server with graceful shutdown fucntion on " + s.srv.Addr)
	s.srv.Handler = stack(s.middlewares)(s.Mux.mux)
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
	s.srv.Handler = stack(s.middlewares)(s.Mux.mux)
	return s.srv.ListenAndServe()
}
