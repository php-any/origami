package context

import (
	contextsrc "context"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewContextClassFrom(source contextsrc.Context) data.ClassStmt {
	return &ContextClass{
		source:   source,
		deadline: &ContextDeadlineMethod{source: source},
		done:     &ContextDoneMethod{source: source},
		err:      &ContextErrMethod{source: source},
		value:    &ContextValueMethod{source: source},
	}
}

type ContextClass struct {
	node.Node
	source   contextsrc.Context
	deadline data.Method
	done     data.Method
	err      data.Method
	value    data.Method
}

func (s *ContextClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}
func (s *ContextClass) GetName() string                            { return "Context\\Context" }
func (s *ContextClass) GetExtend() *string                         { return nil }
func (s *ContextClass) GetImplements() []string                    { return nil }
func (s *ContextClass) AsString() string                           { return "Context{}" }
func (s *ContextClass) GetSource() any                             { return s.source }
func (s *ContextClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (s *ContextClass) GetPropertyList() []data.Property           { return nil }
func (s *ContextClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "deadline":
		return s.deadline, true
	case "done":
		return s.done, true
	case "err":
		return s.err, true
	case "value":
		return s.value, true
	}
	return nil, false
}

func (s *ContextClass) GetMethods() []data.Method {
	return []data.Method{
		s.deadline,
		s.done,
		s.err,
		s.value,
	}
}

func (s *ContextClass) GetConstruct() data.Method { return nil }
