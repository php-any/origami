package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewHeaderClass() data.ClassStmt {
	return &HeaderClass{
		source:      nil,
		get:         &HeaderGetMethod{source: nil},
		set:         &HeaderSetMethod{source: nil},
		values:      &HeaderValuesMethod{source: nil},
		write:       &HeaderWriteMethod{source: nil},
		writeSubset: &HeaderWriteSubsetMethod{source: nil},
		add:         &HeaderAddMethod{source: nil},
		clone:       &HeaderCloneMethod{source: nil},
		del:         &HeaderDelMethod{source: nil},
	}
}

func NewHeaderClassFrom(source *httpsrc.Header) data.ClassStmt {
	return &HeaderClass{
		source:      source,
		write:       &HeaderWriteMethod{source: source},
		writeSubset: &HeaderWriteSubsetMethod{source: source},
		add:         &HeaderAddMethod{source: source},
		clone:       &HeaderCloneMethod{source: source},
		del:         &HeaderDelMethod{source: source},
		get:         &HeaderGetMethod{source: source},
		set:         &HeaderSetMethod{source: source},
		values:      &HeaderValuesMethod{source: source},
	}
}

type HeaderClass struct {
	node.Node
	source      *httpsrc.Header
	values      data.Method
	write       data.Method
	writeSubset data.Method
	add         data.Method
	clone       data.Method
	del         data.Method
	get         data.Method
	set         data.Method
}

func (s *HeaderClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewProxyValue(s, ctx.CreateBaseContext()), nil
}

func (s *HeaderClass) GetName() string         { return "Net\\Http\\Header" }
func (s *HeaderClass) GetExtend() *string      { return nil }
func (s *HeaderClass) GetImplements() []string { return nil }
func (s *HeaderClass) AsString() string        { return "Header{}" }
func (s *HeaderClass) GetSource() any          { return s.source }
func (s *HeaderClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "add":
		return s.add, true
	case "clone":
		return s.clone, true
	case "del":
		return s.del, true
	case "get":
		return s.get, true
	case "set":
		return s.set, true
	case "values":
		return s.values, true
	case "write":
		return s.write, true
	case "writeSubset":
		return s.writeSubset, true
	}
	return nil, false
}

func (s *HeaderClass) GetMethods() []data.Method {
	return []data.Method{
		s.values,
		s.write,
		s.writeSubset,
		s.add,
		s.clone,
		s.del,
		s.get,
		s.set,
	}
}

func (s *HeaderClass) GetConstruct() data.Method { return nil }

func (s *HeaderClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (s *HeaderClass) GetPropertyList() []data.Property {
	return []data.Property{}
}
