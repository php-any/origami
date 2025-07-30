package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"net/http"
)

func NewResponseClass(w http.ResponseWriter, r *http.Request) data.ClassStmt {
	source := &Response{
		r: r,
		w: w,
	}
	return &ResponseClass{
		write: &ResponseWriteMethod{source},
	}
}

type ResponseClass struct {
	node.Node
	write data.Method
}

func (s *ResponseClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *ResponseClass) GetName() string {
	return "Net\\Http\\Response"
}

func (s *ResponseClass) GetExtend() *string {
	return nil
}

func (s *ResponseClass) GetImplements() []string {
	return nil
}

func (s *ResponseClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	}
	return nil, false
}

func (s *ResponseClass) GetProperties() map[string]data.Property {
	return map[string]data.Property{}
}

func (s *ResponseClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "write":
		return s.write, true
	}
	return nil, false
}

func (s *ResponseClass) GetMethods() []data.Method {
	return []data.Method{
		s.write,
	}
}

func (s *ResponseClass) GetConstruct() data.Method {
	return nil
}
