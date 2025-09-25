package http

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	httpsrc "net/http"
)

func NewServeMuxClass() data.ClassStmt {
	return &ServeMuxClass{
		source:     nil,
		handle:     &ServeMuxHandleMethod{source: nil},
		handleFunc: &ServeMuxHandleFuncMethod{source: nil},
		handler:    &ServeMuxHandlerMethod{source: nil},
		serveHTTP:  &ServeMuxServeHTTPMethod{source: nil},
	}
}

func NewServeMuxClassFrom(source *httpsrc.ServeMux) data.ClassStmt {
	return &ServeMuxClass{
		source:     source,
		handle:     &ServeMuxHandleMethod{source: source},
		handleFunc: &ServeMuxHandleFuncMethod{source: source},
		handler:    &ServeMuxHandlerMethod{source: source},
		serveHTTP:  &ServeMuxServeHTTPMethod{source: source},
	}
}

type ServeMuxClass struct {
	node.Node
	source     *httpsrc.ServeMux
	handle     data.Method
	handleFunc data.Method
	handler    data.Method
	serveHTTP  data.Method
}

func (s *ServeMuxClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewServeMuxClassFrom(&httpsrc.ServeMux{}), ctx.CreateBaseContext()), nil
}

func (s *ServeMuxClass) GetName() string         { return "http\\ServeMux" }
func (s *ServeMuxClass) GetExtend() *string      { return nil }
func (s *ServeMuxClass) GetImplements() []string { return nil }
func (s *ServeMuxClass) AsString() string        { return "ServeMux{}" }
func (s *ServeMuxClass) GetSource() any          { return s.source }
func (s *ServeMuxClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "handler":
		return s.handler, true
	case "serveHTTP":
		return s.serveHTTP, true
	case "handle":
		return s.handle, true
	case "handleFunc":
		return s.handleFunc, true
	}
	return nil, false
}

func (s *ServeMuxClass) GetMethods() []data.Method {
	return []data.Method{
		s.handle,
		s.handleFunc,
		s.handler,
		s.serveHTTP,
	}
}

func (s *ServeMuxClass) GetConstruct() data.Method { return nil }

func (s *ServeMuxClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	}
	return nil, false
}

func (s *ServeMuxClass) GetProperties() map[string]data.Property {
	return map[string]data.Property{}
}

func (s *ServeMuxClass) SetProperty(name string, value data.Value) data.Control {
	if s.source == nil {
		return data.NewErrorThrow(nil, errors.New("无法设置属性，source 为 nil"))
	}

	switch name {
	default:
		return data.NewErrorThrow(nil, errors.New("属性不存在: "+name))
	}
}
