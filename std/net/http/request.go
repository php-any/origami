package http

import (
	"github.com/php-any/origami/data"
	"net/http"
)

type Request struct {
	w http.ResponseWriter
	r *http.Request
}

func (req *Request) Method() string {
	return http.MethodGet //TODO
}

func (req *Request) Path() string {
	return "TODO"
}

func (req *Request) RequestURI() string {
	return "TODO"
}
func (req *Request) Input(name string) data.Value {
	return nil
}

func (req *Request) Get(name string) data.Value {
	return nil
}
