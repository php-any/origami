package support

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const appendableClassName = "Illuminate\\View\\AppendableAttributeValue"

type AppendableClass struct {
	node.Node
	methods map[string]data.Method
}

func NewAppendableClass() data.ClassStmt {
	c := &AppendableClass{methods: map[string]data.Method{}}
	c.methods["__construct"] = newInstanceMethod("__construct", []string{"value"}, appendableConstruct, false)
	c.methods["__tostring"] = newInstanceMethod("__toString", nil, appendableToString, false)
	return c
}

func (c *AppendableClass) GetName() string { return appendableClassName }
func (c *AppendableClass) GetExtend() *string {
	return nil
}
func (c *AppendableClass) GetImplements() []string { return []string{"Stringable"} }
func (c *AppendableClass) GetProperty(name string) (data.Property, bool) {
	if name == "value" {
		return node.NewProperty(nil, "value", "public", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *AppendableClass) GetPropertyList() []data.Property {
	return []data.Property{node.NewProperty(nil, "value", "public", false, data.NewNullValue())}
}
func (c *AppendableClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *AppendableClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *AppendableClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[strings.ToLower(name)]
	return m, ok
}
func (c *AppendableClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *AppendableClass) GetStaticMethod(string) (data.Method, bool) { return nil, false }

func newAppendableValue(ctx data.Context, v data.Value) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewNullValue(), nil
	}
	cls, _ := vm.GetClass(appendableClassName)
	if cls == nil {
		cls = NewAppendableClass()
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	if v == nil {
		v = data.NewNullValue()
	}
	_ = cv.SetProperty("value", v)
	return cv, nil
}

func appendableConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := overlayReceiver(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	v := indexVal(ctx, 0)
	if v == nil {
		v = data.NewNullValue()
	}
	_ = cv.SetProperty("value", v)
	return cv, nil
}

func appendableToString(ctx data.Context) (data.GetValue, data.Control) {
	cv := overlayReceiver(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	v, ctl := cv.GetProperty("value")
	if ctl != nil {
		return nil, ctl
	}
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return data.NewStringValue(v.AsString()), nil
}
