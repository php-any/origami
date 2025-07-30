package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"net/http"
)

func NewRequestClass(w http.ResponseWriter, r *http.Request) data.ClassStmt {
	source := &Request{
		r: r,
		w: w,
	}
	return &RequestClass{
		get:        &RequestGetMethod{source},
		input:      &RequestInputMethod{source},
		method:     &RequestMethodMethod{source},
		path:       &RequestPathMethod{source},
		requestUri: &RequestRequestURIMethod{source},
	}
}

type RequestClass struct {
	node.Node
	get        data.Method
	input      data.Method
	method     data.Method
	path       data.Method
	requestUri data.Method
}

func (s *RequestClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *RequestClass) GetName() string {
	return "Net\\Http\\Request"
}

func (s *RequestClass) GetExtend() *string {
	return nil
}

func (s *RequestClass) GetImplements() []string {
	return nil
}

func (s *RequestClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	}
	return nil, false
}

func (s *RequestClass) GetProperties() map[string]data.Property {
	return map[string]data.Property{}
}

func (s *RequestClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "get":
		return s.get, true
	case "input":
		return s.input, true
	case "method":
		return s.method, true
	case "path":
		return s.path, true
	case "requestUri":
		return s.requestUri, true
	}
	return nil, false
}

func (s *RequestClass) GetMethods() []data.Method {
	return []data.Method{
		s.get,
		s.input,
		s.method,
		s.path,
		s.requestUri,
	}
}

func (s *RequestClass) GetConstruct() data.Method {
	return nil
}
