package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewServerClass() data.ClassStmt {
	return &ServerClass{
		source: httpsrc.NewServeMux(),
	}
}

func NewServerClassFromGroup(prefix string, server *ServerClass) data.ClassStmt {
	return &ServerClass{
		source:      server.source,
		Prefix:      prefix,
		Host:        server.Host,
		Port:        server.Port,
		Middlewares: append([]Middleware{}, server.Middlewares...),
	}
}

type ServerClass struct {
	node.Node
	source *httpsrc.ServeMux

	Prefix string

	Host        string
	Port        int
	Middlewares []Middleware
}

func (s *ServerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(NewServerClass(), ctx.CreateBaseContext()), nil
}

func (s *ServerClass) GetName() string         { return "Net\\Http\\Server" }
func (s *ServerClass) GetExtend() *string      { return nil }
func (s *ServerClass) GetImplements() []string { return nil }
func (s *ServerClass) AsString() string        { return "Server{}" }
func (s *ServerClass) GetSource() any          { return s.source }
func (s *ServerClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "get", "post", "put", "delete", "head", "options", "patch", "trace":
		return &ServerHandleMethod{server: s, name: name}, true
	case "group":
		return &ServerGroupMethod{server: s}, true
	case "middleware":
		return &ServerMiddlewareMethod{server: s}, true
	case "run":
		return &ServerRunMethod{server: s}, true
	}
	return nil, false
}

func (s *ServerClass) GetMethods() []data.Method {
	return []data.Method{
		&ServerHandleMethod{server: s, name: "get"},
		&ServerHandleMethod{server: s, name: "post"},
		&ServerHandleMethod{server: s, name: "put"},
		&ServerHandleMethod{server: s, name: "delete"},
		&ServerHandleMethod{server: s, name: "head"},
		&ServerHandleMethod{server: s, name: "options"},
		&ServerHandleMethod{server: s, name: "patch"},
		&ServerHandleMethod{server: s, name: "trace"},
		&ServerGroupMethod{server: s},
		&ServerMiddlewareMethod{server: s},
		&ServerRunMethod{server: s},
	}
}

func (s *ServerClass) GetConstruct() data.Method { return &ServerConstructMethod{source: s} }

func (s *ServerClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	}
	return nil, false
}

func (s *ServerClass) GetPropertyList() []data.Property {
	return []data.Property{}
}
