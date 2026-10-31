// Package mongodb implements the supported subset of PHP ext-mongodb using
// MongoDB's official Go driver. Unsupported options are errors, never ignored.
package mongodb

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/exception"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const driverNS = "MongoDB\\Driver\\"
const bsonNS = "MongoDB\\BSON\\"
const exceptionNS = driverNS + "Exception\\"

type nativeClass struct {
	node.Node
	name       string
	interfaces []string
	methods    []*nativeMethod
	constants  map[string]data.Value
}

func (c *nativeClass) GetName() string                          { return c.name }
func (c *nativeClass) GetExtend() *string                       { return nil }
func (c *nativeClass) GetImplements() []string                  { return c.interfaces }
func (c *nativeClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *nativeClass) GetPropertyList() []data.Property         { return nil }
func (c *nativeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *nativeClass) GetMethod(name string) (data.Method, bool) {
	for _, m := range c.methods {
		if strings.EqualFold(name, m.name) {
			return m, true
		}
	}
	return nil, false
}
func (c *nativeClass) GetMethods() []data.Method {
	r := make([]data.Method, len(c.methods))
	for i, m := range c.methods {
		r[i] = m
	}
	return r
}
func (c *nativeClass) GetConstruct() data.Method { m, _ := c.GetMethod("__construct"); return m }
func (c *nativeClass) GetStaticProperty(name string) (data.Value, bool) {
	v, ok := c.constants[name]
	return v, ok
}

type nativeMethod struct {
	name      string
	params    []data.GetValue
	variables []data.Variable
	call      func(data.Context) (data.GetValue, data.Control)
}

