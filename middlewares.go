package myhttp

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func stack(m []Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			h = m[i](h)
		}
		return h
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		wWrapper := responseWriter{ResponseWriter: w, StatusCode: 200}
		next.ServeHTTP(&wWrapper, r)
		e := time.Since(s)
		log.Printf("%d %s %s %s\n", wWrapper.StatusCode, r.Method, r.URL.Path, e.String())
	})
}

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (r *responseWriter) WriteHeader(s int) {
	r.ResponseWriter.WriteHeader(s)
	r.StatusCode = s
}
