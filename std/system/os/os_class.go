package os

import (
	"runtime"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewOSClass() data.ClassStmt {
	source := newOs()
	return &OSClass{
		eol:      data.NewStringValue(source.EOL),
		exit:     &OSExitMethod{source},
		hostname: &OSHostnameMethod{source},
		path:     &OSPathMethod{source},
	}
}

type OSClass struct {
	node.Node
	eol      data.Value
	exit     data.Method
	hostname data.Method
	path     data.Method
}

func (s *OSClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(s, ctx), nil
}

func (s *OSClass) GetName() string {
	return "OS"
}

func (s *OSClass) GetExtend() *string {
	return nil
}

func (s *OSClass) GetImplements() []string {
	return nil
}

func (s *OSClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (s *OSClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "EOL":
		return s.eol, true
	case "GOOS":
		return data.NewStringValue(runtime.GOOS), true

	}
	return nil, false
}

func (s *OSClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (s *OSClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "exit":
		return s.exit, true
	case "hostname":
		return s.hostname, true
	case "path":
		return s.path, true
	}
	return nil, false
}

func (s *OSClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "exit":
		return s.exit, true
	case "hostname":
		return s.hostname, true
	case "path":
		return s.path, true
	}
	return nil, false
}

func (s *OSClass) GetMethods() []data.Method {
	return []data.Method{
		s.exit,
		s.hostname,
		s.path,
	}
}

func (s *OSClass) GetConstruct() data.Method {
	return nil
}
