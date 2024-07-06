package myhttp

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (c *Context) Json(status int, data any) error {
	c.Writer.WriteHeader(status)
	c.Writer.Header().Add("ContentType", "application/json")
	return json.NewEncoder(c.Writer).Encode(status)
}

func (c *Context) SaveFile() error {
	return nil
}
