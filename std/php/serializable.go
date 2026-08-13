package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SerializableInterface 表示 PHP 的 Serializable 接口
type SerializableInterface struct {
	*node.InterfaceStatement
}

// NewSerializableInterface 创建一个新的 Serializable 接口
func NewSerializableInterface() *SerializableInterface {
	return &SerializableInterface{
		InterfaceStatement: node.NewInterfaceStatement(
			nil,
			"Serializable",
			nil,
			[]data.Method{
				&SerializableSerializeMethod{},
				&SerializableUnserializeMethod{},
			},
		),
	}
}

// SerializableSerializeMethod 表示 Serializable::serialize(): string|null
type SerializableSerializeMethod struct{}

func (m *SerializableSerializeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func (m *SerializableSerializeMethod) GetName() string            { return "serialize" }
func (m *SerializableSerializeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SerializableSerializeMethod) GetIsStatic() bool          { return false }

func (m *SerializableSerializeMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *SerializableSerializeMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *SerializableSerializeMethod) GetReturnType() data.Types {
	return data.NewNullableType(data.NewBaseType("string"))
}

// SerializableUnserializeMethod 表示 Serializable::unserialize(string $data): void
type SerializableUnserializeMethod struct{}

func (m *SerializableUnserializeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewNullValue(), nil
}

func (m *SerializableUnserializeMethod) GetName() string            { return "unserialize" }
func (m *SerializableUnserializeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SerializableUnserializeMethod) GetIsStatic() bool          { return false }

func (m *SerializableUnserializeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, nil),
	}
}

func (m *SerializableUnserializeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.NewBaseType("string")),
	}
}

func (m *SerializableUnserializeMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
