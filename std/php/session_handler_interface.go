package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SessionHandlerInterface 表示 PHP 的 SessionHandlerInterface 接口
type SessionHandlerInterface struct {
	*node.InterfaceStatement
}

func NewSessionHandlerInterface() *SessionHandlerInterface {
	return &SessionHandlerInterface{
		InterfaceStatement: node.NewInterfaceStatement(
			nil,
			"SessionHandlerInterface",
			nil,
			[]data.Method{
				&SessionHandlerOpenMethod{},
				&SessionHandlerCloseMethod{},
				&SessionHandlerReadMethod{},
				&SessionHandlerWriteMethod{},
				&SessionHandlerDestroyMethod{},
				&SessionHandlerGcMethod{},
			},
		),
	}
}

type SessionHandlerOpenMethod struct{}

func (m *SessionHandlerOpenMethod) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}
func (m *SessionHandlerOpenMethod) GetName() string            { return "open" }
func (m *SessionHandlerOpenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SessionHandlerOpenMethod) GetIsStatic() bool          { return false }
func (m *SessionHandlerOpenMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "save_path", 0, nil, nil),
		node.NewParameter(nil, "session_name", 1, nil, nil),
	}
}
func (m *SessionHandlerOpenMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "save_path", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "session_name", 1, data.NewBaseType("string")),
	}
}
func (m *SessionHandlerOpenMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }

type SessionHandlerCloseMethod struct{}

func (m *SessionHandlerCloseMethod) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}
func (m *SessionHandlerCloseMethod) GetName() string               { return "close" }
func (m *SessionHandlerCloseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SessionHandlerCloseMethod) GetIsStatic() bool             { return false }
func (m *SessionHandlerCloseMethod) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *SessionHandlerCloseMethod) GetVariables() []data.Variable { return []data.Variable{} }
func (m *SessionHandlerCloseMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }

type SessionHandlerReadMethod struct{}

func (m *SessionHandlerReadMethod) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(""), nil
}
func (m *SessionHandlerReadMethod) GetName() string            { return "read" }
func (m *SessionHandlerReadMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SessionHandlerReadMethod) GetIsStatic() bool          { return false }
func (m *SessionHandlerReadMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "id", 0, nil, nil)}
}
func (m *SessionHandlerReadMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "id", 0, data.NewBaseType("string"))}
}
func (m *SessionHandlerReadMethod) GetReturnType() data.Types { return data.NewBaseType("string") }

type SessionHandlerWriteMethod struct{}

func (m *SessionHandlerWriteMethod) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}
func (m *SessionHandlerWriteMethod) GetName() string            { return "write" }
func (m *SessionHandlerWriteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SessionHandlerWriteMethod) GetIsStatic() bool          { return false }
func (m *SessionHandlerWriteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "id", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, nil),
	}
}
func (m *SessionHandlerWriteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "id", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "data", 1, data.NewBaseType("string")),
	}
}
func (m *SessionHandlerWriteMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }

type SessionHandlerDestroyMethod struct{}

func (m *SessionHandlerDestroyMethod) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(true), nil
}
func (m *SessionHandlerDestroyMethod) GetName() string            { return "destroy" }
func (m *SessionHandlerDestroyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SessionHandlerDestroyMethod) GetIsStatic() bool          { return false }
func (m *SessionHandlerDestroyMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "id", 0, nil, nil)}
}
func (m *SessionHandlerDestroyMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "id", 0, data.NewBaseType("string"))}
}
func (m *SessionHandlerDestroyMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }

type SessionHandlerGcMethod struct{}

func (m *SessionHandlerGcMethod) Call(data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(0), nil
}
func (m *SessionHandlerGcMethod) GetName() string            { return "gc" }
func (m *SessionHandlerGcMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SessionHandlerGcMethod) GetIsStatic() bool          { return false }
func (m *SessionHandlerGcMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "max_lifetime", 0, nil, nil)}
}
func (m *SessionHandlerGcMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "max_lifetime", 0, data.NewBaseType("int"))}
}
func (m *SessionHandlerGcMethod) GetReturnType() data.Types { return data.NewBaseType("int") }
