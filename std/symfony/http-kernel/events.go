package httpkernel

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)


const kernelEventsName = "Symfony\\Component\\HttpKernel\\KernelEvents"
const requestEventName = "Symfony\\Component\\HttpKernel\\Event\\RequestEvent"
const exceptionEventName = "Symfony\\Component\\HttpKernel\\Event\\ExceptionEvent"

type KernelEventsClass struct{ node.Node }

func NewKernelEventsClass() data.ClassStmt { return &KernelEventsClass{} }
func (c *KernelEventsClass) GetName() string {
	return kernelEventsName
}
func (c *KernelEventsClass) GetExtend() *string                       { return nil }
func (c *KernelEventsClass) GetImplements() []string                  { return nil }
func (c *KernelEventsClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *KernelEventsClass) GetPropertyList() []data.Property         { return nil }
func (c *KernelEventsClass) GetConstruct() data.Method                { return nil }
func (c *KernelEventsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *KernelEventsClass) GetMethod(string) (data.Method, bool) { return nil, false }
func (c *KernelEventsClass) GetMethods() []data.Method            { return nil }
func (c *KernelEventsClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "REQUEST":
		return data.NewStringValue("kernel.request"), true
	case "EXCEPTION":
		return data.NewStringValue("kernel.exception"), true
	case "RESPONSE":
		return data.NewStringValue("kernel.response"), true
	case "TERMINATE":
		return data.NewStringValue("kernel.terminate"), true
	case "CONTROLLER":
		return data.NewStringValue("kernel.controller"), true
	case "VIEW":
		return data.NewStringValue("kernel.view"), true
	case "FINISH_REQUEST":
		return data.NewStringValue("kernel.finish_request"), true
	}
	return nil, false
}

type RequestEventClass struct{ node.Node }

func NewRequestEventClass() data.ClassStmt { return &RequestEventClass{} }
func (c *RequestEventClass) GetName() string {
	return requestEventName
}
func (c *RequestEventClass) GetExtend() *string {
	parent := "Symfony\\Component\\HttpKernel\\Event\\KernelEvent"
	return &parent
}
func (c *RequestEventClass) GetImplements() []string                  { return nil }
func (c *RequestEventClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *RequestEventClass) GetPropertyList() []data.Property         { return nil }
func (c *RequestEventClass) GetConstruct() data.Method {
	return &hkMethod{name: "__construct", params: []string{"kernel", "request", "requestType"}, fn: func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewNullValue(), nil
	}}
}
func (c *RequestEventClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *RequestEventClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.GetConstruct(), true
	}
	return nil, false
}
func (c *RequestEventClass) GetMethods() []data.Method { return []data.Method{c.GetConstruct()} }

type ExceptionEventClass struct{ node.Node }

func NewExceptionEventClass() data.ClassStmt { return &ExceptionEventClass{} }
func (c *ExceptionEventClass) GetName() string {
	return exceptionEventName
}
func (c *ExceptionEventClass) GetExtend() *string {
	parent := "Symfony\\Component\\HttpKernel\\Event\\RequestEvent"
	return &parent
}
func (c *ExceptionEventClass) GetImplements() []string                  { return nil }
func (c *ExceptionEventClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *ExceptionEventClass) GetPropertyList() []data.Property         { return nil }
func (c *ExceptionEventClass) GetConstruct() data.Method {
	return &hkMethod{name: "__construct", params: []string{"kernel", "request", "requestType", "e"}, fn: func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewNullValue(), nil
	}}
}
func (c *ExceptionEventClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *ExceptionEventClass) GetMethod(name string) (data.Method, bool) {
	if name == "__construct" {
		return c.GetConstruct(), true
	}
	return nil, false
}
func (c *ExceptionEventClass) GetMethods() []data.Method { return []data.Method{c.GetConstruct()} }

type hkMethod struct {
	name   string
	params []string
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *hkMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *hkMethod) GetName() string                                     { return m.name }
func (m *hkMethod) GetModifier() data.Modifier                           { return data.ModifierPublic }
func (m *hkMethod) GetIsStatic() bool                                   { return false }
func (m *hkMethod) GetReturnType() data.Types                           { return nil }
func (m *hkMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p, i, nil, nil)
	}
	return out
}
func (m *hkMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p, i, nil)
	}
	return out
}
