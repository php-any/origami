package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewHandlerClass() data.ClassStmt {
	return &HandlerClass{
		source:    nil,
		serveHTTP: &HandlerServeHTTPMethod{source: nil},
	}
}

func NewHandlerClassFrom(source httpsrc.Handler) data.ClassStmt {
	return &HandlerClass{
		source:    source,
		serveHTTP: &HandlerServeHTTPMethod{source: source},
	}
}

type HandlerClass struct {
	node.Node
	source    httpsrc.Handler
	serveHTTP data.Method
}

func (s *HandlerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewHandlerClassFrom(nil), ctx.CreateBaseContext()), nil
}

func (s *HandlerClass) GetName() string         { return "Net\\Http\\Handler" }
func (s *HandlerClass) GetExtend() *string      { return nil }
func (s *HandlerClass) GetImplements() []string { return nil }
func (s *HandlerClass) AsString() string        { return "Handler{}" }
func (s *HandlerClass) GetSource() any          { return s.source }
func (s *HandlerClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "serveHTTP":
		return s.serveHTTP, true
	}
	return nil, false
}

func (s *HandlerClass) GetMethods() []data.Method {
	return []data.Method{
		s.serveHTTP,
	}
}

func (s *HandlerClass) GetConstruct() data.Method { return nil }

func (s *HandlerClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (s *HandlerClass) GetProperties() map[string]data.Property {
	return map[string]data.Property{}
}
