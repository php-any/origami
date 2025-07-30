package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewServerClass() *ServerClass {
	source := newServer()
	return &ServerClass{
		init: &ServerConstructMethod{source},
		get:  &ServerGetMethod{source},
		run:  &ServerRunMethod{source},
	}
}

type ServerClass struct {
	node.Node
	init data.Method
	get  data.Method
	run  data.Method
}

func (s *ServerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := newServer()

	return data.NewClassValue(&ServerClass{
		init: &ServerConstructMethod{source},
		get:  &ServerGetMethod{source},
		run:  &ServerRunMethod{source},
	}, ctx.CreateBaseContext()), nil
}

func (s *ServerClass) GetName() string {
	return "Net\\Http\\Server"
}

func (s *ServerClass) GetExtend() *string {
	return nil
}

func (s *ServerClass) GetImplements() []string {
	return nil
}

func (s *ServerClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *ServerClass) GetProperties() map[string]data.Property {
	return nil
}

func (s *ServerClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "get":
		return s.get, true
	case "start":
		return s.run, true
	}
	return nil, false
}

func (s *ServerClass) GetMethods() []data.Method {
	return []data.Method{
		s.get,
	}
}

func (s *ServerClass) GetConstruct() data.Method {
	return s.init
}
