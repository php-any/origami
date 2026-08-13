package container

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewCircularDependencyExceptionClass() data.ClassStmt {
	source := &circularDependencyExceptionSource{}
	return &CircularDependencyExceptionClass{
		getMessage: &circularDependencyGetMessageMethod{source: source},
	}
}

type circularDependencyExceptionSource struct {
	msg string
}

type CircularDependencyExceptionClass struct {
	node.Node
	getMessage data.Method
}

func (c *CircularDependencyExceptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *CircularDependencyExceptionClass) GetName() string {
	return "Container\\CircularDependencyException"
}
func (c *CircularDependencyExceptionClass) GetExtend() *string {
	extend := "Exception"
	return &extend
}
func (c *CircularDependencyExceptionClass) GetImplements() []string { return nil }
func (c *CircularDependencyExceptionClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}
func (c *CircularDependencyExceptionClass) GetPropertyList() []data.Property { return nil }
func (c *CircularDependencyExceptionClass) GetMethod(name string) (data.Method, bool) {
	if name == "getMessage" {
		return c.getMessage, true
	}
	return nil, false
}
func (c *CircularDependencyExceptionClass) GetMethods() []data.Method {
	return []data.Method{c.getMessage}
}
func (c *CircularDependencyExceptionClass) GetConstruct() data.Method { return nil }

type circularDependencyGetMessageMethod struct {
	source *circularDependencyExceptionSource
}

func (m *circularDependencyGetMessageMethod) GetName() string            { return "getMessage" }
func (m *circularDependencyGetMessageMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *circularDependencyGetMessageMethod) GetIsStatic() bool          { return false }
func (m *circularDependencyGetMessageMethod) GetParams() []data.GetValue { return nil }
func (m *circularDependencyGetMessageMethod) GetVariables() []data.Variable {
	return nil
}
func (m *circularDependencyGetMessageMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
func (m *circularDependencyGetMessageMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(m.source.msg), nil
}

func circularDependencyError(abstract string) data.Control {
	msg := "Circular dependency detected while resolving [" + abstract + "]"
	return data.NewErrorThrowByName(nil, &circularDependencyErr{msg: msg}, "Container\\CircularDependencyException")
}

type circularDependencyErr struct{ msg string }

func (e *circularDependencyErr) Error() string { return e.msg }