// A trailing ? means an optional parameter with a null default.
func method(name string, names []string, call func(data.Context) (data.GetValue, data.Control)) *nativeMethod {
	m := &nativeMethod{name: name, call: call}
	for i, n := range names {
		var def data.GetValue
		if strings.HasSuffix(n, "?") {
			n = strings.TrimSuffix(n, "?")
			def = data.NewNullValue()
		}
		m.params = append(m.params, node.NewParameter(nil, n, i, def, nil))
		m.variables = append(m.variables, node.NewVariable(nil, n, i, nil))
	}
	return m
}
func (m *nativeMethod) GetName() string                                     { return m.name }
func (m *nativeMethod) GetModifier() data.Modifier                          { return data.ModifierPublic }
func (m *nativeMethod) GetIsStatic() bool                                   { return false }
func (m *nativeMethod) GetReturnType() data.Types                           { return nil }
func (m *nativeMethod) GetParams() []data.GetValue                          { return m.params }
func (m *nativeMethod) GetVariables() []data.Variable                       { return m.variables }
func (m *nativeMethod) Call(ctx data.Context) (data.GetValue, data.Control) { return m.call(ctx) }
func arg(ctx data.Context, i int) data.Value                                { v, _ := ctx.GetIndexValue(i); return v }
func isNull(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}
func receiver(ctx data.Context) *data.PropertyBag {
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		return c.ObjectValue
	case *data.ClassValue:
		return c.ObjectValue
	}
	return nil
}
func object(ctx data.Context, name string, state any) *data.ClassValue {
	c, ok := ctx.GetVM().GetClass(name)
	if !ok {
		panic("unregistered MongoDB class: " + name)
	}
	v := data.NewClassValue(c, ctx)
	v.InstanceSource = state
	return v
}
func state[T any](ctx data.Context) (T, bool) {
	var zero T
	r := receiver(ctx)
	if r == nil {
		return zero, false
	}
	s, ok := r.InstanceSource.(T)
	return s, ok
}
func objectState[T any](v data.Value) (T, bool) {
	var zero T
	c, ok := v.(*data.ClassValue)
	if !ok {
		return zero, false
	}
	s, ok := c.InstanceSource.(T)
	return s, ok
}
func invalid(err error) (data.GetValue, data.Control) {
	return nil, data.NewErrorThrowByName(nil, err, exceptionNS+"InvalidArgumentException")
}
func invalidf(format string, a ...any) (data.GetValue, data.Control) {
	return invalid(fmt.Errorf(format, a...))
}
func logic(message string) (data.GetValue, data.Control) {
	return nil, data.NewErrorThrowByName(nil, fmt.Errorf("%s", message), exceptionNS+"LogicException")
}
func stringArg(v data.Value, name string) (string, error) {
	s, ok := v.(*data.StringValue)
	if !ok {
		return "", fmt.Errorf("%s must be a string", name)
	}
	return s.Value, nil
}
func intArg(v data.Value, name string) (int64, error) {
	s, ok := v.(*data.IntValue)
	if !ok {
		return 0, fmt.Errorf("%s must be an integer", name)
	}
	return int64(s.Value), nil
}
func optionMap(v data.Value, allowed ...string) (map[string]data.Value, error) {
	r := map[string]data.Value{}
	if isNull(v) {
		return r, nil
	}
	a, ok := v.(*data.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("options must be an array")
	}
	for i, s := range a.Range() {
		k := s.PHPArrayKey(i).AsString()
		found := false
		for _, n := range allowed {
			if k == n {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("unsupported MongoDB option %q", k)
		}
		r[k] = s.ReadValue()
	}
	return r, nil
}
func boolOption(m map[string]data.Value, key string, def bool) (bool, error) {
	v, ok := m[key]
	if !ok {
		return def, nil
	}
	b, ok := v.(*data.BoolValue)
	if !ok {
		return false, fmt.Errorf("%s must be a boolean", key)
	}
	return b.Value, nil
}

// Embed the ordinary exception implementation so properties and methods follow
// the runtime's existing Throwable contract, including inheritance.
type driverException struct {
	*exception.ExceptionClass
	name, parent string
}

func (c *driverException) GetName() string         { return exceptionNS + c.name }
func (c *driverException) GetExtend() *string      { return &c.parent }
func (c *driverException) GetImplements() []string { return []string{exceptionNS + "Exception"} }
func (c *driverException) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *driverException) GetMethod(name string) (data.Method, bool) {
	if c.name == "BulkWriteException" && strings.EqualFold(name, "getWriteResult") {
		return method("getWriteResult", nil, func(ctx data.Context) (data.GetValue, data.Control) {
			r, ok := state[*mongo.BulkWriteResult](ctx)
			if !ok || r == nil {
				return data.NewNullValue(), nil
			}
			return object(ctx, driverNS+"WriteResult", r), nil
		}), true
	}
	return c.ExceptionClass.GetMethod(name)
}
func (c *driverException) GetMethods() []data.Method {
	m := c.ExceptionClass.GetMethods()
	if c.name == "BulkWriteException" {
		r, _ := c.GetMethod("getWriteResult")
		m = append(m, r)
	}
	return m
}

func Load(vm data.VM) {
	vm.SetConstant("MONGODB_VERSION", data.NewStringValue("0.1.0-origami"))
	vm.AddInterface(node.NewInterfaceStatement(nil, bsonNS+"Type", nil, nil))
	vm.AddInterface(node.NewInterfaceStatement(nil, exceptionNS+"Exception", []string{"Throwable"}, nil))
	for _, e := range [][2]string{{"RuntimeException", "RuntimeException"}, {"LogicException", "LogicException"}, {"InvalidArgumentException", exceptionNS + "LogicException"}, {"UnexpectedValueException", exceptionNS + "RuntimeException"}, {"ConnectionException", exceptionNS + "RuntimeException"}, {"ConnectionTimeoutException", exceptionNS + "ConnectionException"}, {"AuthenticationException", exceptionNS + "ConnectionException"}, {"BulkWriteException", exceptionNS + "RuntimeException"}, {"CommandException", exceptionNS + "RuntimeException"}} {
		vm.AddClass(&driverException{exception.NewExceptionClass(), e[0], e[1]})
	}
	for _, c := range bsonClasses() {
		vm.AddClass(c)
	}
	for _, c := range driverClasses() {
		vm.AddClass(c)
	}
	for _, f := range bsonFunctions() {
		vm.AddFunc(f)
	}
}
