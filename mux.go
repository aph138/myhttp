package myhttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) error

func newHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if e, ok := err.(*Error); ok {
				WriteJson(w, e.Status, e.Message)
			} else {
				log.Printf("http error: %s\n", err.Error())
			}
		}
	}
}
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("ContentType", "application/json")
	return json.NewEncoder(w).Encode(v)
}
func (e *Error) Error() string {
	return e.Message
}

type Error struct {
	Status  int
	Message string
}

func (s *Server) Handle(method string, path string, handler Handler) {
	s.mux.HandleFunc(fmt.Sprintf("%s %s", method, path), newHandler(handler))
}
func (s *Server) Get(path string, handler Handler) {
	s.Handle(http.MethodGet, path, handler)
}
func (s *Server) Post(path string, handler Handler) {
	s.Handle(http.MethodPost, path, handler)
}

func (s *Server) ServeFile(path string, file http.FileSystem) {
	s.mux.Handle(path, http.StripPrefix(path, http.FileServer(file)))
}
