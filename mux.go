package myhttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Mux struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewMux() *Mux {
	return &Mux{
		mux: http.NewServeMux(),
	}
}

type Handler func(c *Context) error

func newHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := Context{
			Writer:  w,
			Request: r,
		}
		if err := h(&c); err != nil {
			if e, ok := err.(*Error); ok {
				writeJson(w, e.Status, e.Message)
			} else {
				//TODO
				log.Printf("http error: %s\n", err.Error())
			}
		}
	}
}
func writeJson(w http.ResponseWriter, status int, v any) error {
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

func (m *Mux) Handle(method string, path string, handler Handler) {
	m.mux.HandleFunc(fmt.Sprintf("%s %s", method, path), newHandler(handler))
}
func (m *Mux) Get(path string, handler Handler) {
	m.Handle(http.MethodGet, path, handler)
}
func (m *Mux) Post(path string, handler Handler) {
	m.Handle(http.MethodPost, path, handler)
}

func (m *Mux) ServeFolder(path string, file http.FileSystem) {
	m.mux.Handle(path, http.StripPrefix(path, http.FileServer(file)))
}
func (m *Mux) ServeFile(path string, file string) {
	m.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	})
}

// you can't add middleware after calling AddSubRouter
func (m *Mux) Use(middle Middleware) {
	m.middlewares = append(m.middlewares, middle)
}
