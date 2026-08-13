package signal

import (
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type SignalChannelClass struct {
	channel *SignalChannel
}

func NewSignalChannelClass() data.ClassStmt {
	return &SignalChannelClass{
		channel: NewSignalChannel(),
	}
}

func (c *SignalChannelClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(&SignalChannelClass{
		channel: NewSignalChannel(),
	}, ctx.CreateBaseContext()), nil
}

func (c *SignalChannelClass) GetFrom() data.From { return nil }

func (c *SignalChannelClass) GetName() string { return "Signal\\Channel" }

func (c *SignalChannelClass) GetExtend() *string { return nil }

func (c *SignalChannelClass) GetImplements() []string { return nil }

func (c *SignalChannelClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *SignalChannelClass) GetPropertyList() []data.Property {
	return nil
}

func (c *SignalChannelClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "receive":
		return &SignalChannelReceiveMethod{source: c}, true
	case "close":
		return &SignalChannelCloseMethod{source: c}, true
	}
	return nil, false
}

func (c *SignalChannelClass) GetMethods() []data.Method {
	return []data.Method{
		&SignalChannelReceiveMethod{source: c},
		&SignalChannelCloseMethod{source: c},
	}
}

func (c *SignalChannelClass) GetConstruct() data.Method {
	return &SignalChannelConstructMethod{source: c}
}

func (c *SignalChannelClass) AddAnnotations(a *data.ClassValue) {}

type SignalChannelConstructMethod struct {
	source *SignalChannelClass
}

func (m *SignalChannelConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	capacity := 1
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if iv, ok := v.(*data.IntValue); ok {
			capacity, _ = iv.AsInt()
		}
	}
	m.source.channel.Construct(capacity)
	return data.NewNullValue(), nil
}

func (m *SignalChannelConstructMethod) GetName() string            { return "__construct" }
func (m *SignalChannelConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SignalChannelConstructMethod) GetIsStatic() bool          { return false }
func (m *SignalChannelConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "buffer", 0, data.NewIntValue(1), data.NewBaseType("int")),
	}
}
func (m *SignalChannelConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "buffer", 0, data.NewBaseType("int")),
	}
}
func (m *SignalChannelConstructMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

type SignalChannelReceiveMethod struct {
	source *SignalChannelClass
}

func (m *SignalChannelReceiveMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sig, ok := m.source.channel.Receive()
	if !ok {
		return data.NewNullValue(), nil
	}
	if s, ok := sig.(syscall.Signal); ok {
		return data.NewIntValue(int(s)), nil
	}
	return data.NewIntValue(0), nil
}

func (m *SignalChannelReceiveMethod) GetName() string            { return "receive" }
func (m *SignalChannelReceiveMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SignalChannelReceiveMethod) GetIsStatic() bool          { return false }
func (m *SignalChannelReceiveMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (m *SignalChannelReceiveMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *SignalChannelReceiveMethod) GetReturnType() data.Types { return data.NewBaseType("int") }

type SignalChannelCloseMethod struct {
	source *SignalChannelClass
}

func (m *SignalChannelCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	m.source.channel.Close()
	return nil, nil
}

func (m *SignalChannelCloseMethod) GetName() string            { return "close" }
func (m *SignalChannelCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SignalChannelCloseMethod) GetIsStatic() bool          { return false }
func (m *SignalChannelCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}
func (m *SignalChannelCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}
func (m *SignalChannelCloseMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
