package node

import "github.com/php-any/origami/data"

type readonlyCloneState struct {
	object  *data.ObjectValue
	written map[string]bool
	active  bool
}
type readonlyCloneContext struct {
	data.Context
	data.CallRecorder
	state *readonlyCloneState
}

func (c *readonlyCloneContext) CreateContext(vars []data.Variable) data.Context {
	return newReadonlyCloneContext(c.Context.CreateContext(vars), c.state)
}
func (c *readonlyCloneContext) CreateBaseContext() data.Context {
	return newReadonlyCloneContext(c.Context.CreateBaseContext(), c.state)
}
func readonlyClonePermission(ctx data.Context, object *data.ObjectValue, name string) *readonlyCloneState {
	for ctx != nil {
		switch c := ctx.(type) {
		case *data.ClassMethodContext:
			ctx = c.Context
		case *data.BoundContext:
			ctx = c.Context
		case *readonlyCloneContext:
			if c.state.active && c.state.object == object && !c.state.written[name] {
				return c.state
			}
			return nil
		default:
			return nil
		}
	}
	return nil
}

func newReadonlyCloneContext(inner data.Context, state *readonlyCloneState) *readonlyCloneContext {
	return &readonlyCloneContext{Context: inner, CallRecorder: inner.(data.CallRecorder), state: state}
}
func (c *readonlyCloneContext) MarkEscaped() {
	if inner, ok := c.Context.(data.ContextEscaper); ok {
		inner.MarkEscaped()
	}
}
func (c *readonlyCloneContext) IsEscaped() bool {
	if inner, ok := c.Context.(data.ContextEscaper); ok {
		return inner.IsEscaped()
	}
	return false
}
func (c *readonlyCloneContext) ReleasePooled() {
	if inner, ok := c.Context.(interface{ ReleasePooled() }); ok {
		inner.ReleasePooled()
	}
}
