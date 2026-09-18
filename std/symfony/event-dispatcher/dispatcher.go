package eventdispatcher

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)


const eventDispatcherName = "Symfony\\Component\\EventDispatcher\\EventDispatcher"

type EventDispatcherClass struct{ node.Node }

func NewEventDispatcherClass() data.ClassStmt { return &EventDispatcherClass{} }
func (c *EventDispatcherClass) GetName() string {
	return eventDispatcherName
}
func (c *EventDispatcherClass) GetExtend() *string { return nil }
func (c *EventDispatcherClass) GetImplements() []string {
	return []string{"Symfony\\Contracts\\EventDispatcher\\EventDispatcherInterface"}
}
func (c *EventDispatcherClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *EventDispatcherClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "listeners", "private", false, data.NewArrayValue(nil))}
}
func (c *EventDispatcherClass) GetConstruct() data.Method {
	return &edMethod{name: "__construct", fn: func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewNullValue(), nil
	}}
}
func (c *EventDispatcherClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *EventDispatcherClass) GetMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "__construct":
		return c.GetConstruct(), true
	case "dispatch":
		return &edMethod{name: "dispatch", params: []string{"event", "eventName"}, fn: edDispatch}, true
	case "addlistener":
		return &edMethod{name: "addListener", params: []string{"eventName", "listener", "priority"}, fn: edAddListener}, true
	}
	return nil, false
}
func (c *EventDispatcherClass) GetMethods() []data.Method { return []data.Method{c.GetConstruct()} }

type edMethod struct {
	name   string
	params []string
	fn     func(data.Context) (data.GetValue, data.Control)
}

func (m *edMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.fn(ctx) }
func (m *edMethod) GetName() string                                     { return m.name }
func (m *edMethod) GetModifier() data.Modifier                           { return data.ModifierPublic }
func (m *edMethod) GetIsStatic() bool                                   { return false }
func (m *edMethod) GetReturnType() data.Types                           { return nil }
func (m *edMethod) GetParams() []data.GetValue {
	out := make([]data.GetValue, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewParameter(nil, p, i, nil, nil)
	}
	return out
}
func (m *edMethod) GetVariables() []data.Variable {
	out := make([]data.Variable, len(m.params))
	for i, p := range m.params {
		out[i] = node.NewVariable(nil, p, i, nil)
	}
	return out
}

func edSelf(ctx data.Context) *data.ClassValue {
	if c, ok := ctx.(*data.ClassMethodContext); ok {
		return c.ClassValue
	}
	return nil
}

func edAddListener(ctx data.Context) (data.GetValue, data.Control) {
	cv := edSelf(ctx)
	name, _ := ctx.GetIndexValue(0)
	listener, _ := ctx.GetIndexValue(1)
	listeners, _ := cv.GetProperty("listeners")
	arr, _ := listeners.(*data.ArrayValue)
	if arr == nil {
		arr = data.NewArrayValue(nil).(*data.ArrayValue)
	}
	key := ""
	if name != nil {
		key = name.AsString()
	}
	arr.List = append(arr.List, data.NewNamedZVal(key, listener))
	_ = cv.SetProperty("listeners", arr)
	return data.NewNullValue(), nil
}

func edDispatch(ctx data.Context) (data.GetValue, data.Control) {
	event, _ := ctx.GetIndexValue(0)
	if event == nil {
		return data.NewNullValue(), nil
	}
	return event, nil
}
