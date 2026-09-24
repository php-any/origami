package bus

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

const dispatcherClassName = "Illuminate\\Bus\\Dispatcher"

type DispatcherClass struct {
	node.Node
	methods map[string]data.Method
}

func NewDispatcherClass() data.ClassStmt {
	c := &DispatcherClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *DispatcherClass) GetName() string    { return dispatcherClassName }
func (c *DispatcherClass) GetExtend() *string { return nil }
func (c *DispatcherClass) GetImplements() []string {
	return nil
}
func (c *DispatcherClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "container", "pipes", "handlers", "queueResolver", "allowsDispatchingAfterResponses", "pipeline":
		return node.NewProperty(nil, name, "protected", false, data.NewNullValue()), true
	}
	return nil, false
}
func (c *DispatcherClass) GetPropertyList() []data.Property {
	names := []string{"container", "pipes", "handlers", "queueResolver", "allowsDispatchingAfterResponses", "pipeline"}
	out := make([]data.Property, 0, len(names))
	for _, n := range names {
		p, _ := c.GetProperty(n)
		out = append(out, p)
	}
	return out
}
func (c *DispatcherClass) GetConstruct() data.Method { return c.methods["__construct"] }
func (c *DispatcherClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *DispatcherClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *DispatcherClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *DispatcherClass) GetStaticMethod(name string) (data.Method, bool) { return c.GetMethod(name) }

func (c *DispatcherClass) register() {
	c.methods["__construct"] = kit.InstanceMethodOpt("__construct", []string{"container", "queueResolver"}, 0, busConstruct)
	c.methods["dispatch"] = kit.InstanceMethod("dispatch", []string{"command"}, busDispatch)
	c.methods["dispatchsync"] = kit.InstanceMethodOpt("dispatchSync", []string{"command", "handler"}, 1, busDispatchSync)
	c.methods["dispatchnow"] = kit.InstanceMethodOpt("dispatchNow", []string{"command", "handler"}, 1, busDispatchNow)
}

func busConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Bus\\Dispatcher missing $this"))
	}
	if c := kit.Arg(ctx, 0); c != nil {
		_ = cv.SetProperty("container", c)
	}
	if q := kit.Arg(ctx, 1); q != nil {
		_ = cv.SetProperty("queueResolver", q)
	}
	_ = cv.SetProperty("pipes", data.NewArrayValue(nil))
	_ = cv.SetProperty("handlers", data.NewArrayValue(nil))
	_ = cv.SetProperty("allowsDispatchingAfterResponses", data.NewBoolValue(true))
	return cv, nil
}

func busDispatchCommand(ctx data.Context, command, handler data.Value) (data.GetValue, data.Control) {
	if handler != nil && !kit.IsNull(handler) {
		if cv, ok := kit.Unwrap(handler).(*data.ClassValue); ok && cv != nil {
			if m, ok := cv.GetMethod("handle"); ok && m != nil {
				return kit.CallInstanceMethod(ctx, cv, "handle", command)
			}
		}
		return kit.Call(ctx, handler, command)
	}
	cmd := kit.Unwrap(command)
	if fv, ok := cmd.(*data.FuncValue); ok && fv != nil {
		return kit.Call(ctx, fv)
	}
	if cv, ok := cmd.(*data.ClassValue); ok && cv != nil {
		if m, ok := cv.GetMethod("handle"); ok && m != nil {
			return kit.CallInstanceMethod(ctx, cv, "handle")
		}
		if m, ok := cv.GetMethod("__invoke"); ok && m != nil {
			return kit.CallInstanceMethod(ctx, cv, "__invoke")
		}
	}
	return command, nil
}

func busDispatch(ctx data.Context) (data.GetValue, data.Control) {
	return busDispatchCommand(ctx, kit.Arg(ctx, 0), nil)
}

func busDispatchSync(ctx data.Context) (data.GetValue, data.Control) {
	return busDispatchCommand(ctx, kit.Arg(ctx, 0), kit.Arg(ctx, 1))
}

func busDispatchNow(ctx data.Context) (data.GetValue, data.Control) {
	return busDispatchSync(ctx)
}
